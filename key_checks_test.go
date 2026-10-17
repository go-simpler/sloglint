package sloglint

import "testing"

func Test_checkKeyCase(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		caseName  string
		wantFixed string
		wantMatch bool
	}{
		// Snake case tests
		{name: "snake: valid single word", key: "foo", caseName: keyNamingCaseSnake, wantFixed: "", wantMatch: true},
		{name: "snake: valid multi word", key: "foo_bar", caseName: keyNamingCaseSnake, wantFixed: "", wantMatch: true},
		{name: "snake: valid dotted otel", key: "http.status_code", caseName: keyNamingCaseSnake, wantFixed: "", wantMatch: true},
		{name: "snake: valid dotted 3 parts", key: "service.instance.id", caseName: keyNamingCaseSnake, wantFixed: "", wantMatch: true},
		{name: "snake: invalid kebab key", key: "foo-bar", caseName: keyNamingCaseSnake, wantFixed: "foo_bar", wantMatch: false},
		{name: "snake: invalid camel key", key: "fooBar", caseName: keyNamingCaseSnake, wantFixed: "foo_bar", wantMatch: false},
		{name: "snake: invalid pascal key", key: "FooBar", caseName: keyNamingCaseSnake, wantFixed: "foo_bar", wantMatch: false},
		{name: "snake: invalid dotted kebab part", key: "http.status-code", caseName: keyNamingCaseSnake, wantFixed: "http.status_code", wantMatch: false},
		{name: "snake: invalid dotted camel part", key: "http.statusCode", caseName: keyNamingCaseSnake, wantFixed: "http.status_code", wantMatch: false},
		{name: "snake: invalid dotted pascal part", key: "Http.Status_Code", caseName: keyNamingCaseSnake, wantFixed: "http.status_code", wantMatch: false},
		{name: "snake: invalid dotted multi part", key: "service.instance-id.region_name", caseName: keyNamingCaseSnake, wantFixed: "service.instance_id.region_name", wantMatch: false},

		// Kebab case tests
		{name: "kebab: valid single word", key: "foo", caseName: keyNamingCaseKebab, wantFixed: "", wantMatch: true},
		{name: "kebab: valid multi word", key: "foo-bar", caseName: keyNamingCaseKebab, wantFixed: "", wantMatch: true},
		{name: "kebab: valid dotted", key: "http.status-code", caseName: keyNamingCaseKebab, wantFixed: "", wantMatch: true},
		{name: "kebab: valid dotted 3 parts", key: "service.instance.id", caseName: keyNamingCaseKebab, wantFixed: "", wantMatch: true},
		{name: "kebab: invalid snake key", key: "foo_bar", caseName: keyNamingCaseKebab, wantFixed: "foo-bar", wantMatch: false},
		{name: "kebab: invalid dotted snake part", key: "http.status_code", caseName: keyNamingCaseKebab, wantFixed: "http.status-code", wantMatch: false},
		{name: "kebab: invalid dotted camel part", key: "http.statusCode", caseName: keyNamingCaseKebab, wantFixed: "http.status-code", wantMatch: false},

		// Camel case tests
		{name: "camel: valid single word", key: "foo", caseName: keyNamingCaseCamel, wantFixed: "", wantMatch: true},
		{name: "camel: valid multi word", key: "fooBar", caseName: keyNamingCaseCamel, wantFixed: "", wantMatch: true},
		{name: "camel: valid dotted", key: "http.statusCode", caseName: keyNamingCaseCamel, wantFixed: "", wantMatch: true},
		{name: "camel: valid dotted 3 parts", key: "service.instanceId.nodeName", caseName: keyNamingCaseCamel, wantFixed: "", wantMatch: true},
		{name: "camel: invalid snake key", key: "foo_bar", caseName: keyNamingCaseCamel, wantFixed: "fooBar", wantMatch: false},
		{name: "camel: invalid dotted snake part", key: "http.status_code", caseName: keyNamingCaseCamel, wantFixed: "http.statusCode", wantMatch: false},
		{name: "camel: invalid dotted kebab part", key: "http.status-code", caseName: keyNamingCaseCamel, wantFixed: "http.statusCode", wantMatch: false},

		// Pascal case tests
		{name: "pascal: valid single word", key: "Foo", caseName: keyNamingCasePascal, wantFixed: "", wantMatch: true},
		{name: "pascal: valid multi word", key: "FooBar", caseName: keyNamingCasePascal, wantFixed: "", wantMatch: true},
		{name: "pascal: valid dotted", key: "Http.StatusCode", caseName: keyNamingCasePascal, wantFixed: "", wantMatch: true},
		{name: "pascal: invalid snake key", key: "foo_bar", caseName: keyNamingCasePascal, wantFixed: "FooBar", wantMatch: false},
		{name: "pascal: invalid dotted snake parts", key: "http.status_code", caseName: keyNamingCasePascal, wantFixed: "Http.StatusCode", wantMatch: false},
		{name: "pascal: invalid dotted lowercase part", key: "http.StatusCode", caseName: keyNamingCasePascal, wantFixed: "Http.StatusCode", wantMatch: false},

		// Edge cases
		{name: "empty key", key: "", caseName: keyNamingCaseSnake, wantFixed: "", wantMatch: true},
		{name: "trailing dot", key: "http.status_code.", caseName: keyNamingCaseSnake, wantFixed: "http.status_code", wantMatch: false},
		{name: "leading dot", key: ".http.status_code", caseName: keyNamingCaseSnake, wantFixed: "http.status_code", wantMatch: false},
		{name: "consecutive dots", key: "http..status_code", caseName: keyNamingCaseSnake, wantFixed: "http.status_code", wantMatch: false},
		{name: "single dot", key: ".", caseName: keyNamingCaseSnake, wantFixed: "", wantMatch: false},
		{name: "double dot", key: "..", caseName: keyNamingCaseSnake, wantFixed: "", wantMatch: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotFixed, gotMatch := checkKeyCase(tc.key, tc.caseName)
			if gotMatch != tc.wantMatch {
				t.Errorf("checkKeyCase(%q, %q) match = %v; want %v", tc.key, tc.caseName, gotMatch, tc.wantMatch)
			}
			if gotFixed != tc.wantFixed {
				t.Errorf("checkKeyCase(%q, %q) fixed = %q; want %q", tc.key, tc.caseName, gotFixed, tc.wantFixed)
			}
		})
	}
}
