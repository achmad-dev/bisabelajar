package config

type Env struct {
	// PostgreSQL Environment Variables
	PostgresUser     string `env:"POSTGRES_USER" json:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" json:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB" json:"POSTGRES_DB"`
	PostgresHost     string `env:"POSTGRES_HOST" json:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT" json:"POSTGRES_PORT"`

	// RabbitMQ Environment Variables
	RabbitMQUser           string `env:"RABBITMQ_DEFAULT_USER" json:"RABBITMQ_DEFAULT_USER"`
	RabbitMQPass           string `env:"RABBITMQ_DEFAULT_PASS" json:"RABBITMQ_DEFAULT_PASS"`
	RabbitMQHost           string `env:"RABBITMQ_HOST" json:"RABBITMQ_HOST"`
	RabbitMQPort           string `env:"RABBITMQ_PORT" json:"RABBITMQ_PORT"`
	RabbitMQManagementPort string `env:"RABBITMQ_MANAGEMENT_PORT" json:"RABBITMQ_MANAGEMENT_PORT"`

	// Redis Environment Variables
	RedisPassword string `env:"REDIS_PASSWORD" json:"REDIS_PASSWORD"`
	RedisPort     string `env:"REDIS_PORT" json:"REDIS_PORT"`
	RedisHost     string `env:"REDIS_HOST" json:"REDIS_HOST"`

	// App
	Port string `env:"PORT" json:"PORT"`
}
