package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const usage = `用法：
  todo add <標題>   新增待辦事項
  todo list         列出所有待辦事項
  todo done <id>    把待辦事項標成完成
  todo delete <id>  刪除待辦事項`

// run 執行一次 todo 指令，待辦清單存放在 dataDir。
// 回傳的錯誤只說明原因；沒給指令或指令不認得時，另外把用法寫到 stderr。
func run(args []string, dataDir string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		fmt.Fprintln(stderr, usage)
		return errors.New("缺少指令")
	}
	l, err := load(dataDir)
	if err != nil {
		return err
	}
	switch args[0] {
	case "add":
		title := strings.TrimSpace(strings.Join(args[1:], " "))
		if title == "" {
			return errors.New("請提供待辦事項的標題")
		}
		it := l.add(title, time.Now())
		if err := l.save(dataDir); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "已新增 %d  %s\n", it.ID, it.Title)
	case "list":
		if len(l.Items) == 0 {
			fmt.Fprintln(stdout, "沒有待辦事項")
		}
		for _, it := range l.Items {
			mark := " "
			if it.Done {
				mark = "x"
			}
			fmt.Fprintf(stdout, "[%s] %d  %s\n", mark, it.ID, it.Title)
		}
	case "done":
		id, err := parseID(args[1:])
		if err != nil {
			return err
		}
		i, err := l.indexOf(id)
		if err != nil {
			return err
		}
		it := &l.Items[i]
		if it.Done {
			fmt.Fprintf(stdout, "%d  %s 已經完成了\n", it.ID, it.Title)
			return nil
		}
		it.Done = true
		if err := l.save(dataDir); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "已完成 %d  %s\n", it.ID, it.Title)
	case "delete":
		id, err := parseID(args[1:])
		if err != nil {
			return err
		}
		it, err := l.remove(id)
		if err != nil {
			return err
		}
		if err := l.save(dataDir); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "已刪除 %d  %s\n", it.ID, it.Title)
	default:
		fmt.Fprintln(stderr, usage)
		return fmt.Errorf("未知的指令：%s", args[0])
	}
	return nil
}

// parseID 從指令參數取出待辦事項的 id。
func parseID(args []string) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("請提供待辦事項的 id")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("id 必須是數字：%s", args[0])
	}
	return id, nil
}
