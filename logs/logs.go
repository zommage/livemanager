package logs

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Log *zap.SugaredLogger

func InitLog(logPath, logLevel string) error {
	var err error = nil
	Logger, Log, err = InitZapLog(logPath, logLevel)
	if err != nil {
		fmt.Println("Init zap log is err: ", err)
		return fmt.Errorf("Init zap log is err: %v", err)
	}

	return nil
}

/************************************
@ path: the path of log file
@ level: log level(DEBUG, INFO, ERROR)
@ return: *zap.Logger:  critical log
@ *zap.SugaredLogger: not critical log
@***********************************/
func InitZapLog(path string, level string) (*zap.Logger, *zap.SugaredLogger, error) {
	if path == "" || level == "" {
		fmt.Println("log path is nil or log level is nil")
		return nil, nil, fmt.Errorf("log path is nil or log level is nil")
	}

	var js string
	js = fmt.Sprintf(`{
             "level": "%s",
             "encoding": "json",
             "outputPaths": ["%s"],
             "errorOutputPaths": ["%s"]
             }`, level, path, path)

	var cfg zap.Config
	if err := json.Unmarshal([]byte(js), &cfg); err != nil {
		return nil, nil, fmt.Errorf("Init log json is err: %v", err)
	}
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	Logger, err := cfg.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("zap logger build error: %v", err)
	}

	Log := Logger.Sugar()

	return Logger, Log, nil
}
