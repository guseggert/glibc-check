package glibccheck

import (
	"strings"
	"testing"
)

func TestGLIBCVersions_FindViolations(t *testing.T) {
	cases := []struct {
		name              string
		versions          GLIBCVersions
		expr              string
		expected          GLIBCVersions
		expectErrContains string
	}{
		{
			name: "happy case",
			versions: GLIBCVersions([]GLIBCVersion{
				{Full: "2.1", Major: 2, Minor: 1},
				{Full: "2.3.4", Major: 2, Minor: 3, Patch: 4},
			}),
			expr: "major == 2 && minor >= 1",
		},
		{
			name: "finds violations",
			versions: GLIBCVersions([]GLIBCVersion{
				{Full: "2.1", Major: 2, Minor: 1},
				{Full: "2.3.4", Major: 2, Minor: 3, Patch: 4},
			}),
			expr: "major != 2",
			expected: GLIBCVersions([]GLIBCVersion{
				{Full: "2.1", Major: 2, Minor: 1},
				{Full: "2.3.4", Major: 2, Minor: 3, Patch: 4},
			}),
		},
		{
			name:              "returns error on invalid expr",
			expr:              "&",
			versions:          GLIBCVersions([]GLIBCVersion{{Full: "2.1", Major: 2, Minor: 1}}),
			expectErrContains: "parsing error",
		},
		{
			name:              "returns error if eval result is not a bool",
			expr:              "3",
			versions:          GLIBCVersions([]GLIBCVersion{{Full: "2.1", Major: 2, Minor: 1}}),
			expectErrContains: "did not evaluate to a boolean",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			violations, err := c.versions.FindViolations(c.expr)
			if c.expectErrContains != "" {
				if err == nil {
					t.Fatalf("expected error but none returned")
				}
				if !strings.Contains(err.Error(), c.expectErrContains) {
					t.Fatalf("expected error '%s' to contain '%s'", err.Error(), c.expectErrContains)
				}
			}
			if len(violations) != len(c.expected) {
				t.Fatalf("expected %d violations but got %d", len(c.expected), len(violations))
			}
			for i, v := range violations {
				if v.String() != c.expected[i].String() {
					t.Fatalf("expected violation %d to be %s, but was %s", i, c.expected[i], v)
				}
			}
		})
	}
}

func TestParseGLIBCVersion(t *testing.T) {
	cases := []struct {
		name              string
		str               string
		expected          GLIBCVersion
		expectErrContains string
	}{
		{
			name:     "happy case",
			str:      "1.2.3",
			expected: GLIBCVersion{Full: "1.2.3", Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:     "no patch version",
			str:      "1.2",
			expected: GLIBCVersion{Full: "1.2", Major: 1, Minor: 2, Patch: 0},
		},
		{
			name:              "invalid format",
			str:               "&&&&&",
			expectErrContains: "unable to parse",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, err := ParseGLIBCVersion(c.str)
			if c.expectErrContains != "" {
				if err == nil {
					t.Fatalf("expected error but none returned")
				}
				if !strings.Contains(err.Error(), c.expectErrContains) {
					t.Fatalf("expected error '%s' to contain '%s'", err.Error(), c.expectErrContains)
				}
			}
			if v.Full != c.expected.Full {
				t.Fatalf("expected full to be %s but was %s", c.expected.Full, v.Full)
			}
			if v.Major != c.expected.Major {
				t.Fatalf("expected major to be %d but was %d", c.expected.Major, v.Major)
			}

			if v.Minor != c.expected.Minor {
				t.Fatalf("expected minor to be %d but was %d", c.expected.Minor, v.Minor)
			}

			if v.Patch != c.expected.Patch {
				t.Fatalf("expected patch to be %d but was %d", c.expected.Patch, v.Patch)
			}

		})
	}
}
