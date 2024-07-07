package nysenateapi

import (
	"context"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

func (a NYSenateAPI) Bills(ctx context.Context, session string, offset int) (*BillsResponse, error) {
	if session == "" {
		return nil, nil
	}
	params := &url.Values{"offset": []string{fmt.Sprintf("%d", offset)}, "limit": []string{"1000"}}
	path := fmt.Sprintf("/api/3/bills/%s", url.PathEscape(session))
	var data BillsResponse
	log.WithContext(ctx).WithField("offset", offset).WithField("session", session).Debugf("bills session:%s", session)
	err := a.get(ctx, path, params, &data)
	return &data, err
}

const timeFormat = "2006-01-02T15:04:05"

// GetBillUpdates returns a list of bills that have been updated in the given time range.
// https://legislation.nysenate.gov/static/docs/html/bills.html#detailed-update-digests
func (a NYSenateAPI) GetBillUpdates(ctx context.Context, from, to time.Time) (*BillUpdateResponse, error) {
	// /api/3/bills/updates/{fromDateTime}
	// should be inputted as 2014-12-10T13:30:02.
	// The fromDateTime and toDateTime range is exclusive/inclusive respectively.
	log.WithContext(ctx).WithField("from", from).WithField("to", to).Debugf("bill updates")
	path := fmt.Sprintf("/api/3/bills/updates/%s/%s", from.Format(timeFormat), to.Format(timeFormat))
	params := &url.Values{}
	params.Set("type", "processed")
	params.Set("detail", "false")
	var data BillUpdateResponse
	err := a.get(ctx, path, params, &data)
	return &data, err
}

type BillReference struct {
	BillID
	BillType struct {
		Chamber    string `json:"chamber"`
		Desc       string `json:"desc"`
		Resolution bool   `json:"resolution"`
	} `json:"billType"`
	Title             string `json:"title"`
	ActiveVersion     string `json:"activeVersion"`
	Year              int    `json:"year"`
	PublishedDateTime string `json:"publishedDateTime"`
}

type BillID struct {
	BasePrintNo    string `json:"basePrintNo"`
	Session        int    `json:"session"`
	BasePrintNoStr string `json:"basePrintNoStr,omitempty"`
	PrintNo        string `json:"printNo,omitempty"`
	Version        string `json:"version,omitempty"`
}

type BillUpdate struct {
	ID              BillID `json:"id"`
	ContentType     string `json:"contentType"`     // i.e. "BILL"
	SourceID        string `json:"sourceId"`        // i.e. "2019-02-13-09.01.14.643609_LDSPON_S01826.XML-1-LDSPON",
	SourceDateTime  string `json:"sourceDateTime"`  // i.e. "2019-02-13T09:01:14.643609",
	ProcessDateTime string `json:"processDateTime"` // i.e. "2019-02-13T09:06:09.796845"
}

type BillUpdateResponse struct {
	Envelope
	Result struct {
		Items []BillUpdate `json:"items"`
	}
}

type Envelope struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ResponseType string `json:"responseType"`
	Total        int    `json:"total"`
	OffsetStart  int    `json:"offsetStart"`
	OffsetEnd    int    `json:"offsetEnd"`
	Limit        int    `json:"limit"`
}

type BillsResponse struct {
	Envelope
	Result struct {
		Items []BillReference `json:"items"`
		Size  int             `json:"size"`
	} `json:"result"`
}

func (a NYSenateAPI) GetBill(ctx context.Context, session, printNo string) (*Bill, error) {
	if session == "" || printNo == "" {
		return nil, nil
	}
	params := &url.Values{}
	params.Set("view", "with_refs")
	path := fmt.Sprintf("/api/3/bills/%s/%s", url.PathEscape(session), url.PathEscape(printNo))
	var data BillResponse
	log.WithContext(ctx).WithField("session", session).WithField("printNo", printNo).Debugf("looking up bill %s-%s", session, printNo)
	err := a.get(ctx, path, params, &data)
	return &(data.Bill), err
}

type BillResponse struct {
	Envelope
	Bill Bill `json:"result"`
}

// from https://legislation.nysenate.gov/static/docs/html/bills.html
type Bill struct {
	BasePrintNo string `json:"basePrintNo"`
	Session     int    `json:"session"`
	PrintNo     string `json:"printNo"`
	BillType    struct {
		Chamber    string `json:"chamber"`
		Desc       string `json:"desc"`
		Resolution bool   `json:"resolution"`
	} `json:"billType"`
	Title             string `json:"title"`
	ActiveVersion     string `json:"activeVersion"`
	Year              int    `json:"year"`
	PublishedDateTime string `json:"publishedDateTime"`
	SubstitutedBy     BillID `json:"substitutedBy"`
	Sponsor           struct {
		Member MemberEntry `json:"member"`
		Budget bool        `json:"budget"`
		Rules  bool        `json:"rules"`
	} `json:"sponsor"`
	Summary string `json:"summary"`
	Signed  bool   `json:"signed"`
	Status  struct {
		StatusType    string      `json:"statusType"`
		StatusDesc    string      `json:"statusDesc"`
		ActionDate    string      `json:"actionDate"`
		CommitteeName string      `json:"committeeName"`
		BillCalNo     interface{} `json:"billCalNo"`
	} `json:"status"`
	Milestones struct {
		Items []struct {
			StatusType    string      `json:"statusType"`
			StatusDesc    string      `json:"statusDesc"`
			ActionDate    string      `json:"actionDate"`
			CommitteeName string      `json:"committeeName"`
			BillCalNo     interface{} `json:"billCalNo"`
		} `json:"items"`
		Size int `json:"size"`
	} `json:"milestones"`
	ProgramInfo struct {
		Name       string `json:"name"`
		SequenceNo int    `json:"sequenceNo"`
	} `json:"programInfo"`
	Amendments struct {
		Items map[string]struct {
			BillID
			PublishDate string `json:"publishDate"`
			SameAs      struct {
				Items []BillID `json:"items"`
				Size  int      `json:"size"`
			} `json:"sameAs"`
			Memo             string          `json:"memo"`
			LawSection       string          `json:"lawSection"`
			LawCode          string          `json:"lawCode"`
			ActClause        string          `json:"actClause"`
			FullTextFormats  []string        `json:"fullTextFormats"`
			FullText         string          `json:"fullText"`
			FullTextHTML     interface{}     `json:"fullTextHtml"`
			FullTextTemplate interface{}     `json:"fullTextTemplate"`
			CoSponsors       MemberEntryList `json:"coSponsors"`
			MultiSponsors    MemberEntryList `json:"multiSponsors"`
			UniBill          bool            `json:"uniBill"`
			Stricken         bool            `json:"stricken"`
		} `json:"items"`
		Size int `json:"size"`
	} `json:"amendments"`
	Votes struct {
		Items []BillVote `json:"items"`
		Size  int        `json:"size"`
	} `json:"votes"`
	VetoMessages struct {
		Items []struct {
			BillID     BillID      `json:"billId"`
			Year       int         `json:"year"`
			VetoNumber int         `json:"vetoNumber"`
			MemoText   string      `json:"memoText"`
			VetoType   string      `json:"vetoType"`
			Chapter    int         `json:"chapter"`
			BillPage   int         `json:"billPage"`
			LineStart  int         `json:"lineStart"`
			LineEnd    int         `json:"lineEnd"`
			Signer     string      `json:"signer"`
			SignedDate interface{} `json:"signedDate"`
		} `json:"items"`
		Size int `json:"size"`
	} `json:"vetoMessages"`
	ApprovalMessage struct {
		BillID         BillID `json:"billId"`
		Year           int    `json:"year"`
		ApprovalNumber int    `json:"approvalNumber"`
		Chapter        int    `json:"chapter"`
		Signer         string `json:"signer"`
		Text           string `json:"text"`
	} `json:"approvalMessage"`
	AdditionalSponsors MemberEntryList `json:"additionalSponsors"`
	PastCommittees     struct {
		Items []struct {
			Chamber       string `json:"chamber"`
			Name          string `json:"name"`
			SessionYear   int    `json:"sessionYear"`
			ReferenceDate string `json:"referenceDate"`
		} `json:"items"`
		Size int `json:"size"`
	} `json:"pastCommittees"`
	Actions struct {
		Items []struct {
			BillID     BillID `json:"billId"`
			Date       string `json:"date"`
			Chamber    string `json:"chamber"`
			SequenceNo int    `json:"sequenceNo"`
			Text       string `json:"text"`
		} `json:"items"`
		Size int `json:"size"`
	} `json:"actions"`
	PreviousVersions struct {
		Items []BillID `json:"items"`
		Size  int      `json:"size"`
	} `json:"previousVersions"`
	CommitteeAgendas struct {
		Items []struct {
			AgendaID struct {
				Number int `json:"number"`
				Year   int `json:"year"`
			} `json:"agendaId"`
			CommitteeID struct {
				Chamber string `json:"chamber"`
				Name    string `json:"name"`
			} `json:"committeeId"`
		} `json:"items"`
		Size int `json:"size"`
	} `json:"committeeAgendas"`
	Calendars struct {
		Items []struct {
			Year           int `json:"year"`
			CalendarNumber int `json:"calendarNumber"`
		} `json:"items"`
		Size int `json:"size"`
	} `json:"calendars"`
	BillInfoRefs struct {
		Items interface{} `json:"items"`
		Size  int         `json:"size"`
	} `json:"billInfoRefs"`
}

type BillVote struct {
	Version   string `json:"version"`
	VoteType  string `json:"voteType"`
	VoteDate  string `json:"voteDate"`
	Committee struct {
		Chamber string `json:"chamber"`
		Name    string `json:"name"`
	} `json:"committee"`
	MemberVotes struct {
		Items MemberVotes `json:"items"`
		Size  int         `json:"size"`
	} `json:"memberVotes"`
}

type MemberVotes struct {
	Aye                 MemberEntryList `json:"AYE"`
	AyeWithReservations MemberEntryList `json:"AYEWR"`
	Nay                 MemberEntryList `json:"NAY"` // ?
	Excused             MemberEntryList `json:"EXC"` // excused
	Absent              MemberEntryList `json:"Absent"`
}
