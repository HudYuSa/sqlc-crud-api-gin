package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/HudYuSa/sqlc-crud-api-gin/config"
	"github.com/HudYuSa/sqlc-crud-api-gin/controllers"
	dbConn "github.com/HudYuSa/sqlc-crud-api-gin/db/sqlc"
	"github.com/HudYuSa/sqlc-crud-api-gin/routes"
	"github.com/HudYuSa/sqlc-crud-api-gin/utils"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

var (
	server *gin.Engine
	db     *dbConn.Queries

	AuthController controllers.AuthController
	AuthRoutes     routes.AuthRoutes

	UserController controllers.UserController
	UserRoutes     routes.UserRoutes
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	conn, err := sql.Open(config.PostgreDriver, config.PostgresSource)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}

	db = dbConn.New(conn)

	fmt.Println("PostgreSQL connected successfully...")

	AuthController = controllers.NewAuthController(db)
	AuthRoutes = routes.NewAuthRoutes(AuthController, db)

	UserController = controllers.NewUserController(db)
	UserRoutes = routes.NewUserRoutes(UserController, db)

	server = gin.Default()

}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	router := server.Group("/api")

	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, utils.WebResponse{
			Status:  "success",
			Message: "Welcome to Golang with PostgreSQL",
		})
	})

	// setup the routes
	AuthRoutes.SetupRoutes(router)
	UserRoutes.SetupRoutes(router)

	//  if there's no route matching found
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, utils.WebResponse{
			Status:  "fail",
			Message: fmt.Sprintf("Route %s not found", ctx.Request.URL),
		})
	})

	log.Fatal(server.Run(":" + config.Port))
}
