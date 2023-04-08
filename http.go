package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/shadiestgoat/log"
)

func Fetch(method, url string, body, resp any, headers http.Header, attempt int, maxAttempts int) {
	var reqBody io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err == nil {
			reqBody = bytes.NewReader(b)
			headers.Set("Content-Type", "application/json")
		}
	}
	req, _ := http.NewRequest(method, url, reqBody)
	req.Header = headers
	res, err := http.DefaultClient.Do(req)
	if err != nil || res == nil || res.StatusCode != 200 || (resp != nil && res.Body == nil) {
		if attempt > maxAttempts {
			status := 0
			body := ""
			if res != nil {
				status = res.StatusCode
				if res.Body != nil {
					tmp, _ := io.ReadAll(res.Body)
					body = string(tmp)
				}
			}
			log.Fatal("Couldn't %s '%s': err: %v, status: %d, body: '%s'", method, url, err, status, body)
		}
		time.Sleep(7 * time.Second)
		Fetch(method, url, body, resp, headers, attempt + 1, maxAttempts)
		return
	}
	if resp != nil {
		json.NewDecoder(res.Body).Decode(resp)
	}
}

func DonFetch(method, path string, body any, resp any) {
	headers := http.Header{}
	headers.Set("Authorization", DONATION_TOKEN)
	Fetch(method, "https://" + DONATION_HOST + "/api" + path, body, resp, headers, 1, 1)
}

func DiscordFetch(method, path string, body any, resp any) {
	headers := http.Header{}
	headers.Set("Authorization", "Bot " + DISCORD_TOKEN)
	Fetch(method, "https://discord.com/api/v10" + path, body, resp, headers, 1, 1)
}

func SEFetch(method string, path string, body any, respJson any, max int) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer " + STREAM_ELEMENTS_JWT)
	Fetch(method, "https://api.streamelements.com/kappa/v2" + path, body, respJson, headers, 1, 1)
}
