package lp

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

func mustUnmarshalYaml(configPath string, v interface{}) {
	log.Printf("Reading %s\n", configPath)
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("err opening %s   #%v\n", configPath, err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, v)
	if err != nil {
		log.Printf("unmarhsal error   #%v\n", err)
		os.Exit(1)
	}
}

// processTemplate returns the template object
func processTemplate(templateName string, goTemplate string) *template.Template {
	return template.Must(template.New(templateName).Parse(goTemplate))
}

// write pages
// https://yourbasic.org/golang/append-to-file/
// update to write generated tmpdir and then serve from that
func writePages(siteData *SiteData, tDir string) {

	var fileName string

	for _, page := range siteData.Template.Pages {
		if page.IsIndex {
			fileName = fmt.Sprintf(tDir + "/index.html")
		} else {
			fileName = fmt.Sprintf(tDir + "/" + page.Name + ".html")
		}
		log.Printf("Creating %s\n", fileName)

		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to open %s #%v\n", page.Name, err)
			os.Exit(1)
		}

		err = os.Truncate(fileName, 0)
		if err != nil {
			log.Printf("Unable to truncate file %s", fileName)
		}

		// render common
		t := processTemplate("common", commonTemplate)
		err = t.Execute(file, siteData)
		if err != nil {
			log.Printf("common template render error #%v\n", err)
			os.Exit(1)
		}

		// render navbar
		t = processTemplate("navbar", navbarTemplate)
		err = t.Execute(file, siteData)
		if err != nil {
			log.Printf("navbar template render error #%v\n", err)
			os.Exit(1)
		}

		// render body
		t = processTemplate("body", bodyTemplate)
		err = t.Execute(file, page)
		if err != nil {
			log.Printf("body template render error #%v\n", err)
			os.Exit(1)
		}

	}
}

// serveLP
func serveLP(htmlDir string, port string) {
	log.Fatal(http.ListenAndServe(port, http.FileServer(http.Dir(htmlDir))))
}

func getWatcher(siteTemplate []string) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Unable to start fsnotify watcher %s", err)
	}

	for _, site := range siteTemplate {
		log.Printf("Watching %s for changes\n", site)
		watcher.Add(site)
	}

	return watcher
}

func monitorChanges(w *fsnotify.Watcher, siteTemplate []string, lp LpConfig) {
	for {
		select {
		// Read Errors
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			log.Printf("Fsnotify ERROR: %s", err)
		// No matter what the return event is we should rerun writePages
		case _, _ = <-w.Events:
			// This should be tuned more, but in some cases we are too quick and the file is not finalized yet
			// This seems to specifically happen when working with git functions
			time.Sleep(1 * time.Second)
			lp.sitedata = mergePages(siteTemplate)
			lp.writePages()
			log.Printf("Change detected. Regenerating...")
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
			log.Println("unable to create temporary directory")
		}

		defer os.RemoveAll(config.Lpconfig.RootDir)
	}

	log.Printf("Using %s as html root\n", config.Lpconfig.RootDir)

	// Generate html pages
	config.writePages()

	// If we are called by generate, return without serving page
	switch action {
	case "generate":
		if !follow {
			return
		}
	case "serve":
		log.Printf("Serving LP on :%d\n", config.Lpconfig.Port)
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
		log.Printf("Received %s. Cleaning up %s\n", sig, config.Lpconfig.RootDir)
		done <- true
	}()

	<-done
}
