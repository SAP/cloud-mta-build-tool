package handlers

import (
	"fmt"
	"net/http"
)

type Port struct {
	Port string
}

func (p *Port) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%s", p.Port)))
}
