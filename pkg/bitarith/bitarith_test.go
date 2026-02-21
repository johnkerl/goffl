package bitarith

import "testing"

func TestMsb32(t *testing.T) {
	tests := []struct {
		x, want uint32
	}{
		{0x00000000, 0},
		{0x00000001, 1},
		{0x00000002, 2},
		{0x00000003, 2},
		{0x00000004, 4},
		{0x000000AA, 0x80},
		{0x000000FF, 0x80},
		{0x0000BEEF, 0x8000},
		{0xDEADBEEF, 0x80000000},
	}
	for _, tt := range tests {
		if got := Msb32(tt.x); got != tt.want {
			t.Errorf("Msb32(0x%08X) = 0x%08X, want 0x%08X", tt.x, got, tt.want)
		}
	}
}

func TestLsb32(t *testing.T) {
	tests := []struct {
		x, want uint32
	}{
		{0, 0},
		{1, 1},
		{3, 1},
		{0xAA, 2},
	}
	for _, tt := range tests {
		if got := Lsb32(tt.x); got != tt.want {
			t.Errorf("Lsb32(%d) = %d, want %d", tt.x, got, tt.want)
		}
	}
}

func TestOnes32(t *testing.T) {
	tests := []struct {
		x    uint32
		want int
	}{
		{0, 0},
		{3, 2},
		{0xFF, 8},
		{0xDEADBEEF, 24},
	}
	for _, tt := range tests {
		if got := Ones32(tt.x); got != tt.want {
			t.Errorf("Ones32(0x%X) = %d, want %d", tt.x, got, tt.want)
		}
	}
}

func TestMsbPos(t *testing.T) {
	tests := []struct {
		x    uint64
		want int
	}{
		{0, -1},
		{1, 0},
		{0x13, 4},
		{0xDEADBEEF, 31},
	}
	for _, tt := range tests {
		if got := MsbPos(tt.x); got != tt.want {
			t.Errorf("MsbPos(0x%X) = %d, want %d", tt.x, got, tt.want)
		}
	}
}

func TestLsbPos(t *testing.T) {
	if got := LsbPos(0); got != -1 {
		t.Errorf("LsbPos(0) = %d, want -1", got)
	}
	if got := LsbPos(0x0C); got != 2 {
		t.Errorf("LsbPos(0x0C) = %d, want 2", got)
	}
}
