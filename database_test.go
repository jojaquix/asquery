package asquery

import (
	"sort"
	"testing"

	"gopkg.in/sqle/sqle.v0/sql"
	fixtures "gopkg.in/src-d/go-git-fixtures.v3"

	"github.com/stretchr/testify/assert"
	//"gopkg.in/src-d/go-git-fixtures.v3"
)

func init() {
	fixtures.RootFolder = "vendor/gopkg.in/src-d/go-git-fixtures.v3/"
}

const (
	testDBName = "sysInfo"
)

func TestDatabase_Tables(t *testing.T) {
	assert := assert.New(t)

	f := fixtures.Basic().One()
	db := getDB(assert, f, testDBName)

	tables := db.Tables()
	var tableNames []string
	for key := range tables {
		tableNames = append(tableNames, key)
	}

	sort.Strings(tableNames)
	expected := []string{
		osVersionTableName,
		//commitsTableName,
		//referencesTableName,
		//treeEntriesTableName,
		//tagsTableName,
		//blobsTableName,
		//objectsTableName,
	}
	sort.Strings(expected)

	assert.Equal(expected, tableNames)
}

func TestDatabase_Name(t *testing.T) {
	assert := assert.New(t)

	f := fixtures.Basic().One()
	db := getDB(assert, f, testDBName)
	assert.Equal(testDBName, db.Name())
}

func getDB(assert *assert.Assertions, fixture *fixtures.Fixture,
	name string) sql.Database {

	///s, err := filesystem.NewStorage(fixture.DotGit())
	//assert.NoError(err)

	//TODO add sysInfo or something object.
	//r, err := git.Open(s, memfs.New())
	//assert.NoError(err)

	db := NewDatabase(name)
	assert.NotNil(db)

	return db
}

func getTable(assert *assert.Assertions, fixture *fixtures.Fixture,
	name string) sql.Table {

	db := getDB(assert, fixture, "foo")
	assert.NotNil(db)
	assert.Equal(db.Name(), "foo")

	tables := db.Tables()
	table, ok := tables[name]
	assert.True(ok, "table %s does not exist", table)
	assert.NotNil(table)

	return table
}
