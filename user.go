package wrap

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// UserCan func to check user perms
type UserCan func(Userer) bool

// Userer is an interface for authed user which is stored in ctx
type Userer interface {
	fmt.Stringer
	ID() uuid.UUID
	Can(UserCan) bool
}
