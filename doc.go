package mango

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/spec"
)

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

//go:embed dist
var distFS embed.FS

func (u *Service) hostDocs(data []byte, r *chi.Mux) {
	var swagger spec.Swagger

	err := json.Unmarshal(data, &swagger)
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.FS(distFS))

	r.Get("/api/*", http.StripPrefix("/api", fs).ServeHTTP)

	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(swagger)
	})
}
