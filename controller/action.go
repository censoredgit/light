package controller

import (
	"cmp"
	"context"
	"fmt"
	"github.com/censoredgit/light/validator"
	"log/slog"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"
)

type ActionHandler func(ctx *Ctx) (Response, error)

type Action struct {
	handler     ActionHandler
	validator   *RequestValidator
	middlewares []Middleware
	name        *string
	uri         *string
	isReady     bool
}

func (a *Action) Next(ctx *Ctx) (Response, error) {
	var response Response
	var err error

	if a.validator != nil {
		reqValidator := validator.New()
		reqValidator.SetRuleCollection(a.validator.RuleCollection())
		isValid := false
		if ctx.IsJson() {
			isValid = reqValidator.ValidateByJsonRequest(ctx.request)
		} else {
			isValid = reqValidator.
				ValidateByRequestForms(&ctx.request.Form, ctx.request.MultipartForm)
		}

		if !isValid {
			if !ctx.IsJson() {
				response = ctx.BackRedirectResponse()
			} else {
				response = ctx.JsonResponse(a.validator.JsonResponseBody(), a.validator.JsonResponseCode())
			}
			response.With(func(response ResponseExtendData) {
				response.SetOldInput(ctx.request.Form)
				response.SetErrors(reqValidator.Errors())
			})

			response.ExceptInputs(a.validator.ProtectedFields()...)
		} else {
			ctx.requestValidatorInput = reqValidator.InputData()
		}
	}

	if response == nil {
		response, err = a.handler(ctx)
	}

	return response, err
}

func (a *Action) Priority() uint {
	return math.MaxInt
}

func NewAction(handler ActionHandler) *Action {
	return &Action{handler: handler}
}

func NewRootAction(handler ActionHandler) *Action {
	if rootAction != nil {
		panic("should be one root action")
	}
	rootAction = &Action{handler: handler}

	return rootAction
}

func (a *Action) WithMiddleware(m ...Middleware) *Action {
	for _, newMiddleware := range m {
		if !slices.Contains(a.middlewares, newMiddleware) {
			a.middlewares = append(a.middlewares, newMiddleware)
		}
	}

	slices.SortFunc(a.middlewares, func(a, b Middleware) int {
		return cmp.Compare(int64(a.Priority())-math.MaxInt, int64(b.Priority())-math.MaxInt)
	})

	return a
}

func (a *Action) SetValidator(rv *RequestValidator) *Action {
	a.validator = rv

	return a
}

func (a *Action) middlewaresName() []string {
	names := make([]string, len(a.middlewares))
	for i, m := range a.middlewares {
		mNameChunks := strings.Split(reflect.TypeOf(m).String(), ".")
		names[i] = mNameChunks[len(mNameChunks)-1]
	}

	return names
}

func (a *Action) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var response Response
	var err error

	ctx := &Ctx{
		Context:      context.Background(),
		request:      req,
		httpResponse: res,
		log:          iLog.With(slog.String("Url", req.RequestURI)),
		routeName:    *a.name,
		routeUri:     *a.uri,
		form: &CtxForm{
			Values: url.Values{},
			Files:  make(map[string][]*multipart.FileHeader),
		},
		extras:               make(map[string]string),
		action:               a,
		middlewareStackIndex: -1,
	}

	if config.Debug {
		defer func() {
			if err := req.ParseForm(); err == nil {
				ctx.Log().Info("request: ", slog.AnyValue(req.Form).String(), "")
			}
		}()
	}

	ctx.session, err = config.SessionManager.InitByRequest(ctx.request)
	if err != nil {
		handleError(ctx, err)
		return
	}
	defer func() {
		err := config.SessionManager.Close(ctx.Session())
		if err != nil {
			iLogError(fmt.Errorf("close session err: %w", err).Error())
		}
	}()

	ctx.flashStorage = NewContextSessionFlashStorage(ctx.session)
	defer ctx.flashStorage.Flush()

	handleUser(ctx)

	err = ctx.parseForm()
	if err != nil {
		handleError(ctx, err)
		return
	}

	// resolve /* uri
	if req.URL.Path != "/" && *a.uri == "/" {
		if config.Templates.Page404 != "" {
			rsp := ctx.TemplateResponse(config.Templates.Page404)
			rsp.SetCode(http.StatusNotFound)
			response = rsp
		} else {
			response = ctx.CodeResponse(http.StatusNotFound)
		}
	}

	if response == nil {
		response, err = ctx.Next()
		if err != nil {
			handleError(ctx, err)
			return
		}
	}

	if response == nil {
		iLogError("empty response")
		response = ctx.CodeResponse(http.StatusInternalServerError)
	}

	handleBackAfterAuth(ctx, response)
	handleCookie(ctx)
	response.Process(ctx)
}

func handleCookie(ctx *Ctx) {
	if ctx.session.IsNew() {
		http.SetCookie(ctx.httpResponse, config.SessionManager.ToCookie(ctx.session))
	}
}

func handleBackAfterAuth(ctx *Ctx, response Response) {
	if ctx.IsAuth() && ctx.Session().Has(backRedirectKey) {
		if uri := ctx.Session().Get(backRedirectKey); uri != "" {
			response = ctx.RedirectResponse(uri)
		}
		ctx.Session().Delete(backRedirectKey)
	}
}

func handleUser(ctx *Ctx) {
	if config.UserProvider == nil {
		return
	}

	if !ctx.IsAuth() {
		return
	}

	user, err := config.UserProvider.GetAuthIdentification(ctx, ctx.AuthIdentification())
	if err != nil {
		iLogError(err.Error())
		ctx.Logout()
		return
	}

	if !user.IsActive() {
		iLogError(fmt.Sprintf("User %s no more active. Logout", ctx.AuthIdentification()))
		ctx.Logout()
		return
	}

	role, err := config.UserProvider.GetRoleSupport(ctx, ctx.AuthIdentification())
	if err != nil {
		iLogError(err.Error())
	} else {
		ctx.SetRoleSupport(role)
	}
}

func handleError(ctx *Ctx, err error) {
	var errResponse Response

	iLogError(err.Error())

	if ctx.flashStorage != nil {
		ctx.flashStorage.Errors().SetRaw("error", err)
	}

	if ctx.IsJson() {
		errResponse = ctx.JsonResponse(nil, http.StatusInternalServerError)
	} else if config.Templates.Page500 != "" {
		errResponse = ctx.TemplateResponse(config.Templates.Page500)
		errResponse.(*TemplateResponse).SetCode(http.StatusInternalServerError)
	} else {
		http.Error(ctx.httpResponse, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	errResponse.Process(ctx)
}
