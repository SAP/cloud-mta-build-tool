package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cloudfoundry-samples/test-app/helpers"
)

type Exit struct {
	Time time.Time
}

func (p *Exit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	index, _ := helpers.FetchIndex()

	w.WriteHeader(http.StatusOK)
	styledTemplate.Execute(w, Body{
		Class: "goodbye",
		Body: fmt.Sprintf(`
<div class="hello">
	Shutting Down
</div>

<div class="my-index">My Index Is</div>

<div class="index">%d</div>
<div class="mid-color">Uptime: %s</div>
<div class="bottom-color"></div>
    `, index, time.Since(p.Time)),
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Test App shutting down")
		os.Exit(1)
	}()
}
