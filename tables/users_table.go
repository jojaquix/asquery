package tables

import (
	"container/list"
	"io"

	"asquery/extraction"

	"gopkg.in/sqle/sqle.v0/sql"
)

type UsersTable struct {
}

func NewUsersTable() sql.Table {
	return &UsersTable{}
}

func (UsersTable) Resolved() bool {
	return true
}

func (UsersTable) Name() string {
	//return osVersionTableName
	return "users"
}

func (UsersTable) Schema() sql.Schema {
	return sql.Schema{
		{Name: "uid", Type: sql.BigInteger, Nullable: false},
		{Name: "gid", Type: sql.BigInteger, Nullable: false},
		{Name: "uid_signed", Type: sql.BigInteger, Nullable: false},
		{Name: "gid_signed", Type: sql.BigInteger, Nullable: false},
		{Name: "username", Type: sql.String, Nullable: true},
		{Name: "description", Type: sql.String, Nullable: true},
		{Name: "directory", Type: sql.String, Nullable: true},
		{Name: "shell", Type: sql.String, Nullable: true},
		{Name: "uuid", Type: sql.String, Nullable: false},
		{Name: "type", Type: sql.String, Nullable: false},
	}
}

func (r *UsersTable) TransformUp(f func(sql.Node) sql.Node) sql.Node {
	return f(r)
}

func (r *UsersTable) TransformExpressionsUp(f func(sql.Expression) sql.Expression) sql.Node {
	return r
}

func (r UsersTable) RowIter() (sql.RowIter, error) {

	var iter sql.RowIter = &usersIter{}
	return iter, nil
}

func (UsersTable) Children() []sql.Node {
	return []sql.Node{}
}

type usersIter struct {
	fetched bool
	info    list.List
	rowPtr  *list.Element
}

func (iter *usersIter) Next() (sql.Row, error) {

	//test return just thow users
	//change to actual extractor call
	if iter.rowPtr == nil && !iter.fetched {
		iter.info = extraction.GetUsers()
		iter.rowPtr = iter.info.Front()
		iter.fetched = true
	}

	if iter.rowPtr != nil {
		r := iter.rowPtr.Value
		iter.rowPtr = iter.rowPtr.Next()
		return userInfoToRow(r.(extraction.Row)), nil
	} else {
		err := io.EOF
		return nil, err
	}

}

func (iter *usersIter) Close() error {
	iter.info.Init()
	iter.rowPtr = nil
	return nil
}

func userInfoToRow(info extraction.Row) sql.Row {

	//TODO why queries only works with Type.BigInteger <-> int64

	return sql.NewRow(
		info["uid"].(int64),
		info["gid"].(int64),
		info["uid_signed"].(int64),
		info["gid_signed"].(int64),
		info["username"].(string),
		info["description"].(string),
		info["directory"].(string),
		info["shell"].(string),
		info["uuid"].(string),
		info["type"].(string),
	)
}
