package song

import (
	"music-library/internal/handlers/interfaces/song"
)

type songHandler struct{}

func NewSongHandler() song.SongInterface {
	return &songHandler{}
}
