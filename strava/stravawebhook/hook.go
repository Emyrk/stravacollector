package stravawebhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Emyrk/strava/strava/stravalimit"

	"github.com/Emyrk/strava/strava"
)

const hookURL = "https://www.strava.com/api/v3/push_subscriptions"

func CreateWebhook(ctx context.Context, clientID string, clientSecret string, callbackURL string, verifyToken string) (int, error) {
	r, err := http.NewRequest(http.MethodPost, hookURL, bytes.NewBufferString(url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"callback_url":  {callbackURL},
		"verify_token":  {verifyToken},
	}.Encode()))
	if err != nil {
		return -1, fmt.Errorf("create request: %w", err)
	}
	r = r.WithContext(ctx)

	cli := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := cli.Do(r)
	if err != nil {
		return -1, fmt.Errorf("do request: %w", err)
	}

	var hook Webhook
	err = WebhookResponse(http.StatusCreated, resp, &hook)
	if err != nil {
		return -1, err
	}
	return hook.ID, nil
}

func ViewWebhook(ctx context.Context, clientID string, clientSecret string) ([]Webhook, error) {
	q := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", hookURL, q.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	r = r.WithContext(ctx)

	cli := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := cli.Do(r)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	var hook []Webhook
	return hook, WebhookResponse(http.StatusOK, resp, &hook)
}

func DeleteWebhook(ctx context.Context, clientID string, clientSecret string, id int) error {
	r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", hookURL, id), bytes.NewBufferString(url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	r = r.WithContext(ctx)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println(r.URL.String())

	cli := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := cli.Do(r)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	return WebhookResponse(http.StatusNoContent, resp, nil)
}

func WebhookResponse(expectedCode int, resp *http.Response, into any) error {
	stravalimit.Update(resp.Header)

	if resp.StatusCode != expectedCode {
		standardErr := fmt.Errorf("status code not ok: %d", resp.StatusCode)
		var e strava.Error
		err := json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return standardErr
		}

		if e.Message == "" {
			return standardErr
		}
		return e
	}

	if into == nil {
		d, _ := io.ReadAll(resp.Body)
		fmt.Println(string(d))
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(into)
}

type Webhook struct {
	ID            int       `json:"id"`
	ResourceState int       `json:"resource_state"`
	ApplicationID int       `json:"application_id"`
	CallbackURL   string    `json:"callback_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
