package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestBaseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *BaseError
		contains []string
	}{
		{
			name: "basic short message",
			err: &BaseError{
				Name:         "TestError",
				ShortMessage: "Something went wrong.",
			},
			contains: []string{"Something went wrong."},
		},
		{
			name: "with details",
			err: &BaseError{
				Name:         "TestError",
				ShortMessage: "Failed.",
				Details:      "underlying reason",
			},
			contains: []string{"Failed.", "Details: underlying reason"},
		},
		{
			name: "with docs path",
			err: &BaseError{
				Name:         "TestError",
				ShortMessage: "Failed.",
				DocsPath:     "/docs/foo",
			},
			contains: []string{"Failed.", "Docs: https://viem.sh/docs/foo"},
		},
		{
			name: "with meta messages",
			err: &BaseError{
				Name:         "TestError",
				ShortMessage: "Failed.",
				MetaMessages: []string{"Hint 1", "Hint 2"},
			},
			contains: []string{"Failed.", "Hint 1", "Hint 2"},
		},
		{
			name: "empty short message defaults",
			err: &BaseError{
				Name: "TestError",
			},
			contains: []string{"An error occurred."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			for _, sub := range tt.contains {
				if !strings.Contains(got, sub) {
					t.Errorf("Error() = %q, want to contain %q", got, sub)
				}
			}
		})
	}
}

func TestBaseError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := &BaseError{
		ShortMessage: "wrapped",
		Cause:        cause,
	}
	if got := errors.Unwrap(err); got != cause {
		t.Errorf("Unwrap() = %v, want %v", got, cause)
	}
}

func TestWalk(t *testing.T) {
	leaf := errors.New("leaf")
	mid := &BaseError{ShortMessage: "mid", Cause: leaf}
	top := &BaseError{ShortMessage: "top", Cause: mid}

	t.Run("nil fn returns leaf", func(t *testing.T) {
		got := Walk(top, nil)
		if got != leaf {
			t.Errorf("Walk(top, nil) = %v, want %v", got, leaf)
		}
	})

	t.Run("fn matches top", func(t *testing.T) {
		got := Walk(top, func(e error) bool {
			_, ok := e.(*BaseError)
			return ok
		})
		if got != top {
			t.Errorf("Walk(top, fn) = %v, want top", got)
		}
	})

	t.Run("fn matches mid", func(t *testing.T) {
		got := Walk(top, func(e error) bool {
			be, ok := e.(*BaseError)
			return ok && be.ShortMessage == "mid"
		})
		if got != mid {
			t.Errorf("Walk(top, fn) = %v, want mid", got)
		}
	})

	t.Run("fn matches none returns nil", func(t *testing.T) {
		got := Walk(top, func(e error) bool { return false })
		if got != nil {
			t.Errorf("Walk(top, fn=false) = %v, want nil", got)
		}
	})

	t.Run("nil err returns nil", func(t *testing.T) {
		if got := Walk(nil, nil); got != nil {
			t.Errorf("Walk(nil, nil) = %v, want nil", got)
		}
	})
}
