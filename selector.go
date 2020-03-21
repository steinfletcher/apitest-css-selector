package selector

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/steinfletcher/apitest"
)

type selectionMatcher func(i int, selection *goquery.Selection) bool

func FirstTextValue(selection string, expectedTextValue string) apitest.Assert {
	return newSelectionAssert(selection, expectedTextValue, func(i int, selection *goquery.Selection) bool {
		if i == 0 {
			if selection.Text() == expectedTextValue {
				return true
			}
		}
		return false
	})
}

func NthTextValue(n int, selection string, expectedTextValue string) apitest.Assert {
	return newSelectionAssert(selection, expectedTextValue, func(i int, selection *goquery.Selection) bool {
		if i == n {
			if selection.Text() == expectedTextValue {
				return true
			}
		}
		return false
	})
}

func ContainsTextValue(selection string, expectedTextValue string) apitest.Assert {
	return newSelectionAssert(selection, expectedTextValue, func(i int, selection *goquery.Selection) bool {
		if selection.Text() == expectedTextValue {
			return true
		}
		return false
	})
}

func newSelectionAssert(selection string, expectedTextValue string, matcher selectionMatcher) apitest.Assert {
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
			return fmt.Errorf("result did not contain expected value '%s' for selector '%s'",
				expectedTextValue, selection)
		}

		return nil
	}
}
