package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runCLI 以暫存資料夾當資料位置執行一次 todo 指令，回傳 stdout、stderr 與錯誤。
func runCLI(t *testing.T, dataDir string, args ...string) (string, string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	err := run(args, dataDir, &stdout, &stderr)
	return stdout.String(), stderr.String(), err
}

// mustRun 執行指令並要求成功，回傳 stdout。
func mustRun(t *testing.T, dataDir string, args ...string) string {
	t.Helper()
	out, errOut, err := runCLI(t, dataDir, args...)
	if err != nil {
		t.Fatalf("todo %v 失敗：%v（stderr：%q）", args, err, errOut)
	}
	return out
}

func TestAddedItemAppearsInList(t *testing.T) {
	dir := t.TempDir()
	mustRun(t, dir, "add", "買", "牛奶")

	got := mustRun(t, dir, "list")
	want := "[ ] 1  買 牛奶\n"
	if got != want {
		t.Errorf("list 輸出 = %q，想要 %q", got, want)
	}
}

func TestAddsGetIncreasingIdsAndListIsInIdOrder(t *testing.T) {
	dir := t.TempDir()
	if got := mustRun(t, dir, "add", "買牛奶"); got != "已新增 1  買牛奶\n" {
		t.Errorf("第一次 add 輸出 = %q", got)
	}
	if got := mustRun(t, dir, "add", "繳電費"); got != "已新增 2  繳電費\n" {
		t.Errorf("第二次 add 輸出 = %q", got)
	}
	mustRun(t, dir, "add", "倒垃圾")

	got := mustRun(t, dir, "list")
	want := "[ ] 1  買牛奶\n[ ] 2  繳電費\n[ ] 3  倒垃圾\n"
	if got != want {
		t.Errorf("list 輸出 = %q，想要 %q", got, want)
	}
}

func TestEmptyListSaysNothingToDo(t *testing.T) {
	got := mustRun(t, t.TempDir(), "list")
	if got != "沒有待辦事項\n" {
		t.Errorf("空清單 list 輸出 = %q", got)
	}
}

func TestMissingDataDirIsCreatedOnFirstAdd(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data")
	mustRun(t, dir, "add", "買牛奶")

	if _, err := os.Stat(filepath.Join(dir, "todos.json")); err != nil {
		t.Fatalf("資料檔沒有建立：%v", err)
	}
	if got := mustRun(t, dir, "list"); got != "[ ] 1  買牛奶\n" {
		t.Errorf("list 輸出 = %q", got)
	}
}

func TestInvalidCommandsFailWithoutTouchingTheList(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		want      string // 錯誤訊息應包含的文字
		wantUsage bool   // stderr 是否應印出用法
	}{
		{"沒給指令", nil, "缺少指令", true},
		{"未知指令", []string{"fly"}, "未知的指令：fly", true},
		{"add 沒給標題", []string{"add"}, "請提供待辦事項的標題", false},
		{"add 標題只有空白", []string{"add", " ", ""}, "請提供待辦事項的標題", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			mustRun(t, dir, "add", "買牛奶")

			_, errOut, err := runCLI(t, dir, tt.args...)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("錯誤 = %v，想要包含 %q", err, tt.want)
			}
			if got := strings.Contains(errOut, "用法"); got != tt.wantUsage {
				t.Errorf("stderr 有沒有用法 = %v，想要 %v（stderr：%q）", got, tt.wantUsage, errOut)
			}
			if got := mustRun(t, dir, "list"); got != "[ ] 1  買牛奶\n" {
				t.Errorf("失敗的指令改動了清單：%q", got)
			}
		})
	}
}

func TestCorruptDataFileIsReportedInChinese(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "todos.json"), []byte("{壞掉"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err := runCLI(t, dir, "list")
	if err == nil || !strings.Contains(err.Error(), "讀取資料檔失敗") {
		t.Fatalf("錯誤 = %v，想要包含「讀取資料檔失敗」", err)
	}
}

func TestDoneMarksItemAsDoneInList(t *testing.T) {
	dir := t.TempDir()
	mustRun(t, dir, "add", "買牛奶")
	mustRun(t, dir, "add", "繳電費")
	mustRun(t, dir, "add", "倒垃圾")

	if got := mustRun(t, dir, "done", "2"); got != "已完成 2  繳電費\n" {
		t.Errorf("done 輸出 = %q", got)
	}
	got := mustRun(t, dir, "list")
	want := "[ ] 1  買牛奶\n[x] 2  繳電費\n[ ] 3  倒垃圾\n"
	if got != want {
		t.Errorf("list 輸出 = %q，想要 %q", got, want)
	}
}

func TestDoneAgainOnlyPrintsHintAndLeavesDataFileUnchanged(t *testing.T) {
	dir := t.TempDir()
	mustRun(t, dir, "add", "買牛奶")
	mustRun(t, dir, "done", "1")
	before := readDataFile(t, dir)

	if got := mustRun(t, dir, "done", "1"); got != "1  買牛奶 已經完成了\n" {
		t.Errorf("再次 done 輸出 = %q", got)
	}
	if after := readDataFile(t, dir); after != before {
		t.Errorf("再次 done 改動了資料檔：\n前：%s\n後：%s", before, after)
	}
}

// readDataFile 讀出資料檔的原始內容，用來確認指令沒有改動它。
func readDataFile(t *testing.T, dir string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, "todos.json"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestDoneWithBadIdFailsWithoutTouchingDataFile(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string // 錯誤訊息應包含的文字
	}{
		{"沒給 id", []string{"done"}, "請提供待辦事項的 id"},
		{"id 不是數字", []string{"done", "abc"}, "id 必須是數字：abc"},
		{"id 不存在", []string{"done", "9"}, "找不到 id 為 9 的待辦事項"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			mustRun(t, dir, "add", "買牛奶")
			before := readDataFile(t, dir)

			_, _, err := runCLI(t, dir, tt.args...)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("錯誤 = %v，想要包含 %q", err, tt.want)
			}
			if after := readDataFile(t, dir); after != before {
				t.Errorf("失敗的 done 改動了資料檔：\n前：%s\n後：%s", before, after)
			}
		})
	}
}
