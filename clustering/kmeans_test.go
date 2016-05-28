package clustering

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestKmeans(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	m := [][]float64{
		[]float64{0.1, 0.0},
		[]float64{0.9, 1.0},
		[]float64{-0.1, 0.0},
		[]float64{0.0, -0.1},
		[]float64{1.1, 1.0},
		[]float64{1.0, 1.1},
		[]float64{1.0, 0.9},
		[]float64{0.0, 0.1},
	}

	means, tags := Kmeans(m, 2)

	if tags[0] == 0 {
		assertEqual(tags, []int{0, 1, 0, 0, 1, 1, 1, 0}, t)
		assertEqual(means, [][]float64{{0, 0}, {1, 1}}, t)
	} else {
		assertEqual(tags, []int{1, 0, 1, 1, 0, 0, 0, 1}, t)
		assertEqual(means, [][]float64{{1, 1}, {0, 0}}, t)
	}
}

func assertEqual(act, exp interface{}, t *testing.T) {
	if !reflect.DeepEqual(act, exp) {
		t.Fatalf("Wrong value: %v, expected %v", act, exp)
	}
}
