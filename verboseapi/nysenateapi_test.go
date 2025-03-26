package verboseapi

import (
	"context"
	"os"
	"testing"
	"time"
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

func TestBillUpdates(t *testing.T) {
	a := NewAPI(os.Getenv("NY_SENATE_TOKEN"))
	from := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	updates, err := a.GetBillUpdates(context.Background(), from, from.Add(time.Hour*24*7), 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", updates.Result.Items[0])
}
