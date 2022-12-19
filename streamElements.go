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
