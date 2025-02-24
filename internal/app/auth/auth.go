// Package auth для работы с сессиями пользователей.
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ndreyserg/ushort/internal/app/logger"
)

// Session интерфейс взаимодействия с сессией
type Session interface {
	// Open открывает новую сессию, если сессии не найдено, возвращает ID пользователя в случае успеха.
	Open(w http.ResponseWriter, r *http.Request) (string, error)
	// GetID возвращает ID пользователя из сессии.
	GetID(r *http.Request) (string, error)
}

const tokenKey = "token"

func generateID() string {
	id := make([]byte, 10)
	rand.Read(id)
	return hex.EncodeToString(id)
}

type claims struct {
	jwt.RegisteredClaims
	UserID string
}

type jwtSession struct {
	secretKey string
}

func (j *jwtSession) newToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
		},
		UserID: userID,
	})
	return token.SignedString([]byte(j.secretKey))
}

// Open открывает новую сессию, если сессии не найдено, возвращает ID пользователя в случае успеха.
func (j *jwtSession) Open(w http.ResponseWriter, r *http.Request) (string, error) {
	id, err := j.GetID(r)

	if err != nil {
		id = generateID()
	}

	strToken, err := j.newToken(id)

	if err != nil {
		logger.Log.Error(err)
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     tokenKey,
		Value:    strToken,
		HttpOnly: true,
	})

	return id, nil
}

// GetID возвращает ID пользователя из сессии.
func (j *jwtSession) GetID(r *http.Request) (string, error) {

	var strToken string

	cookies := r.Cookies()

	for _, cookie := range cookies {
		if cookie.Name == tokenKey {
			strToken = cookie.Value
			break
		}
	}
	if strToken == "" {
		return "", errors.New("empty token")
	}
	claims := &claims{}
	token, err := jwt.ParseWithClaims(strToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

// NewJWTSession создает объект работы с сессиями на основе JWT
func NewJWTSession(secret string) Session {
	return &jwtSession{
		secretKey: secret,
	}
}
