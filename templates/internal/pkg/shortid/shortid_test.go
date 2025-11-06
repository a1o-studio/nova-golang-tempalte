package shortid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	testCases := []int64{
		0,
		1,
		42,
		123456,
		99999999,
	}

	for _, id := range testCases {
		code, err := Encode(id)
		assert.NoError(t, err, "unexpected error encoding %d", id)
		assert.NotEmpty(t, code, "encoded string should not be empty for id %d", id)

		decoded, err := Decode(code)
		assert.NoError(t, err, "unexpected error decoding %s", code)
		assert.Equal(t, id, decoded, "expected %d after decode, got %d", id, decoded)
	}
}

func TestEncodeNegativeID(t *testing.T) {
	_, err := Encode(-1)
	assert.NoError(t, err, "encoding negative id should not error, but may not be meaningful")
}

func TestDecodeInvalidString(t *testing.T) {
	invalidCodes := []string{
		"",
		"abc$def", // invalid character
		"!!!!!!",  // not in alphabet
	}

	for _, code := range invalidCodes {
		_, err := Decode(code)
		assert.Error(t, err, "expected error for invalid code: %s", code)
	}
}
