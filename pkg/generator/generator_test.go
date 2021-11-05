package generator

import (
	"context"
	"strings"
	"testing"

	"github.com/operator-framework/combo/pkg/combination"
	"github.com/operator-framework/combo/test/assets/generatorTestData"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expected struct {
	err        error
	evaluation string
}

func TestEvaluate(t *testing.T) {
	for _, tt := range []struct {
		name         string
		template     string
		combinations combination.Stream
		expected     expected
	}{
		{
			name:     "can process a template",
			template: generatorTestData.EvaluateInput,
			expected: expected{
				err:        nil,
				evaluation: generatorTestData.EvaluateOutput,
			},
			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"NAMESPACE": {"foo", "bar"},
					"NAME":      {"baz"},
				}),
				combination.WithSolveAhead(),
			),
		},
		{
			name:     "processes an empty template",
			template: ``,
			expected: expected{
				err:        nil,
				evaluation: ``,
			},
			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"NAMESPACE": {"foo", "bar"},
					"NAME":      {"baz"},
				}),
				combination.WithSolveAhead(),
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			evaluation, err := Evaluate(ctx, tt.template, tt.combinations)
			if err != nil {
				t.Fatal("received an error during evaluation:", err)
			}

			assert.Equal(t, tt.expected.err, err)
			require.ElementsMatch(
				t,
				strings.Split(string(tt.expected.evaluation), "---"),
				strings.Split(string(evaluation), "---"),
				"Document evaluations generated incorrectly")
		})
	}
}