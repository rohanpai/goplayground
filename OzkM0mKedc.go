// This program solves &#34;easy&#34; (i.e. fully constrained) Sudoku puzzles.
// It uses go functions to have each square concurrently &#34;solve&#34; itself
// in terms of finding a constrained answer.
// It is based upon a program from Andrew Gerrand&#39;s blog.
//
//  Nick Brandaleone - November 2012

package main

import &#34;fmt&#34;

// Sample constrained puzzle
// The coordinates are puzzle[row][column] or puzzle[y][x]

var puzzle = [9][9]int{{0, 3, 0, 9, 0, 0, 0, 8, 0}, {0, 0, 6, 2, 0, 3, 7, 9, 0}, {0, 0, 0, 1, 0, 0, 0, 0, 0},
	{0, 2, 0, 3, 0, 0, 0, 7, 0}, {0, 0, 0, 0, 7, 0, 0, 6, 4}, {1, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 5, 0, 0, 0, 4, 9, 0, 0}, {0, 7, 2, 0, 0, 0, 0, 0, 0}, {0, 9, 0, 0, 5, 0, 8, 3, 0}}

func main() {
	fmt.Println(&#34;Starting position&#34;)
	printBoard(puzzle)

	var sol = NewSolver()
	sol.InitBoard()

	solved := sol.Solve()
	fmt.Println(&#34;Finishing position&#34;)
	printBoard(solved)
}

///////////////////////////////////////////////////////////////

// Print out Sudoku board
func printBoard(board [9][9]int) {
	fmt.Printf(&#34;=====================\n&#34;)
	for y := 0; y &lt; 9; y&#43;&#43; {
		for x := 0; x &lt; 9; x&#43;&#43; {
			fmt.Printf(&#34;%d &#34;, board[y][x])
			if x%3 == 2 &amp;&amp; x &lt; 8 {
				fmt.Printf(&#34;| &#34;)
			}
		}
		fmt.Printf(&#34;\n&#34;)
		if y%3 == 2 &amp;&amp; y &lt; 8 {
			fmt.Printf(&#34;-----   -----   -----\n&#34;)
		}
	}
	fmt.Printf(&#34;=====================\n\n&#34;)
}

// Each square has its own go-routine. It determines what value it has by
// listening to see what values it CAN&#39;T be, and if it sees 8 such values,
// it knows it must be the remaining 9th value. It then sends the answer
// along the &#34;done&#34; channel.
func square(x, y int, elim &lt;-chan int, done chan&lt;- solution) {
	var eliminated [9]bool
	for n := range elim {
		eliminated[n-1] = true
		var c, s int
		for i, v := range eliminated {
			if v {
				c&#43;&#43;
			} else {
				s = i &#43; 1
			}
		}	// inner for
		if c == 8 {
			done &lt;- solution{x, y, s}
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
	s = &amp;Solver{done: make(chan solution)}
	for y := 0; y &lt; 9; y&#43;&#43; {
		for x := 0; x &lt; 9; x&#43;&#43; {
			s.elim[y][x] = make(chan int)
			go square(x, y, s.elim[y][x], s.done)
		}
	}
	return
}

// Seeds the puzzle with known values
func (s *Solver) InitBoard() {
	for y := 0; y &lt; 9; y&#43;&#43; {
		for x := 0; x &lt; 9; x&#43;&#43; {
			s.Set(x, y, puzzle[y][x])
		}
	}
}

// Sends the known values to the appropriate &#34;elim&#34; channel
func (s *Solver) Set(x, y, v int) {
	go func() {
		for i := 1; i &lt;= 9; i&#43;&#43; {
			if i != v &amp;&amp; v != 0 {
				s.elim[y][x] &lt;- i
			}
		}
	}()
}

// Loop through all 81 &#34;done&#34; channels (one for each square). Sends out determined
// values to the &#34;eliminate&#34; function to create new constraints.
func (s *Solver) Solve() [9][9]int {
	for n := 0; n &lt; 9*9; n&#43;&#43; {
		u := &lt;-s.done
		go s.eliminate(u)
		s.solution[u.y][u.x] = u.v

		fmt.Printf(&#34;Done[%d]: %#v\n&#34;, n, u)
		printBoard(s.solution)
	}
	return s.solution
}

// Send answers to all associated rows/columns/boxes so that squares
// have new information to determine constraints
func (s *Solver) eliminate(u solution) {
	// row
	for x := 0; x &lt; 9; x&#43;&#43; {
		if x != u.x &amp;&amp; u.v != 0 {
			s.elim[u.y][x] &lt;- u.v
		}
	}
	// column
	for y := 0; y &lt; 9; y&#43;&#43; {
		if y != u.y &amp;&amp; u.v != 0 {
			s.elim[y][u.x] &lt;- u.v
		}
	}
	// 3x3 group
	sX, sY := u.x/3*3, u.y/3*3	// group start coordinates
	for y := sY; y &lt; sY&#43;3; y&#43;&#43; {
		for x := sX; x &lt; sX&#43;3; x&#43;&#43; {
			if x != u.x &amp;&amp; y != u.y &amp;&amp; u.v != 0 {
				s.elim[y][x] &lt;- u.v
			}
		}
	}
}
