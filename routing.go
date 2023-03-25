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
	"github.com/go-openapi/spec"
	"github.com/mangopkg/mango/types"
)

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

					var opProps spec.PathItemProps

					infoChunk, err := json.Marshal(js.Info)

					if err != nil {
						log.Panicln(err)
					}

					if err := json.Unmarshal(infoChunk, &opProps); err != nil {
						log.Panicln(err)
					}

					if js.Info == nil {

						resp := spec.ResponsesProps{
							StatusCodeResponses: make(map[int]spec.Response),
						}

						resp.StatusCodeResponses[200] = spec.Response{
							ResponseProps: spec.ResponseProps{
								Description: "No description provided",
							},
						}

						defProp := spec.OperationProps{
							Tags: []string{js.MountAt},
							Responses: &spec.Responses{
								ResponsesProps: resp,
							},
						}

						if js.Method == "GET" {
							js.Info = spec.PathItemProps{
								Get: &spec.Operation{
									OperationProps: defProp,
								},
							}
						}
						if js.Method == "POST" {
							js.Info = spec.PathItemProps{
								Post: &spec.Operation{
									OperationProps: defProp,
								},
							}
						}
						if js.Method == "PUT" {
							js.Info = spec.PathItemProps{
								Put: &spec.Operation{
									OperationProps: defProp,
								},
							}
						}
						if js.Method == "DELETE" {
							js.Info = spec.PathItemProps{
								Delete: &spec.Operation{
									OperationProps: defProp,
								},
							}
						}
						if js.Method == "OPTIONS" {
							js.Info = spec.PathItemProps{
								Options: &spec.Operation{
									OperationProps: defProp,
								},
							}
						}
						if js.Method == "HEAD" {
							js.Info = spec.PathItemProps{
								Head: &spec.Operation{
									OperationProps: defProp,
								},
							}
						}
						if js.Method == "PATCH" {
							js.Info = spec.PathItemProps{
								Patch: &spec.Operation{
									OperationProps: defProp,
								},
							}
						}

					}

					u.Routes = append(u.Routes, &js)
				}
			}
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
