package collection_test

import (
	"testing"

	"github.com/becheran/roumon/internal/collection"
	"github.com/stretchr/testify/assert"
)

func Test_SliceContains_Contains(t *testing.T) {
	assert.True(t, collection.SliceContains([]string{"foo", "bar"}, "bar"))
}

func Test_SliceContains_DoNotContain(t *testing.T) {
	assert.False(t, collection.SliceContains([]string{"foo", "bar"}, "baz"))
}
