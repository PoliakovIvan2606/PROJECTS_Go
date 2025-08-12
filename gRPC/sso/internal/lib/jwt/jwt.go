package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"gRPC/internal/domain/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	app.Secret = strings.TrimSpace(app.Secret)
	app.Secret = strings.ReplaceAll(app.Secret, `\n`, "\n")

	tokenEcdsa, err := ParseECDSAPrivateKeyFromPEM(app.Secret)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString(tokenEcdsa)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}


func ParseECDSAPrivateKeyFromPEM(pemStr string) (*ecdsa.PrivateKey, error) {
    block, _ := pem.Decode([]byte(pemStr))
    if block == nil {
        return nil, errors.New("failed to decode PEM block")
    }

    // Попытка распарсить как EC PRIVATE KEY
    key, err := x509.ParseECPrivateKey(block.Bytes)
    if err == nil {
        return key, nil
    }

    // Если не получилось, пытаемся как PKCS8 PRIVATE KEY
    parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    ecdsaKey, ok := parsedKey.(*ecdsa.PrivateKey)
    if !ok {
        return nil, errors.New("key is not ECDSA private key")
    }

    return ecdsaKey, nil
}

