package auth

import (
	"client-runaway-zenoti/internal/db/models"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
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

	profile := models.Profile{
		Name: profileName,
	}
	res := models.DB.Create(&profile)
	if res.Error != nil {
		c.Data(lvn.Res(400, "", "Unable to create profile"))
		return
	}

	user := models.User{
		Email:     email,
		Password:  password,
		ProfileID: profile.ID,
	}

	result = models.DB.Create(&user)

	if result.Error != nil {
		lvn.GinErr(c, 400, result.Error, "Unable to create user")
		return
	}

	profile.OwnerID = user.ID
	models.DB.Save(&profile)
	c.Data(lvn.Res(200, result, "User registered successfully"))
}
