package nysenateapi

import (
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/civil"
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
	Date      civil.Date
	Committee string `json:"Committee,omitempty"`
}

type Sponsor struct {
	ID    int
	Name  string
	Short string
}
type Vote struct {
	VoteType  string // COMMITTEE, FLOOR
	Date      civil.Date
	Version   string `json:"Version,omitempty"`
	Chamber   string `json:"Chamber,omitempty"`
	Committee string `json:"Committee,omitempty"`
	Votes     []VoteEntry
}
type VoteEntry struct {
	ID    int
	Name  string `json:"Name,omitempty"`
	Short string `json:"Short,omitempty"`
	Vote  string // Aye, Nay, Excused
}

type Action struct {
	Text    string     `json:"Text,omitempty"`
	Date    civil.Date `json:"Date,omitempty"`
	Chamber string     `json:"Chamber,omitempty"`
	Version string     `json:"Version,omitempty"`
}

var newlineReplacer = strings.NewReplacer("\n", " ")

func newBill(b *verboseapi.Bill) *Bill {
	bill := &Bill{
		PrintNo:    b.BasePrintNo,
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
		ActClause:  newlineReplacer.Replace(b.Amendments.Items[b.ActiveVersion].ActClause),
	}
	for _, m := range b.Milestones.Items {
		bill.Milestones = append(bill.Milestones, Milestone{
			Type:      m.StatusType,
			Date:      civil.DateOf(parseTime(m.ActionDate)),
			Committee: m.CommitteeName,
		})
	}
	for _, m := range b.Actions.Items {
		bill.Actions = append(bill.Actions, Action{
			Text:    m.Text,
			Date:    civil.DateOf(parseTime(m.Date)),
			Chamber: m.Chamber,
			Version: m.BillID.Version,
		})
	}
	// sponsors (including multi-sponsors, etc)
	seen := map[int]bool{}
	if b.Sponsor.Member.MemberID > 0 {
		seen[b.Sponsor.Member.MemberID] = true
		bill.Sponsors = append(bill.Sponsors, Sponsor{
			ID:    b.Sponsor.Member.MemberID,
			Name:  b.Sponsor.Member.FullName,
			Short: b.Sponsor.Member.ShortName,
		})
	}
	for _, s := range b.Amendments.Items[b.ActiveVersion].MultiSponsors.Items {
		if seen[s.MemberID] {
			continue
		}
		seen[s.MemberID] = true
		bill.Sponsors = append(bill.Sponsors, Sponsor{
			ID:    s.MemberID,
			Name:  s.FullName,
			Short: s.ShortName,
		})
	}
	for _, s := range b.Amendments.Items[b.ActiveVersion].CoSponsors.Items {
		if seen[s.MemberID] {
			continue
		}
		seen[s.MemberID] = true
		bill.Sponsors = append(bill.Sponsors, Sponsor{
			ID:    s.MemberID,
			Name:  s.FullName,
			Short: s.ShortName,
		})
	}
	if b.Amendments.Items[b.ActiveVersion].SameAs.Size > 0 {
		bill.SameAsPrintNo = b.Amendments.Items[b.ActiveVersion].SameAs.Items[0].BasePrintNoStr
	}
	bill.Votes = newVotes(b.Votes.Items)
	// previousVersions
	sort.Slice(b.PreviousVersions.Items, func(i, j int) bool { return b.PreviousVersions.Items[i].Session < b.PreviousVersions.Items[j].Session })
	for _, v := range b.PreviousVersions.Items {
		bill.PreviousVersions = append(bill.PreviousVersions, v.BasePrintNoStr) // i.e. S1234-2020
	}

	return bill
}

func newVotes(bv []verboseapi.BillVote) []Vote {
	var o []Vote
	for _, v := range bv {
		o = append(o, Vote{
			VoteType:  v.VoteType,
			Date:      civil.DateOf(parseTime(v.VoteDate)),
			Version:   v.Version,
			Chamber:   v.Committee.Chamber,
			Committee: v.Committee.Name,
			Votes:     newVoteEntries(v.MemberVotes.Items),
		})
	}
	return o
}

func newVoteEntries(v verboseapi.MemberVotes) []VoteEntry {
	var o []VoteEntry
	// TODO: dedupe
	for _, m := range v.Aye.Items {
		o = append(o, VoteEntry{
			ID:    m.MemberID,
			Short: m.ShortName,
			Name:  m.FullName,
			Vote:  "Aye",
		})
	}
	for _, m := range v.AyeWithReservations.Items {
		o = append(o, VoteEntry{
			ID:    m.MemberID,
			Short: m.ShortName,
			Name:  m.FullName,
			Vote:  "Aye",
			// TODO: add note "with reservations"
		})
	}
	for _, m := range v.Nay.Items {
		o = append(o, VoteEntry{
			ID:    m.MemberID,
			Short: m.ShortName,
			Name:  m.FullName,
			Vote:  "Nay",
		})
	}
	for _, m := range v.Excused.Items {
		o = append(o, VoteEntry{
			ID:    m.MemberID,
			Short: m.ShortName,
			Name:  m.FullName,
			Vote:  "Excused",
		})
	}
	for _, m := range v.Absent.Items {
		o = append(o, VoteEntry{
			ID:    m.MemberID,
			Short: m.ShortName,
			Name:  m.FullName,
			Vote:  "Absent",
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
