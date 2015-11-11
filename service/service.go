package service

import (
	"errors"
	"fmt"
	"github.com/SudoQ/crisp/external"
	"github.com/SudoQ/crisp/item"
	"github.com/SudoQ/crisp/resources"
	"github.com/SudoQ/crisp/storage"
	//"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	//"net/http"
	"net/url"
	"os"
	"time"
)

type Service struct {
	URL             string
	Port            string
	Period          time.Duration
	Limit           uint
	Cache           *storage.Store
	logger          *log.Logger
	ext             *external.Ext
	resourceManager *resources.Manager
}

func New(target, port string, limit uint) *Service {
	period, err := LimitToPeriod(limit)
	if err != nil {
		limit = 1
		period, _ = LimitToPeriod(limit)
	}
	srv := &Service{
		URL:             target,
		Port:            port,
		Period:          period,
		Limit:           limit,
		Cache:           storage.New(),
		logger:          nil,
		ext:             external.New(target, period),
		resourceManager: nil,
	}
	srv.initLogger()
	srv.initResouceManager()
	return srv
}

func (this *Service) initResouceManager() {
	this.resourceManager = resources.New(this.Cache)
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

func (this *Service) LoadCache(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	newItem, err := item.NewFromJSON(content)
	if err != nil {
		return err
	}
	this.Cache.Add(newItem)
	return nil
}

func (this *Service) Info() string {
	return fmt.Sprintf("Crisp API caching service v0.1")
}

func (this *Service) Run() {
	dataCh := this.ext.DataChannel()
	defer this.ext.Close()
	go this.ext.Collect()
	go func() {
		for payload := range dataCh {
			newItem := item.New(time.Now(), payload)
			this.Cache.Add(newItem)
			cacheFilename := "cache.json"
			latestItem, err := this.Cache.Get()
			if err != nil {
				this.logger.Fatal(err)
			}
			err = latestItem.WriteFile(cacheFilename)
			if err != nil {
				this.logger.Fatal(err)
			}
			this.logger.Printf("Saved cache to %s\n", cacheFilename)
		}
	}()

	this.resourceManager.Run(this.Port)
}
