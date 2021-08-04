// Copyright Contributors to the Open Cluster Management project

package helpers

import (
	"bytes"
	"fmt"
	"math"
	"time"
)

func TimeDiff(t time.Time, precision time.Duration) string {
	diff := time.Since(t)
	days := diff / (24 * time.Hour)
	hours := diff % (24 * time.Hour)
	minutes := hours % time.Hour
	seconds := math.Mod(minutes.Seconds(), 60)
	var buffer bytes.Buffer
	printed := false
	if days > 0 || printed {
		buffer.WriteString(fmt.Sprintf("%dd", days))
		printed = true
	}
	if precision == time.Hour*24 {
		return buffer.String()
	}
	if hours/time.Hour > 0 || printed {
		buffer.WriteString(fmt.Sprintf("%dh", hours/time.Hour))
		printed = true
	}
	if precision == time.Hour {
		return buffer.String()
	}
	if minutes/time.Minute > 0 || printed {
		buffer.WriteString(fmt.Sprintf("%dm", minutes/time.Minute))
		printed = true
	}
	if precision == time.Minute {
		return buffer.String()
	}
	if seconds > 0 || printed {
		buffer.WriteString(fmt.Sprintf("%.1fs", seconds))
		printed = true
	}
	return buffer.String()
}
