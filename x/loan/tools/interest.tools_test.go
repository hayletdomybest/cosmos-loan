package tools_test

import (
	"loan/tools"
	xtools "loan/x/loan/tools"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateInterest(t *testing.T) {
	res := xtools.CalculateInterest(330000.123, 1, 18)

	res = tools.Round(res, 8)

	require.EqualValues(t, res, 0.00188356)
}
