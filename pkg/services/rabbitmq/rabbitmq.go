package rabbitmq

import (
	"github.com/golang/protobuf/proto"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Service struct {
	config *Config
	conn   *amqp.Connection
}

type Config struct {
	Addr string `envconfig:"addr" default:"amqp://expected:expected@localhost/expected"`
}

func New(config *Config) *Service {
	return &Service{
		config: config,
		conn:   nil,
	}
}

func NewFromEnv() *Service {
	config := &Config{}
	if err := envconfig.Process("RABBITMQ", config); err != nil {
		logrus.WithError(err).Fatalln("unable to process environment configuration")
	}
	return New(config)
}

func (srv *Service) Name() string {
	return "rabbitmq"
}

func (srv *Service) Start() error {
	conn, err := amqp.Dial(srv.config.Addr)
	if err != nil {
		return err
	}
	srv.conn = conn
	return nil
}

func (srv *Service) Stop() error {
	return srv.conn.Close()
}

func (srv *Service) Started() bool {
	return srv.conn != nil && !srv.conn.IsClosed()
}

func (srv *Service) Client() *amqp.Connection {
	return srv.conn
}

type Message struct {
	// Properties
	DeliveryMode  uint8     // Transient (0 or 1) or Persistent (2)
	Priority      uint8     // 0 to 9
	CorrelationId string    // correlation identifier
	ReplyTo       string    // address to to reply to (ex: RPC)
	Expiration    string    // message expiration spec
	MessageId     string    // message identifier
	Timestamp     time.Time // message timestamp
	Type          string    // message type name
	UserId        string    // creating user id - ex: "guest"
	AppId         string    // creating application id

	// The application specific payload of the message
	Body proto.Message
}

func (srv *Service) Publish(ch *amqp.Channel, exchange, routingKey string, message Message) error {
	b, err := proto.Marshal(message.Body)
	if err != nil {
		return err
	}
	return ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType:   "application/vnd.google.protobuf",
		DeliveryMode:  message.DeliveryMode,
		Priority:      message.Priority,
		CorrelationId: message.CorrelationId,
		ReplyTo:       message.ReplyTo,
		Expiration:    message.Expiration,
		MessageId:     message.MessageId,
		Timestamp:     message.Timestamp,
		Type:          message.Type,
		UserId:        message.UserId,
		AppId:         message.AppId,
		Body:          b,
	})
}
