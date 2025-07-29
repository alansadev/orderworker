package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;"`
	OrderNumber   int       `gorm:"autoIncrement"`
	Status        string
	TotalAmount   int64
	PaymentMethod string
	CreatedAt     time.Time
	Client        User
	ClientID      uuid.UUID
	Seller        User `gorm:"foreignKey:SellerID"`
	SellerID      *uuid.UUID
	Items         []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;"`
	Quantity        int
	PriceAtPurchase int64
	OrderID         uuid.UUID
	Product         Product
	ProductID       uuid.UUID
}
