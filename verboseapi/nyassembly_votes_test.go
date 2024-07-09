package verboseapi

import (
	"context"
	"os"
	"testing"
)

func TestAssemblyVotes(t *testing.T) {
	a := NewAPI(os.Getenv("NY_SENATE_TOKEN"))
	ctx := context.Background()
	m, err := a.GetMembers(ctx, "2023", AssemblyChamber)
	if err != nil {
		t.Fatal(err)
	}

	assemblyVotes, err := a.AssemblyVotes(ctx, m, "2021", "A09275")
	if err != nil {
		t.Fatal(err)
	}
	if len(assemblyVotes) != 4 {
		t.Fatalf("expected 4 votes got %d", len(assemblyVotes))
	}
	bill := &Bill{}
	bill.Votes.Items = assemblyVotes

	votes := bill.GetVotes()
	for i, v := range votes {
		if v.MemberID == 0 {
			t.Logf("[%d] unknown member %#v", i, v)
		}
	}

	// this bill has "Held for consideration" votes that should be skipped
	assemblyVotes, err = a.AssemblyVotes(ctx, m, "2023", "A06141")
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("%#v", bill.Votes.Items)
	if len(assemblyVotes) < 1 {
		t.Fatalf("expected 1 votes got %d", len(assemblyVotes))
	}
	if assemblyVotes[0].VoteType != "Held for Consideration" {
		t.Fatalf("expected Held for Consideration got %s", assemblyVotes[0].VoteType)
	}

}
