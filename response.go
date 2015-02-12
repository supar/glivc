package main

import (
	"reflect"
	"bytes"
	"encoding/json"
)


// Default response structure
type Response struct {
	Success	bool					`json:"success"`
	Results interface{}				`json:"results,omitempty"`
	Options interface{}				`json:"options,omitempty"`
	Data	[]byte					`json:"-"`
	Error	*Error					`json:"error,omitempty"`
}


// Common error structure in response
type Error struct {
	Code	int		`json:"code"`
	Message	string	`json:"message"`
}


func (this *Response) Compile() (int, string) {
	var (
		code int = 200
		enc []byte
		err error
	)
	
	if this.Error != nil {
		code = this.Error.Code
		
		this.Success = false
		this.Results = nil
		this.Options = nil
	}
	
	enc, err = json.Marshal(this)
	if err != nil {
		log.Warn(err.Error())
	}
	
	if DEBUG {
		var out bytes.Buffer
		json.Indent(&out, enc, "", "\t")
		return code, string(out.Bytes())
	} else {
		return code, string(enc)
	}
}


func (this *Response) SetData(data []byte) {
	this.Data = data
	
	// Debug response data
	log.Debug("%s", string(this.Data))
}

// Fail response 
// Set HTTP status according to the code. Default is 500
// Send response with filled error attribute
func (this *Response) SetError(code int, err ...interface{}) {
	var (
		msg string
	)
	
	switch len(err) {
		case 0: 
			msg = "Unknown error";
		
		break;
		
		default:
			for i, v := range err {
				if i == 0 {
					msg = getError(v)
				}
				
				log.Error(getError(v))
			}
	}
	
	this.Success = false
	this.Error = &Error{
		Code: code,
		Message: msg,
	}
}

// Describe error
func getError(err interface{}) string {
	if msg, ok := err.(string); ok {
		return msg
	} else {
		e := reflect.ValueOf(err)
		
		if e.MethodByName("Error").IsValid() {
			return e.MethodByName("Error").Call(nil)[0].String()
		}
	}
	
	return "Unknown error"
}

// Get indirect type of element
func indirectType(v reflect.Type) reflect.Type {
	switch v.Kind() {
		case reflect.Ptr:
			return indirectType(v.Elem())
		default:
			return v
	}
}