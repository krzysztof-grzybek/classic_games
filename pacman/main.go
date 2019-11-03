package main
import (
	"fmt"
	"os"
	"os/exec"
	"bufio"
	"log"
	"math/rand"
	"time"
	"encoding/json"
	"flag"
)
var (
    configFile = flag.String("config-file", "config.json", "path to custom configuration file")
    mazeFile   = flag.String("maze-file", "maze01.txt", "path to a custom maze file")
)
var maze []string

type Player struct {
	row int
	col int
}

type Ghost struct {
	row int
	col int
}

type Config struct {
    Player   string `json:"player"`
    Ghost    string `json:"ghost"`
    Wall     string `json:"wall"`
    Dot      string `json:"dot"`
    Pill     string `json:"pill"`
    Death    string `json:"death"`
    Space    string `json:"space"`
    UseEmoji bool   `json:"use_emoji"`
}

var cfg Config
var player Player
var ghosts []*Ghost
var score int
var numDots int
var lives = 1

func loadConfig() error {
	f, err := os.Open(*configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	decodingErr := decoder.Decode(&cfg)

	if decodingErr != nil { 
		return decodingErr
	}

	return nil
}

func loadMaze() error {
	file, err := os.Open(*mazeFile);	
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
			case 'G':
				ghosts = append(ghosts, &Ghost{row, col})
			case '.':
				numDots++
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
				fmt.Printf(cfg.Wall)
			case '.':
				fmt.Printf(cfg.Dot)
			default:
				fmt.Printf(cfg.Space)
			}

		}
		fmt.Printf("\n")
		
	}

	moveCursor(player.row, player.col)
	fmt.Printf(cfg.Player)

	for _, ghost := range ghosts {
		moveCursor(ghost.row, ghost.col);
		fmt.Printf(cfg.Ghost)
	}

	moveCursor(len(maze) + 1, 0)
	fmt.Printf("Score: %v, Lives: %v", score, lives)
}

func initialize() {
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
	if cfg.UseEmoji {
		fmt.Printf("\x1b[%d;%df", row, col * 2 + 1)
	} else {
		fmt.Printf("\x1b[%d;%df", row, col + 1)
	}
}

func drawDirection() string {
	dir := rand.Intn(4)
	dirMap := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return dirMap[dir]
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
	switch maze[player.row][player.col] {
	case '.':
		numDots--
		score++
		maze[player.row] = maze[player.row][0:player.col] + " " + maze[player.row][player.col + 1:]
	}
}

func moveGhosts() {
	for _, ghost := range ghosts {
		dir := drawDirection()
		ghost.row, ghost.col = makeMove(ghost.row, ghost.col, dir)
	}
}

func main() {
	flag.Parse()

	initialize()
	defer cleanup()

	err := loadMaze()
    if err != nil {
		fmt.Printf("Error loading maze: %v\n", err)
	}

	cfgErr := loadConfig()
	if cfgErr != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
	}

	inputCh := make(chan string)

	go func(ch chan<- string) {
		for {
			input, err := readinput()
			if err != nil {
				ch <- "ESC"
			}
			ch <-input
		}
	}(inputCh)

	for {
		printScreen()




		select {
		case inp := <-inputCh:
			if inp == "ESC" {
				lives = 0
			}
			movePlayer(inp)
		default:
		}

		moveGhosts()

		for _, ghost := range ghosts {
			if ghost.row == player.row && ghost.col == player.col {
				lives = 0
			}
		}

		if lives == 0 || numDots == 0 {
			if lives == 0 {
				moveCursor(player.row, player.col)
				fmt.Printf(cfg.Death)
				moveCursor(len(maze)+2, 0)
			}
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
}