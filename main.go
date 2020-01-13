package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	"github.com/eiannone/keyboard"
)

const TILE_MAP_BG = ' '
const TILE_SNAKE_BODY = '*'
const TILE_HORIZ_BORDER = '-'
const TILE_VERT_BORDER = '|'

type coord struct {
	x,y int
}

type snake_t struct {
	direction direction_t
	vel_x     int
	vel_y     int
	body      []coord
}

type direction_t int
const (
	dirRight direction_t = 1
	dirLeft  direction_t = 2
	dirUp    direction_t = 3
	dirDown  direction_t = 4
)

const WIDTH = 40
const HEIGHT = 20

func __cleanScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func mapFlush(gameMap *[WIDTH][HEIGHT]rune) {
	for x, val := range gameMap {
		for y := range val {
			gameMap[x][y] = TILE_MAP_BG
		}
	}
}

func mapPrint(gameMap [WIDTH][HEIGHT]rune) {
	for i := 0; i < WIDTH + 2; i++ {
		fmt.Printf("%c", TILE_HORIZ_BORDER)
	}
	fmt.Println()

	for y := 0; y < HEIGHT; y++ {
		fmt.Printf("%c", TILE_VERT_BORDER)
		for x := 0; x < WIDTH; x++ {
			fmt.Printf("%c", gameMap[x][y])
		}
		fmt.Printf("%c\n", TILE_VERT_BORDER)
	}

	for i := 0; i < WIDTH + 2; i++ {
		fmt.Printf("%c", TILE_HORIZ_BORDER)
	}
	fmt.Println()
}

func snakeInit(snake *snake_t) {
	snake.vel_x = 1
	snake.vel_y = 0

	snake.body = make([]coord, 3)

	snake.body[0].x = 2
	snake.body[0].y = 3
	snake.body[1].x = 3
	snake.body[1].y = 3
	snake.body[2].x = 4
	snake.body[2].y = 3
}

func snakeDraw(snake snake_t, gameMap *[WIDTH][HEIGHT]rune) {
	for _, val := range snake.body {
		fmt.Println(val)
		gameMap[val.x][val.y] = TILE_SNAKE_BODY
	}
}

func snakeMove(snake *snake_t) {
	snake.body = snake.body[1:len(snake.body)]
	segment := snake.body[len(snake.body) - 1]

	segment.x += snake.vel_x
	segment.y += snake.vel_y

	if segment.x < 0 {
		segment.x = WIDTH - 1
	} else if segment.x >= WIDTH {
		segment.x = 0
	} else if segment.y < 0 {
		segment.y = HEIGHT - 1
	} else if segment.y >= HEIGHT {
		segment.y = 0
	}

	snake.body = append(snake.body, segment)
}

func snakeSwitchDirection(snake *snake_t, direction direction_t) {
	switch direction {
	case dirRight:
		snake.vel_x = 1
		snake.vel_y = 0
	case dirLeft:
		snake.vel_x = -1
		snake.vel_y = 0
	case dirUp:
		snake.vel_x = 0
		snake.vel_y = -1
	case dirDown:
		snake.vel_x = 0
		snake.vel_y = 1
	default:
		panic("Wrong argument")
	}
}

func main() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)
	var gameMap [WIDTH][HEIGHT]rune
	var snake   snake_t

	snakeInit(&snake)

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Input gorutine
	go func() {
		for {
			select {
			case <- done:
				fmt.Println("Over")
				return
			case <-ticker.C:
				char, key, err := keyboard.GetKey()
				if err != nil {
					panic(err)
				}

				if key == keyboard.KeyEsc {
					panic("ESC")
				}

				switch char {
				case 'd', 'D':
					snakeSwitchDirection(&snake, dirRight)
				case 'a', 'A':
					snakeSwitchDirection(&snake, dirLeft)
				case 'w', 'W':
					snakeSwitchDirection(&snake, dirUp)
				case 's', 'S':
					snakeSwitchDirection(&snake, dirDown)
				}
			}
		}
	} ()

	// Draw gorutine
	go func() {
		for {
			select {
			case <- done:
				fmt.Println("Over")
				return
			case <-ticker.C:
				__cleanScreen()

				mapFlush(&gameMap)
				snakeDraw(snake, &gameMap)
				snakeMove(&snake)

				mapPrint(gameMap)
			}
		}
	} ()

	for {
		select {
		case <- done:
			break
		}
	}

	done <- true
}
