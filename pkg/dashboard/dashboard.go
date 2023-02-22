package dashboard

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/template"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/sudermanjr/led-controller/pkg/color"
	"github.com/sudermanjr/led-controller/pkg/neopixel"
	"github.com/sudermanjr/led-controller/pkg/screen"
	"github.com/sudermanjr/led-controller/pkg/utils"
)

var (
	//go:embed templates/*
	templates embed.FS

	//go:embed static
	static embed.FS
)

// App encapsulates all the config for the server
type App struct {
	Router    *chi.Mux
	Port      int
	Array     *neopixel.LEDArray
	Screen    *screen.Display
	ButtonPin int64
	Logger    *zap.SugaredLogger
}

func getBaseTemplate() (*template.Template, error) {
	tmpl := template.New("main").Funcs(template.FuncMap{
		"getUUID": getUUID,
	})

	templateFileNames := []string{
		"main.gohtml",
		"head.gohtml",
		"navbar.gohtml",
		"dashboard.gohtml",
		"footer.gohtml",
	}
	return parseTemplateFiles(tmpl, templateFileNames)
}

func parseTemplateFiles(tmpl *template.Template, templateFileNames []string) (*template.Template, error) {
	for _, fname := range templateFileNames {
		templateFile, err := templates.ReadFile("templates/" + fname)
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.Parse(string(templateFile))
		if err != nil {
			return nil, err
		}
	}
	return tmpl, nil
}

func (a *App) writeTemplate(tmpl *template.Template, data string, w http.ResponseWriter) {
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, data)
	if err != nil {
		a.Logger.Errorw("error executing template", "error", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		a.Logger.Errorw("error writing template", "error", err)
	}
}

// Initialize sets up an instance of App
func (a *App) Initialize() {
	a.Router = chi.NewRouter()
	a.Router.Use(middleware.Recoverer)
	a.Router.Use(LoggingMiddleware(a.Logger))

	//API
	a.Router.MethodFunc("GET", "/health", a.healthHandler)
	a.Router.MethodFunc("POST", "/control", a.controlHandler)
	a.Router.MethodFunc("POST", "/demo", a.demoHandler)
	a.Router.MethodFunc("POST", "/button/power", a.powerButtonHandler)
	a.Router.MethodFunc("GET", "/", a.rootHandler)

	// Static Files
	a.Router.Handle("/static/*", http.FileServer(http.FS(static)))

	if a.Screen != nil {
		// Display Info On Screen
		err := a.Screen.InfoDisplay()
		if err != nil {
			a.Logger.Errorw("error displaying screen", "error", err)
		}
	}

	a.ButtonPin = 4 // TODO: this probably should be more dynamic
}

func (a *App) Run() error {
	a.Logger.Infow("starting server", "port", a.Port)
	go a.WatchButton()
	defer a.Array.WS.Fini()

	// https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
	// The HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", a.Port),
		Handler: a.Router,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				a.Logger.Fatalw("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			a.Logger.Errorw("error shutting down", "error", err)
			return
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	return nil
}

// rootHandler gets template data and renders the dashboard with it.
func (a *App) rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := getBaseTemplate()
	if err != nil {
		a.Logger.Errorw("error getting template data", "error", err)
		http.Error(w, "Error getting template data", 500)
		return
	}
	a.writeTemplate(tmpl, "{}", w)
}

// health is a healthcheck endpoint
func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("healthy"))
	if err != nil {
		a.Logger.Errorw("error writing healthcheck", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) controlHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	colorValue := color.HexToColor(r.Form["color"][0])
	brightness, err := strconv.ParseInt(r.Form["brightness"][0], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.Array.Color = colorValue
	a.Array.Brightness = utils.ScaleBrightness(int(brightness), a.Array.MinBrightness, a.Array.MaxBrightness)
	err = a.Array.Display(0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 302)
}

func (a *App) demoHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delay, err := strconv.ParseInt(r.Form["delay"][0], 10, 32)
	if err != nil {
		a.Logger.Errorw("error processing delay", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	brightness, err := strconv.ParseInt(r.Form["brightness"][0], 10, 32)
	if err != nil {
		a.Logger.Errorw("error processing brightness", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	steps, err := strconv.ParseInt(r.Form["gradient-steps"][0], 10, 32)
	if err != nil {
		a.Logger.Errorw("error processing gradient stpes", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.Array.Brightness = utils.ScaleBrightness(int(brightness), a.Array.MinBrightness, a.Array.MaxBrightness)
	a.Array.Demo(1, int(delay), int(steps))

	http.Redirect(w, r, "/", 302)
}
