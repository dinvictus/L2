package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

// HTTPServer сервер для работы с календарём
type HTTPServer struct {
	httpCalendar Calendar
	config
	createupdateEventFields []string
	deleteEventFields       []string
	getFields               []string
}

type config struct {
	Port string
}

type valfields struct {
	userid, eventid uint
	message, date   string
}

type errorObj struct {
	Error string
}

type resultObj struct {
	Result interface{}
}

func (s *HTTPServer) startHTTPServer() {
	s.createupdateEventFields = []string{"userid", "eventid", "message", "date"}
	s.deleteEventFields = []string{"userid", "eventid"}
	s.getFields = []string{"userid"}
	logr := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	logrErr := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	errConfig := s.initConfig()
	if errConfig != nil {
		logrErr.Println(errConfig)
	}
	s.httpCalendar = CreateCalendar()
	mux := http.NewServeMux()
	mux.Handle("/create_event/", s.postEvent(logr, logrErr))
	mux.Handle("/update_event/", s.postEvent(logr, logrErr))
	mux.Handle("/delete_event/", s.postEvent(logr, logrErr))
	mux.Handle("/events_for_day/", s.getEvent(logr, logrErr))
	mux.Handle("/events_for_week/", s.getEvent(logr, logrErr))
	mux.Handle("/events_for_month/", s.getEvent(logr, logrErr))
	mux.Handle("/", errHandler(logrErr))
	logr.Println("Http server starting in", s.Port)
	err := http.ListenAndServe(s.Port, mux)
	if err != nil {
		logrErr.Println(err)
		os.Exit(1)
	}
}

func (s *HTTPServer) initConfig() error {
	confi := config{}
	conf, errOpen := os.Open("config.txt")
	if errOpen != nil {
		return errOpen
	}
	confBytes, errRead := ioutil.ReadAll(conf)
	if errRead != nil {
		return errRead
	}
	errParse := json.Unmarshal(confBytes, &confi)
	if errParse != nil {
		return errParse
	}
	confi.Port = ":" + confi.Port
	s.config = confi
	return nil
}

func errHandler(logrErr *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, errorJSON("404 not found"), http.StatusNotFound)
		logrErr.Println("Request:", r.RequestURI, "not found", "Status code: ", http.StatusNotFound)
	})
}

func (s HTTPServer) getEvent(logr, logrErr *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, errorJSON("Method not allowed"), http.StatusInternalServerError)
			logrErr.Println("Request:", r.RequestURI, "Method not allowed", "Status code: ", http.StatusInternalServerError)
			return
		}
		fields, err := validateData(r.URL.Query(), s.getFields)
		if err != nil {
			http.Error(w, errorJSON(err.Error()), http.StatusBadRequest)
			logrErr.Println("Request:", r.RequestURI, err, "Status code: ", http.StatusBadRequest)
			return
		}
		var obj interface{}
		var errGet error
		switch {
		case strings.Contains(r.RequestURI, "/events_for_day/"):
			obj, errGet = s.httpCalendar.evForDay(fields.userid)
		case strings.Contains(r.RequestURI, "/events_for_week/"):
			obj, errGet = s.httpCalendar.evForWeek(fields.userid)
		case strings.Contains(r.RequestURI, "/events_for_month/"):
			obj, errGet = s.httpCalendar.evForMonth(fields.userid)
		}
		if errGet != nil {
			http.Error(w, errorJSON(errGet.Error()), http.StatusServiceUnavailable)
			logrErr.Println("Request:", r.RequestURI, errGet, "Status code: ", http.StatusServiceUnavailable)
			return
		}
		json, errSer := serialize(obj)
		if errSer != nil {
			http.Error(w, errorJSON(errSer.Error()), http.StatusServiceUnavailable)
			logrErr.Println("Request:", r.RequestURI, errSer, "Status code: ", http.StatusServiceUnavailable)
			return
		}
		w.Write(json)
		logr.Println("Request:", r.RequestURI, "Status code: ", http.StatusOK)
	})
}

func (s HTTPServer) postEvent(logr, logrErr *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, errorJSON("Method not allowed"), http.StatusInternalServerError)
			logrErr.Println("Request:", r.RequestURI, "Method not allowed", "Status code: ", http.StatusInternalServerError)
			return
		}
		r.ParseForm()
		var fields valfields
		var errValidate error
		switch {
		case strings.Contains(r.RequestURI, "/create_event/") || strings.Contains(r.RequestURI, "/update_event/"):
			fields, errValidate = validateData(r.Form, s.createupdateEventFields)
		case strings.Contains(r.RequestURI, "/delete_event/"):
			fields, errValidate = validateData(r.Form, s.deleteEventFields)
		}
		if errValidate != nil {
			http.Error(w, errorJSON(errValidate.Error()), http.StatusBadRequest)
			logrErr.Println("Request:", r.RequestURI, errValidate, "Status code: ", http.StatusBadRequest)
			return
		}
		var errPost error
		switch {
		case strings.Contains(r.RequestURI, "/create_event/"):
			errPost = s.httpCalendar.createEv(fields.message, fields.date, fields.eventid, fields.userid)
		case strings.Contains(r.RequestURI, "/update_event/"):
			errPost = s.httpCalendar.updateEv(fields.message, fields.date, fields.eventid, fields.userid)
		case strings.Contains(r.RequestURI, "/delete_event/"):
			errPost = s.httpCalendar.deleteEv(fields.eventid, fields.userid)
		}
		if errPost != nil {
			http.Error(w, errorJSON(errPost.Error()), http.StatusServiceUnavailable)
			logrErr.Println("Request:", r.RequestURI, errPost, "Status code: ", http.StatusServiceUnavailable)
			return
		}
		logr.Println("Request:", r.RequestURI, "Status code: ", http.StatusOK)
	})
}

func serialize(obj interface{}) ([]byte, error) {
	resObj := resultObj{Result: obj}
	json, err := json.Marshal(resObj)
	if err != nil {
		return nil, err
	}
	return json, nil
}

func errorJSON(err string) string {
	er := errorObj{Error: err}
	json, errJSON := json.Marshal(er)
	if errJSON != nil {
		return ""
	}
	return string(json)
}

func validateData(form url.Values, rFields []string) (valfields, error) {
	validateFields := valfields{}
	for _, field := range rFields {
		value := form.Get(field)
		switch field {
		case "userid":
			userid, errUserid := strconv.Atoi(value)
			if errUserid != nil || userid < 0 {
				return validateFields, errors.New("error userid format")
			}
			validateFields.userid = uint(userid)
		case "eventid":
			eventid, errEventid := strconv.Atoi(value)
			if errEventid != nil || eventid < 0 {
				return validateFields, errors.New("error eventid format")
			}
			validateFields.eventid = uint(eventid)
		case "message":
			if value == "" {
				return validateFields, errors.New("message empty")
			}
			validateFields.message = value
		case "date":
			if value == "" {
				return validateFields, errors.New("date empty")
			}
			validateFields.date = value
		}
	}
	return validateFields, nil
}

func main() {
	server := HTTPServer{}
	server.startHTTPServer()
}
