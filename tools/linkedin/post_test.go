package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWithLinkedIn(t *testing.T) {
	url := "https://www.linkedin.com/feed/update/urn:li:share:1/"
	cases := map[string]struct{ in, want string }{
		"adds params table": {
			"+++\ntitle = \"A\"\ntags = [\"x\"]\n+++\n\nBody\n",
			"+++\ntitle = \"A\"\ntags = [\"x\"]\n[params]\n  linkedin = \"" + url + "\"\n+++\n\nBody\n",
		},
		"joins existing params": {
			"+++\ntitle = \"A\"\n[params]\n  other = 1\n+++\nBody",
			"+++\ntitle = \"A\"\n[params]\n  linkedin = \"" + url + "\"\n  other = 1\n+++\nBody",
		},
	}
	for name, c := range cases {
		if got := withLinkedIn(c.in, url); got != c.want {
			t.Errorf("%s:\n got %q\nwant %q", name, got, c.want)
		}
		if frontMatterLinkedIn(withLinkedIn(c.in, url)) != url {
			t.Errorf("%s: URL not read back", name)
		}
	}
}

func TestEscapeCommentary(t *testing.T) {
	if got := escapeCommentary("3x PRs (Aug–Sep) #ai @me"); got != `3x PRs \(Aug–Sep\) \#ai \@me` {
		t.Errorf("got %q", got)
	}
}

// TestPostFlow runs `post` end to end against a fake LinkedIn API.
func TestPostFlow(t *testing.T) {
	site := t.TempDir()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.MkdirAll(filepath.Join(site, "content", "story"), 0o755))
	must(os.MkdirAll(filepath.Join(site, "public", "og"), 0o755))
	md := filepath.Join(site, "content", "story", "new-one.md")
	must(os.WriteFile(md, []byte("+++\ntitle = \"New one\"\n+++\n\nHello\n"), 0o644))
	must(os.WriteFile(filepath.Join(site, "public", "og", "story-new-one.png"), []byte("PNG"), 0o644))
	index := `[{"title":"New one","url":"/story/new-one/","kind":"Story","summary":"Why (really) #1"}]`
	must(os.WriteFile(filepath.Join(site, "public", "search-index.json"), []byte(index), 0o644))

	var uploaded string
	var posted map[string]any
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == "GET" && r.URL.Path == "/v2/userinfo":
			io.WriteString(w, `{"sub":"abc123"}`)
		case r.Method == "POST" && r.URL.Path == "/rest/images":
			if r.Header.Get("Linkedin-Version") == "" || r.Header.Get("X-Restli-Protocol-Version") != "2.0.0" {
				t.Error("missing versioned REST headers")
			}
			io.WriteString(w, `{"value":{"uploadUrl":"`+srv.URL+`/upload","image":"urn:li:image:IMG"}}`)
		case r.Method == "PUT" && r.URL.Path == "/upload":
			b, _ := io.ReadAll(r.Body)
			uploaded = string(b)
		case r.Method == "POST" && r.URL.Path == "/rest/posts":
			json.NewDecoder(r.Body).Decode(&posted)
			w.Header().Set("X-Restli-Id", "urn:li:share:999")
			w.WriteHeader(http.StatusCreated)
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	t.Setenv("LINKEDIN_API_BASE", srv.URL)
	t.Setenv("LINKEDIN_ACCESS_TOKEN", "tok")
	t.Setenv("LINKEDIN_PERSON_URN", "urn:li:person:abc123")

	must(runPost([]string{"-slug", "new-one", "-site", site}))

	if uploaded != "PNG" {
		t.Errorf("card not uploaded, got %q", uploaded)
	}
	if posted["author"] != "urn:li:person:abc123" || posted["commentary"] != `Why \(really\) \#1` {
		t.Errorf("bad post body: %v", posted)
	}
	art := posted["content"].(map[string]any)["article"].(map[string]any)
	if art["source"] != siteURL+"/story/new-one/" || art["thumbnail"] != "urn:li:image:IMG" || art["title"] != "New one" {
		t.Errorf("bad article: %v", art)
	}
	b, _ := os.ReadFile(md)
	if !strings.Contains(string(b), `linkedin = "https://www.linkedin.com/feed/update/urn:li:share:999/"`) {
		t.Errorf("front matter not updated:\n%s", b)
	}
	// A second run must refuse rather than double-post.
	if err := runPost([]string{"-slug", "new-one", "-site", site}); err == nil || !strings.Contains(err.Error(), "already on LinkedIn") {
		t.Errorf("expected double-post refusal, got %v", err)
	}
}
