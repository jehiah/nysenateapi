package nysenateapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

const apiDomain = "https://legislation.nysenate.gov"

func NewAPI(token string) *NYSenateAPI {
	if token == "" {
		panic("missing token")
	}
	return &NYSenateAPI{
		token:     token,
		UserAgent: "https://github.com/jehiah/nysenateapi",
	}
}

type NYSenateAPI struct {
	token     string
	UserAgent string
}

// var Sessions = Sessions{
// 	{2023, 2024},
// 	{2021, 2022},
// 	{2019, 2020},
// 	{2017, 2018},
// 	{2015, 2016},
// 	{2013, 2014},
// 	{2011, 2012},
// 	{2009, 2010},
// 	{2007, 2008},
// }

type Chamber string

const SenateChamber Chamber = "senate"
const AssemblyChamber Chamber = "assembly"

func (a NYSenateAPI) get(ctx context.Context, path string, params *url.Values, v interface{}) error {
	if params == nil {
		params = &url.Values{}
	}
	params.Set("key", a.token)
	params.Set("view", "with_refs")
	u := apiDomain + path
	log.WithContext(ctx).WithField("nysenate_api", u+"?"+params.Encode()).Debug("NYSenateAPI.get")
	req, err := http.NewRequestWithContext(ctx, "GET", u+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", a.UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(&v)
}
