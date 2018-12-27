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
	return "programs"
}

func (ProgramsTable) Schema() sql.Schema {
	return sql.Schema{
		{Name: "name", Type: sql.String, Nullable: false},
		{Name: "version", Type: sql.String, Nullable: true},
		{Name: "install_location", Type: sql.String, Nullable: true},
		{Name: "install_source", Type: sql.String, Nullable: true},
		{Name: "install_source", Type: sql.String, Nullable: true},
		{Name: "language", Type: sql.String, Nullable: true},
		{Name: "publisher", Type: sql.String, Nullable: true},
		{Name: "uninstall_string", Type: sql.String, Nullable: true},
		{Name: "install_date", Type: sql.String, Nullable: true},
		{Name: "identifying_number", Type: sql.String, Nullable: true},
	}
}

func (r *ProgramsTable) TransformUp(f func(sql.Node) sql.Node) sql.Node {
	return f(r)
}

func (r *ProgramsTable) TransformExpressionsUp(f func(sql.Expression) sql.Expression) sql.Node {
	return r
}

func (r ProgramsTable) RowIter() (sql.RowIter, error) {

	var iter sql.RowIter = &programsIter{}
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

	return sql.NewRow(
		info["name"].(string),
		info["version"].(string),
		info["install_location"].(string),
		info["install_source"].(string),
		info["install_source"].(string),
		info["language"].(string),
		info["publisher"].(string),
		info["uninstall_string"].(string),
		info["install_date"].(string),
		info["identifying_number"].(string),
	)
}
