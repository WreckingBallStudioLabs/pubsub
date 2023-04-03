//////
// Shared resources for tests.
//////

package shared

import (
	"errors"
	"net/http"
	"testing"
)

func TestCreateHTTPTestServer(t *testing.T) {
	type args struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestCreateHTTPTestServer - should work",
			args: args{
				statusCode: http.StatusOK,
				body:       "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CreateHTTPTestServer(tt.args.statusCode, nil, nil, tt.args.body)
			defer s.Close()

			if s == nil {
				t.Errorf("CreateHTTPTestServer() = %v, want %v", s, !tt.wantErr)
			}

			if s.URL == "" {
				t.Errorf("CreateHTTPTestServer() = %v, want %v", s.URL, !tt.wantErr)
			}

			if s.Client() == nil {
				t.Errorf("CreateHTTPTestServer() = %v, want %v", s.Client(), !tt.wantErr)
			}

			r, err := s.Client().Get(s.URL)
			if err != nil {
				t.Errorf("CreateHTTPTestServer() = %v, want %v", err, !tt.wantErr)
			}

			defer r.Body.Close()

			if r.StatusCode != tt.args.statusCode {
				t.Errorf("CreateHTTPTestServer() = %v, want %v", r.StatusCode, tt.args.statusCode)
			}

			if r.Body == nil {
				t.Errorf("CreateHTTPTestServer() = %v, want %v", r.Body, !tt.wantErr)
			}
		})
	}
}

func TestErrorContains(t *testing.T) {
	// Define a test case.
	testCases := []struct {
		name     string
		err      error
		text     string
		expected bool
	}{
		{
			name:     "error contains text",
			err:      errors.New("this is an error message"),
			text:     "error",
			expected: true,
		},
		{
			name:     "error does not contain text",
			err:      errors.New("this is an error message"),
			text:     "foo",
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			text:     "error",
			expected: false,
		},
	}

	// Test each case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Check if the error contains the expected text.
			result := ErrorContains(tc.err, tc.text)

			// Check that the result matches the expected value.
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestIsEnvironment(t *testing.T) {
	type args struct {
		environments []string
	}
	tests := []struct {
		name string
		args args
		env  string
		want bool
	}{
		{
			name: "TestIsEnvironment - should work",
			args: args{
				environments: []string{"testing", "integration"},
			},
			env:  "integration",
			want: true,
		},
		{
			name: "TestIsEnvironment - should not work",
			args: args{
				environments: []string{"production"},
			},
			env: "integration",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ENVIRONMENT", tt.env)

			if got := IsEnvironment(tt.args.environments...); got != tt.want {
				t.Errorf("IsEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
