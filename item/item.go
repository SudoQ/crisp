package item

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Item struct {
	Timestamp time.Time `json:timestamp`
	Payload   []byte    `json:payload`
}

func New(timestamp time.Time, payload []byte) *Item {
	return &Item{
		Timestamp: timestamp,
		Payload:   payload,
	}
}

func NewFromJSON(jsonBlob []byte) (*Item, error) {
	var item Item
	err := json.Unmarshal(jsonBlob, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (this *Item) JSON() ([]byte, error) {
	blob, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	return blob, nil
}

func (this *Item) WriteFile(filename string) error {
	blob, err := this.JSON()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, blob, 0644)
	if err != nil {
		return err
	}
	return nil
}
