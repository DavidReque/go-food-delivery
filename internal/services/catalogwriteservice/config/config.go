package config

// GormOptions representa la configuración de la base de datos PostgreSQL.
type GormOptions struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
	SSLMode  bool   `json:"sslMode"`
}

// RabbitMQHostOptions representa la configuración de conexión al host de RabbitMQ.
type RabbitMQHostOptions struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	HostName string `json:"hostName"`
	Port     int    `json:"port"`
	HTTPPort int    `json:"httpPort"`
}

// RabbitMQOptions representa las opciones generales de RabbitMQ.
type RabbitMQOptions struct {
	AutoStart           bool                `json:"autoStart"`
	Reconnecting        bool                `json:"reconnecting"`
	RabbitMQHostOptions RabbitMQHostOptions `json:"rabbitmqHostOptions"`
}

// Config es la estructura principal que agrupa todas las opciones de configuración.
type Config struct {
	GormOptions     GormOptions     `json:"gormOptions"`
	RabbitMQOptions RabbitMQOptions `json:"rabbitmqOptions"`
}