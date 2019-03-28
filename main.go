package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

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

			fmt.Println("Enter the Amount of Keys to be generated")
			fmt.Printf("(Enter a value between 1 and 1000000) > ")

			color.Set(color.FgRed)
			reader := bufio.NewReader(os.Stdin)
			keyString, _ := reader.ReadString('\n')
			color.Unset()

			keyString = keyString[:len(keyString)-1]

			keyNum, err := strconv.Atoi(keyString)
			if err != nil {
				log.Println(err)
				log.Println("non numneric value entered, please enter a valid number")
				continue
			}

			if keyNum <= 0 {
				log.Println("Please enter a value greater than 0.")
				continue
			}

			fmt.Println("Start Seeding Now?")
			fmt.Printf("Confirm [%s] or [%s]\n", red("Y"), red("y"))

			reader = bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()
			fmt.Println()

			if char == 'Y' || char == 'y' {
				seed.CreateJsonSeed(keyNum)
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
				// TODO: This
				fmt.Println("still in dev")
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
