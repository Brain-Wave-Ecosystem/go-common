package rabbits

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap/buffer"
)

var _ IPublisher = (*Publisher)(nil)

type Publisher struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel

	contentType string
}

func NewPublisher(url, contentType string, opts ...Option) (*Publisher, error) {
	var p Publisher

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	p.ch = ch
	p.conn = conn
	p.contentType = contentType

	for _, opt := range opts {
		err = opt.apply(ch)
		if err != nil {
			return nil, err
		}
	}

	return &p, nil
}

func (p *Publisher) ExchangeDeclare(name, kind string, durable, autoDelete bool, args amqp091.Table) error {
	err := p.ch.ExchangeDeclare(name, kind, durable, autoDelete, false, false, args)
	if err != nil {
		return fmt.Errorf("exchangeDeclare: %w", err)
	}

	return nil
}

func (p *Publisher) QueueDeclare(name string, durable, autoDelete, exclusive bool, args amqp091.Table) error {
	_, err := p.ch.QueueDeclare(name, durable, autoDelete, exclusive, false, args)
	if err != nil {
		return fmt.Errorf("queueDeclare error: %s", err.Error())
	}
	return nil
}

func (p *Publisher) QueueDeclareAndBind(queueName, routingKey, exchange string, durable, autoDelete bool, args amqp091.Table) error {
	err := p.QueueDeclare(queueName, durable, autoDelete, false, args)
	if err != nil {
		return err
	}

	err = p.ch.QueueBind(queueName, routingKey, exchange, false, args)
	if err != nil {
		return fmt.Errorf("queueDeclareAndBind error: %s", err.Error())
	}

	return nil
}

func (p *Publisher) Publish(ctx context.Context, exchange, routingKey string, data *buffer.Buffer) error {
	return p.ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: p.contentType,
		Body:        data.Bytes(),
	})
}

func (p *Publisher) Close() error {
	return p.conn.Close()
}
