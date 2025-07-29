package models

import (
	"github.com/google/uuid"
	"time"
)

type EventPayload struct {
	EventType string    `json:"eventType"`
	OrderID   uuid.UUID `json:"orderId"`
	NewStatus string    `json:"newStatus,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
