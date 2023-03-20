package funcutils_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xeptore/wireuse/pkg/funcutils"
)

func TestMap(t *testing.T) {
	result := funcutils.Map([]int{1, 2, 3, 4, 5, 6, 7}, func(i int) string {
		return strconv.Itoa(i)
	})

	require.Equal(t, []string{"1", "2", "3", "4", "5", "6", "7"}, result)
	require.Len(t, result, 7)
	require.Equal(t, 7, cap(result), "expected mapped result cap to be 7")
}
