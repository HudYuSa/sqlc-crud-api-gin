package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HudYuSa/sqlc-crud-api-gin/config"
	db "github.com/HudYuSa/sqlc-crud-api-gin/db/sqlc"
	"github.com/HudYuSa/sqlc-crud-api-gin/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeserializeUser(db *db.Queries) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token from cookie or header
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")

		fields := strings.Fields(authorizationHeader)
		// if auth header is not empty
		if len(fields) == 2 && fields[0] == "Bearer" {
			accessToken = fields[1]
			// else if cookie is not returning error
		} else if err == nil {
			accessToken = cookie
		}

		// send error: if there's no token from either cookie or auth header then send error response
		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, utils.WebResponse{
				Status:  "fail",
				Message: "Cannot access this resource",
			})
			return
		}

		config, _ := config.LoadConfig(".")

		// validate the token and get the user id from the sub/subject
		sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.WebResponse{
				Status:  "fail",
				Message: err.Error(),
			})
			return
		}

		user, err := db.GetUserById(context.TODO(), uuid.MustParse(fmt.Sprint(sub)))

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.WebResponse{
				Status:  "fail",
				Message: "The user belonging to this token no longer exists",
			})
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
