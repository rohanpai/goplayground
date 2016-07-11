package main

import (
	&#34;flag&#34;
	&#34;fmt&#34;
	&#34;io/ioutil&#34;
	&#34;log&#34;
	&#34;net/http&#34;
	&#34;os&#34;
	&#34;time&#34;

	httprouter &#34;github.com/julienschmidt/httprouter&#34;
	yamlLib &#34;gopkg.in/yaml.v2&#34;
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
	//	STATIC_FOLDER = &#34;/opt/simmetrica/static&#34;
	//	TEMPLATE_FOLDER = &#34;/opt/simmetrica/templates&#34;
	//	DEFAULT_CONFIG_FILE = &#34;/opt/simmetrica/config/config.yml&#34;

	STATIC_FOLDER       = &#34;/var/golang/src/github.com/feyyazesat/simmetrica/static&#34;
	TEMPLATE_FOLDER     = &#34;/var/golang/src/github.com/feyyazesat/simmetrica/templates&#34;
	DEFAULT_CONFIG_FILE = &#34;/var/golang/src/github.com/feyyazesat/simmetrica/config/config.yml&#34;
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
			&#34;%s\t%s\t%s\t%s&#34;,
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

	return &amp;configContent, nil
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
	fmt.Fprint(w, &#34;Welcome It&#39;s Index!\n&#34;)
}

func Push(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, params.ByName(&#34;event&#34;))
}

func Query(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, params.ByName(&#34;event&#34;), params.ByName(&#34;start&#34;), params.ByName(&#34;end&#34;))
}

func Graph(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, &#34;Welcome It&#39;s Graph!\n&#34;)
}

func CreateRoutes(router *httprouter.Router) *httprouter.Router {
	router.GET(&#34;/&#34;, LogRoutes(Index, &#34;Index&#34;))
	router.GET(&#34;/push/:event&#34;, LogRoutes(Push, &#34;Push&#34;))
	router.GET(&#34;/query/:event/:start/:end&#34;, LogRoutes(Query, &#34;Query&#34;))
	router.GET(&#34;/graph&#34;, LogRoutes(Graph, &#34;Graph&#34;))

	fs := justFilesFilesystem{http.Dir(STATIC_FOLDER)}
	router.Handler(&#34;GET&#34;, &#34;/static/*filepath&#34;, http.StripPrefix(&#34;/static&#34;, http.FileServer(fs)))

	return router
}

func init() {
	flag.BoolVar(
		&amp;debug,
		&#34;debug&#34;,
		false,
		&#34;Run the app in debug mode&#34;)

	flag.StringVar(
		&amp;config,
		&#34;config&#34;,
		DEFAULT_CONFIG_FILE,
		fmt.Sprintf(&#34;Run with the specified config file (default: %s)&#34;, DEFAULT_CONFIG_FILE))

	flag.StringVar(
		&amp;redis_host,
		&#34;redis_host&#34;,
		&#34;&#34;,
		&#34;Connect to redis on the specified host&#34;)

	flag.StringVar(
		&amp;redis_port,
		&#34;redis_port&#34;,
		&#34;&#34;,
		&#34;Connect to redis on the specified port&#34;)

	flag.StringVar(
		&amp;redis_db,
		&#34;redis_db&#34;,
		&#34;&#34;,
		&#34;Connect to the specified db in redis&#34;)

	flag.StringVar(
		&amp;redis_password,
		&#34;redis_password&#34;,
		&#34;&#34;,
		&#34;Authorization password of redis&#34;)

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

	log.Fatal(http.ListenAndServe(&#34;:8080&#34;, router))
}
