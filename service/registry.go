package service

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

var (
	// Precompute the reflect.Type of error and http.Request
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
)

// ----------------------------------------------------------------------------
// service
// ----------------------------------------------------------------------------

type ServiceMeta struct {
	name    string                 // name of service
	rcvrV   reflect.Value          // receiver of methods for the service
	rcvrT   reflect.Type           // type of the receiver
	methods map[string]*MethodMeta // registered methods
}

func (r *ServiceMeta) ReceiverType() reflect.Type {
	return r.rcvrT
}

func (r *ServiceMeta) ReceiverValue() reflect.Value {
	return r.rcvrV
}

type MethodMeta struct {
	method     reflect.Method // receiver method
	paramTypes []reflect.Type // type of the request argument
	returnType reflect.Type   // type of the response argument
}

func (mm *MethodMeta) Call(in []reflect.Value) []reflect.Value {
	return mm.method.Func.Call(in)
}

func (mm *MethodMeta) ReturnType() reflect.Type {
	return mm.returnType
}

func (mm *MethodMeta) ParamValues() (values []reflect.Value, instances []interface{}) {
	if nil == mm.paramTypes || 0 == len(mm.paramTypes) {
		return nil, nil
	}

	pCount := len(mm.paramTypes)
	values = make([]reflect.Value, pCount)
	instances = make([]interface{}, pCount)

	for indexI := 0; indexI < pCount; indexI++ {
		values[indexI] = getValue(mm.paramTypes[indexI])
		instances[indexI] = values[indexI].Interface()
	}

	return
}

func getValue(t reflect.Type) reflect.Value {
	rt := t
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	// rv := reflect.New(rt)
	// if rt.Kind() != reflect.Struct {
	// 	rv = reflect.Indirect(rv)
	// }

	// var rv reflect.Value

	// 	switch rt.Kind() {
	// 	case reflect.Slice:
	// 		rv = reflect.New(reflect.SliceOf(rt.Elem()))
	// 	default:
	// 		rv = reflect.New(rt)
	// 	}

	return reflect.New(rt)
}

// ----------------------------------------------------------------------------
// Registry
// ----------------------------------------------------------------------------

// Registry is a registry for services.
type Registry struct {
	mutex    sync.RWMutex
	services map[string]*ServiceMeta
}

func (r *Registry) GetService(name string) interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if r.services == nil {
		return nil
	}
	return r.services[name]
}

// register adds a new service using reflection to extract its methods.
func (r *Registry) Register(rcvr interface{}, name string) error {
	// Setup service.
	s := &ServiceMeta{
		name:    name,
		rcvrV:   reflect.ValueOf(rcvr),
		rcvrT:   reflect.TypeOf(rcvr),
		methods: make(map[string]*MethodMeta),
	}
	if name == "" {
		s.name = reflect.Indirect(s.rcvrV).Type().Name()
		if !isExported(s.name) {
			return fmt.Errorf("Registry: type %q is not exported", s.name)
		}
	}
	if s.name == "" {
		return fmt.Errorf("Registry: no service name for type %q",
			s.rcvrT.String())
	}

	var err error
	// Setup methods.
Loop:
	for i := 0; i < s.rcvrT.NumMethod(); i++ {
		m := s.rcvrT.Method(i)
		mt := m.Type
		// Method must be exported.
		if m.PkgPath != "" {
			continue
		}

		var paramTypes []reflect.Type
		var returnType reflect.Type

		pCount := mt.NumIn() - 1

		if 0 < pCount {
			paramTypes = make([]reflect.Type, pCount)

			for indexI := 0; indexI < pCount; indexI++ {
				pt := mt.In(indexI + 1)
				if err = validateType(mt.In(indexI + 1)); nil != err {
					return err
				}
				paramTypes[indexI] = pt
			}
		}

		switch mt.NumOut() {
		case 1:
			if t := mt.Out(0); t != typeOfError {
				continue Loop
			}
		case 2:
			if t := mt.Out(0); !isExportedOrBuiltin(t) {
				continue Loop
			}

			if t := mt.Out(1); t != typeOfError {
				continue Loop
			}
			rt := mt.Out(0)
			if err = validateType(rt); nil != err {
				return err
			}
			returnType = rt
		default:
			continue
		}

		s.methods[m.Name] = &MethodMeta{
			method:     m,
			paramTypes: paramTypes,
			returnType: returnType,
		}
	}
	if len(s.methods) == 0 {
		return fmt.Errorf("Registry: %q has no exported methods of suitable type", s.name)
	}
	// Add to the map.
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.services == nil {
		r.services = make(map[string]*ServiceMeta)
	} else if _, ok := r.services[s.name]; ok {
		return fmt.Errorf("Registry: service already defined: %q", s.name)
	}
	r.services[s.name] = s
	return nil
}

func validateType(t reflect.Type) error {
	if t.Kind() == reflect.Struct {
		return fmt.Errorf("Type is Struct. Pass by reference, i.e. *%s", t)
	}
	return nil
}

// get returns a registered service given a method name.
//
// The method name uses a dotted notation as in "Service.Method".
func (r *Registry) Get(method string) (*ServiceMeta, *MethodMeta, error) {
	parts := strings.Split(method, ".")
	if len(parts) != 2 {
		err := fmt.Errorf("Registry: service/method request ill-formed: %q", method)
		return nil, nil, err
	}
	r.mutex.Lock()
	service := r.services[parts[0]]
	r.mutex.Unlock()
	if service == nil {
		err := fmt.Errorf("Registry: can't find service %q", method)
		return nil, nil, err
	}
	MethodMeta := service.methods[parts[1]]
	if MethodMeta == nil {
		err := fmt.Errorf("Registry: can't find method %q", method)
		return nil, nil, err
	}
	return service, MethodMeta, nil
}

// isExported returns true of a string is an exported (upper case) name.
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// isExportedOrBuiltin returns true if a type is exported or a builtin.
func isExportedOrBuiltin(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}
