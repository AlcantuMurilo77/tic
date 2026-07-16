package main

import (
	"testing"
	"strconv"
)

func TestMoveDiag(t *testing.T){
	game := NewTicTacToe(3) 
	game.Move(0, 0, 1)
	game.Move(1, 1, 1)
	won := game.Move(2, 2, 1)

	if !won {
		t.Error("expected diagonal victory")
	}
}

func TestMoveAntiDiag(t *testing.T){
	game := NewTicTacToe(3)

	game.Move(0, 2, 1)
	game.Move(1, 1, 1)
	won := game.Move(2, 0, 1)


	if !won {
		t.Error("expected diagonal victory")
	}

}

func TestVictory(t *testing.T){
	game := NewTicTacToe(3)
	want := true

	game.Move(0,0, 1)
	game.Move(1,0,1)
	result := game.Move(2,0,1)

	if result != want {
		t.Errorf(`game.CheckIfLegalMove(0, 0) resulted in %s, wanted result was true`, strconv.FormatBool(result))
	}


}

func TestCheckLegalMove(t *testing.T){
	game := NewTicTacToe(3)
	
	game.Move(0, 0, 1)

	if game.Move(0, 0, 1){
		t.Fatal("expected move to be illegal")
	}
}
