package logs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/phuslu/log"
	"github.com/robfig/cron/v3"
	"github.com/vanilla-os/sdk/pkg/v1/app/types"
	logsTypes "github.com/vanilla-os/sdk/pkg/v1/logs/types"
)

func getLogPath() (string, error) {
	var logPath string

	// if running as root, we should use /var/vlogs/, while users should store
	// logs in their home directory (e.g. ~/.vlogs)
	if os.Geteuid() == 0 {
		logPath = "/var/vlogs/"
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %v", err)
		}
		logPath = filepath.Join(homeDir, ".vlogs")
	}

	// we have to create the directory if it doesn't exist
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.MkdirAll(logPath, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create log directory: %v", err)
		}
	}

	return logPath, nil
}

// NewLogger creates a new logger for the application, each logger has
// a file logger and a console logger. The file logger is used to log
// to the vlogs directory, while the console logger is used to log to
// the console.
//
// Example:
//
//	logger, err := logs.NewLogger(app)
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	logger.File.Info().Msg("Batman reached the file logger")
//	logger.Console.Info().Msg("Batman reached the console logger")
//
// since we use structured logging, we can also log with fields:
//
//	logger.File.Info().Str("where", "file").Msg("Batman is saving Gotham")
//	logger.Console.Info().Str("where", "console").Msg("Batman is saving Gotham")
func NewLogger(app *types.App) (logsTypes.Logger, error) {
	vLogger := logsTypes.Logger{}

	// preparing the file logger
	logPath, err := getLogPath()
	if err != nil {
		return vLogger, err
	}

	vLogFile := filepath.Join(logPath, app.Sign, "log.json")

	vLogger.File = log.Logger{
		Level: log.ParseLevel("info"),
		Writer: &log.FileWriter{
			Filename:     vLogFile,
			FileMode:     0600,
			MaxSize:      500 * 1024 * 1024,
			MaxBackups:   7,
			EnsureFolder: true,
			LocalTime:    true,
			TimeFormat:   "15:04:05",
			Cleaner: func(filename string, maxBackups int, matches []os.FileInfo) {
				var dir = filepath.Dir(filename)
				for i, fi := range matches {
					filename := filepath.Join(dir, fi.Name())
					switch {
					case i > maxBackups:
						os.Remove(filename)
					case !strings.HasSuffix(filename, ".gz"):
						go exec.Command("nice", "gzip", filename).Run()
					}
				}
			},
		},
	}

	// setting up the rotation for the file logger
	runner := cron.New(cron.WithLocation(time.Local))
	runner.AddFunc("0 0 * * *", func() { vLogger.File.Writer.(*log.FileWriter).Rotate() })
	go runner.Run()

	// preparing the console logger
	vLogger.Console = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}

	return vLogger, nil
}
