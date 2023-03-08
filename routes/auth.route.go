package routes

import (
	"github.com/HudYuSa/sqlc-crud-api-gin/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type authRoutes struct {
	authController controllers.AuthController
}

func NewAuthRoutes(authController controllers.AuthController) AuthRoutes {
	return &authRoutes{
		authController: authController,
	}
}

func (ar *authRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", ar.authController.SignUpUser)
}
