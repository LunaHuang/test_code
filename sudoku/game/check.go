package game

// Cell 是盤面上一格的位置：Row、Col 皆為 0–8。
type Cell struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// Result 是確認作答的結果。
type Result struct {
	Wrong    []Cell `json:"wrong"`    // 錯格，依列、行順序排列
	Complete bool   `json:"complete"` // 所有空格都已作答且沒有錯格
}

// Check 比對作答與題目的解答。entries 中題目已有提示數的格子會被忽略，0 代表未作答。
func Check(puzzle, entries Grid) (Result, error) {
	solution, err := Solve(puzzle)
	if err != nil {
		return Result{}, err
	}
	res := Result{Wrong: []Cell{}}
	filled := true
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if puzzle[r][c] != 0 {
				continue
			}
			if entries[r][c] == 0 {
				filled = false
				continue
			}
			if entries[r][c] != solution[r][c] {
				res.Wrong = append(res.Wrong, Cell{r, c})
			}
		}
	}
	res.Complete = filled && len(res.Wrong) == 0
	return res, nil
}
