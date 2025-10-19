package filters

import (
	"time"

	"github.com/google/uuid"
)

type SubFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	From        *time.Time
	To          *time.Time
}
