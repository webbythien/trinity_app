package worker

import (
	"github.com/hrshadhin/fiber-go-boilerplate/app/task"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/workers"
)

func Setting(broker, resultBackend string) *config.WorkerConfig {
	return &config.WorkerConfig{
		Broker:        broker,
		ResultBackend: resultBackend,
		Workers: map[string]config.Task{
			"core_api_queue": map[string]interface{}{
				"Worker.HealthCheck": task.HealthCheck,
				// "Worker.QuizAnswer":  tasks.QuizAnswer,
			},
		},
	}
}

func WorkerExecute(queueName, consume string, concurrency int) error {
	getConfig := config.WrkCfg()
	wcf := Setting(getConfig.Broker, getConfig.ResultBackend)
	cnf := config.NewWorker(queueName, wcf)
	return workers.Execute(cnf, consume, concurrency)
}
