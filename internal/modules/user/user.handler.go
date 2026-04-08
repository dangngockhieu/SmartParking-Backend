package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	appErrors "backend/internal/common/errors"
	"backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// FindWithPagination godoc
// @Summary Lấy danh sách người dùng
// @Description Lấy danh sách người dùng có phân trang và tìm kiếm
// @Tags user
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang hiện tại" default(1)
// @Param limit query int false "Số lượng mỗi trang" default(10)
// @Param search query string false "Từ khóa tìm kiếm"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users [get]
func (h *Handler) FindWithPagination(c *gin.Context) {
	// Lấy người dùng và phân trang
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	search := c.Query("search")

	result, err := h.service.FindWithPagination(page, pageSize, search)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy danh sách người dùng thành công", result)
}

// GetMyInfo godoc
// @Summary Lấy thông tin người dùng hiện tại
// @Description Lấy thông tin của người dùng đang đăng nhập
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/my-info [get]
func (h *Handler) GetMyInfo(c *gin.Context) {
	// Lấy thông tin người dùng hiện tại
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}
	userID, ok := userIDValue.(uint64)
	if !ok {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}

	user, err := h.service.GetMyInfo(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy thông tin người dùng thành công", user)
}

// CreateByAdmin godoc
// @Summary Admin tạo người dùng
// @Description Tạo mới người dùng bởi quản trị viên
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "Thông tin người dùng"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users [post]
func (h *Handler) CreateByAdmin(c *gin.Context) {
	// Tạo người dùng mới bởi admin
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	if err := h.service.CreateUserByAdmin(req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "Tạo người dùng thành công", nil)
}

// ChangePassword godoc
// @Summary Đổi mật khẩu
// @Description Người dùng hiện tại đổi mật khẩu của chính mình
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Thông tin đổi mật khẩu"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/change-password [patch]
func (h *Handler) ChangePassword(c *gin.Context) {
	// Đổi mật khẩu
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}

	if err := h.service.ChangePassword(userID, req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Đổi mật khẩu thành công", nil)
}

// ChangeRole godoc
// @Summary Đổi vai trò người dùng
// @Description Quản trị viên đổi vai trò của người dùng theo ID
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID người dùng"
// @Param request body ChangeRoleRequest true "Vai trò mới"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/change-role/{id} [patch]
func (h *Handler) ChangeRole(c *gin.Context) {
	// Đổi vai trò người dùng
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequest("id không hợp lệ"))
		return
	}

	var req ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	if err := h.service.ChangeRole(uint64(id64), req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Đổi vai trò thành công", nil)
}

// UpdateProfile godoc
// @Summary Cập nhật thông tin cá nhân
// @Description Người dùng hiện tại cập nhật thông tin cá nhân của chính mình
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangeProfileRequest true "Thông tin cập nhật"
// @Success 200 {object} map[string]interface{}
func (h *Handler) UpdateProfile(c *gin.Context) {
	var req ChangeProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}
	userID, ok := userIDValue.(uint64)
	if !ok {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}
	updatedUser, err := h.service.ChangeProfile(userID, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "Cập nhật thông tin thành công", updatedUser)
}

// GetWalletBalance godoc
// @Summary Lấy số dư ví
// @Description Người dùng hiện tại lấy số dư trong ví của mình
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/wallet/balance [get]
func (h *Handler) GetMyWallet(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}
	userID, ok := userIDValue.(uint64)
	if !ok {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}
	balance, err := h.service.GetWalletBalance(userID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "Lấy số dư ví thành công", map[string]int64{"balance": balance})
}
