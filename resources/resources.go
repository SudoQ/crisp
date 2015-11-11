package resources

import (
	"fmt"
	"github.com/SudoQ/crisp/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Manager struct {
	cache *storage.Store
}

func New(store *storage.Store) *Manager {
	return &Manager{
		cache: store,
	}
}

func (this *Manager) Run(rawport string) {
	r := mux.NewRouter()
	r.HandleFunc("/", this.HomeHandler)
	r.HandleFunc("/info", this.InfoHandler)
	r.HandleFunc("/cache.json", this.CacheHandler)
	port := fmt.Sprintf(":%s", rawport)
	http.ListenAndServe(port, r)
}

func (this *Manager) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	latestItem, err := this.cache.Get()
	if err != nil {
		log.Print(err)
		w.WriteHeader(404)
		return
	}
	w.Write(latestItem.Payload)
}

func (this *Manager) Info() string {
	return fmt.Sprintf("Crisp API caching service v0.1")
}

func (this *Manager) InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(this.Info()))
}

func (this *Manager) CacheHandler(w http.ResponseWriter, r *http.Request) {
	latestItem, err := this.cache.Get()
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}

	response, err := latestItem.JSON()
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}
