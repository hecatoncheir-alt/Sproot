package engine

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	// sql driver
	_ "github.com/cayleygraph/cayley/graph/sql"
)

// Engine is a main object of engine pkg
type Engine struct {
	Store       *cayley.Handle
	SQLDataBase *sql.DB
}

// New is a constructor for Engine
func New() *Engine {
	engine := Engine{}
	return &engine
}

// DatabaseSetUp is a method for setup SQL database for graph engine
func (engine *Engine) DatabaseSetUp(user, host string, port int, ssl string, baseName string) (*sql.DB, *cayley.Handle, error) {
	var err error
	baseType := "postgres"

	dbAddr := fmt.Sprintf("%vql://%v@%v:%v?sslmode=%v", baseType, user, host, port, ssl)

	db, err := sql.Open(baseType, dbAddr)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	engine.SQLDataBase = db

	db.Query("drop database items")
	createDbQuery := fmt.Sprintf("create database if not exists %v", baseName)

	_, err = db.Query(createDbQuery)
	if err != nil {
		fmt.Println("daaa")
		log.Fatal(err)
		return db, nil, err
	}

	tableAddr := fmt.Sprintf("%vql://%v@%v:%v/%v?sslmode=%v", baseType, user, host, port, baseName, ssl)

	err = graph.InitQuadStore("sql", tableAddr, graph.Options{"flavor": "cockroach"})
	if err != nil {
		// log.Println(err)
	}

	store, err := cayley.NewGraph("sql", tableAddr, graph.Options{"flavor": "cockroach"})
	if err != nil {
		log.Fatal(err)
		return db, nil, err
	}

	// defer store.Close()

	engine.Store = store

	return db, store, nil
}

// SavePriceForProductOfCompany method for save subject, predicate and object in graph database
func (engine *Engine) SavePriceForProductOfCompany(item *Item) (*PriceOfProduct, error) {
	price := PriceOfProduct{Name: item.Name}
	return &price, nil
}
