package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          App{Version: "v0.1"},
			expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: errors.New("wrong length")}},
		},
		{
			in: User{
				ID:     "560d2dc4-6889-4de0-9783-2ef32078895e",
				Name:   "Somebody",
				Age:    42,
				Email:  "a@b.c",
				Role:   "admin",
				Phones: []string{"123-456-789", "123--456789"},
				meta:   []byte{},
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 500,
				Body: "012345",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
			// Place your code here.
			_ = tt
		})
	}
}
