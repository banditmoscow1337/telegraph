package telegraph

import "github.com/goccy/go-json"

type Query struct {
	Name  string
	Value string
}

type Page struct {
	Path        string        `json:"path"`
	Url         string        `json:"url"`
	Title       string        `json:"title"`
	Description string        `json:"decription"`
	AuthorName  string        `json:"author_name"`
	AuthorUrl   string        `json:"author_url"`
	ImageUrl    string        `json:"image_url"`
	Content     []NodeElement `json:"content"`
	Views       int           `json:"views"`
	CanEdit     bool          `json:"can_edit"`
}

type PageList struct {
	TotalCount int    `json:"total_count"`
	Pages      []Page `json:"pages"`
}

type PageViews struct {
	Views int `json:"views"`
}

type cPage struct {
	AccessToken  string        `json:"access_token"`
	Title        string        `json:"title"`
	Description  string        `json:"decription"`
	AuthorName   string        `json:"author_name"`
	AuthorUrl    string        `json:"author_url"`
	Content      []NodeElement `json:"content"`
	ReturnConent bool          `json:"return_content"`
}

type Account struct {
	ShortName   string `json:"short_name"`
	AuthorName  string `json:"author_name"`
	AuthorUrl   string `json:"author_url"`
	AccessToken string `json:"access_token"`
	AuthUrl     string `json:"auth_url"`
	PageCount   int    `json:"page_count"`
}

type NodeElement struct {
	Tag      string            `json:"tag"`
	Attrs    map[string]string `json:"attrs,omitempty"`
	Children []any             `json:"children"`
}

type Telegraph struct {
	a    Account
	Name string
}

type Result struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}
