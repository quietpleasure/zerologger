package zerologger

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type Hook struct {
	hooker hooker
}

type hooker interface {
	SendLog(ctx context.Context, ts time.Time, lvl string, msg []byte) error
}

func (h Hook) Run(e *zerolog.Event, lvl zerolog.Level, msg string) {
	h.sendLog(time.Now(), lvl.String(), []byte(msg))

}

func (h Hook) sendLog(ts time.Time, lvl string, msg []byte) error {
	return h.hooker.SendLog(
		context.TODO(),
		ts,
		lvl,
		msg,
	)
}
