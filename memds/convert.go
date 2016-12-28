package memds

func Uint8ArrayToString(a []uint8) string {
	b := make([]byte, 0, len(a))
	for _, e := range a {
		b = append(b, byte(e))
	}
	return string(b)
}
