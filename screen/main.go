package screen

import "fmt"

func Clear() {
	fmt.Print("\033[2J")
	// Clears the scrollback as well
	fmt.Print("\033[3J\033[0;0H")
}

func MoveTopLeft() {
	fmt.Print("\033[H")
}
