package sh

import "unsafe"

func bytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
