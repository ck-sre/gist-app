package main

import (
	"bytes"
	"fmt"
	"gistapp.ck89.net/internal/dblayer/mocks"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"
)

var csrfTokenRegex = regexp.MustCompile(`<input type='hidden' name='csrf' value='(.+)'>`)

type tServer struct {
	*httptest.Server
}

func newTServer(t *testing.T, c http.Handler) *tServer {
	tsvr := httptest.NewTLSServer(c)

	cJar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	tsvr.Client().Jar = cJar

	tsvr.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &tServer{tsvr}
}
func (tsvr *tServer) retrieve(t *testing.T, uPath string) (int, http.Header, string) {
	r, err := tsvr.Client().Get(tsvr.URL + uPath)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	bd, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	bd = bytes.TrimSpace(bd)
	return r.StatusCode, r.Header, string(bd)

}

func newTestMission(t *testing.T) *mission {
	tmplCache, err := newTmplCache()
	if err != nil {
		t.Fatal(err)
	}

	frmDcdr := form.NewDecoder()
	snMngr := scs.New()
	snMngr.Lifetime = 12 * time.Hour
	snMngr.Cookie.Secure = true

	return &mission{
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		gists:     &mocks.Gistdblayer{},
		usrs:      &mocks.UserLayer{},
		tmplCache: tmplCache,
		formDcdr:  frmDcdr,
		snMgr:     snMngr,
	}
}

func newTestServer(t *testing.T) *httptest.Server {
	msn := newTestMission(t)
	return httptest.NewServer(msn.paths())
}

func getCSRFToken(t *testing.T, bd string) string {
	matches := csrfTokenRegex.FindStringSubmatch(bd)
	fmt.Println(len(matches))
	if len(matches) < 2 {
		t.Fatal("no csrf token found")
	}
	return html.UnescapeString(string(matches[1]))
}

func (tsvr *tServer) postForm(t *testing.T, uPath string, form url.Values) (int, http.Header, string) {

	r, err := tsvr.Client().PostForm(tsvr.URL+uPath, form)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	bd, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	bd = bytes.TrimSpace(bd)
	return r.StatusCode, r.Header, string(bd)
}
