package handler

import (
	"encoding/json"
	"net/http"

	"github.com/exam-5/Car-Wash-Auth-Service/genproto/user"
	"github.com/gin-gonic/gin"
)

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieves a user's profile by ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string false "User ID"
// @Param email query string false "User Email"
// @Success 200 {object} user.GetProfileResponse
// @Failure 400 {object} string "Invalid request payload"
// @Router /user/profile/{id} [get]
func (h *Handler) GetProfile(c *gin.Context) {
	id := c.Query("id")
	email := c.Query("email")
	req := &user.GetProfileRequest{Id: id, Email: email}
	res, err := h.User.GetProfile(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Updates a user's profile
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} user.UpdateProfileResponse
// @Failure 400 {object} string "Invalid request payload"
// @Router /user/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	req := &user.UpdateProfileRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid struct"})
		return
	}
	_, err := h.User.UpdateProfile(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Success updated User")
}

// ChangePassword godoc
// @Summary Change user password
// @Description Changes a user's password
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} user.ChangePasswordResponse
// @Failure 400 {object} string "Invalid request payload"
// @Router /user/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	req := &user.ChangePasswordRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid struct"})
		return
	}

	if len(req.NewPassword) <= 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	res, err := json.Marshal(req)
	if err != nil {
		c.JSON(400, "Error Marshal Struct")
	}
	err = h.Kafka.ProduceMessages("user", res)
	if err != nil {
		c.JSON(400, "Error change")
	}
	c.JSON(http.StatusOK, "Success")
}

// @Summary Get all users
// @Description Retrieves all users
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} user.GetAllUsersResponse
// @Failure 400 {object} string "Invalid request payload"
// @Router /user/all [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	req := &user.GetAllUsersRequest{}
	res, err := h.User.GetAllUsers(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
