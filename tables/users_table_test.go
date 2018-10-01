package tables

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"gopkg.in/sqle/sqle.v0/sql"
)

func TestConvertions(t *testing.T) {
	assert := assert.New(t)

	userTable := NewUsersTable()
	rows, err := sql.NodeToRows(userTable)
	assert.Nil(err)
	assert.NotZero(len(rows))

	for _, r := range rows {
		t.Log(r[0], r[1])
	}

}
