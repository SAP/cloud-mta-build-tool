package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type Env struct {
}

func (p *Env) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	formatJson := r.URL.Query().Get("json")
	if formatJson == "" {
		out := "<dl>"
		for _, env := range os.Environ() {
			kv := strings.Split(env, "=")
			out += "<dt>" + kv[0] + "</dt>"
			out += "<dd>" + kv[1] + "</dd>"
		}
		out += "</dl>"
		styledTemplate.Execute(w, Body{Body: `<div class="envs">` + out + `</div>`})
	} else {
		envs := [][]string{}
		for _, env := range os.Environ() {
			envs = append(envs, strings.Split(env, "="))
		}
		json.NewEncoder(w).Encode(envs)
	}
}
