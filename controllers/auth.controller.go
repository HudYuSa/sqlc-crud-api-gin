package controllers

import (
	"net/http"
	"time"

	db "github.com/HudYuSa/sqlc-crud-api-gin/db/sqlc"
	"github.com/HudYuSa/sqlc-crud-api-gin/utils"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	SignUpUser(ctx *gin.Context)
}

type authController struct {
	db *db.Queries
}

func NewAuthController(db *db.Queries) AuthController {
	return &authController{db}
}

func (ac *authController) SignUpUser(ctx *gin.Context) {
	credentials := db.User{}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	hashedPassword := utils.HashPassword(credentials.Password)

	args := &db.CreateUserParams{
		Name:      credentials.Name,
		Email:     credentials.Email,
		Password:  hashedPassword,
		Photo:     "default.jpeg",
		Verified:  true,
		Role:      "user",
		UpdatedAt: time.Now(),
	}

	user, err := ac.db.CreateUser(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, utils.WebResponse{
		Status: "success",
		Data:   gin.H{"user": user},
	})

}
