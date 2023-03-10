package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HudYuSa/sqlc-crud-api-gin/config"
	db "github.com/HudYuSa/sqlc-crud-api-gin/db/sqlc"
	"github.com/HudYuSa/sqlc-crud-api-gin/models"
	"github.com/HudYuSa/sqlc-crud-api-gin/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController interface {
	SignUpUser(ctx *gin.Context)
	SignInUser(ctx *gin.Context)
	RefreshAccessToken(ctx *gin.Context)
	LogoutUser(ctx *gin.Context)
}

type authController struct {
	db *db.Queries
}

func NewAuthController(db *db.Queries) AuthController {
	return &authController{db}
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

func (ac *authController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	user, err := ac.db.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	if err := utils.ComparePassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.WebResponse{
			Status:  "fail",
			Message: "invalid email or password",
		})
		return
	}

	config, _ := config.LoadConfig(".")

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	// set the cookie

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	// sent access token as response
	ctx.JSON(http.StatusOK, utils.WebResponse{
		Status: "success",
		Data: gin.H{
			"access_token": accessToken,
		},
	})
}

func (ac *authController) RefreshAccessToken(ctx *gin.Context) {

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	config, _ := config.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	user, err := ac.db.GetUserById(ctx, uuid.MustParse(fmt.Sprint(sub)))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, utils.WebResponse{
			Status:  "fail",
			Message: "the user belonging to this token no longer exist",
		})
		return
	}

	// create new access token
	accessToken, err := utils.CreateToken(time.Duration(config.AccessTokenExpiresIn), user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, utils.WebResponse{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, utils.WebResponse{
		Status: "success",
		Data: gin.H{
			"access_token": accessToken,
		},
	})
}

func (ac *authController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, utils.WebResponse{
		Status: "success",
	})
}
