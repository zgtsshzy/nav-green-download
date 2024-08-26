package isoline

func Max(nums []float32) (int, float32) {
	max := nums[0]
	idx := 0
	for i := range nums {
		if nums[i] > max {
			idx = i
			max = nums[i]
		}
	}
	return idx, max
}

func Min(nums []float32) (int, float32) {
	min := nums[0]
	idx := 0
	for i := range nums {
		if nums[i] < min {
			idx = i
			min = nums[i]
		}
	}
	return idx, min
}

func Max64(nums []float64) (int, float64) {
	max := nums[0]
	idx := 0
	for i := range nums {
		if nums[i] > max {
			idx = i
			max = nums[i]
		}
	}
	return idx, max
}

func Min64(nums []float64) (int, float64) {
	min := nums[0]
	idx := 0
	for i := range nums {
		if nums[i] < min {
			idx = i
			min = nums[i]
		}
	}
	return idx, min
}

func Max_2d(nums [][]float32) (int, int, float32) {
	max := nums[0][0]
	idx_i, idx_j := 0, 0
	for i := range nums {
		idx, ma := Max(nums[i])
		if ma > max {
			idx_i = i
			idx_j = idx
			max = ma
		}
	}
	return idx_i, idx_j, max
}

func Min_2d(nums [][]float32) (int, int, float32) {
	min := nums[0][0]
	idx_i, idx_j := 0, 0
	for i := range nums {
		idx, mi := Min(nums[i])
		if mi < min {
			idx_i = i
			idx_j = idx
			min = mi
		}
	}
	return idx_i, idx_j, min
}
