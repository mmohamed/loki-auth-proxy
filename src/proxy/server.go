package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/mmohamed/loki-auth-proxy/src/pkg"
	"github.com/urfave/cli"
)

func Serve(c *cli.Context) error {
	lokiServerURL, _ := url.Parse(c.String("loki-server"))
	serveAt := fmt.Sprintf(":%d", c.Int("port"))
	authConfigLocation := c.String("auth-config")
	orgCheck := c.Bool("org-check")
	authConfig, err := pkg.ParseConfig(&authConfigLocation)
	
	log.Printf("Loki multi tenant proxy is starting for %s on port %d ...", c.String("loki-server"), c.Int("port"))

	if authConfig == nil {
		log.Fatalf("Starting failed, unable to load auth-config file : %v", err)
		return err
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	http.HandleFunc("/", createHandler(lokiServerURL, authConfig, orgCheck))
	if err := http.ListenAndServe(serveAt, nil); err != nil {
		log.Fatalf("Loki multi tenant proxy can not start %v", err)
		return err
	}
	return nil
}

func createHandler(lokiServerURL *url.URL, authConfig *pkg.Authn, orgCheck bool) http.HandlerFunc {
	reverseProxy := httputil.NewSingleHostReverseProxy(lokiServerURL)
	return LogRequest(BasicAuth(ReverseLoki(reverseProxy, lokiServerURL), authConfig, orgCheck))
}
