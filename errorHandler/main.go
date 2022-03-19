package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	tc "errorHandler/errorHandler/core/category"

	tb "errorHandler/errorHandler/core/blog"

	us "errorHandler/errorHandler/core/user"

	tpc "errorHandler/gunk/v1/category"

	tpb "errorHandler/gunk/v1/blog"

	tpu "errorHandler/gunk/v1/user"

	//	err "errorHandler/gunk/v1/errorHandler"

	"errorHandler/errorHandler/storage/postgres"

	"errorHandler/errorHandler/services/blog"
	"errorHandler/errorHandler/services/category"
	"errorHandler/errorHandler/services/user"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("category/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("error loading configuration: %v", err)
	}

	grpcServer := grpc.NewServer()

	store, err := newDBFromConfig(config)
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	cs := tc.NewCoreSve(store)

	bc := tb.NewCoreSve(store)

	usr := us.NewCoreSve(store)

	u := user.NewUserServer(usr)

	s := category.NewCategoryServer(cs)

	b := blog.NewCategoryServer(bc)

	tpu.RegisterUserRegServiceServer(grpcServer, u)

	tpc.RegisterCategoryServiceServer(grpcServer, s)

	tpb.RegisterBlogServiceServer(grpcServer, b)

	host, port := config.GetString("server.host"), config.GetString("server.port")

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	log.Printf("Server is starting on: http://%s:%s\n", host, port)

	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}

func newDBFromConfig(config *viper.Viper) (*postgres.Storage, error) {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("user")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("dbname")
	if password := cf("password"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslMode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransacionTimeout")
	db, err := postgres.NewStorage(dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}
