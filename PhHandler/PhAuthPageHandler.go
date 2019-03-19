package PhHandler

import (
		"net/http"
	"github.com/julienschmidt/httprouter"
	"os"
	"log"
)

type PhAuthPageHandler struct {
	Method     string
	HttpMethod string
	Args       []string
}

func (h PhAuthPageHandler) NewAuthPageHandler(args ...interface{}) PhAuthPageHandler {
	var hm string
	var md string
	var ag []string
	for i, arg := range args {
		if i == 0 {
		} else if i == 1 {
			md = arg.(string)
		} else if i == 2 {
			hm = arg.(string)
		} else if i == 3 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		} else {
		}
	}

	return PhAuthPageHandler{Method: md, HttpMethod: hm, Args: ag}
}

func (h PhAuthPageHandler) Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	log.Println("Start ===> Load `Auth` Page")
	file, err := os.Open(h.Args[0])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return 1
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, r, file.Name(), fi.ModTime(), file)
	return 0
}

func (h PhAuthPageHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhAuthPageHandler) GetHandlerMethod() string {
	return h.Method
}
