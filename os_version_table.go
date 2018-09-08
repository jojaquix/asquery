package gitquery

import (
	"io"

	"gopkg.in/sqle/sqle.v0/sql"
	//"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type OsVersionTable struct {
	//r *git.Repository
}

func newOsVersionTable() sql.Table {
	return &OsVersionTable{}
}

func (OsVersionTable) Resolved() bool {
	return true
}

func (OsVersionTable) Name() string {
	return osVersionTableName
}

func (OsVersionTable) Schema() sql.Schema {
	return sql.Schema{
		{Name: "name", Type: sql.String, Nullable: false},
		{Name: "version", Type: sql.String, Nullable: false},
		{Name: "major", Type: sql.Integer, Nullable: true},
		{Name: "minor", Type: sql.Integer, Nullable: true},
		{Name: "patch", Type: sql.Integer, Nullable: true},
		{Name: "build", Type: sql.String, Nullable: true},
		{Name: "platform", Type: sql.String, Nullable: true},
		{Name: "platform_like", Type: sql.String, Nullable: true},
		{Name: "codename", Type: sql.String, Nullable: true},
	}
}

func (r *OsVersionTable) TransformUp(f func(sql.Node) sql.Node) sql.Node {
	return f(r)
}

func (r *OsVersionTable) TransformExpressionsUp(f func(sql.Expression) sql.Expression) sql.Node {
	return r
}

func (r OsVersionTable) RowIter() (sql.RowIter, error) {

	var iter sql.RowIter = &OsVersionInfoIter{}
	var err error

	if err != nil {
		return nil, err
	}

	return iter, nil
}

func (OsVersionTable) Children() []sql.Node {
	return []sql.Node{}
}

//TODO going to other file ??

//OsVersionInfo ...
type OsVersionInfo struct {
	name    string
	version string
}

//OsVersinInfoIterable ... is a generic closable interface for iterating over OsVersionInfo
type OsVersionInfoIterable interface {
	Next() (*OsVersionInfo, error)
	ForEach(func(*OsVersionInfo) error) error
	Close()
}

type OsVersionInfoIter struct {
	//xname string
	count int
}

func (i *OsVersionInfoIter) Next() (sql.Row, error) {
	var osVersionInfo = &OsVersionInfo{name: "jhonjamesOs", version: "1"}
	var err error
	//if err != nil {
	if i.count > 0 {
		err = io.EOF
		return nil, err
	}

	i.count++
	return osVersionInfoToRow(osVersionInfo), nil
}

func (i *OsVersionInfoIter) Close() error {
	//i.i.Close()
	return nil
}

func osVersionInfoToRow(info *OsVersionInfo) sql.Row {
	return sql.NewRow(
		info.name,
		info.version,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}
