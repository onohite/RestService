package main

import (
	"RestService/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const HTTPS = "https://"

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				errorLog.Println("Body close")
			}
		}(r.Body)
		if err != nil {
			errorLog.Println("ReadAll")
			http.Error(w, "Bad Request ", 400)
			return
		}
		requestAddress := &model.RequestModel{}
		err = requestAddress.UnmarshalJSON(body)
		if err != nil {
			errorLog.Println("Unmarshal")
			http.Error(w, "Bad Request ", 400)
			return
		}
		body, err = getSiteBody(*requestAddress)
		responseModel := &model.ResponseModel{HTML: body}
		json, err := responseModel.MarshalJSON()
		if err != nil {
			errorLog.Println("Unmarshal")
			http.Error(w, "Internal Server Error", 500)
			return
		}
		_, err = w.Write(json)
		if err != nil {
			errorLog.Println("response Write Error")
			http.Error(w, "Connection Timed Out", 522)
			return
		}
	} else {
		http.Error(w, "Method "+r.Method+" Not Allowed", 405)
		return
	}
}

func getSiteBody(requestModel model.RequestModel) ([]byte, error) {
	site, err := http.Get(HTTPS + requestModel.Address)
	if err != nil {
		errorLog.Println("http.Get error")
		return nil, err
	}
	infoLog.Println("connection to request Address complete")
	body, err := ioutil.ReadAll(site.Body)
	if err != nil {
		errorLog.Println("ReadAll error")
		return nil, err
	}
	return body, nil
}

func main() {
	http.HandleFunc("/getHTML", handler)
	infoLog.Println("starting server: localhost:8080")

	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		errorLog.Println(err)
	}
}
