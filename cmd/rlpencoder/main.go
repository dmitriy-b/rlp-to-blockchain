package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/rlp"
)

// Function to convert nested byte slices to hexadecimal strings
func convertToHexString(data interface{}) interface{} {
	switch v := data.(type) {
	case []byte:
		return hex.EncodeToString(v) // Convert byte slice to hex string
	case []interface{}:
		for i, elem := range v {
			v[i] = convertToHexString(elem)
		}
		return v
	default:
		return data
	}
}

// Function to convert hexadecimal strings back to byte slices
func convertToByteSlice(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		bytes, err := hex.DecodeString(v)
		if err != nil {
			log.Fatalf("Error decoding hex string: %v", err)
		}
		return bytes
	case []interface{}:
		for i, elem := range v {
			v[i] = convertToByteSlice(elem)
		}
		return v
	default:
		return data
	}
}

func main() {
	// Define command-line flags
	outputFile := flag.String("output", "output.json", "File to save JSON output")
	loadFile := flag.String("load", "", "JSON file to load and convert to RLP")
	flag.Parse()

	if *loadFile != "" {
		// Load data from JSON file and save to RLP
		loadAndSaveRLP(*loadFile)
	} else {
		// Decode RLP file and save to JSON
		decodeAndSaveJSON(*outputFile)
	}
}

func decodeAndSaveJSON(outputFile string) {
	// Open the RLP file
	file, err := os.Open("chain.rlp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Use a stream to avoid memory overload on large files
	stream := rlp.NewStream(file, 1024*1024) // 1 MB buffer

	// List to hold all decoded blocks
	var allBlocks []interface{}

	// Loop to decode multiple blocks
	for {
		// Define a variable of type interface{} to hold the decoded data
		var decodedData interface{}

		// Decode the next block
		err = stream.Decode(&decodedData)
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			log.Fatalf("Error decoding block: %v", err)
		}

		// Convert decoded data to hexadecimal strings
		decodedData = convertToHexString(decodedData)

		// Append the decoded block to the list
		allBlocks = append(allBlocks, decodedData)
	}

	// Marshal the list of blocks into JSON
	jsonData, err := json.MarshalIndent(allBlocks, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %v", err)
	}

	// Save the JSON data to the specified output file
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Printf("Decoded data saved to %s\n", outputFile)
}

func loadAndSaveRLP(loadFile string) {
	// Read JSON data from file
	jsonData, err := os.ReadFile(loadFile)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Unmarshal JSON data into a slice of interfaces
	var allBlocks []interface{}
	err = json.Unmarshal(jsonData, &allBlocks)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Convert hexadecimal strings back to byte slices
	for i, block := range allBlocks {
		allBlocks[i] = convertToByteSlice(block)
	}

	// Open the RLP file for writing
	file, err := os.Create("output.rlp")
	if err != nil {
		log.Fatalf("Error creating RLP file: %v", err)
	}
	defer file.Close()

	// Encode each block back to RLP
	for _, block := range allBlocks {
		err = rlp.Encode(file, block)
		if err != nil {
			log.Fatalf("Error encoding block to RLP: %v", err)
		}
	}

	fmt.Println("Data loaded from JSON and saved to output.rlp")
}
