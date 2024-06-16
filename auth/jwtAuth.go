package auth

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secret_key = []byte("Synapse_Rocks")
)

type claims struct {
	Id string
	*jwt.RegisteredClaims
}

func IsAuthorizedToken(tokenStr string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) { return secret_key, nil })

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("signature is invalid...")
			return false, nil
		}
		log.Println("Bad Request...")
		return false, err
	}

	if !token.Valid {
		log.Println("Unauthorized token...")
		return false, nil
	}

	return true, nil
}

func GetUserIdFromToken(tokenStr string) string {
	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) { return secret_key, nil })
	if err != nil {
		log.Println("Error in parsing Token to extract id")
		return ""
	}

	claim := token.Claims.(*claims)
	return claim.Id
}
