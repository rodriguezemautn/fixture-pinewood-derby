// Package auth provee autenticación simple para administración.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func getPIN() string {
	if pin := os.Getenv("ADMIN_PIN"); pin != "" {
		return pin
	}
	return "1234"
}

func secretKey() []byte {
	h := sha256.Sum256([]byte(getPIN()))
	return h[:]
}

// GenerarToken crea un token HMAC expirable.
func GenerarToken() string {
	exp := time.Now().Add(24 * time.Hour).Unix()
	payload := "admin:" + strconv.FormatInt(exp, 10)
	mac := hmac.New(sha256.New, secretKey())
	mac.Write([]byte(payload))
	sig := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + hex.EncodeToString(sig)
}

func verificarToken(token string) bool {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return false
	}
	payloadB, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	payload := string(payloadB)

	// Parse "role:unix_timestamp"
	idx := strings.LastIndexByte(payload, ':')
	if idx < 0 {
		return false
	}
	exp, err := strconv.ParseInt(payload[idx+1:], 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return false
	}

	sigB, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, secretKey())
	mac.Write([]byte(payload))
	return hmac.Equal(sigB, mac.Sum(nil))
}

// Middleware protege rutas que requieren autenticación admin.
// GET requests y /api/auth/login quedan exentas.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/auth/login" ||
			r.URL.Path == "/health" ||
			r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"se requiere autenticación"}`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if !verificarToken(token) {
			http.Error(w, `{"error":"token inválido o expirado"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
