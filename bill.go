package nysenateapi

import (
	"context"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func (a NYSenateAPI) Bills(ctx context.Context, session string, offset int) (*BillsResponse, error) {
	if session == "" {
		return nil, nil
	}
	params := &url.Values{"offset": []string{fmt.Sprintf("%d", offset)}, "limit": []string{"1000"}}
	path := fmt.Sprintf("/api/3/bills/%s", url.PathEscape(session))
	var data BillsResponse
	log.WithContext(ctx).WithField("session", session).Infof("looking for bills %s", session)
	err := a.get(ctx, path, params, &data)
	return &data, err
}

type BillReference struct {
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
}

type BillsResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ResponseType string `json:"responseType"`
	Total        int    `json:"total"`
	OffsetStart  int    `json:"offsetStart"`
	OffsetEnd    int    `json:"offsetEnd"`
	Limit        int    `json:"limit"`
	Result       struct {
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
	log.WithContext(ctx).WithField("session", session).WithField("printNo", printNo).Infof("looking up bill %s-%s", session, printNo)
	err := a.get(ctx, path, params, &data)
	return &(data.Bill), err
}

type BillResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ResponseType string `json:"responseType"`
	Bill         Bill   `json:"result"`
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
	SubstitutedBy     struct {
		BasePrintNo    string `json:"basePrintNo"`
		Session        int    `json:"session"`
		BasePrintNoStr string `json:"basePrintNoStr"`
	} `json:"substitutedBy"`
	Sponsor struct {
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
			BasePrintNo    string `json:"basePrintNo"`
			Session        int    `json:"session"`
			BasePrintNoStr string `json:"basePrintNoStr"`
			PrintNo        string `json:"printNo"`
			Version        string `json:"version"`
			PublishDate    string `json:"publishDate"`
			SameAs         struct {
				Items []struct {
					BasePrintNo string `json:"basePrintNo"`
					Session     int    `json:"session"`
					PrintNo     string `json:"printNo"`
					Version     string `json:"version"`
				} `json:"items"`
				Size int `json:"size"`
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
			BillID struct {
				BasePrintNo string `json:"basePrintNo"`
				Session     int    `json:"session"`
				PrintNo     string `json:"printNo"`
				Version     string `json:"version"`
			} `json:"billId"`
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
		BillID struct {
			BasePrintNo string `json:"basePrintNo"`
			Session     int    `json:"session"`
			PrintNo     string `json:"printNo"`
			Version     string `json:"version"`
		} `json:"billId"`
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
			BillID struct {
				BasePrintNo string `json:"basePrintNo"`
				Session     int    `json:"session"`
				PrintNo     string `json:"printNo"`
				Version     string `json:"version"`
			} `json:"billId"`
			Date       string `json:"date"`
			Chamber    string `json:"chamber"`
			SequenceNo int    `json:"sequenceNo"`
			Text       string `json:"text"`
		} `json:"items"`
		Size int `json:"size"`
	} `json:"actions"`
	PreviousVersions struct {
		Items []struct {
			BasePrintNo string `json:"basePrintNo"`
			Session     int    `json:"session"`
			PrintNo     string `json:"printNo"`
			Version     string `json:"version"`
		} `json:"items"`
		Size int `json:"size"`
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
