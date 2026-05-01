package wallet

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type UserWallet struct {
	ID    uint64 `gorm:"column:id"`
	Money int64  `gorm:"column:money"`
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

// CreateTransaction tạo một giao dịch mới trong ví người dùng
func (r *Repository) CreateTransaction(ctx context.Context, tx *gorm.DB, item *WalletTransaction) error {
	return tx.WithContext(ctx).Create(item).Error
}

// UpdateOrderCode cập nhật order code cho giao dịch
func (r *Repository) UpdateOrderCode(ctx context.Context, tx *gorm.DB, id uint64, orderCode uint64) error {
	return tx.WithContext(ctx).
		Model(&WalletTransaction{}).
		Where("id = ?", id).
		Update("order_code", orderCode).Error
}

// MarkTransactionFailed đánh dấu giao dịch là thất bại
func (r *Repository) MarkTransactionFailed(ctx context.Context, id uint64, reason string) error {
	return r.db.WithContext(ctx).
		Model(&WalletTransaction{}).
		Where("id = ? AND status = ?", id, WalletStatusPending).
		Updates(map[string]any{
			"status":      WalletStatusFailed,
			"description": reason,
		}).Error
}

// MarkTransactionCanceled đánh dấu giao dịch là hủy
func (r *Repository) MarkTransactionCanceled(ctx context.Context, orderCode uint64, reason string) error {
	return r.db.WithContext(ctx).
		Model(&WalletTransaction{}).
		Where("order_code = ? AND status = ?", orderCode, WalletStatusPending).
		Updates(map[string]any{
			"status":      WalletStatusCanceled,
			"description": reason,
		}).Error
}

// LockTransactionByOrderCode khóa giao dịch theo order code
func (r *Repository) LockTransactionByOrderCode(ctx context.Context, tx *gorm.DB, orderCode uint64) (*WalletTransaction, error) {
	var item WalletTransaction

	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("order_code = ?", orderCode).
		First(&item).Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// LockUserWallet khóa bản ghi ví của người dùng để tránh xung đột khi cập nhật số tiền
func (r *Repository) LockUserWallet(ctx context.Context, tx *gorm.DB, userID uint64) (*UserWallet, error) {
	var user UserWallet

	err := tx.WithContext(ctx).
		Table("users").
		Select("id, money").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserMoney cập nhật số tiền trong ví của người dùng
func (r *Repository) UpdateUserMoney(ctx context.Context, tx *gorm.DB, userID uint64, money int64) error {
	return tx.WithContext(ctx).
		Table("users").
		Where("id = ?", userID).
		Update("money", money).Error
}

// MarkTransactionSuccess đánh dấu giao dịch là thành công và cập nhật số tiền trong ví của người dùng
func (r *Repository) MarkTransactionSuccess(
	ctx context.Context,
	tx *gorm.DB,
	id uint64,
	balanceBefore int64,
	balanceAfter int64,
	transactionID string,
) error {
	updates := map[string]any{
		"status":         WalletStatusSuccess,
		"balance_before": balanceBefore,
		"balance_after":  balanceAfter,
	}

	if transactionID != "" {
		updates["transaction_id"] = transactionID
	}

	return tx.WithContext(ctx).
		Model(&WalletTransaction{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// FindUserTransactionsWithCursor lấy danh sách giao dịch của người dùng theo cursor-based pagination
func (r *Repository) FindUserTransactionsWithCursor(
	ctx context.Context,
	userID uint64,
	cursorCreatedAt *time.Time,
	cursorID *uint64,
	limit int,
) ([]WalletTransaction, error) {
	db := r.db.WithContext(ctx).
		Where("user_id = ?", userID)

	if cursorCreatedAt != nil && cursorID != nil {
		db = db.Where(
			"(created_at < ? OR (created_at = ? AND id < ?))",
			*cursorCreatedAt,
			*cursorCreatedAt,
			*cursorID,
		)
	}

	var rows []WalletTransaction

	err := db.
		Order("created_at DESC").
		Order("id DESC").
		Limit(limit + 1).
		Find(&rows).Error

	return rows, err
}
