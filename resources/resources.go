package resources

import (
	"fmt"
	"github.com/SudoQ/crisp/logging"
	"github.com/SudoQ/crisp/storage"
	"github.com/gorilla/mux"
	"net/http"
)

type Manager struct {
	cache *storage.Store
	port  string
}

func New(store *storage.Store, port string) *Manager {
	return &Manager{
		cache: store,
		port:  "8080",
	}
}

func (this *Manager) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/", this.HomeHandler)
	r.HandleFunc("/info", this.InfoHandler)
	r.HandleFunc("/latest.json", this.CacheHandler)
	port := fmt.Sprintf(":%s", this.port)
	http.ListenAndServe(port, r)
}

func (this *Manager) logAccess(r *http.Request) {
	logging.Info(fmt.Sprintf("%s accessed by %s", r.URL, r.RemoteAddr))
}

func (this *Manager) Info() string {
	return fmt.Sprintf("Crisp Service v0.1")
}

func (this *Manager) HomeHandler(w http.ResponseWriter, r *http.Request) {
	this.logAccess(r)
	latestItem, err := this.cache.Get()
	if err != nil {
		logging.Error(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(latestItem.Payload)
}

func (this *Manager) InfoHandler(w http.ResponseWriter, r *http.Request) {
	this.logAccess(r)
	w.WriteHeader(200)
	w.Write([]byte(this.Info()))
}

func (this *Manager) CacheHandler(w http.ResponseWriter, r *http.Request) {
	this.logAccess(r)
	latestItem, err := this.cache.Get()
	if err != nil {
		logging.Error(err)
		w.WriteHeader(500)
		return
	}

	response, err := latestItem.JSON()
	if err != nil {
		logging.Error(err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}
