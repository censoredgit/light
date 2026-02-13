package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/censoredgit/light/session"
	"github.com/censoredgit/light/utils"
	"github.com/censoredgit/light/validator/input"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var ErrUnexpectedMiddlewareRun = errors.New("unexpected middleware run")

var ctxCfg struct {
	authFieldName string
	csrfFieldName string
}

func setupCtx(authFieldName, csrfFieldName string) {
	ctxCfg.authFieldName = authFieldName
	ctxCfg.csrfFieldName = csrfFieldName
}

type CtxForm struct {
	Values url.Values
	Files  map[string][]*multipart.FileHeader
}

type ContextFlashStorage interface {
	Errors() *ErrorBag
	Inputs() *InputBag
	Flush()
}

type Ctx struct {
	context.Context
	request               *http.Request
	httpResponse          http.ResponseWriter
	session               *session.Data
	log                   *slog.Logger
	requestValidatorInput *input.Data
	role                  RoleSupport
	routeName             string
	routeUri              string
	form                  *CtxForm
	extras                map[string]string
	middlewareStackIndex  int8
	action                *Action
	flashStorage          ContextFlashStorage
}

func (ctx *Ctx) Next() (Response, error) {
	if len(ctx.action.middlewares) == 0 {
		ctx.middlewareStackIndex++
		if ctx.middlewareStackIndex > 0 {
			return nil, ErrUnexpectedMiddlewareRun
		}

		return ctx.action.handler(ctx)
	}

	index := int(ctx.middlewareStackIndex)
	if index > len(ctx.action.middlewares) {
		return nil, ErrUnexpectedMiddlewareRun
	} else if index == len(ctx.action.middlewares) {
		ctx.middlewareStackIndex++

		return ctx.action.handler(ctx)
	}

	ctx.middlewareStackIndex++
	return ctx.action.middlewares[index+1].Next(ctx)
}

type AuthIdentification interface {
	AuthId() string
	IsActive() bool
}

type RoleSupport interface {
	Role(is string) bool
	Can(action string) bool
}

func (ctx *Ctx) Request() *http.Request {
	return ctx.request
}

func (ctx *Ctx) SetRoleSupport(rs RoleSupport) {
	ctx.role = rs
}

func (ctx *Ctx) RedirectResponse(url string, args ...interface{}) *RedirectResponse {
	if newUrl, found := strings.CutPrefix(url, ":"); found {
		url = ctx.Route(newUrl, args...)
	}
	return newRedirect(url, ctx.flashStorage)
}

func (ctx *Ctx) BackRedirectResponse() *BackRedirectResponse {
	return newBackRedirectResponse(ctx.request, ctx.flashStorage)
}

func (ctx *Ctx) HtmlResponse(view string, code int) *HtmlResponse {
	return newHtmlResponse(view, code, ctx.flashStorage)
}

func (ctx *Ctx) TemplateResponse(view string) *TemplateResponse {
	return newTemplateResponse(view, ctx.flashStorage)
}

func (ctx *Ctx) TemplateInlineResponse(data string) *TemplateInlineResponse {
	return newTemplateInlineResponse(data, ctx.flashStorage)
}

func (ctx *Ctx) TextResponse(data string, code int) *TextResponse {
	return newTextResponse(data, code, ctx.flashStorage)
}

func (ctx *Ctx) CodeResponse(code int) *CodeResponse {
	return newCodeResponse(code, ctx.flashStorage)
}

func (ctx *Ctx) JsonResponse(data interface{}, code int) *JsonResponse {
	return newJsonResponse(data, code)
}

func (ctx *Ctx) Session() *session.Data {
	return ctx.session
}

func (ctx *Ctx) Log() *slog.Logger {
	return ctx.log
}

func (ctx *Ctx) Login(id AuthIdentification) {
	ctx.session.Set(ctxCfg.authFieldName, id.AuthId())
}

func (ctx *Ctx) AuthIdentification() string {
	return ctx.session.Get(ctxCfg.authFieldName)
}

func (ctx *Ctx) IsAuth() bool {
	return ctx.session.Has(ctxCfg.authFieldName)
}

func (ctx *Ctx) IsGuest() bool {
	return !ctx.session.Has(ctxCfg.authFieldName)
}

func (ctx *Ctx) Logout() {
	ctx.session.Empty()
}

func (ctx *Ctx) RequestValidatorInput() *input.Data {
	return ctx.requestValidatorInput
}

func (ctx *Ctx) Role(role string) bool {
	if ctx.role == nil {
		return false
	}
	return ctx.role.Role(role)
}

func (ctx *Ctx) Can(permission string) bool {
	if ctx.role == nil {
		return false
	}
	return ctx.role.Can(permission)
}

func (ctx *Ctx) IntPathValue(name string) int64 {
	value, err := strconv.Atoi(ctx.request.PathValue(name))
	if err != nil {
		return 0
	}
	return int64(value)
}

func (ctx *Ctx) PathValue(name string) string {
	return ctx.request.PathValue(name)
}

func (ctx *Ctx) Route(name string, args ...interface{}) string {
	if uri, exists := namedRouterMap[name]; exists {
		result, err := ctx.composeUri(uri, args)
		if err != nil {
			iLog.Warn(err.Error())
			return *rootAction.uri
		}
		return result
	}
	iLog.Warn("Route not found: " + name)
	return *rootAction.uri
}

func (ctx *Ctx) Url(routeName string, args ...interface{}) string {
	uri := ctx.Route(routeName, args...)

	if config.ExternalHost != "" {
		return fmt.Sprintf("%s%s", config.ExternalHost, uri)
	}

	return fmt.Sprintf("%s://%s:%s%s", config.Protocol, config.Host, config.Port, uri)
}

func (ctx *Ctx) InternalUrl(routeName string, args ...interface{}) string {
	uri := ctx.Route(routeName, args...)

	if config.InternalHost != "" {
		return fmt.Sprintf("%s%s", config.InternalHost, uri)
	}

	return fmt.Sprintf("%s%s:%s%s", config.Protocol, config.Host, config.Port, uri)
}

func (ctx *Ctx) IsCurrentRoute(name string) bool {
	return ctx.routeName == name
}

func (ctx *Ctx) HasCurrentRoutePrefix(name string, args ...interface{}) bool {
	haystack := strings.TrimPrefix(ctx.Route(name, args...), "/")
	needle := strings.TrimPrefix(ctx.routeUri, "/")
	return strings.HasPrefix(needle, haystack)
}

func (ctx *Ctx) Form() *CtxForm {
	return ctx.form
}

func (ctx *Ctx) CsrfToken() string {
	if !ctx.session.Has(ctxCfg.csrfFieldName) {
		ctx.session.Set(ctxCfg.csrfFieldName, utils.UUID())
	}

	return ctx.session.Get(ctxCfg.csrfFieldName)
}

func (ctx *Ctx) CsrfFieldName() string {
	return ctxCfg.csrfFieldName
}

func (ctx *Ctx) IsJson() bool {
	if ctx.request == nil {
		return false
	}

	return ctx.request.Header.Get("Content-Type") == "application/json"
}

func (ctx *Ctx) composeUri(uri string, args []interface{}) (string, error) {
	braceCount := strings.Count(uri, "{")
	if braceCount == 0 {
		return uri, nil
	}

	if braceCount > len(args) {
		return "", fmt.Errorf("uri %s parameters more than number of args %d", uri, len(args))
	}

	pattern := regexp.MustCompile("{(.*?)}")
	allMatches := pattern.FindAllString(uri, braceCount)

	var paramValue string
	for index := range braceCount {
		paramValue = fmt.Sprint(args[index])
		uri = strings.Replace(uri, allMatches[index], paramValue, 1)
	}

	return uri, nil
}

func (ctx *Ctx) IsPost() bool {
	return ctx.request.Method == http.MethodPost
}

func (ctx *Ctx) IsPut() bool {
	return ctx.request.Method == http.MethodPut
}

func (ctx *Ctx) IsGet() bool {
	return ctx.request.Method == http.MethodGet
}

func (ctx *Ctx) IsDelete() bool {
	return ctx.request.Method == http.MethodDelete
}

func (ctx *Ctx) IsPatch() bool {
	return ctx.request.Method == http.MethodPatch
}

func (ctx *Ctx) parseForm() error {
	var err error

	if ctx.IsPost() || ctx.IsPut() {
		err = ctx.request.ParseMultipartForm(config.MaxUploadSize)
		if err == nil {
			ctx.form.Files = ctx.request.MultipartForm.File
			ctx.form.Values = ctx.request.MultipartForm.Value
		} else {
			err = ctx.request.ParseForm()
			if err != nil {
				return err
			}

			ctx.form.Values = ctx.request.Form
		}
	} else {
		err = ctx.request.ParseForm()
		if err != nil {
			return err
		}
		ctx.form.Values = ctx.request.Form
	}

	return nil
}

func (ctx *Ctx) ErrBag() *ErrorBag {
	return ctx.flashStorage.Errors()
}
func (ctx *Ctx) InputBag() *InputBag {
	return ctx.flashStorage.Inputs()
}
