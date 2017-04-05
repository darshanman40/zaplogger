package logger

import (
	"errors"

	"github.com/BurntSushi/toml"
)

//Log ...
type Log struct {
	Tracelevel  string
	Stacktrace  bool
	Erroroutput bool
	Caller      bool
	CallerSkip  int `toml:"caller_skip"`
	Async       bool
}

//App ...
type App struct {
	Log           map[string]*Log
	LogBufferSize int `toml:"log_buffer_size"`
}

//LogConfig ...
type LogConfig map[string]*Log

var app map[string]*Log

var bufferSize int

//LoadLogConfig ...
func LoadLogConfig(data string, env string) (map[string]*Log, error) {

	if data == "" {
		return nil, errors.New("No data provided")
	}
	var apps map[string]App
	if _, err := toml.Decode(data, &apps); err != nil {
		return nil, err
	}
	tempapp, ok := apps[env]
	bufferSize = tempapp.LogBufferSize
	if !ok {
		return nil, errors.New("Environment not found in configuration: " + env)
	}
	app = tempapp.Log
	return app, nil
}

//GetLogBufferSize ...
func getLogBufferSize() int {
	return bufferSize
}
