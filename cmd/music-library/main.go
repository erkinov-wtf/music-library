package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "music-library/docs"
	"music-library/internal/config"
	"music-library/internal/handlers/dependencyInjection"
	"music-library/internal/routes/song"
	"music-library/internal/storage/database"
	"music-library/pkg/utils/environment"
)

func main() {
	config.MustLoad()
	environment.SetupEnv(config.Cfg.General.Env)
	database.DBConnect()

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiPrefixRoute := router.Group("/api")
	container := dependencyInjection.NewContainer()

	song.SetupSongRoutes(apiPrefixRoute, container)

	err := router.Run(fmt.Sprintf(":%v", config.Cfg.General.Port))
	if err != nil {
		panic(err.Error())
	}
}
