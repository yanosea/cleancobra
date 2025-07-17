package model

import (
	"testing"
)

func TestMode_String(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
		want string
	}{
		{
			name: "positive testing - ModeNormal",
			mode: ModeNormal,
			want: "Normal",
		},
		{
			name: "positive testing - ModeInput",
			mode: ModeInput,
			want: "Input",
		},
		{
			name: "positive testing - ModeConfirmation",
			mode: ModeConfirmation,
			want: "Confirmation",
		},
		{
			name: "positive testing - ModeEdit",
			mode: ModeEdit,
			want: "Edit",
		},
		{
			name: "positive testing - unknown mode",
			mode: Mode(999),
			want: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mode.String()
			if got != tt.want {
				t.Errorf("Mode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}