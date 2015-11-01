package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SudoQ/ganache/item"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Service struct {
	URL    string        `json:"url"`
	Port   string        `json:"port"`
	Period time.Duration `json:"-"`
	Limit  uint          `json:"limit"`
	Cache  *item.Item    `json:"-"`
}

func New(url, port string, limit uint) *Service {
	period, err := LimitToDuration(limit)
	if err != nil {
		limit = 1
		period, _ = LimitToDuration(limit)
	}
	return &Service{
		URL:    url,
		Port:   port,
		Period: period,
		Limit:  limit,
		Cache:  nil,
	}
}

func LimitToDuration(limit uint) (time.Duration, error) {
	if limit == 0 {
		return 0, errors.New("Division with zero")
	}
	return time.Duration(time.Duration(60/limit) * time.Minute), nil
}

func NewFromJSON(jsonBlob []byte) (*Service, error) {
	var srv Service
	err := json.Unmarshal(jsonBlob, &srv)
	if err != nil {
		return nil, err
	}
	srv.Period, err = LimitToDuration(srv.Limit)
	if err != nil {
		srv.Limit = 1
		srv.Period, _ = LimitToDuration(srv.Limit)
	}
	return &srv, nil
}

func NewFromFile(filename string) (*Service, error) {
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewFromJSON(jsonBlob)
}

func (this *Service) LoadConfig(filename string) error {
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBlob, this)
	return err
}

func (this *Service) JSON() ([]byte, error) {
	// Return service configuration as JSON
	blob, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	return blob, nil
}

func (this *Service) Chaos(err error) bool {
	if err != nil {
		// TODO Log error
		log.Println(err)
		return true
	}
	return false
}

func (this *Service) SaveConfig(filename string) error {
	jsonBlob, err := this.JSON()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, jsonBlob, 0644)
	return err
}

func (this *Service) Collect() {
	for {
		now := time.Now()
		timeDelta := now.Sub(this.Cache.Timestamp)
		if this.Cache.Timestamp.Before(now) && timeDelta > this.Period {
			log.Printf("GET %s\n", this.URL)

			resp, err := http.Get(this.URL)
			if this.Chaos(err) {
				// TODO handle error
				log.Fatal(err)
			}

			payload, err := ioutil.ReadAll(resp.Body)
			if this.Chaos(err) {
				// TODO handle error
				log.Fatal(err)
			}

			// TODO Store header info?
			newItem := item.New(now, payload)
			this.Cache = newItem
			cacheFilename := "cache.json"
			err = this.Cache.WriteFile(cacheFilename)
			if this.Chaos(err) {
				log.Fatal(err)
			}
			log.Printf("Saved cache to %s\n", cacheFilename)

		}
		select {
		case <-time.After(this.Period):
			continue
		}
	}
}

func (this *Service) LoadCache(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	newItem, err := item.NewFromJSON(content)
	if err != nil {
		return err
	}
	this.Cache = newItem
	return nil
}

func (this *Service) Info() string {
	return fmt.Sprintf("Ganache API caching service v0.1")
}

func (this *Service) Run() {
	err := this.LoadCache("cache.json")
	if err != nil {
		log.Fatal(err)
	}
	go this.Collect()
	r := mux.NewRouter()
	r.HandleFunc("/", this.HomeHandler)
	r.HandleFunc("/info", this.InfoHandler)
	r.HandleFunc("/cache.json", this.CacheHandler)
	port := fmt.Sprintf(":%s", this.Port)
	http.ListenAndServe(port, r)
}

func (this *Service) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(this.Cache.Payload)
}

func (this *Service) InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(this.Info()))
}

func (this *Service) CacheHandler(w http.ResponseWriter, r *http.Request) {
	response, err := this.Cache.JSON()
	if err != nil {
		// TODO
		log.Fatal(err)
	}
	w.Write(response)
}
