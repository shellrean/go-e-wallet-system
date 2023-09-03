package service

import (
	"encoding/json"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/config"
)

type emailService struct {
	cnf          *config.Config
	queueService domain.QueueService
}

func NewEmail(cnf *config.Config, queueService domain.QueueService) domain.EmailService {
	return &emailService{cnf, queueService}
}

func (e emailService) Send(to, subject, body string) error {
	payload := dto.EmailSendReq{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return e.queueService.Enqueue("send:email", data, 3)
}
