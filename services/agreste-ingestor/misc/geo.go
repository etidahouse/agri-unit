package misc

import (
	"math/rand"
	"time"
)

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GenerateRandomCoordinates(minLat, maxLat, minLon, maxLon float64) (latitude, longitude float64) {

	if minLat == 0 && maxLat == 0 {
		minLat = -90.0
		maxLat = 90.0
	}
	if minLon == 0 && maxLon == 0 {
		minLon = -180.0
		maxLon = 180.0
	}

	latitude = minLat + rng.Float64()*(maxLat-minLat)

	longitude = minLon + rng.Float64()*(maxLon-minLon)

	return latitude, longitude
}
