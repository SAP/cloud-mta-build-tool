package dir

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetPath(t *testing.T) {

	path := GetPath()
	assert.NotEmpty(t, path)
}

func TestProjectPath(t *testing.T) {
	path := ProjectPath()
	assert.NotEmpty(t, path)
}
