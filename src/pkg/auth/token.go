// Package auth 提供JWT Token的生成、验证与刷新
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	bizErr "github.com/euler/mtap/pkg/errors"
)

// Claims JWT自定义声明
type Claims struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	DepartmentID string `json:"department_id,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair Token对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// JWTManager JWT管理器
type JWTManager struct {
	Secret         []byte
	AccessExpire   time.Duration
	RefreshExpire  time.Duration
	Issuer         string
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secret string, accessExpire, refreshExpire time.Duration) *JWTManager {
	return &JWTManager{
		Secret:        []byte(secret),
		AccessExpire:  accessExpire,
		RefreshExpire: refreshExpire,
		Issuer:        "mtap",
	}
}

// GenerateTokenPair 生成AccessToken + RefreshToken
func (m *JWTManager) GenerateTokenPair(userID, username, role, deptID string) (*TokenPair, error) {
	now := time.Now()
	accessExpiry := now.Add(m.AccessExpire)

	accessClaims := &Claims{
		UserID:       userID,
		Username:     username,
		Role:         role,
		DepartmentID: deptID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    m.Issuer,
			Subject:   userID,
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(m.Secret)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    m.Issuer,
			Subject:   userID,
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(m.Secret)
	if err != nil {
		return nil, bizErr.Wrap(bizErr.ErrInternal, err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry.Unix(),
	}, nil
}

// ValidateToken 验证Token并提取Claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, bizErr.New(bizErr.ErrUnauthorized)
		}
		return m.Secret, nil
	})
	if err != nil {
		return nil, bizErr.NewWithDetail(bizErr.ErrUnauthorized, "token无效或已过期")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, bizErr.New(bizErr.ErrUnauthorized)
	}

	return claims, nil
}

// RefreshAccessToken 用RefreshToken刷新AccessToken
func (m *JWTManager) RefreshAccessToken(refreshTokenStr string) (*TokenPair, error) {
	claims, err := m.ValidateToken(refreshTokenStr)
	if err != nil {
		return nil, err
	}
	// 仅用于刷新，Role可能为空，此处生成新的token pair需要完整信息
	// 实际业务中应从DB重新加载用户信息
	return m.GenerateTokenPair(claims.UserID, claims.Username, claims.Role, claims.DepartmentID)
}
