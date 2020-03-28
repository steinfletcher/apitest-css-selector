# apitest-css-selector

Assertions for [apitest](https://github.com/steinfletcher/apitest) using css selectors.

## Examples

### `selector.FirstTextValue`

```go
apitest.New().
	Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<div class="myClass">content</div>`))
		w.WriteHeader(http.StatusOK)
	})).
	Get("/").
	Expect(t).
	Status(http.StatusOK).
	Assert(selector.FirstTextValue(`.myClass`, "content")).
	End()
```

see also `selector.NthTextValue` and `selector.ContainsTextValue`

### `selector.Exists`

```go
apitest.New().
	Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<div data-test-id="product-5">first</div>`))
		w.WriteHeader(http.StatusOK)
	})).
	Get("/").
	Expect(t).
	Status(http.StatusOK).
	Assert(selector.Exists(`div[data-test-id^="product-"]`)).
	End()
```

`selector.Exists` is a variadic function so multiple selectors can be provided, e.g.

```go
selector.Exists(".myClass", `div[data-test-id^="product-"]`, "#myId")
```

### `selector.Selection`

This exposes `goquery`'s Selection api and offers more flexibility over the previous methods

```go
Assert(selector.Selection(".outerClass", func(selection *goquery.Selection) error {
	if test.expectedText != selection.Find(".innerClass").Text() {
	    return fmt.Errorf("text did not match")
	}
	return nil
})).
```

