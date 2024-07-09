package nysenateapi

import (
	"time"

	"github.com/jehiah/nysenateapi/verboseapi"
)

type BillReference struct {
	PrintNo string `json:"PrintNo,omitempty"`
	Session int    `json:"Session,omitempty"`
}

type Bill struct {
	PrintNo    string `json:"PrintNo"`
	Version    string `json:"Version,omitempty"`
	Session    int    `json:"Session"`
	Chamber    string `json:"Chamber,omitempty"`
	BillType   string `json:"BillType,omitempty"`
	Resolution bool   `json:"Resolution,omitempty"`

	Published  time.Time   `json:"Published,omitempty"`
	Status     string      `json:"Status,omitempty"`
	Committee  string      `json:"Committee,omitempty"`
	Milestones []Milestone `json:"Milestones,omitempty"`
	Actions    []Action    `json:"Actions,omitempty"`
	Votes      []Vote      `json:"Votes,omitempty"`
	Sponsors   []Sponsor   `json:"Sponsors,omitempty"`

	Title      string `json:"Title,omitempty"`
	Summary    string `json:"Summary,omitempty"`
	LawSection string `json:"LawSection,omitempty"`
	LawCode    string `json:"LawCode,omitempty"`
	ActClause  string `json:"ActClause,omitempty"`
	// TODO: BodyURL

	SameAsPrintNo    string   `json:"SameAsPrintNo,omitempty"`
	PreviousVersions []string `json:"PreviousVersions,omitempty"`
}

type Milestone struct {
	Type      string
	Date      time.Time
	Committee string `json:"Committee,omitempty"`
}

type Sponsor struct {
	ID        int
	FullName  string
	ShortName string
}
type Vote struct {
	VoteType  string // COMMITTEE, FLOOR
	Date      time.Time
	Version   string
	Chamber   string `json:"Chamber,omitempty"`
	Committee string
	Votes     []VoteEntry
}
type VoteEntry struct {
	ID        int
	FullName  string `json:"FullName,omitempty"`
	ShortName string `json:"ShortName,omitempty"`
	Vote      string // Aye, Nay, Excused
}

type Action struct {
	Text    string    `json:"Text,omitempty"`
	Date    time.Time `json:"Date,omitempty"`
	Chamber string    `json:"Chamber,omitempty"`
	Version string    `json:"Version,omitempty"`
}

func newBill(b *verboseapi.Bill) *Bill {
	bill := &Bill{
		PrintNo:    b.PrintNo,
		Version:    b.ActiveVersion,
		Session:    b.Session,
		Chamber:    b.BillType.Chamber,
		BillType:   b.BillType.Desc,
		Resolution: b.BillType.Resolution,
		Published:  parseTime(b.PublishedDateTime),
		Status:     b.Status.StatusType,
		Committee:  b.Status.CommitteeName,
		Title:      b.Title,
		Summary:    b.Summary,
		LawSection: b.Amendments.Items[b.ActiveVersion].LawSection,
		LawCode:    b.Amendments.Items[b.ActiveVersion].LawCode,
		ActClause:  b.Amendments.Items[b.ActiveVersion].ActClause,
	}
	for _, m := range b.Milestones.Items {
		bill.Milestones = append(bill.Milestones, Milestone{
			Type:      m.StatusType,
			Date:      parseTime(m.ActionDate),
			Committee: m.CommitteeName,
		})
	}
	for _, m := range b.Actions.Items {
		bill.Actions = append(bill.Actions, Action{
			Text:    m.Text,
			Date:    parseTime(m.Date),
			Chamber: m.Chamber,
			Version: m.BillID.Version,
		})
	}
	// sponsors (including multi-sponsors, etc)
	seen := map[int]bool{b.Sponsor.Member.MemberID: true}
	bill.Sponsors = append(bill.Sponsors, Sponsor{
		ID:        b.Sponsor.Member.MemberID,
		FullName:  b.Sponsor.Member.FullName,
		ShortName: b.Sponsor.Member.ShortName,
	})
	for _, s := range b.Amendments.Items[b.ActiveVersion].MultiSponsors.Items {
		if seen[s.MemberID] {
			continue
		}
		seen[s.MemberID] = true
		bill.Sponsors = append(bill.Sponsors, Sponsor{
			ID:        s.MemberID,
			FullName:  s.FullName,
			ShortName: s.ShortName,
		})
	}
	for _, s := range b.Amendments.Items[b.ActiveVersion].CoSponsors.Items {
		if seen[s.MemberID] {
			continue
		}
		seen[s.MemberID] = true
		bill.Sponsors = append(bill.Sponsors, Sponsor{
			ID:        s.MemberID,
			FullName:  s.FullName,
			ShortName: s.ShortName,
		})
	}
	if b.Amendments.Items[b.ActiveVersion].SameAs.Size > 0 {
		bill.SameAsPrintNo = b.Amendments.Items[b.ActiveVersion].SameAs.Items[0].BasePrintNoStr
	}
	// votes
	for _, v := range b.Votes.Items {
		bill.Votes = append(bill.Votes, Vote{
			VoteType:  v.VoteType,
			Date:      parseTime(v.VoteDate),
			Version:   v.Version,
			Chamber:   v.Committee.Chamber,
			Committee: v.Committee.Name,
			Votes:     newVoteEntries(v.MemberVotes.Items),
		})
	}

	return bill
}

func newVoteEntries(v verboseapi.MemberVotes) []VoteEntry {
	var o []VoteEntry
	// TODO: dedupe
	for _, m := range v.Excused.Items {
		o = append(o, VoteEntry{
			ShortName: m.ShortName,
			FullName:  m.FullName,
			ID:        m.MemberID,
			Vote:      "Excused",
		})
	}
	for _, m := range v.Aye.Items {
		o = append(o, VoteEntry{
			ID:        m.MemberID,
			ShortName: m.ShortName,
			FullName:  m.FullName,
			Vote:      "Aye",
		})
	}
	for _, m := range v.Nay.Items {
		o = append(o, VoteEntry{
			ID:        m.MemberID,
			ShortName: m.ShortName,
			FullName:  m.FullName,
			Vote:      "Nay",
		})
	}
	for _, m := range v.AyeWithReservations.Items {
		o = append(o, VoteEntry{
			ID:        m.MemberID,
			ShortName: m.ShortName,
			FullName:  m.FullName,
			Vote:      "Aye",
			// TODO: add note "with reservations"
		})
	}
	for _, m := range v.Absent.Items {
		o = append(o, VoteEntry{
			ID:        m.MemberID,
			ShortName: m.ShortName,
			FullName:  m.FullName,
			Vote:      "Absent",
		})
	}
	// TODO: Abstained ?
	return o
}

func parseTime(s string) time.Time {
	for _, pattern := range []string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.999999",
		"2006-01-02T15:04",
		"2006-01-02",
	} {
		t, err := time.Parse(pattern, s)
		if err == nil {
			return t
		}
	}
	return time.Unix(0, 0)
}
