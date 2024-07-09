package nysenateapi

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jehiah/nysenateapi/verboseapi"
)

type API struct {
	api *verboseapi.NYSenateAPI

	mutex           sync.Mutex
	assemblyMembers map[int][]verboseapi.MemberEntry
}

func NewAPI(token string) *API {
	return &API{
		api:             verboseapi.NewAPI(token),
		assemblyMembers: make(map[int][]verboseapi.MemberEntry),
	}
}

func (a *API) GetBill(ctx context.Context, session, printNo string) (*Bill, error) {
	bill, err := a.api.GetBill(ctx, session, printNo)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, nil
	}

	out := newBill(bill)

	if out.Chamber == "ASSEMBLY" {
		members, err := a.getAssemblyMembers(ctx, bill.Session)
		if err != nil {
			return nil, err
		}
		votes, err := a.api.AssemblyVotes(ctx, members, session, printNo)
		if err != nil {
			return nil, err
		}
		out.Votes = append(out.Votes, newVotes(votes)...)
	}

	return out, nil
}

func (a *API) getAssemblyMembers(ctx context.Context, session int) ([]verboseapi.MemberEntry, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if members, ok := a.assemblyMembers[session]; ok {
		return members, nil
	}
	sessionStr := fmt.Sprintf("%d", session)
	members, err := a.api.GetMembers(ctx, sessionStr, verboseapi.AssemblyChamber)
	if err != nil {
		return nil, err
	}
	a.assemblyMembers[session] = members
	return members, nil
}

type BillsResponse struct {
	Bills []BillReference
}

func (a API) GetBillUpdates(ctx context.Context, from, to time.Time) (BillsResponse, error) {
	var out BillsResponse
	resp, err := a.api.GetBillUpdates(ctx, from, to)
	if err != nil {
		return out, err
	}
	for _, bill := range resp.Result.Items {
		out.Bills = append(out.Bills, BillReference{
			PrintNo: bill.ID.BasePrintNo,
			Session: bill.ID.Session,
		})
	}
	return out, nil
}

func (a API) Bills(ctx context.Context, session string, offset int) (BillsResponse, error) {
	var out BillsResponse
	resp, err := a.api.Bills(ctx, session, offset)
	if err != nil {
		return out, err
	}
	for _, bill := range resp.Result.Items {
		out.Bills = append(out.Bills, BillReference{
			PrintNo: bill.BasePrintNo,
			Session: bill.Session,
		})
	}
	return out, nil
}
