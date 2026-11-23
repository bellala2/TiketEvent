package clock

import "time"

// Clock membungkus "waktu sekarang" supaya bisa diganti saat test
// (misalnya untuk menguji token kedaluwarsa tanpa menunggu 1 jam).
type Clock interface {
	Now() time.Time
}

// System adalah Clock asli yang memakai jam komputer, selalu dalam UTC.
type System struct{}

func (System) Now() time.Time {
	return time.Now().UTC()
}
