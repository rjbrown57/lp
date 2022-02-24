# lp (link page)

Serve a webpage of useful links based on a yaml document

![intro](images/screenshot.png)

## How to use 

Clone the repo and run `docker-compose up -d`. This will use the example [site.yaml](config/site.yaml). Please customize [site.yaml](config/site.yaml) as you see fit!

You could also run with `docker run -i -t -p 8080:8080 -v "${PWD}/config:/config" ghcr.io/rjbrown57/lp:latest` to get the same result

Last option is to grab a release and simply run `./lp_0.0.1_version_here` locally

```
./lp --help
A yaml based static link page for every day work use.

Usage:
  lp [flags]

Flags:
  -h, --help                  help for lp
  -l, --lpConfig string       base config for lp see https://github.com/rjbrown57/lp/blob/main/config/lp.yaml (default "config/lp.yaml")
  -s, --siteTempalte string   site template see https://github.com/rjbrown57/lp/blob/main/config/site.yaml (default "config/site.yaml")
  -t, --toggle                Help message for toggle
```