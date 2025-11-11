package domain

import "time"

type Wallet struct {
	ID        string    `json:"wallet_id"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ApplyOperationRequest struct {
	WalletID      string `json:"wallet_id" validate:"required"`
	OperationType string `json:"operation_type" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
}
