package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	erHanG "errorHandler/gunk/v1/errorHandler"
	"github.com/spf13/viper"

	//"google.golang.org/grpc/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//"google.golang.org/grpc@v1.42.0/codes/codes.go"
)

type ErrorHandler struct {
	ID              string
	ErrorCode       string
	ErrorDetails    string
	EnvType         string
	CreatedAt       time.Time
	CreatedBy       string
	DeletedAt       time.Time
	DeleteByEnvType string
}

type ErrorHandlerData struct {
	Form ErrorHandler
}

func (s *Server) errorHandlerGet(rw http.ResponseWriter, r *http.Request) {

	template := s.lookupTemplate("error.html")
	if template == nil {
		errMsg := "unable to load template"
		log.Error(errMsg)
		http.Error(rw, errMsg, http.StatusInternalServerError)
		return
	}

	queryParams := r.URL.Query()
	id, err := url.PathUnescape(queryParams.Get("code"))
	if err != nil {
		s.logger.Error("unable to decode url type param")
	}
	var data ErrorHandlerData
	if id == "" {
		data = ErrorHandlerData{
			Form: ErrorHandler{
				ID:           "",
				ErrorCode:    "404 Not Found",
				ErrorDetails: "The page you are looking for might be changed, removed or not exists. Go back and try other links",
				EnvType:      "",
			},
		}

		if err := template.Execute(rw, data); err != nil {
			http.Redirect(rw, r, s.ErrorHandlerAdd(rw, r, err), http.StatusSeeOther)
		}

	} else {
		res, err := s.erHanG.GetErrorHandler(r.Context(), &erHanG.GetErrorHandlerRequest{
			ID: id,
		})
		if err != nil {
			http.Redirect(rw, r, s.ErrorHandlerAdd(rw, r, err), http.StatusSeeOther)
		}
		data = ErrorHandlerData{
			Form: ErrorHandler{
				ID:           res.ErrorHandler.ID,
				ErrorCode:    res.ErrorHandler.ErrorCode,
				ErrorDetails: res.ErrorHandler.ErrorDetails,
				EnvType:      res.ErrorHandler.EnvType,
			},
		}
	}

	if err := template.Execute(rw, data); err != nil {
		log.Infof("error with template execution: %+v", err)
	}
}

func (s *Server) ErrorHandlerAdd(rw http.ResponseWriter, r *http.Request, data error) string {

	value := fmt.Sprintf("", data)

	i := status.Code(data)
	htt := http.StatusNotFound
	switch i {
	case codes.OK:
		htt = http.StatusOK
	case codes.Canceled:
		htt = http.StatusCreated
	case codes.Unknown:
		htt = http.StatusUnauthorized
	case codes.InvalidArgument:
		htt = http.StatusUnauthorized
	case codes.DeadlineExceeded:
		htt = http.StatusAlreadyReported
	case codes.NotFound:
		htt = http.StatusNotFound
	case codes.AlreadyExists:
		htt = http.StatusNotAcceptable
	case codes.PermissionDenied:
		htt = http.StatusBadRequest
	case codes.ResourceExhausted:
		htt = http.StatusNonAuthoritativeInfo
	case codes.FailedPrecondition:
		htt = http.StatusBadRequest
	case codes.Aborted:
		htt = http.StatusNotModified
	case codes.OutOfRange:
		htt = http.StatusBadGateway
	case codes.Unimplemented:
		htt = http.StatusNotImplemented
	case codes.Internal:
		htt = http.StatusInternalServerError
	case codes.Unavailable:
		htt = http.StatusNotExtended
	case codes.DataLoss:
		htt = http.StatusFailedDependency
	case codes.Unauthenticated:
		htt = http.StatusNoContent
	default:
		htt = http.StatusNotFound
	}

	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	if value == "" {
		var desk []string

		result := strings.SplitN(value, ":", 3)
		for v := range result {
			code := result[v]
			eode := strings.SplitN(code, "code =", 3)

			for v := range eode {
				des := eode[v]
				desc := strings.SplitN(des, "desc =", 3)
				for _, v := range desc {
					desk = append(desk, v)
				}
			}
		}

		ErrorCode := desk[2]
		ErrorDetails := desk[3]

		id, err := s.erHanG.CreateErrorHandler(r.Context(), &erHanG.CreateErrorHandlerRequest{
			ErrorHandler: &erHanG.ErrorHandler{
				ErrorCode:    ErrorCode,
				ErrorDetails: ErrorDetails,
				EnvType:      config.GetString("runtime.environment"),
				CreatedBy:    "5c2b110d-105d-4c43-8663-f309fb9e0768",
			},
		})

		if err != nil {
			http.Redirect(rw, r, ErrorPath, http.StatusSeeOther)
		}

		baseUrl, err := url.Parse("/error")
		if err != nil {
			fmt.Println("Malformed URL: ", err.Error())
		}
		params := url.Values{}
		params.Add("code", id.ID)
		baseUrl.RawQuery = params.Encode() // Escape Query Parameters

		return baseUrl.String()
	} else {

		value, code := s.ErrorFormet(htt)

		id, err := s.erHanG.CreateErrorHandler(r.Context(), &erHanG.CreateErrorHandlerRequest{
			ErrorHandler: &erHanG.ErrorHandler{
				ErrorCode:    code,
				ErrorDetails: value,
				EnvType:      config.GetString("runtime.environment"),
				CreatedBy:    "5c2b110d-105d-4c43-8663-f309fb9e0768",
			},
		})

		if err != nil {
			http.Redirect(rw, r, ErrorPath, http.StatusSeeOther)
		}

		baseUrl, err := url.Parse("/error")
		if err != nil {
			fmt.Println("Malformed URL: ", err.Error())
		}
		params := url.Values{}
		params.Add("code", id.ID)
		baseUrl.RawQuery = params.Encode() // Escape Query Parameters

		return baseUrl.String()
	}

}

func (s *Server) ErrorFormet(code int) (string, string) {

	codeA := http.StatusText(code)
	str := strconv.Itoa(code)

	return codeA, str
}
