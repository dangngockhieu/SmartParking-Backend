package auth

import (
	"strings"

	"backend/internal/auth/token"
	"backend/internal/modules/user"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindUserByEmail tìm user theo email, trả về nil nếu không tìm thấy
func (r *Repository) FindUserByEmail(email string) (*user.User, error) {
	var u user.User

	err := r.db.
		Where("email = ?", strings.TrimSpace(strings.ToLower(email))).
		Take(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

// CreateUser tạo user mới trong database
func (r *Repository) CreateUser(u *user.User) error {
	return r.db.Create(u).Error
}

// UpdateUserVerified cập nhật trạng thái verified của user
func (r *Repository) UpdateUserVerified(userID uint64, verified bool) error {
	return r.db.Model(&user.User{}).
		Where("id = ?", userID).
		Update("is_verified", verified).Error
}

// UpdateUserPassword cập nhật mật khẩu đã hash của user
func (r *Repository) UpdateUserPassword(userID uint64, hashedPassword string) error {
	return r.db.Model(&user.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error
}

// CreateRefreshToken lưu refresh token vào database
func (r *Repository) CreateRefreshToken(rt *token.RefreshToken) error {
	return r.db.Create(rt).Error
}

// FindRefreshTokensByUserID lấy tất cả refresh token của user, sắp xếp theo thời gian tạo tăng dần
func (r *Repository) FindRefreshTokensByUserID(userID uint64) ([]token.RefreshToken, error) {
	var tokens []token.RefreshToken
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&tokens).Error
	return tokens, err
}

// DeleteRefreshTokenByID xóa refresh
func (r *Repository) DeleteRefreshTokenByID(id uint64) error {
	return r.db.Delete(&token.RefreshToken{}, id).Error
}
