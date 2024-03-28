package core

import (
	log "github.com/sirupsen/logrus"
)

var DEBUG = false

func GetLogLevel() log.Level {
	if DEBUG {
		return log.DebugLevel
	}
	return log.InfoLevel
}
