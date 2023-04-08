package main

func GetOwnID() string {
	var resp *struct {
		ID string `json:"_id"`
	}

	SEFetch(`GET`, `/channels/me`, nil, &resp, 1)
	if resp == nil {
		panic("Couldn't identify self!")
	}

	return resp.ID
}

type SETip struct {
	User     *SEUser `json:"user"`
	Provider string  `json:"provider"`
	Message  string  `json:"message"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Imported bool    `json:"imported"`
}

type SEUser struct {
	Name string `json:"username"`
}