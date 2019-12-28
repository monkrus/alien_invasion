package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// invasion simulation
type Invasion struct {
	cities map[string][]string
	aliens map[string][]int
}

// starts new invasion, populates cities, populates aliens, writes map state to a file
func main() {
	data := readCitiesFile()
	defer data.Close()

	cities := populateCities(data)

	// Get number of aliens from program arg
	numAliens, err := strconv.Atoi(os.Args[1:2][0])
	if err != nil {
		log.Fatal(err)
	}

	aliens := populateAliens(numAliens, cities)

	state := Invasion{
		cities: cities,
		aliens: aliens,
	}

	state.runSimulation()

	writeMapState(state.cities)
}

// function writes the end state of the invasion simulation to a file
// in the same format as the input file
// Assumptions: can't specify result filename
func writeMapState(cities map[string][]string) {
	f, err := os.Create("result")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for city := range cities {
		var buffer bytes.Buffer
		buffer.WriteString(city + " ")

		for i, c := range cities[city] {
			if i == len(cities[city])-1 {
				buffer.WriteString(c)
			} else {
				buffer.WriteString(c + " ")
			}
		}

		buffer.WriteString("\n")

		_, err := w.WriteString(buffer.String())
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()
}

// function reads a file named "cities" and returns a File struct
// with the contents
func readCitiesFile() *os.File {
	data, err := os.Open("cities")
	if err != nil {
		log.Fatal(err)
	}

	return data
}

// function reads data from a file called "cities"
// and populates the cities map
// Assumptions:
// 1. File is NOT specified by program argument - spec only mentions # of aliens
// 2. File will always be correctly formatted, e.g. no cities w/ spaces
func populateCities(data *os.File) map[string][]string {
	populated := make(map[string][]string)
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	// Iterate over lines of cities file
	for scanner.Scan() {
		// Split city line by spaces
		currentCity := strings.Split(scanner.Text(), " ")

		// If this city does not already exist
		if len(populated[currentCity[0]]) == 0 {
			// Create new city with name of current line
			populated[currentCity[0]] = make([]string, 0)

			// Add to uniqueCities for use in alien assignments
			// invasion.uniqueCities = append(invasion.uniqueCities, currentCity[0])
		}

		// Add neighbors of current city
		for _, neighbor := range currentCity[1:] {
			// Skip neighbor if it is the same as current city
			if strings.Contains(neighbor, currentCity[0]) {
				continue
			}

			populated[currentCity[0]] = append(populated[currentCity[0]], neighbor)
		}
	}

	// Throw error and exit if no cities were created
	if len(populated) == 0 {
		log.Fatalf("populateAliens: must populate cities first")
	}

	return populated
}

// function creates numAliens and randomly places them in a city
// Assumption: no more than 2 aliens may begin in the same city
func populateAliens(numAliens int, cities map[string][]string) map[string][]int {
	if numAliens > len(cities)*2 {
		log.Fatalf("Number of aliens cannot be > 2x number of cities!")
	}

	aliens := make(map[string][]int)
	rand.Seed(time.Now().Unix())

	uniqueCities := make([]string, 0)

	for city := range cities {
		uniqueCities = append(uniqueCities, city)
	}

	for i := 0; i < numAliens; i++ {
		// Pick random city
		city := uniqueCities[rand.Intn(len(uniqueCities))]

		// Ensure no cities have more than 2 aliens
		for len(aliens[city]) == 2 {
			city = uniqueCities[rand.Intn(len(uniqueCities))]
		}

		// Append alien to city
		aliens[city] = append(aliens[city], i)
	}

	return aliens
}

// function deletes a city, specified by a string, from the cities map
// and alien locations map
// When destroying a city, visits its own neighbors and removes itself from
// their list of neighbors
func (invasion *Invasion) destroyCity(city string) {
	// Regex to filter out direction prefix
	reg, err := regexp.Compile(`^(.*?)=`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v has been destroyed by alien %v and alien %v!\n", city, invasion.aliens[city][0], invasion.aliens[city][1])

	// Go to neighbors of deleted and delete itself from their lists
	for _, neighbor := range invasion.cities[city] {
		neighbor = reg.ReplaceAllString(neighbor, "")
		for i, n := range invasion.cities[neighbor] {
			// Neighbors contain the direction prefix, so check that city
			// is a substring of the neighbor instead of direct comparison
			if strings.Contains(n, city) {
				invasion.cities[neighbor][i] = invasion.cities[neighbor][len(invasion.cities[neighbor])-1]
				invasion.cities[neighbor] = invasion.cities[neighbor][:len(invasion.cities[neighbor])-1]
			}
		}
	}

	delete(invasion.aliens, city)
	delete(invasion.cities, city)
}

// function runs the invasion simulation until one of two conditions is
// met: all aliens have died, or 10k turns have passed
func (invasion *Invasion) runSimulation() {
	turns := 0
	for len(invasion.aliens) != 0 && turns < 10000 {
		// First pass: destroy cities with 2 aliens in them
		// This prevents case of having 3 aliens in a city
		for city := range invasion.aliens {
			if len(invasion.aliens[city]) >= 2 {
				invasion.destroyCity(city)
			}
		}

		// Second pass: handle moving aliens
		for city := range invasion.aliens {
			// If this city has no neighbors to move to this alien
			// does nothing
			if len(invasion.cities[city]) == 0 {
				continue
			} else {
				rand.Seed(time.Now().Unix())

				// Regex to filter out direction prefix so we can properly move
				// This is because format of cities map is a non-prefixed
				// cites mapped to an array of direction-prefixed cities
				reg, err := regexp.Compile(`^(.*?)=`)
				if err != nil {
					log.Fatal(err)
				}

				// Get random neighbor of current city for alien to move to
				newCity := invasion.cities[city][rand.Intn(len(invasion.cities[city]))]
				filtered := reg.ReplaceAllString(newCity, "")

				// Move this alien
				invasion.aliens[filtered] = append(invasion.aliens[filtered], invasion.aliens[city][0])
				invasion.aliens[city] = invasion.aliens[city][1:]

				// Remove city from aliens if it was the only one there
				if len(invasion.aliens[city]) == 0 {
					delete(invasion.aliens, city)
				}
			}
		}

		turns++
	}
}
