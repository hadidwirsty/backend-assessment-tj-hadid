package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewConsumer(url string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"fleet.events", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"geofence_alerts", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,         // queue name
		"geofence",     // routing key
		"fleet.events", // exchange
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn: conn,
		ch:   ch,
	}, nil
}

func (c *Consumer) Consume() error {
	msgs, err := c.ch.Consume(
		"geofence_alerts", // queue
		"geofence-worker", // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return err
	}

	for delivery := range msgs {
		log.Printf("Geofence alert received: %s", string(delivery.Body))

		var event GeofenceEvent
		if err := json.Unmarshal(delivery.Body, &event); err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			delivery.Nack(false, false)
			continue
		}

		log.Printf("ALERT: Vehicle %s entered geofence at lat: %f, lon: %f, timestamp: %d",
			event.VehicleID, event.Location.Latitude, event.Location.Longitude, event.Timestamp)

		delivery.Ack(false)
	}

	return nil
}

func (c *Consumer) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
