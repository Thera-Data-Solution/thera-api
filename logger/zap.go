package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	var err error
	Log, err = zap.NewProduction() // pakai NewDevelopment() untuk dev
	if err != nil {
		panic(err)
	}
	defer Log.Sync() // flush buffer jika ada
}
