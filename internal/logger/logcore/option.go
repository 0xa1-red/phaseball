package logcore

import "github.com/google/uuid"

type LoggerOpt func(GameLog)

func WithTimestamp() LoggerOpt {
	return func(g GameLog) {
		g.SetWithTimestamp(true)
	}
}

func WithGameID(id uuid.UUID) LoggerOpt {
	return func(g GameLog) {
		g.SetGameID(id)
	}
}
