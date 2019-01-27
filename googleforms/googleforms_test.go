package googleforms

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func equals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestFormCombos(t *testing.T) {
	tests := []struct {
		ExpectedAction string
		ExpectedFields []string
		Response       string
	}{
		{
			ExpectedAction: "https://docs.google.com/forms/d/e/flasdjfaljsdfjsa/formResponse",
			ExpectedFields: []string{"entry.1390350291"},
			Response: `<form action="https://docs.google.com/forms/d/e/flasdjfaljsdfjsa/formResponse">
			<input name="entry.1390350291">`,
		},
		{
			ExpectedAction: "https://docs.google.com/forms/d/e/flasdjfaljsdfjsa/formResponse",
			ExpectedFields: []string{"entry.1390350291", "entry.1390350292"},
			Response: `<form action="https://docs.google.com/forms/d/e/flasdjfaljsdfjsa/formResponse">
			<input name="unrelated">
			<input name="entry.1390350291">
			<input name="entry.1390350292">`,
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, test.Response)
		}))
		defer ts.Close()
		form := NewForm(ts.URL)
		err := form.Init()
		if err != nil {
			t.Error(err)
			continue
		}
		if form.ActionURL != test.ExpectedAction {
			t.Error(errors.New(fmt.Sprintf("mismatched action:\nhave:%v\nwant:%v", form.ActionURL, test.ExpectedAction)))
		}
		if !equals(form.FieldKeys, test.ExpectedFields) {
			t.Error(errors.New(fmt.Sprintf("mismatched fields:\nhave:%v\nwant:%v", form.FieldKeys, test.ExpectedFields)))
		}
	}
}
