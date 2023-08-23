package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/omakoto/bashcomp"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/src/cmd/oauther/oauth"
)

var (
	clientId     = flag.String("client-id", "", "Client ID")
	clientSecret = flag.String("client-secret", "", "Client secret")
	scopes       = flag.String("scopes", "", "Oauth scopes")
	newtoken     = flag.Bool("new", false, "Clear cache and request new token (useful to change accounts)")
	cachePath    = flag.String("cache-path", "", "Cache path")
)

const (
	REDIRECT_URL = "http://localhost:8080/"
)

// openURL opens a browser window to the specified location.
// This code originally appeared at:
//
//	http://stackoverflow.com/questions/10377243/how-can-i-launch-a-process-that-is-not-a-file-in-go
func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func buildConfig() (*oauth.Config, error) {
	if *clientId == "" {
		log.Fatalf("You must provide an oauth client ID with -client-id")
	}

	if *clientSecret == "" {
		log.Fatalf("You must provide an oauth client secret with -client-secret")
	}

	if *scopes == "" {
		log.Fatalf("You must provide an oauth client secret with -scopes")
	}

	openFlag := os.O_RDWR | os.O_CREATE
	if *newtoken {
		openFlag |= os.O_TRUNC
	}

	if *cachePath == "" {
		*cachePath = filepath.Join(getHomeDir(), ".oauther")
	}
	err := os.MkdirAll(*cachePath, 0700)
	common.Check(err, "MkdirAll failed")

	cacheFile := filepath.Join(*cachePath, *clientId)
	_, err = os.OpenFile(cacheFile, openFlag, 0600)
	common.Check(err, "OpenFile failed")

	return &oauth.Config{
		ClientId:       *clientId,
		ClientSecret:   *clientSecret,
		Scope:          *scopes,
		AuthURL:        "https://accounts.google.com/o/oauth2/auth",
		TokenURL:       "https://accounts.google.com/o/oauth2/token",
		RedirectURL:    REDIRECT_URL,
		TokenCache:     oauth.CacheFile(cacheFile),
		AccessType:     "offline",
		ApprovalPrompt: "force",
	}, nil
}

// startWebServer starts a web server that listens on http://localhost:8080.
// The webserver waits for an oauth code in the three-legged auth flow.
func startWebServer() (codeCh chan string, err error) {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}
	codeCh = make(chan string)
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeCh <- code // send code to OAuth flow
		listener.Close()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Received code: %v\r\nYou can now safely close this browser window.", code)
	}))

	return codeCh, nil
}

func getToken() (string, error) {
	config, err := buildConfig()
	if err != nil {
		msg := fmt.Sprintf("Cannot read configuration file: %v", err)
		return "", errors.New(msg)
	}

	transport := &oauth.Transport{Config: config}

	// Try to read the token from the cache file.
	// If an error occurs, do the three-legged OAuth flow because
	// the token is invalid or doesn't exist.
	token, err := config.TokenCache.Token()
	if err != nil {
		// Start web server.
		// This is how this program receives the authorization code
		// when the browser redirects.
		codeCh, err := startWebServer()
		if err != nil {
			return "", err
		}

		// Open url in browser
		url := config.AuthCodeURL("")
		// fmt.Println("URL=" + url)
		err = openURL(url)
		if err != nil {
			log.Println("Visit the URL below to get a code.",
				" This program will pause until the site is visted.")
		} else {
			log.Println("Your browser has been opened to an authorization URL.",
				" This program will resume once authorization has been provided.")
		}
		log.Println(url)
		// fmt.Println(url)

		// Wait for the web server to get the code.
		code := <-codeCh

		// This code caches the authorization code on the local
		// filesystem, if necessary, as long as the TokenCache
		// attribute in the config is set.
		token, err = transport.Exchange(code)
		if err != nil {
			return "", err
		}
	}

	transport.Token = token

	newToken, err := transport.GetAccessToken()
	common.Check(err, "GetAccessToken failed")

	common.Dump("token=", newtoken)

	return newToken, nil
}

func main() {
	flag.Parse()

	bashcomp.HandleBashCompletion()

	token, err := getToken()
	if err != nil {
		log.Fatalf("Error building OAuth client: %v", err)
	}
	fmt.Println(token)
}
