package tools

import (
	"fmt"
	go_kit_log "github.com/go-kit/kit/log"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/k0kubun/pp"
	"github.com/labstack/gommon/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	log.SetPrefix(Highlight(Blue("[CC-Table] ")))
	//log.SetFlags(log.Llongfile)

	initZapLog()
}

// go native logger
// -----------------------------------

func CurrentMethod() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name() + "()"
}

func MethodError() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name() + "() is error"
}

func MethodOk() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name() + "() is ok"
}

func MethodSuccess() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name() + "() is success"
}

func MethodFailure() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name() + "() is failure"
}

func ParameterError() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name() + "()'s parameter is error"
}

// uber zap logger
// ----------------------------------

var my_logger *zap.Logger
var disable_debug_log = false
var disable_error_log = false

func Get_mylogger() zap.Logger {
	return *my_logger
}

func Disable_debug_logger() {
	disable_debug_log = true
}

func Disable_error_logger() {
	disable_error_log = true
}

func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}

func initZapLog() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "t",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "trace",
			LineEnding:    zapcore.DefaultLineEnding,
			//EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeLevel: nil,
			//EncodeTime:     formatEncodeTime,
			EncodeTime:     nil,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   nil,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		//InitialFields: map[string]interface{}{
		//	"app": "iTable",
		//},
	}
	var err error
	my_logger, err = cfg.Build()
	if err != nil {
		panic("Zap logger init fail:" + err.Error())
	}
}

func Debug_log(msg string, args ...interface{}) {

	if disable_debug_log == true {
		return
	}

	if len(args) > 0 {
		str := color.Green("\n----------------------------------------------------------------------------------------\n")
		message := str + color.White("调试: ") + color.Green(msg) + str + pp.Sprint(args) + str
		my_logger.Sugar().Debug(message)
	} else {
		str := color.Green("\n----------------------------------------------------------------------------------------\n")
		message := str + color.White("调试: ") + color.Green(msg) + str
		my_logger.Sugar().Debug(message)
	}

}

func Error_log(msg string, args ...interface{}) {
	if disable_error_log == true {
		return
	}

	if len(args) > 0 {
		str := color.Red("\n----------------------------------------------------------------------------------------\n")
		message := str + color.Yellow("错误: ") + color.Yellow(msg) + str + pp.Sprint(args) + str
		my_logger.Sugar().Debug(message)
	} else {
		str := color.Red("\n----------------------------------------------------------------------------------------\n")
		message := str + color.Yellow("错误: ") + color.Yellow(msg) + str
		my_logger.Sugar().Debug(message)
	}
}

func Log(args ...interface{}) {

	if disable_debug_log == true {
		return
	}

	pc, _, _, _ := runtime.Caller(1)
	func_name := runtime.FuncForPC(pc).Name() + "()"
	file_name, file_line := runtime.FuncForPC(pc).FileLine(pc)

	func_name_split := strings.Split(func_name, "/")
	func_name_split_length := len(func_name_split)
	short_func_name := func_name_split[func_name_split_length-1]

	file_name_split := strings.Split(file_name, "/")
	file_name_split_length := len(file_name_split)

	short_file_name := ""
	switch file_name_split_length {
	case 0:
		short_file_name = ""
	case 1:
		short_file_name = file_name_split[0]
	case 2:
		short_file_name = file_name_split[0] + "/" + file_name_split[1]
	default:
		start := file_name_split_length - 3
		short_file_name = file_name_split[start] + "/" + file_name_split[start+1] + "/" + file_name_split[start+2]
	}

	line := color.Grey("\n--------------------------------------\n")
	line_in := color.Grey("\n--------------------------------------\n")
	if len(args) > 0 {
		message := line +
			color.White("Function: ") + color.Yellow(short_func_name) + "\n" +
			color.White("Location: ") + color.Yellow(short_file_name+" [") +
			color.Red(strconv.Itoa(file_line)) + color.Green("] ") +
			line_in + pp.Sprint(args) + "\n"

		my_logger.Sugar().Debug(message)
	}
}

// go kit log
// --------------------------------------

func Get_gokit_logger() go_kit_log.Logger {
	var kit_logger go_kit_log.Logger
	kit_logger = go_kit_log.NewLogfmtLogger(os.Stderr)
	kit_logger = go_kit_log.With(kit_logger, "ts", go_kit_log.DefaultTimestampUTC)
	kit_logger = go_kit_log.With(kit_logger, "caller", go_kit_log.DefaultCaller)
	return kit_logger
}
