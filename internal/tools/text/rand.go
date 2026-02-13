package text

import "math/rand"

const (
	seed = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = seed[rand.Intn(len(seed))]
	}
	return string(b)
}
