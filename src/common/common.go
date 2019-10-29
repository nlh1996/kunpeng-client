package common

// ComputeDistance .
func ComputeDistance(x1 int, y1 int,x2 int, y2 int) int {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}