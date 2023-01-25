package numhelper

func Digits(n int) int {
	if n < 0 {
		n *= -1
	}
	if n == 0 {
		return 1
	}
	count := 0
	for n > 0 {
		n = n / 10
		count++
	}
	return count
}
