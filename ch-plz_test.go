package chplz

import (
	"testing"
)

func TestFind(t *testing.T) {
	want := Zip(5000)
	city, ok := FindFirst("Aarau")
	if !ok || city.Zip != want {
		t.Errorf("got %v, wanted %d", city, want)
	}
}

func TestDoubleZip(t *testing.T) {
	want := 2
	cities := FindAll("5405")
	if len(cities) != want {
		t.Errorf("got %v, wanted 2 entries", cities)
	}

}

func TestMultipleCities(t *testing.T) {
	want := 6
	cities := FindAll("Baden")
	if len(cities) != want {
		t.Errorf("got %d, %v, wanted %d entries", len(cities), cities, want)
	}
}

func TestRegex(t *testing.T) {
	want := 3
	cities, err := FindAllRegex(`Oberdorf.*`)
	if err != nil || len(cities) != want {
		t.Errorf("got error %v, %d, %v, wanted %d", err, len(cities), cities, want)
	}
}
