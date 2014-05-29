package main

import (
    "fmt"
    "time"
    "math/rand"
)

type Board struct {
    rows, cols int
    cells [][]int
}

func NewBoard(rows, cols int) *Board {
    b := new(Board)
    b.rows = rows
    b.cols = cols

    b.cells = make([][]int, b.rows)
    for i := range b.cells {
        b.cells[i] = make([]int, b.cols)
    }

    return b
}

func (board *Board) print_row_divider() {
    for c := 0; c < board.cols; c++ {
        fmt.Print("+––––")
    }
    fmt.Println("+")
}

func (board *Board) display() {

    for _, row := range board.cells {
        board.print_row_divider()

        for _, cell := range row {
            fmt.Print("|")
            if cell == 0 {
                fmt.Print("    ")
            } else if cell > 999 {
                fmt.Printf("%4d", cell)
            } else {
                fmt.Printf("%3d ", cell)
            }
        }
        fmt.Println("|")
    }

    board.print_row_divider()
}

func (board *Board) move_up() {
}
func (board *Board) move_left() {
}
func (board *Board) move_right() {
}
func (board *Board) move_down() {
}

func (board *Board) is_full() bool {
    for _, row := range board.cells {
        for _, cell := range row {
            if cell == 0 {
                return false
            }
        }
    }
    return true
}

func (board *Board) add_number() {
    if board.is_full() {
        panic("add_number: board is full")
    }

    // TODO: determine circumstances when we should add a number > 2

    to_add := 2

    for {
        row_index := rand.Intn(board.rows)
        col_index := rand.Intn(board.cols)
        if board.cells[row_index][col_index] == 0 {
            board.cells[row_index][col_index] = to_add
            return
        }
    }
}

type Game struct {
    board *Board
    turn int
    start_time time.Time
}

func NewGame() *Game {
    g := new(Game)
    g.board = NewBoard(4, 4)
    g.start_time = time.Now()
    return g
}

func init() {
    // fmt.Println("Start time:", game.start_time)
}

func main() {
    var game = NewGame()

    for !game.board.is_full() {
        game.board.display()
        game.board.add_number()
    }

    game.board.display()
}
