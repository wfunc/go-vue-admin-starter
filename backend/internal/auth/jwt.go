package auth

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      int      `json:"user_id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	RoleID      int      `json:"role_id"`
	RoleCode    string   `json:"role_code"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewManager(cfg config.JWTConfig) *Manager {
	return &Manager{secret: []byte(cfg.Secret), issuer: cfg.Issuer, ttl: cfg.TTL}
}

func (m *Manager) Generate(user CurrentUser) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.ttl)
	claims := Claims{
		UserID:      user.UserID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		RoleID:      user.RoleID,
		RoleCode:    user.RoleCode,
		Permissions: user.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	return signed, expiresAt, err
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
