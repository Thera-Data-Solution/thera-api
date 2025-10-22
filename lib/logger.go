package lib

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger sekarang adalah SugaredLogger untuk sintaks yang lebih mudah
var Logger *zap.SugaredLogger

// InitLogger menginisialisasi Sugared Logger Zap.
func InitLogger() {
	var err error
	// NewDevelopment() tetap digunakan untuk output di konsol yang mudah dibaca
	baseLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Zap logger: %v", err))
	}

	// Konversi baseLogger ke SugaredLogger
	Logger = baseLogger.Sugar()
}
