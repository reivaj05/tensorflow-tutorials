package GoLogger

import (
	"errors"
	"os"

	"github.com/Sirupsen/logrus"
)

const (
	PANIC = iota
	FATAL
	ERROR
	WARNING
	INFO
	DEBUG
)

type LoggerOptions struct {
	OutputFile string
	Path       string
	LogLevel   int
}

// Init initializes logger with the options passed as parameter
func Init(options *LoggerOptions) error {
	if !wasOutputFileProvided(options.OutputFile) {
		return errors.New("No file provided")
	}
	return setLoggerOptions(options)
}

func wasOutputFileProvided(outputFile string) bool {
	return outputFile != ""
}

func setLoggerOptions(options *LoggerOptions) error {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	file, err := openOrCreateFile(options.Path, options.OutputFile)
	if err != nil {
		return err
	}
	logrus.SetOutput(file)
	SetLogLevel(options.LogLevel)
	return nil
}

func openOrCreateFile(path, outputFile string) (*os.File, error) {
	if err := os.MkdirAll(path, 0744); err != nil {
		return nil, err
	}
	mode := os.O_APPEND | os.O_WRONLY | os.O_CREATE
	return os.OpenFile(path+outputFile, mode, 0755)
}

// SetLogLevel sets a new log level
func SetLogLevel(logLevel int) {
	logLevel = parseLogLevel(logLevel)
	logrus.SetLevel(logrus.Level(logLevel))
}

func parseLogLevel(logLevel int) int {
	if logLevel < PANIC || logLevel > DEBUG {
		return DEBUG
	}
	return logLevel
}

// GetLogLevel returns the current log level
func GetLogLevel() int {
	return int(logrus.GetLevel())
}

// LogPanic logs panic msg
func LogPanic(msg string, fields map[string]interface{}) {
	if fields != nil {
		logrus.WithFields(fields).Panic(msg)
	} else {
		logrus.Panic(msg)
	}
}

// LogFatal logs fatal msg
func LogFatal(msg string, fields map[string]interface{}) {
	if fields != nil {
		logrus.WithFields(fields).Fatal(msg)
	} else {
		logrus.Fatal(msg)
	}
}

// LogError logs error msg
func LogError(msg string, fields map[string]interface{}) {
	if fields != nil {
		logrus.WithFields(fields).Error(msg)
	} else {
		logrus.Error(msg)
	}
}

// LogWarning logs warning msg
func LogWarning(msg string, fields map[string]interface{}) {
	if fields != nil {
		logrus.WithFields(fields).Warning(msg)
	} else {
		logrus.Warning(msg)
	}
}

// LogInfo logs info msg
func LogInfo(msg string, fields map[string]interface{}) {
	if fields != nil {
		logrus.WithFields(fields).Info(msg)
	} else {
		logrus.Info(msg)
	}
}

// LogDebug logs debug msg
func LogDebug(msg string, fields map[string]interface{}) {
	if fields != nil {
		logrus.WithFields(fields).Debug(msg)
	} else {
		logrus.Debug(msg)
	}
}
