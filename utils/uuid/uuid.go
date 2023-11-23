package uuid

import (
	"github.com/rs/xid"
)

// NewUuid Create uuid(20chars)
func NewUuid() string {
	guid := xid.New()
	return guid.String()
}
