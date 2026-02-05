package auth

import (
	"client-runaway-zenoti/internal/db/models"
	"fmt"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Auth(c *gin.Context) {

	// Assuming the header is in the format "Basic base64encodedstring"
	email, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Data(lvn.Res(400, "", "No login info"))
		c.Abort()
		return
	}

	user := models.User{}
	result := models.DB.Where("email = ? AND password = ?", email, password).Preload("Profile").First(&user)
	if result.Error != nil {
		c.Data(lvn.Res(400, "", "Incorrect email or password"))
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}

func Login(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	c.Data(lvn.Res(200, user, ""))
}

func Register(c *gin.Context) {

	email := c.PostForm("email")
	password := c.PostForm("password")
	profileName := c.PostForm("profile_name")

	// check if user already exists
	var existingUser models.User
	result := models.DB.Where("email = ?", email).First(&existingUser)
	if result.Error == nil {
		c.Data(lvn.Res(400, "", "User already exists"))
		return
	}

	tx := models.DB.Begin()
	if tx.Error != nil {
		lvn.GinErr(c, 500, tx.Error, "Unable to start registration")
		return
	}

	profile := models.Profile{
		Name: profileName,
	}
	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		c.Data(lvn.Res(400, "", "Unable to create profile"))
		return
	}

	user := models.User{
		Email:     email,
		Password:  password,
		ProfileID: profile.ID,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		lvn.GinErr(c, 400, err, "Unable to create user")
		return
	}

	if err := tx.Model(&profile).Update("owner_id", user.ID).Error; err != nil {
		tx.Rollback()
		lvn.GinErr(c, 400, err, "Unable to update profile owner")
		return
	}

	if err := createInternalMCPKey(tx, profile.ID); err != nil {
		tx.Rollback()
		lvn.GinErr(c, 500, err, "Unable to create internal MCP key")
		return
	}

	if err := tx.Commit().Error; err != nil {
		lvn.GinErr(c, 500, err, "Unable to complete registration")
		return
	}

	c.Data(lvn.Res(200, user, "User registered successfully"))
}

func createInternalMCPKey(tx *gorm.DB, profileID uint) error {
	if profileID == 0 {
		return nil
	}

	plainKey, keyHash, keyPrefix, err := models.GenerateMCPApiKey()
	if err != nil {
		return err
	}

	key := models.MCPApiKey{
		Name:       fmt.Sprintf("internal-profile-%d", profileID),
		PlainKey:   plainKey,
		KeyHash:    keyHash,
		KeyPrefix:  keyPrefix,
		ProfileID:  profileID,
		IsActive:   true,
		IsInternal: true,
	}

	return tx.Create(&key).Error
}
