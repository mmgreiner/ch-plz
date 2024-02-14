package chplz

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//go:embed zipcodes.ch.csv
var csvfile []byte

type Zip int
type City struct {
	Zip          Zip
	Name         string
	KantonCode   string
	Kanton       string
	Bezirk       string
	BezirkCode   int
	Gemeinde     string
	GemeindeCode int
	Latitute     string
	Longitute    string
}

func (c City) FullName() string {
	return fmt.Sprintf("%d %s %s", c.Zip, c.Name, c.Kanton)
}

var (
	plzMap  = map[Zip]City{}
	nameMap = map[string]City{}
	cities  = []City{}
)

func atoi(s string) int {
	if i, err := strconv.Atoi(s); err != nil {
		panic(fmt.Errorf("cannot convert %s to integer", s))
	} else {
		return i
	}
}

func init() {

	reader := csv.NewReader(bytes.NewBuffer(csvfile))
	data, _ := reader.ReadAll()

	// skip header
	for _, row := range data[1:] {

		// sequence of columns:
		// country_code,zipcode,place,state,state_code,province,province_code,community,community_code,latitude,longitude

		city := City{
			Zip:          Zip(atoi(row[1])),
			Name:         row[2],
			Kanton:       row[3],
			KantonCode:   row[4],
			Bezirk:       row[5],
			BezirkCode:   atoi(row[6]),
			Gemeinde:     row[7],
			GemeindeCode: atoi(row[8]),
			Latitute:     row[9],
			Longitute:    row[10],
		}

		cities = append(cities, city)
		plzMap[city.Zip] = city

		// what if two entries have same city name?
		c2, ok := nameMap[city.Name]
		if ok && city.Kanton == c2.Kanton {
			// only keep the name with the lowest PLZ (8000 for Zürich)
			if city.Zip < c2.Zip {
				nameMap[city.Name] = city
			}
		} else {
			nameMap[city.Name] = city
			nameMap[strings.ToUpper(city.Name)] = city
		}
	}
}

// FindFirst returns the first city with the given name or zip code
func FindFirst(plzOrName string) (City, bool) {
	plz, err := strconv.Atoi(plzOrName)
	if err == nil {
		// is plz
		city, ok := plzMap[Zip(plz)]
		return city, ok
	}

	// now it is a name
	city, ok := nameMap[plzOrName]
	if ok {
		return city, ok
	}

	// almost last attempt with all caps
	city, ok = nameMap[strings.ToUpper(plzOrName)]
	return city, ok
}

// FindAll returns all PLZ where the city name matches
// plzOrName can also be a regular expression
func FindAll(plzOrName string) []City {
	plz, err := strconv.Atoi(plzOrName)
	if err == nil {
		return findAll(func(c City) bool {
			return c.Zip == Zip(plz)
		})
	}
	return findAll(func(c City) bool {
		return strings.Compare(c.Name, plzOrName) == 0
	})
}

// FindAllRegex finds all cities where the city name matches the given expression.
func FindAllRegex(ex string) ([]City, error) {
	pat, err := regexp.Compile(ex)
	if err != nil {
		return []City{}, err
	}
	return findAll(func(c City) bool {
		return pat.Match([]byte(c.Name))
	}), nil
}

func findAll(isequal func(city City) bool) []City {
	results := []City{}
	for _, c := range cities {
		if isequal(c) {
			results = append(results, c)
		}
	}
	return results
}
