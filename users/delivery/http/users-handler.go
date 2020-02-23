package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Systemnick/user-service/users"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
	"gopkg.in/go-playground/validator.v10"
)

type Router interface {
	Get(path string, handlers ...routing.Handler) *routing.Route
	Post(path string, handlers ...routing.Handler) *routing.Route
	Put(path string, handlers ...routing.Handler) *routing.Route
	Patch(path string, handlers ...routing.Handler) *routing.Route
	Delete(path string, handlers ...routing.Handler) *routing.Route
	Head(path string, handlers ...routing.Handler) *routing.Route
	Options(path string, handlers ...routing.Handler) *routing.Route
	Any(path string, handlers ...routing.Handler) *routing.Route
	Group(prefix string, handlers ...routing.Handler) *routing.RouteGroup
	Use(handlers ...routing.Handler)
}

type UsersHandler struct {
	router     Router
	users      users.UseCase
	middleware routing.Handler
	validator  *validator.Validate
}

func NewHandler(r Router, uc users.UseCase, middleware routing.Handler) *UsersHandler {
	h := &UsersHandler{
		router:     r,
		users:      uc,
		middleware: middleware,
		validator:  validator.New(),
	}

	r.Use(middleware)

	// r.Get("/users", h.GetUserList)
	r.Post("/users", h.CreateUser)
	// r.Get("/users/<id>", h.GetUser)
	// r.Patch("/users/<id>", h.UpdateUser)
	// r.Delete("/users/<id>", h.DeleteUser)

	r.Post("/auth", h.AuthorizeUser)

	return h
}

func (h *UsersHandler) CreateUser(ctx *routing.Context) error {
	up := &users.UserParams{}
	err := json.Unmarshal(ctx.Request.Body(), &up)
	if err != nil {
		h.RespondWithError(ctx, errors.Wrap(err, "json.Unmarshal"))
		return nil
	}

	err = h.validator.Struct(up)
	if err != nil {
		h.RespondWithError(ctx, err)
		return nil
	}

	err = h.users.Register(ctx, up)
	if err != nil {
		h.RespondWithError(ctx, errors.Wrap(err, "users.Register"))
		return nil
	}

	return err
}

func (h *UsersHandler) AuthorizeUser(ctx *routing.Context) error {
	var resp []byte

	cred := &users.UserCredentials{}

	err := json.Unmarshal(ctx.Request.Body(), &cred)
	if err != nil {
		h.RespondWithError(ctx, errors.Wrap(err, "json.Unmarshal"))
		return nil
	}

	err = h.validator.Struct(cred)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		h.RespondWithError(ctx, errs)
		return nil
	}

	response := users.Response{}
	response.StatusCode = 200

	_, err = h.users.Login(ctx, cred)
	if err != nil {
		response.StatusCode = 403
		response.Errors = append(response.Errors, err.Error())
	}

	response.StatusMessage = http.StatusText(response.StatusCode)

	resp, err = json.Marshal(response)
	if err != nil {
		h.RespondWithError(ctx, err)
		return nil
	}

	ctx.Response.SetStatusCode(response.StatusCode)
	ctx.Response.AppendBody(resp)

	return nil
}

func (h *UsersHandler) RespondWithError(ctx *routing.Context, rErr interface{}) {
	var response users.Response
	var err error
	var resp []byte

	switch v := rErr.(type) {
	case validator.ValidationErrors:
		var errSlice []string
		for _, e := range v {
			errSlice = append(errSlice, fmt.Sprintf("%s is %s, type: %s", e.Field(), e.Tag(), e.Type()))
		}
		response.StatusCode = 400
		response.StatusMessage = http.StatusText(response.StatusCode)
		response.Errors = errSlice
	case error:
		response.StatusCode = 500
		response.StatusMessage = http.StatusText(response.StatusCode)
		response.Errors = append(response.Errors, v.Error())
	}

	resp, err = json.Marshal(response)
	if err != nil {
		h.RespondWithError(ctx, err)
		return
	}

	ctx.Response.SetStatusCode(response.StatusCode)
	ctx.Response.AppendBody(resp)
}
