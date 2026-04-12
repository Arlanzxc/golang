package v1

import (
	"net/http"
	"time"

	"GOLANG/internal/entity"
	"GOLANG/internal/usecase"
	"GOLANG/utils"

	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	t usecase.UserInterface
}

func NewUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface) {
	r := &userRoutes{t}

	h := handler.Group("/users")

	h.Use(utils.RateLimitMiddleware(3, 20*time.Second))
	h.POST("/", r.RegisterUser)
	h.POST("/login", r.LoginUser)

	protected := h.Group("/")
	protected.Use(utils.JWTAuthMiddleware())
	{
		protected.GET("/me", r.GetMe)

		admin := protected.Group("/")
		admin.Use(utils.RoleMiddleware("admin"))
		{
			admin.PATCH("/promote/:id", r.PromoteUser)
		}
	}
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var input entity.CreateUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := entity.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role,
	}

	createdUser, sessionID, err := r.t.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully",
		"session_id": sessionID,
		"user":       createdUser,
	})
}

func (r *userRoutes) LoginUser(c *gin.Context) {
	var input entity.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.t.LoginUser(&input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := r.t.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func (r *userRoutes) PromoteUser(c *gin.Context) {
	targetID := c.Param("id")

	if err := r.t.PromoteUser(targetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}