package service

import (
	"github.com/hibiken/asynq"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/internal/component"
	"shellrean.id/belajar-auth/internal/config"
)

type queueService struct {
	queueClient *asynq.Client
}

func NewQueue(cnf *config.Config) domain.QueueService {
	redisConnection := asynq.RedisClientOpt{
		Addr:     cnf.Queue.Addr,
		Password: cnf.Queue.Pass,
	}
	client := asynq.NewClient(redisConnection)

	return &queueService{
		queueClient: client,
	}
}

func (q queueService) Enqueue(name string, data []byte, retry int) error {
	task := asynq.NewTask(name, data, asynq.MaxRetry(retry))

	info, err := q.queueClient.Enqueue(task)
	if err != nil {
		component.Log.Error("error when enqueue: ", err.Error())
		return err
	}

	component.Log.Info("enqueue-client id: ", info.ID)
	return err
}
