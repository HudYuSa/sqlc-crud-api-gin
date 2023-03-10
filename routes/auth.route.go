package routes

import (
	"github.com/HudYuSa/sqlc-crud-api-gin/controllers"
	db "github.com/HudYuSa/sqlc-crud-api-gin/db/sqlc"
	"github.com/HudYuSa/sqlc-crud-api-gin/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type authRoutes struct {
	authController controllers.AuthController
	db             *db.Queries
}

func NewAuthRoutes(authController controllers.AuthController, db *db.Queries) AuthRoutes {
	return &authRoutes{
		authController: authController,
		db:             db,
	}
}

func (ar *authRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", ar.authController.SignUpUser)
	router.POST("login", ar.authController.SignInUser)
	router.GET("/refresh", ar.authController.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(ar.db), ar.authController.LogoutUser)
}
