package main

import (
	"fmt"
	"github.com/haggis-io/registry/pkg/api"
	"github.com/haggis-io/registry/pkg/repository"
	"github.com/haggis-io/registry/pkg/server"
	"github.com/haggis-io/registry/pkg/service"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strconv"
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

	var (
		address  = net.JoinHostPort(addr, strconv.Itoa(port))
		lis, err = net.Listen("tcp", address)
		srv      = grpc.NewServer()
	)

	if err != nil {
		log.Fatal(err)
	}

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
		documentRepository = repository.NewDocumentRepository()
		documentService    = service.NewRegistryService(gormDB, documentRepository)
		documentServer     = server.NewRegistryServer(documentService)
	)

	api.RegisterRegistryServer(srv, documentServer)

	go func() {
		log.Infof("GRPC server listening on %s", address)
		if err := srv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			log.Infof("Exiting %s", name)
			srv.GracefulStop()
			os.Exit(0)
		}
	}

}
