package controller

import (
	"errors"
	"fmt"
	"github.com/censoredgit/light/locker"
	"log/slog"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/flosch/pongo2/v6"
)

const (
	defaultMaxUploadSize  = 20 * 1024
	defaultLoginRouteName = "login"
	backRedirectKey       = "_back_redirect"
	defaultCsrfTokenField = "_csrf_field"
	defaultAuthFieldName  = "_auth"
	defaultStaticPath     = "/static"
)

var config Config

var iLog *slog.Logger
var rootAction *Action
var namedRouterMap = make(map[string]string)

var templateSet *pongo2.TemplateSet
var templateFuncMap = make(map[string]func(args ...any) string)

type Controller struct {
	*Mount
	groups []*Group
	static []*static
	*Container
}

func newController() *Controller {
	c := &Controller{}
	c.Container = &Container{items: make([]any, 0)}
	c.Mount = newMount(c.Container)

	return c
}

func MustSetup(cfg *Config) *Controller {
	if cfg.Logger == nil {
		panic("logger required.")
	}
	iLog = cfg.Logger.With(slog.String("package", "controller"))

	if cfg.SessionManager == nil {
		panic("session driver required")
	}

	config = *cfg

	if config.Templates.RootPath != "" {
		templateSet = pongo2.NewSet("base", pongo2.MustNewLocalFileSystemLoader(config.Templates.RootPath))
	}

	if config.MaxUploadSize == 0 {
		config.MaxUploadSize = defaultMaxUploadSize
	}

	if strings.TrimSpace(config.CsrfFieldName) == "" {
		config.CsrfFieldName = defaultCsrfTokenField
	}

	if strings.TrimSpace(config.LoginRouteName) == "" {
		config.LoginRouteName = defaultLoginRouteName
	}

	if strings.TrimSpace(config.StaticPath) == "" {
		config.StaticPath = defaultStaticPath
	}

	return newController()
}

func (c *Controller) RegisterTemplateFunc(name string, fn func(args ...any) string) {
	templateFuncMap[name] = fn
}

func (c *Controller) Group() *Group {
	g := &Group{
		m: newMount(c.Container),
	}
	g.m.belong = g
	c.groups = append(c.groups, g)
	return g
}

func (c *Controller) AddStatic(uri, targetPath string) {
	c.static = append(
		c.static,
		&static{
			uri:     http.MethodGet + " " + uri,
			handler: http.StripPrefix(uri, http.FileServer(newStaticFileSystem(http.Dir(targetPath)))),
		},
	)
}

func (c *Controller) Serve() error {
	err := c.composeRouters()
	if err != nil {
		return err
	}

	err = c.checkRootAction()
	if err != nil {
		return err
	}

	loginUri, err := c.loginRouteUri()
	if err != nil {
		iLog.Info(err.Error())
	}

	setupCtx(
		defaultAuthFieldName,
		config.CsrfFieldName,
	)
	setupAuthMiddleware(
		loginUri,
		config.LoginRouteName,
		backRedirectKey,
	)
	setupCsrfMiddleware(config.CsrfFieldName)
	setupLockMiddleware(locker.New(&locker.Config{}))

	return http.ListenAndServe(net.JoinHostPort(config.Host, config.Port), nil)
}

func (c *Controller) composeRouters() error {
	mountInfo := make(map[string]*MountInfo)

	for uri, info := range c.info {
		if _, exists := mountInfo[uri]; exists {
			return errors.New("duplicated uri: " + uri)
		}

		mountInfo[uri] = info
	}

	for _, group := range c.groups {
		if err := c.resolveGroups(mountInfo, group); err != nil {
			return err
		}
	}

	for uri, info := range mountInfo {
		if info.action.isReady {
			return errors.New(fmt.Sprintf("action already for other route. [%v]", info.action))
		}

		http.Handle(uri, info.action)

		info.action.middlewares = append(info.action.middlewares, info.action)
		info.routeUri = getRawUri(uri)
		info.action.isReady = true

		if config.Debug {
			iLog.Info(fmt.Sprintf("%s%s%v", "Mounted: ", uri, info.action.middlewaresName()))
		}

		if info.routeName == "" {
			continue
		}

		if _, exists := namedRouterMap[info.routeName]; exists {
			return errors.New("duplicated route name: " + info.routeName)
		}

		namedRouterMap[info.routeName] = getRawUri(uri)
	}

	for _, _static := range c.static {
		http.Handle(_static.uri, _static.handler)
	}

	return nil
}

func (c *Controller) resolveGroups(mountInfo map[string]*MountInfo, group *Group) error {
	if group == nil {
		return nil
	}

	for uri, info := range group.m.info {
		expsUri := strings.SplitN(uri, " ", 2)
		uri = expsUri[0] + " " + path.Join(group.composePrefix(), expsUri[1])

		if _, exists := mountInfo[uri]; exists {
			return errors.New("duplicated uri: " + uri)
		}

		info.action.WithMiddleware(group.composeMiddleware()...)

		mountInfo[uri] = info
	}

	for _, g := range group.m.groups {
		if err := c.resolveGroups(mountInfo, g); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) checkRootAction() error {
	if rootAction == nil {
		return errors.New("should be an one root action")
	}

	for _, info := range c.info {
		if info.action == rootAction && (info.method != http.MethodGet || info.hasParameters) {
			return errors.New("only get method and no parameters are allowed for root action")
		}
	}

	for _, group := range c.groups {
		for _, info := range group.m.info {
			if info.action == rootAction && (info.method != http.MethodGet || group.hasParameters() || info.hasParameters) {
				return errors.New("only get method and no parameters are allowed for root action")
			}
		}
	}

	return nil
}

func (c *Controller) loginRouteUri() (string, error) {
	if existsLoginUri, ok := namedRouterMap[config.LoginRouteName]; ok {
		return existsLoginUri, nil
	}

	return "", errors.New("login uri not found by route name " + config.LoginRouteName)
}

func getRawUri(uri string) string {
	return strings.NewReplacer(
		http.MethodGet+" ", "",
		http.MethodPost+" ", "",
		http.MethodPut+" ", "",
		http.MethodPatch+" ", "",
		http.MethodDelete+" ", "").Replace(uri)
}

func iLogError(msg string) {
	iLog.Error(msg)
}
