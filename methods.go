package telegraph

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/banditmoscow1337/utils/crypto"
	"github.com/goccy/go-json"
)

func (t *Telegraph) SetImage(data []byte, name string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", name)
	io.Copy(part, bytes.NewReader(data))
	writer.Close()

	r, _ := http.NewRequest("POST", "https://telegra.ph/upload", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}

	bod, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	var res []struct {
		Src string `json:"src"`
	}

	if err = json.Unmarshal(bod, &res); err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", errors.New("file not uploaded")
	}

	return strings.ReplaceAll(res[0].Src, "/file/", ""), nil
}

func (t *Telegraph) GetImage(link string) ([]byte, error) {
	c := &http.Client{}
	req, _ := http.NewRequest("GET", "https://telegra.ph/file/"+link, nil)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func CreateAccount(names ...string) (t *Telegraph, err error) {
	var short, full string

	t = &Telegraph{}

	switch len(names) {
	case 0:
		short = crypto.RandString(16)
		full = "Anonymous"
		break
	case 1:
		full = names[0]
		short = crypto.RandString(16)
		break
	default:
		short = names[0]
		full = names[1]
	}

	data, err := t.call("createAccount", []Query{{"short_name", short}, {"author_name", full}}, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &t.a)

	t.Name = t.a.ShortName

	return
}

func (t *Telegraph) CreatePage(title string, description string, content []NodeElement, author ...string) (np Page, err error) {
	if len(title) == 0 {
		err = errors.New("title is empty")
		return
	}

	if len(content) == 0 {
		err = errors.New("content is empty")
		return
	}

	var p cPage

	p.AccessToken = t.a.AccessToken
	p.Title = title
	p.Content = content

	if len(description) > 0 {
		p.Description = description
	}

	switch len(author) {
	case 0:
		break
	case 1:
		p.AuthorName = author[0]
		break
	default:
		p.AuthorName = author[0]
		p.AuthorUrl = author[1]
	}

	body, err := json.Marshal(p)
	if err != nil {
		return
	}

	data, err := t.call("createPage", nil, body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &np)
	if err != nil {
		return
	}

	return
}

func GetIntoAccount(token string) (t *Telegraph, err error) {
	if len(token) == 0 {
		err = errors.New("token not found")
		return
	}

	t = &Telegraph{}

	t.a.AccessToken = token

	data, err := t.call("getAccountInfo", nil, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &t.a)

	t.Name = t.a.ShortName

	return
}

func (t *Telegraph) GetPage(path string, returnConent bool) (p Page, err error) {
	if len(path) == 0 {
		err = errors.New("path is empty")
		return
	}

	uri := "getPage/" + path
	var q []Query
	if returnConent {
		q = append(q, Query{"return_content", "true"})
	}

	data, err := t.call(uri, q, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &p)
	if err != nil {
		return
	}

	return
}

func (t *Telegraph) GetPageList(limit, offset int) (pl PageList, err error) {
	if limit == 0 {
		limit = 50
	}
	data, err := t.call("getPageList", []Query{{"limit", strconv.Itoa(limit)}, {"offset", strconv.Itoa(offset)}}, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &pl)
	if err != nil {
		return
	}

	return
}

func (t *Telegraph) GetViews(path string, year, month, day, hour int) (v PageViews, err error) {
	if len(path) == 0 {
		err = errors.New("path is empty")
		return
	}

	var q []Query

	if year >= 2000 {
		q = append(q, Query{"year", strconv.Itoa(year)})
	}

	if month > 0 {
		q = append(q, Query{"month", strconv.Itoa(month)})
	}

	if day > 0 {
		q = append(q, Query{"day", strconv.Itoa(hour)})
	}

	if hour > 0 {
		q = append(q, Query{"hour", strconv.Itoa(hour)})
	}

	data, err := t.call("getViews/"+path, q, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return
	}

	return
}

func (t *Telegraph) EditAccountInfo(short, author, url string) (a Account, err error) {
	if len(short) > 0 {
		t.a.ShortName = short
	}

	if len(author) > 0 {
		t.a.AuthorName = author
	}

	if len(url) > 0 {
		t.a.AuthorUrl = url
	}

	body, err := json.Marshal(t.a)
	if err != nil {
		return
	}

	data, err := t.call("editAccountInfo/", nil, body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &a)
	if err != nil {
		return
	}

	t.Name = t.a.ShortName

	return
}

func (t *Telegraph) EditPage(op Page, returnConent bool) (np Page, err error) {
	if len(op.Path) == 0 {
		err = errors.New("path is empty")
		return
	}

	body, err := json.Marshal(cPage{
		AccessToken: t.a.AccessToken,
		Title:       op.Title,
		Description: op.Description,
		AuthorName:  op.AuthorName,
		AuthorUrl:   op.AuthorUrl,
		Content:     op.Content,
	})
	if err != nil {
		return
	}

	data, err := t.call("editPage/"+op.Path, nil, body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &np)
	if err != nil {
		return
	}
	return
}

func (t *Telegraph) GetToken() string {
	return t.a.AccessToken
}

func (t *Telegraph) RevokeAccessToken(short, author, url string) (a Account, err error) {
	data, err := t.call("revokeAccessToken", nil, nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &a)
	if err != nil {
		return
	}

	t.a = a

	return
}
