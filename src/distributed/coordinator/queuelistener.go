package coordinator

import (
	"bytes"
	//"distributed/dto"
	//"distributed/qutils"
	"encoding/gob"
	"fmt"

	"github.com/gopheramit/distributed-go-with-rabbitmq/src/distributed/dto"
	"github.com/gopheramit/distributed-go-with-rabbitmq/src/distributed/qutils"
	"github.com/streadway/amqp"
)

const url = "amqp://guest:guest@localhost:5672"

type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
	ea      *EventAggregator
}

func NewQueueListener() *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
		ea:      NewEventAggregator(),
	}

	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

func (ql *QueueListener) DiscoverSensors() {
	ql.ch.ExchangeDeclare(
		qutils.SensorDiscoveryExchange, //name string,
		"fanout",                       //kind string,
		false,                          //durable bool,
		false,                          //autoDelete bool,
		false,                          //internal bool,
		false,                          //noWait bool,
		nil)                            //args amqp.Table)

	ql.ch.Publish(
		qutils.SensorDiscoveryExchange, //exchange string,
		"",                             //key string,
		false,                          //mandatory bool,
		false,                          //immediate bool,
		amqp.Publishing{})              //msg amqp.Publishing)
}

func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch)
	ql.ch.QueueBind(
		q.Name,       //name string,
		"",           //key string,
		"amq.fanout", //exchange string,
		false,        //noWait bool,
		nil)          //args amqp.Table)

	msgs, _ := ql.ch.Consume(
		q.Name, //queue string,
		"",     //consumer string,
		true,   //autoAck bool,
		false,  //exclusive bool,
		false,  //noLocal bool,
		false,  //noWait bool,
		nil)    //args amqp.Table)

	ql.DiscoverSensors()

	fmt.Println("listening for new sources")

	// updated the if guard below to surround all
	// of the for-loops contents to prevent
	// same sensor being registered multiple
	// times with RabbitMQ
	for msg := range msgs {
		if ql.sources[string(msg.Body)] == nil {
			fmt.Println("new source discovered")
			sourceChan, _ := ql.ch.Consume(
				string(msg.Body), //queue string,
				"",               //consumer string,
				true,             //autoAck bool,
				false,            //exclusive bool,
				false,            //noLocal bool,
				false,            //noWait bool,
				nil)              //args amqp.Table)

			ql.sources[string(msg.Body)] = sourceChan

			go ql.AddListener(sourceChan)
		}
	}
}

func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)

		fmt.Printf("Received message: %v\n", sd)

		ed := EventData{
			Name:      sd.Name,
			Timestamp: sd.Timestamp,
			Js:        sd.Js,
		}

		ql.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)
	}
}
