package service

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
	"github.com/mangopkg/mango/types"
)

type Service struct {
	Routes []*types.RInfo
	Response
	R    *chi.Mux
	Ext  map[string]interface{}
	Docs *types.Docs
}

type ServiceInit struct {
	R    *chi.Mux
	Ext  map[string]interface{}
	Info spec.InfoProps
}

func (u *Service) SetupHandler(mountAt string, f reflect.Type, v reflect.Value) {
	r := chi.NewMux()
	u.createRoutes(mountAt)
	u.loadMethods(f, v, mountAt)
	u.loadRoutes(r, mountAt)
	u.R.Mount(mountAt, r)
}

func (u *Service) mapkey(value string) (key string, ok bool) {
	for _, v := range u.Routes {
		if v.Func == value {
			key = v.Pattern
			ok = true
			return
		}
	}
	return
}

func (u *Service) findRinfo(mount string, method func(http.ResponseWriter, *http.Request), key string, mName string) {
	for i, v := range u.Routes {
		if v.MountAt == mount && v.Func == mName {
			u.Routes[i].Handler = method

		}
	}
}

func (u *Service) loadMethods(f reflect.Type, v reflect.Value, mount string) {
	for i := 0; i < f.NumMethod(); i++ {
		method := f.Method(i)
		key, ok := u.mapkey(method.Name)
		if ok {
			v := v.MethodByName(method.Name).Call([]reflect.Value{})
			y := v[0].Interface().(func(http.ResponseWriter, *http.Request))
			u.findRinfo(mount, y, key, method.Name)
		}
	}
}

func (u *Service) loadRoutes(r *chi.Mux, subR string) {
	for _, v := range u.Routes {
		switch v.Method {
		case "POST":
			r.Post(v.Pattern, v.Handler)
		case "GET":
			r.Get(v.Pattern, v.Handler)
		case "DELETE":
			r.Delete(v.Pattern, v.Handler)
		case "HEAD":
			r.Head(v.Pattern, v.Handler)
		case "PUT":
			r.Put(v.Pattern, v.Handler)
		case "CONNECT":
			r.Connect(v.Pattern, v.Handler)
		case "OPTIONS":
			r.Options(v.Pattern, v.Handler)
		case "TRACE":
			r.Trace(v.Pattern, v.Handler)
		case "PATCH":
			r.Patch(v.Pattern, v.Handler)
		}
	}
}

func (u *Service) createRoutes(p string) {
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
					var js types.RInfo
					if err := json.Unmarshal([]byte(v[1]), &js); err != nil {
						log.Panicln(err)
					}
					js.MountAt = p
					u.Routes = append(u.Routes, &js)
				}
			}
		}
	}
}

func (u *Service) GenerateDocs() {
	for _, v := range u.Routes {
		nD, err := json.Marshal(v.Info)

		if err != nil {
			log.Panicln(err)
		}

		if string(nD) == "null" {
			continue
		}

		el, ok := u.Docs.Paths[v.MountAt+v.Pattern]

		if ok {

			var result map[string]interface{}

			if err := json.Unmarshal(el, &result); err != nil {
				log.Panicln(err)
			}

			if err := json.Unmarshal(nD, &result); err != nil {
				log.Panicln(err)
			}

			rD, err := json.Marshal(result)

			if err != nil {
				log.Panicln(err)
			}

			u.Docs.Paths[v.MountAt+v.Pattern] = json.RawMessage(string(rD))

		} else {
			u.Docs.Paths[v.MountAt+v.Pattern] = json.RawMessage(string(nD))
		}
	}

	u.Docs.Swagger = "2.0"

	b, err := json.Marshal(*u.Docs)

	if err != nil {
		log.Panicln(err)
	}

	u.hostDocs(b, u.R)
}

func (u *Service) hostDocs(data []byte, r *chi.Mux) {

	var swagger spec.Swagger

	aerr := json.Unmarshal(data, &swagger)
	if aerr != nil {
		log.Fatal(aerr)
	}

	root := "./dist/swagger-ui"

	fs := http.FileServer(http.Dir(root))

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})

	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(swagger)
	})
}

func New(s ServiceInit) *Service {
	ns := &Service{}
	ns.Routes = []*types.RInfo{}
	ns.Response = Response{}
	ns.R = s.R
	ns.Ext = s.Ext

	docs := types.Docs{
		Paths: make(map[string]json.RawMessage),
	}

	ns.Docs = &docs

	ns.Docs.Info = s.Info

	return ns
}
