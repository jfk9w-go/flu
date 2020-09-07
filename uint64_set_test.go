package flu

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint64Set(t *testing.T) {
	set := make(Uint64Set)
	assert.False(t, set.Has(1))
	set.Add(1)
	assert.True(t, set.Has(1))
	bytes, err := json.Marshal(set)
	assert.Nil(t, err)
	set.Delete(1)
	assert.False(t, set.Has(1))
	set = make(Uint64Set)
	assert.Nil(t, json.Unmarshal(bytes, &set))
	assert.True(t, set.Has(1))
	copy := set.Copy()
	assert.True(t, copy.Has(1))
	set.Delete(1)
	assert.False(t, set.Has(1))
	assert.True(t, copy.Has(1))
}
