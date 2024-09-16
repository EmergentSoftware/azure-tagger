package main

import "log"

type Logger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}
