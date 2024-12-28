package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	camect_go "github.com/kaiiorg/camect-go"
)

var (
	logLevels = map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
	}
)

var (
	json  = flag.Bool("json", false, "print JSON logs")
	level = flag.String("level", "INFO", "log level: DEBUG, INFO, WARN, ERROR")

	ip       = flag.String("ip", "0.0.0.0", "ip address of hub")
	username = flag.String("username", "admin", "username of hub local admin")
	password = flag.String("password", "this isn't a real password, provide your own", "password of hub local admin")

	mode = flag.String("mode", "", "if set, will attempt to set hub to given mode (HOME or DEFAULT only)")
)

func main() {
	flag.Parse()
	initSlog()

	slog.Info("camect-go", "goVersion", runtime.Version())

	hub := camect_go.New(
		*ip,
		*username,
		*password,
		nil, // logger, nil means use global default
	)

	info, err := hub.Info()
	if err != nil {
		panic(err)
	}

	slog.Info(
		"got hub info",
		"hubName", info.Name,
		"hubId", info.Id,
		"mode", info.Mode,
	)

	ctx, ctxCancel := context.WithCancel(context.Background())
	eventsChan, err := hub.Events(ctx, 1)
	if err != nil {
		slog.Error("failed to listen for events", "error", err.Error())
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go setHubMode(hub)

	for {
		select {
		case sig := <-signalChan:
			slog.Info("got signal to exit", "signal", sig)
			ctxCancel()
			time.Sleep(time.Second)
		case e := <-eventsChan.AlertChan:
			slog.Info("got alert", "data", fmt.Sprintf("%#v", e))
		case e := <-eventsChan.ModeChangeChan:
			slog.Info("got mode changed event", "data", fmt.Sprintf("%#v", e))
		case e := <-eventsChan.AlertDisabledChan:
			slog.Info("got alert disabled event", "data", fmt.Sprintf("%#v", e))
		case e := <-eventsChan.AlertEnabledChan:
			slog.Info("got alert enabled event", "data", fmt.Sprintf("%#v", e))
		case e := <-eventsChan.CameraOnlineChan:
			slog.Info("got camera online event", "data", fmt.Sprintf("%#v", e))
		case e := <-eventsChan.CameraOfflineChan:
			slog.Info("got camera offline event", "data", fmt.Sprintf("%#v", e))
		case e := <-eventsChan.UnknownEventChan:
			slog.Warn("got unknown event", "raw", string(e))
		}

		if ctx.Err() != nil {
			break
		}
	}
}

func initSlog() {
	level, ok := logLevels[strings.ToUpper(*level)]
	if !ok {
		panic("invalid log level")
	}

	options := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if *json {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	slog.SetDefault(slog.New(handler))
}

func setHubMode(hub *camect_go.Hub) {
	var newMode camect_go.HubMode

	switch *mode {
	case "HOME":
		newMode = camect_go.ModeHome
	case "DEFAULT":
		newMode = camect_go.ModeDefault
	default:
		return
	}

	err := hub.SetMode(newMode, "CLI Example")
	if err != nil {
		slog.Error("failed to set hub mode", "err", err.Error())
	}
}
