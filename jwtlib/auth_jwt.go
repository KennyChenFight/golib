package jwtlib

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

type Config struct {
	Payload map[string]interface{}
	claims jwt.MapClaims
	SignALG SignALG
	SecretKey []byte
	TokenTimeout time.Duration
}

func newJWTAuth(c Config) *JWTAuth {
	return &JWTAuth{config: c}
}

type JWTAuth struct {
	config Config
}

func (j *JWTAuth) Sign() (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range j.config.Payload {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(j.config.TokenTimeout).Unix()
	j.config.claims = claims

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, j.config.claims)
	return tokenClaims.SignedString(j.config.SecretKey)
}

func (j *JWTAuth) Verify(token string) (map[string]interface{}, error) {
	token = strings.Replace(token, "Bearer ", "", -1)
	// parse and verify the token string
	tokenClaims, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return j.config.SecretKey, nil
	})
	// detail for jwt token err message
	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				message = "token is malformed"
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				message = "token could not be verified because of signing problems"
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				message = "signature validation failed"
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "token is expired"
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				message = "token is not yet valid before sometime"
			} else {
				message = "can not handle this token"
			}
		}
		err = errors.New(message)
		return nil, err
	}
	claims, _ := tokenClaims.Claims.(jwt.MapClaims)
	return claims, nil
}

