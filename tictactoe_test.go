package main

import (
	"testing"
	"strconv"
)

func TestMoveDiag(t *testing.T){
	game := NewTicTacToe(3) 
	want := false 
	result := game.Move(0, 0, 2)

	if result != want{
		t.Errorf(`game.move(1, 1, 2) resulted in %s, wanted result was false`, strconv.FormatBool(result))
	}
}

func TestMoveAntiDiag(t *testing.T){
	game := NewTicTacToe(3)
	want := false
	result := game.Move(0,2,1)

	if result != want{
		t.Errorf(`game.move(1, 1, 2) resulted in %s, wanted result was false`, strconv.FormatBool(result))
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
	want := true
	result := game.CheckIfLegalMove(0, 0)

	if result != want {
		t.Errorf(`game.CheckIfLegalMove(0, 0) resulted in %s, wanted result was true`, strconv.FormatBool(result))
	}
}

func TestCheckIlegalMove(t *testing.T){
	game := NewTicTacToe(3)
	want := false

	game.Move(0, 0, 1)

	result := game.CheckIfLegalMove(0, 0)

	if result != want {
		t.Errorf(`game.CheckIfLegalMove(0, 0) resulted in %s, wanted result was false`, strconv.FormatBool(result))
	}

}
