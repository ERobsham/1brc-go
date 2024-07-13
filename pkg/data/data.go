package data

import "fmt"

type StationData struct {
	Min   int16
	Max   int16
	Count int32
	Sum   int64
}

func (d StationData) String() string {
	return fmt.Sprintf("%0.1f/%0.1f/%0.1f", (float64(d.Min) / 10.0), (float64(d.Sum) / float64(d.Count) / 10.0), float64(d.Max)/10.0)
}
