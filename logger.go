package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	infostring  = "info"
	warnstring  = "warn"
	errstring   = "err"
	panicstring = "panic"
	debugstring = "debug"
)

var (
	zapEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
)

//Logger to implement zap logger
type Logger interface {
	Info(string, map[string]interface{})
	Warning(string, map[string]interface{})
	Error(string, map[string]interface{})
	Panic(string, map[string]interface{})
	Debug(string, map[string]interface{})
	CloseAll()
}

type logger struct {
	logInfo    *zap.Logger
	infoAsync  bool
	logWarn    *zap.Logger
	warnAsync  bool
	logErr     *zap.Logger
	errAsync   bool
	logPanic   *zap.Logger
	panicAsync bool
	logDebug   *zap.Logger
	debugAsync bool
	fileClose  func()
}

type logMessages struct {
	log    *zap.Logger
	msg    string
	fields map[string]interface{}
}

var (
	logInfoChan, logWarnChan, logErrChan, logPanicChan,
	logDebugChan chan logMessages
	env  string
	logr Logger
)

func (l *logger) Info(msg string, opts map[string]interface{}) {
	if l.infoAsync {
		logInfoChan <- logMessages{log: l.logInfo, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logInfo.Info(msg, zapOpts...)
}

func (l *logger) Warning(msg string, opts map[string]interface{}) {
	if l.warnAsync {
		logWarnChan <- logMessages{log: l.logWarn, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logWarn.Warn(msg, zapOpts...)

}

func (l *logger) Error(msg string, opts map[string]interface{}) {
	if l.errAsync {
		logErrChan <- logMessages{log: l.logErr, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logErr.Error(msg, zapOpts...)

}

func (l *logger) Panic(msg string, opts map[string]interface{}) {
	if l.panicAsync {
		logPanicChan <- logMessages{log: l.logPanic, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logPanic.Panic(msg, zapOpts...)
}

func (l *logger) Debug(msg string, opts map[string]interface{}) {
	if l.debugAsync {
		logDebugChan <- logMessages{log: l.logDebug, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logDebug.Debug(msg, zapOpts...)
}

func (l *logger) CloseAll() {
	l.fileClose()
}

//GetFields conver maps into zapcore.Fields
func GetFields(fields map[string]interface{}) []zapcore.Field {
	if fields == nil {
		return make([]zapcore.Field, 0)
	}
	zapFields := make([]zapcore.Field, len(fields))
	i := 0
	for k, v := range fields {
		switch v := v.(type) {
		case int:
			zapFields[i] = zap.Int(k, v)
		case []int:
			zapFields[i] = zap.Ints(k, v)
		case int8:
			zapFields[i] = zap.Int8(k, v)
		case []int8:
			zapFields[i] = zap.Int8s(k, v)
		case int16:
			zapFields[i] = zap.Int16(k, v)
		case []int16:
			zapFields[i] = zap.Int16s(k, v)
		case int32:
			zapFields[i] = zap.Int32(k, v)
		case []int32:
			zapFields[i] = zap.Int32s(k, v)
		case int64:
			zapFields[i] = zap.Int64(k, v)
		case []int64:
			zapFields[i] = zap.Int64s(k, v)
		case string:
			zapFields[i] = zap.String(k, v)
		case []string:
			zapFields[i] = zap.Strings(k, v)
		case float32:
			zapFields[i] = zap.Float32(k, v)
		case []float32:
			zapFields[i] = zap.Float32s(k, v)
		case float64:
			zapFields[i] = zap.Float64(k, v)
		case []float64:
			zapFields[i] = zap.Float64s(k, v)
		case bool:
			zapFields[i] = zap.Bool(k, v)
		case []bool:
			zapFields[i] = zap.Bools(k, v)
		case error:
			zapFields[i] = zap.Error(v)
		case *error:
			zapFields[i] = zap.Error(*v)
		default:
			log.Fatal("No match found for logger zap fields type:", v)
		}
		i++
	}
	return zapFields
}

//GetInstance to retrieve single instance of Logger
func GetInstance() Logger {
	if logr == nil {
		logr = NewLogger("", "")
	}
	return logr
}

//GetZapLogger ..
func GetZapLogger(ws zapcore.WriteSyncer, l *Log) *zap.Logger {
	var z []zap.Option

	var zc zapcore.Core
	var zl zapcore.Level
	if l.Tracelevel != "" {
		zl = getLevel(l.Tracelevel)
	}

	if l.Erroroutput {
		z = append(z, zap.ErrorOutput(ws))
	}

	if l.Stacktrace {
		z = append(z, zap.AddStacktrace(zl))
	}

	if l.Caller {
		z = append(z, zap.AddCaller())
	}

	z = append(z, zap.AddCallerSkip(l.CallerSkip))

	zc = zapCoreConfig(ws, zl)
	newZap := zap.New(zc, z...)
	return newZap

}

func getLevel(s string) zapcore.Level {
	switch s {
	case "infolevel":
		return zapcore.InfoLevel
	case "warnlevel":
		return zapcore.WarnLevel
	case "errorlevel":
		return zapcore.ErrorLevel
	case "paniclevel":
		return zapcore.PanicLevel
	case "debuglevel":
		return zapcore.DebugLevel
	}
	return 10
}

func zapCoreConfig(ws zapcore.WriteSyncer, l zapcore.Level) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zapEncoderConfig),
		ws,
		l,
	)
}

//NewLogger Get new instance of Logger
func NewLogger(data, env string) Logger {

	l, err := LoadLogConfig(data, env)
	if err != nil {
		log.Println("Load LogConfig config failed: " + err.Error() + ".\nLoading default log config")
	}
	ws, f, err := zap.Open("stderr")
	if err != nil {
		log.Fatal("Error at logger: ", err)
	}
	if l == nil {
		logr = &logger{
			logInfo:   zap.NewNop(),
			logWarn:   zap.NewNop(),
			logErr:    zap.NewNop(),
			logPanic:  zap.NewNop(),
			logDebug:  zap.NewNop(),
			fileClose: f,
		}
	} else {

		logr = &logger{
			logInfo:    GetZapLogger(ws, l[infostring]),
			infoAsync:  l[infostring].Async,
			logWarn:    GetZapLogger(ws, l[warnstring]),
			warnAsync:  l[warnstring].Async,
			logErr:     GetZapLogger(ws, l[errstring]),
			errAsync:   l[errstring].Async,
			logPanic:   GetZapLogger(ws, l[panicstring]),
			panicAsync: l[panicstring].Async,
			logDebug:   GetZapLogger(ws, l[debugstring]),
			debugAsync: l[debugstring].Async,
			fileClose:  f,
		}
	}

	logInfoChan = make(chan logMessages)
	logWarnChan = make(chan logMessages)
	logErrChan = make(chan logMessages)
	logPanicChan = make(chan logMessages)
	logDebugChan = make(chan logMessages)
	LogRoutine()
	return logr
}

//LogRoutine ...
func LogRoutine() {
	go func() {
		for {
			select {
			case l := <-logInfoChan:
				opts := GetFields(l.fields)
				l.log.Info(l.msg, opts...)

			case l := <-logWarnChan:
				opts := GetFields(l.fields)
				l.log.Warn(l.msg, opts...)

			case l := <-logErrChan:
				opts := GetFields(l.fields)
				l.log.Error(l.msg, opts...)

			case l := <-logPanicChan:
				opts := GetFields(l.fields)
				l.log.Panic(l.msg, opts...)

			case l := <-logDebugChan:
				opts := GetFields(l.fields)
				l.log.Debug(l.msg, opts...)
			}
		}
	}()
}
