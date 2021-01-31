package utils

import (
	"5g-v2x-user-service/internal/config"
	"crypto/sha512"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

var SEPARATOR string = ":"
var CFG *config.Config = config.NewConfig()

// AccessToken ...
type AccessToken struct {
	Username  string
	Expires   string
	Signature string
}

// GenerateAccessToken ...
func GenerateAccessToken(username, hashedPassword string) (string, error) {
	lifetimeDuration, err := time.ParseDuration(CFG.AccessTokenLifetime)
	if err != nil {
		return "", err
	}
	expires := fmt.Sprintf("%d", int64(time.Now().Add(lifetimeDuration).Unix()))
	signature, err := generateSignature(username, expires, hashedPassword)
	if err != nil {
		return "", err
	}
	token := strings.Join([]string{username, expires, signature}, SEPARATOR)
	token = base64encode(token)
	return token, nil
}

// ExtractAccessToken ...
func ExtractAccessToken(token string) (*AccessToken, error) {
	decoded, err := base64decode(token)
	if err != nil {
		return nil, err
	}

	arr := strings.Split(decoded, SEPARATOR)
	if len(arr) != 3 {
		return nil, errors.New("Invalid access token length")
	}

	// fmt.Println(accessToken)

	return &AccessToken{
		Username:  arr[0],
		Expires:   arr[1],
		Signature: arr[2],
	}, nil
}

// VerifyAccessToken ...
func VerifyAccessToken(accessToken *AccessToken, hashedPassword string) error {
	expectedSignature, err := generateSignature(accessToken.Username, accessToken.Expires, hashedPassword)
	if err != nil {
		return err
	}
	if accessToken.Signature != expectedSignature {
		return errors.New("Invalid access token")
	}

	return nil
}

func generateSignature(username string, expires string, hashedPassword string) (string, error) {
	signature := strings.Join([]string{username, expires, hashedPassword, CFG.AccessTokenSecret}, SEPARATOR)
	signature, err := hash(signature)
	if err != nil {
		return "", err
	}
	return signature, nil
}

func hash(text string) (string, error) {
	hashed, err := HashAndSalt([]byte(text))
	if err != nil {
		return "", err
	}
	return hashed, nil
}

func base64encode(text string) string {
	sEnc := b64.StdEncoding.EncodeToString([]byte(text))
	return sEnc
}

func base64decode(text string) (string, error) {
	sDec, err := b64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	return string(sDec), nil
}

// HashAndSalt ...
func HashAndSalt(text []byte) (string, error) {
	hash := []byte(text)
	return fmt.Sprintf("%x", sha512.Sum512(hash)), nil
}

// WasExpired ...
func WasExpired(timestamp time.Time) bool {
	return time.Now().After(timestamp)
}
