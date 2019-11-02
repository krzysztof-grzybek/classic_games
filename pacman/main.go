package main
import (
	"fmt"
	"os"
	"os/exec"
	"bufio"
	"log"
)

var maze []string

type Player struct {
	row int
	col int
}

var player Player

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

	for row, line := range maze {
		for col, char := range line {
			switch char {
			case 'P': 
				player = Player{row, col}
			}
		}
	}
	return nil
}

func printScreen() {
	clearScreen()
	for _, line := range maze {
		for _, chr := range line {
			switch chr {
			case '#':
				fmt.Printf("%c", chr)
			default:
				fmt.Printf(" ")
			}

		}
		fmt.Printf("\n")
		
	}

	moveCursor(player.row, player.col)
	fmt.Printf("P")
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
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
					return "UP", nil
			case 'B':
					return "DOWN", nil
			case 'C':
					return "RIGHT", nil
			case 'D':
					return "LEFT", nil
			}
		}
	}

	return "", nil
}

func clearScreen() {
	fmt.Printf("\x1b[2J")
	moveCursor(0, 0)
}

func moveCursor(row, col int) {
	fmt.Printf("\x1b[%d;%df", row + 1, col + 1)
}

func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol

	switch dir {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(maze) - 1
		}	
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(maze) {
			newRow = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(maze[0]) - 1
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(maze[0]) {
			newCol = 0
		}
	}

	if (maze[newRow][newCol] == '#') {
		newRow = oldRow
		newCol = oldCol
	}

	return
}

func movePlayer(dir string) {
	player.row, player.col = makeMove(player.row, player.col, dir)
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

		movePlayer(input)

		if input == "ESC" {
			break
		}
	}
}