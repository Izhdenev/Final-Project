package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("final-project-secret-key")

func passwordHash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	var creds struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	if authPassword == "" {
		writeJSON(w, http.StatusOK, map[string]string{"token": ""})
		return
	}

	if creds.Password != authPassword {
		writeError(w, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": passwordHash(authPassword),
	})
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка формирования токена")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	if authPassword == "" {
		return next
	}

	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		hash, ok := claims["hash"].(string)
		if !ok || hash != passwordHash(authPassword) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
