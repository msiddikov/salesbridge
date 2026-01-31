package svc_googleads

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/packages/googleads"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
)

// GetAuthURL returns a profile-scoped OAuth URL. Requires auth to bind to the caller's profile.
func GetAuthURL(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	state := fmt.Sprintf("profile:%d", user.ProfileID)
	url, err := Svc.AuthURL(state)
	lvn.GinErr(c, 400, err, "unable to build oauth url")
	if err != nil {
		return
	}

	c.Data(lvn.Res(200, gin.H{
		"url":   url,
		"state": state,
	}, "OK"))
}

func parseProfileID(state string) (uint, error) {
	if !strings.HasPrefix(state, "profile:") {
		return 0, fmt.Errorf("unexpected state format")
	}
	raw := strings.TrimPrefix(state, "profile:")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse profile id: %w", err)
	}
	return uint(id), nil
}

// OAuthCallback exchanges the code and stores the connection.
// On GET, it redirects to /oauth/callback?google_ads_connected=true (or error=<slug>).
// On POST (JSON), it returns the stored connection.
func OAuthCallback(c *gin.Context) {
	method := c.Request.Method
	var code, state string

	if method == http.MethodPost {
		var body struct {
			Code  string `json:"code"`
			State string `json:"state"`
		}
		if err := c.BindJSON(&body); err != nil {
			lvn.GinErr(c, 400, err, "invalid payload")
			return
		}
		code = body.Code
		state = body.State
	} else {
		code = c.Query("code")
		state = c.Query("state")
	}

	profileID, err := parseProfileID(state)
	if err != nil {
		handleCallbackError(c, method, "state_mismatch")
		return
	}

	if code == "" {
		handleCallbackError(c, method, "missing_code")
		return
	}

	conn, err := Svc.ExchangeCode(c, code, googleads.Connection{
		ProfileID: profileID,
	})
	if err != nil {
		handleCallbackError(c, method, "exchange_failed")
		return
	}

	if method == http.MethodPost {
		c.Data(lvn.Res(200, conn, "Connection saved"))
		return
	}

	redirectURL := config.Confs.Settings.AppDomain + "/oauth/callback?google_ads_connected=true"
	c.Redirect(http.StatusFound, redirectURL)
}

func handleCallbackError(c *gin.Context, method, code string) {
	if method == http.MethodPost {
		lvn.GinErr(c, 400, errors.New(code), code)
		return
	}
	c.Redirect(http.StatusFound, "/oauth/callback?error="+code)
}
