package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// item 是一筆待辦事項。
type item struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// list 是整份待辦清單；NextID 是下一筆要發的 id，永遠遞增，刪除後也不重複使用。
type list struct {
	NextID int    `json:"next_id"`
	Items  []item `json:"items"`
}

const dataFile = "todos.json"

// load 讀取 dir 裡的待辦清單；檔案不存在時回傳空清單。
func load(dir string) (*list, error) {
	b, err := os.ReadFile(filepath.Join(dir, dataFile))
	if errors.Is(err, fs.ErrNotExist) {
		return &list{NextID: 1}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("讀取資料檔失敗：%w", err)
	}
	var l list
	if err := json.Unmarshal(b, &l); err != nil {
		return nil, fmt.Errorf("讀取資料檔失敗：%w", err)
	}
	return &l, nil
}

// save 把整份待辦清單寫回 dir，必要時建立資料夾。
func (l *list) save(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, dataFile), b, 0o644)
}

// add 新增一筆未完成的待辦事項並回傳它。
func (l *list) add(title string, now time.Time) item {
	it := item{ID: l.NextID, Title: title, CreatedAt: now}
	l.NextID++
	l.Items = append(l.Items, it)
	return it
}
