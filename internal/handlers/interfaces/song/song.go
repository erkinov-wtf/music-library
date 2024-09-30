package song

import "github.com/gin-gonic/gin"

type SongInterface interface {
	Create(context *gin.Context)
	Index(context *gin.Context)
	Show(context *gin.Context)
	Update(context *gin.Context)
	Lyrics(context *gin.Context)
	Info(context *gin.Context)
	Delete(context *gin.Context)
}
