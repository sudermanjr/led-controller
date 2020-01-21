package dashboard

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getUUID(t *testing.T) {
	var validUUID = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	match := validUUID.MatchString(getUUID())
	assert.Equal(t, match, true)
}
