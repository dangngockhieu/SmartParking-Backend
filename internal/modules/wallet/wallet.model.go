package wallet

import "time"

const (
	WalletTypeDeposit = "DEPOSIT"
	WalletTypeDeduct  = "DEDUCT"

	WalletStatusPending  = "PENDING"
	WalletStatusSuccess  = "SUCCESS"
	WalletStatusFailed   = "FAILED"
	WalletStatusCanceled = "CANCELED"
)

type WalletTransaction struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"column:user_id;not null" json:"userId"`
	Type          string    `gorm:"column:type;not null" json:"type"`
	Amount        int64     `gorm:"column:amount;not null" json:"amount"`
	BalanceBefore int64     `gorm:"column:balance_before;not null" json:"balanceBefore"`
	BalanceAfter  int64     `gorm:"column:balance_after;not null" json:"balanceAfter"`
	Status        string    `gorm:"column:status;not null" json:"status"`
	OrderCode     *uint64   `gorm:"column:order_code" json:"orderCode,omitempty"`
	TransactionID *string   `gorm:"column:transaction_id" json:"transactionId,omitempty"`
	Description   *string   `gorm:"column:description" json:"description,omitempty"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}
