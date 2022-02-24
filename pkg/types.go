package lp

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
}

// SiteData is the meat and potatoes of LP. Used to create html page of links.
type SiteData struct {
	Template struct {
		Theme string `yaml:"Theme,omitempty"`
		Pages []struct {
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
