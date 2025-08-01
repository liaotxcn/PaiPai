package ctxdata

import "github.com/golang-jwt/jwt/v4"

const IdentityKey = "peninsula"

// GetJwtToken 生成 JWT 令牌
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[IdentityKey] = uid

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
