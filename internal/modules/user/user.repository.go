package user

import (
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindByEmail tìm người dùng theo email
func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", strings.TrimSpace(strings.ToLower(email))).Take(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID tìm người dùng theo ID
func (r *Repository) FindByID(id uint64) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create thêm mới người dùng vào database
func (r *Repository) Create(user *User) error {
	return r.db.Create(user).Error
}

// UpdatePassword cập nhật mật khẩu đã được hash cho người dùng
func (r *Repository) UpdatePassword(id uint64, hashedPassword string) error {
	return r.db.Model(&User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// UpdateRole cập nhật vai trò cho người dùng
func (r *Repository) UpdateRole(id uint64, role Role) error {
	return r.db.Model(&User{}).Where("id = ?", id).Update("role", role).Error
}

// FindWithPagination tìm người dùng với phân trang và tìm kiếm
func (r *Repository) FindWithPagination(page, pageSize int, search string) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&User{}).Where("is_verified = ?", true)

	search = strings.TrimSpace(search)
	if search != "" {
		like := search + "%"
		query = query.Where(
			r.db.Where("email LIKE ?", like),
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Select("id", "first_name", "last_name", "email", "role").
		Order("id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateProfile cập nhật thông tin cá nhân cho người dùng
func (r *Repository) UpdateProfile(id uint64, first_name, last_name *string) (*User, error) {
	updateData := make(map[string]interface{})

	if first_name != nil {
		updateData["first_name"] = *first_name
	}
	if last_name != nil {
		updateData["last_name"] = *last_name
	}

	err := r.db.Model(&User{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

// Lấy số tiền hiện tại trong ví của người dùng
func (r *Repository) GetWalletBalance(userID uint64) (int64, error) {
	var user User
	err := r.db.Select("money").First(&user, userID).Error
	if err != nil {
		return 0, err
	}
	return user.Money, nil
}

// DepositToWallet nạp tiền vào ví của người dùng
func (r *Repository) DepositToWallet(userID uint64, amount int64) error {
	return r.db.Model(&User{}).Where("id = ?", userID).UpdateColumn("money", gorm.Expr("money + ?", amount)).Error
}

// WithdrawFromWallet rút tiền từ ví của người dùng
func (r *Repository) WithdrawFromWallet(userID uint64, amount int64) error {
	return r.db.Model(&User{}).Where("id = ?", userID).UpdateColumn("money", gorm.Expr("money - ?", amount)).Error
}
