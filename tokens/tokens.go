package tokens

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var secret string = "mySecretKey123"

func CreateAccessToken(uId interface{}) string {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":        uId,
		"ExpiresAt": time.Now().Add(time.Minute * 15),
		"IssuedAt":  time.Now(),
	})
	signedToken, _ := accessToken.SignedString([]byte(secret))

	return signedToken
}

func CreateRefreshToken(accessToken string) string {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":          primitive.NewObjectID(),
		"accessToken": accessToken,
		"ExpiresAt":   time.Now().Add(time.Hour * 2),
	})
	signedToken, _ := refreshToken.SignedString([]byte(secret))
	return signedToken
}

func HashToken(refreshToken string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	return string(hash)
}

func GetIdFromToken(tokenStr string) interface{} {
	claims := jwt.MapClaims{}
	var id interface{}
	jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	for key, val := range claims {
		if key == "id" {
			id = val
		}
	}
	return id
}
func CheckToken(refToken string, hashedToken string) bool {
	match := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(refToken))
	if match == nil {
		return true
	} else {
		return false
	}

}

func BeareringAccessToken(tk string) string {
	return "Bearer " + tk
}

func CheckTokensLifetime(tk string) bool {
	claims := jwt.MapClaims{}
	var expAt interface{}
	jwt.ParseWithClaims(tk, claims, func(t *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	for key, val := range claims {
		if key == "ExpiresAt" {
			expAt = val
		}
	}
	expAt1, _ := time.Parse(time.RFC3339, expAt.(string))
	if expAt1.After(time.Now()) {
		return true
	} else {
		return false
	}
}

func DeBeareringToken(tk string) string {
	return strings.Replace(tk, "Bearer ", "", 1)
}
