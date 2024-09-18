package zerologger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Option func(option *options) error

type options struct {
	level           *string
	timeFormat      *string
	filePath        *string
	pretty          *bool
	caller          *bool
	fullCaller      *bool
	rotateAtStartup *bool
	maxSize         *int
	maxBackups      *int
	maxAge          *int
	localtime       *bool
	compress        *bool
}

func New(opts ...Option) (*zerolog.Logger, error) {
	var (
		err    error
		logctx zerolog.Context
	)
	var (
		writer io.Writer
		opt    options
	)
	for _, option := range opts {
		if err := option(&opt); err != nil {
			return nil, err
		}
	}
	//check and set level
	var lvl zerolog.Level

	if opt.level == nil || (opt.level != nil && *opt.level == "") {
		lvl = zerolog.NoLevel
	} else {
		lvl, err = zerolog.ParseLevel(*opt.level)
		if err != nil {
			return nil, err
		}
	}
	zerolog.SetGlobalLevel(lvl)
	// zerolog.ErrorStackMarshaler = func(err error) interface{} {
	// 	return pkgerrors.MarshalStack(errors.WithStack(err))
	// }
	//set time format
	if opt.timeFormat == nil || (opt.timeFormat != nil && *opt.timeFormat == "") {
		zerolog.TimeFieldFormat = time.RFC3339Nano
	} else {
		zerolog.TimeFieldFormat = *opt.timeFormat
	}

	if (opt.caller != nil && *opt.caller) || (opt.fullCaller != nil && *opt.fullCaller) {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return fmt.Sprintf("%s:%d", file, line)
		}
		if (opt.caller != nil && *opt.caller) && (opt.fullCaller == nil || (opt.fullCaller != nil && !*opt.fullCaller)) {
			zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
				return filepath.Base(fmt.Sprintf("%s:%d", file, line))
			}
		}
	}

	consoleWriter := func(out io.Writer) zerolog.ConsoleWriter {
		console := zerolog.ConsoleWriter{
			Out:        out,
			TimeFormat: zerolog.TimeFieldFormat,
		}
		if (opt.caller != nil && *opt.caller) || (opt.fullCaller != nil && *opt.fullCaller) {
			console.FormatCaller = func(i interface{}) string {
				return fmt.Sprintf("%s >", i)
			}
			if (opt.caller != nil && *opt.caller) && (opt.fullCaller == nil || (opt.fullCaller != nil && !*opt.fullCaller)) {
				console.FormatCaller = func(i interface{}) string {
					return filepath.Base(fmt.Sprintf("%s >", i))
				}
			}
		}
		return console
	}

	console := consoleWriter(os.Stdout)

	if opt.filePath == nil || (opt.filePath != nil && *opt.filePath == "") {
		writer = console
	} else {
		file := newRollingFile(*opt.filePath, opt)
		if opt.pretty != nil && *opt.pretty {
			prettyWriter := consoleWriter(file)
			prettyWriter.NoColor = true
			writer = zerolog.MultiLevelWriter(console, prettyWriter)
		} else {
			writer = zerolog.MultiLevelWriter(console, file) // так будет всегда выводить в консоль претти
			// writer = zerolog.MultiLevelWriter(os.Stdout, file) //так будет выводить в консоль претти или json
		}
	}

	logctx = zerolog.New(writer).With().Timestamp()

	if (opt.caller != nil && *opt.caller) || (opt.fullCaller != nil && *opt.fullCaller) {
		logctx = logctx.Caller()
	}

	l := logctx.Logger()
	return &l, nil
}

// Without option or  "", "disabled" level=NoLevel logger disabled.
// Possible values: "trace", "debug", "info", "warn", "error", "fatal", "panic"
func WithLevel(level string) Option {
	return func(options *options) error {
		options.level = &level
		return nil
	}
}

// Default format time.RFC3339Nano 2023-12-28T18:33:58.9954552+02:00
func WithCustomTimestamp(timeformat string) Option {
	return func(options *options) error {
		options.timeFormat = &timeformat
		return nil
	}
}

func WithPretty(with bool) Option {
	return func(options *options) error {
		options.pretty = &with
		return nil
	}
}

func WithCaller(with bool) Option {
	return func(options *options) error {
		options.caller = &with
		return nil
	}
}

func WithFullCaller(with bool) Option {
	return func(options *options) error {
		options.fullCaller = &with
		return nil
	}
}

// default console writer
func WithFile(filepath string) Option {
	return func(options *options) error {
		options.filePath = &filepath
		return nil
	}
}

func WithRotateAtStartup(with bool) Option {
	return func(options *options) error {
		options.rotateAtStartup = &with
		return nil
	}
}

func WithCompress(with bool) Option {
	return func(options *options) error {
		options.compress = &with
		return nil
	}
}

func WithLocalTime(with bool) Option {
	return func(options *options) error {
		options.localtime = &with
		return nil
	}
}

func WithMaxSize(size int) Option {
	return func(options *options) error {
		if size < 0 {
			return fmt.Errorf("file size cannot be less than zero")
		}
		options.maxSize = &size
		return nil
	}
}

func WithMaxBackups(backups int) Option {
	return func(options *options) error {
		if backups < 0 {
			return fmt.Errorf("number of files cannot be less than zero")
		}
		options.maxBackups = &backups
		return nil
	}
}

func WithMaxAge(age int) Option {
	return func(options *options) error {
		if age < 0 {
			return fmt.Errorf("number of days cannot be less than zero")
		}
		options.maxAge = &age
		return nil
	}
}

func newRollingFile(pathfile string, opts options) io.Writer {
	var (
		localtime, compress         bool
		maxsize, maxbackups, maxage int
	)

	if opts.localtime != nil {
		localtime = *opts.localtime
	}
	if opts.compress != nil {
		compress = *opts.compress
	}
	if opts.maxSize != nil {
		maxsize = *opts.maxSize
	}
	if opts.maxBackups != nil {
		maxbackups = *opts.maxBackups
	}
	if opts.maxAge != nil {
		maxage = *opts.maxAge
	}

	rotator := &lumberjack.Logger{
		Filename:   pathfile,
		LocalTime:  localtime,
		Compress:   compress,
		MaxSize:    maxsize,
		MaxBackups: maxbackups,
		MaxAge:     maxage,
	}
	if opts.rotateAtStartup != nil && *opts.rotateAtStartup {
		rotator.Rotate()
	}
	return rotator
}
