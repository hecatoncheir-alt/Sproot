package storage

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	dataBaseClient "github.com/dgraph-io/dgraph/client"
	dataBaseAPI "github.com/dgraph-io/dgraph/protos/api"
)

var (
	databaseHost = "192.168.99.100"
	databasePort = 9080
)

// Storage is a object with database resource
type Storage struct {
	GraphAddress  string
	GraphGRPCHost string
	GraphGRPCPort int

	Client     *dataBaseClient.Dgraph
	Categories *Categories
	Companies  *Companies
}

// New is a constructor for Storage objects
func New(host string, port int) *Storage {
	storage := &Storage{}

	storage.GraphGRPCHost = host
	storage.GraphGRPCPort = port
	storage.GraphAddress = fmt.Sprintf("%v:%v", host, port)

	return storage
}

func (storage *Storage) prepareDataBaseClient() (*dataBaseClient.Dgraph, error) {
	conn, err := grpc.Dial(storage.GraphAddress, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	baseClient := dataBaseAPI.NewDgraphClient(conn)
	databaseGraph := dataBaseClient.NewDgraphClient(baseClient)

	return databaseGraph, nil
}

// SetUp is a method of storage for prepare database client and objects of resource of database.
func (storage *Storage) SetUp() (err error) {
	storage.Client, err = storage.prepareDataBaseClient()
	if err != nil {
		return err
	}

	storage.Categories = NewCategoriesResourceForStorage(storage)
	err = storage.Categories.SetUp()
	if err != nil {
		return err
	}

	storage.Companies = NewCompaniesResourceForStorage(storage)
	err = storage.Companies.SetUp()
	if err != nil {
		return err
	}

	return nil
}

// DeleteAll drop all records in database
func (storage *Storage) DeleteAll() error {
	return storage.Client.Alter(context.Background(), &dataBaseAPI.Operation{DropAll: true})
}
