package utils

func MinValue(n []int) int {
	if len(n) == 0 {
		return 0
	}

	min := n[0]
	for _, value := range n {
		if value < min {
			min = value
		}
	}
	return min
}

func MaxValue(n []int) int {
	if len(n) == 0 {
		return 0
	}

	max := n[0]
	for _, value := range n {
		if value > max {
			max = value
		}
	}
	return max
}
