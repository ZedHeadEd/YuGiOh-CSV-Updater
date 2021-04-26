/*

1. Create a new .csv file for the output to go to.
2. Create the headers in the new file.
3. Open the default .csv collection file.
4. Grab the card ID and card quantity.
5. Send the response of the requesting card id to the api.
6. Receive and store the response from the api of the requested card.
7. Write the response and the card quanity to the new file.
8. Wait 1 second to not be IP blocked by the API for an hour.
9. Jump back to Step 4 for the next card until all the card data has been collected.

*/

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// A Response struct to map the incoming response to
type Response struct {
	IncomingResponse []CardData `json:"data"`
}

// A CardData Struct to map every card to.
type CardData struct {
	CardId      int          `json:"id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Description string       `json:"desc"`
	Attack      int          `json:"atk"`
	Defence     int          `json:"def"`
	Level       int          `json:"level"`
	Race        string       `json:"race"`
	Attribute   string       `json:"attribute"`
	Sets        []CardSets   `json:"card_sets"`
	Images      []CardImages `json:"card_images"`
	Prices      []CardPrices `json:"card_prices"`
}

// A CardSets struct to hold info specific to each card set.
type CardSets struct {
	SetName       string  `json:"set_name"`
	SetCode       string  `json:"set_code"`
	SetRarity     string  `json:"set_rarity"`
	SetRarityCode string  `json:"set_rarity_code"`
	SetPrice      float32 `json:"set_price"`
}

// A CardImages struct to hold info specific to the images of the card.
type CardImages struct {
	ImageID       string `json:"id"`
	ImageURL      string `json:"image_url"`
	ImageURLSmall string `json:"image_url_small"`
}

// A CardPrices struct to hold info specific to the IRL prices of the card.
type CardPrices struct {
	PriceCardmarket   string `json:"cardmarket_price"`
	PriceTGCPlayer    string `json:"tcgplayer_price"`
	PriceEbay         string `json:"ebay_price"`
	PriceAmazon       string `json:"amazon_price"`
	PriceCoolstuffinc string `json:"coolstuffinc_price"`
}

// An OGcsv struct to hold info of the card data from the original csv card.
type OGcsv struct {
	cardname      string
	cardq         string
	cardid        string
	cardrarity    string
	cardcondition string
	card_edition  string
	cardset       string
	cardcode      string
}

func main() {

	//URL Default for API.
	ApiURL := "https://db.ygoprodeck.com/api/v7/cardinfo.php?id="

	//Headers for the new file in order.
	NewcsvPrep := []string{"CardID", "Name", "Quanity", "Type", "Level", "Attack", "Defence",
		"Race & Spell/Trap Type", "Attribute", "Card Rarity", "Card Set", "Description"}

	//Read original CSV File.
	OGcsvRecords, err := readData("Card Collection.csv")
	if err != nil {
		fmt.Println("Unable to find 'Card Collection.csv' in this directory...")
		time.Sleep(3000 * time.Millisecond)
		log.Fatal(err)
	}
	fmt.Println("Successfully Opened/Read Card Collection CSV file")

	//Create the new CSV file.
	outFile, err := os.Create("New Collection.csv")
	if err != nil {

		log.Fatalln("failed to open file", err)
	}
	defer outFile.Close()

	//Prepare the new CSV file.
	outWriter := csv.NewWriter(outFile)
	err = outWriter.Write(NewcsvPrep)
	outWriter.Flush()
	if err != nil {
		log.Fatal(err)
	}

	//Vars used for the upcoming loop for user output
	var loopCounter = 0
	var resposneNumber = 0
	var recordLength = len(OGcsvRecords)

	///Loop through each line of the old records, grab the api response for the line and write it to a new file.
	for _, record := range OGcsvRecords {
		CardInfo := OGcsv{
			cardname:      record[0],
			cardq:         record[1],
			cardid:        record[2],
			cardrarity:    record[3],
			cardcondition: record[4],
			card_edition:  record[5],
			cardset:       record[6],
			cardcode:      record[7],
		}

		//Get the API response for the card.
		HttpResponse, err := http.Get(ApiURL + CardInfo.cardid)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		//Read the body of the response.
		responseData, err := ioutil.ReadAll(HttpResponse.Body)
		if err != nil {
			log.Fatal(err)
		}

		//Format the response.
		var responseObject Response
		json.Unmarshal(responseData, &responseObject)

		//Output to the terminal which card you are currently working on.
		resposneNumber = resposneNumber + 1
		fmt.Println("Grabbed card data " + fmt.Sprint(resposneNumber) + " of " + fmt.Sprint(recordLength))

		//Write the line to the new .csv file.
		NewLineToWrite := []string{fmt.Sprint(responseObject.IncomingResponse[0].CardId), responseObject.IncomingResponse[0].Name,
			CardInfo.cardq, responseObject.IncomingResponse[0].Type, fmt.Sprint(responseObject.IncomingResponse[0].Level),
			fmt.Sprint(responseObject.IncomingResponse[0].Attack), fmt.Sprint(responseObject.IncomingResponse[0].Defence),
			responseObject.IncomingResponse[0].Race, responseObject.IncomingResponse[0].Attribute, CardInfo.cardrarity,
			CardInfo.cardset, responseObject.IncomingResponse[0].Description}
		outWriter.Write(NewLineToWrite)
		outWriter.Flush()

		//Yup
		loopCounter = loopCounter + 1

		//Stall for 0.5 seconds to not be locked out of the API for a hour.
		if loopCounter > 19 {
			loopCounter = 0
			time.Sleep(500 * time.Millisecond)
		}

	}
}

//Reads the data from the old .csv file
func readData(fileName string) ([][]string, error) {

	csvFile, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	//Skips the first line
	if _, err := reader.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := reader.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
