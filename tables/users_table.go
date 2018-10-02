package tables

import (
	"container/list"
	"io"
	"strconv"

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
		{Name: "uuid", Type: sql.BigInteger, Nullable: true},
		{Name: "username", Type: sql.String, Nullable: false},
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
		var err error
		iter.info, err = extraction.GetUsers()
		if err == nil {
			iter.rowPtr = iter.info.Front()
			iter.fetched = true
		}
		//else {
		//	//TODO logging
		//}

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
	imajor, _ := strconv.Atoi(string(info["uuid"]))

	return sql.NewRow(
		int64(imajor),
		string(info["username"]),
	)
}
