package web

import (
	"fmt"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/coocood/freecache"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/sentry"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	. "github.com/tj/go-debug"

	"github.com/liut/staffio/pkg/backends"
	"github.com/liut/staffio/pkg/settings"
)

var (
	cache *freecache.Cache
	svr   *server
	debug = Debug("staffio:web")
)

func IsKeeper(uid string) bool {
	return Default().service.InGroup("keeper", uid)
}

type server struct {
	*gin.Engine
	service backends.Servicer
	osvr    *osin.Server
}

func (s *server) IsKeeper(uid string) bool {
	return s.InGroup("keeper", uid)
}

func (s *server) InGroup(gn, uid string) bool {
	return s.service.InGroup(gn, uid)
}

// Default returns current server instance
func Default() *server {
	if svr != nil {
		return svr
	}
	service := backends.NewService()
	osvr := osin.NewServer(newOsinConfig(), service.OSIN())
	var err error
	osvr.AccessTokenGen, err = getTokenGenJWT()
	if err != nil {
		panic(err)
	}

	svr = &server{
		Engine:  gin.New(),
		service: service,
		osvr:    osvr,
	}

	if settings.IsDevelop() {
		fmt.Printf("In Developing(Debug) mode, gin: %s\n", gin.Mode())
		svr.Use(gin.Logger())
		svr.Use(gin.Recovery())
	} else {
		fmt.Printf("In Release mode, gin: %s\n", gin.Mode())
		if settings.SentryDSN != "" {
			raven.SetDSN(settings.SentryDSN)
			onlyCrashes := false
			svr.Use(sentry.Recovery(raven.DefaultClient, onlyCrashes))
		}
	}

	store := sessionStore()
	svr.Use(sessions.Sessions("session", store))
	svr.StrapRouter(svr.Engine)

	cache = freecache.NewCache(settings.CacheSize)

	return svr
}

func (s *server) HandleFunc(path string, hf http.HandlerFunc) {
	h := func(c *gin.Context) {
		hf(c.Writer, c.Request)
	}
	s.GET(path, h)
	s.POST(path, h)
}

func (s *server) Run(addr string) error {
	return http.ListenAndServe(addr, s.Engine)
}

func newOsinConfig() *osin.ServerConfig {
	return &osin.ServerConfig{
		AuthorizationExpiration: 900,
		AccessExpiration:        3600 * 24,
		TokenType:               "bearer",
		AllowedAuthorizeTypes: osin.AllowedAuthorizeType{
			osin.CODE,
			osin.TOKEN,
		},
		AllowedAccessTypes: osin.AllowedAccessType{
			osin.AUTHORIZATION_CODE,
			osin.IMPLICIT,
			// osin.REFRESH_TOKEN,
			osin.PASSWORD,
			// osin.CLIENT_CREDENTIALS,
		},
		ErrorStatusCode:           200,
		AllowClientSecretInParams: true,
		AllowGetAccessRequest:     false,
	}
}
