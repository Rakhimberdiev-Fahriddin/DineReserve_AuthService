package handler

import (
	"auth-service/auth/token"
	pb "auth-service/generated/auth_service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param Register body auth_service.RegisterRequest true "User Registration"
// @Success 201 {object} models.Success
// @Failure 400 {object} models.Errors
// @Router /auth/register [post]
func (h *Handler) RegisterHandler(ctx *gin.Context) {
	h.Logger.Info("Handling RegisterHandler request")

	user := pb.RegisterRequest{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.Logger.Error("Error binding JSON:", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("Error generating hashed password", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	user.Password = string(hashedPassword)

	resp, err := h.UserRepo.CreateUser(&user)
	if err != nil {
		h.Logger.Error("Error register user", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	h.Logger.Info(resp.Message)
	ctx.JSON(http.StatusCreated, gin.H{
		"Message": resp.Message,
	})
}

// LoginHandler handles user login
// @Summary Login a user
// @Description Login a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param Login body auth_service.LoginRequest true "User Login"
// @Success 200 {object} models.Token
// @Failure 400 {object} models.Errors
// @Failure 404 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Router /auth/login [post]
func (h *Handler) LoginHandler(ctx *gin.Context) {
	h.Logger.Info("Handling LoginHandler request")

	user := pb.LoginRequest{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.Logger.Error("Error binding JSON:", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	storedUser, err := h.UserRepo.GetByEmail(user.Email)
	if err != nil {
		h.Logger.Error("Error getting user by email", "error", err.Error())
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		h.Logger.Error("Invalid password", "error", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	accessToken, err := token.GenerateAccessJWT(storedUser)
	if err != nil {
		h.Logger.Error("Error generating access token:", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := token.GenerateRefreshJWT(storedUser)
	if err != nil {
		h.Logger.Error("Error generating refresh token:", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	h.Logger.Info("user login successfully")
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// LogoutUserHandler handles the logout of a user.
// @Summary Logout User
// @Description Logout the authenticated user
// @Tags Auth
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param user-id path string true "User ID"
// @Success 200 {object} auth_service.LogoutResponse
// @Failure 400 {object} string "Bad Request"
// @Router /api/auth/logout/{user-id} [post]
func (h *Handler) LogoutUserHandler(ctx *gin.Context) {
	h.Logger.Info("Handling LogoutUserHandler request")

	id := ctx.Param("user-id")

	resp, err := h.UserRepo.LogoutUser(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetUserProfileHandler retrieves the user profile.
// @Summary Get User Profile
// @Description Get profile of the authenticated user
// @Tags Auth
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} auth_service.GetUserProfileResponse
// @Failure 400 {object} string "Bad Request"
// @Router /api/auth/profile/{username} [get]
func (h *Handler) GetUserProfileHandler(ctx *gin.Context) {
	h.Logger.Info("Handling GetUserProfileHandler request")

	username := ctx.Param("username")

	resp, err := h.UserRepo.GetUserProfile(username)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateUserProfile updates the user profile.
// @Summary Update User Profile
// @Description Update the profile of the authenticated user
// @Tags Auth
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param user-id path string true "User ID"
// @Param profile body auth_service.UpdateUserProfileRequest true "Profile"
// @Success 200 {object} auth_service.UpdateUserProfileResponse
// @Failure 400 {object} string "Bad Request"
// @Router /api/auth/profile/{user-id} [put]
func (h *Handler) UpdateUserProfile(ctx *gin.Context) {
	h.Logger.Info("Handling UpdateUserProfile request")

	id := ctx.Param("user-id")

	profile := pb.UpdateUserProfileRequest{}
	if err := ctx.ShouldBindJSON(&profile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	profile.UserId = id

	resp, err := h.UserRepo.UpdateUserProfile(&profile)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// RefreshToken generates a new access token using a refresh token
// @Summary Refresh access token
// @Description Refresh the access token using the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Refresh token"
// @Success 200 {object} models.Request
// @Failure 400 {object} models.Errors
// @Failure 401 {object} models.Errors
// @Failure 500 {object} models.Errors
// @Security ApiKeyAuth
// @Router /auth/refresh_token [get]
func (h *Handler) RefreshToken(c *gin.Context) {
	h.Logger.Info("Handling RefreshToken request")

	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	claims, err := token.ExtractClaim(refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		return
	}

	newAccessToken, err := token.GenerateAccessJWT(&pb.LoginResponse{
		UserId:   claims.UserId,
		Username: claims.Username,
		Email:    claims.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
