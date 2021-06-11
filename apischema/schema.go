package apischema

import (
	"bytes"
	"encoding/json"
)

type Schema struct {
	SchemaVersion int          `json:"schema_version"`
	Host          string       `json:"host,omitempty"`
	Namespace     string       `json:"namespace,omitempty"`
	API           string       `json:"api,omitempty"`
	RawObjects    []*RawObject `json:"objects,omitempty"`

	APIs               map[string]*API               `json:"-"`
	ArgumentSets       map[string]*ArgumentSet       `json:"-"`
	Controllers        map[string]*Controller        `json:"-"`
	Endpoints          map[string]*Endpoint          `json:"-"`
	Enums              map[string]*Enum              `json:"-"`
	Errors             map[string]*Error             `json:"-"`
	LookupArgumentSets map[string]*LookupArgumentSet `json:"-"`
	Objects            map[string]*Object            `json:"-"`
	Polymorphs         map[string]*Polymorph         `json:"-"`
	Scalars            map[string]*Scalar            `json:"-"`
}

//nolint:gocyclo
func (s *Schema) UnmarshalJSON(b []byte) error {
	type alias Schema

	err := strictJSONUnmarshal(b, (*alias)(s))
	if err != nil {
		return err
	}

	s.APIs = map[string]*API{}
	s.ArgumentSets = map[string]*ArgumentSet{}
	s.Controllers = map[string]*Controller{}
	s.Endpoints = map[string]*Endpoint{}
	s.Enums = map[string]*Enum{}
	s.Errors = map[string]*Error{}
	s.LookupArgumentSets = map[string]*LookupArgumentSet{}
	s.Objects = map[string]*Object{}
	s.Polymorphs = map[string]*Polymorph{}
	s.Scalars = map[string]*Scalar{}

	var newObjects []*RawObject
	for _, obj := range s.RawObjects {
		switch obj.Type {
		case "api":
			v := &API{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.APIs[v.ID] = v
		case "argument_set":
			v := &ArgumentSet{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.ArgumentSets[v.ID] = v
		case "controller":
			v := &Controller{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.Controllers[v.ID] = v
		case "endpoint":
			v := &Endpoint{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.Endpoints[v.ID] = v
		case "enum":
			v := &Enum{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.Enums[v.ID] = v
		case "error":
			v := &Error{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.Errors[v.ID] = v
		case "lookup_argument_set":
			v := &LookupArgumentSet{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.LookupArgumentSets[v.ID] = v
		case "object":
			v := &Object{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.Objects[v.ID] = v
		case "polymorph":
			v := &Polymorph{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.Polymorphs[v.ID] = v
		case "scalar":
			v := &Scalar{}
			err2 := strictJSONUnmarshal(obj.Value, v)
			if err2 != nil {
				return err2
			}
			s.Scalars[v.ID] = v
		default:
			newObjects = append(newObjects, obj)
		}
	}
	s.RawObjects = newObjects

	return nil
}

type RawObject struct {
	Type  string
	Value json.RawMessage
}

type API struct {
	ID            string      `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	Description   string      `json:"description,omitempty"`
	Authenticator string      `json:"authenticator,omitempty"`
	RouteSet      *RouteSet   `json:"route_set,omitempty"`
	Scopes        []*APIScope `json:"scopes,omitempty"`
}

type APIScope struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type RouteSet struct {
	Routes []*Route `json:"routes,omitempty"`
	Groups []*Group `json:"groups,omitempty"`
}

type Route struct {
	Path          string `json:"path,omitempty"`
	RequestMethod string `json:"request_method,omitempty"`
	Controller    string `json:"controller,omitempty"`
	Endpoint      string `json:"endpoint,omitempty"`
	Group         string `json:"group,omitempty"`
}

type Group struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Groups      []*Group `json:"groups,omitempty"`
}

type Enum struct {
	ID          string       `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Values      []*EnumValue `json:"values,omitempty"`
}

type EnumValue struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Error struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Code        string   `json:"code,omitempty"`
	HTTPStatus  int      `json:"http_status,omitempty"`
	Fields      []*Field `json:"fields,omitempty"`
}

type ArgumentSet struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Arguments   []*Argument `json:"arguments,omitempty"`
}

type LookupArgumentSet struct {
	ID              string      `json:"id,omitempty"`
	Name            string      `json:"name,omitempty"`
	Description     string      `json:"description,omitempty"`
	Arguments       []*Argument `json:"arguments,omitempty"`
	PotentialErrors []string    `json:"potential_errors,omitempty"`
}

type Argument struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Array       bool   `json:"array,omitempty"`
	Default     string `json:"default,omitempty"`
}

type Controller struct {
	ID            string                `json:"id,omitempty"`
	Name          string                `json:"name,omitempty"`
	Description   string                `json:"description,omitempty"`
	Authenticator string                `json:"authenticator,omitempty"`
	Endpoints     []*ControllerEndpoint `json:"endpoints,omitempty"`
}

type ControllerEndpoint struct {
	Name     string `json:"name,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

type Endpoint struct {
	ID              string       `json:"id,omitempty"`
	Name            string       `json:"name,omitempty"`
	Description     string       `json:"description,omitempty"`
	HTTPStatus      int          `json:"http_status,omitempty"`
	Authenticator   string       `json:"authenticator,omitempty"`
	ArgumentSet     *ArgumentSet `json:"argument_set,omitempty"`
	Fields          []*Field     `json:"fields,omitempty"`
	PotentialErrors []string     `json:"potential_errors,omitempty"`
	Scopes          []string     `json:"scopes,omitempty"`
}

type Object struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Fields      []*Field `json:"fields,omitempty"`
}

type Polymorph struct {
	ID          string             `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Options     []*PolymorphOption `json:"options,omitempty"`
}

type PolymorphOption struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type Scalar struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Field struct {
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Type        string     `json:"type,omitempty"`
	Null        bool       `json:"null,omitempty"`
	Array       bool       `json:"array,omitempty"`
	Spec        *FieldSpec `json:"spec,omitempty"`
}

type FieldSpec struct {
	All  bool   `json:"all,omitempty"`
	Spec string `json:"spec,omitempty"`
}

func strictJSONUnmarshal(b []byte, v interface{}) error {
	r := bytes.NewBuffer(b)
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	return dec.Decode(v)
}
