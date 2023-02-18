package lp

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/rjbrown57/lp/pkg/logging"
	"gopkg.in/yaml.v2"
)

var cssUrl string = "https://unpkg.com/@primer/css@^20.2.4/dist/primer.css"

func mustUnmarshalYaml(configPath string, v interface{}) {
	log.Infof("Reading %s\n", configPath)
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("err opening %s   #%v\n", configPath, err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, v)
	if err != nil {
		log.Fatalf("unmarhsal error   #%v\n", err)
		os.Exit(1)
	}
}

// processTemplate returns the template object
func processTemplate(templateName string, goTemplate string) *template.Template {
	return template.Must(template.New(templateName).Parse(goTemplate))
}

// serveLP
func serveLP(htmlDir string, port string) {
	log.Fatalf("Error serving page %s", http.ListenAndServe(port, http.FileServer(http.Dir(htmlDir))))
}

func getWatcher(siteTemplate []string) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Unable to start fsnotify watcher %s", err)
	}

	for _, site := range siteTemplate {
		log.Infof("Watching %s for changes\n", site)
		watcher.Add(site)
	}

	return watcher
}

func DownloadCss(path string) error {
	log.Infof("Downloading %s", cssUrl)
	resp, err := http.Get(cssUrl)
	if err != nil || resp.StatusCode > 200 {
		return fmt.Errorf("failed to download from %s - %v , %d", cssUrl, err, resp.StatusCode)
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath.Clean(path))
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(io.MultiWriter(out), resp.Body)
	if err != nil {
		return err
	}

	log.Debugf("Download %s complete", cssUrl)
	return nil
}

func monitorChanges(w *fsnotify.Watcher, siteTemplate []string, lp LpConfig) {
	for {
		select {
		// Read Errors
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Fatalf("Fsnotify ERROR: %s", err)
		// Read events and writePages if appropriate
		case event, ok := <-w.Events:
			if !ok {
				log.Fatalf("Channel closed %s", event.Name)
			}

			switch event.Op {
			// Ignore chmod
			case fsnotify.Chmod, fsnotify.Rename:
				continue
			// On create/remove we need to wait a second for the dust to settle. Otherwise we might read a a file that does not exist
			// We need to add the new file to be monitored
			case fsnotify.Create, fsnotify.Remove:
				time.Sleep(1 * time.Second)
				w.Add(event.Name)
				fallthrough
			default:
				log.Infof("Change detected %s %s. Regenerating...", event.Op.String(), event.Name)
				lp.sitedata = mergePages(siteTemplate)
				lp.writePages()
			}

		}
	}
}

// mergePages will allow user to supply multiple page templates and merge them together
func mergePages(siteTemplate []string) SiteData {

	mergedSiteData := &SiteData{}

	var sites []SiteData

	for _, st := range siteTemplate {
		siteData := &SiteData{}
		mustUnmarshalYaml(st, siteData)
		sites = append(sites, *siteData)
	}

	//initialize mergedSiteData with first element
	mergedSiteData.Template = sites[0].Template

	// Skip the first element since we take it above
	for _, st := range sites[1:] {
		mergedSiteData.Template.Pages = append(mergedSiteData.Template.Pages, st.Template.Pages...)
	}

	return *mergedSiteData
}

// Lp calls mustUnmarshalYaml for configs, writePages to write appropriate files, serveLP to host
func Lp(action string, follow bool, lpconfig string, siteTemplate []string) {

	// Create basic config object and
	// Merge all site templates supplied by user
	config := &LpConfig{
		sitedata: mergePages(siteTemplate),
	}

	mustUnmarshalYaml(lpconfig, config)

	var err error

	if config.Lpconfig.RootDir == "" {
		config.Lpconfig.RootDir, err = os.MkdirTemp("/tmp/", "lp")
		if err != nil {
			log.Fatalf("unable to create temporary directory")
		}

		defer os.RemoveAll(config.Lpconfig.RootDir)
	}

	log.Infof("Using %s as html root\n", config.Lpconfig.RootDir)

	// Grab css
	if err := DownloadCss(config.Lpconfig.RootDir + "/primer.css"); err != nil {
		log.Fatalf("Unable to download Css - %s", err)
	}

	// Generate html pages
	config.writePages()

	// If we are called by generate, return without serving page
	switch action {
	case "generate":
		if !follow {
			return
		}
	case "serve":
		log.Infof("Serving LP on :%d\n", config.Lpconfig.Port)
		go serveLP(config.Lpconfig.RootDir, fmt.Sprintf(":%d", config.Lpconfig.Port))
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	w := getWatcher(siteTemplate)
	defer w.Close()

	go monitorChanges(w, siteTemplate, *config)

	go func() {
		sig := <-sigs
		log.Infof("Received %s. Cleaning up %s\n", sig, config.Lpconfig.RootDir)
		done <- true
	}()

	<-done
}
