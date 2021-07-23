package main

import (
	"errors"
	"gorm.io/gorm"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/logger"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jsannemo/omogenhost/storage"
	"github.com/jsannemo/omogenhost/webapi/problems"
	apipb "github.com/jsannemo/omogenhost/webapi/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type dbConfig struct {
	Server string
	Port   int
}

type apiConfig struct {
	Server string
	Port   int
}

type config struct {
	Database dbConfig
	Webapi   apiConfig
}

func main() {
	defer logger.Init("webapi", true, false, ioutil.Discard).Close()
	data, err := ioutil.ReadFile("/etc/omogenjudge/webapi.toml")
	if err != nil {
		panic(err)
	}
	var conf config
	if _, err := toml.Decode(string(data), &conf); err != nil {
		panic(err)
	}
	connStr := fmt.Sprintf("postgres://omogenjudge:omogenjudge@%s:%d/omogenjudge?sslmode=disable", conf.Database.Server, conf.Database.Port)
	if err := storage.Init(connStr); err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	http.HandleFunc("/problems/img/", handleAttachment)
	httpServer := &http.Server{
		Addr: fmt.Sprintf("%s:%d", conf.Webapi.Server, conf.Webapi.Port),
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if wrappedGrpc.IsGrpcWebRequest(req) {
				wrappedGrpc.ServeHTTP(resp, req)
				return
			}
			// Fall back to other servers.
			http.DefaultServeMux.ServeHTTP(resp, req)
		}),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	apipb.RegisterProblemServiceServer(grpcServer, problems.InitProblemService())
	log.Fatal(httpServer.ListenAndServe())
}

func handleAttachment(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path[len("/problems/img/"):]
	slash := strings.IndexRune(path, '/')
	if slash == -1 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	problem := path[:slash]
	filePath := path[slash+1:]
	var file storage.StoredFile
	if res := storage.GormDB.Debug().Select("FileContents").Joins(
		"LEFT JOIN problem_statement_file psf ON psf.statement_file_hash = stored_file.file_hash").Joins(
			"LEFT JOIN Problem on Problem.problem_id = psf.problem_id").Where("problem.short_name = ?", problem).Where("psf.file_path = ?", filePath).Find(&file); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			writer.WriteHeader(http.StatusNotFound)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(file.FileContents)
}
