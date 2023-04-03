//////
// Shared utils.
//////

package shared

import (
	"errors"
	"testing"
)

func TestPrintErrorMessages(t *testing.T) {
	type args struct {
		errors []error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestPrintErrorMessages",
			args: args{
				errors: []error{
					errors.New("failed to process"),
					errors.New("failed to print"),
					errors.New("failed, no internet connection"),
				},
			},
			want: "failed to process. failed to print. failed, no internet connection",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrintErrorMessages(tt.args.errors...); got != tt.want {
				t.Errorf("PrintErrorMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}
