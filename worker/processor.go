package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/terenjit/simplebank/db/sqlc"
)

const (
	Queue_Critical = "critical"
	Queue_Default  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redis asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redis, asynq.Config{
		Queues: map[string]int{
			Queue_Critical: 10,
			Queue_Default:  5,
		},
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
