package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/grvbrk/async0_server/internal/models"
	"github.com/grvbrk/async0_server/internal/utils"
	"github.com/rbcervilla/redisstore/v9"
)

var allowedOrigins = []string{
	"http://localhost:3000",
	"http://localhost:3001",
}

type contextKey string

const UserContextKey contextKey = "user"
const AdminContextKey contextKey = "admin"

type MiddlewareHandler struct {
	SessionStore *redisstore.RedisStore
	logger       *log.Logger
}

func NewMiddlewareHandler(logger *log.Logger, store *redisstore.RedisStore) *MiddlewareHandler {
	return &MiddlewareHandler{
		logger:       logger,
		SessionStore: store,
	}
}

func (mh *MiddlewareHandler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := mh.SessionStore.Get(r, "session")
		if err != nil {
			mh.logger.Println("Error getting session in auth middleware:", err)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Not Authorized"})
			return
		}

		if session.IsNew {
			mh.logger.Println("New session found in auth middleware (not authenticated)")
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Not Authorized"})
			return
		}

		userEmail, emailOk := session.Values["user_email"].(string)
		userIDStr, idOk := session.Values["user_id"].(string)

		if !emailOk || !idOk || userEmail == "" || userIDStr == "" {
			mh.logger.Println("Invalid or missing user data in session")
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Not Authorized"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			mh.logger.Println("Invalid user ID format in session:", err)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Not Authorized"})
			return
		}

		user := &models.User{
			ID:    userID,
			Email: userEmail,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mh *MiddlewareHandler) AuthenticateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := mh.SessionStore.Get(r, "session")
		if err != nil {
			mh.logger.Println("Error getting admin session in auth middleware:", err)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Admin access required"})
			return
		}

		if session.IsNew {
			mh.logger.Println("New admin session found in auth middleware (not authenticated)")
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Admin access required"})
			return
		}

		adminEmail, emailOk := session.Values["admin_email"].(string)
		adminIDStr, idOk := session.Values["admin_id"].(string)

		if !emailOk || !idOk || adminEmail == "" || adminIDStr == "" {
			mh.logger.Println("Invalid or missing admin data in session")
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Admin access required"})
			return
		}

		adminID, err := uuid.Parse(adminIDStr)
		if err != nil {
			mh.logger.Println("Invalid admin ID format in session:", err)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Admin access required"})
			return
		}

		admin := &models.User{
			ID:    adminID,
			Email: adminEmail,
			Role:  "ADMIN",
		}

		ctx := context.WithValue(r.Context(), UserContextKey, admin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mh *MiddlewareHandler) Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && !isOriginAllowed(origin) {
			mh.logger.Printf("Origin not allowed: %s", origin)
			utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "Origin not allowed"})
			return
		}

		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (mh *MiddlewareHandler) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		mh.logger.Printf("Request: %s %s | Origin: %s",
			r.Method, r.URL.Path, origin)

		next.ServeHTTP(w, r)
	})
}

func (mh *MiddlewareHandler) Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}

func isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}

func GetUserFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	return user, ok
}

func GetAdminFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(AdminContextKey).(*models.User)
	return user, ok
}
