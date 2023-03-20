package mango

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"reflect"
	"regexp"

	"github.com/go-chi/chi/v5"
)

type RInfo struct {
	Pattern string   `json:"pattern"`
	Func    string   `json:"func"`
	Method  string   `json:"method"`
	Auth    string   `json:"auth"`
	Roles   []string `json:"roles"`
	ReqBody string   `json:"reqBody"`
	Handle  func(http.ResponseWriter, *http.Request)
}

type Response struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Error      bool        `json:"error"`
}

func (r *Response) Send(w http.ResponseWriter) {
	w.WriteHeader(r.StatusCode)
	if r.Error {
		r.Data = nil
		if err := json.NewEncoder(w).Encode(r); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(r); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type Service struct {
	Routes map[string]RInfo
	Response
	R   *chi.Mux
	Ext map[string]interface{}
}

type ServiceInit struct {
	R   *chi.Mux
	Ext map[string]interface{}
}

func (u *Service) mapkey(m map[string]RInfo, value string) (key string, ok bool) {
	for _, v := range m {
		if v.Func == value {
			key = v.Pattern
			ok = true
			return
		}
	}
	return
}

func (u *Service) LoadMethods(f reflect.Type, v reflect.Value, m *map[string]func(http.ResponseWriter, *http.Request), mount string) {
	for i := 0; i < f.NumMethod(); i++ {
		method := f.Method(i)
		key, ok := u.mapkey(u.Routes, method.Name)
		if ok {
			v := v.MethodByName(method.Name).Call([]reflect.Value{})
			y := v[0].Interface().(func(http.ResponseWriter, *http.Request))
			(*m)[key] = y
		}
	}
}

func (u *Service) SetupHandler(mountAt string, f reflect.Type, v reflect.Value, m map[string]func(http.ResponseWriter, *http.Request)) {
	r := chi.NewMux()
	u.CreateRoutes(mountAt)
	u.LoadMethods(f, v, &m, mountAt)
	u.LoadRoutes(m, r, mountAt)
	u.R.Mount(mountAt, r)
}

func (u *Service) LoadRoutes(m map[string]func(http.ResponseWriter, *http.Request), r *chi.Mux, subR string) {
	for _, v := range u.Routes {
		switch v.Method {
		case "POST":
			r.Post(v.Pattern, m[v.Pattern])
		case "GET":
			r.Get(v.Pattern, m[v.Pattern])
		case "DELETE":
			r.Delete(v.Pattern, m[v.Pattern])
		case "HEAD":
			r.Head(v.Pattern, m[v.Pattern])
		case "PUT":
			r.Put(v.Pattern, m[v.Pattern])
		case "CONNECT":
			r.Connect(v.Pattern, m[v.Pattern])
		case "OPTIONS":
			r.Options(v.Pattern, m[v.Pattern])
		case "TRACE":
			r.Trace(v.Pattern, m[v.Pattern])
		case "PATCH":
			r.Patch(v.Pattern, m[v.Pattern])
		}
	}
}

func (u *Service) CreateRoutes(p string) {
	fset := token.NewFileSet()

	d, err := parser.ParseDir(fset, "."+p, nil, parser.ParseComments)
	if err != nil {
		log.Panicln(err)
	}

	for _, f := range d {
		for _, f := range f.Files {
			for _, c := range f.Comments {
				left := "<@route"
				right := ">"
				rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(left) + `(.*?)` + regexp.QuoteMeta(right))
				matches := rx.FindAllStringSubmatch(c.Text(), -1)
				for _, v := range matches {
					var js = RInfo{}
					if err := json.Unmarshal([]byte(v[1]), &js); err != nil {
						log.Panicln(err)
					}
					u.Routes[p+js.Pattern] = js
				}
			}
		}

	}
}

func New(s ServiceInit) Service {

	ns := Service{}

	ns.Routes = make(map[string]RInfo)
	ns.Response = Response{}
	ns.R = s.R
	ns.Ext = s.Ext

	return ns
}
