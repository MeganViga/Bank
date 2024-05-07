package db

import (
	"context"

	"os"
	"testing"
	"log"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testStore *Store
const (
	//dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/bank?sslmode=disable"
)
func TestMain(m *testing.M){
	//conn, err := pgx.Connect(context.Background(), dbSource)
	pool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil{
		log.Fatal("cannot connect to database:", err)
	}
	testQueries = New(pool)
	testStore = NewStore(pool)
	os.Exit(m.Run())
}