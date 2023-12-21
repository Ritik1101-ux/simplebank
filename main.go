package main

import (
	"database/sql"
	"log"

	"github.com/Ritik1101-ux/simplebank/api"
	db "github.com/Ritik1101-ux/simplebank/db/sqlc"
	"github.com/Ritik1101-ux/simplebank/utils"
	_ "github.com/lib/pq"
)


func main(){

	config,err:=utils.LoadConfig(".") //Because we have app.env in same directory

	if err!=nil{
		log.Fatal("Cannot Load Config: ",err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

    if err != nil {
        log.Fatal("Cannot Connect to the Database:", err)
    }

	store:=db.NewStore(conn)
	server,err:=api.NewServer(config,store)

	if err!=nil{
		log.Fatal(err)
	}

	err=server.Start(config.ServerAddress)

	if err!=nil{
		log.Fatal("Cannot Start server: ",err)
	}

	
}