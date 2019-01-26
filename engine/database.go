package engine

import (
	"asquery/tables"
	"gopkg.in/sqle/sqle.v0/sql"
	//"gopkg.in/src-d/go-git.v4"
)

const (
// TODO 'references' is a reserved keyword into the parser
//osVersionTableName = "os_version"
//referencesTableName  = "refs"
//commitsTableName     = "commits"
//tagsTableName        = "tags"
//blobsTableName       = "blobs"
//treeEntriesTableName = "tree_entries"
//objectsTableName     = "objects"
)

type Database struct {
	name           string
	osVersionTable sql.Table
	usersTable     sql.Table
	programsTable  sql.Table
}

func NewDatabase(name string) sql.Database {
	return &Database{
		name:           name,
		osVersionTable: tables.NewOsVersionTable(),
		usersTable:     tables.NewUsersTable(),
		programsTable:  tables.NewProgramsTable(),
	}
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) Tables() map[string]sql.Table {
	return map[string]sql.Table{
		d.osVersionTable.Name(): d.osVersionTable,
		d.usersTable.Name():     d.usersTable,
		d.programsTable.Name():  d.programsTable,
	}
}
