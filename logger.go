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

	emptyString = ""

	infolevelString  = "infolevel"
	warnlevelString  = "warnlevel"
	errlevelString   = "errorlevel"
	paniclevelString = "paniclevel"
	debuglevelString = "debuglevel"

	noMatchFound = "No match found for logger zap fields type:"

	stderr = "stderr"

	staringLogroutine = "Starting Log Gorutine"

	loadingFail1  = "Load LogConfig config failed: "
	loadingFail2  = ".\nLoading default log config"
	errorAtLogger = "Error at logger: "
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

	// zapMap map[string]func(string, ...zapcore.Field)
	zapInfo, zapWarn, zapErr, zapPanic, zapDebug func(string, ...zapcore.Field)
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
	log     *zap.Logger
	msg     string
	logType string
	fields  map[string]interface{}
}

var (
	logInfoChan, logWarnChan, logErrChan, logPanicChan,
	logDebugChan, logChan chan logMessages
	env  string
	logr Logger
)

func (l *logger) Info(msg string, opts map[string]interface{}) {
	if l.logInfo == nil {
		return
	}
	if l.infoAsync {
		logChan <- logMessages{logType: infostring, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logInfo.Info(msg, zapOpts...)
}

func (l *logger) Warning(msg string, opts map[string]interface{}) {
	if l.logWarn == nil {
		return
	}
	if l.warnAsync {
		logChan <- logMessages{logType: warnstring, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logWarn.Warn(msg, zapOpts...)

}

func (l *logger) Error(msg string, opts map[string]interface{}) {
	if l.logErr == nil {
		return
	}
	if l.errAsync {
		logChan <- logMessages{logType: errstring, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logErr.Error(msg, zapOpts...)

}

func (l *logger) Panic(msg string, opts map[string]interface{}) {
	if l.logPanic == nil {
		return
	}
	if l.panicAsync {
		logChan <- logMessages{logType: panicstring, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logPanic.Panic(msg, zapOpts...)
}

func (l *logger) Debug(msg string, opts map[string]interface{}) {
	if l.logDebug == nil {
		return
	}
	if l.debugAsync {
		logChan <- logMessages{logType: debugstring, msg: msg, fields: opts}
		return
	}
	zapOpts := GetFields(opts)
	l.logDebug.Debug(msg, zapOpts...)
}

func (l *logger) CloseAll() {
	l.fileClose()
	close(logChan)
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
		case uint:
			zapFields[i] = zap.Uint(k, v)
		case []uint:
			zapFields[i] = zap.Uints(k, v)
		case uint8:
			zapFields[i] = zap.Uint8(k, v)
		case []uint8:
			zapFields[i] = zap.Uint8s(k, v)
		case uint16:
			zapFields[i] = zap.Uint16(k, v)
		case []uint16:
			zapFields[i] = zap.Uint16s(k, v)
		case uint32:
			zapFields[i] = zap.Uint32(k, v)
		case []uint32:
			zapFields[i] = zap.Uint32s(k, v)
		case uint64:
			zapFields[i] = zap.Uint64(k, v)
		case []uint64:
			zapFields[i] = zap.Uint64s(k, v)
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
			log.Fatal(noMatchFound, v)
		}
		i++
	}
	return zapFields
}

//GetInstance to retrieve single instance of Logger
func GetInstance() Logger {
	if logr == nil {
		logr, _ = NewLogger(emptyString, emptyString)
	}
	return logr
}

//GetZapLogger ..
func GetZapLogger(ws zapcore.WriteSyncer, l *Log, logLevel string) (*zap.Logger, bool) {
	if l == nil {
		log.Println("no config found for ", logLevel)
		log.Println("ignoring ", logLevel)
		return nil, false
	}
	var z []zap.Option

	var zc zapcore.Core
	var zl zapcore.Level

	if l.Tracelevel != emptyString {
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
	return newZap, l.Async

}

func getLevel(s string) zapcore.Level {
	switch s {
	case infolevelString:
		return zapcore.InfoLevel
	case warnlevelString:
		return zapcore.WarnLevel
	case errlevelString:
		return zapcore.ErrorLevel
	case paniclevelString:
		return zapcore.PanicLevel
	case debuglevelString:
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
func NewLogger(data, env string) (Logger, error) {

	l, err := LoadLogConfig(data, env)
	if err != nil {
		log.Println(loadingFail1 + err.Error() + loadingFail2)
	}
	ws, f, err := zap.Open(stderr)
	if err != nil {
		log.Fatal(errorAtLogger, err)
	}
	var loger *logger
	// zapMap = make(map[string]func(string, ...zapcore.Field))
	if l == nil {
		loger = &logger{
			logInfo:    zap.NewNop(),
			infoAsync:  false,
			logWarn:    zap.NewNop(),
			warnAsync:  false,
			logErr:     zap.NewNop(),
			errAsync:   false,
			logPanic:   zap.NewNop(),
			panicAsync: false,
			logDebug:   zap.NewNop(),
			debugAsync: false,
			fileClose:  f,
		}
	} else {

		loger = &logger{}
		loger.logInfo, loger.infoAsync = GetZapLogger(ws, l[infostring], infostring)
		loger.logWarn, loger.warnAsync = GetZapLogger(ws, l[warnstring], warnstring)
		loger.logErr, loger.errAsync = GetZapLogger(ws, l[errstring], errstring)
		loger.logPanic, loger.panicAsync = GetZapLogger(ws, l[panicstring], panicstring)
		loger.logDebug, loger.debugAsync = GetZapLogger(ws, l[debugstring], debugstring)
		loger.fileClose = f
		//
		// logInfo:    GetZapLogger(ws, l[infostring]),
		// infoAsync:  l[infostring].Async,
		// logWarn:    GetZapLogger(ws, l[warnstring]),
		// warnAsync:  l[warnstring].Async,
		// logErr:     GetZapLogger(ws, l[errstring]),
		// errAsync:   l[errstring].Async,
		// logPanic:   GetZapLogger(ws, l[panicstring]),
		// panicAsync: l[panicstring].Async,
		// logDebug:   GetZapLogger(ws, l[debugstring]),
		// debugAsync: l[debugstring].Async,
		// fileClose:  f,
		// }
	}

	logr = loger

	// zapMap[infostring]
	zapInfo = loger.logInfo.Info
	// zapMap[warnstring]
	zapWarn = loger.logWarn.Warn
	// zapMap[errstring]
	zapErr = loger.logErr.Error
	// zapMap[panicstring]
	zapPanic = loger.logPanic.Panic
	// zapMap[debugstring]
	zapDebug = loger.logDebug.Debug

	logChan = make(chan logMessages) //, getLogBufferSize())

	log.Println(staringLogroutine)
	go logRoutine()
	return logr, nil
}

//LogRoutine ...
func logRoutine() {

	// for l := range logChan {
	// 	opts := GetFields(l.fields)
	// 	zapMap[l.logType](l.msg, opts...)
	// }
	for {
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		log.Println("Recovered in LogRoutine")
		// 		if logChan == nil {
		// 			logChan = make(chan logMessages)
		// 			go logRoutine()
		// 		}
		// 	}
		// }()
		select {
		case l := <-logChan:
			opts := GetFields(l.fields)
			switch l.logType {
			case infostring:
				zapInfo(l.msg, opts...)
			case warnstring:
				zapWarn(l.msg, opts...)
			case errstring:
				zapErr(l.msg, opts...)
			case panicstring:
				zapPanic(l.msg, opts...)
			case debugstring:
				zapDebug(l.msg, opts...)
			}

			//zapMap[l.logType](l.msg, opts...)
		}
	}
}
