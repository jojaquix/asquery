package tables

import (
	"container/list"
	"io"

	"asquery/extraction"

	"gopkg.in/sqle/sqle.v0/sql"
)

type ProgramsTable struct {
}

func NewProgramsTable() sql.Table {
	return &ProgramsTable{}
}

func (ProgramsTable) Resolved() bool {
	return true
}

func (ProgramsTable) Name() string {
	//return osVersionTableName
	return "users"
}

func (ProgramsTable) Schema() sql.Schema {
	return sql.Schema{
		{Name: "uid", Type: sql.BigInteger, Nullable: false},
		{Name: "uuid", Type: sql.String, Nullable: false},
		{Name: "username", Type: sql.String, Nullable: false},
	}
}

func (r *ProgramsTable) TransformUp(f func(sql.Node) sql.Node) sql.Node {
	return f(r)
}

func (r *ProgramsTable) TransformExpressionsUp(f func(sql.Expression) sql.Expression) sql.Node {
	return r
}

func (r ProgramsTable) RowIter() (sql.RowIter, error) {

	var iter sql.RowIter = &usersIter{}
	return iter, nil
}

func (ProgramsTable) Children() []sql.Node {
	return []sql.Node{}
}

type programsIter struct {
	fetched bool
	info    list.List
	rowPtr  *list.Element
}

func (iter *programsIter) Next() (sql.Row, error) {

	if iter.rowPtr == nil && !iter.fetched {
		iter.info = extraction.GetPrograms()
		iter.rowPtr = iter.info.Front()
		iter.fetched = true
	}

	if iter.rowPtr != nil {
		r := iter.rowPtr.Value
		iter.rowPtr = iter.rowPtr.Next()
		return programsRowToSqlRow(r.(extraction.Row)), nil
	} else {
		err := io.EOF
		return nil, err
	}

}

func (iter *programsIter) Close() error {
	iter.info.Init()
	iter.rowPtr = nil
	return nil
}

func programsRowToSqlRow(info extraction.Row) sql.Row {

	//TODO why queries only works with Type.BigInteger <-> int64

	return sql.NewRow(
		info["uid"].(int64),
		info["uuid"].(string),
		info["username"].(string),
	)
}
