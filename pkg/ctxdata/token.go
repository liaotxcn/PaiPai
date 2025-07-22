package ctxdata

import "github.com/golang-jwt/jwt/v4"

const Identitf = "paipai"

func GetJwtToken(secreKey string, iat, seconds int64, uid string) (string, error) {
	claims := jwt.MapClaims{}
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[Identitf] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secreKey))
}
