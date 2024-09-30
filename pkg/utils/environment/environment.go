package environment

import (
	"github.com/gin-gonic/gin"
	"music-library/pkg/utils/environment/variable"
	"music-library/pkg/utils/logger"
)

func SetupEnv(mode string) {
	switch mode {
	case variable.Debug:
		logger.SetupLogger(variable.Debug)
		logger.Logger.Info("logger initiated")
		logger.Logger.Debug("DEBUG mode set")

		gin.ForceConsoleColor()
		gin.SetMode(gin.DebugMode)

	case variable.Release:
		logger.SetupLogger(variable.Release)
		logger.Logger.Info("logger initiated")
		logger.Logger.Info("RELEASE mode set")

		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}
}
