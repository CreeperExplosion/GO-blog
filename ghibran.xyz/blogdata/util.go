package blogdata

func Reverse(n []Comment) []Comment {
	for i := 0; i < len(n)/2; i++ {
		j := len(n) - i - 1
		n[i], n[j] = n[j], n[i]
	}
	return n
}
