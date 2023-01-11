package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"time"
	"zip_service/internal/handler"
	"zip_service/internal/service"
)

type App struct {
	router     *httprouter.Router
	httpServer *http.Server
}

func NewApp() *App {
	router := httprouter.New()

	log.Println("zip routes initializing")
	zh := handler.NewZipHandler(service.NewZipStreamer())
	zh.Register(router)

	router.NotFound = http.FileServer(http.Dir("root"))

	return &App{
		router: router,
	}
}

func (a *App) Run(ctx context.Context) error {

	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP(ctx)
	})
	err := grp.Wait()

	return err
}

func (a *App) startHTTP(ctx context.Context) error {
	log.Print("HTTP Server initializing")

	//host := fmt.Sprintf("%s:%s", a.cfg.LocalIP, a.cfg.HTTPConfig.Port)
	// TODO: config
	host := fmt.Sprintf("%s:%s", "0.0.0.0", "8080")
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal("failed to create listener")
	}

	// TODO: find timeouts
	a.httpServer = &http.Server{
		Handler: a.router,
		//WriteTimeout: 10 * time.Second,
		//ReadTimeout:  10 * time.Second,
	}

	go func() {
		if err = a.httpServer.Serve(listener); err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed):
				log.Print("server shutdown")
			default:
				log.Fatal(err)
			}
		}
	}()

	// graceful shutdown
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if shutdownErr := a.httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
		log.Printf("error shutting down server %s", shutdownErr)
	} else {
		log.Print("server shutdown gracefully", shutdownErr)
	}

	return err
}
