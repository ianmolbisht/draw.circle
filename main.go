package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    _ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("sqlite", "./circle.db")
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS scores (
        name TEXT PRIMARY KEY,
        score INTEGER
    )`)
    if err != nil {
        log.Fatal(err)
    }
}

func submitScore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var payload struct {
        Name  string `json:"name"`
        Score int    `json:"score"`
    }

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    // Check if user exists and get current score
    var currentScore int
    err := db.QueryRow("SELECT score FROM scores WHERE name = ?", payload.Name).Scan(&currentScore)
    
    if err == sql.ErrNoRows {
        // User doesn't exist, insert new record
        _, err = db.Exec("INSERT INTO scores(name, score) VALUES(?, ?)", payload.Name, payload.Score)
    } else if err != nil {
        http.Error(w, "Database error", 500)
        return
    } else {
        // User exists, update only if new score is higher
        if payload.Score > currentScore {
            _, err = db.Exec("UPDATE scores SET score = ? WHERE name = ?", payload.Score, payload.Name)
        }
    }
    if err != nil {
        http.Error(w, "Database error", 500)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query(`SELECT name, score FROM scores ORDER BY score DESC LIMIT 10`)
    if err != nil {
        http.Error(w, "DB query error", 500)
        return
    }
    defer rows.Close()

    var leaderboard []map[string]interface{}
    for rows.Next() {
        var name string
        var score int
        rows.Scan(&name, &score)
        leaderboard = append(leaderboard, map[string]interface{}{"name": name, "score": score})
    }
    json.NewEncoder(w).Encode(leaderboard)
}

func main() {
    initDB()
    http.Handle("/", http.FileServer(http.Dir("static")))
    http.HandleFunc("/score", submitScore)
    http.HandleFunc("/leaderboard", getLeaderboard)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
