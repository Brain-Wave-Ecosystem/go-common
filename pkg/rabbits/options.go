package rabbits

import (
	"github.com/rabbitmq/amqp091-go"
)

var _ Option = (*funcOption)(nil)

type Option interface {
	apply(*amqp091.Channel) error
}

func WithExchange(name, kind string, durable, autoDelete bool, args amqp091.Table) Option {
	return newOptFunc(func(c *amqp091.Channel) error {
		return c.ExchangeDeclare(name, kind, durable, autoDelete, false, false, args)
	})
}

func WithQueue(name string, durable, autoDelete, exclusive bool, args amqp091.Table) Option {
	return newOptFunc(func(c *amqp091.Channel) error {
		_, err := c.QueueDeclare(name, durable, autoDelete, exclusive, false, args)
		return err
	})
}

func WithQueueAndBind(queueName, routingKey, exchange string, durable, autoDelete bool, args amqp091.Table) Option {
	return newOptFunc(func(c *amqp091.Channel) error {
		_, err := c.QueueDeclare(queueName, durable, autoDelete, false, false, args)
		if err != nil {
			return err
		}

		return c.QueueBind(queueName, routingKey, exchange, false, args)
	})
}

type funcOption struct {
	f func(*amqp091.Channel) error
}

func (fdo *funcOption) apply(ch *amqp091.Channel) error {
	return fdo.f(ch)
}

func newOptFunc(f func(*amqp091.Channel) error) *funcOption {
	return &funcOption{f: f}
}
