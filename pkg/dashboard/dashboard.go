package dashboard

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"k8s.io/klog"

	"github.com/sudermanjr/led-controller/pkg/color"
	"github.com/sudermanjr/led-controller/pkg/neopixel"
)

// App encapsulates all the config for the server
type App struct {
	Router *mux.Router
	Port   int
	Array  *neopixel.LEDArray
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
		templateFile, err := pkger.Open("/pkg/dashboard/templates/" + fname)
		if err != nil {
			return nil, err
		}
		defer templateFile.Close()

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(templateFile)
		if err != nil {
			klog.Error(err)
		}
		s := buf.String()

		tmpl, err = tmpl.Parse(string(s))
		if err != nil {
			return nil, err
		}
	}
	return tmpl, nil
}

func writeTemplate(tmpl *template.Template, data string, w http.ResponseWriter) {
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, data)
	if err != nil {
		klog.Errorf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		klog.Errorf("Error writing template: %v", err)
	}
}

// Initialize sets up an instance of App
func (a *App) Initialize() {
	router := mux.NewRouter()
	router.NotFoundHandler = Handle404()

	//API
	router.HandleFunc("/health", a.health).Methods("GET")
	router.HandleFunc("/control", a.control).Methods("POST")
	router.HandleFunc("/demo", a.demo).Methods("POST")

	// HTML Dashboard
	fileServer := http.FileServer(pkger.Dir("/pkg/dashboard/assets"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		rootHandler(w, r)
	})

	a.Router = router
}

// Run starts the http server
func (a *App) Run() {
	http.Handle("/", a.Router)
	klog.Infof("Starting dashboard server on port %d", a.Port)
	defer a.Array.WS.Fini()
	klog.Fatalf("%v", http.ListenAndServe(fmt.Sprintf(":%d", a.Port), nil))
}

// Handle404 handles the not found error and logs it
func Handle404() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		klog.V(8).Info(r)
		http.Error(w, "Not Found", http.StatusNotFound)
	})
}

// rootHandler gets template data and renders the dashboard with it.
func rootHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := getBaseTemplate()
	if err != nil {
		klog.Errorf("Error getting template data %v", err)
		http.Error(w, "Error getting template data", 500)
		return
	}
	writeTemplate(tmpl, "{}", w)
}

// health is a healthcheck endpoint
func (a *App) health(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("healthy"))
	if err != nil {
		klog.Errorf("Error writing healthcheck: %v", err)
	}
}

func (a *App) control(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	colorValue := color.HexToColor(color.ColorMap[r.Form["color"][0]])
	brightness, err := strconv.ParseInt(r.Form["brightness"][0], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	a.Array.Color = colorValue
	a.Array.Brightness = int(brightness)
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
		klog.Errorf("error processing delay: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	brightness, err := strconv.ParseInt(r.Form["brightness"][0], 10, 32)
	if err != nil {
		klog.Errorf("error processing brightness: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	a.Array.Brightness = int(brightness)
	a.Array.Demo(1, int(delay), 1000)

	http.Redirect(w, r, "/", 302)
}
