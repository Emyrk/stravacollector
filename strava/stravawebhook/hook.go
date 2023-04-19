package stravawebhook

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const hookURL = "https://www.strava.com/api/v3/push_subscriptions"

func CreateWebhook(ctx context.Context, clientID string, clientSecret string, callbackURL string, verifyToken string) error {
	r, err := http.NewRequest(http.MethodPost, hookURL, bytes.NewBufferString(url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"callback_url":  {callbackURL},
		"verify_token":  {verifyToken},
	}.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	r = r.WithContext(ctx)

	cli := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := cli.Do(r)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	d, _ := io.ReadAll(resp.Body)
	fmt.Println(string(d))

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code not ok: %d", resp.StatusCode)
	}

	return nil
}

func ViewWebhook(ctx context.Context, clientID string, clientSecret string) error {
	r, err := http.NewRequest(http.MethodGet, hookURL, bytes.NewBufferString(url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	r = r.WithContext(ctx)

	cli := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := cli.Do(r)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	d, _ := io.ReadAll(resp.Body)
	fmt.Println(string(d))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code not ok: %d", resp.StatusCode)
	}

	return nil
}

func DeleteWebhook(ctx context.Context, clientID string, clientSecret string, id string) error {
	r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", hookURL, id), bytes.NewBufferString(url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	r = r.WithContext(ctx)

	cli := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := cli.Do(r)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code not ok: %d", resp.StatusCode)
	}
	d, _ := io.ReadAll(resp.Body)
	fmt.Println(string(d))
	return nil
}
