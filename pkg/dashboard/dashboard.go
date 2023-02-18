package dashboard

import (
	"bytes"
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

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

	//go:embed assets
	assets embed.FS
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
	router := chi.NewRouter()

	//API
	router.Get("/health", a.health)
	router.Post("/control", a.control)
	router.Post("/demo", a.demo)

	// HTML Dashboard
	fileServer := http.FileServer(http.FS(assets))
	router.Handle("/static/*", fileServer)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		a.rootHandler(w, r)
	})

	if a.Screen != nil {
		// Display Info On Screen
		err := a.Screen.InfoDisplay()
		if err != nil {
			a.Logger.Errorw("error displaying screen", "error", err)
		}
	}

	a.ButtonPin = 4

	a.Router = router
}

// Run starts the http server
func (a *App) Run() {
	a.Logger.Infow("starting server", "port", a.Port)
	go a.WatchButton()
	defer a.Array.WS.Fini()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.Port), nil); err != nil {
		a.Logger.Fatalw("failed to start server", "error", err)
	}
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
func (a *App) health(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("healthy"))
	if err != nil {
		a.Logger.Errorw("error writing healthcheck", "error", err)
	}
}

func (a *App) control(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	colorValue := color.HexToColor(r.Form["color"][0])
	brightness, err := strconv.ParseInt(r.Form["brightness"][0], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	a.Array.Color = colorValue
	a.Array.Brightness = utils.ScaleBrightness(int(brightness), a.Array.MinBrightness, a.Array.MaxBrightness)
	err = a.Array.Display(0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", 302)
}

func (a *App) demo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	delay, err := strconv.ParseInt(r.Form["delay"][0], 10, 32)
	if err != nil {
		a.Logger.Errorw("error processing delay", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	brightness, err := strconv.ParseInt(r.Form["brightness"][0], 10, 32)
	if err != nil {
		a.Logger.Errorw("error processing brightness", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	a.Array.Brightness = utils.ScaleBrightness(int(brightness), a.Array.MinBrightness, a.Array.MaxBrightness)
	a.Array.Demo(1, int(delay), 1000)

	http.Redirect(w, r, "/", 302)
}
