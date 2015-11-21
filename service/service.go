package service

import (
	"errors"
	"github.com/SudoQ/satchel/external"
	"github.com/SudoQ/satchel/item"
	"github.com/SudoQ/satchel/resources"
	"github.com/SudoQ/satchel/storage"
	"io/ioutil"
	"time"
)

type Service struct {
	URL             string
	Port            string
	Period          time.Duration
	Limit           uint
	Cache           *storage.Store
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
		ext:             external.New(target, period),
		resourceManager: nil,
	}
	srv.initResouceManager()
	return srv
}

func (this *Service) initResouceManager() {
	this.resourceManager = resources.New(this.Cache, this.Port)
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

func (this *Service) Run() {
	dataCh := this.ext.DataChannel()
	defer this.ext.Close()
	go this.ext.Collect()
	go func() {
		for {
			select {
			case payload := <-dataCh:
				newItem := item.New(time.Now(), payload)
				this.Cache.Add(newItem)
			}
		}
	}()

	this.resourceManager.Run()
}

func LimitToPeriod(limit uint) (time.Duration, error) {
	if limit == 0 {
		return 0, errors.New("Division with zero")
	}

	period := (60.0 / float64(limit)) * 60
	return (time.Duration(period) * time.Second), nil
}
