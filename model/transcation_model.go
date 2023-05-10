package model

import (
	"time"

	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type TransactionEntity struct {
	ID                int            `json:"id"`
	UserId            int            `json:"user_id"`
	Amount            int            `json:"amount"`
	SawerUserId       int            `json:"sawer_user_id"`
	TransactionStatus string         `json:"transaction_status"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at"`
}

type TransactionResponse struct {
	ID          int            `json:"id"`
	UserId      int            `json:"user_id"`
	Amount      int            `json:"amount"`
	SawerUserId int            `json:"sawer_user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Midtrans    *snap.Response `json:"midtrans,omitempty"`
}

type TransactionRequest struct {
	UserId      int `json:"user_id,omitempty"`
	Amount      int `json:"amount"`
	SawerUserId int `json:"sawer_user_id"`
}

type NotificationPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionId     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	OrderId           string `json:"order_id"`
}
