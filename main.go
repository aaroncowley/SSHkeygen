package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	db "./db"
	seed "./jsonSeed"
)

func menu() {
	for {
		fmt.Println("What would you like to do?")

		fmt.Println(`
		(1) - Seed Artifacts
		(2) - Insert Artifacts into db 
		(3) - Read Artifact from db
		(4) - Reinit Database
		(5) - Exit`)

		//read one character
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()

		switch userInput := char; userInput {
		case '1':
			fmt.Println("Seed Arifacts chosen")

			fmt.Println("Start Seeding Now?")
			fmt.Println("Confirm [Y] or [N]")
			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()

			if char == 'Y' || char == 'y' {
				seed.CreateJsonSeed()
			} else {
				continue
			}

		case '2':
			fmt.Println("Insert Artifacts into db chosen")

			fmt.Println("Database will be loaded from jsonSeed/output.json.")
			fmt.Println("Confirm [Y] or [N]")
			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()

			if char == 'Y' || char == 'y' {
				if err := db.InsertKeys(); err != nil {
					log.Fatal(err)
				}
			} else {
				continue
			}

		case '3':
			fmt.Println("Read Artifact from db chosen")
			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()

			if char == 'Y' || char == 'y' {
				seed.CreateJsonSeed()
			} else {
				continue
			}

		case '4':
			fmt.Println("Reinit Database chosen")

			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Confirm [Y] or [N]")
			char, _, _ := reader.ReadRune()

			if char == 'Y' || char == 'y' {
				if err := db.InitDB(); err != nil {
					log.Fatal(err)
				}
			} else {
				continue
			}

		case '5':
			os.Exit(0)

		default:
			fmt.Println("Invalid Input")
		}
	}
}

func main() {
	menu()
}
