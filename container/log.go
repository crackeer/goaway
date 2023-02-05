package container

import (
	"fmt"

	syslog "log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	writerHook "github.com/sirupsen/logrus/hooks/writer"
)

const (
	LogTypeAPI    = "api"
	LogTypeSVC    = "svc"
	LogTypeError  = "error"
	LogTypeRouter = "router"
)

var useLogTypes []string = []string{
	LogTypeAPI, LogTypeSVC, LogTypeError, LogTypeRouter,
}

var (
	loggers map[string]*logrus.Logger = map[string]*logrus.Logger{}
)

// InitLogger ...
func InitLogger(cfg *AppConfig) error {
	for _, flag := range useLogTypes {
		logger := logrus.New()
		logger.SetFormatter(&logrus.JSONFormatter{})

		writer, err := getLogOutput(cfg.LogDir, flag)
		if err != nil {
			return fmt.Errorf("get log `%s` output error, %s", flag, err.Error())
		}
		logger.SetOutput(writer)
		if cfg.Debug {
			logger.AddHook(&writerHook.Hook{
				Writer:    syslog.Default().Writer(),
				LogLevels: logrus.AllLevels,
			})
		}
		loggers[flag] = logger
	}
	return nil
}

func getLogOutput(rawLogPath string, flag string) (*rotatelogs.RotateLogs, error) {
	logPath := fmt.Sprintf("%s/%s.%s.log", rawLogPath, "%Y-%m-%d_%H", flag)
	return rotatelogs.New(
		logPath,
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
}

// DestroyLogger
func DestroyLogger() {
	for _, v := range useLogTypes {
		if l, exists := loggers[v]; exists {
			l.Writer().Close()
		}
	}
}

func GetLogger(flag string) *logrus.Logger {
	if l, exists := loggers[flag]; exists {
		return l
	}
	return nil
}

func log(flag string, fields map[string]interface{}, message interface{}) error {
	logger := GetLogger(flag)
	if logger == nil {
		return fmt.Errorf("%s logger nil", flag)
	}
	logger.WithFields(logrus.Fields{
		"channel": flag,
	}).WithFields(logrus.Fields(fields)).Info(message)
	return nil
}

func Log(fields map[string]interface{}, message interface{}) error {
	return log(LogTypeSVC, fields, message)
}

// LogError ...
func LogError(fields map[string]interface{}, message interface{}) error {
	return log(LogTypeError, fields, message)
}
