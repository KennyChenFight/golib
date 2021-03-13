// Package jwtlib is for encapsulating github.com/dgrijalva/jwt-go any operations
//
// As a quick start:
// 	userID := "user1"
//	customPayload := map[string]interface{}{
//		"userID": userID,
//	}
//	jwtLib := jwtlib.NewJWTLib(jwtlib.Config{
//		Payload:      customPayload,
//		SignALG:      jwtlib.HS256,
//		SecretKey:    []byte("secret"),
//		TokenTimeout: time.Hour,
//	})
//	token, err := jwtLib.Sign()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(token)
//
//	payload, err := jwtLib.Verify(token)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(payload)
package jwtlib

func NewJWTLib(c Config) Auth {
	return newJWTAuth(c)
}
