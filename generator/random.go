package generator

import "math/rand"

func GetRandomKey() string {
	AlphaBet := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Shuffle(len(AlphaBet), func(i, j int) {
		AlphaBet[i], AlphaBet[j] = AlphaBet[j], AlphaBet[i]
	})
	id := string(AlphaBet[:5])
	return id
}
