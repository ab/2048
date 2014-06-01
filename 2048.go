package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    "time"
    "math/rand"
)

type Move int

const (
    MOVE_UP Move = iota
    MOVE_DOWN
    MOVE_RIGHT
    MOVE_LEFT
)

const DEBUG = true
var debug *bool = flag.Bool("debug", false, "enable debug logging")

func Debugf(format string, args ...interface{}) {
    if DEBUG || *debug {
        log.Printf("Debug: " + format, args...)
    }
}
func Errorf(format string, args ...interface{}) {
    log.Printf("Error: " + format, args...)
}

func AllMoves() []Move {
    moves := make([]Move, 4)
    moves = append(moves, MOVE_UP)
    moves = append(moves, MOVE_DOWN)
    moves = append(moves, MOVE_RIGHT)
    moves = append(moves, MOVE_LEFT)
    return moves
}

type Board struct {
    rows, cols int
    cells [][]*Cell
}

type Cell struct {
    value int
    merged bool
    just_merged bool
    just_created bool
}

func NewCell() *Cell {
    return new(Cell)
}

func (c *Cell) isEmpty() bool {
    return c.value == 0
}

func (c *Cell) clear() {
    c.value = 0
    c.merged = false
}

func (c *Cell) newTurn() {
    c.just_merged = c.merged
    c.merged = false
    c.just_created = false
}

func (c *Cell) printColor() {
    var escape string
    var lpad, rpad int

    if c.isEmpty() {
        fmt.Printf("    ")
        return
    }

    if c.value < 10 {
        lpad = 3
        rpad = 1
    } else if c.value < 100 {
        lpad = 2
        rpad = 1
    } else if c.value < 1000 {
        lpad = 1
        rpad = 1
    } else if c.value < 10000 {
        lpad = 1
        rpad = 0
    } else {
        lpad = 0
        rpad = 0
    }

    switch c.value {
    case 2, 4: escape = "1;30;47m"
    case 8, 16: escape = "1;34;47m"
    case 32, 64: escape = "1;41m"
    case 128, 256: escape = "1;43m"
    case 512, 1024: escape = "1;43m"
    case 2048, 4096: escape = "1;43m"
    case 8192, 16384: escape = "1;43m"
    case 32768, 65536: escape = "1;43m"
    default:
        escape = "m"
    }

    fmt.Printf("\033[%s%s%d%s\033[m", escape, strings.Repeat(" ", lpad),
               c.value, strings.Repeat(" ", rpad))
}

func NewBoard(rows, cols int) *Board {
    b := new(Board)
    b.rows = rows
    b.cols = cols

    b.cells = make([][]*Cell, b.rows)
    for row_i := range b.cells {
        b.cells[row_i] = make([]*Cell, b.cols)
        for col_i := range b.cells[row_i] {
            b.cells[row_i][col_i] = NewCell()
        }
    }

    return b
}

func (board *Board) print_row_divider() {
    for c := 0; c < board.cols; c++ {
        fmt.Print("+––––")
        if DEBUG {
            fmt.Print("–")
        }
    }
    fmt.Println("+")
}

func (board *Board) display() {

    for _, row := range board.cells {
        board.print_row_divider()

        for _, cell := range row {

            fmt.Print("|")
            // cell.printColor()

            if cell.isEmpty() {
                fmt.Print("    ")
            } else if cell.just_created {
                fmt.Printf(" <%d>", cell.value)
            } else if cell.value > 999 {
                fmt.Printf("%4d", cell.value)
            } else {
                fmt.Printf("%3d ", cell.value)
            }
            if DEBUG {
                if cell.just_merged {
                    fmt.Print("+")
                } else {
                    fmt.Print(" ")
                }
            }
        }
        fmt.Println("|")
    }

    board.print_row_divider()
}

func (board *Board) newTurn() {
    // clear merged state on all cells
    for _, row := range board.cells {
        for _, cell := range row {
            cell.newTurn()
        }
    }
}

func (board *Board) move_direction(row_offset, col_offset int,
                                   dry_run bool) bool {
    made_changes := false

    // exactly one of row_offset, col_offset should be +/- 1
    if !((row_offset == 1 || row_offset == -1) && col_offset == 0) &&
       !((col_offset == 1 || col_offset == -1) && row_offset == 0) {
        panic(fmt.Sprintf("Invalid move: %v, %v", row_offset, col_offset))
    }

    var row_iteration, col_iteration, row_start, col_start int

    if row_offset == 1 {
        row_iteration = -1
        row_start = board.rows - 1
    } else {
        row_iteration = 1
        row_start = 0
    }
    if col_offset == 1 {
        col_iteration = -1
        col_start = board.cols - 1
    } else {
        col_iteration = 1
        col_start = 0
    }

    for row_i := row_start; row_i >= 0 && row_i < board.rows;
            row_i += row_iteration {

        // for each column in the row, look for places to move to
        for col_i := col_start; col_i >= 0 && col_i < board.cols;
                col_i += col_iteration {

            cell := board.cells[row_i][col_i]

            if cell.isEmpty() {
                // empty cell
                continue
            }

            Debugf("Current cell: (%d,%d) value %d", row_i, col_i, cell.value)

            var neighbor *Cell
            var action string
            var neighbor_row, neighbor_col int

            for row_look, col_look := row_i + row_offset, col_i + col_offset ;;
            {

                if row_look < 0 || row_look >= board.rows ||
                   col_look < 0 || col_look >= board.cols {
                    break
                }

                // three cases
                look_cell := board.cells[row_look][col_look]

                // case 1: target is empty => should move to it
                if look_cell.isEmpty() {

                    // assert that we haven't merged before
                    if cell.merged {
                        panic("Somehow cell has merged before")
                    }

                    action = "move"
                    neighbor = look_cell
                    neighbor_row = row_look
                    neighbor_col = col_look

                // case 2: target is equal => check if we can merge
                } else if look_cell.value == cell.value {
                    // can only merge if neither cell has merged this turn
                    if look_cell.merged || cell.merged {
                        break
                    } else {
                        action = "merge"
                        neighbor = look_cell
                        neighbor_row = row_look
                        neighbor_col = col_look
                    }
                // case 3: target is not equal, cannot move
                } else {
                    // TODO: log.debug(cannot merge with XX)
                    break
                }

                row_look += row_offset
                col_look += col_offset
            }

            switch action {
                case "move":
                    // swap neighbor and cell
                    made_changes = true
                    if !dry_run {
                        board.cells[row_i][col_i] = neighbor
                        board.cells[neighbor_row][neighbor_col] = cell
                    }
                case "merge":
                    // merge neighbor and cell
                    made_changes = true
                    if !dry_run {
                        neighbor.value *= 2
                        neighbor.merged = true
                        cell.clear()
                    }
            }
        }
    }

    return made_changes
}

func (board *Board) move_up(dry_run bool) bool {
    return board.move_direction(-1, 0, dry_run)
}
func (board *Board) move_down(dry_run bool) bool {
    return board.move_direction(1, 0, dry_run)
}
func (board *Board) move_right(dry_run bool) bool {
    return board.move_direction(0, 1, dry_run)
}
func (board *Board) move_left(dry_run bool) bool {
    return board.move_direction(0, -1, dry_run)
}

func (board *Board) tryMove(move Move, dry_run bool) bool {
    switch move {
    case MOVE_UP: return board.move_up(dry_run)
    case MOVE_DOWN: return board.move_down(dry_run)
    case MOVE_RIGHT: return board.move_right(dry_run)
    case MOVE_LEFT: return board.move_left(dry_run)
    default:
        panic(fmt.Sprintf("Unexpected move: %#v", move))
    }
}

func (board *Board) isGameOver() bool {
    if !board.isFull() {
        return false
    }

    for _, move := range AllMoves() {
        if board.tryMove(move, true) {
            // successful move dry run
            return false
        }
    }
    return true
}

func (board *Board) isFull() bool {
    for _, row := range board.cells {
        for _, cell := range row {
            if cell.isEmpty() {
                return false
            }
        }
    }
    return true
}

func (board *Board) add_number() {
    if board.isFull() {
        panic("add_number: board is full")
    }

    // TODO: determine circumstances when we should add a number > 2

    to_add := 2

    for {
        row_index := rand.Intn(board.rows)
        col_index := rand.Intn(board.cols)
        if board.cells[row_index][col_index].isEmpty() {
            board.cells[row_index][col_index].value = to_add
            board.cells[row_index][col_index].just_created = true
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

func (g *Game) newTurn() {
    g.board.newTurn()
    g.turn += 1
}

func (g *Game) display() {
    fmt.Printf("Turn %d:\n", g.turn)
    g.board.display()
    fmt.Println()
}

func (g *Game) promptMove() Move {
    for {
        fmt.Print("Your move: ")
        reader := bufio.NewReader(os.Stdin)
        line, _ := reader.ReadString('\n')
        line = strings.TrimSpace(line)

        switch line {
        case "u", "U", "\x1b[A":
            return MOVE_UP
        case "d", "D", "\x1b[B":
            return MOVE_DOWN
        case "r", "R", "\x1b[C":
            return MOVE_RIGHT
        case "l", "L", "\x1b[D":
            return MOVE_LEFT
        }

        fmt.Printf("Invalid move: %#v\n", line)
    }
}

func (g *Game) doMove() {
    for {
        move := g.promptMove()
        if g.board.tryMove(move, false) {
            break
        }
        fmt.Println("Cannot move that way")
    }
}

func promptContinue() {
    fmt.Printf("Press enter to continue...")
    reader := bufio.NewReader(os.Stdin)
    reader.ReadString('\n')
}

func init() {
    // fmt.Println("Start time:", game.start_time)
}

func makeSampleBoard(b *Board) {
    v := 1
    for _, row := range b.cells {
        for _, cell := range row {
            v = v << 1
            cell.value = v
        }
    }
}

func main() {
    var samplegame = NewGame()

    makeSampleBoard(samplegame.board)
    samplegame.display()

    promptContinue()

    var game = NewGame()

    game.board.add_number()

    for !game.board.isGameOver() {
        // promptContinue()

        game.board.add_number()
        game.display()

        game.doMove()
        game.newTurn()

        game.display()
    }
}
