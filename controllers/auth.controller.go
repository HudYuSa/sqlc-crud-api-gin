package controllers

import (
	"context"
	"net/http"
	"time"

	db "github.com/HudYuSa/sqlc-crud-api-gin/db/sqlc"
	"github.com/HudYuSa/sqlc-crud-api-gin/models"
	"github.com/HudYuSa/sqlc-crud-api-gin/utils"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	SignUpUser(ctx *gin.Context)
}

type authController struct {
	db  *db.Queries
	ctx context.Context
}

func NewAuthController(db *db.Queries, ctx context.Context) AuthController {
	return &authController{db, ctx}
}

func (ac *authController) SignUpUser(ctx *gin.Context) {
	credentials := models.SignUpInput{}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	hashedPassword := utils.HashPassword(credentials.Password)

	var role string

	if _, ok := ctx.GetQuery("admin"); ok {
		role = "admin"
	} else {
		role = "user"
	}

	args := &db.CreateUserParams{
		Name:      credentials.Name,
		Email:     credentials.Email,
		Password:  hashedPassword,
		Photo:     credentials.Photo,
		Verified:  true,
		Role:      role,
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

	userResponse := models.FilteredResponse(user)

	ctx.JSON(http.StatusCreated, utils.WebResponse{
		Status: "success",
		Data:   gin.H{"user": userResponse},
	})

}
