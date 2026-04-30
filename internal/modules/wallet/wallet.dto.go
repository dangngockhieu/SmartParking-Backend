package wallet

import "time"

type CreateDepositRequest struct {
	Amount int64 `json:"amount" binding:"required,min=1000"`
}

type CreateDepositResponse struct {
	TransactionID uint64 `json:"transactionId"`
	OrderCode     uint64 `json:"orderCode"`
	Amount        int64  `json:"amount"`
	PaymentURL    string `json:"paymentUrl"`
}

// PayOSWebhookRequest khớp với cấu trúc payos.Webhook từ SDK.
// SDK's VerifyData() nhận interface{}, nên ta có thể truyền thẳng struct này.
type PayOSWebhookRequest struct {
	Code      string                 `json:"code"`
	Desc      string                 `json:"desc"`
	Success   bool                   `json:"success"`
	Data      map[string]interface{} `json:"data"`
	Signature string                 `json:"signature"`
}

// PayOSWebhookData là type nội bộ, dùng sau khi SDK đã verify.
type PayOSWebhookData struct {
	OrderCode              int64  `json:"orderCode"`
	Amount                 int64  `json:"amount"`
	Description            string `json:"description,omitempty"`
	AccountNumber          string `json:"accountNumber,omitempty"`
	Reference              string `json:"reference,omitempty"`
	TransactionDateTime    string `json:"transactionDateTime,omitempty"`
	Currency               string `json:"currency,omitempty"`
	PaymentLinkID          string `json:"paymentLinkId,omitempty"`
	Code                   string `json:"code,omitempty"`
	Desc                   string `json:"desc,omitempty"`
	CounterAccountBankID   string `json:"counterAccountBankId,omitempty"`
	CounterAccountBankName string `json:"counterAccountBankName,omitempty"`
	CounterAccountName     string `json:"counterAccountName,omitempty"`
	CounterAccountNumber   string `json:"counterAccountNumber,omitempty"`
	VirtualAccountName     string `json:"virtualAccountName,omitempty"`
	VirtualAccountNumber   string `json:"virtualAccountNumber,omitempty"`
}

type DeductWalletRequest struct {
	UserID           uint64  `json:"userId" binding:"required"`
	Amount           int64   `json:"amount" binding:"required,min=1"`
	ParkingSessionID *uint64 `json:"parkingSessionId"`
	Description      string  `json:"description"`
}

type CursorTransactionResponse struct {
	Data        []WalletTransaction `json:"data"`
	NextCursor  *WalletCursor       `json:"nextCursor"`
	HasNextPage bool                `json:"hasNextPage"`
}

type WalletCursor struct {
	CreatedAt time.Time `json:"createdAt"`
	ID        uint64    `json:"id"`
}
