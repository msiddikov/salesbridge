package db

import (
	"client-runaway-zenoti/internal/config"
	"client-runaway-zenoti/internal/db/models"
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	host := config.Confs.DB.Host
	port := config.Confs.DB.Port
	user := config.Confs.DB.User
	password := config.Confs.DB.Password
	dbname := config.Confs.DB.DbName

	// getting db for gorm
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", user, password, host, port, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}
	db.Debug()
	ctx := context.Background()
	DB = db.WithContext(ctx)
	models.DB = DB
	Migrate()
}
