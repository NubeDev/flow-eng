package uuid

import (
	"fmt"
	"time"
)

func TimestampUUID() {
	now := time.Now().Format(time.RFC3339Nano)

	fmt.Println(now)
}
