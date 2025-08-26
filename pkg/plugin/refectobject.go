package plugin

import (
	"errors"
	"reflect"
	"unsafe"
    "time"
    "fmt"
    "strings"
    "github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/aws/aws-sdk-go/service/workspaces"
)

type ReflectObject struct {
    typeReflect reflect.Type
    valueReflect reflect.Value
    fieldNames []string
    fieldTypes map[string]string
    isString bool
}

func NewReflectObject(v reflect.Value, t reflect.Type) *ReflectObject {
    
    wo := new(ReflectObject)
    wo.typeReflect = t
    wo.valueReflect = v
    if wo.valueReflect.Type().String() == "string" {
        wo.fieldNames = []string{"value"}
        wo.fieldTypes = map[string]string{"value": "string"}
        wo.isString = true
        return wo
    }
    numfield := wo.valueReflect.NumField()
    wo.fieldNames = make([]string, numfield)
    wo.fieldTypes = make(map[string]string)
    for i := 0; i < numfield; i++ {
        wo.fieldNames[i] = wo.typeReflect.Field(i).Name
        wo.fieldTypes[wo.fieldNames[i]] = wo.valueReflect.Field(i).Type().String()
    }
    wo.isString = false
    return wo
}   

func (w ReflectObject) newField(name string) (*data.Field, error) {
    var fieldType string
    if w.isString {
        return data.NewField(name, nil, []string{}), nil
    }
    if temp, ok := w.fieldTypes[name]; ok {
        fieldType = temp
    } else {
        return nil, errors.New("Invalid Field Name")
    }
    switch fieldType {
    case "*string":
        return data.NewField(name, nil, []*string{}), nil
    case "*float64":
        return data.NewField(name, nil, []*float64{}), nil
    case "*int64":
        return data.NewField(name, nil, []*int64{}), nil
    case "*time.Time":
        return data.NewField(name, nil, []*time.Time{}), nil
    case "[]*string":
        return data.NewField(name, nil, []string{}), nil
    case "string":
        return data.NewField(name, nil, []string{}), nil
    default:
        return data.NewField(name, nil, []string{}), nil
        //return nil, errors.New(fmt.Sprintf("Invalid Field Type %s", fieldType))
    }
}

type Test interface {
    String() string
}

func (w ReflectObject) AddFrameFields(frame *data.Frame, fieldlist []string) error {
    var field *data.Field
    for i, fieldName := range fieldlist {
        if len(frame.Fields) < len(fieldlist) {
            var err error
            field, err = w.newField(fieldName)
            if err != nil {
                return err
            }
            frame.Fields = append(frame.Fields, field)
        } else {
            field = frame.Fields[i]
        }
        if w.isString {
            field.Append(w.valueReflect.String())
        } else {
            v := w.valueReflect.FieldByName(fieldName).Pointer()
            switch w.fieldTypes[fieldName] {
            case "*string":
                field.Append((*string)(unsafe.Pointer(v)))
            case "*float64":
                field.Append((*float64)(unsafe.Pointer(v)))
            case "*int64":
                field.Append((*int64)(unsafe.Pointer(v)))
            case "*time.Time":
                field.Append((*time.Time)(unsafe.Pointer(v)))
            case "[]*string":
                sliceLen := w.valueReflect.FieldByName(fieldName).Len()
                v1 := w.valueReflect.FieldByName(fieldName).Slice(0, sliceLen)
                temp := make([]string, sliceLen)
                for i := 0; i < sliceLen; i++ {
                    temp[i] = *(*string)(unsafe.Pointer(v1.Index(i).Pointer()))
                }
                field.Append(strings.Join(temp, ","))
            default:
                test := w.valueReflect.FieldByName(fieldName).Elem().Addr().Interface()
                if stringer, ok := test.(*workspaces.ActiveDirectoryConfig); ok {
                    field.Append(stringer.String())
                } else {
                    field.Append(fmt.Sprintf("%s", fieldName))
                }
                //return errors.New(fmt.Sprintf("Invalid Field Type %s", w.fieldTypes[fieldName]))
            }
        }
    }
    return nil
}