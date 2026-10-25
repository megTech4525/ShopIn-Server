package auth

import "github.com/gin-gonic/gin"

//RegisterRoutes attaches every auth endpoint to the router
//signupLimit and LoginLimit are passed in (not created here) so main.go
//decides the limits, and so signup and login can have separate ticket bags.

func (h *Handler) RegisterRoutes(r *gin.Engine, secret []byte, signupLimit, loginLimit gin.HandlerFunc) {
	//Everything under /auth gets Nostore, so responses holding tokens
	//are never cached by browsers or proxies

	authGroup := r.Group("/auth", NoStore())
	// POST / auth/signup/client | shop | rider
	signup := authGroup.Group("/signup", signupLimit)
	{
		signup.POST("/client", h.SignupClient)
		// signup.POST("/client", h.SignupShop)
		// signup.POST("/client", h.SignupRider)
	}

	//POST /auth/login
	authGroup.POST("/login", loginLimit, h.login)

	//Get /auth/me (needs a valid token)
	protected := authGroup.Group("", Authenticate(secret))
	{
		protected.GET("/me", h.Me)
	}

}
