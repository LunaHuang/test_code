package game

import (
	"fmt"
	"math/rand/v2"
)

// Difficulty 是題目的難度，以提示數的多寡區分。
type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

var givensByDifficulty = map[Difficulty]int{
	Easy:   40,
	Medium: 32,
	Hard:   26,
}

// ParseDifficulty 把字串轉成難度，不認得時回傳錯誤。
func ParseDifficulty(s string) (Difficulty, error) {
	d := Difficulty(s)
	if _, ok := givensByDifficulty[d]; !ok {
		return "", fmt.Errorf("未知的難度 %q", s)
	}
	return d, nil
}

// Generate 產生一個只有唯一解、提示數符合難度的題目。d 必須是已知的難度。
func Generate(d Difficulty, rng *rand.Rand) Grid {
	target := givensByDifficulty[d]
	for {
		if p, ok := tryGenerate(target, rng); ok {
			return p
		}
	}
}

// tryGenerate 從一個隨機填滿的盤面依隨機順序挖空，
// 每挖一格都確認仍是唯一解；挖不到剩 target 個提示數時回傳 false。
func tryGenerate(target int, rng *rand.Rand) (Grid, bool) {
	var p Grid
	randomFill(&p, 0, rng)
	givens := 81
	for _, i := range rng.Perm(81) {
		if givens == target {
			break
		}
		r, c := i/9, i%9
		d := p[r][c]
		p[r][c] = 0
		if n, _ := countSolutions(p); n == 1 {
			givens--
		} else {
			p[r][c] = d
		}
	}
	return p, givens == target
}

// randomFill 以隨機的數字順序回溯填滿盤面，產生一個隨機的完整解答。
func randomFill(g *Grid, pos int, rng *rand.Rand) bool {
	if pos == 81 {
		return true
	}
	r, c := pos/9, pos%9
	for _, d := range rng.Perm(9) {
		if canPlace(g, r, c, d+1) {
			g[r][c] = d + 1
			if randomFill(g, pos+1, rng) {
				return true
			}
		}
	}
	g[r][c] = 0
	return false
}
