package main

import (
	"bet25-calendar-sync/global_state"
	"bet25-calendar-sync/helpers"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/s00500/env_logger"
	"github.com/sirupsen/logrus"
)

func main() {
	stop := make(chan os.Signal, 1)

	helpers.CreateFolder("./_logs", "")
	logFile, err := os.OpenFile(fmt.Sprintf("./_logs/log_%s.log", time.Now().Format("2006-01-02")), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// write new line to log file as we start
	logFile.Write([]byte("\n\nStarting up\n"))

	// Create a multi writer
	mw := io.MultiWriter(os.Stdout, logFile)

	// Redirect the custom logger's output
	logger := logrus.New()
	logger.SetOutput(mw)
	debugConfig, _ := os.LookupEnv("LOG")
	if debugConfig == "" {
		debugConfig, _ = os.LookupEnv("GOLANG_LOG")
	}
	logger.Formatter.(*logrus.TextFormatter).EnvironmentOverrideColors = true
	logger.SetOutput(mw)
	log.ConfigureAllLoggers(logger, debugConfig)

	// Load the .env file
	var state *global_state.State
	state = getEnv()

	_=state

	// Stop the program
	<-stop
	log.Info("Exiting")

	// close log file as the last thing
	logFile.Close()
	os.Exit(0)
}