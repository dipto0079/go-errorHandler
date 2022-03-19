package main

import (
	"errorHandler/cms/handler"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	tpc "errorHandler/gunk/v1/blog"
	tpb "errorHandler/gunk/v1/category"
	tpu "errorHandler/gunk/v1/user"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("cms/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("error loading configuration: %v", err)
	}

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	store := sessions.NewCookieStore([]byte(config.GetString("session.secret")))

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", config.GetString("blog.host"), config.GetString("blog.port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	tc := tpb.NewCategoryServiceClient(conn)

	tb := tpc.NewBlogServiceClient(conn)

	us := tpu.NewUserRegServiceClient(conn)

	r := handler.New(decoder, store, tc, tb, us)

	host, port := config.GetString("server.host"), config.GetString("server.port")

	log.Printf("Server Starting no : http://%s:%s", host, port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Fatal(err)
	}
}
