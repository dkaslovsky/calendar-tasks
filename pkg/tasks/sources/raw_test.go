package sources

import "testing"

func TestLoadLine(t *testing.T) {
	tests := map[string]struct {
		line     string
		expected RawLine
	}{
		"valid": {
			line: "foo:bar",
			expected: RawLine{
				Date: "foo",
				Text: "bar",
			},
		},
		"valid with spaces": {
			line: "foo:  bar",
			expected: RawLine{
				Date: "foo",
				Text: "bar",
			},
		},
		"multiple delimiters": {
			line: "foo:bar:baz",
			expected: RawLine{
				Date: "foo",
				Text: "bar:baz",
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result, err := LoadLine(test.line)
			if err != nil {
				t.Fatalf("unexpected non-nil error: %v", err)
			}
			if result.Date != test.expected.Date || result.Text != test.expected.Text {
				t.Fatalf("result %v does not equal expected %v", result, test.expected)
			}
		})
	}
}

func TestLoadLineError(t *testing.T) {
	tests := map[string]struct {
		line string
	}{
		"empty": {
			line: "",
		},
		"empty with spaces": {
			line: "    ",
		},
		"no delimiter": {
			line: "foobar",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			_, err := LoadLine(test.line)
			if err == nil {
				t.Fatal("unexpected nil error")
			}
		})
	}
}
