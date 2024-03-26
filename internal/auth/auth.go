package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
	jwtSecretKey           = "some-secret-key"
	jwtRefreshSecretKey    = "some-refresh-secret-key"
)

var user *User

func GetJWTSecret() string {
	return jwtSecretKey
}

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func GetRefreshJWTSecret() string {
	return jwtRefreshSecretKey
}

func GenerateTokensAndSetCookies(user *User, c echo.Context) error {
	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(user, exp, c)
	refreshToken, exp, err := generateRefreshToken(user)
	if err != nil {
		return err
	}
	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	return nil
}

func SetUser(newUser *User) {
	user = newUser
}

func GetUserID() string {
	if user != nil {
		return user.ID
	}

	return ""
}
func JWTErrorChecker(c echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
}

func generateRefreshToken(user *User) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func generateAccessToken(user *User) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

func generateToken(user *User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	claims := &Claims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	zap.L().Info(
		"generateToken",
		zap.String("userID", user.ID),
		zap.String("time", expirationTime.String()),
		zap.String("secret", string(secret)),
		zap.String("token", tokenString),
	)

	return tokenString, expirationTime, nil
}

func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func setUserCookie(user *User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.ID
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}
