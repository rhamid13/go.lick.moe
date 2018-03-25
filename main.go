package main

import (
	"fmt"
	"net/http"
	"database/sql"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func PrintPaste(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	var paste string
	err := db.QueryRow("SELECT paste FROM pastemap WHERE id = ?", params["id"]).Scan(&paste)
	if (err == sql.ErrNoRows) {
		fmt.Fprint(w, "paste not found.")
	} else {
		fmt.Fprintf(w, paste)
	}
}

func Paste(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	text := r.Form["paste"][0]
	sum := sha256.Sum256([]byte(text))

	// string representation of truncated sha256
	shortsum := hex.EncodeToString(sum[0:4])
	fmt.Fprint(w, "lick.moe/p/"+shortsum)

	stmt, _ := db.Prepare("INSERT INTO pastemap(id, paste) values(?,?)")
	stmt.Exec(shortsum, text)
}

var db *sql.DB

func main() {
	// database shit
	db, _ = sql.Open("sqlite3", "pastedb")
	stmt, _ := db.Prepare("CREATE TABLE IF NOT EXISTS pastemap (id VARCHAR(8), paste TEXT)")
	stmt.Exec()

	router := mux.NewRouter()
	router.HandleFunc("/p/{id}", PrintPaste).Methods("GET")
	router.HandleFunc("/paste", Paste).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", router)
	http.ListenAndServe(":8000", router)
}
