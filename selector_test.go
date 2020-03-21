package selector_test

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"

	selector "github.com/steinfletcher/apitest-css-selector"
)

func TestSelector_FirstTextValue(t *testing.T) {
	tests := map[string]struct {
		selector     string
		responseBody string
		expected     string
	}{
		"first text value": {
			selector: "h1",
			responseBody: `<html>
				<head>
					<title>My document</title>
				</head>
				<body>
					<h1>Header</h1>
				</body>
			</html>`,
			expected: "Header",
		},
		"first text value with class": {
			selector:     ".myClass",
			responseBody: `<div class="myClass">content</div>`,
			expected:     "content",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apitest.New().
				Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(test.responseBody))
					w.WriteHeader(http.StatusOK)
				})).
				Get("/").
				Expect(t).
				Status(http.StatusOK).
				Assert(selector.FirstTextValue(test.selector, test.expected)).
				End()
		})
	}
}

func TestSelector_NthTextValue(t *testing.T) {
	tests := map[string]struct {
		selector     string
		responseBody string
		expected     string
		n            int
	}{
		"second text value": {
			selector: ".myClass",
			responseBody: `<div>
				<div class="myClass">first</div>
				<div class="myClass">second</div>
			</div>`, expected: "first",
			n: 0,
		},
		"last text value": {
			selector: ".myClass",
			responseBody: `<div>
				<div class="myClass">first</div>
				<div class="myClass">second</div>
			</div>`, expected: "second",
			n: 1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apitest.New().
				Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(test.responseBody))
					w.WriteHeader(http.StatusOK)
				})).
				Get("/").
				Expect(t).
				Status(http.StatusOK).
				Assert(selector.NthTextValue(test.n, test.selector, test.expected)).
				End()
		})
	}
}

func TestSelector_TextValueContains(t *testing.T) {
	tests := map[string]struct {
		selector     string
		responseBody string
		expected     string
	}{
		"text value contains": {
			selector: ".myClass",
			responseBody: `<div>
				<div class="myClass">first</div>
				<div class="myClass">second</div>
			</div>`,
			expected: "second",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apitest.New().
				Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(test.responseBody))
					w.WriteHeader(http.StatusOK)
				})).
				Get("/").
				Expect(t).
				Status(http.StatusOK).
				Assert(selector.ContainsTextValue(test.selector, test.expected)).
				End()
		})
	}
}

func TestSelector_ElementExists(t *testing.T) {
	tests := map[string]struct {
		selector     string
		responseBody string
	}{
		"element exists": {
			selector:     `div[data-test-id^="product-"]`,
			responseBody: `<div data-test-id="product-5">first</div>`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apitest.New().
				Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(test.responseBody))
					w.WriteHeader(http.StatusOK)
				})).
				Get("/").
				Expect(t).
				Status(http.StatusOK).
				Assert(selector.ElementExists(test.selector)).
				End()
		})
	}
}
