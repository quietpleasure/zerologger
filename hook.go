package zerologger

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type Hook struct {
	Hooker hooker
}

type hooker interface {
	SendLog(ctx context.Context, ts time.Time, lvl string, msg []byte) error
}

func (h *Hook) Run(e *zerolog.Event, lvl zerolog.Level, msg string) {
	if err := h.Hooker.SendLog(
		context.Background(),
		time.Now(),
		lvl.String(),
		[]byte(msg),
	); err != nil {
		e.Err(err).Msg("hook send log")
	}
}
