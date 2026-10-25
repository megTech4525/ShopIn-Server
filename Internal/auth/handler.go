package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) SignupClient(c *gin.Context) {
	var reg ClientRegister
	if err := c.ShouldBindJSON(&reg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request"})
		return
	}

	resp, err := h.svc.SignupClient(c.Request.Context(), reg)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) login(c *gin.Context) {
	var log Login
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "login failed"})
		return
	}

	resp, err := h.svc.Login(c.Request.Context(), log)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, resp)
}

func (h *Handler) Me(c *gin.Context) {
	profile, err := h.svc.Me(c.Request.Context(), AccountIDFrom(c), RoleForm(c))
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			//valid token but the account no longer exists
			c.JSON(http.StatusUnauthorized, gin.H{"error": "account no longer exists"})
			return
		}
		c.JSON(http.StatusOK, profile)
	}
}
func writeError(c *gin.Context, err error) {
	switch {
	case
		errors.Is(err, ErrPasswordTooLong),
		errors.Is(err, ErrPasswordTooShort),
		errors.Is(err, ErrWeakPassword):

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	case errors.Is(err, ErrPhoneTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	case errors.Is(err, ErrinvalidPhoneNumber):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrEmailTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	case errors.Is(err, ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

	case errors.Is(err, ErrAccountNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	default:
		log.Printf("internal Error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	}

}
