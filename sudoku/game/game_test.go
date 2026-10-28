package game

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// parse 把 9 行、每行 9 字元的字串轉成盤面，'.' 代表空格。
func parse(t *testing.T, s string) Grid {
	t.Helper()
	var g Grid
	rows := strings.Fields(s)
	if len(rows) != 9 {
		t.Fatalf("需要 9 行，得到 %d 行", len(rows))
	}
	for r, row := range rows {
		for c, ch := range row {
			if ch != '.' {
				g[r][c] = int(ch - '0')
			}
		}
	}
	return g
}

// 維基百科「數獨」條目的範例題目與其解答。
const wikiPuzzle = `
53..7....
6..195...
.98....6.
8...6...3
4..8.3..1
7...2...6
.6....28.
...419..5
....8..79`

const wikiSolution = `
534678912
672195348
198342567
859761423
426853791
713924856
961537284
287419635
345286179`

func TestSolveReturnsTheUniqueSolution(t *testing.T) {
	got, err := Solve(parse(t, wikiPuzzle))
	if err != nil {
		t.Fatalf("Solve 回傳錯誤：%v", err)
	}
	if want := parse(t, wikiSolution); got != want {
		t.Errorf("解答錯誤\n得到 %v\n預期 %v", got, want)
	}
}

func TestSolveRejectsPuzzleWithConflictingGivens(t *testing.T) {
	p := parse(t, wikiSolution)
	p[0][1] = 5 // 第一列已經有 5，盤面雖已填滿但不合法
	if _, err := Solve(p); err != ErrNoSolution {
		t.Errorf("提示數衝突應回傳 ErrNoSolution，得到 %v", err)
	}
}

func TestSolveRejectsOutOfRangeDigits(t *testing.T) {
	p := parse(t, wikiPuzzle)
	p[0][2] = 10
	if _, err := Solve(p); err != ErrNoSolution {
		t.Errorf("超出 1–9 的數字應回傳 ErrNoSolution，得到 %v", err)
	}
}

func TestSolveRejectsPuzzleWithMultipleSolutions(t *testing.T) {
	if _, err := Solve(Grid{}); err != ErrMultipleSolutions {
		t.Errorf("空白盤面應回傳 ErrMultipleSolutions，得到 %v", err)
	}
}

func TestGenerateGivesUniquePuzzleWithDifficultyGivens(t *testing.T) {
	cases := []struct {
		difficulty Difficulty
		givens     int
	}{
		{Easy, 40},
		{Medium, 32},
		{Hard, 26},
	}
	for _, tc := range cases {
		for seed := uint64(1); seed <= 5; seed++ {
			p := Generate(tc.difficulty, rand.New(rand.NewPCG(seed, seed)))
			n := 0
			for r := 0; r < 9; r++ {
				for c := 0; c < 9; c++ {
					if p[r][c] != 0 {
						n++
					}
				}
			}
			if n != tc.givens {
				t.Errorf("%s seed %d：提示數 %d，預期 %d", tc.difficulty, seed, n, tc.givens)
			}
			if _, err := Solve(p); err != nil {
				t.Errorf("%s seed %d：題目不是唯一解：%v", tc.difficulty, seed, err)
			}
		}
	}
}

func TestCheckReportsWrongCellsButNotEmptyOnes(t *testing.T) {
	puzzle := parse(t, wikiPuzzle)
	var entries Grid
	entries[0][2] = 4 // 正確
	entries[0][3] = 9 // 錯，解答是 6
	entries[8][0] = 1 // 錯，解答是 3
	got, err := Check(puzzle, entries)
	if err != nil {
		t.Fatalf("Check 回傳錯誤：%v", err)
	}
	want := []Cell{{0, 3}, {8, 0}}
	if !slices.Equal(got.Wrong, want) {
		t.Errorf("錯格 %v，預期 %v", got.Wrong, want)
	}
	if got.Complete {
		t.Error("還有空格，不應判定完成")
	}
}

func TestCheckReportsCompleteWhenAllEntriesAreCorrect(t *testing.T) {
	puzzle := parse(t, wikiPuzzle)
	entries := parse(t, wikiSolution) // 提示數的格子也有值，應被忽略
	got, err := Check(puzzle, entries)
	if err != nil {
		t.Fatalf("Check 回傳錯誤：%v", err)
	}
	if len(got.Wrong) != 0 || !got.Complete {
		t.Errorf("應判定完成且沒有錯格，得到 %+v", got)
	}
}

func TestCheckIsNotCompleteWhenFilledBoardHasWrongCells(t *testing.T) {
	puzzle := parse(t, wikiPuzzle)
	entries := parse(t, wikiSolution)
	entries[0][2], entries[0][3] = 6, 4 // 互換兩格，盤面填滿但有錯格
	got, err := Check(puzzle, entries)
	if err != nil {
		t.Fatalf("Check 回傳錯誤：%v", err)
	}
	want := []Cell{{0, 2}, {0, 3}}
	if !slices.Equal(got.Wrong, want) || got.Complete {
		t.Errorf("應有錯格 %v 且未完成，得到 %+v", want, got)
	}
}
