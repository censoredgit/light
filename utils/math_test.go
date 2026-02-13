package utils

import "testing"

func TestMath_ManyMaxArgs(t *testing.T) {
	x := Max[int](0, 2, 4, 1, -5, 111111, 43)
	if x != 111111 {
		t.Error("x: ", x)
	}
}

func TestMath_OneMaxArg(t *testing.T) {
	x := Max[int](2)
	if x != 2 {
		t.Error("x: ", x)
	}
}

func TestMath_TwoMaxArgs(t *testing.T) {
	x := Max[int](2, 3)
	if x != 3 {
		t.Error("x: ", x)
	}
}

func TestMath_ManyMinArgs(t *testing.T) {
	x := Min[int](0, 2, 4, 1, -5, 111111, 43)
	if x != -5 {
		t.Error("x: ", x)
	}
}

func TestMath_OneMinArg(t *testing.T) {
	x := Min[int](2)
	if x != 2 {
		t.Error("x: ", x)
	}
}

func TestMath_TwoMinArgs(t *testing.T) {
	x := Min[int](2, 3)
	if x != 2 {
		t.Error("x: ", x)
	}
}
