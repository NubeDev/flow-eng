package helpers

import (
	"crypto/rand"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func PingMessage(ip string, status bool, lastOk time.Time) string {
	var msg string
	switch status {
	case true:
		msg = fmt.Sprintf("✓ %s OK since %s ", ip, ttime.TimeSince(lastOk))
	case false:
		msg = fmt.Sprintf("⚠ %s ERROR since %s", ip, ttime.TimeSince(lastOk))
	}
	return msg

}

// CommandPing a host with a single packet (Linux and Windows).
// Returns false if ping exits with error code or with "Destination host unreachable".
// Returns true in case of successful ping.
func CommandPing(host string) bool {
	cmd := exec.Command("ping", "-c", "1", host)
	var status bool
	if IsWindows() {
		cmd = exec.Command("ping", "-n", "1", host)
	}
	out, err := cmd.Output()
	if err != nil || strings.Contains(string(out), "destination host unreachable") {
		status = false
	} else {
		status = true
	}
	return status
}

func IsWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	} else {
		return false
	}
}

func RandomName(prefix ...string) string {
	u := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, u)
	if n != len(u) || err != nil {
		return "-error-uuid-"
	}
	uuid := fmt.Sprintf("%x%x", u[2:4], u[4:4])
	if len(prefix) > 0 {
		uuid = fmt.Sprintf("%s_%s", prefix[0], uuid)
	}
	return uuid
}

func ShortUUID(prefix ...string) string {
	u := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, u)
	if n != len(u) || err != nil {
		return "-error-uuid-"
	}
	uuid := fmt.Sprintf("%x%x", u[0:4], u[4:4])
	if len(prefix) > 0 {
		uuid = fmt.Sprintf("%s_%s", prefix[0], uuid)
	}
	return uuid
}

func UUID(prefix ...string) string {
	u := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, u)
	if n != len(u) || err != nil {
		return "-error-uuid-"
	}
	uuid := fmt.Sprintf("%x%x", u[0:4], u[4:6])
	if len(prefix) > 0 {
		uuid = fmt.Sprintf("%s_%s", prefix[0], uuid)
	}
	return uuid
}
