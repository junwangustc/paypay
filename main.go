package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logfilePath string
	configPath  string
	globalLog   *log.Logger
)

func init() {
	logfilePath = flag.String("logfile", "/tmp/alipay.log", "alipay logfile path")
	configPath = flag.String("config", "/tmp/alipay.conf", "alipay config file")

	flag.Parse()
}

func initLog() {
	output := &lumberjack.Logger{
		Filename:   logfilePath,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     7,
	}
	log.SetOutput(output)
	globalLog = log.New(output, "", log.Ldefault)

}

func main() {
	initLog()
	cfg, err := ParseConfig(configPath)
	if err != nil {
		log.Fatalf("parseconfig error  start: %v", err)
	}
	srv := NewServer(cfg)
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("alipay start: %v", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	log.Printf("listening for signals ...")
	select {
	case <-signalCh:
		log.Printf("signal received, shutdown...")
	}
}
