package rest_server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type RestServerContext struct {
	ContentType string
	StatusCode  int
}

type RestResponse struct {
	Status      bool        `json:"status"`
	Description interface{} `json:"description"`
	Response    interface{} `json:"response"`
}

type RestServerError struct {
	StatusCode  int
	Description string
}

func (e *RestServerError) Error() string {
	return fmt.Sprintf("parse %v: internal error", e.Description)
}

func CallMethod(m map[string]interface{}, name string, params map[string]interface{}, context *RestServerContext) (result reflect.Value, err error) {
	function := reflect.ValueOf(m[name])
	functionType := reflect.TypeOf(m[name])
	firstArgument := functionType.In(0)
	args := reflect.New(firstArgument).Interface()

	s := reflect.ValueOf(args).Elem()
	if s.Kind() == reflect.Struct {
		// Fill context
		contextField := s.FieldByName("Context")
		if contextField.IsValid() {
			if contextField.CanSet() {
				contextField.Set(reflect.ValueOf(context))
			}
		}

		for k, v := range params {
			f := s.FieldByName(strings.Title(k))

			if f.IsValid() {
				if f.CanSet() {
					// Fill string types
					if f.Kind() == reflect.String && reflect.TypeOf(v).Kind() == reflect.String {
						f.SetString(v.(string))
					}

					// Fill int types
					if f.Kind() == reflect.Int64 {
						// From int
						if reflect.TypeOf(v).Kind() == reflect.Int64 {
							f.SetInt(v.(int64))
						}
						// From string
						if reflect.TypeOf(v).Kind() == reflect.String {
							i, _ := strconv.ParseInt(v.(string), 10, 64)
							f.SetInt(i)
						}
					}

					// Fill maps
					/*if f.Kind() == reflect.Map {
						if reflect.TypeOf(v).Kind() == reflect.Map && strings.Title(k) == "Files" {

							mapData := v.(map[string][]*multipart.FileHeader)
							fmt.Println("XXXXX", mapData)
							for k, v := range mapData {
								f.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
							}
						}
					}*/

					// Fill int types
					if f.Kind() == reflect.Bool {
						// From int
						if reflect.TypeOf(v).Kind() == reflect.Bool {
							f.SetBool(v.(bool))
						}
						// From string
						if reflect.TypeOf(v).Kind() == reflect.String {
							if v.(string) == "true" {
								f.SetBool(true)
							} else {
								f.SetBool(false)
							}
						}
					}
				}
			}

		}
	}

	// Call function
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(s.Interface())
	result = function.Call(in)[0]

	return
}

func ErrorMessage(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	if err := recover(); err != nil {
		switch e := err.(type) {
		case *RestServerError:
			if e.StatusCode == 0 {
				rw.WriteHeader(500)
			} else {
				rw.WriteHeader(e.StatusCode)
			}

			responseData := RestResponse{Status: false, Description: e.Description}
			finalData, _ := json.Marshal(responseData)
			fmt.Fprintf(rw, "%+v", string(finalData))
		default:
			fmt.Fprintf(rw, "%+v", e)
		}
	}
}

func Error(code int, description string) {
	panic(&RestServerError{StatusCode: code, Description: description})
}

func Start(addr string, controller map[string]interface{}) {
	fmt.Printf("Starting server at port 8080\n")

	// Set handler
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		defer ErrorMessage(rw, r)

		// Fuck options
		if r.Method == "OPTIONS" {
			rw.WriteHeader(200)
			fmt.Fprintf(rw, "")
			return
		}

		// Collect args
		args := map[string]interface{}{}
		for key, element := range r.URL.Query() {
			args[key] = element[0]
		}

		// Parse body
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			// Parse multipart body and collect args
			r.ParseMultipartForm(0)
			for key, element := range r.MultipartForm.Value {
				args[key] = element[0]
			}
			if len(r.MultipartForm.File) > 0 {
				args["files"] = r.MultipartForm.File
			}
		} else {
			// Parse json body and collect args
			jsonMap := make(map[string]interface{})
			err := json.NewDecoder(r.Body).Decode(&jsonMap)
			if err != nil {
				Error(500, "Invalid JSON")
			}
			for key, element := range jsonMap {
				args[key] = element
			}
		}
		fmt.Println(args)

		// Get function
		path := strings.ToLower(r.Method) + "_" + strings.Split(r.URL.Path, "/")[1]
		if controller[path] == nil {
			rw.WriteHeader(http.StatusNotFound)

			fmt.Fprintf(rw, "Path not found")
			return
		}

		// Call method
		context := new(RestServerContext)
		context.ContentType = "application/json"
		context.StatusCode = 200
		response, err := CallMethod(controller, path, args, context)

		if err != nil {
			context.StatusCode = 500
			fmt.Fprintf(rw, "Error")
			return
		}

		// Response
		rw.Header().Add("Access-Control-Allow-Origin", "*")
		rw.Header().Add("Access-Control-Allow-Methods", "*")
		rw.Header().Add("Access-Control-allow-Headers", "*")
		rw.Header().Add("Content-Type", context.ContentType)

		if context.ContentType == "application/json" {
			responseData := RestResponse{Status: true}
			responseData.Response = response.Interface()
			finalData, _ := json.Marshal(responseData)
			fmt.Fprintf(rw, "%+v", string(finalData))
		} else {
			fmt.Fprintf(rw, "%+v", response)
		}
	})

	// Start server
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
		return
	}
}
