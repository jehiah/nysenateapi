package nysenateapi

import (
	"context"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
)

// Note: response might have duplicates
func (a NYSenateAPI) GetMembers(ctx context.Context, session string, c Chamber) ([]MemberEntry, error) {
	if session == "" || c == "" {
		return nil, nil
	}
	log.WithContext(ctx).WithField("session", session).WithField("chamber", c).Debugf("GetMembers session:%s", session)
	path := fmt.Sprintf("/api/3/members/%s/%s", url.PathEscape(session), url.PathEscape(string(c)))
	// senate is 63, assembly is 150
	params := &url.Values{"full": []string{"true"}, "limit": []string{"200"}}
	var data MemberListResponse
	err := a.get(ctx, path, params, &data)
	if err != nil {
		return nil, err
	}
	var out []MemberEntry
	for _, m := range data.Result.Items {
		out = append(out, m.MemberEntry())
		for memberSession, mm := range m.Sessions {
			for _, mmm := range mm {
				// could have a different short name for this session
				// i.e. 2021: BICHOTTE, BICHOTTE HERMELYN
				if memberSession == session && mmm.ShortName != m.ShortName {
					out = append(out,
						MemberEntry{
							MemberID:     m.MemberID,
							FullName:     m.FullName,
							ShortName:    mmm.ShortName,
							Chamber:      mmm.Chamber,
							DistrictCode: mmm.DistrictCode,
							Alternate:    mmm.Alternate,
						})
				}
			}
		}
	}
	return out, nil
}

// https://legislation.nysenate.gov/static/docs/html/members.html
type MemberSessionResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	ResponseType string        `json:"responseType"` // "member-sessions"
	Result       MemberSession `json:"result"`
}

type MemberListResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ResponseType string `json:"responseType"` // "member-sessions list"
	Result       struct {
		Items []MemberSession `json:"items"`
	} `json:"result"`
}

type MemberSession struct {
	MemberID     int                      `json:"memberId"` // memberId
	Chamber      string                   `json:"chamber"`  // SENATE
	Incumbent    bool                     `json:"incumbent"`
	FullName     string                   `json:"fullName"` // "James L. Seward"
	ShortName    string                   `json:"shortName"`
	DistrictCode int                      `json:"districtCode"`
	Sessions     map[string][]MemberEntry `json:"sessionShortNameMap"` // year: [...]
	Person       Person                   `json:"person"`
}

func (m MemberSession) MemberEntry() MemberEntry {
	return MemberEntry{
		MemberID:     m.MemberID,
		FullName:     m.FullName,
		ShortName:    m.ShortName,
		Chamber:      m.Chamber,
		DistrictCode: m.DistrictCode,
	}
}

type MemberEntry struct {
	MemberID        int    `json:"memberId"`
	FullName        string `json:"fullName,omitempty"`
	ShortName       string `json:"shortName"`
	Chamber         string `json:"chamber"` // SENATE
	DistrictCode    int    `json:"districtCode"`
	Alternate       bool   `json:"alternate"`
	SessionYear     int    `json:"sessionYear"`
	SessionMemberID int    `json:"sessionMemberId,omitempty"`
}

type MemberEntryList struct {
	Items []MemberEntry `json:"items"`
	Size  int           `json:"size"`
}

type Person struct {
	PersonID   int         `json:"personId"`
	FullName   string      `json:"fullName"`
	FirstName  string      `json:"firstName"`
	MiddleName string      `json:"middleName"`
	LastName   string      `json:"lastName"`
	Email      string      `json:"email"`
	Prefix     string      `json:"prefix"`
	Suffix     interface{} `json:"suffix"`
	Verified   bool        `json:"verified"`
	ImgName    string      `json:"imgName"`
}
