package main
import (
	"fmt"
	"os"
	"os/exec"
	"bufio"
	"log"
)

var maze []string

func loadMaze() error {
	file, err := os.Open("maze01.txt");	
	if err != nil {
		return err
	}

	defer file.Close();

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}

	return nil
}

func printScreen() {
	clearScreen()
	for _, line := range maze {
		fmt.Println(line)
	}
}

func init() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	err := cbTerm.Run()
	if err != nil {
	    log.Fatalf("Cannot enable cbreak mode in terminal: %v\n", err)
	}
}

func cleanup() {
	cookedTerm := exec.Command("stty", "-cbreak", "echo")
	cookedTerm.Stdin = os.Stdin

	err := cookedTerm.Run()
	if err != nil {
		log.Fatalf("Unable to enable cooked terminal: %v\n", err)
	}
}

func readinput() (string, error) {
	buffer := make([]byte, 100)
	cnt, err := os.Stdin.Read(buffer)

	if err != nil {
		return "", err
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	}

	return "", nil
}

func clearScreen() {
	fmt.Printf("\x1b[2J")
	moveCursor(0, 0)
}

func moveCursor(row, col int) {
	fmt.Printf("\x1b[%d;%df", row, col)
}

func main() {
	defer cleanup()

	err := loadMaze()
    if err != nil {
		fmt.Printf("Error loading maze: %v\n", err)
	}
	for {
		printScreen()

		input, err := readinput()
		if err != nil {
			log.Fatalf("Error reading input: %v\n", err)
			break
		}

		if input == "ESC" {
			break
		}
	}
}