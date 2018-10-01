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
	//cr   sql.Table
	//tr   sql.Table
	//rr   sql.Table
	//ter  sql.Table
	//br   sql.Table
	//or   sql.Table
}

func NewDatabase(name string) sql.Database {
	return &Database{
		name:           name,
		osVersionTable: tables.NewOsVersionTable(),
		usersTable:     tables.NewUsersTable(),
		//cr:   newCommitsTable(r),
		//rr:   newReferencesTable(r),
		//tr:   newTagsTable(r),
		//br:   newBlobsTable(r),
		//ter:  newTreeEntriesTable(r),
		//or:   newObjectsTable(r),
	}
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) Tables() map[string]sql.Table {
	return map[string]sql.Table{
		d.osVersionTable.Name(): d.osVersionTable,
		d.usersTable.Name():     d.usersTable,
		//commitsTableName:     d.cr,
		//tagsTableName:        d.tr,
		//referencesTableName:  d.rr,
		//blobsTableName:       d.br,
		//treeEntriesTableName: d.ter,
		//objectsTableName:     d.or,
	}
}
