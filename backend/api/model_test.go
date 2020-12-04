package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringPatch_UnmarshalJSON(t *testing.T) {
	t.Run("explicitly null", func(t *testing.T) {
		var got struct {
			Test StringPatch `json:"test"`
		}

		err := json.Unmarshal([]byte(`{"test": null}`), &got)
		require.NoError(t, err)

		assert.True(t, got.Test.ExplicitlyNull)
		assert.Nil(t, got.Test.Value)
	})

	t.Run("implicitly null", func(t *testing.T) {
		var got struct {
			Test StringPatch `json:"test"`
		}

		err := json.Unmarshal([]byte(`{}`), &got)
		require.NoError(t, err)

		assert.False(t, got.Test.ExplicitlyNull)
		assert.Nil(t, got.Test.Value)
	})

	t.Run("contains a value", func(t *testing.T) {
		var got struct {
			Test StringPatch `json:"test"`
		}

		err := json.Unmarshal([]byte(`{"test": "hello"}`), &got)
		require.NoError(t, err)

		assert.False(t, got.Test.ExplicitlyNull)
		require.NotNil(t, got.Test.Value)
		assert.Equal(t, "hello", *got.Test.Value)
	})
}
