package types

import (
	"fmt"

	"github.com/DavidReque/go-food-delivery/internal/pkg/logger/defaultlogger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/rabbitmq/config"
	errorUtils "github.com/DavidReque/go-food-delivery/internal/pkg/utils/errorutils"

	"emperror.dev/errors"
	"github.com/rabbitmq/amqp091-go"
)

type internalConnection struct {
	cfg *config.RabbitmqOptions
	*amqp091.Connection
	isConnected       bool
	errConnectionChan chan error
	errChannelChan    chan error
	reconnectedChan   chan struct{}
}

type IConnection interface {
	IsClosed() bool
	IsConnected() bool
	// Channel gets a new channel on this internalConnection
	Channel() (*amqp091.Channel, error)
	Close() error
	ReConnect() error
	NotifyClose(receiver chan *amqp091.Error) chan *amqp091.Error
	Raw() *amqp091.Connection
	ErrorConnectionChannel() chan error
	ReconnectedChannel() chan struct{}
}

func NewRabbitMQConnection(cfg *config.RabbitmqOptions) (IConnection, error) {
	// https://levelup.gitconnected.com/connecting-a-service-in-golang-to-a-rabbitmq-server-835294d8c914
	if cfg.RabbitmqHostOptions == nil {
		return nil, errors.New("rabbitmq host options is nil")
	}

	c := &internalConnection{
		cfg:               cfg,
		errConnectionChan: make(chan error),
		// errChannelChan:    make(chan error),
		reconnectedChan: make(chan struct{}),
	}

	err := c.connect()
	if err != nil {
		return nil, err
	}

	if cfg.Reconnecting {
		go c.handleReconnecting()
	}

	return c, err
}

func (c *internalConnection) Close() error {
	return c.Connection.Close()
}

func (c *internalConnection) IsConnected() bool {
	return c.isConnected
}

func (c *internalConnection) ErrorConnectionChannel() chan error {
	return c.errConnectionChan
}

func (c *internalConnection) ReconnectedChannel() chan struct{} {
	return c.reconnectedChan
}

func (c *internalConnection) ReConnect() error {
	if c.Connection.IsClosed() == false {
		return nil
	}

	return c.connect()
}

func (c *internalConnection) Raw() *amqp091.Connection {
	return c.Connection
}

func (c *internalConnection) Channel() (*amqp091.Channel, error) {
	ch, err := c.Connection.Channel()
	//notifyChannelClose := ch.NotifyClose(make(chan *amqp091.Error))
	//go func() {
	//	<-notifyChannelClose //Listen to notifyChannelClose
	//	c.errChannelChan <- errors.New("Channel Closed")
	//}()

	return ch, err
}

func (c *internalConnection) connect() error {
	conn, err := amqp091.Dial(c.cfg.RabbitmqHostOptions.AmqpEndPoint())
	if err != nil {
		return errors.WrapIf(
			err,
			fmt.Sprintf(
				"Error in connecting to rabbitmq with host: %s",
				c.cfg.RabbitmqHostOptions.AmqpEndPoint(),
			),
		)
	}

	c.Connection = conn
	c.isConnected = true

	// https://stackoverflow.com/questions/41991926/how-to-detect-dead-rabbitmq-connection
	// Crea un canal de Go. Este canal se usará para recibir una notificación.
	notifyClose := c.Connection.NotifyClose(make(chan *amqp091.Error))

	// Lanza un proceso en segundo plano (un "vigilante") que escucha el canal de notificación.
	go func() {
		defer errorUtils.HandlePanic()

		// 3. La goroutine se queda "atascada" aquí, esperando recibir algo del canal.
		chanErr := <-notifyClose //Listen to notifyClose

		// La ejecución solo continúa cuando la conexión se ha caído
		// En este punto, `chanErr` contiene el error que causó el cierre.

		c.isConnected = false // Actualiza el estado para que el resto de la app sepa que no hay conexión.

		// Envía el error a otro canal (`c.errConnectionChan`).
		//    Probablemente hay otra parte del código escuchando en `errConnectionChan`
		//    que se encargará de iniciar el proceso de reconexión.
		c.errConnectionChan <- errors.WrapIf(chanErr, "rabbitmq Connection Closed with an error.")
	}()

	return nil
}

func (c *internalConnection) handleReconnecting() {
	defer errorUtils.HandlePanic()
	for {
		select {
		case err := <-c.errConnectionChan:
			if err != nil {
				defaultlogger.GetLogger().
					Info("Rabbitmq Connection Reconnecting...")
				err := c.connect()
				if err != nil {
					defaultlogger.GetLogger().
						Error(fmt.Sprintf("Error in reconnecting, %s", err))
					continue
				}

				defaultlogger.GetLogger().
					Info("Rabbitmq Connection Reconnected")
				c.isConnected = true
				c.reconnectedChan <- struct{}{}
				continue
			}
		}
	}
}
