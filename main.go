package main

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
	"github.com/saranrapjs/gift-helper/googleforms"
)

type Config struct {
	Port         string `default:"8080"`
	TemplateDir  string `default:"example"`
	TrackingPath string `default:"./tracking.json"`
	GoogleForm   string
}

type Registry struct {
	Data         map[string]RegistryRecord
	TrackingPath string
}

func (r *Registry) Read() {
	dat, _ := ioutil.ReadFile(r.TrackingPath)
	if dat != nil {
		json.Unmarshal(dat, &r.Data)
	}
}

func (r *Registry) Write() {
	data, err := json.Marshal(r.Data)
	if err == nil {
		err = ioutil.WriteFile(r.TrackingPath, data, 0644)
	}
}

func (r *Registry) Update(req *http.Request, form googleforms.Form) Selection {
	selection := Selection{}
	req.ParseForm()
	purchaseURL := req.Form.Get("purchase")
	logger.Println("attempting", purchaseURL)
	if rec, found := r.Data[purchaseURL]; purchaseURL != "" && found && !rec.Bought {
		newRec := RegistryRecord{
			Bought: true,
			URL:    purchaseURL,
			Name:   req.Form.Get("name"),
			Email:  req.Form.Get("email"),
			Notes:  req.Form.Get("notes"),
		}
		r.Data[purchaseURL] = newRec
		if err := form.Post([]string{newRec.Name, newRec.Notes, newRec.URL}...); err != nil {
			logger.Println(err)
			return selection
		}
		r.Write()
		selection.Made = true
		if parsed, err := url.Parse(purchaseURL); err == nil {
			selection.URL = parsed
		}
	}
	return selection
}

func NewRegistry(trackingPath string) *Registry {
	r := &Registry{
		Data:         map[string]RegistryRecord{},
		TrackingPath: trackingPath,
	}
	r.Read()
	return r
}

type RegistryRecord struct {
	Bought bool
	Name   string
	URL    string
	Email  string
	Notes  string
}

type Selection struct {
	Made bool
	URL  *url.URL
}

type TemplateData struct {
	Name      string
	Address   string
	Email     string
	Selection Selection
	Registry  *Registry
}

func (t *TemplateData) Bought(url string) template.HTMLAttr {
	bought := "false"
	rec, exists := t.Registry.Data[url]
	switch {
	case exists && rec.Bought:
		bought = "true"
	case !exists:
		t.Registry.Data[url] = RegistryRecord{}
		t.Registry.Write()
	}
	return template.HTMLAttr("data-bought=\"" + bought + "\" href=\"" + url + "\"")
}

var logger = log.New(os.Stderr, "", 0)

func RenderPage(w http.ResponseWriter, r *http.Request, c *Config, reg *Registry, form googleforms.Form) {
	selection := Selection{}
	if r.Method == "POST" {
		selection = reg.Update(r, form)
	}
	tData := &TemplateData{
		Selection: selection,
		Registry:  reg,
	}
	t, _ := template.ParseFiles(c.TemplateDir + "/template.html")
	err := t.Execute(w, tData)
	if err != nil {
		logger.Println(err)
	}
}

func main() {
	var c Config
	envconfig.Process("gift", &c)
	if c.GoogleForm == "" {
		logger.Println("missing Google form")
		return
	}
	form := googleforms.NewForm(c.GoogleForm)
	if err := form.Init(); err != nil {
		logger.Println("problem with Google form:" + err.Error())
		return
	}
	registry := NewRegistry(c.TrackingPath)
	r := http.NewServeMux()
	fs := http.FileServer(http.Dir(c.TemplateDir))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			fs.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == "/tracking.json" {
			http.Error(w, "no dice", http.StatusForbidden)
			return
		}
		RenderPage(w, r, &c, registry, form)
	})

	srv := &http.Server{
		Addr:    ":" + c.Port,
		Handler: handlers.CombinedLoggingHandler(os.Stderr, handlers.RecoveryHandler()(r)),
	}
	defer srv.Close()
	logger.Println("Listening on port " + c.Port)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	logger.Println("attempting to shut down on port " + c.Port)
	<-idleConnsClosed
	logger.Println("shutting down on port " + c.Port)
}
