package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/HudYuSa/sqlc-crud-api-gin/config"
	"github.com/HudYuSa/sqlc-crud-api-gin/controllers"
	dbConn "github.com/HudYuSa/sqlc-crud-api-gin/db/sqlc"
	"github.com/HudYuSa/sqlc-crud-api-gin/routes"
	"github.com/HudYuSa/sqlc-crud-api-gin/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
	db     *dbConn.Queries
	ctx    context.Context

	AuthController controllers.AuthController
	UserController controllers.UserController

	AuthRoutes routes.AuthRoutes
	UserRoutes routes.UserRoutes
)

func init() {
	ctx = context.TODO()
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("couldn not load config: %v", err)
	}

	conn, err := sql.Open(config.PostgreDriver, config.PostgresSource)
	if err != nil {
		log.Fatalf("could not connect to postgres database: %v", err)
	}

	db = dbConn.New(conn)

	fmt.Println("PostgreSQL connected successfully...")

	AuthController = controllers.NewAuthController(db)
	UserController = controllers.NewUserController(db)

	AuthRoutes = routes.NewAuthRoutes(AuthController, db)
	UserRoutes = routes.NewUserRoutes(UserController, db)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Origin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")

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
