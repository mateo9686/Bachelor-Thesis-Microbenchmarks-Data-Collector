package model

type License struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	URL     string `json:"url,omitempty"`
	SPDXID  string `json:"spdx_id,omitempty"`
	NodeID  string `json:"node_id"`
	HtmlUrl string `json:"html_url"`
}
