package compute

import (
	"testing"
)

func TestParser_ParseGet(t *testing.T) {
	parser := Parser{}

	tests := []struct {
		name      string
		expectErr bool
	}{
		{
			name:      "GET 1",
			expectErr: false,
		},
		{
			name:      "get 1",
			expectErr: false,
		},
		{
			name:      "get 1 2 3",
			expectErr: true,
		},
		{
			name:      "get",
			expectErr: true,
		},
		{
			name:      "got",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.name)

			if tt.expectErr && err == nil {
				t.Fatalf("expected error but found none")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("did not expect error but got: %v", err)
			}
		})
	}
}

func TestParser_ParseSet(t *testing.T) {
	parser := Parser{}
	tests := []struct {
		name      string
		expectErr bool
	}{
		{
			name:      "set 1 2",
			expectErr: false,
		},
		{
			name:      "SET 1 2",
			expectErr: false,
		},
		{
			name:      "SET 1",
			expectErr: true,
		},
		{
			name:      "set",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.name)

			if tt.expectErr && err == nil {
				t.Fatalf("expected error but found none")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("did not expect error but got: %v", err)
			}
		})
	}
}

func TestParser_ParseDelete(t *testing.T) {
	parser := Parser{}
	tests := []struct {
		name      string
		expectErr bool
	}{
		{
			name:      "del 1",
			expectErr: false,
		},
		{
			name:      "DEL 1",
			expectErr: false,
		},
		{
			name:      "del",
			expectErr: true,
		},
		{
			name:      "del 1 2 3",
			expectErr: true,
		},
		{
			name:      "delete 1",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.name)

			if tt.expectErr && err == nil {
				t.Fatalf("expected error but found none")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("did not expect error but got: %v", err)
			}
		})
	}
}
