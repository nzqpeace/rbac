package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

const DefaultConfigFile = "config.json"

var (
	showUsage  bool
	configFile string
	cpuprofile string
)

func init() {
	flag.BoolVar(&showUsage, "h", false, "show help message")
	flag.StringVar(&configFile, "f", "", "config file")

	validate = validator.New()
}

func setLogLevel(level string) {
	l, err := log.ParseLevel(level)
	if err != nil {
		fmt.Printf("parse log level failed, %v, use debug instead\n", err)
		l = log.DebugLevel
	}

	log.SetLevel(l)
}

func setOutput(output string) {
	logPath := path.Clean(output)
	logDir := path.Dir(output)

	var f *os.File
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Printf("mkdir[%s] fail, %v\n", logDir, err)
	} else {
		f, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("open file[%s] failed, %v\n", logPath, err)
		}
	}
	if f == nil {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(f)
	}
}

var validate *validator.Validate

func main() {
	flag.Parse()
	if showUsage {
		flag.Usage()
		return
	}

	if configFile == "" {
		configFile = DefaultConfigFile
	}

	// load config
	config := loadConfig(configFile)
	setLogLevel(config.Log.Level)
	setOutput(config.Log.Output)

	formatter := new(log.TextFormatter)
	formatter.DisableColors = true
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)

	log.WithFields(log.Fields{
		"logfile": config.Log.Output,
		"level":   config.Log.Level,
	}).Info("create log success")

	app := iris.New()
	err := registerRoute(app, config)
	if err != nil {
		log.Error(err)
		return
	}

	startHttpServer(app, config.Http)
}

func startHttpServer(app *iris.Application, config *HttpServerConfig) {
	app.Run(iris.Addr(config.Address))
}
