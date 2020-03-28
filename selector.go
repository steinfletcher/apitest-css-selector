package selector

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/steinfletcher/apitest"
)

type selectionMatcher func(i int, selection *goquery.Selection) bool

func FirstTextValue(selection string, expectedTextValue string) apitest.Assert {
	return newAssertSelection(selection, func(i int, selection *goquery.Selection) bool {
		if i == 0 {
			if selection.Text() == expectedTextValue {
				return true
			}
		}
		return false
	})
}

func NthTextValue(n int, selection string, expectedTextValue string) apitest.Assert {
	return newAssertSelection(selection, func(i int, selection *goquery.Selection) bool {
		if i == n {
			if selection.Text() == expectedTextValue {
				return true
			}
		}
		return false
	})
}

func ContainsTextValue(selection string, expectedTextValue string) apitest.Assert {
	return newAssertSelection(selection, func(i int, selection *goquery.Selection) bool {
		if strings.Contains(selection.Text(), expectedTextValue) {
			return true
		}
		return false
	})
}

func Selection(selection string, selectionFunc func(*goquery.Selection) error) apitest.Assert {
	return func(response *http.Response, request *http.Request) error {
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			return err
		}
		return selectionFunc(doc.Find(selection))
	}
}

func Exists(selections ...string) apitest.Assert {
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

			if !found {
				return fmt.Errorf("did not find expected value for selector '%s'", selection)
			}
		}

		return nil
	}
}

func newAssertSelection(selection string, matcher selectionMatcher) apitest.Assert {
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
