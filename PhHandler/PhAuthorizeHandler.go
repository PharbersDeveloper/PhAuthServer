package PhHandler

import (
	"log"
	"net/http"
	"reflect"
	"gopkg.in/oauth2.v3/server"
	"github.com/julienschmidt/httprouter"
	"github.com/PharbersDeveloper/PhAuthServer/PhServer"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"net/url"
)

type PhAuthorizeHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
	srv        *server.Server
}

func (h PhAuthorizeHandler) NewAuthorizeHandler(args ...interface{}) PhAuthorizeHandler {
	var m *BmMongodb.BmMongodb
	var r *BmRedis.BmRedis
	var hm string
	var md string
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmMongodb" {
					m = dm.(*BmMongodb.BmMongodb)
				}
				if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
				}
			}
		} else if i == 1 {
			md = arg.(string)
		} else if i == 2 {
			hm = arg.(string)
		} else if i == 3 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		}
	}
	sv := PhServer.GetInstance(m, r)

	return PhAuthorizeHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r, srv: sv}
}

func (h PhAuthorizeHandler) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	log.Println("Start ===> Authorize Validation")

	var form url.Values
	redisDriver := h.rd.GetRedisClient()
	defer redisDriver.Close()
	returnUri, err := redisDriver.Get("ReturnUri").Result()
	if returnUri != "" {
		form, _ = url.ParseQuery(returnUri)
	}
	r.Form = form
	redisDriver.Del("ReturnUri")

	err = h.srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return 0
}

func (h PhAuthorizeHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PhAuthorizeHandler) GetHandlerMethod() string {
	return h.Method
}
