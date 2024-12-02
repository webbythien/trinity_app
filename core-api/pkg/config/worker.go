package config

import (
	"os"

	"github.com/RichardKnop/machinery/v1/config"
)

type Worker struct {
	Task   map[string]interface{}
	Config *config.Config
}

type Task map[string]interface{}

type WorkerConfig struct {
	Broker        string
	ResultBackend string
	Workers       map[string]Task
}

func NewWorker(queueName string, wcf *WorkerConfig) *Worker {
	return &Worker{
		Config: &config.Config{
			Broker:          wcf.Broker,
			DefaultQueue:    queueName,
			ResultBackend:   wcf.ResultBackend,
			ResultsExpireIn: 3600,
			AMQP: &config.AMQPConfig{
				Exchange:      "machinery_exchange",
				ExchangeType:  "direct",
				BindingKey:    queueName,
				PrefetchCount: 3,
			},
		},
		Task: wcf.Workers[queueName],
	}
}

var workerConfig = &WorkerConfig{}

func WrkCfg() *WorkerConfig {
	return workerConfig
}

// LoadDBCfg loads DB configuration
func LoadWrkCfg() {
	workerConfig.Broker = os.Getenv("BROKER")
	workerConfig.ResultBackend = os.Getenv("RESULT_BACKEND")
}
