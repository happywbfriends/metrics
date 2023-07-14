package metrics

const (
	second = 1_000.0
	minute = 60_000.0
)

var DefaultDurationMsBuckets = []float64{
	1.0, 5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0,
	1 * second, 3 * second, 5 * second, 10 * second, 30 * second,
	1 * minute, 5 * minute, 10 * minute}
