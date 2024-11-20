package main

// struct to match the schedules JSON body resp structure
type SchedulesBody struct {
	Data []struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Name     string   `json:"name"`
			Tags     []string `json:"tags"`
			Timezone string   `json:"timezone"`
		} `json:"attributes"`
		Relationships map[string]any `json:"relationships"`
	} `json:"data"`
	// not interested in the rest of the body
	Included []any `json:"included"`
}

// struct to match the on-call JSON body resp structure
type OnCallBody struct {
	Data struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"attributes"`
		Relationships struct {
			User struct {
				Data struct {
					Id   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"user"`
		} `json:"relationships"`
	} `json:"data"`
	Included []struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		} `json:"attributes"`
	} `json:"included"`
}
