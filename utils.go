package memaccess

import "bytes"

func toint8(arr []uint16) []uint8 {
	var out []uint8
	for _, v := range arr {
		out = append(out, uint8(v))
	}

	return out
}

func parseint8(arr []uint8) string {
	n := bytes.Index(arr, []uint8{0})

	return string(arr[:n])
}
