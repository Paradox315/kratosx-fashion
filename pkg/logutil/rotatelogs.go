package logutil

import (
	"os"
	"path"
	"time"

	"go.uber.org/zap/zapcore"

	zapRotateLogs "github.com/lestrrat-go/file-rotatelogs"
)

func GetWriteSyncer(dir string, logInConsole bool) (zapcore.WriteSyncer, error) {
	fileWriter, err := zapRotateLogs.New(
		path.Join(dir, "%Y-%m-%d.log"),
		zapRotateLogs.WithMaxAge(7*24*time.Hour),
		zapRotateLogs.WithRotationTime(24*time.Hour),
	)

	if logInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
