package misc

import (
	"math"
	"testing"
)

func TestGenerateRandomCoordinates(t *testing.T) {

	t.Run("GlobalDefaults", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			lat, lon := GenerateRandomCoordinates(0, 0, 0, 0)
			if lat < -90.0 || lat > 90.0 {
				t.Errorf("Latitude (%f) hors de la plage globale attendue [-90, 90]", lat)
			}
			if lon < -180.0 || lon > 180.0 {
				t.Errorf("Longitude (%f) hors de la plage globale attendue [-180, 180]", lon)
			}
		}
	})

	t.Run("CustomPositiveRange", func(t *testing.T) {
		minLat, maxLat := 10.0, 20.0
		minLon, maxLon := 30.0, 40.0
		for i := 0; i < 100; i++ {
			lat, lon := GenerateRandomCoordinates(minLat, maxLat, minLon, maxLon)
			if lat < minLat || lat > maxLat {
				t.Errorf("Latitude (%f) hors de la plage personnalisée [%f, %f]", lat, minLat, maxLat)
			}
			if lon < minLon || lon > maxLon {
				t.Errorf("Longitude (%f) hors de la plage personnalisée [%f, %f]", lon, minLon, maxLon)
			}
		}
	})

	t.Run("CustomNegativeRange", func(t *testing.T) {
		minLat, maxLat := -20.0, -10.0
		minLon, maxLon := -40.0, -30.0
		for i := 0; i < 100; i++ {
			lat, lon := GenerateRandomCoordinates(minLat, maxLat, minLon, maxLon)
			if lat < minLat || lat > maxLat {
				t.Errorf("Latitude (%f) hors de la plage personnalisée [%f, %f]", lat, minLat, maxLat)
			}
			if lon < minLon || lon > maxLon {
				t.Errorf("Longitude (%f) hors de la plage personnalisée [%f, %f]", lon, minLon, maxLon)
			}
		}
	})

	t.Run("TinyRange", func(t *testing.T) {
		const epsilon = 0.000000001
		fixedLat, fixedLon := 45.12345, 5.67890
		lat, lon := GenerateRandomCoordinates(fixedLat, fixedLat, fixedLon, fixedLon)

		if math.Abs(lat-fixedLat) > epsilon {
			t.Errorf("Latitude (%f) n'est pas suffisamment proche de la valeur fixe attendue (%f)", lat, fixedLat)
		}
		if math.Abs(lon-fixedLon) > epsilon {
			t.Errorf("Longitude (%f) n'est pas suffisamment proche de la valeur fixe attendue (%f)", lon, fixedLon)
		}
	})
}
