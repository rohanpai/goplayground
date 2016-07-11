// This program solves "easy" (i.e. fully constrained) Sudoku puzzles.
// It uses go functions to have each square concurrently "solve" itself
// in terms of finding a constrained answer.
// It is based upon a program from Andrew Gerrand's blog.
//
//  Nick Brandaleone - November 2012

package main

import "fmt"

// Sample constrained puzzle
// The coordinates are puzzle[row][column] or puzzle[y][x]

var puzzle = [9][9]int{{0, 3, 0, 9, 0, 0, 0, 8, 0}, {0, 0, 6, 2, 0, 3, 7, 9, 0}, {0, 0, 0, 1, 0, 0, 0, 0, 0},
	{0, 2, 0, 3, 0, 0, 0, 7, 0}, {0, 0, 0, 0, 7, 0, 0, 6, 4}, {1, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 5, 0, 0, 0, 4, 9, 0, 0}, {0, 7, 2, 0, 0, 0, 0, 0, 0}, {0, 9, 0, 0, 5, 0, 8, 3, 0}}

func main() {
	fmt.Println("Starting position")
	printBoard(puzzle)

	var sol = NewSolver()
	sol.InitBoard()

	solved := sol.Solve()
	fmt.Println("Finishing position")
	printBoard(solved)
}

///////////////////////////////////////////////////////////////

// Print out Sudoku board
func printBoard(board [9][9]int) {
	fmt.Printf("=====================\n")
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			fmt.Printf("%d ", board[y][x])
			if x%3 == 2 && x < 8 {
				fmt.Printf("| ")
			}
		}
		fmt.Printf("\n")
		if y%3 == 2 && y < 8 {
			fmt.Printf("-----   -----   -----\n")
		}
	}
	fmt.Printf("=====================\n\n")
}

// Each square has its own go-routine. It determines what value it has by
// listening to see what values it CAN'T be, and if it sees 8 such values,
// it knows it must be the remaining 9th value. It then sends the answer
// along the "done" channel.
func square(x, y int, elim <-chan int, done chan<- solution) {
	var eliminated [9]bool
	for n := range elim {
		eliminated[n-1] = true
		var c, s int
		for i, v := range eliminated {
			if v {
				c++
			} else {
				s = i + 1
			}
		}	// inner for
		if c == 8 {
			done <- solution{x, y, s}
			//          close(elim) or break or return ?
			/********** HERE ***********/
		}
	}	// outer for
}

type solution struct{ x, y, v int }

type Solver struct {
	done		chan solution
	elim		[9][9]chan int
	solution	[9][9]int
}

// Initializes the Solver structure with appropriate channels and go-routines.
func NewSolver() (s *Solver) {
	s = &Solver{done: make(chan solution)}
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			s.elim[y][x] = make(chan int)
			go square(x, y, s.elim[y][x], s.done)
		}
	}
	return
}

// Seeds the puzzle with known values
func (s *Solver) InitBoard() {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			s.Set(x, y, puzzle[y][x])
		}
	}
}

// Sends the known values to the appropriate "elim" channel
func (s *Solver) Set(x, y, v int) {
	go func() {
		for i := 1; i <= 9; i++ {
			if i != v && v != 0 {
				s.elim[y][x] <- i
			}
		}
	}()
}

// Loop through all 81 "done" channels (one for each square). Sends out determined
// values to the "eliminate" function to create new constraints.
func (s *Solver) Solve() [9][9]int {
	for n := 0; n < 9*9; n++ {
		u := <-s.done
		go s.eliminate(u)
		s.solution[u.y][u.x] = u.v

		fmt.Printf("Done[%d]: %#v\n", n, u)
		printBoard(s.solution)
	}
	return s.solution
}

// Send answers to all associated rows/columns/boxes so that squares
// have new information to determine constraints
func (s *Solver) eliminate(u solution) {
	// row
	for x := 0; x < 9; x++ {
		if x != u.x && u.v != 0 {
			s.elim[u.y][x] <- u.v
		}
	}
	// column
	for y := 0; y < 9; y++ {
		if y != u.y && u.v != 0 {
			s.elim[y][u.x] <- u.v
		}
	}
	// 3x3 group
	sX, sY := u.x/3*3, u.y/3*3	// group start coordinates
	for y := sY; y < sY+3; y++ {
		for x := sX; x < sX+3; x++ {
			if x != u.x && y != u.y && u.v != 0 {
				s.elim[y][x] <- u.v
			}
		}
	}
}
