package lp

import (
	"fmt"
	"os"

	log "github.com/rjbrown57/lp/pkg/logging"
)

// LpConfig set config for site hosting
type LpConfig struct {
	Lpconfig struct {
		RootDir  string `yaml:"rootDir"`
		Port     int    `yaml:"port"`
		Sitename string `yaml:"sitename"`
		Tls      struct {
			Key  string `yaml:"key"`
			Cert string `yaml:"cert"`
			Ca   string `yaml:"ca"`
		} `yaml:"tls"`
	} `yaml:"lpconfig"`
	sitedata SiteData
}

// write pages
func (lp LpConfig) writePages() {

	var fileName string

	for _, page := range lp.sitedata.Template.Pages {
		if page.IsIndex {
			fileName = fmt.Sprintf(lp.Lpconfig.RootDir + "/index.html")
		} else {
			fileName = fmt.Sprintf(lp.Lpconfig.RootDir + "/" + page.Name + ".html")
		}
		log.Infof("Creating %s\n", fileName)

		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open %s #%v\n", page.Name, err)
		}

		err = os.Truncate(fileName, 0)
		if err != nil {
			log.Fatalf("Unable to truncate file %s", fileName)
		}

		// render common
		t := processTemplate("common", commonTemplate)
		err = t.Execute(file, lp.sitedata)
		if err != nil {
			log.Fatalf("common template render error #%v\n", err)
		}

		// render navbar
		t = processTemplate("navbar", navbarTemplate)
		err = t.Execute(file, lp.sitedata)
		if err != nil {
			log.Fatalf("navbar template render error #%v\n", err)
		}

		// render body
		t = processTemplate("body", bodyTemplate)
		err = t.Execute(file, page)
		if err != nil {
			log.Fatalf("body template render error #%v\n", err)
		}

	}
}

// SiteData is the meat and potatoes of LP. Used to create html page of links.
type SiteData struct {
	Template struct {
		Theme string `yaml:"Theme,omitempty"`
		Pages []struct {
			IsIndex  bool   `yaml:"IsIndex,omitempty"`
			Name     string `yaml:"Name"`
			Headings []struct {
				Name  string `yaml:"Name"`
				Links []struct {
					Name string              `yaml:"Name"`
					Url  string              `yaml:"Url,omitempty"`
					Urls []map[string]string `yaml:"Urls,omitempty"`
				} `yaml:"Links"`
			} `yaml:"Headings"`
		} `yaml:"Pages"`
	} `yaml:"Template"`
}
