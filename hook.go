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
	if err := json.Unmarshal([]byte(ev), &logData); err != nil {
		e.Err(err).Msg("hook send unmarshal fields")
		return
	}

	// now you can either access a map of the data (logData) or string of the data (ev)

	logData["msg"] = msg

	// var l log
	data, err := json.Marshal(logData)
	if err != nil {
		e.Err(err).Msg("hook send marshal fields")
		return
	}
	if err := h.Hooker.SendLog(
		context.Background(),
		time.Now(),
		lvl.String(),
		data,
	); err != nil {
		e.Err(err).Msg("hook send log")
	}
	// {"ts": "13.09.2024 14:57:56.853640", "msg": "Prepare", "pid": 15136, "sql": "INSERT INTO executor_panel.refresh_tokens (key_token,val_token,fingerprint,created_at) VALUES ($1,$2,$3,$4) \n\tON CONFLICT ON CONSTRAINT uq_refresh_token DO UPDATE SET val_token=$2,created_at=$4", "name": "stmtcache_678b68ca1c6923ce1e6a6e324f9b30c961ed68d38b998d01", "time": 0.0007081, "level": "info", "alreadyPrepared": false}
}
