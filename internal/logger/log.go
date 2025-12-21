package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

// NewLogger creates a zap logger based on the mode ("dev" or "prod").
func NewLogger(mode string) (*Logger, error) {
	var z *zap.Logger
	var err error

	if mode == "dev" {
		z, err = zap.NewDevelopment()
	} else {
		z, err = zap.NewProduction()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return &Logger{Logger: z}, nil
}
