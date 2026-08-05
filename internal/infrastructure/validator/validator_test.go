package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JrNovaEX/DCB/internal/infrastructure/validator"
)

func TestNew(t *testing.T) {
	v, err := validator.New()
	require.NoError(t, err)
	assert.NotNil(t, v)
}

func TestValidator_Struct(t *testing.T) {
	type sample struct {
		Name  string `validate:"required"`
		Count int    `validate:"min=1"`
	}

	v, err := validator.New()
	require.NoError(t, err)

	t.Run("valid struct", func(t *testing.T) {
		s := sample{Name: "test", Count: 1}
		assert.NoError(t, v.Struct(s))
	})

	t.Run("missing required field", func(t *testing.T) {
		s := sample{Count: 1}
		err := v.Struct(s)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Name")
	})

	t.Run("min violation", func(t *testing.T) {
		s := sample{Name: "test", Count: 0}
		err := v.Struct(s)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Count")
	})
}
