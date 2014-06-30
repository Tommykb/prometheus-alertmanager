// Copyright 2013 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/http/pprof"

	"code.google.com/p/gorest"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/alertmanager/web/api"
	"github.com/prometheus/alertmanager/web/blob"
)

// Commandline flags.
var (
	listenAddress  = flag.String("listenAddress", ":9090", "Address to listen on for web interface.")
	useLocalAssets = flag.Bool("useLocalAssets", false, "Read assets/templates from file instead of binary.")
)

type WebService struct {
	AlertManagerService *api.AlertManagerService
	AlertsHandler       *AlertsHandler
	SilencesHandler     *SilencesHandler
	StatusHandler       *StatusHandler
}

func (w WebService) ServeForever() error {
	gorest.RegisterService(w.AlertManagerService)

	http.Handle("/favicon.ico", prometheus.InstrumentHandlerFunc("favicon", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", 404)
	}))

	http.Handle("/debug/pprof/", prometheus.InstrumentHandlerFunc("debug", pprof.Index))
	http.Handle("/debug/pprof/cmdline", prometheus.InstrumentHandlerFunc("debug", pprof.Cmdline))
	http.Handle("/debug/pprof/profile", prometheus.InstrumentHandlerFunc("debug", pprof.Profile))
	http.Handle("/debug/pprof/symbol", prometheus.InstrumentHandlerFunc("debug", pprof.Symbol))

	http.Handle("/", prometheus.InstrumentHandler("root", w.AlertsHandler))
	http.Handle("/alerts", prometheus.InstrumentHandler("alerts", w.AlertsHandler))
	http.Handle("/silences", prometheus.InstrumentHandler("silences", w.SilencesHandler))
	http.Handle("/status", prometheus.InstrumentHandler("status", w.StatusHandler))

	http.Handle("/api/", prometheus.InstrumentHandler("api", compressionHandler{handler: gorest.Handle()}))
	http.Handle("/metrics", prometheus.Handler())
	if *useLocalAssets {
		http.Handle("/static/", prometheus.InstrumentHandler("static", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static")))))
	} else {
		http.Handle("/static/", prometheus.InstrumentHandler("static", http.StripPrefix("/static/", new(blob.Handler))))
	}

	glog.Info("listening on ", *listenAddress)

	return http.ListenAndServe(*listenAddress, nil)
}

func getLocalTemplate(name string) (*template.Template, error) {
	t := template.New("_base.html")
	t.Funcs(webHelpers)
	return t.ParseFiles(
		"web/templates/_base.html",
		fmt.Sprintf("web/templates/%s.html", name),
	)
}

func getEmbeddedTemplate(name string) (*template.Template, error) {
	t := template.New("_base.html")
	t.Funcs(webHelpers)

	file, err := blob.GetFile(blob.TemplateFiles, "_base.html")
	if err != nil {
		glog.Error("Could not read base template: ", err)
		return nil, err
	}
	t.Parse(string(file))

	file, err = blob.GetFile(blob.TemplateFiles, name+".html")
	if err != nil {
		glog.Errorf("Could not read %s template: %s", name, err)
		return nil, err
	}
	t.Parse(string(file))

	return t, nil
}

func getTemplate(name string) (t *template.Template, err error) {
	if *useLocalAssets {
		t, err = getLocalTemplate(name)
	} else {
		t, err = getEmbeddedTemplate(name)
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}

func executeTemplate(w http.ResponseWriter, name string, data interface{}) {
	tpl, err := getTemplate(name)
	if err != nil {
		glog.Error("Error preparing layout template: ", err)
		return
	}
	err = tpl.Execute(w, data)
	if err != nil {
		glog.Error("Error executing template: ", err)
	}
}
