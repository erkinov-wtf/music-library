package song

import (
	"github.com/gin-gonic/gin"
	"music-library/internal/handlers/dependencyInjection"
)

func SetupSongRoutes(router *gin.RouterGroup, container *dependencyInjection.Container) {
	songRoutes := router.Group("/songs")
	{
		songRoutes.POST("/", container.SongHandler.Create)
		songRoutes.GET("/", container.SongHandler.Index)
		songRoutes.GET("/:id", container.SongHandler.Show)
		songRoutes.PUT("/:id", container.SongHandler.Update)
		songRoutes.GET("/:id/lyrics", container.SongHandler.Lyrics)
		songRoutes.GET("/info", container.SongHandler.Info)
		songRoutes.DELETE("/:id", container.SongHandler.Delete)
	}
}
