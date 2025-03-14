package expr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapr/dapr/pkg/expr"
)

func TestEval(t *testing.T) {
	var e expr.Expr
	code := `(has(input.test) && input.test == 1234) || (has(result.test) && result.test == 5678)`
	err := e.DecodeString(code)
	require.NoError(t, err)
	assert.Equal(t, code, e.String())
	result, err := e.Eval(map[string]interface{}{
		"input": map[string]interface{}{
			"test": 1234,
		},
		"result": map[string]interface{}{
			"test": 5678,
		},
	})
	require.NoError(t, err)
	assert.True(t, result.(bool))
}

func TestJSONMarshal(t *testing.T) {
	var e expr.Expr
	exprBytes := []byte(`"(has(input.test) && input.test == 1234) || (has(result.test) && result.test == 5678)"`)
	err := e.UnmarshalJSON(exprBytes)
	require.NoError(t, err)
	assert.Equal(t, `(has(input.test) && input.test == 1234) || (has(result.test) && result.test == 5678)`, e.Expr())
	_, err = e.MarshalJSON()
	require.NoError(t, err)
}

func TestEmptyProgramNoPanic(t *testing.T) {
	var e expr.Expr
	r, err := e.Eval(map[string]interface{}{})

	assert.Nil(t, r)
	require.Error(t, err)
}

var result interface{}

func BenchmarkEval(b *testing.B) {
	var e expr.Expr
	err := e.DecodeString(`(has(input.test) && input.test == 1234) || (has(result.test) && result.test == 5678)`)
	require.NoError(b, err)
	data := map[string]interface{}{
		"input": map[string]interface{}{
			"test": 1234,
		},
		"result": map[string]interface{}{
			"test": 5678,
		},
	}
	var r interface{}
	for range b.N {
		r, _ = e.Eval(data)
	}
	result = r
}
