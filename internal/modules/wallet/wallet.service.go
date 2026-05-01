package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidWebhook      = errors.New("invalid webhook")
)

type Service struct {
	repo                         *Repository
	payosClient                  PayOSClient
	returnURL                    string
	paymentUpdateStatusCancelURL string
}

func NewService(
	repo *Repository,
	payosClient PayOSClient,
	returnURL string,
	paymentUpdateStatusCancelURL string,
) *Service {
	return &Service{
		repo:                         repo,
		payosClient:                  payosClient,
		returnURL:                    returnURL,
		paymentUpdateStatusCancelURL: paymentUpdateStatusCancelURL,
	}
}

// Nạp tiền vào ví qua PayOS
func (s *Service) CreateDeposit(ctx context.Context, userID uint64, req CreateDepositRequest) (*CreateDepositResponse, error) {
	if req.Amount <= 0 {
		return nil, errors.New("invalid deposit amount")
	}

	if s.returnURL == "" || s.paymentUpdateStatusCancelURL == "" {
		return nil, errors.New("PAYOS_RETURN_URL or PAYMENT_UPDATE_STATUS_CANCEL_URL is not configured")
	}

	var itemID uint64
	var orderCode uint64

	// Chỉ tạo transaction PENDING trong DB transaction.
	// Không gọi PayOS trong transaction để tránh giữ lock lâu.
	err := s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.LockUserWallet(ctx, tx, userID)
		if err != nil {
			return err
		}

		desc := fmt.Sprintf("Wallet deposit: %d", req.Amount)

		item := &WalletTransaction{
			UserID:        user.ID,
			Type:          WalletTypeDeposit,
			Amount:        req.Amount,
			BalanceBefore: user.Money,
			BalanceAfter:  user.Money,
			Status:        WalletStatusPending,
			Description:   &desc,
		}

		if err := s.repo.CreateTransaction(ctx, tx, item); err != nil {
			return err
		}

		itemID = item.ID
		orderCode = uint64(10000100) + item.ID

		if err := s.repo.UpdateOrderCode(ctx, tx, item.ID, orderCode); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Gọi PayOS sau khi DB đã commit.
	paymentLink, err := s.payosClient.CreatePaymentLink(ctx, CreatePaymentLinkInput{
		OrderCode:   int64(orderCode),
		Amount:      int(req.Amount),
		Description: fmt.Sprintf("Nap vi %d", orderCode),
		ReturnURL:   s.returnURL,
		CancelURL:   s.paymentUpdateStatusCancelURL,
	})

	if err != nil {
		_ = s.repo.MarkTransactionFailed(ctx, itemID, err.Error())
		return nil, err
	}

	return &CreateDepositResponse{
		TransactionID: itemID,
		OrderCode:     orderCode,
		Amount:        req.Amount,
		PaymentURL:    paymentLink.CheckoutURL,
	}, nil
}

func (s *Service) UpdateWalletCancel(ctx context.Context, orderCode uint64) error {
	reason := "User canceled the payment" + fmt.Sprintf(" (order_code: %d)", orderCode)
	return s.repo.MarkTransactionCanceled(ctx, orderCode, reason)
}

// Xử lý webhook từ PayOS
func (s *Service) HandlePayOSWebhook(ctx context.Context, req PayOSWebhookRequest) error {
	// SDK verify signature, trả error nếu không hợp lệ.
	data, err := s.payosClient.VerifyWebhook(ctx, req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidWebhook, err)
	}

	// Chỉ cộng tiền khi PayOS báo thành công.
	if data.Code != "00" {
		return nil
	}

	orderCode := uint64(data.OrderCode)
	amount := data.Amount
	transactionID := data.Reference
	if transactionID == "" {
		transactionID = data.PaymentLinkID
	}

	if orderCode == 0 {
		return errors.New("webhook orderCode is empty")
	}

	if amount <= 0 {
		return errors.New("webhook amount is invalid")
	}

	return s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.LockTransactionByOrderCode(ctx, tx, orderCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// PayOS có thể gửi webhook test/confirm với orderCode không có trong hệ thống.
				// Trả nil để endpoint phản hồi 2xx, giúp xác nhận webhook URL thành công.
				return nil
			}
			return err
		}

		// Chống webhook retry làm cộng tiền nhiều lần.
		if item.Status == WalletStatusSuccess {
			return nil
		}

		if item.Type != WalletTypeDeposit {
			return errors.New("transaction type is not DEPOSIT")
		}

		if item.Status != WalletStatusPending {
			return nil
		}

		if item.Amount != amount {
			return errors.New("webhook amount mismatch")
		}

		user, err := s.repo.LockUserWallet(ctx, tx, item.UserID)
		if err != nil {
			return err
		}

		balanceBefore := user.Money
		balanceAfter := user.Money + item.Amount

		if err := s.repo.UpdateUserMoney(ctx, tx, user.ID, balanceAfter); err != nil {
			return err
		}

		if err := s.repo.MarkTransactionSuccess(
			ctx,
			tx,
			item.ID,
			balanceBefore,
			balanceAfter,
			transactionID,
		); err != nil {
			return err
		}

		return nil
	})
}

// Trừ tiền từ ví người dùng
func (s *Service) DeductWallet(ctx context.Context, req DeductWalletRequest) error {
	if req.UserID == 0 {
		return errors.New("userId is required")
	}

	if req.Amount <= 0 {
		return errors.New("invalid deduction amount")
	}

	return s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user, err := s.repo.LockUserWallet(ctx, tx, req.UserID)
		if err != nil {
			return err
		}

		if user.Money < req.Amount {
			return ErrInsufficientBalance
		}

		balanceBefore := user.Money
		balanceAfter := user.Money - req.Amount

		if err := s.repo.UpdateUserMoney(ctx, tx, user.ID, balanceAfter); err != nil {
			return err
		}

		desc := req.Description
		if desc == "" {
			desc = "Parking fee payment"
		}

		item := &WalletTransaction{
			UserID:        user.ID,
			Type:          WalletTypeDeduct,
			Amount:        req.Amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Status:        WalletStatusSuccess,
			Description:   &desc,
		}

		if err := s.repo.CreateTransaction(ctx, tx, item); err != nil {
			return err
		}

		return nil
	})
}

// Lấy lịch sử giao dịch của người dùng với pagination kiểu cursor
func (s *Service) GetMyTransactions(
	ctx context.Context,
	userID uint64,
	cursorCreatedAt *time.Time,
	cursorID *uint64,
	limit int,
) (*CursorTransactionResponse, error) {
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	rows, err := s.repo.FindUserTransactionsWithCursor(
		ctx,
		userID,
		cursorCreatedAt,
		cursorID,
		limit,
	)

	if err != nil {
		return nil, err
	}

	hasNextPage := len(rows) > limit

	if hasNextPage {
		rows = rows[:limit]
	}

	var nextCursor *WalletCursor

	if hasNextPage && len(rows) > 0 {
		last := rows[len(rows)-1]

		nextCursor = &WalletCursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		}
	}

	return &CursorTransactionResponse{
		Data:        rows,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
