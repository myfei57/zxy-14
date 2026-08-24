package mask

import (
	"time"

	"github.com/google/uuid"
)

func randomID() string {
	return uuid.NewString()
}

type timeType = time.Time

func timeNow() time.Time {
	return time.Now().UTC()
}
