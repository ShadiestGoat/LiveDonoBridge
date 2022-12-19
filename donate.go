package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type DonEventType int

const (
	DET_NONE DonEventType = iota
	DET_NEW_DON
	DET_NEW_FUND
	DET_PING
)

type Donation struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"ppOrderID"`
	CaptureID string  `json:"ppCaptureID"`
	Donor     string  `json:"donor"`
	Message   string  `json:"message"`
	Amount    float64 `json:"amount"`
	FundID    string  `json:"fundID"`
}

type DonEvent struct {
	Type DonEventType    `json:"event"`
	Body json.RawMessage `json:"body"`
}

func DonationsRestart() {
	time.Sleep(30 * time.Second)
	InitDonation()
}

type DonDonor struct {
	ID        string `json:"id"`
	DiscordID string `json:"discordID"`
	PayPal    string `json:"PayPal"`
	CycleDay  int    `json:"payCycle"`
}

type DonProfileResponse struct {
	Donors []*DonDonor `json:"donors"`
}

func DonorID(userID string, resolve bool) *DonProfileResponse {
	q := ""
	if resolve {
		q = "?resolve=true"
	}
	resp := &DonProfileResponse{}
	DonFetch(http.MethodGet, "/donors/donor/"+userID+q, nil, resp)

	return resp
}

func InitDonation() {
	headers := http.Header{
		"Authorization": []string{DONATION_TOKEN},
	}
	conn, resp, err := websocket.DefaultDialer.Dial(fmt.Sprintf(`wss://%s/api/ws`, DONATION_HOST), headers)

	if err != nil || conn == nil || resp.StatusCode != 101 {
		body := ""
		if resp != nil && resp.Body != nil {
			b, _ := io.ReadAll(resp.Body)
			body = string(b)
		}
		status := 0
		if resp != nil {
			status = resp.StatusCode
		}
		PrintErr("Couldn't connect to the donation api: '%v', '%v', '%v'", err, status, body)
		go DonationsRestart()

		return
	}

	conn.SetPingHandler(func(appData string) error {
		err = conn.WriteControl(websocket.PongMessage, []byte{}, time.Time{})
		return err
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				PrintErr("The Donation WS Connection was closed! Restarting in 20s...")
			} else {
				PrintErr("Unknown Don WS error: %v", err)
			}
			conn.Close()
			conn = nil
			go DonationsRestart()
			return
		}
		evRaw := DonEvent{}

		if json.Unmarshal(message, &evRaw) != nil {
			PrintErr("Couldn't parse WS Donation: '%v'", string(message))
			continue
		}

		switch evRaw.Type {
		case DET_NEW_DON:
			donation := &Donation{}

			err := json.Unmarshal(evRaw.Body, donation)
			if err != nil {
				PrintErr("Bad donation parse! '%v'", string(evRaw.Body))
				continue
			}

			donor := DonorID(donation.Donor, false)

			name := "Someone"

			if len(donor.Donors) >= 1 {
				discordID := donor.Donors[0].DiscordID
				if discordID != "anon" {
					var resp *struct {
						Name string `json:"username"`
					}
					DiscordFetch("GET", `/users/`+discordID, nil, &resp)
					if resp != nil {
						name = resp.Name
					}
				}
			}

			SEFetch(`POST`, `/tips/`+TWITCH_ID, struct {
				User struct {
					Name string `json:"username"`
				} `json:"user"`
				Provider string  `json:"provider"`
				Message  string  `json:"message"`
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Imported bool    `json:"imported"`
			}{
				User: struct {
					Name string `json:"username"`
				}{
					Name: name,
				},
				Provider: "donate.shadygoat.eu",
				Message:  donation.Message,
				Amount:   donation.Amount,
				Currency: "EUR",
				Imported: true,
			}, nil, 10)
		case DET_PING:
			conn.WriteMessage(1, []byte{'P'})
		}
	}
}
