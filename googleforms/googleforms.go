package googleforms

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var fieldFinder = regexp.MustCompile(`name="(entry\.[^\"]+)"`)
var actionFinder = regexp.MustCompile(`form action="([^\"]+)`)

type Form struct {
	URL       string
	ActionURL string
	FieldKeys []string
}

func (f *Form) Init() error {
	resp, err := http.Get(f.URL)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	dataString := string(data)
	action := actionFinder.FindStringSubmatch(dataString)
	if action == nil {
		return errors.New("no form action found")
	}
	f.ActionURL = action[1]
	matches := fieldFinder.FindAllStringSubmatch(dataString, -1)
	if len(matches) == 0 {
		return errors.New("no fields found")
	}
	for _, match := range matches {
		f.FieldKeys = append(f.FieldKeys, match[1])
	}
	return nil
}

func (f *Form) Post(formvals ...string) error {
	form := url.Values{}
	for i, val := range formvals {
		form.Add(f.FieldKeys[i], val)
	}
	req, err := http.NewRequest("POST", f.ActionURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.86 Safari/537.36")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("non-200 status")
	}
	return nil
}

func NewForm(formurl string) Form {
	return Form{
		URL: formurl,
	}
}
