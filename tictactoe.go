package main

type TicTacToe struct {
	Size int
	Board [][]int
	Rows []int
	Cols []int
	Diagonal int
	AntiDiagonal int
}

//to update the board we are keeping
//we just do t.Board[row][col] = player

func NewTicTacToe(size int) *TicTacToe {
	board := make([][]int, size)
	for i := range board {
		board[i] = make([]int, size)
	}
	return &TicTacToe {
		Size: size,
		Board: board,
		Rows: make([]int, size),
		Cols: make([]int, size),
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

	value := 1 //if X player, the value it increments on the array is 1
	if player == 2 {
		value = -1 //if O, it decreases by -1
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
		return true //checks if game ended
	}
	
	return false

	//by checking if game ended, since we received which players move it is
	//to determinate the winner we just hvae to return the last to play
	//e.g won := game.Move(0, 2, 1) player 1 won, if this returned true.

}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
