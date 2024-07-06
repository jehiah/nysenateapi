package nysenateapi

import (
	"context"
	"os"
	"testing"
)

func TestGetBill(t *testing.T) {
	a := NewAPI(os.Getenv("NY_SENATE_TOKEN"))
	b, err := a.GetBill(context.Background(), "2021", "S5130")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", b)
}
func TestBills(t *testing.T) {
	a := NewAPI(os.Getenv("NY_SENATE_TOKEN"))
	bills, err := a.Bills(context.Background(), "2021", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", bills.Result.Items[0])
}
