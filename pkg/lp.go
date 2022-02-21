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

// write pages
//https://yourbasic.org/golang/append-to-file/
// update to write generated tmpdir and then serve from that
func writePages(siteData *SiteData, tDir string) {
	for _, page := range siteData.Template.Pages {

		fileName := fmt.Sprintf(tDir + "/" + page.Name + ".html")
		log.Printf("Creating %s\n", fileName)

		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to open %s #%v\n", page.Name, err)
			os.Exit(1)
		}

		// render common
		t := template.Must(template.New("common").Parse(commonTemplate))
		err = t.Execute(file, siteData)
		if err != nil {
			log.Printf("common template render error #%v\n", err)
			os.Exit(1)
		}

		// render navbar
		t = template.Must(template.New("navbar").Parse(navbarTemplate))
		err = t.Execute(file, siteData)
		if err != nil {
			log.Printf("navbar template render error #%v\n", err)
			os.Exit(1)
		}

		// render body
		t = template.Must(template.New("body").Parse(bodyTemplate))
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

func Lp(lpconfig string, siteTempalte string) {
	config := &LpConfig{}
	siteData := &SiteData{}

	tDir, err := os.MkdirTemp("/tmp/", "lp")
	if err != nil {
		log.Println("unable to create temporary directory")
	}

	defer os.RemoveAll(tDir)

	mustUnmarshalYaml(lpconfig, config)
	mustUnmarshalYaml(siteTempalte, siteData)

	writePages(siteData, tDir)

	go serveLP(tDir, fmt.Sprintf(":%d", config.Lpconfig.Port))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		log.Printf("Recieved %s. Cleaning up %s\n", sig, tDir)
		done <- true
	}()

	<-done

}
