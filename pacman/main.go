package main
import (
	"fmt"
	"os"
	"bufio"
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
	for _, line := range maze {
		fmt.Println(line)
	}
}

func main() {

	err := loadMaze()
    if err != nil {
		fmt.Printf("Error loading maze: %v\n", err)
	}
	for {
		printScreen()
	    break;
	}
}