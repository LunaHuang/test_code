# 數獨

網頁版數獨遊戲：選擇難度開新局，可隨時確認作答或顯示解答。

## 執行

```sh
cd sudoku
go run .
```

瀏覽器開啟 <http://localhost:8080>。要改用其他 port：`go run . -addr :9000`。

## 玩法

- 從下拉選單選難度（簡單／中等／困難），按「新遊戲」產生題目。
- 點選空格後，用鍵盤 1–9 或畫面上的數字鍵作答；Backspace、0 或「清除」清除；方向鍵移動。
- 「確認答案」：標出錯格；全部填對時顯示完成並鎖住盤面。
- 「顯示答案」：填上整盤解答並鎖住盤面。

## 結構

- [`game/`](game/)：出題、解題與確認作答，測試在 `game_test.go`。
- [`main.go`](main.go)：HTTP 伺服器與 JSON API（`/api/new`、`/api/check`、`/api/solve`）。
- [`web/index.html`](web/index.html)：前端頁面，以 `embed` 打包進執行檔。
- 詞彙：[`GLOSSARY.md`](GLOSSARY.md)；決策：[`docs/adr/`](docs/adr/)。
