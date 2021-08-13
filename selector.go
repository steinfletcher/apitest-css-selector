package selector

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type selectionMatcher func(i int, selection *goquery.Selection) bool

func DataTestID(value string) string {
	return fmt.Sprintf(`[data-test-id="%s"]`, value)
}

func FirstTextValue(selection string, expectedTextValue string) func(*http.Response, *http.Request) error {
	return newAssertSelection(selection, func(i int, selection *goquery.Selection) bool {
		if i == 0 {
			if selection.Text() == expectedTextValue {
				return true
			}
		}
		return false
	})
}

func NthTextValue(n int, selection string, expectedTextValue string) func(*http.Response, *http.Request) error {
	return newAssertSelection(selection, func(i int, selection *goquery.Selection) bool {
		if i == n {
			if selection.Text() == expectedTextValue {
				return true
			}
		}
		return false
	})
}

func ContainsTextValue(selection string, expectedTextValue string) func(*http.Response, *http.Request) error {
	return newAssertSelection(selection, func(i int, selection *goquery.Selection) bool {
		if strings.Contains(selection.Text(), expectedTextValue) {
			return true
		}
		return false
	})
}

func Selection(selection string, selectionFunc func(*goquery.Selection) error) func(*http.Response, *http.Request) error {
	return func(response *http.Response, request *http.Request) error {
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			return err
		}
		return selectionFunc(doc.Find(selection))
	}
}

func Exists(selections ...string) func(*http.Response, *http.Request) error {
	return expectExists(true, selections...)
}

func NotExists(selections ...string) func(*http.Response, *http.Request) error {
	return expectExists(false, selections...)
}

func TextExists(text string) func(*http.Response, *http.Request) error {
	return func(response *http.Response, request *http.Request) error {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if !strings.Contains(string(bodyBytes), text) {
			return fmt.Errorf("document did not contain '%v'", text)
		}

		return nil
	}
}

func expectExists(exists bool, selections ...string) func(*http.Response, *http.Request) error {
	return func(response *http.Response, request *http.Request) error {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		for _, selection := range selections {
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
			if err != nil {
				return err
			}

			var found bool
			doc.Find(selection).Each(func(i int, selection *goquery.Selection) {
				found = true
			})

			if found != exists {
				return fmt.Errorf("expected found='%v' for selector '%s'", exists, selection)
			}
		}

		return nil
	}
}

func newAssertSelection(selection string, matcher selectionMatcher) func(*http.Response, *http.Request) error {
	return func(response *http.Response, request *http.Request) error {
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			return err
		}

		var found bool
		doc.Find(selection).Each(func(i int, selection *goquery.Selection) {
			if matcher(i, selection) {
				found = true
			}
		})

		if !found {
			return fmt.Errorf("did not find expected value for selector '%s'", selection)
		}

		return nil
	}
}
