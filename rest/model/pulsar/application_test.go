package pulsar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBBPulsarApplication(t *testing.T) {
	a := NewApplication("app_name")
	assert.Equal(t, "app_name", a.Name, "Wrong name")
}
