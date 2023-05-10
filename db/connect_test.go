package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitMySQL(t *testing.T) {
	connectDB := NewConnectDB("root", "", "127.0.0.1", "3306", "sawer")
	_, err := connectDB.InitMySQL()
	assert.Nil(t, err)

}
