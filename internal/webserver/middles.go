package webServer

import (
	"fmt"
	"strings"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var corsMiddleware = cors.New(cors.Config{
	//AllowOrigins:     allowedOrigins,
	AllowMethods:     []string{"GET"},
	AllowHeaders:     []string{"Origin", "Content-Length", "Content-type"},
	ExposeHeaders:    []string{"Content-Length", "Content-type"},
	AllowCredentials: true,
	AllowOriginFunc:  checkOrigin,
	MaxAge:           12 * time.Hour,
})

func checkOrigin(origin string) bool {

	return true

}

var allowedOrigins = []string{
	"http://localhost:3000",
	"https://app.clientrunway.com",
}

func recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if strings.Contains(fmt.Sprint(err), "lvn.GinErr panic") {
				return
			}
			c.JSON(500, lvn.Response(nil, fmt.Sprint(err), false))
			fmt.Println(err)
		}
	}()
	c.Next()
}
