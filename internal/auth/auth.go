package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
	jwtSecretKey           = "some-secret-key"
	jwtRefreshSecretKey    = "some-refresh-secret-key"
)

// GetAccessTokenCookieName get access token name
func GetAccessTokenCookieName() string {
	return accessTokenCookieName
}

// GetJWTSecret get secret
func GetJWTSecret() string {
	return jwtSecretKey
}

// GetSigningMethod get signing method name
func GetSigningMethod() *jwt.SigningMethodHMAC {
	return jwt.SigningMethodHS256
}

// Claims struct
type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

// GetRefreshJWTSecret get refresh token name
func GetRefreshJWTSecret() string {
	return jwtRefreshSecretKey
}

// GenerateTokensAndSetCookies generate and set cookie
func GenerateTokensAndSetCookies(c echo.Context, user *User) error {
	accessToken, accessTokenString, exp, err := generateAccessToken(user)
	if err != nil {
		return err
	}

	setTokenCookie(c, accessTokenCookieName, accessTokenString, exp)
	c.Set("user", accessToken)
	setUserCookie(c, user, exp)
	_, refreshTokenString, exp, err := generateRefreshToken(user)
	if err != nil {
		return err
	}
	setTokenCookie(c, refreshTokenCookieName, refreshTokenString, exp)

	return nil
}

// GetUserID get user
func GetUserID(c echo.Context) string {
	if c.Get("user") == nil {
		return ""
	}
	u := c.Get("user").(*jwt.Token)

	claims := u.Claims.(*Claims)

	return claims.ID
}

// JWTErrorChecker function for error handling
func JWTErrorChecker(c echo.Context, err error) error {
	if err != nil {
		zap.L().Error(
			"JWTErrorChecker",
			zap.Error(err),
		)
	}

	return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
}

func generateRefreshToken(user *User) (*jwt.Token, string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func generateAccessToken(user *User) (*jwt.Token, string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

func generateToken(user *User, expirationTime time.Time, secret []byte) (*jwt.Token, string, time.Time, error) {
	claims := &Claims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(GetSigningMethod(), claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, "", time.Now(), err
	}

	return token, tokenString, expirationTime, nil
}

func setTokenCookie(c echo.Context, name, token string, expiration time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func setUserCookie(c echo.Context, user *User, expiration time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.ID
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}
