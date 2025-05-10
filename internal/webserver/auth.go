package webServer

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/runway"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setAuthRoutes(router *gin.Engine) {
	router.GET("/auth/hl", GHLVL)
	router.GET("/auth/hl/update", GHLVL_Update)
}

func GHLVL(c *gin.Context) {
	code := c.Query("code")

	locId, err := runway.TradeAccessCode(code)
	if err != nil {
		panic(err)
	}

	c.Redirect(http.StatusTemporaryRedirect, config.Confs.Settings.SrvDomain+"/app?oper=locSettings&locationId="+locId)
}

func GHLVL_Update(c *gin.Context) {
	code := c.Query("code")

	locId, err := runway.TradeAccessCodeAndSaveTokens(code)
	if err != nil {
		panic(err)
	}

	c.Redirect(http.StatusTemporaryRedirect, config.Confs.Settings.SrvDomain+"/app?oper=locSettings&locationId="+locId)
}
