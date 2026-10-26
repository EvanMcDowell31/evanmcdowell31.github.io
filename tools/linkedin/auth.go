package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

// runAuth performs LinkedIn's 3-legged OAuth flow with a one-shot local callback server.
func runAuth(args []string) error {
	fs := flag.NewFlagSet("auth", flag.ExitOnError)
	port := fs.Int("port", 8723, "local callback port; register http://localhost:<port>/callback in the LinkedIn app")
	fs.Parse(args)

	id, err := env("LINKEDIN_CLIENT_ID")
	if err != nil {
		return err
	}
	secret, err := env("LINKEDIN_CLIENT_SECRET")
	if err != nil {
		return err
	}
	redirect := fmt.Sprintf("http://localhost:%d/callback", *port)
	stateBytes := make([]byte, 16)
	rand.Read(stateBytes)
	state := hex.EncodeToString(stateBytes)

	authURL := "https://www.linkedin.com/oauth/v2/authorization?" + url.Values{
		"response_type": {"code"},
		"client_id":     {id},
		"redirect_uri":  {redirect},
		"state":         {state},
		"scope":         {"openid profile w_member_social"},
	}.Encode()

	codes := make(chan string, 1)
	errs := make(chan error, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case q.Get("state") != state:
			http.Error(w, "state mismatch", http.StatusBadRequest)
			errs <- fmt.Errorf("state mismatch in callback")
		case q.Get("error") != "":
			http.Error(w, q.Get("error_description"), http.StatusBadRequest)
			errs <- fmt.Errorf("LinkedIn returned %s: %s", q.Get("error"), q.Get("error_description"))
		default:
			fmt.Fprintln(w, "Connected. You can close this tab and return to the terminal.")
			codes <- q.Get("code")
		}
	})
	srv := &http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", *port), Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- err
		}
	}()
	defer srv.Shutdown(context.Background())

	fmt.Println("Opening LinkedIn to approve access. If nothing opens, visit:\n\n  " + authURL + "\n")
	exec.Command("open", authURL).Start()

	var code string
	select {
	case code = <-codes:
	case err := <-errs:
		return err
	case <-time.After(10 * time.Minute):
		return fmt.Errorf("timed out waiting for LinkedIn approval")
	}

	resp, err := http.PostForm("https://www.linkedin.com/oauth/v2/accessToken", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirect},
		"client_id":     {id},
		"client_secret": {secret},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var tok struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Error       string `json:"error_description"`
	}
	json.NewDecoder(resp.Body).Decode(&tok)
	if tok.AccessToken == "" {
		return fmt.Errorf("token exchange failed (%s): %s", resp.Status, tok.Error)
	}

	urn, err := newClient(tok.AccessToken).whoami()
	if err != nil {
		return err
	}
	expires := time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second).Format("Jan 2, 2006")
	repo := "EvanMcDowell31/evanmcdowell31.github.io"
	fmt.Printf(`Connected as %s. The token expires %s.

Save both as repo secrets (with gh, or Settings → Secrets and variables → Actions):

  gh secret set LINKEDIN_PERSON_URN -R %s -b '%s'
  gh secret set LINKEDIN_ACCESS_TOKEN -R %s   # paste the token when prompted

LINKEDIN_ACCESS_TOKEN=%s
`, urn, expires, repo, urn, repo, strings.TrimSpace(tok.AccessToken))
	return nil
}
