# light
Golang light framework (not ready for production)

### Features:
 - Routing
 - Auth
 - Roles & Permissions
 - Middlewares
 - Context
 - Responses
 - Templates (pongo2)
 - Flashes
 - Static files
 - Container & Dependency Injection
 - Session
 - Validator
 - etc

### Examples:
#### Hello world
```
c := controller.MustSetup(...)

c.Get("/", controller.NewRootAction(func(ctx *controller.Ctx) (controller.Response, error) {
    return ctx.TextResponse("Hello world!", http.StatusOK), nil
}))

log.Fatal(c.Serve())
```

#### Inject service
```
type Printer interface {
    Print(w io.Writer, msg string)
}

type SimplePrinter struct{}

func (p *SimplePrinter) Print(w io.Writer, msg string) {
    _, _ = fmt.Fprintln(w, msg)
}

c.MustSingleton(&SimplePrinter{})

c.Get("/", func(printer Printer) *controller.Action {
    action := controller.NewRootAction(func(ctx *controller.Ctx) (controller.Response, error) {
        buf := &bytes.Buffer{}
        printer.Print(buf, "Hello World!")
        
        return ctx.TextResponse(buf.String(), http.StatusOK), nil
    })
    
    return action
})
```

#### Inject service & input validation
```
type Printer interface {
    Print(w io.Writer, msg string)
}

type SimplePrinter struct{}

func (p *SimplePrinter) Print(w io.Writer, msg string) {
    _, _ = fmt.Fprintln(w, msg)
}

c.MustSingleton(&SimplePrinter{})

c.Get("/show-query-param", func(printer Printer) *controller.Action {
    action := controller.NewAction(func(ctx *controller.Ctx) (controller.Response, error) {
        buf := &bytes.Buffer{}
        name := ctx.RequestValidatorInput().Value("name")
        printer.Print(buf, fmt.Sprintf("Hello %s!", html.EscapeString(name)))

        return ctx.TextResponse(buf.String(), http.StatusOK), nil
    }).SetValidator(controller.NewRequestValidator(func(ruleCollection *validator.RuleCollection) {
        ruleCollection.AddRule(
            "name",
            rules.Required(),
            rules.Length(10).SetMin(3).AsRunes(),
        )
    }))

    return action
})
```
#### Simple auth by email (without password)
```
c.Group().Prefix("/auth").Middlewares(controller.Guest()).Mount(func(m *controller.Mount) {
    m.Get("/login", controller.NewAction(func(ctx *controller.Ctx) (controller.Response, error) {

        return ctx.TemplateInlineResponse(`
            <html>
            <body>
                <form action="{{ Route("auth_login_process") }}" method="post">
                {{ CsrfTokenInput() }}
                <input type="text" name="login" placeholder="Type login..." value="{{ Input.OldOrDefault("login", "") }}" />
                {% if Err.Has("login") %}
                    <br />
                    <small style="color: red;">{{ Err.Get("login") }}</small>
                    <br />
                {% endif %}
                <input type="submit" value="Submit">
                </form>
            </body>
            </html>`,
        ), nil

    })).Name("auth_login_show")

    m.Post("/login", controller.NewAction(func(ctx *controller.Ctx) (controller.Response, error) {
        login := ctx.RequestValidatorInput().Value("login")
        
        user := GetUserByLogin(login) // implements controller.AuthIdentification interface
        ctx.Login(user)

        return ctx.RedirectResponse(":profile_index").With(func(response controller.ResponseExtendData) {
            response.AddMessage("alert-success", "Welcome, You Have Successfully Logged In.")
        }), nil

    }).SetValidator(controller.NewRequestValidator(func(ruleCollection *validator.RuleCollection) {
        ruleCollection.AddRule("login",
            rules.Required(),
            rules.Email())
    }))).Name("auth_login_process")
})

c.Get("/profile", controller.NewAction(func(ctx *controller.Ctx) (controller.Response, error) {
    return ctx.TextResponse("Profile page: " + ctx.AuthIdentification(), http.StatusOK), nil
    
}).WithMiddleware(controller.Auth())).Name("profile_index")
```

#### Full example
```
package main

import (
    "bytes"
    "fmt"
    "html"
    "io"
    "github.com/censoredgit/light/controller"
    "github.com/censoredgit/light/locker"
    "github.com/censoredgit/light/session"
    "github.com/censoredgit/light/session/driver/memory"
    "github.com/censoredgit/light/session/hasher"
    "github.com/censoredgit/light/validator"
    "github.com/censoredgit/light/validator/rules"
    "log"
    "log/slog"
    "net/http"
    "os"
    "strconv"
    "time"
)

type User struct {
    Name string
}

func (u User) AuthId() string {
    return u.Name
}

func (u User) IsActive() bool {
    return true
}

type Printer interface {
    Print(w io.Writer, msg string)
}

type SimplePrinter struct{}

func (p *SimplePrinter) Print(w io.Writer, msg string) {
    _, _ = fmt.Fprintln(w, msg)
}

func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    
    locks := locker.New(&locker.Config{GCTimeout: time.Hour})
    
    sessionManager := session.MustSetup(&session.Config{
        Salt:       "salt",
        TTL:        time.Hour,
        CookieName: "_test",
        Driver: memory.Setup(
            locks,
            time.Hour,
            time.Hour,
            memory.DefaultGarbageListInitCap,
        ),
        Logger: logger,
        Hasher: &hasher.Md5Hasher{},
    })
    
    c := controller.MustSetup(&controller.Config{
        Protocol:       "http",
        Host:           "127.0.0.1",
        Port:           "8080",
        Logger:         logger,
        SessionManager: sessionManager,
        Debug:          true,
        LoginRouteName: "auth_login_show",
    })
    
    c.RegisterTemplateFunc("inc", func(args ...any) string {
        return strconv.Itoa(args[0].(int) + 1)
    })
    
    c.MustSingleton(&SimplePrinter{})
        
    c.Get("/", func(iPrinter Printer, simplePrinter *SimplePrinter) *controller.Action {
        action := controller.NewRootAction(func(ctx *controller.Ctx) (controller.Response, error) {
            buf := &bytes.Buffer{}
    
            iPrinter.Print(buf, "<html><body>")
    
            if ctx.InputBag().Has("alert-success") {
                simplePrinter.Print(buf, "<div><span style=\"color: green;\">")
                simplePrinter.Print(buf, ctx.InputBag().Get("alert-success"))
                simplePrinter.Print(buf, "</span></div>")
            }
    
            iPrinter.Print(buf, "Hello world!")
            if ctx.IsAuth() {
                iPrinter.Print(buf, fmt.Sprintf(`<a href="%s">Profile</a>`, ctx.Route("profile_index")))
            } else {
                iPrinter.Print(buf, fmt.Sprintf(`<a href="%s">Login</a>`, ctx.Route("auth_login_show")))
            }
    
            simplePrinter.Print(buf, "</body></html>")
    
            return ctx.HtmlResponse(buf.String(), http.StatusOK), nil
        })
    
        return action
    }).Name("index")

    c.Post("/auth/logout", controller.NewAction(func(ctx *controller.Ctx) (controller.Response, error) {
        ctx.Logout()
    
        return ctx.RedirectResponse(":index").With(func(response controller.ResponseExtendData) {
            response.AddMessage("alert-success", "You have been successfully logged out.")
        }), nil
    
    }).WithMiddleware(controller.Auth())).Name("auth_logout")
    
    c.Group().Prefix("/auth").Middlewares(controller.Guest()).Mount(func(m *controller.Mount) {
        m.Get("/login", controller.NewAction(func(ctx *controller.Ctx) (controller.Response, error) {
    
            return ctx.TemplateInlineResponse(`
                <html>
                <body>
                    <form action="{{ Route("auth_login_process") }}" method="post">
                    {{ CsrfTokenInput() }}
                    <input type="text" name="login" placeholder="Type login..." value="{{ Input.OldOrDefault("login", "") }}" />
                    {% if Err.Has("login") %}
                    <br />
                        <small style="color: red;">{{ Err.Get("login") }}</small>
                    <br />
                    {% endif %}
                    <input type="submit" value="Submit">
                    </form>
                </body>
                </html>`,
            ), nil
    
        })).Name("auth_login_show")
    
        m.Post("/login", controller.NewAction(func(ctx *controller.Ctx) (controller.Response, error) {
            ctx.Login(User{
                Name: ctx.RequestValidatorInput().Value("login"),
            })
    
            return ctx.RedirectResponse(":profile_index").With(func(response controller.ResponseExtendData) {
                response.AddMessage("alert-success", "Welcome, You Have Successfully Logged In.")
            }), nil
    
        }).SetValidator(controller.NewRequestValidator(func(ruleCollection *validator.RuleCollection) {
            ruleCollection.AddRule("login",
                rules.Required(),
                rules.Email())
        }))).Name("auth_login_process")
    })
    
    c.Get("/profile", controller.NewAction(func(ctx *controller.Ctx) (controller.Response, error) {
        counter := 0
        if ctx.Session().Has("counter") {
            counter, _ = strconv.Atoi(ctx.Session().Get("counter"))
        }
        
        counter++
        ctx.Session().Set("counter", strconv.Itoa(counter))
    
        return ctx.TemplateInlineResponse(
            `
            <html>
            <body>
                {% if Input.Has("alert-success") %}
                    <div>
                        <span style="color: green;">{{ Input.Get("alert-success") }}</span>
                    </div>
                {% endif %}
                Hi, {{ AuthIdentification() }}!
                <br />
                Counter: {{ Data }}
                <br />
                Next Counter: {{ inc(Data) }}

                <form  method="get">
                <input type="submit" value="Refresh">
                </form>

                <form action="{{ Route("auth_logout") }}" method="post">
                    {{ CsrfTokenInput() }}
                    <input type="submit" value="Logout">
                </form>
            </body>
        </html>`,
        ).With(func(response controller.ResponseExtendData) {
            response.SetData(counter)
        }), nil
    }).WithMiddleware(controller.Auth())).Name("profile_index")
    
    log.Fatal(c.Serve())
}

```