package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"time"
	"zip_service/cmd/main/docs"
	"zip_service/internal/config"
	"zip_service/internal/handler"
	"zip_service/internal/service"
)

type App struct {
	router     *httprouter.Router
	httpServer *http.Server
	cfg        *config.Config
}

func NewApp(cfg *config.Config) *App {
	router := httprouter.New()

	log.Println("swagger docs initializing")
	initSwagger(router, cfg.HTTPConfig.Host, cfg.HTTPConfig.Port)

	log.Println("zip routes initializing")
	zh := handler.NewZipHandler(service.NewZipStreamer(cfg))
	zh.Register(router)

	router.NotFound = http.FileServer(http.Dir("root"))

	return &App{
		router: router,
		cfg:    cfg,
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
	host := fmt.Sprintf("%s:%s", a.cfg.HTTPConfig.Host, a.cfg.HTTPConfig.Port)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal("failed to create listener")
	}

	a.httpServer = &http.Server{
		Handler:      a.router,
		WriteTimeout: time.Duration(a.cfg.HTTPConfig.SendTimeout) * time.Second,
		ReadTimeout:  time.Duration(a.cfg.HTTPConfig.ReadTimeout) * time.Second,
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

func initSwagger(router *httprouter.Router, ip, port string) {
	host := fmt.Sprintf("%s:%s", ip, port)
	docs.SwaggerInfo.Host = host
	router.Handler(http.MethodGet, "/swagger/*filename", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", host)),
	))
}
