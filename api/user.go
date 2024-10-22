package api

import (
	"errors"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/jwt"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
)

type UserHandler struct {
	UserService service.UserService
	JwtService  jwt.Jwt
	logger      *slog.Logger
}

func NewUserHandler(userRepo service.UserService, jwt jwt.Jwt) *UserHandler {
	return &UserHandler{
		UserService: userRepo,
		JwtService:  jwt,
		logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (h *UserHandler) InsertUser(c *gin.Context) {
	type body struct {
		Email     string `json:"email" binding:"required"`
		Username  string `json:"username" binding:"required"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password" binding:"required"`
	}
	var user body

	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	toInsert := domain.User{
		Email:     &user.Email,
		Username:  &user.Username,
		FirstName: &user.FirstName,
		LastName:  &user.LastName,
		Password:  &user.Password,
	}

	tx := h.UserService.Db.Begin()

	insertedUser, err := h.UserService.Repo.Insert(tx, toInsert)
	if err != nil {
		h.logger.Error(err.Error())
		if db := tx.Rollback(); db.Error != nil {
			h.logger.Error(db.Error.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while rollback tx"})
			return
		}
		var contextErr *domain.BadRequestError
		if errors.As(err, &contextErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while inserting user"})
		return
	}

	token, err := h.JwtService.GenerateToken(insertedUser.ID)
	if err != nil {
		h.logger.Error(err.Error())
		if db := tx.Rollback(); db.Error != nil {
			h.logger.Error(db.Error.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while rollback tx"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while generating user token"})
		return
	}

	if db := tx.Commit(); db.Error != nil {
		h.logger.Error(db.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while commiting tx"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User created successfully",
		"username": insertedUser.Username,
		"email":    insertedUser.Email,
		"id":       insertedUser.ID,
		"token":    token,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	type body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var user body

	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.UserService.Repo.FindByEmailAndPassword(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while finding user"})
		return
	}

	token, err := h.JwtService.GenerateToken(res.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while generating JWToken"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged successfully",
		"user":    res.Username,
		"token":   token,
	})
}
