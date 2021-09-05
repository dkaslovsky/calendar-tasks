package sources

import "testing"

func TestLoadLine(t *testing.T) {
	tests := map[string]struct {
		line      string
		expected  RawLine
		shouldErr bool
	}{
		"empty": {
			line:      "",
			shouldErr: true,
		},
		"empty with spaces": {
			line:      "    ",
			shouldErr: true,
		},
		"no delimiter": {
			line:      "foobar",
			shouldErr: true,
		},
		"multiple delimiters": {
			line: "foo:bar:baz",
			expected: RawLine{
				Date: "foo",
				Text: "bar:baz",
			},
			shouldErr: false,
		},
		"valid": {
			line: "foo:bar",
			expected: RawLine{
				Date: "foo",
				Text: "bar",
			},
			shouldErr: false,
		},
		"valid with spaces": {
			line: "foo:  bar",
			expected: RawLine{
				Date: "foo",
				Text: "bar",
			},
			shouldErr: false,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result, err := LoadLine(test.line)
			assertShouldError(t, test.shouldErr, err)
			if test.shouldErr {
				return
			}
			if result.Date != test.expected.Date || result.Text != test.expected.Text {
				t.Fatalf("result %v does not equal expected %v", result, test.expected)
			}
		})
	}
}
