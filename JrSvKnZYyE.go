package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	httprouter "github.com/julienschmidt/httprouter"
	yamlLib "gopkg.in/yaml.v2"
)

type (
	Config struct {
		Graphs []map[string]interface{}
	}

	Yaml Config

	justFilesFilesystem struct {
		fs http.FileSystem
	}
	neuteredReaddirFile struct {
		http.File
	}

	httprouterReturn func(w http.ResponseWriter, r *http.Request, param httprouter.Params)
)

const (
	//	STATIC_FOLDER = "/opt/simmetrica/static"
	//	TEMPLATE_FOLDER = "/opt/simmetrica/templates"
	//	DEFAULT_CONFIG_FILE = "/opt/simmetrica/config/config.yml"

	STATIC_FOLDER       = "/var/golang/src/github.com/feyyazesat/simmetrica/static"
	TEMPLATE_FOLDER     = "/var/golang/src/github.com/feyyazesat/simmetrica/templates"
	DEFAULT_CONFIG_FILE = "/var/golang/src/github.com/feyyazesat/simmetrica/config/config.yml"
)

var (
	err error

	debug          bool
	config         string
	redis_host     string
	redis_port     string
	redis_db       string
	redis_password string

	yaml Yaml

	router *httprouter.Router
)

func Check(err error) {
	if err != nil {
		panic(err)
		os.Exit(1)
	}
}

func LogRoutes(fnHandler httprouterReturn, name string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

		start := time.Now()

		fnHandler(w, r, param)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	}
}

func (yaml *Yaml) UnmarshalYAML(configContent *[]byte) error {

	return yamlLib.Unmarshal(*configContent, yaml)

}

func ReadParametersFile(configPath string) (*[]byte, error) {
	configContent, err := ioutil.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	return &configContent, nil
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	return neuteredReaddirFile{f}, nil
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome It's Index!\n")
}

func Push(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, params.ByName("event"))
}

func Query(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, params.ByName("event"), params.ByName("start"), params.ByName("end"))
}

func Graph(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome It's Graph!\n")
}

func CreateRoutes(router *httprouter.Router) *httprouter.Router {
	router.GET("/", LogRoutes(Index, "Index"))
	router.GET("/push/:event", LogRoutes(Push, "Push"))
	router.GET("/query/:event/:start/:end", LogRoutes(Query, "Query"))
	router.GET("/graph", LogRoutes(Graph, "Graph"))

	fs := justFilesFilesystem{http.Dir(STATIC_FOLDER)}
	router.Handler("GET", "/static/*filepath", http.StripPrefix("/static", http.FileServer(fs)))

	return router
}

func init() {
	flag.BoolVar(
		&debug,
		"debug",
		false,
		"Run the app in debug mode")

	flag.StringVar(
		&config,
		"config",
		DEFAULT_CONFIG_FILE,
		fmt.Sprintf("Run with the specified config file (default: %s)", DEFAULT_CONFIG_FILE))

	flag.StringVar(
		&redis_host,
		"redis_host",
		"",
		"Connect to redis on the specified host")

	flag.StringVar(
		&redis_port,
		"redis_port",
		"",
		"Connect to redis on the specified port")

	flag.StringVar(
		&redis_db,
		"redis_db",
		"",
		"Connect to the specified db in redis")

	flag.StringVar(
		&redis_password,
		"redis_password",
		"",
		"Authorization password of redis")

	flag.Parse()
}

func main() {
	{
		//scope to unset configContent
		var configContent *[]byte
		configContent, err = ReadParametersFile(config)
		Check(err)

		err = yaml.UnmarshalYAML(configContent)
		Check(err)
	}
	router = CreateRoutes(httprouter.New())

	log.Fatal(http.ListenAndServe(":8080", router))
}
