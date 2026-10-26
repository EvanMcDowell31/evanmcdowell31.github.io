package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const siteURL = "https://evanmcdowell31.github.io"

type indexEntry struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Kind    string `json:"kind"`
	Summary string `json:"summary"`
}

var slugRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// runPost shares one built Story post and records the LinkedIn URL in its front matter.
func runPost(args []string) error {
	fs := flag.NewFlagSet("post", flag.ExitOnError)
	slug := fs.String("slug", "", "story slug, e.g. the-job-is-changing")
	text := fs.String("text", "", "intro text for the LinkedIn post (default: the post's summary)")
	dry := fs.Bool("dry-run", false, "print what would be posted and change nothing")
	site := fs.String("site", ".", "Hugo site root (with content/ and a built public/)")
	fs.Parse(args)

	if !slugRe.MatchString(*slug) {
		return fmt.Errorf("-slug must be a story slug like the-job-is-changing")
	}
	mdPath := filepath.Join(*site, "content", "story", *slug+".md")
	md, err := os.ReadFile(mdPath)
	if err != nil {
		return fmt.Errorf("no story at %s", mdPath)
	}
	if existing := frontMatterLinkedIn(string(md)); existing != "" {
		return fmt.Errorf("%s is already on LinkedIn: %s", *slug, existing)
	}

	entry, err := findEntry(filepath.Join(*site, "public", "search-index.json"), "/story/"+*slug+"/")
	if err != nil {
		return err
	}
	cardPath := filepath.Join(*site, "public", "og", "story-"+*slug+".png")
	card, err := os.ReadFile(cardPath)
	if err != nil {
		return fmt.Errorf("share card %s missing; run hugo first", cardPath)
	}
	intro := strings.TrimSpace(*text)
	if intro == "" {
		intro = entry.Summary
	}
	a := article{URL: siteURL + entry.URL, Title: entry.Title, Description: entry.Summary}

	fmt.Printf("Post:     %s\nLink:     %s\nIntro:    %s\nCard:     %s (%d KB)\n", a.Title, a.URL, intro, cardPath, len(card)/1024)
	if *dry {
		fmt.Println("\nDry run: nothing was posted.")
		return nil
	}

	token, err := env("LINKEDIN_ACCESS_TOKEN")
	if err != nil {
		return err
	}
	c := newClient(token)
	author, err := c.whoami() // also proves the token still works before uploading anything
	if err != nil {
		return err
	}
	if want := os.Getenv("LINKEDIN_PERSON_URN"); want != "" && want != author {
		return fmt.Errorf("token belongs to %s, but LINKEDIN_PERSON_URN is %s", author, want)
	}
	if a.Thumbnail, err = c.uploadImage(author, card); err != nil {
		return err
	}
	urn, err := c.createPost(author, intro, a)
	if err != nil {
		return err
	}
	postURL := "https://www.linkedin.com/feed/update/" + urn + "/"
	if err := os.WriteFile(mdPath, []byte(withLinkedIn(string(md), postURL)), 0o644); err != nil {
		return fmt.Errorf("posted %s but could not update %s: %w", postURL, mdPath, err)
	}
	fmt.Println("\nPosted:", postURL)
	if out := os.Getenv("GITHUB_OUTPUT"); out != "" {
		if f, err := os.OpenFile(out, os.O_APPEND|os.O_WRONLY, 0o644); err == nil {
			fmt.Fprintf(f, "post_url=%s\n", postURL)
			f.Close()
		}
	}
	return nil
}

func findEntry(indexPath, url string) (indexEntry, error) {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return indexEntry{}, fmt.Errorf("%s missing; run hugo first", indexPath)
	}
	var entries []indexEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return indexEntry{}, err
	}
	for _, e := range entries {
		if e.URL == url && e.Kind == "Story" {
			return e, nil
		}
	}
	return indexEntry{}, fmt.Errorf("%s isn't in the build. Is it still a draft?", url)
}

var linkedInLine = regexp.MustCompile(`(?m)^\s*linkedin\s*=\s*"([^"]+)"`)

// frontMatterLinkedIn returns the post's recorded LinkedIn URL, if any.
func frontMatterLinkedIn(md string) string {
	fm, _, ok := splitFrontMatter(md)
	if !ok {
		return ""
	}
	if m := linkedInLine.FindStringSubmatch(fm); m != nil {
		return m[1]
	}
	return ""
}

// withLinkedIn records the URL under [params], adding the table at the end of the
// front matter when the post doesn't have one yet.
func withLinkedIn(md, url string) string {
	fm, body, ok := splitFrontMatter(md)
	if !ok {
		return md
	}
	line := fmt.Sprintf("  linkedin = %q", url)
	var out []string
	inserted := false
	sc := bufio.NewScanner(strings.NewReader(fm))
	for sc.Scan() {
		out = append(out, sc.Text())
		if strings.TrimSpace(sc.Text()) == "[params]" && !inserted {
			out = append(out, line)
			inserted = true
		}
	}
	if !inserted {
		out = append(out, "[params]", line)
	}
	return "+++\n" + strings.Join(out, "\n") + "\n+++" + body
}

// splitFrontMatter splits TOML front matter (between +++ lines) from the body.
func splitFrontMatter(md string) (fm, body string, ok bool) {
	if !strings.HasPrefix(md, "+++\n") {
		return "", "", false
	}
	end := strings.Index(md[4:], "\n+++")
	if end < 0 {
		return "", "", false
	}
	return md[4 : 4+end], md[4+end+4:], true
}
