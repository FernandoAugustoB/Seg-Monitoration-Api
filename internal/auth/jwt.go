package auth

import (
	thub.com/golang-jwt/jwt/v5"

	me"
)

var jwtKey = []byte("sua_chave_secreta_aqui")

func GenerateJWT(userID string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 8).Unix(),
	})

	return token.SignedString(jwtKey)
}
