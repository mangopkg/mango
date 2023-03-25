package mango

import (
	"encoding/json"
	"reflect"

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
