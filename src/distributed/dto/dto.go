package dto

import (
	"encoding/gob"
)

type SensorMessage struct {
	Url    string
	Js     bool
	Header bool
	Html   bool
}

func init() {
	gob.Register(SensorMessage{})
}
