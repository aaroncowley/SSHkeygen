package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	db "./db"
	seed "./jsonSeed"
	"github.com/fatih/color"
)

func menu() {
	for {
		fmt.Println("What would you like to do?")
		color.Set(color.FgCyan)
		fmt.Printf("(1) - Seed Artifacts\n" +
			"(2) - Insert Artifacts into db\n" +
			"(3) - Read Artifacts from db\n" +
			"(4) - Reset Database\n" +
			"(5) - Exit\n")
		color.Unset()

		fmt.Printf(color.GreenString("CLI> "))
		//read one character
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()

		red := color.New(color.FgRed).SprintFunc()

		switch userInput := char; userInput {
		case '1':
			fmt.Println("Seed Arifacts chosen")

			fmt.Println("Start Seeding Now?")
			fmt.Printf("Confirm [%s] or [%s]\n", red("Y"), red("y"))

			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()
			fmt.Println()

			if char == 'Y' || char == 'y' {
				seed.CreateJsonSeed()
			} else {
				continue
			}

		case '2':
			fmt.Println("Insert Artifacts into db chosen")

			fmt.Println("Database will be loaded from jsonSeed/output.json.")
			fmt.Printf("Confirm [%s] or [%s]\n", red("Y"), red("y"))

			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()
			fmt.Println()

			if char == 'Y' || char == 'y' {
				if err := db.InsertKeys(); err != nil {
					log.Fatal(err)
				}
			} else {
				continue
			}

		case '3':
			fmt.Println("Read Artifacts from db chosen")
			fmt.Printf("Confirm [%s] or [%s]\n", red("Y"), red("y"))

			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()
			fmt.Println()

			if char == 'Y' || char == 'y' {
				seed.CreateJsonSeed()
			} else {
				continue
			}

		case '4':
			fmt.Println("Reinit Database chosen\n")

			fmt.Println("Are you sure?")
			fmt.Printf("Confirm [%s] or [%s]\n", red("Y"), red("y"))

			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()
			fmt.Println()

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
