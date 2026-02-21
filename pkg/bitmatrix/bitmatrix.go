package bitmatrix

import (
	"fmt"
	"goffl/pkg/bitvector"
)

// BitMatrix is a matrix of bit vectors over GF(2). Supports row echelon form and kernel basis.
type BitMatrix struct {
	numRows, numCols int
	Rows             []*bitvector.BitVector
}

func New(numRows, numCols int) (*BitMatrix, error) {
	if numRows <= 0 || numCols <= 0 {
		return nil, fmt.Errorf("BitMatrix: dimensions must be > 0; got %d x %d", numRows, numCols)
	}
	rows := make([]*bitvector.BitVector, numRows)
	for i := 0; i < numRows; i++ {
		v, err := bitvector.New(numCols)
		if err != nil {
			return nil, err
		}
		rows[i] = v
	}
	return &BitMatrix{numRows: numRows, numCols: numCols, Rows: rows}, nil
}

func (m *BitMatrix) NumRows() int { return m.numRows }
func (m *BitMatrix) NumCols() int { return m.numCols }

func (m *BitMatrix) SetHexOutput()  { bitvector.SetHexOutput() }
func (m *BitMatrix) SetBinaryOutput() { bitvector.SetBinaryOutput() }

func (m *BitMatrix) Row(i int) *bitvector.BitVector { return m.Rows[i] }

// RowEchelonForm reduces the matrix to row echelon form in-place.
func (m *BitMatrix) RowEchelonForm() {
	m.rowReduceBelow()
	for row := 0; row < m.numRows; row++ {
		for row2 := row + 1; row2 < m.numRows; row2++ {
			row2LeaderPos := m.Rows[row2].FindLeaderPos()
			if row2LeaderPos < 0 {
				break
			}
			rowLeaderVal, _ := m.Rows[row].Get(row2LeaderPos)
			if rowLeaderVal == 0 {
				continue
			}
			m.Rows[row].Bits ^= m.Rows[row2].Bits
		}
	}
}

func (m *BitMatrix) Rank() int {
	rr := m.clone()
	rr.rowReduceBelow()
	return rr.RankRR()
}

func (m *BitMatrix) RankRR() int {
	r := 0
	for i := 0; i < m.numRows; i++ {
		if m.Rows[i].Bits == 0 {
			return r
		}
		r++
	}
	return r
}

func (m *BitMatrix) rowReduceBelow() {
	topRow := 0
	leftColumn := 0

	for topRow < m.numRows && leftColumn < m.numCols {
		if topRow < m.numRows-1 {
			pivotRow := topRow
			pivotSuccessful := false
			for !pivotSuccessful && pivotRow < m.numRows {
				val, _ := m.Rows[pivotRow].Get(leftColumn)
				if val != 0 {
					if topRow != pivotRow {
						m.Rows[topRow], m.Rows[pivotRow] = m.Rows[pivotRow], m.Rows[topRow]
					}
					pivotSuccessful = true
				} else {
					pivotRow++
				}
			}
			if !pivotSuccessful {
				leftColumn++
				continue
			}
		}

		val, _ := m.Rows[topRow].Get(leftColumn)
		if val != 0 {
			for row := topRow + 1; row < m.numRows; row++ {
				v, _ := m.Rows[row].Get(leftColumn)
				if v != 0 {
					m.Rows[row].Bits ^= m.Rows[topRow].Bits
				}
			}
		}
		leftColumn++
		topRow++
	}
}

func (m *BitMatrix) clone() *BitMatrix {
	other, _ := New(m.numRows, m.numCols)
	for i := 0; i < m.numRows; i++ {
		other.Rows[i].Bits = m.Rows[i].Bits
	}
	return other
}

// KernelBasis returns a basis for the nullspace, or nil if nullity is zero.
func (m *BitMatrix) KernelBasis() (*BitMatrix, error) {
	rr := m.clone()
	rr.RowEchelonForm()
	rank := rr.RankRR()
	dimker := rr.numCols - rank
	if dimker == 0 {
		return nil, nil
	}

	basis, _ := New(dimker, rr.numCols)
	freeFlags := make([]int, m.numCols)
	for i := range freeFlags {
		freeFlags[i] = 1
	}
	freeIndices := make([]int, m.numCols)
	nfree := 0

	for i := 0; i < rank; i++ {
		depPos := rr.Rows[i].FindLeaderPos()
		if depPos >= 0 {
			freeFlags[depPos] = 0
		}
	}
	for i := 0; i < m.numCols; i++ {
		if freeFlags[i] != 0 {
			freeIndices[nfree] = i
			nfree++
		}
	}
	if nfree != dimker {
		return nil, fmt.Errorf("coding error detected: kernel_basis")
	}

	for i := 0; i < dimker; i++ {
		basis.Rows[i].Set(freeIndices[i], 1)
		for j := 0; j < rank; j++ {
			v, _ := rr.Rows[j].Get(freeIndices[i])
			if v == 0 {
				continue
			}
			depPos := rr.Rows[j].FindLeaderPos()
			if depPos < 0 {
				return nil, fmt.Errorf("coding error detected: kernel_basis")
			}
			basis.Rows[i].Set(depPos, v)
		}
	}
	return basis, nil
}
