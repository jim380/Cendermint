package types

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppConfig struct {
	Chain                string
	RestAddr             string
	RpcAddr              string
	ListeningPort        string
	OperAddr             string
	LogOutput            string
	PollIntervalChain    string
	PollIntervalAsync    string
	PollIntervalBackfill string
	LastUpdatedMoreThan  string
	LogLevel             zapcore.Level
	Logger               *zap.Logger
}
