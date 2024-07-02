package redismq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Package provides headers and handling functions around payloads
type Package struct {
	Payload    string
	CreatedAt  time.Time
	Queue      interface{} `json:"-"`
	Consumer   *Consumer   `json:"-"`
	Collection *[]*Package `json:"-"`
	Acked      bool        `json:"-"`
	//TODO add Headers or smth. when needed
	//wellle suggested error headers for failed packages
}

func unmarshalPackage(input string, queue *Queue, consumer *Consumer) (*Package, error) {
	p := &Package{Queue: queue, Consumer: consumer, Acked: false}
	err := json.Unmarshal([]byte(input), p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (pack *Package) getString() string {
	json, err := json.Marshal(pack)
	if err != nil {
		log.Printf(" Queue failed to marshal content %v [%s]", pack, err.Error())
		// TODO build sensible error handling
		return ""
	}
	return string(json)
}

func (pack *Package) index() int {
	if pack.Collection == nil {
		return 0
	}
	var i int
	for i = range *pack.Collection {
		if (*pack.Collection)[i] == pack {
			break
		}
	}
	return i
}

// Ack removes the packages from the queue
func (pack *Package) Ack() error {
	if pack.Collection != nil {
		return fmt.Errorf("cannot Ack package in multi package answer")
	}
	err := pack.Consumer.ackPackage()
	return err
}

// Requeue moves a package back to input
func (pack *Package) Requeue() error {
	return pack.reject(true)
}

// Fail moves a package to the failed queue
func (pack *Package) Fail() error {
	return pack.reject(false)
}

func (pack *Package) reject(requeue bool) error {
	if pack.Collection != nil && !(*pack.Collection)[pack.index()-1].Acked {
		return fmt.Errorf("cannot reject package while unacked package before it")
	}

	if !requeue {
		err := pack.Consumer.failPackage()
		return err
	}
	err := pack.Consumer.requeuePackage()
	return err
}
