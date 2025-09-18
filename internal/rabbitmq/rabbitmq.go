package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Komilov31/delayed-notifier/internal/config"
	"github.com/Komilov31/delayed-notifier/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

const (
	retries = 3
)

type RabbitMq struct {
	pulisher *rabbitmq.Publisher
	consumer <-chan amqp.Delivery
}

func New() *RabbitMq {
	url := fmt.Sprintf(
		"amqp://guest:guest@%s%s/",
		config.Cfg.RabbitMq.Host,
		config.Cfg.RabbitMq.Port,
	)

	connection, err := rabbitmq.Connect(url, retries, 0)
	if err != nil {
		log.Fatal("could not connect to rabbitmq server: ", err)
	}

	pubCh, err := connection.Channel()
	if err != nil {
		log.Fatal("could not create channel for rabbitmq: ", err)
	}

	qm := rabbitmq.NewQueueManager(pubCh)
	_, err = qm.DeclareQueue("notification")
	if err != nil {
		log.Fatal("could not create queue for rabbitmq: ", err)
	}

	publisher := rabbitmq.NewPublisher(pubCh, "")

	conCh, err := connection.Channel()
	if err != nil {
		log.Fatal("could not create channel for consumer rabbitmq: ", err)
	}

	deliveries, err := conCh.Consume("notification", "", false, false, false, false, amqp.Table{})
	if err != nil {
		log.Fatal("could not create consumer for rabbitmq: ", err)
	}

	return &RabbitMq{
		pulisher: publisher,
		consumer: deliveries,
	}
}

func (r *RabbitMq) Publish(notification model.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		zlog.Logger.Error().Msg("could not marshal notification to send to rabbitmq: " + err.Error())
	}

	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  2,
	}

	return r.pulisher.PublishWithRetry(body, "notification", "application/json", strategy)
}

func (r *RabbitMq) Consume(ctx context.Context) (<-chan []byte, error) {
	messages := make(chan []byte)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(messages)
				return
			default:
				next, ok := <-r.consumer
				if !ok {
					return
				}

				if err := next.Ack(false); err != nil {
					zlog.Logger.Error().Msg("could not acknowledge message consuming: " + err.Error())
				}
				messages <- next.Body
			}
		}
	}()

	return messages, nil
}
