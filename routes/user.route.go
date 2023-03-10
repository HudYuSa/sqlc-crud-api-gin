package routes

import (
	"github.com/HudYuSa/sqlc-crud-api-gin/controllers"
	db "github.com/HudYuSa/sqlc-crud-api-gin/db/sqlc"
	"github.com/HudYuSa/sqlc-crud-api-gin/middleware"
	"github.com/gin-gonic/gin"
)

type UserRoutes interface {
	SetupRoutes(rg *gin.RouterGroup)
}

type userRoutes struct {
	userController controllers.UserController
	db             *db.Queries
}

func NewUserRoutes(userController controllers.UserController, db *db.Queries) UserRoutes {
	return &userRoutes{
		userController: userController,
		db:             db,
	}
}

func (ur *userRoutes) SetupRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/users")

	router.GET("/me", middleware.DeserializeUser(ur.db), ur.userController.GetMe)
}
