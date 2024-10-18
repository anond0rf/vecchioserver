package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/anond0rf/vecchioserver/internal/handlers"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const banner = `
                     _     _       __                          
 /\   /\___  ___ ___| |__ (_) ___ / _\ ___ _ ____   _____ _ __ 
 \ \ / / _ \/ __/ __| '_ \| |/ _ \\ \ / _ \ '__\ \ / / _ \ '__|
  \ V /  __/ (_| (__| | | | | (_) |\ \  __/ |   \ V /  __/ |   
   \_/ \___|\___\___|_| |_|_|\___/\__/\___|_|    \_/ \___|_|    

               Post on vecchiochan through REST                

                            v1.0.0                              
`

func main() {
	var port int
	var userAgent string
	var verbose bool
	flag.IntVar(&port, "p", 8080, "Port to run the server on")
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.StringVar(&userAgent, "u", "", "User-Agent header to be used by the internal client")
	flag.StringVar(&userAgent, "user-agent", "", "User-Agent header to be used by the internal client")
	flag.BoolVar(&verbose, "v", false, "Enables verbose logging")
	flag.BoolVar(&verbose, "verbose", false, "Enables verbose logging")
	flag.Parse()

	e := echo.New()
	if verbose {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.HideBanner = true
	fmt.Printf("%s\nAPI documentation available at http://localhost:%d/swagger/index.html\n", banner, port)

	addServer := getServerAdder(port)

	h := handlers.NewAPIHandler(userAgent, verbose)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/swagger", redirectToIndex)
	e.GET("/swagger//index.html", redirectToIndex)
	e.GET("/swagger/doc.json", addServer)
	e.GET("/swagger/doc.yaml", addServer)
	e.POST("/thread", h.NewThread)
	e.POST("/reply", h.PostReply)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server: %q", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func redirectToIndex(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}

func getServerAdder(port int) func(c echo.Context) error {
	return func(c echo.Context) error {
		swagger, err := handlers.GetSwagger()
		if err != nil {
			return err
		}

		swagger.Servers = openapi3.Servers{
			&openapi3.Server{
				URL: fmt.Sprintf("http://localhost:%s", strconv.Itoa(port)),
			},
		}

		jsonData, err := json.Marshal(swagger)
		if err != nil {
			return err
		}

		return c.JSONBlob(http.StatusOK, jsonData)
	}
}
