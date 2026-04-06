package user

import (
	"strings"

	appErrors "backend/internal/common/errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// FindWithPagination lấy danh sách người dùng với phân trang và tìm kiếm
func (s *Service) FindWithPagination(page, pageSize int, search string) (*UserPaginationResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	users, total, err := s.repo.FindWithPagination(page, pageSize, search)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy danh sách người dùng thất bại")
	}

	result := make([]UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, UserResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Role:      u.Role,
		})
	}

	return &UserPaginationResponse{
		Users: result,
		Total: total,
	}, nil
}

// Get My Info
func (s *Service) GetMyInfo(id uint64) (*MyAccountResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy thông tin người dùng thất bại")
	}
	if user == nil {
		return nil, appErrors.NewNotFound("Người dùng không tồn tại!")
	}
	return &MyAccountResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		Money:     user.Money,
	}, nil
}

// CreateUserByAdmin tạo người dùng mới bởi admin
func (s *Service) CreateUserByAdmin(req CreateUserRequest) error {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	exist, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return appErrors.NewInternal("Kiểm tra email thất bại")
	}
	if exist != nil {
		return appErrors.NewConflict("Email đã được đăng ký!")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return appErrors.NewInternal("Mã hóa mật khẩu thất bại")
	}

	user := &User{
		Email:      req.Email,
		Password:   string(hashed),
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Role:       req.Role,
		IsVerified: true,
	}

	if err := s.repo.Create(user); err != nil {
		return appErrors.NewInternal("Tạo người dùng thất bại")
	}

	return nil
}

// ChangePassword đổi mật khẩu cho người dùng
func (s *Service) ChangePassword(userID uint64, req ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return appErrors.NewBadRequest("Mật khẩu mới và xác nhận mật khẩu không khớp!")
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return appErrors.NewInternal("Lấy thông tin người dùng thất bại")
	}
	if user == nil {
		return appErrors.NewNotFound("Người dùng không tồn tại!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return appErrors.NewBadRequest("Mật khẩu cũ không đúng!")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return appErrors.NewInternal("Mã hóa mật khẩu thất bại")
	}

	if err := s.repo.UpdatePassword(userID, string(hashed)); err != nil {
		return appErrors.NewInternal("Đổi mật khẩu thất bại")
	}

	return nil
}

// ChangeRole đổi vai trò cho người dùng
func (s *Service) ChangeRole(userID uint64, req ChangeRoleRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return appErrors.NewInternal("Lấy thông tin người dùng thất bại")
	}
	if user == nil {
		return appErrors.NewNotFound("Người dùng không tồn tại!")
	}

	if err := s.repo.UpdateRole(userID, req.NewRole); err != nil {
		return appErrors.NewInternal("Đổi vai trò thất bại")
	}

	return nil
}

// ChangeProfile cập nhật thông tin cá nhân cho người dùng
func (s *Service) ChangeProfile(userID uint64, req ChangeProfileRequest) (*UserResponse, error) {
	// Hứng cả updatedUser và err
	updatedUser, err := s.repo.UpdateProfile(userID, req.FirstName, req.LastName)

	if err != nil {
		return nil, appErrors.NewInternal("Cập nhật thông tin người dùng thất bại")
	}

	// Trả data mới nhất về
	return &UserResponse{
		ID:        updatedUser.ID,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Email:     updatedUser.Email,
		Role:      updatedUser.Role,
	}, nil
}

// GetWalletBalance lấy số tiền hiện tại trong ví của người dùng
func (s *Service) GetWalletBalance(userID uint64) (int64, error) {
	balance, err := s.repo.GetWalletBalance(userID)
	if err != nil {
		return 0, appErrors.NewInternal("Lấy số tiền trong ví thất bại")
	}
	return balance, nil
}

// Cộng tiền vào ví của người dùng
func (s *Service) DepositToWallet(userID uint64, amount int64) error {
	if amount <= 0 {
		return appErrors.NewBadRequest("Số tiền phải lớn hơn 0")
	}
	if err := s.repo.DepositToWallet(userID, amount); err != nil {
		return appErrors.NewInternal("Cộng tiền vào ví thất bại")
	}
	return nil
}

// Trừ tiền từ ví của người dùng
func (s *Service) WithdrawFromWallet(userID uint64, amount int64) error {
	if amount <= 0 {
		return appErrors.NewBadRequest("Số tiền phải lớn hơn 0")
	}
	if err := s.repo.WithdrawFromWallet(userID, amount); err != nil {
		return appErrors.NewInternal("Trừ tiền từ ví thất bại")
	}
	return nil
}
