package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cloudfoundry-samples/test-app/helpers"
)

type Hello struct {
	Time time.Time
}

func (p *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	index, _ := helpers.FetchIndex()

	styledTemplate.Execute(w, Body{Body: fmt.Sprintf(`
<div class="hello">
	Test App	
</div>

<div class="my-index">My Index Is</div>

<div class="index">%d</div>
<div class="mid-color">Uptime: %s</div>
<div class="bottom-color"></div>
    `, index, time.Since(p.Time))})
}
