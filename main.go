package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vyas-git/lti_code_test/api"
	docs "github.com/vyas-git/lti_code_test/docs"
	"github.com/vyas-git/lti_code_test/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	LoadEnv()

	// DB creation and migration
	db, err := gorm.Open(postgres.Open("postgres://lti_spotify_user:quO3vtc70dfjv5Bad2Df9n6mWzrh9IVo@dpg-cm87fj0cmk4c7390r3s0-a.oregon-postgres.render.com/lti_spotify"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if db != nil {
		fmt.Println("database created succesfully")
	}
	db.AutoMigrate(&model.Track{})
	docs.SwaggerInfo.BasePath = "/"

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.Run(router, db)
	router.Run(":8080")
}
func LoadEnv() {
	// Load values from sample.env file for local use, not for prod
	err := godotenv.Load("sample.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
