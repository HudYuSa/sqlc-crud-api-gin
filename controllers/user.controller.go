package controllers

import (
	"net/http"

	db "github.com/HudYuSa/sqlc-crud-api-gin/db/sqlc"
	"github.com/HudYuSa/sqlc-crud-api-gin/models"
	"github.com/HudYuSa/sqlc-crud-api-gin/utils"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetMe(ctx *gin.Context)
}

type userController struct {
	db *db.Queries
}

func NewUserController(db *db.Queries) UserController {
	return &userController{
		db: db,
	}
}

func (uc *userController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(db.User)

	userResponse := models.FilteredResponse(currentUser)

	ctx.JSON(http.StatusOK, utils.WebResponse{
		Status: "success",
		Data: gin.H{
			"user": userResponse,
		},
	})
}
