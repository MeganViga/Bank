package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

var testQueries *Queries
const (
	//dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/bank?sslmode=disable"
)
func TestMain(m *testing.M){
	conn, err := pgx.Connect(context.Background(), dbSource)

	if err != nil{
		log.Fatal("cannot connect to database:", err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}