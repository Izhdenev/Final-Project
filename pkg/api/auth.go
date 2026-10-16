package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("final-project-secret-key")

func passwordHash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "метод не поддерживается"})
		return
	}

	var creds struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	want := os.Getenv("TODO_PASSWORD")
	if want == "" {
		writeJSON(w, map[string]string{"token": ""})
		return
	}

	if creds.Password != want {
		writeJSON(w, map[string]string{"error": "Неверный пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": passwordHash(want),
	})
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка формирования токена"})
		return
	}

	writeJSON(w, map[string]string{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) == 0 {
			next(w, r)
			return
		}

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
		if !ok || hash != passwordHash(pass) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
