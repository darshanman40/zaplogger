package logger

import (
	"errors"
	"io/ioutil"
	"log"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	tomlcontent = `[dev]
Env = "dev"
RabbitHost     = "localhost:5672"
RabbitLogin    = "guest"
RabbitPassword = "guest"

[dev.log.info]
tracelevel = "infolevel"
stacktrace =  false
erroroutput = false
caller = false
caller_skip= 1
async = false
[dev.log.warn]
tracelevel = "warnlevel"
stacktrace =  false
erroroutput = true
caller = false
caller_skip= 1
async = false
[dev.log.err]
tracelevel = "errorlevel"
stacktrace =  true
erroroutput = false
caller = false
caller_skip=2
async = false
[dev.log.panic]
tracelevel = "paniclevel"
stacktrace =  true
erroroutput = true
caller = false
caller_skip = 2
async = false

################################################################################################################
################################################################################################################
################################################################################################################

[debug]
Env = "debug"
RabbitHost     = "localhost:5672"
RabbitLogin    = "guest"
RabbitPassword = "guest"


[debug.log.info]
tracelevel = "infolevel"
stacktrace =  false
erroroutput = false
caller = false
caller_skip= 1
async = true
[debug.log.warn]
tracelevel = "warnlevel"
stacktrace =  false
erroroutput = true
caller = false
caller_skip= 1
async = true
[debug.log.err]
tracelevel = "errorlevel"
stacktrace =  true
erroroutput = false
caller = true
caller_skip=2
async = true
[debug.log.panic]
tracelevel = "paniclevel"
stacktrace =  true
erroroutput = true
caller = true
caller_skip = 2
async = true
[debug.log.debug]
tracelevel = "debuglevel"
stacktrace =  false
erroroutput = true
caller = false
caller_skip = 1
async = true
`

	syncLogConfig = `[debug]
Env = "debug"
RabbitHost     = "localhost:5672"
RabbitLogin    = "guest"
RabbitPassword = "guest"

[debug.log.info]
tracelevel = "infolevel"
stacktrace =  false
erroroutput = false
caller = false
caller_skip= 1
async = false
[debug.log.warn]
tracelevel = "warnlevel"
stacktrace =  false
erroroutput = true
caller = false
caller_skip= 1
async = false
[debug.log.err]
tracelevel = "errorlevel"
stacktrace =  true
erroroutput = false
caller = true
caller_skip=2
async = false
[debug.log.panic]
tracelevel = "paniclevel"
stacktrace =  true
erroroutput = true
caller = true
caller_skip = 2
async = false
[debug.log.debug]
tracelevel = "debuglevel"
stacktrace =  false
erroroutput = true
caller = false
caller_skip = 1
async = false
`
)

func TestGetFile(t *testing.T) {
	s := getFile("config.sample.toml")
	if s != tomlcontent {
		t.Fail()
	}
}

func TestLoadLogger(t *testing.T) {

	l := NewLogger(getFile("config.sample.toml"), "debug")
	if l != GetInstance() {
		t.Fail()
	}
}

func TestEmptyLogger(t *testing.T) {

	l := NewLogger("", "")
	if l != GetInstance() {
		t.Fail()
	}
}

func TestIntsGetFields(t *testing.T) {
	origZapInts := zap.Ints("intsKey", []int{int(9)})
	origZapInt8s := zap.Int8s("int8sKey", []int8{int8(9)})
	origZapInt16s := zap.Int16s("int16sKey", []int16{int16(9)})
	origZapInt32s := zap.Int32s("int32sKey", []int32{int32(9)})
	origZapInt64s := zap.Int64s("int64sKey", []int64{int64(9)})

	zapInts := GetFields(map[string]interface{}{
		"intsKey": []int{int(9)},
	})
	zapInt8s := GetFields(map[string]interface{}{
		"int8sKey": []int8{int8(9)},
	})
	zapInt16s := GetFields(map[string]interface{}{
		"intsKey": []int16{int16(9)},
	})
	zapInt32s := GetFields(map[string]interface{}{
		"int8sKey": []int32{int32(9)},
	})
	zapInt64s := GetFields(map[string]interface{}{
		"intsKey": []int64{int64(9)},
	})

	if origZapInts.Type != zapInts[0].Type {
		t.Fail()
	}
	if origZapInt8s.Type != zapInt8s[0].Type {
		t.Fail()
	}
	if origZapInt16s.Type != zapInt16s[0].Type {
		t.Fail()
	}
	if origZapInt32s.Type != zapInt32s[0].Type {
		t.Fail()
	}
	if origZapInt64s.Type != zapInt64s[0].Type {
		t.Fail()
	}
}

func TestFloatsGetFields(t *testing.T) {
	origZapFloat32s := zap.Float32s("float32sKey", []float32{float32(9)})
	origZapFloat64s := zap.Float64s("float64sKey", []float64{float64(9)})

	zapFloat32s := GetFields(map[string]interface{}{
		"float32sKey": []float32{float32(9)},
	})
	zapFloat64s := GetFields(map[string]interface{}{
		"float64sKey": []float64{float64(9)},
	})

	if origZapFloat32s.Type != zapFloat32s[0].Type || origZapFloat32s.Key != zapFloat32s[0].Key {
		t.Fail()
	}
	if origZapFloat64s.Type != zapFloat64s[0].Type || origZapFloat64s.Key != zapFloat64s[0].Key {
		t.Fail()
	}
}

func TestStringsBoolsGetFields(t *testing.T) {
	origZapStrings := zap.Strings("stringsKey", []string{"abc"})
	zapStrings := GetFields(map[string]interface{}{
		"stringsKey": []string{"abc"},
	})
	origZapBools := zap.Bools("boolsKey", []bool{true})
	zapBools := GetFields(map[string]interface{}{
		"boolsKey": []bool{true},
	})
	if origZapStrings.Type != zapStrings[0].Type || origZapStrings.Key != zapStrings[0].Key {
		t.Fail()
	}
	if origZapBools.Type != zapBools[0].Type || origZapBools.Key != zapBools[0].Key {
		t.Fail()
	}
}

func TestGetFields(t *testing.T) {
	err := errors.New("custom error")
	originalFields := []zapcore.Field{
		zap.Int("intKey", 9),
		zap.Int8("int8Key", 1),
		zap.Int16("int16Key", 16),
		zap.Int32("int32Key", 32),
		zap.Int64("int64Key", 64),
		zap.Uint("uintKey", 9),
		zap.Uint8("uint8Key", 1),
		zap.Uint16("uint16Key", 16),
		zap.Uint32("uint32Key", 32),
		zap.Uint64("uint64Key", 64),
		zap.String("stringKey", "stringdata"),
		zap.Float32("float32Key", 2.56455),
		zap.Float64("float64Key", 3.5565645),
		zap.Error(err),
		zap.Bool("boolKey", true),
	}

	m := map[string]interface{}{
		"stringKey":  "stringdata",
		"intKey":     int(9),
		"int8Key":    int8(1),
		"int16Key":   int16(16),
		"int32Key":   int32(32),
		"int64Key":   int64(64),
		"uintKey":    uint(9),
		"uint8Key":   uint8(1),
		"uint16Key":  uint16(16),
		"uint32Key":  uint32(32),
		"uint64Key":  uint64(64),
		"float32Key": float32(2.56455),
		"float64Key": float64(3.5565645),
		"errKey":     err,
		"boolKey":    true,
	}
	fields := GetFields(m)
	if !testEqField(fields, originalFields) {
		t.Fail()
	}
}

func BenchmarkFullAsync(b *testing.B) {
	l := NewLogger(getFile("config.sample.toml"), "debug")
	for i := 0; i < b.N; i++ {
		l.Info("test", map[string]interface{}{
			"test type": "async",
		})
	}
}

func BenchmarkFullSync(b *testing.B) {
	l := NewLogger(syncLogConfig, "debug")
	for i := 0; i < b.N; i++ {
		l.Info("test", map[string]interface{}{
			"test type": "sync",
		})
	}
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

//testEqField is created to compare unorganized slice of zapcore.Field
func testEqField(a, b []zapcore.Field) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	temp := false
	for _, i := range a {
		for _, j := range b {
			if i != j {
				temp = false
				continue
			}
			temp = true
			break
		}
		if !temp {
			return temp
		}
	}

	return temp
}
