// Package middleware provides Gin middleware for Zitadel-based authentication
// and role-based access control.
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

// Authorizer is the concrete Zitadel authorizer type used throughout the app.
type Authorizer = authorization.Authorizer[*oauth.IntrospectionContext]

// NewAuthorizer creates a Zitadel introspection-based authorizer.
//
//   - domain: Zitadel instance hostname, e.g. "my-org.zitadel.cloud"
//   - port:   non-empty only for local/insecure instances, e.g. "8081"
//   - clientID / clientSecret: credentials of the API application created in Zitadel
func NewAuthorizer(ctx context.Context, domain, port, clientID, clientSecret string) (*Authorizer, error) {
	var z *zitadel.Zitadel
	if port != "" {
		z = zitadel.New(domain, zitadel.WithInsecure(port))
	} else {
		z = zitadel.New(domain)
	}
	return authorization.New(
		ctx,
		z,
		oauth.WithIntrospection[*oauth.IntrospectionContext](
			oauth.ClientIDSecretIntrospectionAuthentication(clientID, clientSecret),
		),
	)
}

// Auth returns a Gin middleware that validates Bearer tokens on all /api/* routes
// using Zitadel token introspection. Routes outside /api/ are passed through
// without any auth check.
//
// When authorizer is nil (e.g. ZITADEL env vars not set), all requests are
// allowed and no auth context is stored — useful during local development.
func Auth(authorizer *Authorizer) gin.HandlerFunc {
	if authorizer == nil {
		slog.Warn("auth middleware: no authorizer configured, all requests are permitted")
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next()
			return
		}

		token := extractBearer(c.Request)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Problem{
				Title:  "Unauthorized",
				Status: http.StatusUnauthorized,
				Detail: "Missing or invalid Authorization header",
			})
			return
		}

		authCtx, err := authorizer.CheckAuthorization(c.Request.Context(), token)
		if err != nil {
			slog.WarnContext(c.Request.Context(), "authorization check failed", "err", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, Problem{
				Title:  "Unauthorized",
				Status: http.StatusUnauthorized,
				Detail: "Invalid or expired token",
			})
			return
		}

		// Store in request context so Huma handlers can access it via
		// authorization.Context[*oauth.IntrospectionContext](ctx).
		ctx := authorization.WithAuthContext(c.Request.Context(), authCtx)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GetAuth retrieves the Zitadel auth context from a Huma handler's context.
// Returns nil when auth is disabled (no authorizer configured).
func GetAuth(ctx context.Context) *oauth.IntrospectionContext {
	return authorization.Context[*oauth.IntrospectionContext](ctx)
}

// UserInfo holds the identity fields extracted from an auth token.
type UserInfo struct {
	ZitadelID string
	Email     string
	Username  string
}

// GetUserInfo extracts identity fields from the request context.
// Returns nil when auth is disabled or no token is present.
func GetUserInfo(ctx context.Context) *UserInfo {
	authCtx := GetAuth(ctx)
	if authCtx == nil {
		return nil
	}
	return &UserInfo{
		ZitadelID: authCtx.UserID(),
		Email:     authCtx.Email,
		Username:  authCtx.Username,
	}
}

func extractBearer(r *http.Request) string {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || token == "" {
		return ""
	}
	return token
}

type Problem struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}
