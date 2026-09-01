package game

type TicTacToe struct {
	Size         int
	Board        [][]int
	Rows         []int
	Cols         []int
	Diagonal     int
	AntiDiagonal int
}

func NewTicTacToe(size int) *TicTacToe {
	board := make([][]int, size)
	for i := range board {
		board[i] = make([]int, size)
	}
	return &TicTacToe{
		Size:  size,
		Board: board,
		Rows:  make([]int, size),
		Cols:  make([]int, size),
	}
}

func (t *TicTacToe) CheckIfLegalMove(row int, col int) bool {
	return t.Board[row][col] == 0
}

func (t *TicTacToe) Move(row int, col int, player int) bool {

	if t.Board[row][col] != 0 {
		return false
	}
	t.Board[row][col] = player

	value := 1
	if player == 2 {
		value = -1
	}

	t.Rows[row] += value
	t.Cols[col] += value

	if row == col {
		t.Diagonal += value
	}

	if row+col == t.Size-1 {
		t.AntiDiagonal += value
	}

	if abs(t.Rows[row]) == t.Size ||
		abs(t.Cols[col]) == t.Size ||
		abs(t.Diagonal) == t.Size ||
		abs(t.AntiDiagonal) == t.Size {
		return true
	}

	return false

}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func NewTicTacToeFromBoard(board [][]int) *TicTacToe {
	size := len(board)

	t := &TicTacToe{
		Size:  size,
		Board: board,
		Rows:  make([]int, size),
		Cols:  make([]int, size),
	}

	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			player := board[row][col]

			if player == 0 {
				continue
			}

			value := 1
			if player == 2 {
				value = -1
			}

			t.Rows[row] += value
			t.Cols[col] += value

			if row == col {
				t.Diagonal += value
			}

			if row+col == size-1 {
				t.AntiDiagonal += value
			}
		}
	}

	return t
}
