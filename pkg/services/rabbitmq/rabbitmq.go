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
	Headers       amqp.Table    // Application or header exchange table
	DeliveryMode  uint8         // Transient (0 or 1) or Persistent (2)
	Priority      uint8         // 0 to 9
	CorrelationId string        // Correlation identifier
	ReplyTo       string        // Address to to reply to (ex: RPC)
	Expiration    string        // Message expiration spec
	MessageId     string        // Message identifier
	Timestamp     time.Time     // Message timestamp
	Type          string        // Message type name
	UserId        string        // Creating user id - ex: "guest"
	AppId         string        // Creating application id
	Body          proto.Message // Application specific payload of the message
}

func (srv *Service) Publish(ch *amqp.Channel, exchange, routingKey string, message Message) error {
	b, err := proto.Marshal(message.Body)
	if err != nil {
		return err
	}
	return ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType:   "application/vnd.google.protobuf",
		Headers:       message.Headers,
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
