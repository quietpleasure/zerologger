package zerologger

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
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

	logData := make(map[string]interface{})
	// create a string that appends } to the end of the buf variable you access via reflection
	ev := fmt.Sprintf("%s}", reflect.ValueOf(e).Elem().FieldByName("buf"))
	json.Unmarshal([]byte(ev), &logData)

	// now you can either access a map of the data (logData) or string of the data (ev)

	fmt.Printf("FIELDS: %v\n", logData)

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
