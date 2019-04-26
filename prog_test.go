package wrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	p := New()
	assert.Equal(Development, p.Env)
}

func TestIsMethods(t *testing.T) {
	assert := assert.New(t)
	p := New()
	p.Env = Development
	assert.Equal(true, p.IsDevelopment())
	assert.Equal(false, p.IsTesting())
	assert.Equal(false, p.IsProduction())
	p.Env = Testing
	assert.Equal(false, p.IsDevelopment())
	assert.Equal(true, p.IsTesting())
	assert.Equal(false, p.IsProduction())
	p.Env = Production
	assert.Equal(false, p.IsDevelopment())
	assert.Equal(false, p.IsTesting())
	assert.Equal(true, p.IsProduction())
}
