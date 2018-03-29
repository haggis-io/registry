package main

import (
	"encoding/json"
	"fmt"
	"github.com/haggis-io/registry/pkg/model"
	"github.com/haggis-io/registry/pkg/repository/document"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

const (
	name             = "registry"
	description      = "A Registry which can store Data format trees."
	supportedDb      = "postgres"
	dbName           = "postgres"
	connectionString = "host=%s port=%v user=%s dbname=%s password=%s sslmode=disable"
)

var (
	log                   = logrus.WithField("component", name)
	commit                string
	version               string
	addr                  string
	databaseHost          string
	databasePort          int
	databaseUser          string
	databasePass          string
	migrationFileLocation string
	port                  int
	debugMode             bool
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
			Name:        "database-host, db-host",
			Usage:       "Database connection host",
			Value:       "127.0.0.1",
			Destination: &databaseHost,
			EnvVar:      "DATABASE_HOST",
		},
		cli.IntFlag{
			Name:        "database-port, db-port",
			Usage:       "Database connection port",
			Value:       5432,
			Destination: &databasePort,
			EnvVar:      "DATABASE_PORT",
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
		cli.StringFlag{
			Name:        "migration-file-location, migrate",
			Usage:       "Location of Database migration files",
			Value:       "/migration",
			Destination: &migrationFileLocation,
			EnvVar:      "MIGRATION_FILE_LOCATION",
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

	gormDB, err := gorm.Open(supportedDb, fmt.Sprintf(connectionString, databaseHost, databasePort, databaseUser, dbName, databasePass))
	defer gormDB.Close()

	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(gormDB.DB(), &postgres.Config{})

	log.Infof("Running migrations from %s", migrationFileLocation)

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationFileLocation),
		dbName, driver)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	var (
		documentRepository = &document.DocumentRepository{}
	)

	err = documentRepository.Create(gormDB, &model.Document{
		Description: "some description",
		Author:      "some author",
		Snippet: model.Snippet{
			Text:     "some text",
			TestCase: "s",
		},
		Helper: map[string]interface{}{
			"Name":    "Java",
			"Version": "0.1.0",
		},
	})

	if err != nil {
		panic(err)
	}

	scala := model.Document{
		Description: "some description",
		Author:      "some author",
		Snippet: model.Snippet{
			Text:     "some text",
			TestCase: "s",
		},
		Helper: map[string]interface{}{
			"Name":    "Scala",
			"Version": "0.1.0",
		},
	}

	err = documentRepository.Create(gormDB, &scala)

	if err != nil {
		panic(err)
	}

	err = gormDB.Model(&scala).Association("Dependencies").Append(&model.Document{
		Name:    "Java",
		Version: "0.1.0",
	}).Error

	if err != nil {
		panic(err)
	}

	documents, err := documentRepository.GetDocuments(gormDB, "Scala")

	if err != nil {
		panic(err)
	}

	b, _ := json.Marshal(documents)

	fmt.Println(string(b))

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
