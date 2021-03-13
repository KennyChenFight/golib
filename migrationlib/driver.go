package migrationlib

type SourceDriver string
type DatabaseDriver string

const (
	FileDriver SourceDriver = "file"
)

const (
	PostgresDriver DatabaseDriver = "postgres"
)
