package logger

import (
	"errors"
	"io/ioutil"
	"log"

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
	Log map[string]*Log
}

//LogConfig ...
type LogConfig map[string]*Log

var app map[string]*Log

//LoadLogConfig ...
func LoadLogConfig(data string, env string) (map[string]*Log, error) {

	if data == "" {
		return nil, errors.New("No file provided")
	}
	var apps map[string]App //map[string]Broker

	// data := getFile(filepath)
	if _, err := toml.Decode(data, &apps); err != nil {
		return nil, err
	}
	tempapp, ok := apps[env]
	if !ok {
		return nil, errors.New("Environment not found in configuration: " + env)
	}
	app = tempapp.Log
	return app, nil

}

func getFile(fpath string) string {
	if fpath == "" {
		return ""
	}
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatal("File read error: ", err)
	}
	return string(bs)
}

// //LogConfig ...
// type LogConfig interface {
// }
