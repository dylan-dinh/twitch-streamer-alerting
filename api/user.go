package api

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/domain"
	"github.com/dylan-dinh/twitch-streamer-alerting/internal/repository"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
)

type UserHandler struct {
	UserRepo repository.UserRepo
	logger   *slog.Logger
}

func NewUserHandler(userRepo repository.UserRepo) *UserHandler {
	return &UserHandler{
		UserRepo: userRepo,
		logger:   slog.New(slog.NewTextHandler(os.Stdout, nil)),
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

	err := h.UserRepo.Insert(toInsert)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while inserting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"user":    user.Username,
	})
}
