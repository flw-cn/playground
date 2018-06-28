package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoarding(t *testing.T) {
	var err error

	t.Run("Testing Boarding()", func(t *testing.T) {
		err = Boarding("", "")
		assert.Error(t, err)

		err = Boarding("foo/bar", "/hello/world")
		assert.Error(t, err)

		err = Boarding("/foo/bar", "hello/world")
		assert.Error(t, err)

		err = Boarding("foo/bar", "hello/world")
		assert.Error(t, err)

		err = Boarding("/foo/bar", "/hello/world")
		assert.NoError(t, err)

		err = Boarding("/", "/")
		assert.NoError(t, err)
	})

	t.Run("Testing boarding()", func(t *testing.T) {
		Boarding("/", "/")

		path, err := boarding("/tmp/code.piece")
		assert.NoError(t, err)
		assert.Equalf(t, path, "/tmp/code.piece", "")

		Boarding("/foo/bar", "/hello/world")

		path, err = boarding("/foo/bar/tmp/code.piece")
		assert.NoError(t, err)
		assert.Equalf(t, path, "/hello/world/tmp/code.piece", "")
	})
}
