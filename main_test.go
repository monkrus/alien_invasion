package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
)


// TestPopulateCities verifies that number of cities and neighbors
// matches the input file
func TestPopulateCities(t *testing.T) {
	data := readCitiesFile()
	defer data.Close()

	cities := populateCities(data)

	data, err := os.Open("cities")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	cityCount := 0
	citiesFromFile := make([]string, 0)

	neighborCount := 0
	neighborsFromFile := make(map[string]int)

	for scanner.Scan() {
		cityCount++

		currentCity := strings.Split(scanner.Text(), " ")
		citiesFromFile = append(citiesFromFile, currentCity[0])

		for _, n := range currentCity[1:] {
			neighborCount++
			neighborsFromFile[n] = 1
		}
	}

	createdNeighborCount := 0

	for _, neighbors := range cities {
		for range neighbors {
			createdNeighborCount++
		}
	}

	// Test if number of cities and neighbors match the file
	if len(cities) != cityCount {
		t.Errorf("Number of created cities %v does not match number in cities file %v", len(cities), cityCount)
	} else if createdNeighborCount != neighborCount {
		t.Errorf("Number of created neighbors %v does not match number in cities file %v", createdNeighborCount, neighborCount)
	}

	// Test if city values match the file
	for _, c := range citiesFromFile {
		if len(cities[c]) == 0 {
			t.Errorf("City %v not found in created map", c)
		}
	}

	// Check if neighbor values match the file
	for _, neighbors := range cities {
		for _, n := range neighbors {
			if neighborsFromFile[n] != 1 {
				t.Errorf("Neighbor %v not found in created map", n)
			}
		}
	}
}

// Test creating aliens and placing them in cities
// Assumption: Max number of aliens can only be # of cities * 2
func TestPopulateAliens(t *testing.T) {
	data := readCitiesFile()
	defer data.Close()

	cities := populateCities(data)

	for i := 0; i < len(cities)*2; i++ {
		aliens := populateAliens(i, cities)

		totalAliens := make([]int, 0)
		for _, v := range aliens {
			for alien := range v {
				totalAliens = append(totalAliens, alien)
			}
		}

		if len(totalAliens) != i {
			t.Errorf("Created aliens %v does not match input %v", len(totalAliens), i)
		}
	}
}

// Function puts 2 aliens in every city, then goes through cities
// and destroys each
// Verifies cities map and alien locations map are empty as a result
func TestDestroyCity(t *testing.T) {
	data := readCitiesFile()
	defer data.Close()

	cities := populateCities(data)

	testState := Invasion{
		cities: cities,
		aliens: make(map[string][]int),
	}

	alienID := 0
	for k := range testState.cities {
		testState.aliens[k] = append(testState.aliens[k], alienID)
		alienID++
		testState.aliens[k] = append(testState.aliens[k], alienID)
		alienID++
	}

	for city := range testState.aliens {
		testState.destroyCity(city)
	}

	if len(testState.aliens) != 0 && len(testState.cities) != 0 {
		t.Errorf("Not all cities destroyed!\nAlien occupied cities: %v\nCities: %v", testState.aliens, testState.cities)
	}
}

// Function tests that the values written to the result file
// reflect the state of the game
func TestWriteMapState(t *testing.T) {
	data := readCitiesFile()
	defer data.Close()

	cities := populateCities(data)
	writeMapState(cities)

	data, err := os.Open("result")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	cityCount := 0
	citiesFromFile := make([]string, 0)

	neighborCount := 0
	neighborsFromFile := make(map[string]int)

	for scanner.Scan() {
		cityCount++

		currentCity := strings.Split(scanner.Text(), " ")
		citiesFromFile = append(citiesFromFile, currentCity[0])

		for _, n := range currentCity[1:] {
			neighborCount++
			neighborsFromFile[n] = 1
		}
	}

	createdNeighborCount := 0

	for _, neighbors := range cities {
		for range neighbors {
			createdNeighborCount++
		}
	}

	// Test if number of cities and neighbors match the file
	if len(cities) != cityCount {
		t.Errorf("Number of created cities %v does not match number in cities file %v", len(cities), cityCount)
	} else if createdNeighborCount != neighborCount {
		t.Errorf("Number of created neighbors %v does not match number in cities file %v", createdNeighborCount, neighborCount)
	}

	// Test if city values match the file
	for _, c := range citiesFromFile {
		if len(cities[c]) == 0 {
			t.Errorf("City %v not found in created map", c)
		}
	}

	// Check if neighbor values match the file
	for _, neighbors := range cities {
		for _, n := range neighbors {
			if neighborsFromFile[n] != 1 {
				t.Errorf("Neighbor %v not found in created map", n)
			}
		}
	}
}
