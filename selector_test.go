package selector_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/steinfletcher/apitest"

	selector "github.com/steinfletcher/apitest-css-selector"
)

func TestTextExists(t *testing.T) {
	apitest.New().
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<html>
			<head>
				<title>My document</title>
			</head>
			<body>
			<h1>Header</h1>
			<p>Some text to match on</p>
			</body>
			</html>`,
			))
			w.WriteHeader(http.StatusOK)
		}).
		Get("/").
		Expect(t).
		Status(http.StatusOK).
		Assert(selector.TextExists("Some text to match on")).
		End()
}

func TestWithDataTestID(t *testing.T) {
	apitest.New().
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<html>
		<head>
			<title>My document</title>
		</head>
		<body>
		<h1>Header</h1>
		<div data-test-id="some-test-id">
			<div>some content</div>
		</div>
		</body>
		</html>`,
			))
			w.WriteHeader(http.StatusOK)
		}).
		Get("/").
		Expect(t).
		Status(http.StatusOK).
		Assert(selector.ContainsTextValue(selector.DataTestID("some-test-id"), "some content")).
		End()
}

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
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(test.responseBody))
					w.WriteHeader(http.StatusOK)
				}).
				Get("/").
				Expect(t).
				Status(http.StatusOK).
				Assert(selector.FirstTextValue(test.selector, test.expected)).
				End()
		})
	}
}

func TestSelector_FirstTextValue_NoMatch(t *testing.T) {
	verifier := &mockVerifier{
		EqualMock: func(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
			expectedError := "did not find expected value for selector '.myClass'"
			if actual.(error).Error() != expectedError {
				t.Fatalf("actual was unexpected: %v", actual)
			}
			return true
		},
	}

	apitest.New().
		Verifier(verifier).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<div class="myClass">content</div>`))
			w.WriteHeader(http.StatusOK)
		}).
		Get("/").
		Expect(t).
		Assert(selector.FirstTextValue(".myClass", "notContent")).
		End()
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
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(test.responseBody))
					w.WriteHeader(http.StatusOK)
				}).
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
				<div class="myClass">something second</div>
			</div>`,
			expected: "second",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apitest.New().
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(test.responseBody))
					w.WriteHeader(http.StatusOK)
				}).
				Get("/").
				Expect(t).
				Status(http.StatusOK).
				Assert(selector.ContainsTextValue(test.selector, test.expected)).
				End()
		})
	}
}

func TestSelector_Exists_NoMatch(t *testing.T) {
	verifier := &mockVerifier{
		EqualMock: func(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
			expectedError := "expected found='true' for selector '.myClass'"
			if actual.(error).Error() != expectedError {
				t.Fatal()
			}
			return true
		},
	}

	apitest.New().
		Verifier(verifier).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<div class="someClass">content</div>`))
			w.WriteHeader(http.StatusOK)
		}).
		Get("/").
		Expect(t).
		Assert(selector.Exists(".myClass")).
		End()
}

func TestSelector_Exists(t *testing.T) {
	tests := map[string]struct {
		exists       bool
		selector     []string
		responseBody string
	}{
		"exists": {
			exists:   true,
			selector: []string{`div[data-test-id^="product-"]`},
		},
		"multiple exists": {
			exists:   true,
			selector: []string{`div[data-test-id^="product-"]`, `.otherClass`},
		},
		"not exists": {
			exists:   false,
			selector: []string{`div[data-test-id^="product-4"]`},
		},
		"multiple not exists": {
			exists:   false,
			selector: []string{`div[data-test-id^="product-4"]`, `.notExistClass`},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sel := selector.NotExists(test.selector...)
			if test.exists {
				sel = selector.Exists(test.selector...)
			}
			apitest.New().
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(`<div>
					<div class="myClass">first</div>
					<div class="otherClass">something second</div>
					<div data-test-id="product-5">first</div>
				</div>`))
					w.WriteHeader(http.StatusOK)
				}).
				Get("/").
				Expect(t).
				Assert(sel).
				End()
		})
	}
}

func TestSelector_Selection(t *testing.T) {
	tests := map[string]struct {
		selector      string
		selectionFunc func(*goquery.Selection) error
		responseBody  string
		expectedText  string
	}{
		"with selection": {
			selector: `div[data-test-id^="product-"]`,
			responseBody: `<div>
				<div class="otherClass">something second</div>
				<div data-test-id="product-5">
					<div class="myClass">expectedText</div>
				</div>
			</div>`,
			expectedText: "expectedText",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apitest.New().
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(test.responseBody))
					w.WriteHeader(http.StatusOK)
				}).
				Get("/").
				Expect(t).
				Status(http.StatusOK).
				Assert(selector.Selection(test.selector, func(selection *goquery.Selection) error {
					if test.expectedText != selection.Find(".myClass").Text() {
						return fmt.Errorf("text did not match")
					}
					return nil
				})).
				End()
		})
	}
}

func TestSelector_Selection_NotMatch(t *testing.T) {
	verifier := &mockVerifier{
		EqualMock: func(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
			expectedError := "text did not match"
			if actual.(error).Error() != expectedError {
				t.Fatalf("actual was unexpected: %v", actual)
			}
			return true
		},
	}

	tests := map[string]struct {
		selector      string
		selectionFunc func(*goquery.Selection) error
		responseBody  string
		expectedText  string
	}{
		"with selection": {
			selector: `div[data-test-id^="product-"]`,
			responseBody: `<div>
				<div class="otherClass">something second</div>
				<div data-test-id="product-5">
					<div class="myClass">notExpectedText</div>
				</div>
			</div>`,
			expectedText: "expectedText",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apitest.New().
				Verifier(verifier).
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(test.responseBody))
					w.WriteHeader(http.StatusOK)
				}).
				Get("/").
				Expect(t).
				Assert(selector.Selection(test.selector, func(selection *goquery.Selection) error {
					if test.expectedText != selection.Find(".myClass").Text() {
						return fmt.Errorf("text did not match")
					}
					return nil
				})).
				End()
		})
	}
}

func TestSelector_MultipleExists_NoMatch(t *testing.T) {
	verifier := &mockVerifier{
		EqualMock: func(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
			expectedError := "expected found='true' for selector '.myClass'"
			if actual.(error).Error() != expectedError {
				t.Fatalf("actual was unexpected: %v", actual)
			}
			return true
		},
	}

	apitest.New().
		Verifier(verifier).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`<div class="someClass">content</div>`))
			w.WriteHeader(http.StatusOK)
		}).
		Get("/").
		Expect(t).
		Assert(selector.Exists(".someClass", ".myClass")).
		End()
}

type mockVerifier struct {
	EqualInvoked bool
	EqualMock    func(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool

	JSONEqInvoked bool
	JSONEqMock    func(t *testing.T, expected string, actual string, msgAndArgs ...interface{}) bool

	FailInvoked bool
	FailMock    func(t *testing.T, failureMessage string, msgAndArgs ...interface{}) bool
}

func (m *mockVerifier) Equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	m.EqualInvoked = true
	return m.EqualMock(t, expected, actual, msgAndArgs)
}

func (m *mockVerifier) JSONEq(t *testing.T, expected string, actual string, msgAndArgs ...interface{}) bool {
	m.JSONEqInvoked = true
	return m.JSONEqMock(t, expected, actual, msgAndArgs)
}

func (m *mockVerifier) Fail(t *testing.T, failureMessage string, msgAndArgs ...interface{}) bool {
	m.FailInvoked = true
	return m.FailMock(t, failureMessage, msgAndArgs)
}
