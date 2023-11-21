package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cloudfoundry-samples/test-app/handlers"
	"github.com/cloudfoundry-samples/test-app/helpers"
	"github.com/cloudfoundry-samples/test-app/routes"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/rata"
)

var message string
var quiet bool

var portsFlag = flag.String(
	"ports",
	"",
	"Comma delimited list of ports, where the app will be listening to",
)

func init() {
	flag.StringVar(&message, "message", "Hello", "The Message to Log and Display")
	flag.BoolVar(&quiet, "quiet", false, "Less Verbose Logging")
	flag.Parse()
}

func main() {
	flag.Parse()

	logger := lager.NewLogger("test-app")
	if quiet {
		logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))
	} else {
		logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	}

	ports := getServerPorts()

	logger.Info("test-app.starting", lager.Data{"ports": ports})
	index, err := helpers.FetchIndex()
	appName := fetchAppName()
	go func() {
		t := time.NewTicker(time.Second)
		for {
			<-t.C
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to fetch index: %s\n", err.Error())
			} else {
				fmt.Println(fmt.Sprintf("%s. Says %s. on index: %d", appName, message, index))
			}
		}
	}()

	wg := sync.WaitGroup{}
	for _, port := range ports {
		wg.Add(1)
		go func(wg *sync.WaitGroup, port string) {
			defer wg.Done()
			handler, err := rata.NewRouter(routes.Routes, handlers.New(logger, port))
			if err != nil {
				logger.Fatal("router.creation.failed", err)
			}

			server := ifrit.Envoke(http_server.New(":"+port, handler))
			logger.Info("test-app.up", lager.Data{"port": port})
			err = <-server.Wait()
			if err != nil {
				logger.Error("shutting down server", err, lager.Data{"server port": port})
			}
			logger.Info("shutting down server", lager.Data{"server port": port})
		}(&wg, port)
	}
	wg.Wait()
	logger.Info("shutting latice app")
}

func fetchAppName() string {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		return "test-app"
	}
	return appName
}

func getServerPorts() []string {
	givenPorts := *portsFlag
	if givenPorts == "" {
		givenPorts = os.Getenv("PORT")
	}
	if givenPorts == "" {
		givenPorts = "8080"
	}
	ports := strings.Replace(givenPorts, " ", "", -1)
	return strings.Split(ports, ",")
}
