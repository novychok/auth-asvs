// Package authapiv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.0 DO NOT EDIT.
package authapiv1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ChangePasswordRequest defines model for ChangePasswordRequest.
type ChangePasswordRequest struct {
	CurrentPassword    string              `json:"currentPassword"`
	Email              openapi_types.Email `json:"email"`
	NewPassword        string              `json:"newPassword"`
	NewPasswordConfirm string              `json:"newPasswordConfirm"`
}

// Error defines model for Error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Email    openapi_types.Email `json:"email"`
	Password string              `json:"password"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

// RefreshTokenRequest defines model for RefreshTokenRequest.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// RegisterRequest defines model for RegisterRequest.
type RegisterRequest struct {
	Email           openapi_types.Email `json:"email"`
	Name            string              `json:"name"`
	Password        string              `json:"password"`
	PasswordConfirm string              `json:"passwordConfirm"`
}

// ChangePasswordJSONRequestBody defines body for ChangePassword for application/json ContentType.
type ChangePasswordJSONRequestBody = ChangePasswordRequest

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = LoginRequest

// RefreshTokenJSONRequestBody defines body for RefreshToken for application/json ContentType.
type RefreshTokenJSONRequestBody = RefreshTokenRequest

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody = RegisterRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Change Password
	// (POST /change-password)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	// Login
	// (POST /login)
	Login(w http.ResponseWriter, r *http.Request)
	// Refresh Token
	// (POST /refresh-token)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	// Register
	// (POST /register)
	Register(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Change Password
// (POST /change-password)
func (_ Unimplemented) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Login
// (POST /login)
func (_ Unimplemented) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Refresh Token
// (POST /refresh-token)
func (_ Unimplemented) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Register
// (POST /register)
func (_ Unimplemented) Register(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// ChangePassword operation middleware
func (siw *ServerInterfaceWrapper) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ChangePassword(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Login operation middleware
func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Login(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// RefreshToken operation middleware
func (siw *ServerInterfaceWrapper) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RefreshToken(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Register operation middleware
func (siw *ServerInterfaceWrapper) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Register(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/change-password", wrapper.ChangePassword)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/login", wrapper.Login)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/refresh-token", wrapper.RefreshToken)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/register", wrapper.Register)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xXzW7bOBB+FYG7wF4Uy9nsAoFuTpBDsAs0iNteihwYaWwxlUhmOEriBnr3gtSPJZty",
	"YiBuLz7ZIOfnm5mPM6NXlqhCKwmSDItfmUkyKLj7e5lxuYQbbsyzwvQWHkswZC80Kg1IApxYUiKCpFbO",
	"Hi0UFpxYzHR7GDJaaWAxM4RCLlkVMii4yAfS9YlHVMLz3uZ7OpdKLgQW71StQobwWAqElMXfOlCbYQ5R",
	"ef3ddbbV/QMkZGFdISr0JFGlYH8beSEJloBWoQBj+LJ/OQLUmVjL+5z/r5ZCjhZyj4Lo/aoxktJOYQdW",
	"o5U0sA2WJwkY81l9B+lJTdi/v3rRAsHMaIA25QQnJArwxYewQDDZuPm+wN72N9LRD2XD9UgcY/59ebzt",
	"SY6W/o14NwAPpP0+l8IQ4EdQTfICvCXQ+3YE/RHtwMEJtym8bX47MVXIDCQlClrNbZutc3EPHAFnJWVd",
	"/7VKF+54jScj0qyyNoRcKCuagklQaBJKsphZA8Fs/nUezG6urZqgHDznT4Cm1jidTCdTmxilQXItWMzO",
	"JtPJmYuFMgcuStwQOOknW6u6pEP/9bQIKIOAl5SBJJFwgjQoDeBfJmgtTJhziNzqXaedZq+VYk2bC5Wu",
	"6t4oCaRzybXOrVmhZPRglFxPLPvvT4QFi9kf0XqkRc08i/zDrBqWl7AEd1B3HZeBv6fTvUAMid7r3vDC",
	"C+1K0qIISp26FJnSvfJFmecrLwM3iVSFG9nvTPZNBXXxUlvjf/YMY1cu6xnmQXHB06BLrfV5enifX6Sl",
	"m0Lxow70318R6LUkQMnzYA74BBi0giEzZVFwXK0fRI/XxJfG9pBZA9hBYndWLcrtuBt/XW4abr2c9vQQ",
	"D2awKxzgnbzDdzP7Pfn/9N+R07+D0y3hdjG52Q5OqF0m/Ixu9pKgXXSGzL4dbkGHILhvMTry/MjzlYec",
	"u/leL7u7qN5IbLO8uzgMw4dr+JHdR3avhrQbJbb7VrEW7PUrKzFvvkHiKMpVwvNMGYrPp+enrArfuL+r",
	"fgYAAP//FvLecWYSAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
