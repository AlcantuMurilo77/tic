package main

type TicTacToe struct {
	size int // [ ] 
	rows []int // -
	cols []int // |
	diagonal int // \ 
	antiDiagonal int // /
}

func NewTicTacToe(size int) *TicTacToe {
	return &TicTacToe {
		size: size,
		rows: make([]int, size),
		cols: make([]int, size),
	}
}

func (t *TicTacToe) Move(row, col, player int) bool {
	value := 1 //if X player, the value it increments on the array is 1
	if player == 2 {
		value = -1 //if O, it decreases by -1
	}

	t.rows[row] += value
	t.cols[col] += value 

	if row == col {
		t.diagonal += value
	}

	if row+col == t.size-1 {
		t.antiDiagonal += value
	}

	if abs(t.rows[row]) == t.size ||
		abs(t.cols[col]) == t.size || 
		abs(t.diagonal) == t.size ||
		abs(t.antiDiagonal) == t.size {
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
