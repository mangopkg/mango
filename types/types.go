package types

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/spec"
)

type Docs struct {
	Paths   map[string]json.RawMessage `json:"paths"`
	Info    spec.InfoProps             `json:"info"`
	Swagger string                     `json:"swagger"`
}

type RInfo struct {
	Pattern string                                   `json:"pattern"`
	Func    string                                   `json:"func"`
	Method  string                                   `json:"method"`
	Auth    string                                   `json:"auth"`
	ReqBody string                                   `json:"reqBody"`
	Info    interface{}                              `json:"Info"`
	MountAt string                                   `json:"mountAt"`
	Handler func(http.ResponseWriter, *http.Request) `json:"handler"`
}
