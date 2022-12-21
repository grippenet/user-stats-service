package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grippenet/user-stats-service/internal"
	"github.com/grippenet/user-stats-service/pkg/db"
	"github.com/grippenet/user-stats-service/pkg/stats"
	"github.com/grippenet/user-stats-service/pkg/types"
	"golang.org/x/crypto/bcrypt"
)

var authPassword string
var authHander func(c *gin.Context) = NoOp

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func basicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth && user == "stats" && CheckPasswordHash(password, authPassword) {
		return
	}
	c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	c.AbortWithStatus(http.StatusUnauthorized)
}

func NoOp(c *gin.Context) {
	// Do Nothing
}

func main() {

	authPassword = os.Getenv("AUTH_PASSWORD")

	if authPassword != "" {
		authHander = basicAuth
	}

	config := internal.GetUserDBConfig()

	dbService := db.NewUserDBService(config)

	statsService := stats.NewStatService(dbService)

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("accounts/:instance", authHander, func(c *gin.Context) {

		instance := c.Param("instance")

		filter := types.StatFilter{}

		fromQuery := c.Query("from")
		if fromQuery != "" {
			i, err := strconv.ParseInt(fromQuery, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(401, map[string]string{"error": "From must be a numeric"})
				return
			}
			filter.From = i
		}

		untilQuery := c.Query("until")
		if untilQuery != "" {
			i, err := strconv.ParseInt(untilQuery, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(401, map[string]string{"error": "until must be a numeric"})
				return
			}
			filter.Until = i
		}
		counters, err := statsService.Fetch(instance, filter)
		if err != nil {
			c.AbortWithStatusJSON(401, map[string]string{"error": fmt.Sprintf("error during fetch: %s", err)})
			return
		}
		c.JSON(200, counters)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
