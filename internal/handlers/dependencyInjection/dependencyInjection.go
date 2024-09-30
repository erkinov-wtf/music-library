package dependencyInjection

import (
	songApi "music-library/internal/handlers/api/song"
	"music-library/internal/handlers/interfaces/song"
)

type Container struct {
	SongHandler song.SongInterface
}

func NewContainer() *Container {
	return &Container{
		SongHandler: songApi.NewSongHandler(),
	}
}
