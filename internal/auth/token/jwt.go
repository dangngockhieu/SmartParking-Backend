package token

import (
	"errors"
	"time"

	"backend/configs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	cfg *configs.Config
}

// AccessToken chứa các thông tin về userID, email, role
type AccessClaims struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RefreshToken chỉ chứa userID và email
type RefreshClaims struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewService(cfg *configs.Config) *Service {
	return &Service{cfg: cfg}
}

// Tạo Access Token với thời hạn 15 phút
func (s *Service) CreateAccessToken(userID uint64, email, role string) (string, error) {
	now := time.Now()

	claims := AccessClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTAccessSecret))
}

// Tạo Refresh Token với thời hạn 7 ngày
func (s *Service) CreateRefreshToken(userID uint64, email string) (string, error) {
	now := time.Now()

	claims := RefreshClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTRefreshSecret))
}

// Xác thực Access Token và trả về thông tin trong token nếu hợp lệ
func (s *Service) VerifyAccessToken(raw string) (*ClaimsPayload, error) {
	token, err := jwt.ParseWithClaims(raw, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTAccessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return &ClaimsPayload{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
		JTI:    claims.ID,
		Exp:    claims.ExpiresAt.Time.Unix(),
	}, nil
}

// Xác thực Refresh Token và trả về thông tin trong token nếu hợp lệ
func (s *Service) VerifyRefreshToken(raw string) (*ClaimsPayload, error) {
	token, err := jwt.ParseWithClaims(raw, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTRefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return &ClaimsPayload{
		UserID: claims.UserID,
		Email:  claims.Email,
		JTI:    claims.ID,
		Exp:    claims.ExpiresAt.Time.Unix(),
	}, nil
}
