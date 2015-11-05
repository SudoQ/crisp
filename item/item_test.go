package item

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestNewSuccess(t *testing.T) {
	now := time.Now()
	payload := []byte("payload")
	it := New(now, payload)
	if !it.Timestamp.Equal(now) {
		t.Fail()
	}
	if !(bytes.Equal(it.Payload, payload)) {
		t.Fail()
	}
}

func TestNewFail(t *testing.T) {
	now := time.Now()
	payload := []byte("payload")
	it := New(now, payload)
	if it.Timestamp.Equal(time.Time{}) {
		t.Fail()
	}
	if bytes.Equal(it.Payload, []byte("")) {
		t.Fail()
	}
}

func TestNewFromJSON_1(t *testing.T) {
	payload := []byte("payload")
	timeStr := "2015-11-01T13:32:36.748674638+01:00"
	var t0 time.Time
	_ = t0.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", timeStr)))

	jsonBlob := []byte(`{
		"Timestamp":"2015-11-01T13:32:36.748674638+01:00",
		"Payload":"cGF5bG9hZA=="}`)
	it0, _ := NewFromJSON(jsonBlob)
	if !it0.Timestamp.Equal(t0) {
		t.Fail()
	}
	if !bytes.Equal(it0.Payload, payload) {
		t.Fail()
	}
}

func TestNewFromJSON_2(t *testing.T) {
	now := time.Now()
	payload := []byte("payload")
	it0 := New(now, payload)
	jsonBlob, _ := it0.JSON()
	it1, _ := NewFromJSON(jsonBlob)

	if !it0.Timestamp.Equal(it1.Timestamp) {
		t.Fail()
	}
	if !bytes.Equal(it0.Payload, it1.Payload) {
		t.Fail()
	}
}
