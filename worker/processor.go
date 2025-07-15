package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/terenjit/simplebank/db/sqlc"
	"github.com/terenjit/simplebank/mail"
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
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redis asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(redis, asynq.Config{
		Queues: map[string]int{
			Queue_Critical: 10,
			Queue_Default:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("process task failed")
		}),
		Logger: NewLogger(),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
