package model

type GitRepo struct {
	FullName        string   `json:"full_name"`
	Name            string   `json:"name"`
	IsPrivate       bool     `json:"private"`
	Url             string   `json:"html_url"`
	Description     string   `json:"description,omitempty"`
	IsFork          bool     `json:"fork"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	PushedAt        string   `json:"pushed_at"`
	SshUrl          string   `json:"ssh_url"`
	Homepage        string   `json:"homepage,omitempty"`
	StarsCount      int      `json:"stargazers_count"`
	WatchersCount   int      `json:"watchers_count"`
	HasIssues       bool     `json:"has_issues"`
	HasProjects     bool     `json:"has_projects"`
	HasDownloads    bool     `json:"has_downloads"`
	HasWiki         bool     `json:"has_wiki"`
	HasPages        bool     `json:"has_pages"`
	HasDiscussions  bool     `json:"has_discussions"`
	ForksCount      int      `json:"forks"`
	Language        string   `json:"language,omitempty"`
	IsArchived      bool     `json:"archived"`
	IsDisabled      bool     `json:"disabled"`
	OpenIssuesCount int      `json:"open_issues"`
	Topics          []string `json:"topics"`
	AllowForking    bool     `json:"allow_forking"`
	IsTemplate      bool     `json:"is_template"`
	CloningError    string   `json:"cloning_error"`
	Score           float32  `json:"score"`
	License         License  `json:"license"`
}
