// 網頁版數獨：可選難度、確認作答、顯示解答。
// 伺服器不保存局的狀態，見 docs/adr/0001-stateless-server.md。
package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"math/rand/v2"
	"net/http"

	"sudoku/game"
)

//go:embed web
var webFiles embed.FS

func main() {
	addr := flag.String("addr", ":8080", "伺服器監聽位址")
	flag.Parse()

	web, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServerFS(web))
	mux.HandleFunc("GET /api/new", handleNew)
	mux.HandleFunc("POST /api/check", handleCheck)
	mux.HandleFunc("POST /api/solve", handleSolve)

	log.Printf("數獨伺服器啟動：http://localhost%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

// handleNew 依 difficulty 參數產生新題目。
func handleNew(w http.ResponseWriter, r *http.Request) {
	d, err := game.ParseDifficulty(r.URL.Query().Get("difficulty"))
	if err != nil {
		writeError(w, err)
		return
	}
	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	writeJSON(w, map[string]game.Grid{"puzzle": game.Generate(d, rng)})
}

// handleCheck 比對作答與解答，回傳錯格與是否完成。
func handleCheck(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Puzzle  game.Grid `json:"puzzle"`
		Entries game.Grid `json:"entries"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}
	res, err := game.Check(req.Puzzle, req.Entries)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, res)
}

// handleSolve 回傳題目的解答。
func handleSolve(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Puzzle game.Grid `json:"puzzle"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}
	solution, err := game.Solve(req.Puzzle)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, map[string]game.Grid{"solution": solution})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Print(err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
