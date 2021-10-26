package logging

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createDir() {
	path, _ := os.Getwd()

	if _, err := os.Stat(fmt.Sprintf("%s/logs", path)); os.IsNotExist(err) {
		_ = os.Mkdir("logs", os.ModePerm)
	}
}

func getLogWriter() zapcore.WriteSyncer {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	timestamp := time.Now().UTC().Format("Jan-02-2006")
	file, err := os.OpenFile(path+"/logs/"+timestamp+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(file)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05-0700"))
	})
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func InitLogger() {
	createDir()
	writerSync := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writerSync, zapcore.InfoLevel) // tweak the log level as needed
	logg := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logg)
}
