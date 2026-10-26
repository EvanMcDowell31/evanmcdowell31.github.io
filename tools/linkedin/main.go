// Command linkedin shares Story posts from this Hugo site to LinkedIn through the
// official Posts API, as an article post carrying the page's share card.
//
//	linkedin auth                 one-time OAuth login; prints the token and member URN
//	linkedin post -slug <slug>    share content/story/<slug>.md (run `hugo` first)
//
// Environment: LINKEDIN_CLIENT_ID and LINKEDIN_CLIENT_SECRET for auth;
// LINKEDIN_ACCESS_TOKEN and LINKEDIN_PERSON_URN for post.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "auth":
		err = runAuth(os.Args[2:])
	case "post":
		err = runPost(os.Args[2:])
	case "-h", "--help", "help":
		usage()
		return
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "linkedin:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `usage:
  linkedin auth [-port 8723]
      Opens LinkedIn's consent page and prints LINKEDIN_ACCESS_TOKEN and
      LINKEDIN_PERSON_URN. Needs LINKEDIN_CLIENT_ID and LINKEDIN_CLIENT_SECRET.

  linkedin post -slug <slug> [-text "intro"] [-dry-run] [-site ../..]
      Shares a built Story post. Run hugo first so public/ has the page and its card.
`)
}

func env(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("%s is not set", name)
	}
	return v, nil
}
