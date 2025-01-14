package verboseapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const apiDomain = "https://legislation.nysenate.gov"

func NewAPI(token string) *NYSenateAPI {
	if token == "" {
		panic("missing token")
	}
	return &NYSenateAPI{
		token:     token,
		UserAgent: "https://github.com/jehiah/nysenateapi",
		Limiter:   rate.NewLimiter(rate.Every(5*time.Millisecond), 25),
	}
}

type NYSenateAPI struct {
	token     string
	UserAgent string

	// Limiter throttles requests to the API
	Limiter *rate.Limiter
}

type Chamber string

const SenateChamber Chamber = "senate"
const AssemblyChamber Chamber = "assembly"

func (a NYSenateAPI) get(ctx context.Context, path string, params *url.Values, v interface{}) error {
	err := a.Limiter.Wait(ctx)
	if err != nil {
		return err
	}
	if params == nil {
		params = &url.Values{}
	}
	params.Set("key", a.token)
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
