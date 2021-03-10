package dto

import (
	"encoding/gob"
)

type SensorMessage struct {
	Name   string
	Url    string
	Js     bool
	Header bool
	Html   bool
}

func init() {
	gob.Register(SensorMessage{})
}
