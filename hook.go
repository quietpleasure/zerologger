package zerologger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type Hook struct {
	Hooker hooker
}

type hooker interface {
	SendLog(ctx context.Context, ts time.Time, lvl string, msg []byte) error
}

var appInsightsProperties []string = []string{
	"database",
	"host",
	"port",
	"module",
	"pid",
	"name",
	"alreadyPrepared",
	"args",
	"sql",
	"commandTag",
}

/*
buf, err := zapcore.NewJSONEncoder(cfg).EncodeEntry(en, fields)

	if err != nil {
		return err
	}
	type log struct {
		Timestamp string `json:"ts"`
		Message   string `json:"msg"`
		Level     string `json:"lvl"`
	}
	var l log
	if err := json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(&l); err != nil {
		return err
	}
*/
func (h *Hook) Run(e *zerolog.Event, lvl zerolog.Level, msg string) {
	fields := make(map[string]any)
	e.Fields(func(key string, value interface{}) {
		fields[key] = fmt.Sprintf("%v", value)
	})

	// Add predefined properties
	for _, prop := range appInsightsProperties {
		if value, ok := e.GetCtx().Value(prop).(string); ok {
			fields[prop] = value
		}
	}

	fmt.Printf("FIELDS: %v\n", fields)
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(fields); err != nil {
		e.Err(err).Msg("hook encode fields")
	}
	// fmt.Println("FIELDS:", buf.String())
	// Add predefined properties
	// for _, prop := range appInsightsProperties {
	// 	if value, ok := e.Context[prop].(string); ok {
	// 		telemetry.Properties[prop] = value
	// 	}
	// }
	// buf, err := zapcore.NewJSONEncoder(cfg).EncodeEntry(en, fields)
	// if err != nil {
	// 	return err
	// }
	// type log struct {
	// 	Timestamp string `json:"ts"`
	// 	Message   string `json:"msg"`
	// 	Level     string `json:"lvl"`
	// }
	// var l log
	// buf := new(bytes.Buffer)
	// if err := json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(&l); err != nil {
	// 	e.Err(err).Msg("hook json decode fields")
	// }
	if err := h.Hooker.SendLog(
		context.Background(),
		time.Now(),
		lvl.String(),
		[]byte(msg),
	); err != nil {
		e.Err(err).Msg("hook send log")
	}
}
