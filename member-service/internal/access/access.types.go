package access

import (
	"member_service/internal/cookies"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	JWTSecret  string
	JWTExpires int
}

// Claims : jwt Claims格式
type Claims struct {
	MemberId int
	UserName string
	jwt.RegisteredClaims
}

const CookieKeyJWT cookies.Key = "jwtToken"
