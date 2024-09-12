package config

import "fmt"

type (
	HTTP struct {
		Host    string `env:"HTTP_HOST" envDefault:"localhost"`
		Port    int    `env:"HTTP_PORT" envDefault:"3000"`
		Prefork bool   `env:"HTTP_PREFORK" envDefault:"false"`
	}

	GRPC struct {
		Host string `env:"HOST" envDefault:"localhost"`
		Port int    `env:"PORT" envDefault:"9000"`
	}

	DB struct {
		Host            string `env:"HOST" envDefault:"localhost"`
		Port            int    `env:"PORT" envDefault:"5432"`
		User            string `env:"USER" envDefault:"root"`
		Pass            string `env:"PASS" envDefault:"changeme"`
		Name            string `env:"NAME" envDefault:"db"`
		ApplicationName string `env:"APPLICATION_NAME" envDefault:"peatio"`
		SSLMode         bool
	}

	Postgres DB

	QuestDB DB

	JWT struct {
		PrivateKeyPath string `env:"JWT_PRIVATE_KEY_PATH" envDefault:"/secrets/barong.key"`
		PublicKey      string `env:"JWT_PUBLIC_KEY"`
	}

	Elasticsearch struct {
		URL      []string `env:"ES_URL" envDefault:"http://localhost:9200"`
		Username string   `env:"ES_USERNAME"`
		Password string   `env:"ES_PASSWORD"`
	}

	Redis struct {
		Host     string `env:"REDIS_HOST" envDefault:"localhost"`
		Port     int    `env:"REDIS_PORT" envDefault:"6379"`
		Password string `env:"REDIS_PASSWORD"`
	}

	Kafka struct {
		Brokers           []string `env:"KAFKA_BROKERS" envDefault:"localhost:9092" envSeparator:","`
		Partitions        int32    `env:"KAFKA_PARTITIONS" envDefault:"1"`
		ReplicationFactor int16    `env:"KAFKA_REPLICATION_FACTOR" envDefault:"1"`
	}

	Vault struct {
		Address         string `env:"VAULT_ADDR" envDefault:"localhost:8200"`
		Token           string `env:"VAULT_TOKEN" envDefault:"changeme"`
		ApplicationName string `env:"VAULT_APP_NAME" envDefault:"peatio"`
	}

	Twilio struct {
		PhoneNumber string `env:"TWILIO_PHONE_NUMBER"`
		AccountSid  string `env:"TWILIO_ACCOUNT_SID"`
		AuthToken   string `env:"TWILIO_AUTH_TOKEN"`
	}

	ObjectStorage struct {
		Bucket       string `env:"OBJECT_STORAGE_BUCKET" envDefault:"barong"`
		Region       string `env:"OBJECT_STORAGE_REGION" envDefault:"us-east-1"`
		AccessKey    string `env:"OBJECT_STORAGE_ACCESS_KEY"`
		AccessSecret string `env:"OBJECT_STORAGE_ACCESS_SECRET"`
		Endpoint     string `env:"OBJECT_STORAGE_ENDPOINT" envDefault:""`
		Version      int    `env:"OBJECT_STORAGE_VERSION" envDefault:"2"`
	}
)

func (h HTTP) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

func (g GRPC) Address() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}
