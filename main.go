package main

import (
	"github.com/polyverse/appconfig"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// Specify the arguments
	params := make(map[string]appconfig.Param)
	params["bind"] = appconfig.Param{Type: appconfig.PARAM_STRING, Default: "80", Usage: "bind-to port.", Required: true}
	params["port"] = appconfig.Param{Type: appconfig.PARAM_STRING, Usage: "redirect-to port.", Required: false}
	params["scheme"] = appconfig.Param{Type: appconfig.PARAM_STRING, Default: "https", Usage: "redirect-to http or https.", Required: true}
	params["host"] = appconfig.Param{Type: appconfig.PARAM_STRING, Usage: "redirect-to host.", Required: false}
	params["debug"] = appconfig.Param{Type: appconfig.PARAM_BOOL, Default: false, Usage: "verbose output.", PrefixOverride: "--"}
	params["help"] = appconfig.Param{Type: appconfig.PARAM_USAGE, Default: false, Usage: "print usage.", Required: false, PrefixOverride: "--"}
	config, err := appconfig.NewConfig(params)
	if config.GetBool("help") {
		config.PrintUsage("Listen to a http or https port and redirect (\"HTTP/1.1 301 Moved Permanently\") to '<scheme>://<host>' while preserving the path and query string.\n\n")
		os.Exit(0)
	}

	redirectPort := config.GetString("port")
	redirectScheme := config.GetString("scheme")
	redirectHost := config.GetString("host")
	listenerBindPort := config.GetString("bind")
	debug := config.GetBool("debug")

	// "redirect starting..." to stdout
	log.WithFields(log.Fields{"port": redirectPort, "scheme": redirectScheme, "host": redirectHost, "bind": listenerBindPort, "debug": debug}).Infof("Starting redirect...")

	if config.GetBool("debug") == true {
		log.SetLevel(log.DebugLevel)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "healthy")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			log.WithFields(log.Fields{"err": err}).Errorf("Error parsing url.")
		}
		u.Scheme = redirectScheme

		hostFragments := strings.Split(r.Host, ":")
		if len(hostFragments) > 0 && redirectHost != "" { //meaning we have the host part of the hostname, and we have a redirect host provided.
			hostFragments[0] = redirectHost
		}

		if len(hostFragments) > 1 && redirectPort != "" { //Meaning we have a port part of the hostname, and we have a redirect port provided.
			hostFragments[1] = redirectPort
		} else if redirectPort != "" {
			hostFragments = append(hostFragments, redirectPort)
		}

		//join the fragments to give us the new hostname
		u.Host = strings.Join(hostFragments, ":")

		// find out who's making the request
		if debug {
			if r.URL.Scheme == "" {
				r.URL.Scheme = "http"
			}
			r.URL.Host = r.Host
			log.Debugf("Redirecting '%s' to '%s'", r.URL.String(), u.String())
		}

		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	})

	if listenerBindPort[0:1] != ":" {
		listenerBindPort = ":" + listenerBindPort
	}
	err = http.ListenAndServe(listenerBindPort, nil)
	if err != nil {
		log.WithFields(log.Fields{"port": redirectPort, "scheme": redirectScheme, "host": redirectHost, "bind": listenerBindPort, "debug": debug}).Fatal("ListenAndServe: ", err)
	}
}
