package rabbits

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap/buffer"
)

const (
	JSON = "application/json"

	ExchangeKey string = "mail_service_exchange"

	ExchangeDirect = "direct"

	ConfirmUserEmailKey             string = "mail_service_confirm_user_email"
	SuccessConfirmUserEmailKey      string = "mail_service_success_confirm_user_email"
	ConfirmUserEmailQueueKey        string = "mail_service_confirm_user_email_queue"
	SuccessConfirmUserEmailQueueKey string = "mail_service_success_confirm_user_email_queue"
)

type IPublisher interface {
	Publish(ctx context.Context, exchange, routingKey string, data *buffer.Buffer) error
	Close() error
}

type IConsumer interface {
	Consume(context.Context, string) (<-chan amqp091.Delivery, error)
	Close() error
}
