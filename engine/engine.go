package engine

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/cayleygraph/cayley/graph"
	"github.com/google/cayley"
)

type Engine struct {
	Store cayley.Handle
}

func New() *Engine {
	engine := Engine{}
	return &engine
}

func (engien *Engine) DatabaseSetUp(user, host, ssl string, baseName string) {
	var err error

	db, err := sql.Open("postgres", "postgresq"+"://"+user+"@"+host+"?"+"sslmode="+ssl)
	if err != nil {
		log.Fatal(err)
	}

	err = graph.InitQuadStore("sql", "postgresql://root@192.168.99.100:26257/cayley?sslmode=disable", graph.Options{"flavor": "cockroach"})

	store, err := cayley.NewGraph("sql", "postgresql://root@192.168.99.100:26257/cayley?sslmode=disable", graph.Options{"flavor": "cockroach"})
	if err != nil {
		fmt.Println(err)
	}

	defer store.Close()

}

// SavePriceForProductOfCompany method for save subject, predicate and object in graph database
func (engine *Engine) SavePriceForProductOfCompany(item Item) (*PriceOfProduct, error) {
	return nil, nil
}
