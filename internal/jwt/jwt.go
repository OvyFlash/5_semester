package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//JWT struct that keeps both keys and allows to work with jwt
type JWT struct {
	privateKey []byte
	publicKey  []byte
	verifyKey  *rsa.PublicKey
	signKey    *rsa.PrivateKey
}

//NewJWT returns instance of JWT
func NewJWT() *JWT {
	signBytes, err := ioutil.ReadFile("/home/oleg/Рабочий стол/KPI/5 Семестр/Системи обробки сигналів/NECKNAME SocialNetwork/5_semester-master/internal/config/keys/app.rsa")
	//signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.PrivateKey))
	if err != nil {
		if len(os.Args) != 4 || os.Args[2] == "" {
			panic(err)
		}
		signBytes, err = ioutil.ReadFile(os.Args[2])
		if err != nil {
			panic(err)
		}
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}
	verifyBytes, err := ioutil.ReadFile("/home/oleg/Рабочий стол/KPI/5 Семестр/Системи обробки сигналів/NECKNAME SocialNetwork/5_semester-master/internal/config/keys/app.rsa.pub")
	//verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(config.PublicKey))
	if err != nil {
		if len(os.Args) != 4 || os.Args[3] == "" {
			panic(err)
		}
		verifyBytes, err = ioutil.ReadFile(os.Args[3])
		if err != nil {
			panic(err)
		}
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}
	return &JWT{
		privateKey: signBytes,
		publicKey:  verifyBytes,
		signKey:    signKey,
		verifyKey:  verifyKey,
	}
}

//GetToken generates token
func (j *JWT) GetToken(userid int32) string {
	token := jwt.New(jwt.SigningMethodRS512)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	claims["userid"] = userid
	token.Claims = claims

	tokenString, _ := token.SignedString(j.signKey)

	return tokenString
}

//IsTokenValid val is a value that comes with request
func (j *JWT) IsTokenValid(val string) (int32, error) {
	token, err := jwt.Parse(val, func(token *jwt.Token) (interface{}, error) {
		return j.verifyKey, nil
	})

	switch err.(type) {
	case nil:
		if !token.Valid {
			return -1, errors.New("Token is invalid")
		}

		var (
			userID int32
		)

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return -1, errors.New("Token is invalid")
		}

		userID = int32(claims["userid"].(float64))
		
		if userID == 0 {
			return -1, nil
		}
		return userID, nil

	case *jwt.ValidationError:
		vErr := err.(*jwt.ValidationError)

		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			return -1, errors.New("Token expired, get a new one")
		default:
			fmt.Println(vErr)
			return -1, errors.New("Error while parsing token")
		}
	default:
		return -1, errors.New("Unable to parse token")
	}
}
