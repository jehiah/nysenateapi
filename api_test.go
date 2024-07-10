package nysenateapi

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBill(t *testing.T) {
	ctx := context.Background()
	a := NewAPI(os.Getenv("NY_SENATE_TOKEN"))
	b, err := a.GetBill(ctx, "2023", "S2304")
	require.NoError(t, err)
	assert.Equal(t, "S2304", b.PrintNo)

	t.Logf("%#v", b)
	b, err = a.GetBill(ctx, "2023", "A1610")
	require.NoError(t, err)
	t.Logf("%#v", b)
}
