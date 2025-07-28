package options

type RabbitMQQueueOptions struct {
	Name       string
	Durable    bool           // if true, the queue will survive a broker restart
	Exclusive  bool           // if true, the queue will be deleted when the connection that created it closes
	AutoDelete bool           // if true, the queue will be deleted when the last consumer unbinds
	Args       map[string]any // additional properties for the queue
}
