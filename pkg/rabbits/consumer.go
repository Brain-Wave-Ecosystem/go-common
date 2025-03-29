package rabbits

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

var _ IConsumer = (*Consumer)(nil)

type Consumer struct {
	conn        *amqp091.Connection
	ch          *amqp091.Channel
	contentType string
}

func NewConsumer(url, contentType string, opts ...Option) (*Consumer, error) {
	var c Consumer

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("create consumer clients error: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("create consumer channel error: %w", err)
	}

	c.ch = ch
	c.conn = conn
	c.contentType = contentType

	for _, opt := range opts {
		err = opt.apply(ch)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c *Consumer) ExchangeDeclare(name, kind string, durable, autoDelete bool, args amqp091.Table) error {
	err := c.ch.ExchangeDeclare(name, kind, durable, autoDelete, false, false, args)
	if err != nil {
		return fmt.Errorf("exchangeDeclare: %w", err)
	}

	return nil
}

func (c *Consumer) QueueDeclare(name string, durable, autoDelete, exclusive bool, args amqp091.Table) error {
	_, err := c.ch.QueueDeclare(name, durable, autoDelete, exclusive, false, args)
	if err != nil {
		return fmt.Errorf("queueDeclare error: %s", err.Error())
	}
	return nil
}

func (c *Consumer) QueueDeclareAndBind(queueName, routingKey, exchange string, durable, autoDelete bool, args amqp091.Table) error {
	err := c.QueueDeclare(queueName, durable, autoDelete, false, args)
	if err != nil {
		return err
	}

	err = c.ch.QueueBind(queueName, routingKey, exchange, false, args)
	if err != nil {
		return fmt.Errorf("queueDeclareAndBind error: %s", err.Error())
	}

	return nil
}

func (c *Consumer) Consume(ctx context.Context, routingKey string) (<-chan amqp091.Delivery, error) {
	ch, err := c.ch.ConsumeWithContext(ctx, routingKey, "", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("consume error: %s", err.Error())
	}

	return ch, nil
}

func (c *Consumer) Close() error {
	return c.conn.Close()
}
