package dto

import "sync"

type Hub struct {
	mu                  sync.RWMutex
	NotificationChannel map[int64]chan NotificationData
}

func NewHub() *Hub {
	return &Hub{
		NotificationChannel: make(map[int64]chan NotificationData),
	}
}

func (h *Hub) Register(userID int64) chan NotificationData {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.NotificationChannel == nil {
		h.NotificationChannel = make(map[int64]chan NotificationData)
	}
	if oldCh, exists := h.NotificationChannel[userID]; exists {
		close(oldCh)
	}

	ch := make(chan NotificationData, 16)
	h.NotificationChannel[userID] = ch
	return ch
}

func (h *Hub) Unregister(userID int64, ch chan NotificationData) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.NotificationChannel == nil {
		return
	}
	if current, exists := h.NotificationChannel[userID]; exists && current == ch {
		delete(h.NotificationChannel, userID)
		close(ch)
	}
}

func (h *Hub) Publish(userID int64, data NotificationData) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.NotificationChannel == nil {
		return
	}
	if ch, ok := h.NotificationChannel[userID]; ok {
		select {
		case ch <- data:
		default:
		}
	}
}
