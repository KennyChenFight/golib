package migrationlib

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"

	"github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

type Config struct {
	DatabaseDriver DatabaseDriver
	DatabaseURL    string
	SourceDriver   SourceDriver
	SourceURL      string

	TableName string
}

func newMigrateCmd(c Config) *MigrateCommand {
	return &MigrateCommand{
		config: c,
	}
}

type MigrateCommand struct {
	config Config
}

func (m *MigrateCommand) prepare() (*migrate.Migrate, error) {
	db, err := sql.Open(string(m.config.DatabaseDriver), m.config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	var sourceDriver source.Driver
	switch m.config.SourceDriver {
	case FileDriver:
		sourceDriver = &file.File{}
	default:
		return nil, errors.New("not supported source driver")
	}
	sourceDriver, err = sourceDriver.Open(m.config.SourceURL)
	if err != nil {
		return nil, err
	}

	var databaseDriver database.Driver
	switch m.config.DatabaseDriver {
	case PostgresDriver:
		databaseDriver, err = postgres.WithInstance(db, &postgres.Config{MigrationsTable: m.config.TableName})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("not supported database driver")
	}
	mi, err := migrate.NewWithInstance(string(m.config.SourceDriver),
		sourceDriver,
		string(m.config.DatabaseDriver),
		databaseDriver)
	if err != nil {
		return nil, err
	}
	return mi, nil
}

func (m MigrateCommand) Up() error {
	mi, err := m.prepare()
	if err != nil {
		return err
	}
	return mi.Up()
}

func (m MigrateCommand) UpTo(limit uint) error {
	mi, err := m.prepare()
	if err != nil {
		return err
	}
	return mi.Steps(int(limit))
}

func (m MigrateCommand) Down() error {
	mi, err := m.prepare()
	if err != nil {
		return err
	}
	return mi.Down()
}

func (m MigrateCommand) DownTo(limit int) error {
	mi, err := m.prepare()
	if err != nil {
		return err
	}
	if limit > 0 {
		return errors.New("limit should be less than 0")
	}
	return mi.Steps(limit)
}

func (m MigrateCommand) Drop() error {
	mi, err := m.prepare()
	if err != nil {
		return err
	}
	return mi.Drop()
}

func (m MigrateCommand) GoTo(version uint) error {
	mi, err := m.prepare()
	if err != nil {
		return err
	}
	return mi.Migrate(version)
}

func (m MigrateCommand) Force(version uint) error {
	mi, err := m.prepare()
	if err != nil {
		return err
	}
	return mi.Force(int(version))
}

func (m MigrateCommand) CurrentVersion() (uint, bool, error) {
	mi, err := m.prepare()
	if err != nil {
		return 0, false, err
	}
	return mi.Version()
}

func (m MigrateCommand) PrintUsageInfo() {
	log.Println(`
Usage: migrate OPTIONS COMMAND [arg...]
migrate [ -version | -help ]

Options:
  -source          Location of the migrations (driver://url)
  -path            Shorthand for -source=file://path
  -database        Run migrations against this database (driver://url)
  -prefetch N      Number of migrations to load in advance before executing (default 10)
  -lock-timeout N  Allow N seconds to acquire database lock (default 15)
  -verbose         Print verbose logging
  -version         Print version
  -help            Print usage

Commands:
  create [-ext E] [-dir D] [-seq] [-digits N] [-format] NAME
	   Create a set of timestamped up/down migrations titled NAME, in directory D with extension E.
	   Use -seq option to generate sequential up/down migrations with N digits.
	   Use -format option to specify a Go time format string. Note: migrations with the same time cause "duplicate migration version" error. 

  goto V       Migrate to version V
  up [N]       Apply all or N up migrations
  down [N]     Apply all or N down migrations
  drop [-f] [-all]    Drop everything inside database
	Use -f to bypass confirmation
	Use -all to apply all down migrations
  force V      Set version V but don't run migration (ignores dirty state)
  version      Print current migration version

Source drivers: godoc-vfs, gcs, file, bitbucket, gitlab, github-ee, go-bindata, s3, github
Database drivers: firebirdsql, mysql, redshift, sqlserver, stub, spanner, cockroachdb, crdb-postgres, firebird, mongodb, postgres, postgresql, cassandra, clickhouse, cockroach, mongodb+srv, neo4j`)
}
