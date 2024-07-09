package nysenateapi

import (
	"context"
	"time"

	"github.com/jehiah/nysenateapi/verboseapi"
)

type API struct {
	api *verboseapi.NYSenateAPI
}

func NewAPI(token string) API {
	return API{
		api: verboseapi.NewAPI(token),
	}
}

func (a API) GetBill(ctx context.Context, session, printNo string) (*Bill, error) {
	bill, err := a.api.GetBill(ctx, session, printNo)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, nil
	}
	return newBill(bill), nil
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
			PrintNo: bill.ID.PrintNo,
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
			PrintNo: bill.PrintNo,
			Session: bill.Session,
		})
	}
	return out, nil
}
