package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/alecthomas/kingpin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DB  *sql.DB
	err error

	mysqlHost = kingpin.Flag("host", "Connect to host.").Default("localhost").Short('h').String()
	mysqlUser = kingpin.Flag("user", "User for login.").Default("root").Short('u').String()
	mysqlPass = kingpin.Flag("password", "Password to use when connecting to server.").Short('h').String()
)

func main() {
	kingpin.Parse()
	DB, err = sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:3306)/mysql", *mysqlUser, *mysqlPass, *mysqlHost),
	)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", statusHandler)
	http.ListenAndServe(":9199", nil)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	err = DB.Ping()
	if err != nil {
		http.Error(w, "Galera Node is *down*. ("+err.Error()+")", 503)
		return
	}

	var key string
	var value int64
	err = DB.QueryRow("SHOW STATUS LIKE 'wsrep_local_state'").Scan(&key, &value)
	if err != nil {
		http.Error(w, "Galera Node is *down*. ("+err.Error()+")", 503)
		return
	}

	switch value {
	case 4, 2:
		fmt.Fprintf(w, "Galera Node is running.")
		return
	default:
		http.Error(w, "Galera Node is *down*. (State Mismatch)", 503)
		return
	}
}
