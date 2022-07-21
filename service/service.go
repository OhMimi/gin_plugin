package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
)

type APIType int

const (
	Get APIType = iota + 1
	Delete
	Put
	Patch
	Post
	Options
	Any
	Head
)

type Service struct {
	name    string                  // name of Service
	typ     reflect.Type            // type of the receiver
	val     reflect.Value           // receiver of methods for the Service
	methods map[string]*handlerFunc // registered methods
}

type handlerFunc struct {
	httpMethod APIType
	path       string
	fn         gin.HandlerFunc
}

func New(handler interface{}, openLog bool) *Service {
	service := new(Service)
	service.methods = make(map[string]*handlerFunc, 0)
	service.typ = reflect.TypeOf(handler)
	service.val = reflect.ValueOf(handler)
	for i := 0; i < service.typ.NumMethod(); i++ {
		m := service.typ.Method(i)

		// handler's method must be exported
		if m.PkgPath != "" {
			if openLog {
				log.Printf("handler: %s fn: %s needs to be exported\n", service.typ.Name(), m.Name)
			}
			continue
		}

		// handler param must be empty
		if m.Type.NumIn() != 1 {
			if openLog {
				log.Printf("handler: %s fn: %s ins needs to be empty\n", service.typ.Name(), m.Name)
			}
			continue
		}

		// call this func, get result
		result := m.Func.Call([]reflect.Value{service.val})
		// get param no.1
		apiType := result[0].Interface().(APIType)
		// get param no.2
		fn := result[1].Interface().(gin.HandlerFunc)
		// get router path
		path := CamelCaseToSnakeCase(m.Name)

		service.methods[m.Name] = &handlerFunc{apiType, path, fn}
	}
	return service
}

func (s *Service) Bind(r gin.IRoutes) {
	for mName, m := range s.methods {
		switch m.httpMethod {
		case Get:
			{
				r.GET(m.path, m.fn)
			}
		case Delete:
			{
				r.DELETE(m.path, m.fn)
			}
		case Put:
			{
				r.PUT(m.path, m.fn)
			}
		case Patch:
			{
				r.PATCH(m.path, m.fn)
			}
		case Post:
			{
				r.POST(m.path, m.fn)
			}
		case Options:
			{
				r.OPTIONS(m.path, m.fn)
			}
		case Any:
			{
				r.Any(m.path, m.fn)
			}
		case Head:
			{
				r.HEAD(m.path, m.fn)
			}
		default:
			log.Printf("service name: %s method: %s path ,not found match api type to register\n", s.name, mName)
		}
	}
}
