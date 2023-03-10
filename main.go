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

	ctx := gin.Context{}

	AuthController = controllers.NewAuthController(db, &ctx)
	AuthRoutes = routes.NewAuthRoutes(AuthController)

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

	AuthRoutes.SetupRoutes(router)

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, utils.WebResponse{
			Status:  "fail",
			Message: fmt.Sprintf("Route %s not found", ctx.Request.URL),
		})
	})

	log.Fatal(server.Run(":" + config.Port))
}
