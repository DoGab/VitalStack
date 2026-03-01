package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	// UserIDContextKey is the key used to store the user ID in the context
	UserIDContextKey contextKey = "user_id"
)

// DevAuthMiddleware injects a dummy Dev User UUID into the request context if DevMode is true.
// This acts as a placeholder for a real authentication middleware (like Supabase Auth).
func DevAuthMiddleware(devMode bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if devMode {
			// For local development, inject the seeded dummy user ID
			ctx := context.WithValue(c.Request.Context(), UserIDContextKey, "00000000-0000-0000-0000-000000000000")
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// GetUserFromContext extracts the UserID string from the provided context.
// It returns the user ID and a boolean indicating whether it was found and valid.
func GetUserFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(UserIDContextKey).(string)
	if !ok || uid == "" {
		return "", false
	}
	return uid, true
}
