# apitest-css-selector

Assertions for [apitest](https://github.com/steinfletcher/apitest) using css selectors.

## Examples

`selector.FirstTextValue`

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

`selector.ElementExists`

```go
apitest.New().
	Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<div data-test-id="product-5">first</div>`))
		w.WriteHeader(http.StatusOK)
	})).
	Get("/").
	Expect(t).
	Status(http.StatusOK).
	Assert(selector.ElementExists(`div[data-test-id^="product-"]`)).
	End()
```