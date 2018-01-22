package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Result struct {
	Id        int      `json:"id"`
	Title     string   `json"title"`
	URI       string   `json:"uri"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
	isStarred bool     `json:"isStarred"`
}

func (r Result) Name() string {
	return strings.Replace(strings.ToLower(r.Title), " ", "_", -1)
}

type Dashboard struct {
	Meta struct {
		Slug string `json: "slug"`
	}
	Dashboard json.RawMessage `json:"dashboard"`
}

type Client struct {
	Endpoint string
	Token    string
	client   *http.Client
}

func NewClient(endpoint, token string) Client {
	return Client{
		Endpoint: endpoint,
		Token:    token,
		client:   &http.Client{},
	}
}

func (c Client) Do(method, path string) (*http.Response, error) {
	target := fmt.Sprintf("%s%s", c.Endpoint, path)
	log.Printf("Requesting %s\n", target)
	req, err := http.NewRequest(method, target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	return c.client.Do(req)
}

func (c Client) Search() ([]*Result, error) {
	resp, err := c.Do("GET", "/api/search")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	results := []*Result{}
	return results, json.Unmarshal(raw, &results)
}

func (c Client) Dashboard(result *Result) (*Dashboard, error) {
	resp, err := c.Do("GET", fmt.Sprintf("/api/dashboards/%s", result.URI))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	dashboard := &Dashboard{}
	return dashboard, json.Unmarshal(raw, dashboard)
}

func maybe(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		url   = flag.String("url", "", "Grafana URL")
		token = flag.String("token", "", "Grafana API Token")
		path  = flag.String("path", "./dashboards", "Path to save dashboards")
	)
	flag.Parse()
	client := NewClient(*url, *token)
	results, err := client.Search()
	maybe(err)
	*path = strings.TrimRight(*path, "/")
	os.Mkdir(*path, 0755)
	for _, result := range results {
		dashboard, err := client.Dashboard(result)
		maybe(err)
		fp := fmt.Sprintf("%s/%s.json", *path, result.Name())
		log.Printf("Writing dashboard to %s\n", fp)
		maybe(ioutil.WriteFile(fp, dashboard.Dashboard, 0644))
	}
}
