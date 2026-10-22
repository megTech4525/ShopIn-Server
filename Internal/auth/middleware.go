package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ctxAccountID = "auth.account_id"
	ctxRole      = "auth.role"
)

// authenticat function rejects any request that does nt carry a valid token
func Authenticate(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			unauthorized(c, "missing or malformed token")
			return
		}

		claims, err := ParseToken(secret, tokenString)
		if err != nil {
			unauthorized(c, "Invalid or expired token")
			return
		}

		c.Set(ctxAccountID, claims.Subject)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}

}

//Require role allows omly the listed roles through. it must run after Authenticate

func RequireRole(allowed ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(ctxRole)
		if !ok {
			//Aunthenticate did not run, so we do not know who this is.
			unauthorized(c, "authentication required")
			return

		}
		current, _ := role.(Role)
		for _, a := range allowed {
			if current == a {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"ERROR": "you do not have access to this resource"})

	}
}

// ACCOUNTID FROM RETURNS THE AUTHENTICATED CALLER'S ACCOUNT ID("" IF NONE).
func AccountIDFrom(c *gin.Context) string {
	return c.GetString(ctxAccountID)
}

// RoleFrom returns the authenticated caller's role ("" if none)

func RoleForm(c *gin.Context) Role {
	v, _ := c.Get(ctxRole)
	r, _ := v.(Role)
	return r
}

//bearerToken Extracts the token from an "Authorization: Bearer<token>" header.

func bearerToken(header string) (string, bool) {
	const prefix = "bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}
	token := strings.TrimSpace(header[len(prefix):])
	return token, token != ""
}

func unauthorized(c *gin.Context, msg string) {
	c.Header("WWW-Authenticate", "Bearer")
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": msg})
}
