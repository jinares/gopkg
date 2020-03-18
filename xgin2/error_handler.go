package xgin2

import (
	"fmt"
	"net/http"
)

func errorHandler() func(writer http.ResponseWriter, request *http.Request, err error) {

	return func(writer http.ResponseWriter, request *http.Request, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
		writer.WriteHeader(http.StatusOK)

		writer.Write([]byte(err.Error()))

	}
}
