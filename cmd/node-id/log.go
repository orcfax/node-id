package main

import (
	"fmt"
	"log"
	"time"
)

type logWriter struct{}

// Write enables us to format a logging prefix for the application. The
// text will appear before the log message output by the caller.
func (lw *logWriter) Write(logString []byte) (int, error) {
	// 2023-11-27 11:36:57 ERROR :: validator_node.py:100:main() :: uvicorn import string not correctly configured: api (No module named 'api')
	return fmt.Print(
		time.Now().UTC().Format(timeFormat),
		fmt.Sprintf(" :: %s :: ", appname),
		string(logString),
	)
}

func init() {
	// Configure logging to use a custom log writer with Orcfax defaults.
	log.SetFlags(0 | log.Lshortfile | log.LUTC)
	log.SetOutput(new(logWriter))
}
