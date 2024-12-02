package workers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"errors"

	"github.com/google/uuid"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"

	tracers "github.com/RichardKnop/machinery/example/tracers"
	opentracing "github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"
)

var (
	WorkerConfig *config.WorkerConfig
)

func StartServer(cnf *config.Worker) (*machinery.Server, error) {
	_server, err := machinery.NewServer(cnf.Config)
	if err != nil {
		return nil, err
	}

	return _server, _server.RegisterTasks(cnf.Task)
}

func Execute(cnf *config.Worker, tag string, concurrency int) error {
	server, err := StartServer(cnf)
	if err != nil {
		return err
	}

	consumerTag := tag

	cleanup, err := tracers.SetupTracer(consumerTag)
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, concurrency)

	// Here we inject some custom code for error handling,
	// start and end of task hooks, useful for metrics for example.
	errorHandler := func(err error) {
		log.ERROR.Println("error handler:", err)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("start of task handler for:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("end of task handler for:", signature.Name)
	}

	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()
}

/*
	asyncResult returns value

asyncResult, err := server.SendTaskWithContext(ctx, task)

	if err != nil {
		return fmt.Errorf("could not send task: %s", err.Error())
	}

results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))

	if err != nil {
		return fmt.Errorf("getting task result failed with error: %s", err.Error())
	}

log.INFO.Printf(tasks.HumanReadableResults(results))
*/
func Delay(queue, taskName string, fn interface{}, args ...interface{}) error {
	if WorkerConfig == nil {
		WorkerConfig = config.WrkCfg()
		if WorkerConfig == nil {
			return errors.New("worker config is required to launch")
		}
	}
	worker := config.NewWorker(queue, WorkerConfig)
	server, err := StartServer(worker)
	if err != nil {
		return err
	}
	// Check if the function has the correct type and the corresponding number of arguments
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return errors.New("invalid function")
	}
	if fnType.NumIn() != len(args) {
		return errors.New("incorrect number of arguments")
	}

	var _args []tasks.Arg
	for _, arg := range args {
		_arg := tasks.Arg{
			Type:  reflect.ValueOf(arg).Type().Name(),
			Value: reflect.ValueOf(arg).Interface(),
		}
		_args = append(_args, _arg)
	}
	task := &tasks.Signature{
		Name: taskName,
		Args: _args,
	}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "delay")
	defer span.Finish()

	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))

	log.INFO.Println("Starting batch:", batchID)
	asyncResult, err := server.SendTaskWithContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not send task: %s", err.Error())
	}
	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return fmt.Errorf("getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf(tasks.HumanReadableResults(results))
	return nil
}
