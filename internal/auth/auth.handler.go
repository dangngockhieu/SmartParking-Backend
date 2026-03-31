package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	appErrors "backend/internal/common/errors"
	"backend/pkg/response"
)

type verifiedPageRenderer interface {
	RenderVerifiedPage(year int) (string, error)
}

type Handler struct {
	service     *Service
	mailService verifiedPageRenderer
}

// NewHandler tạo mới instance Handler với service và mailService đã được khởi tạo
func NewHandler(service *Service, mailService verifiedPageRenderer) *Handler {
	return &Handler{
		service:     service,
		mailService: mailService,
	}
}

// Register godoc
// @Summary Đăng ký tài khoản
// @Description Tạo tài khoản mới và gửi email xác thực
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Thông tin đăng ký"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	// Register Tài khoản mới
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	if err := h.service.Register(req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "Đăng ký thành công. Vui lòng kiểm tra email để xác thực tài khoản.", nil)
}

// VerifyEmail godoc
// @Summary Xác thực email
// @Description Xác thực tài khoản bằng email và mã code từ link gửi qua email
// @Tags auth
// @Produce html
// @Param email query string true "Email"
// @Param code query string true "Mã xác thực"
// @Success 200 {string} string "HTML verified page"
// @Failure 400 {object} map[string]interface{}
// @Router /auth/verify [get]
func (h *Handler) VerifyEmail(c *gin.Context) {
	// Xác thực tài khoản người dùng sau đăng ký
	email := c.Query("email")
	code := c.Query("code")

	if email == "" || code == "" {
		c.Error(appErrors.NewBadRequest("Thiếu email hoặc code"))
		return
	}

	if err := h.service.VerifyEmail(VerifyEmailRequest{
		Email: email,
		Code:  code,
	}); err != nil {
		c.Error(err)
		return
	}

	html, err := h.mailService.RenderVerifiedPage(time.Now().Year())
	if err != nil {
		c.Error(appErrors.NewInternal("Render trang xác thực thất bại"))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// ResendVerification godoc
// @Summary Gửi lại email xác thực
// @Description Gửi lại email xác thực tài khoản
// @Tags auth
// @Accept json
// @Produce json
// @Param request body EmailRequest true "Email người dùng"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/resend [post]
func (h *Handler) ResendVerification(c *gin.Context) {
	// Gửi lại email xác thực tài khoản
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	if err := h.service.ResendVerificationEmail(req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Resend email successful", nil)
}

// Login godoc
// @Summary Đăng nhập
// @Description Đăng nhập bằng email và mật khẩu
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Thông tin đăng nhập"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	// Đăng nhập và trả về access token, refresh token được lưu trong cookie
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	ip := c.ClientIP()

	// Validate thông tin đăng nhập và lấy thông tin người dùng
	u, err := h.service.ValidateUser(req.Email, req.Password, ip)
	if err != nil {
		c.Error(err)
		return
	}

	device := c.GetHeader("User-Agent")

	data, refreshToken, err := h.service.Login(u, device, ip)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", false, true)

	response.Success(c, http.StatusOK, "Đăng nhập thành công", data)
}

// Logout godoc
// @Summary Đăng xuất
// @Description Đăng xuất tài khoản hiện tại và xóa refresh token cookie
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	// Đăng xuất tài khoản hiện tại và xóa refresh token cookie
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

	jtiValue, _ := c.Get("jti")
	jti, _ := jtiValue.(string)

	expValue, _ := c.Get("exp")
	var ttl time.Duration
	if expUnix, ok := expValue.(int64); ok {
		ttl = time.Until(time.Unix(expUnix, 0))
	}

	refreshToken, _ := c.Cookie("refresh_token")

	if err := h.service.Logout(userID, refreshToken, jti, ttl); err != nil {
		c.Error(err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	response.Success(c, http.StatusOK, "Đăng xuất thành công", nil)
}

// RefreshToken godoc
// @Summary Làm mới access token
// @Description Cấp lại access token từ refresh token trong cookie
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh-token [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	// Cấp lại access token từ refresh token trong cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.Error(appErrors.NewUnauthorized("Missing refresh token"))
		return
	}

	data, err := h.service.Refresh(refreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Làm mới token thành công", data)
}

// SendResetPassword godoc
// @Summary Gửi email đặt lại mật khẩu
// @Description Gửi email chứa hướng dẫn hoặc mã để đặt lại mật khẩu
// @Tags auth
// @Accept json
// @Produce json
// @Param request body EmailRequest true "Email người dùng"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/send-reset-password [post]
func (h *Handler) SendResetPassword(c *gin.Context) {
	// Gửi email chứa hướng dẫn hoặc mã để đặt lại mật khẩu
	var req EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	if err := h.service.SendResetPasswordEmail(req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Password reset email sent successfully", nil)
}

// ResetPassword godoc
// @Summary Đặt lại mật khẩu
// @Description Đặt lại mật khẩu bằng mã xác thực hoặc token reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Thông tin đặt lại mật khẩu"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/reset-password [patch]
func (h *Handler) ResetPassword(c *gin.Context) {
	// Đặt lại mật khẩu bằng mã xác thực hoặc token reset
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	if err := h.service.ResetPassword(req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Password reset successfully", nil)
}
