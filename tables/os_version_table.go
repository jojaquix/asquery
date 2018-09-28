package tables

import (
	"strconv"
	"strings"

	"asquery/extraction"

	"gopkg.in/sqle/sqle.v0/sql"
)

type OsVersionTable struct {
	//r *git.Repository
}

func NewOsVersionTable() sql.Table {
	return &OsVersionTable{}
}

func (OsVersionTable) Resolved() bool {
	return true
}

func (OsVersionTable) Name() string {
	//return osVersionTableName
	return "os_version"
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

	var iter sql.RowIter = &osVersionInfoIter{osExtractor: extraction.NewOsExtractor()}
	var err error

	if err != nil {
		return nil, err
	}

	return iter, nil
}

func (OsVersionTable) Children() []sql.Node {
	return []sql.Node{}
}

type osVersionInfoIter struct {
	osExtractor extraction.OsExtractor
}

func (iter *osVersionInfoIter) Next() (sql.Row, error) {

	osVersionInfo, err := iter.osExtractor.Next()
	if err != nil {
		return nil, err
	}

	return osVersionInfoToRow(osVersionInfo), nil
}

func (iter *osVersionInfoIter) Close() error {
	return nil
}

func osVersionInfoToRow(info *extraction.OsVersionInfo) sql.Row {

	version_splited := strings.Split(info.Version, ".")
	var version, major, minor, patch, build string = "", "0", "0", "0", "0"

	switch len(version_splited) {
	case 3:
		build = version_splited[2]
		fallthrough
	case 2:
		minor = version_splited[1]
		fallthrough
	case 1:
		major = version_splited[0]
	default:
		break
	}

	version = major + "." + minor + "." + build

	imajor, _ := strconv.Atoi(major)
	iminor, _ := strconv.Atoi(minor)
	ipatch, _ := strconv.Atoi(patch)

	return sql.NewRow(
		info.Name,
		version,
		int32(imajor),
		int32(iminor),
		int32(ipatch),
		build,
		"windows",
		"windows",
		nil, //codename not implemented not compatible with windows xp
	)
}
