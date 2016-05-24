package clustering

// Performs agglomerative clustering on the indexes 0 to n-1. d should return
// the distance between the i'th and j'th element, such that d(i,j)=d(j,i) and
// d(i,i)=0. d will be called exactly n choose 2 times.
func Agglo(n int, d func(int, int) float64) {
	// Calculate pairwise distances.
	ds := make([][]float64, n)
	for i := range ds {
		ds[i] = make([]float64, n)
	}
	for i := range ds {
		for j := i+1; j < n; j++ {
			ds[i][j] = d(i, j)
			ds[j][i] = ds[i][j]
		}
	}

	setd := ds

	
	sets := make([][]int, len(words))
	for i := range sets {
		sets[i] = []int{i}
	}

	for {
		// Find minimal distance.
		min, mini, minj := 2.0, -1, -1
		for i := range sets {
			for j := range sets[i+1:] {
				j += i+1
				// Check if better than min.
				if setd[i][j] < min {
					min, mini, minj = setd[i][j], i, j
				}
			}
		}

		// Break if min-distance is big enough.
		if min > distCutoff {
			break
		}

		// Move distances of new set in position j.
		for i := range setd {
			li, lj := float64(len(sets[mini])), float64(len(sets[minj]))
			setd[i][mini] = (setd[i][mini] * li + setd[i][minj] * lj) / (li + lj)
			setd[mini][i] = setd[i][mini]
		}
		for i := range setd {
			setd[i][minj] = setd[i][len(setd[i]) - 1]
			setd[i] = setd[i][:len(setd[i]) - 1]
		}
		setd[minj] = setd[len(setd) - 1]
		setd = setd[:len(setd) - 1]

		// Unite i'th and j'th sets.
		sets[mini] = append(sets[mini], sets[minj]...)
		sets[minj] = sets[len(sets) - 1]
		sets = sets[:len(sets) - 1]
	}
}

