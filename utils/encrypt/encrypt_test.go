package encrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	tests := []struct {
		origin string
	}{
		{
			origin: "123456",
		},
		{
			origin: "123456..",
		},
	}

	for _, v := range tests {
		encryptedData := BcryptEncrypt(v.origin)
		result := BcryptCheck(v.origin, encryptedData)
		assert.Equal(t, result, true)
	}
}
