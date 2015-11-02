package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SudoQ/crisp/item"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Service struct {
	URL    string        `json:"url"`
	Port   string        `json:"port"`
	Period time.Duration `json:"-"`
	Limit  uint          `json:"limit"`
	Subscribers []string `json:"subscribers"`
	Cache  *item.Item    `json:"-"`
	logger *log.Logger   `json:"-"`
	cacheCh chan *item.Item
}

func New(targetUrl, port string, limit uint) *Service {
	period, err := LimitToPeriod(limit)
	if err != nil {
		limit = 1
		period, _ = LimitToPeriod(limit)
	}
	srv := &Service{
		URL:    targetUrl,
		Port:   port,
		Period: period,
		Limit:  limit,
		Subscribers: make([]string, 0),
		Cache:  nil,
		logger: nil,
		cacheCh: make(chan *item.Item),
	}
	srv.initLogger()
	return srv
}

func (this *Service) initLogger() {
	u, err := url.Parse(this.URL)
	label := u.Host
	if err != nil {
		label = "?"
	}
	this.logger = log.New(os.Stdin, fmt.Sprintf("crisp[%s]: ", label), log.Lshortfile)
}

func LimitToPeriod(limit uint) (time.Duration, error) {
	if limit == 0 {
		return 0, errors.New("Division with zero")
	}

	period := (60.0 / float64(limit)) * 60
	return (time.Duration(period) * time.Second), nil
}

func NewFromJSON(jsonBlob []byte) (*Service, error) {
	var srv Service
	err := json.Unmarshal(jsonBlob, &srv)
	if err != nil {
		return nil, err
	}
	srv.Period, err = LimitToPeriod(srv.Limit)
	if err != nil {
		srv.Limit = 1
		srv.Period, _ = LimitToPeriod(srv.Limit)
	}
	srv.cacheCh = make(chan *item.Item)
	srv.initLogger()
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
	this.initLogger()
	return err
}

func (this *Service) JSON() ([]byte, error) {
	blob, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	return blob, nil
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
		func() {
			now := time.Now()
			timeDelta := now.Sub(this.Cache.Timestamp)
			if this.Cache.Timestamp.Before(now) && timeDelta > this.Period {
				this.logger.Printf("GET %s\n", this.URL)

				resp, err := http.Get(this.URL)
				if err != nil {
					this.logger.Println(err)
					return
				}

				payload, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					this.logger.Println(err)
					return
				}

				newItem := item.New(now, payload)
				this.Cache = newItem
				cacheFilename := "cache.json"
				err = this.Cache.WriteFile(cacheFilename)
				if err != nil {
					this.logger.Fatal(err)
				}
				this.logger.Printf("Saved cache to %s\n", cacheFilename)
				this.cacheCh <- this.Cache
			}
		}()
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

func (this *Service) Subscribe(subUrl string) {
	// Subscribers provide thier url where the service should send the response
	// This method adds the subscriber to the subscriber list
	// TODO Lock here
	this.Subscribers = append(this.Subscribers, subUrl)
}

func (this *Service) Publish() {
	// Listen for new event on a channel from collect
	// Send out to subscribers
	for {
		select {
		case data := <-this.cacheCh:
			// Send to subscribers
			// TODO Lock here
			for _, sub := range this.Subscribers {
				log.Printf("POST data %v to url %s", data, sub)
			}
		}
	}
}

func (this *Service) Info() string {
	return fmt.Sprintf("Crisp API caching service v0.1")
}

func (this *Service) Run() {
	err := this.LoadCache("cache.json")
	if err != nil {
		this.logger.Fatal(err)
	}
	go this.Collect()
	r := mux.NewRouter()
	r.HandleFunc("/", this.HomeHandler)
	r.HandleFunc("/info", this.InfoHandler)
	r.HandleFunc("/cache.json", this.CacheHandler)
	r.HandleFunc("/subscribe", this.SubscriptionHandler)
	port := fmt.Sprintf(":%s", this.Port)
	http.ListenAndServe(port, r)
}

func (this *Service) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(this.Cache.Payload)
}

func (this *Service) InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(this.Info()))
}

func (this *Service) CacheHandler(w http.ResponseWriter, r *http.Request) {
	response, err := this.Cache.JSON()
	if err != nil {
		this.logger.Println(err)
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

func (this *Service) SubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}
