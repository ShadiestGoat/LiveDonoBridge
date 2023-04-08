package main

import (
	"time"

	"github.com/shadiestgoat/donation-api-wrapper"
	"github.com/shadiestgoat/log"
)

func InitDonations() {
	c := donations.NewClient(DONATION_TOKEN, donations.WithCustomLocation(DONATION_HOST))
	log.Debug("Opening a new WS conn...")

	c.AddHandler(func (c *donations.Client, e *donations.EventClose)  {
		if e.Err != nil {
			log.Error("Closed connection with error: %v", e.Err)
		} else {
			log.Error("Closed connection with no error??")
		}

		time.Sleep(30 * time.Second)
		log.FatalIfErr(c.OpenWS(), "opening WS")
	})

	c.AddHandler(func (c *donations.Client, e *donations.EventOpen)  {
		log.Success("Connected!")
	})

	c.AddHandler(func (c *donations.Client, e *donations.EventNewDonation)  {
		donor, err := c.DonorByID(e.Donor, false)
		if log.ErrorIfErr(err, "fetching donor by id") {
			return
		}

		name := "Someone"

		for _, d := range donor.Donors {
			discordID := d.DiscordID

			if discordID != "anon" && discordID != "" {
				var resp *struct {
					Name string `json:"username"`
				}
				DiscordFetch("GET", `/users/`+discordID, nil, &resp)
				if resp != nil {
					name = resp.Name
					break
				}
			}
		}

		SEFetch(`POST`, `/tips/`+TWITCH_ID, &SETip{
			User: &SEUser{
				Name: name,
			},
			Provider: "donate.shadygoat.eu",
			Message:  e.Message,
			Amount:   e.Amount,
			Currency: "EUR",
			Imported: true,
		}, nil, 10)
	})

	log.FatalIfErr(c.OpenWS(), "opening WS")
}
