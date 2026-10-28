// Package game 實作數獨的出題、解題與確認作答。
package game

import "errors"

// Grid 是 9×9 盤面，0 代表空格。
type Grid [9][9]int

var (
	// ErrNoSolution 表示題目沒有解答。
	ErrNoSolution = errors.New("題目沒有解答")
	// ErrMultipleSolutions 表示題目有不只一組解答。
	ErrMultipleSolutions = errors.New("題目有多組解答")
)

// Solve 回傳題目唯一的解答。
func Solve(puzzle Grid) (Grid, error) {
	if !validGivens(puzzle) {
		return Grid{}, ErrNoSolution
	}
	switch n, solution := countSolutions(puzzle); n {
	case 0:
		return Grid{}, ErrNoSolution
	case 1:
		return solution, nil
	default:
		return Grid{}, ErrMultipleSolutions
	}
}

// countSolutions 回傳 g 的解答數（0、1，或 2 代表兩組以上），以及只有一組時的解答。
func countSolutions(g Grid) (int, Grid) {
	var found Grid
	return search(&g, &found, 2), found
}

// validGivens 回報盤面上的數字是否都在 0–9 之間，且同列、同行、同宮沒有重複。
func validGivens(g Grid) bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			d := g[r][c]
			if d < 0 || d > 9 {
				return false
			}
			if d == 0 {
				continue
			}
			g[r][c] = 0
			ok := canPlace(&g, r, c, d)
			g[r][c] = d
			if !ok {
				return false
			}
		}
	}
	return true
}

// search 以回溯法計算 g 的解答數，最多數到 limit 為止；
// 找到的解答寫入 found（只有一組解時才有意義）。
// 每次挑候選數最少的空格來填，以縮小搜尋範圍。
func search(g *Grid, found *Grid, limit int) int {
	r, c, best := -1, -1, 10
	for i := 0; i < 81 && best > 1; i++ {
		if g[i/9][i%9] != 0 {
			continue
		}
		n := 0
		for d := 1; d <= 9; d++ {
			if canPlace(g, i/9, i%9, d) {
				n++
			}
		}
		if n < best {
			r, c, best = i/9, i%9, n
		}
	}
	if r < 0 {
		*found = *g
		return 1
	}
	count := 0
	for d := 1; d <= 9 && count < limit; d++ {
		if canPlace(g, r, c, d) {
			g[r][c] = d
			count += search(g, found, limit-count)
		}
	}
	g[r][c] = 0
	return count
}

// canPlace 回報 d 能否放在 (r, c) 而不與同列、同行、同宮重複。
func canPlace(g *Grid, r, c, d int) bool {
	br, bc := r/3*3, c/3*3
	for i := 0; i < 9; i++ {
		if g[r][i] == d || g[i][c] == d || g[br+i/3][bc+i%3] == d {
			return false
		}
	}
	return true
}
