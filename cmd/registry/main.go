package main

import (
	"fmt"
	"github.com/haggis-io/registry/pkg/repository"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"upper.io/db.v3/postgresql"
	"os/signal"
	"syscall"
)

const (
	name        = "registry"
	description = "A Registry which can store Data format trees."
)

var (
	log          = logrus.WithField("component", name)
	commit       string
	version      string
	addr         string
	databaseURL  string
	databaseUser string
	databasePass string
	port         int
	debugMode    bool
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = description
	app.Version = fmt.Sprintf("%s (%s)", version, commit)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "addr,a",
			Usage:       "TCP address to listen on",
			Value:       "127.0.0.1",
			Destination: &addr,
		},
		cli.StringFlag{
			Name:        "database-url, db-url",
			Usage:       "Database connection URL",
			Value:       "127.0.0.1:5432",
			Destination: &databaseURL,
			EnvVar:      "DATABASE_URL",
		},
		cli.StringFlag{
			Name:        "database-user, db-user",
			Usage:       "Database connection user",
			Destination: &databaseUser,
			EnvVar:      "DATABASE_USER",
		},
		cli.StringFlag{
			Name:        "database-pass, db-pass",
			Usage:       "Database connection password",
			Destination: &databasePass,
			EnvVar:      "DATABASE_PASS",
		},
		cli.IntFlag{
			Name:        "port,p",
			Usage:       "Port to listen on",
			Value:       8080,
			Destination: &port,
		},
		cli.BoolFlag{
			Name:        "debug,d",
			Usage:       "Debug mode",
			Destination: &debugMode,
		},
	}
	app.Action = func(_ *cli.Context) { run() }
	app.Run(os.Args)
}

func run() {

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	log.Infof("Starting %s", name)

	var (
		settings = postgresql.ConnectionURL{
			Host:     databaseURL,
			User:     databaseUser,
			Password: databasePass,
		}
		db, err = postgresql.Open(settings)
	)

	if err != nil {
		log.Fatal(err)
	}

	var (
		ss = repository.NewEntityRepository(db)
	)

	ss.GetEntityByName("", "")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			log.Infof("Exiting %s", name)
			os.Exit(0)
		}
	}

}
