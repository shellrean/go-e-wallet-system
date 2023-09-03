package domain

type QueueService interface {
	Enqueue(name string, data []byte, retry int) error
}
