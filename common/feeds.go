package common

// FeedConfigurationRequest request
type FeedConfigurationRequest struct {
	Data struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	} `json:"data"`
}

// FeedConfiguration is type for feed
type FeedConfiguration struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
