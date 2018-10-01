package extraction

import (
	"container/list"
)

type Data string
type Row map[string]Data
type ColumnNames []string
type QueryData list.List

//InfoIterable ... is a generic closable interface for iterating over Infos
type InfoIterable interface {
	Next() (*InfoIterable, error)
	ForEach(func(*InfoIterable) error) error
	Close()
}
