package core

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CreateAPIRequest comment
func CreateAPIRequest(cfg *Config, organizationID string, urlString string, data *url.Values) (err error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	apiRequest := &SystemAPI{
		cfg:            cfg,
		OrganizationID: organizationID,
		URL:            u,
		PostData:       data,
	}

	ch := make(chan error)
	go apiRequest.makeRequest(ch)
	if err := <-ch; err != nil {
		return err
	}

	return err
}

// makeRequest comment
func (ctx *SystemAPI) makeRequest(ch chan error) {
	if ctx.PostData == nil {
		ctx.PostData = &url.Values{}
	}

	urlString, err := NormalizeRawURLString(ctx.URL.String())
	if err != nil {
		ch <- err
		return
	}

	req, err := http.NewRequest("POST", urlString, strings.NewReader(ctx.PostData.Encode()))
	if err != nil {
		ch <- err
		return
	}

	req.Header.Set("x-organization-id", ctx.OrganizationID)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: time.Second * 15,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		ch <- err
		return
	}
	defer resp.Body.Close()

	ch <- err
}
