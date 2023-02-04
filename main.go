package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Plant is a struct that represents a plant in Trefle API
type Plant struct {
	CommonName     string `json:"common_name"`
	Slug           string `json:"slug"`
	ScientificName string `json:"scientific_name"`
	Description    string `json:"description"`
	//Year             string `json:"year"`
	//Bibliography     string `json:"bibliography"`
	Rank             string `json:"rank"`
	FamilyCommonName string `json:"family_common_name"`
	Observation      string `json:"observation"`
	Vegetable        bool   `json:"vegetable"`
	Genus            string `json:"genus"`
	Family           string `json:"family"`
	/*CommonName      struct {
		En  []string `json:"en"`
		Eng []string `json:"eng"`
	} `json:"common_name"`
	/*Distribution struct {
		Native     []string `json:"native"`
		Introduced []string `json:"introduced"`
	} `json:"distribution"`

	/*FruitOrSeed struct {
		Conspicuous     string `json:"conspicuous"`
		Color           string `json:"color"`
		Shape           string `json:"shape"`
		SeedPersistence string `json:"seed_persistence"`
	} `json: "fruit_or_seed"`

	Flower struct {
		Color       string `json:"color"`
		Conspicuous string `json:"conspicuous"`
	} `json: "flower"`*/
}

// PlantsResponse is a struct that represents the response from the Trefle API
type PlantsResponse struct {
	Data []Plant `json:"data"`
}

// Garden map to store plants
var Garden map[string]Plant

func init() {
	Garden = make(map[string]Plant)
}

func main() {
	for {
		fmt.Println()
		fmt.Println("Welcome to the Plant Database")

		fmt.Println()
		fmt.Println("1. Search for a plant")
		fmt.Println("2. Add a plant to the garden")
		fmt.Println("3. Edit a plant in the garden")
		fmt.Println("4. Would you like to plant veggies in the garden? ")
		fmt.Println("5. Exit")
		fmt.Println()

		fmt.Print("Enter your choice: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			searchPlant()
		case "2":
			addPlant()
		case "3":
			editPlant()
		case "4":
			vegetablePlant()
		case "5":
			fmt.Println()
			fmt.Println("Thank you for taking the time to search through 100,000 plants to find the one you were looking for!")
			fmt.Println()

			return
		default:
			fmt.Println("Invalid choice. Try again.")
			fmt.Println()
		}
	}
}

func searchPlant() {
	fmt.Println()
	fmt.Print("Enter plant name: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Build the URL to fetch information from the Trefle API
	url := fmt.Sprintf("https://trefle.io/api/v1/plants/search?q=%s&token=eX8zkmOMqHh9Web_qr5vv917_0l2xkkWXo_oxq4wT4s", input)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data from Trefle API:", err)
		return
	}
	defer response.Body.Close()

	// Check if the status code is in the 200s range
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Printf("Error: Got status code %d from Trefle API\n", response.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var plantsResponse PlantsResponse
	err = json.Unmarshal(body, &plantsResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	// Check if there are any plants in the response
	if len(plantsResponse.Data) == 0 {
		fmt.Println("No plants found for the given name.")
		return
	}

	// Loop through the plants in the response
	for i, plant := range plantsResponse.Data {
		fmt.Printf("\nPlant %d\n", i+1)
		fmt.Println("--------------")
		fmt.Printf("Common Name: %s\n", plant.CommonName)
		fmt.Printf("Scientific Name: %s\n", plant.ScientificName)
		fmt.Printf("Slug: %s\n", plant.Slug)
		fmt.Printf("Rank: %s\n", plant.Rank)
		fmt.Println("Description: ", plant.Description)
		/*if plant.FamilyCommonName != "" {
			fmt.Printf("Family Common Name: %s\n", plant.FamilyCommonName)
		} else {
			fmt.Println("Family Common Name: None/ We couldn't find the information")
		}*/
		fmt.Println("Family common name: ", plant.FamilyCommonName)

		fmt.Printf("Observation: %s\n", plant.Observation)
		fmt.Printf("Is Vegetable: %t\n", plant.Vegetable)
		/*fmt.Printf("Distribution (Native): %s\n", plant.Distribution.Native)
		fmt.Printf("Distribution (Introduced): %s\n", plant.Distribution.Introduced)
		/*	fmt.Printf("Flower Color: %s\n", plant.Flower.Color)
			fmt.Printf("Flower Conspicuous: %s\n", plant.Flower.Conspicuous)
			fmt.Printf("Seed Color: %s\n", plant.FruitOrSeed.Color)
			fmt.Printf("Seed Shape: %s\n", plant.FruitOrSeed.Shape)
			fmt.Printf("Seed Persistence: %s\n", plant.FruitOrSeed.SeedPersistence)
		*/
	}
	fmt.Println("Would you like to add any of these plants to your garden? (yes/no)")
	bufio.NewReader(os.Stdin)
	addPlantChoice, _ := reader.ReadString('\n')
	addPlantChoice = strings.TrimSpace(addPlantChoice)

	if addPlantChoice == "yes" {
		fmt.Println("Enter the common name of the plant you would like to add:")
		reader := bufio.NewReader(os.Stdin)
		selectedPlantName, _ := reader.ReadString('\n')
		selectedPlantName = strings.TrimSpace(selectedPlantName)

		var selectedPlant Plant
		foundPlant := false
		for _, plant := range plantsResponse.Data {
			if plant.CommonName == selectedPlantName {
				selectedPlant = plant
				foundPlant = true
				break
			}
		}

		if foundPlant {
			Garden[selectedPlant.CommonName] = Plant{
				CommonName:  selectedPlant.CommonName,
				Description: selectedPlant.Description,
			}
			fmt.Println("Plant added to the garden successfully!")
		} else {
			fmt.Println("Plant not found.")
		}

	} else {
		fmt.Println("Thank you for taking the time to search for the perfect plant, and We hope you have found the one you were searching for!  ")
		fmt.Println()
	}
}
func addPlant() {
	fmt.Println()
	fmt.Print("Enter plant name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter plant description: ")

	var description strings.Builder
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "." {
			break
		}
		description.WriteString(line + "\n")
	}

	Garden[name] = Plant{
		CommonName:  name,
		Description: description.String(),
	}

	fmt.Println()
	fmt.Println("Plant added to the garden successfully!")
	fmt.Println()

}

func editPlant() {
	fmt.Print("Enter the name of the plant you want to edit: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	plant, found := Garden[name]
	if !found {
		fmt.Println("Plant not found.")
		return
	}

	fmt.Print("Enter new description (press enter after each line, type '.' on a new line to stop):\n")
	var descriptionLines []string
	for {
		reader = bufio.NewReader(os.Stdin)
		descriptionLine, _ := reader.ReadString('\n')
		descriptionLine = strings.TrimSpace(descriptionLine)

		if descriptionLine == "." {
			break
		}

		descriptionLines = append(descriptionLines, descriptionLine)
	}

	plant.Description = strings.Join(descriptionLines, "\n")
	Garden[name] = plant

	fmt.Println()
	fmt.Println("Plant description updated successfully!")
	fmt.Println()
}

func vegetablePlant() {
	// Read the CSV file
	file, err := os.Open("vegetables.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV data:", err)
		return
	}

	for {
		// Ask the user for a vegetable to search for
		fmt.Print("Enter a vegetable to search for (or 'q' to quit): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		searchTerm := scanner.Text()

		if searchTerm == "q" {
			break
		}

		// Search the records for a match
		found := false
		for _, record := range records {
			if strings.ToLower(record[0]) == strings.ToLower(searchTerm) {
				fmt.Println("Name:", record[0])
				fmt.Println("Ideal Temperature for the plant:", record[1])
				fmt.Println("PH:", record[2])
				fmt.Println("Soil: ", record[3])
				fmt.Println("Waterlevel: ", record[4], '\n')

				found = true
				break
			}
		}

		if !found {
			fmt.Println("Sorry, the vegetable you entered was not found in the database.")
		}
	}
}
