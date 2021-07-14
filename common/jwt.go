package common

import (
	"admin_demo/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId:         user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //过期时间
			IssuedAt: time.Now().Unix(),      //发放时间
			Issuer: "dwc",                    //发放者
			Subject: "user token",            //主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //SigningMethodHS256，注意别用错加密方法
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}


func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, e error) {
		return jwtKey, nil
	})

	return token, claims, err
}
