package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// client talks to LinkedIn's versioned REST API. base and version are overridable
// (LINKEDIN_API_BASE, LINKEDIN_VERSION) so tests can point it at a fake server.
type client struct {
	token   string
	base    string
	version string
	http    *http.Client
}

func newClient(token string) *client {
	base := os.Getenv("LINKEDIN_API_BASE")
	if base == "" {
		base = "https://api.linkedin.com"
	}
	version := os.Getenv("LINKEDIN_VERSION")
	if version == "" {
		version = "202604"
	}
	return &client{token: token, base: strings.TrimRight(base, "/"), version: version, http: &http.Client{Timeout: 60 * time.Second}}
}

func (c *client) do(method, path string, body any, extra map[string]string) (*http.Response, []byte, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		r = bytes.NewReader(b)
	}
	url := path
	if strings.HasPrefix(path, "/") {
		url = c.base + path
	}
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Linkedin-Version", c.version)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range extra {
		req.Header.Set(k, v)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusUnauthorized {
		return resp, data, fmt.Errorf("LinkedIn rejected the token (401). It lasts 60 days: run `linkedin auth` again and update the LINKEDIN_ACCESS_TOKEN secret")
	}
	if resp.StatusCode >= 300 {
		return resp, data, fmt.Errorf("%s %s: %s: %s", method, path, resp.Status, strings.TrimSpace(string(data)))
	}
	return resp, data, nil
}

// whoami returns the signed-in member's URN; it doubles as a token check.
func (c *client) whoami() (string, error) {
	_, data, err := c.do(http.MethodGet, "/v2/userinfo", nil, nil)
	if err != nil {
		return "", err
	}
	var u struct {
		Sub string `json:"sub"`
	}
	if err := json.Unmarshal(data, &u); err != nil || u.Sub == "" {
		return "", fmt.Errorf("unexpected userinfo response: %s", data)
	}
	return "urn:li:person:" + u.Sub, nil
}

// uploadImage registers an image owned by the member, uploads the bytes, and returns its URN.
func (c *client) uploadImage(owner string, png []byte) (string, error) {
	_, data, err := c.do(http.MethodPost, "/rest/images?action=initializeUpload",
		map[string]any{"initializeUploadRequest": map[string]any{"owner": owner}}, nil)
	if err != nil {
		return "", err
	}
	var init struct {
		Value struct {
			UploadURL string `json:"uploadUrl"`
			Image     string `json:"image"`
		} `json:"value"`
	}
	if err := json.Unmarshal(data, &init); err != nil || init.Value.UploadURL == "" {
		return "", fmt.Errorf("unexpected initializeUpload response: %s", data)
	}
	req, err := http.NewRequest(http.MethodPut, init.Value.UploadURL, bytes.NewReader(png))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "image/png")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("image upload: %s: %s", resp.Status, b)
	}
	return init.Value.Image, nil
}

type article struct {
	URL, Title, Description, Thumbnail string
}

// createPost publishes a public article post and returns its URN (urn:li:share:… or urn:li:ugcPost:…).
func (c *client) createPost(author, commentary string, a article) (string, error) {
	body := map[string]any{
		"author":     author,
		"commentary": escapeCommentary(commentary),
		"visibility": "PUBLIC",
		"distribution": map[string]any{
			"feedDistribution":               "MAIN_FEED",
			"targetEntities":                 []any{},
			"thirdPartyDistributionChannels": []any{},
		},
		"content": map[string]any{
			"article": map[string]any{
				"source":      a.URL,
				"title":       a.Title,
				"description": a.Description,
				"thumbnail":   a.Thumbnail,
			},
		},
		"lifecycleState":            "PUBLISHED",
		"isReshareDisabledByAuthor": false,
	}
	resp, data, err := c.do(http.MethodPost, "/rest/posts", body, nil)
	if err != nil {
		return "", err
	}
	id := resp.Header.Get("X-Restli-Id")
	if id == "" {
		return "", fmt.Errorf("post created but LinkedIn returned no x-restli-id header: %s", data)
	}
	return id, nil
}

// escapeCommentary backslash-escapes the characters LinkedIn's "little text"
// format reserves, so text like "(Aug–Sep)" or "#1" posts literally.
func escapeCommentary(s string) string {
	var b strings.Builder
	for _, r := range s {
		if strings.ContainsRune(`\|{}@[]()<>#*_~`, r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}
