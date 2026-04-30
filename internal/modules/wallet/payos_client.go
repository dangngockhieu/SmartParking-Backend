package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	payossdk "github.com/payOSHQ/payos-lib-golang/v2"

	"backend/configs"
)

// PayOSClient là interface nội bộ để dễ mock trong test.
type PayOSClient interface {
	CreatePaymentLink(ctx context.Context, input CreatePaymentLinkInput) (*CreatePaymentLinkOutput, error)
	// VerifyWebhook trả về (webhookData, error). Nếu signature không hợp lệ thì trả error.
	VerifyWebhook(ctx context.Context, payload PayOSWebhookRequest) (*PayOSWebhookData, error)
}

// CreatePaymentLinkInput và Output giữ nguyên để không phá service layer.
type CreatePaymentLinkInput struct {
	OrderCode   int64
	Amount      int
	Description string
	ReturnURL   string
	CancelURL   string
}

type CreatePaymentLinkOutput struct {
	CheckoutURL string
}

// sdkPayOSClient wraps payos-lib-golang/v2.
type sdkPayOSClient struct {
	client *payossdk.PayOS
}

// disabledPayOSClient được dùng khi chưa cấu hình PayOS.
type disabledPayOSClient struct{}

// NewPayOSClient tạo client từ config. Nếu thiếu credential, trả về disabled client.
func NewPayOSClient(cfg *configs.Config) PayOSClient {
	if cfg == nil || cfg.PayOSClientID == "" || cfg.PayOSAPIKey == "" || cfg.PayOSChecksumKey == "" {
		return &disabledPayOSClient{}
	}

	client, err := payossdk.NewPayOS(&payossdk.PayOSOptions{
		ClientId:    cfg.PayOSClientID,
		ApiKey:      cfg.PayOSAPIKey,
		ChecksumKey: cfg.PayOSChecksumKey,
		Timeout:     15 * time.Second,
	})
	if err != nil {
		// Nếu NewPayOS lỗi (ví dụ credential rỗng), fallback về disabled.
		return &disabledPayOSClient{}
	}

	return &sdkPayOSClient{client: client}
}

// ---- sdkPayOSClient ----

func (c *sdkPayOSClient) CreatePaymentLink(ctx context.Context, input CreatePaymentLinkInput) (*CreatePaymentLinkOutput, error) {
	if input.OrderCode == 0 {
		return nil, errors.New("orderCode is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount is invalid")
	}
	if input.Description == "" {
		return nil, errors.New("description is required")
	}
	if input.ReturnURL == "" || input.CancelURL == "" {
		return nil, errors.New("returnUrl and cancelUrl are required")
	}

	resp, err := c.client.PaymentRequests.Create(ctx, payossdk.CreatePaymentLinkRequest{
		OrderCode:   input.OrderCode,
		Amount:      input.Amount,
		Description: input.Description,
		ReturnUrl:   input.ReturnURL,
		CancelUrl:   input.CancelURL,
	})
	if err != nil {
		return nil, err
	}

	if resp.CheckoutUrl == "" {
		return nil, errors.New("payOS did not return checkoutUrl")
	}

	return &CreatePaymentLinkOutput{CheckoutURL: resp.CheckoutUrl}, nil
}

func (c *sdkPayOSClient) VerifyWebhook(ctx context.Context, payload PayOSWebhookRequest) (*PayOSWebhookData, error) {
	webhookBody := map[string]interface{}{
		"code":      payload.Code,
		"desc":      payload.Desc,
		"success":   payload.Success,
		"data":      payload.Data,
		"signature": payload.Signature,
	}

	result, err := c.client.Webhooks.VerifyData(ctx, webhookBody)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidWebhook, err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok || resultMap == nil {
		return nil, fmt.Errorf("%w: unexpected verify result type %T", ErrInvalidWebhook, result)
	}

	raw, err := json.Marshal(resultMap)
	if err != nil {
		return nil, fmt.Errorf("%w: marshal verified data failed: %v", ErrInvalidWebhook, err)
	}

	var data PayOSWebhookData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("%w: unmarshal verified data failed: %v", ErrInvalidWebhook, err)
	}

	return &data, nil
}

// ---- disabledPayOSClient ----

func (c *disabledPayOSClient) CreatePaymentLink(_ context.Context, _ CreatePaymentLinkInput) (*CreatePaymentLinkOutput, error) {
	return nil, errors.New("payOS is not configured")
}

func (c *disabledPayOSClient) VerifyWebhook(_ context.Context, _ PayOSWebhookRequest) (*PayOSWebhookData, error) {
	return nil, fmt.Errorf("%w: payOS client is disabled (missing/invalid config)", ErrInvalidWebhook)
}
