package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bnquan27/QMQShop/backend/database"
	"github.com/bnquan27/QMQShop/backend/models"
)

type contextKey string

const UserKey contextKey = "user"

// CORS adds CORS headers to allow frontend access
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging logs each request
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// GetUserFromRequest extracts user from Authorization header
func GetUserFromRequest(r *http.Request) *models.User {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		return nil
	}

	session, err := database.GetSessionByToken(token)
	if err != nil || session == nil {
		return nil
	}

	user, err := database.GetUserByID(session.UserID)
	if err != nil || user == nil {
		return nil
	}
	return user
}

// RequireAuth ensures a valid session exists, injects user into context
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromRequest(r)
		if user == nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, user)
		next(w, r.WithContext(ctx))
	}
}

// RequireAdmin ensures user is admin
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(*models.User)
		if user.Role != "admin" {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

// JSON helper
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err, ok := data.(error); ok {
			w.Write([]byte(`{"error":"` + err.Error() + `"}`))
			return
		}
		if bs, ok := data.([]byte); ok {
			w.Write(bs)
			return
		}
			json.NewEncoder(w).Encode(data)
	}
}

// ParseJSON reads JSON body
func ParseJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
