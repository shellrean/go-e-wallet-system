package service

import (
	"context"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
)

type notificationService struct {
	notificationRepository domain.NotificationRepository
}

func NewNotification(notificationRepository domain.NotificationRepository) domain.NotificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
	}
}

func (n notificationService) FindByUser(ctx context.Context, user int64) ([]dto.NotificationData, error) {
	notifications, err := n.notificationRepository.FindByUser(ctx, user)
	if err != nil {
		return nil, err
	}

	var result []dto.NotificationData
	for _, v := range notifications {
		result = append(result, dto.NotificationData{
			ID:        v.ID,
			Title:     v.Title,
			Body:      v.Body,
			Status:    v.Status,
			IsRead:    v.IsRead,
			CreatedAt: v.CreatedAt,
		})
	}
	if result == nil {
		result = make([]dto.NotificationData, 0)
	}

	return result, nil
}
