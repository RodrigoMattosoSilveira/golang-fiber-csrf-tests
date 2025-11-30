package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestCsrf(t *testing.T) {
	app := setupApp()

	// w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/form", nil)
	resp, err := app.Test(req)
    if err != nil {
        t.Fatalf("app.Test error: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected status 200, got %d", resp.StatusCode)
    }

    body, _ := io.ReadAll(resp.Body)
    if !strings.Contains(string(body), "<h1>Fiber CSRF Form</h1>") {
        t.Fatalf("expected 'No books yet.' in response, got: %s", string(body))
    }

	// Read response body
	// bodyBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Fatalf("failed to read response body: %v", err)
	// }
	html := string(body)

	// Extract CSRF token from HTML
	re := regexp.MustCompile(`name="csrf" value="([^"]+)"`)
	matches := re.FindStringSubmatch(html)
	if len(matches) != 2 {
		t.Fatalf("CSRF token not found in HTML: %q", html)
	}
	token := matches[1]

	// Capture cookies set in GET /form (includes CSRF cookie)
	cookies := resp.Cookies()

	// ---- STEP 2: POST /submit with token and cookies ----
	form := "name=Rodrigo&csrf=" + token
	req = httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Attach cookies to request
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("POST /submit failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from POST /submit, got %d", resp.StatusCode)
	}
}
