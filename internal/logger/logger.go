package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() error {
	file, err := os.OpenFile("letStayInn.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	writeSyncer := zapcore.AddSync(file)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)
	Log = zap.New(core)
	return nil
}

func SyncLogger() {
	_ = Log.Sync()
}
