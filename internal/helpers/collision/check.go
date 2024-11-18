package collision

import "math"

func Check(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func CheckCircle(x1, y1, r, x2, y2, w2, h2 float64) bool {
	dx := x1 - math.Max(x2, math.Min(x1, x2+w2))
	dy := y1 - math.Max(y2, math.Min(y1, y2+h2))

	return dx*dx+dy*dy < r*r
}
