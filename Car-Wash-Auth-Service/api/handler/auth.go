package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/Car-Wash/Car-Wash-Auth-Service/docs"
	"github.com/Car-Wash/Car-Wash-Auth-Service/genproto/user"

	"github.com/Car-Wash/Car-Wash-Auth-Service/api/token"
	"github.com/Car-Wash/Car-Wash-Auth-Service/genproto/auth"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Registers a new user and returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Register Request"
// @Success 200 {object} auth.RegisterResponse
// @Failure 400 {object} string "Invalid request payload"
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	req := &auth.RegisterRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid struct"})
	}

	if len(req.Password) <= 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	password, err := token.HashPassword(req.Password)
	if err != nil {
		c.JSON(400, "Error Hash Password")
		return
	}
	req.Password = password

	token := token.GenereteJWTToken(req)
	req.Token = token.RefreshToken
	res, err := json.Marshal(req)
	if err != nil {
		c.JSON(400, err)
		return
	}
	err = h.Kafka.ProduceMessages("auth", res)
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, "Success")
}

// Login godoc
// @Summary Login with existing credentials
// @Description Logs in a user and returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login Request"
// @Success 200 {object} auth.LoginResponse
// @Failure 400 {object} string "Invalid request payload"
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	req := auth.LoginRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid struct"})
		slog.Info(err.Error())
	}

	res, err := h.Auth.Login(c, &req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		slog.Info(err.Error())
	}
	_, err = token.ExtractClaim(res.Token)
	if err != nil {
		slog.Info(err.Error())
		c.JSON(400, err)
	}

	c.JSON(200, res)
}

// Logout godoc
// @Summary Logout a user
// @Description Logs out a user by invalidating the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.LogoutRequest true "Logout Request"
// @Success 200 {object} auth.LogoutResponse
// @Failure 401 {object} string "Invalid request payload"
// @Router /auth/logout [post]'
func (h *Handler) Logout(c *gin.Context) {
	req := &auth.LogoutRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	res, err := h.Auth.Logout(c, req)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(200, res)
}

// ForgotPassword godoc
// @Summary Initiate forgot password flow
// @Description Initiates the forgot password flow by sending a code to the user's email
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.ForgotPasswordRequest true "Email"
// @Success 200 {object} string "Success"
// @Failure 400 {object} string "Invalid request payload"
// @Router /auth/forgot [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	req := &auth.ForgotPasswordRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	res, err := h.User.GetProfile(c, &user.GetProfileRequest{Email: req.Email})
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if res == nil {
		c.JSON(400, gin.H{"error": "Not found"})
		return
	}

	// Random 6 xonali raqam yaratish
	RandomCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	// Redis ga email va code ni saqlash
	status := h.Redis.Set(c, req.Email, RandomCode, 10*time.Minute)
	if status.Err() != nil {
		slog.Info(status.Err().Error())
		c.JSON(500, gin.H{"error": "Failed to save code"})
		return
	}

	c.JSON(200, RandomCode)
}

// ResetPassword godoc
// @Summary Reset user's password
// @Description Resets the user's password using the provided code and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} auth.ResetPasswordResponse
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Error Reset"
// @Router /auth/reset [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	req := &auth.ResetPasswordRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	// Redis dan kodni olish
	codeRes := h.Redis.Get(c, req.Email)
	if codeRes.Err() != nil {
		slog.Info(codeRes.Err().Error())
		c.JSON(400, gin.H{"error": "Failed to retrieve code from Redis"})
		return
	}

	code, err := codeRes.Result()
	if err != nil {
		slog.Info(err.Error())
		c.JSON(400, gin.H{"error": "Failed to retrieve code from Redis"})
		return
	}
	// Kodni tekshirish yani structdan keladigon randomcode bilan
	if code != req.RandomnNum {
		c.JSON(400, gin.H{"error": "Invalid code. Please check your email and try again"})
		return
	}
	// Parolni tiklash
	res, err := h.Auth.ResetPassword(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, res)
}
