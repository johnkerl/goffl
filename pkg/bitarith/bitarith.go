// Package bitarith provides bit-manipulation utilities (msb, lsb, popcount, floor_log2, etc.).
// Algorithms mostly due to aggregate.org/MAGIC
package bitarith

func Msb32(x uint32) uint32 {
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x & ^(x >> 1)
}

func Lsb32(x uint32) uint32 {
	return x & (^(x) + 1)
}

func Ones32(x uint32) int {
	x = (x & 0x55555555) + ((x >> 1) & 0x55555555)
	x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
	x = (x & 0x0F0F0F0F) + ((x >> 4) & 0x0F0F0F0F)
	x = (x & 0x00FF00FF) + ((x >> 8) & 0x00FF00FF)
	x = (x & 0x0000FFFF) + ((x >> 16) & 0x0000FFFF)
	return int(x)
}

func FloorLog2_32(x uint32) int {
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return Ones32(x) - 1
}

func MsbPos32(x uint32) int {
	return FloorLog2_32(Msb32(x))
}

func LsbPos32(x uint32) int {
	return FloorLog2_32(Lsb32(x))
}

func MsbPos(x uint64) int {
	if x == 0 {
		return -1
	}
	count := 0
	for {
		word := uint32(x & 0xFFFFFFFF)
		x >>= 32
		p := MsbPos32(word)
		if x == 0 {
			return p + count
		}
		count += 32
	}
}

func Msb(x uint64) uint64 {
	shiftAmount := uint(1)
	xshift := x >> shiftAmount
	for xshift > 0 {
		x |= xshift
		shiftAmount <<= 1
		xshift = x >> shiftAmount
	}
	return x & ^(x >> 1)
}

func Lsb(x uint64) uint64 {
	return x & (^(x) + 1)
}

func Ones(x uint64) int {
	if x == 0 {
		return 0
	}
	count := 0
	for x != 0 {
		word := uint32(x & 0xFFFFFFFF)
		x >>= 32
		count += Ones32(word)
	}
	return count
}

func ExactLog2(x uint64) int {
	if x == 0 {
		return -1
	}
	count := 0
	for {
		word := uint32(x & 0xFFFFFFFF)
		if word != 0 {
			return count + FloorLog2_32(word)
		}
		x >>= 32
		count += 32
	}
}

func LsbPos(x uint64) int {
	if x == 0 {
		return -1
	}
	return ExactLog2(Lsb(x))
}
