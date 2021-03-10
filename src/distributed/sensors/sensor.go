package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"

	//"math/rand"

	//"src/distributed/dto"
	//	"src/distributed/qutils"

	"github.com/gopheramit/distributed-go-with-rabbitmq/src/distributed/dto"
	"github.com/gopheramit/distributed-go-with-rabbitmq/src/distributed/qutils"
	"github.com/streadway/amqp"
)

var url = "amqp://guest:guest@localhost:5672"

var name = flag.String("name", "sensor", "name of the sensor")

//var freq = flag.Uint("freq", 5, "update frequency in cycles/sec")
//var max = flag.Float64("max", 5., "maximum value for generated readings")
//var min = flag.Float64("min", 1., "minimum value for generated readings")
//var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

//var r = rand.New(rand.NewSource(time.Now().UnixNano()))

//var value = r.Float64()*(*max-*min) + *min
//var nom = (*max-*min)/2 + *min

func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)

	publishQueueName(ch)

	discoveryQueue := qutils.GetQueue("", ch)
	ch.QueueBind(
		discoveryQueue.Name,            //name string,
		"",                             //key string,
		qutils.SensorDiscoveryExchange, //exchange string,
		false,                          //noWait bool,
		nil)                            //args amqp.Table)

	go listenForDiscoverRequests(discoveryQueue.Name, ch)

	//dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	//signal := time.Tick(dur)

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	//
	//	for range signal {
	//	calcValue()
	reading := dto.SensorMessage{
		Name:   *name,
		Url:    url,
		Js:     false,
		Header: false,
		Html:   false,
	}
	buf.Reset()
	enc = gob.NewEncoder(buf)
	enc.Encode(reading)

	msg := amqp.Publishing{
		Body: buf.Bytes(),
	}

	ch.Publish(
		"",             //exchange string,
		dataQueue.Name, //key string,
		false,          //mandatory bool,
		false,          //immediate bool,
		msg)            //msg amqp.Publishing)

	log.Printf("Reading sent. Value: %v\n", msg)
}

func listenForDiscoverRequests(name string, ch *amqp.Channel) {
	msgs, _ := ch.Consume(
		name,  //queue string,
		"",    //consumer string,
		true,  //autoAck bool,
		false, //exclusive bool,
		false, //noLocal bool,
		false, //noWait bool,
		nil)   //args amqp.Table)

	for range msgs {
		log.Println("received discovery request")
		publishQueueName(ch)
	}
}

func publishQueueName(ch *amqp.Channel) {
	msg := amqp.Publishing{Body: []byte(*name)}
	ch.Publish(
		"amq.fanout", //exchange string,
		"",           //key string,
		false,        //mandatory bool,
		false,        //immediate bool,
		msg)          //msg amqp.Publishing)
}

/*
func client() {
	conn, ch, q := GetQueue()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, //queue string,
		"",     //consumer string,
		true,   //autoAck bool,
		false,  //exclusive bool,
		false,  //noLocal bool,
		false,  //noWait bool,
		nil)    //args amqp.Table)

	failOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		log.Printf("Received message with message: %s", msg.Body)
	}
}


func calcValue() {
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}
*/
