package dto

import (
	"encoding/gob"
	"time"
)

type SensorMessage struct {
	Name string
	Url  string
	Js   bool
	//Value     float64
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage{})
}
