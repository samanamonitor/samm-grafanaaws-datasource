package plugin

import (
    "strconv"
    "time"
    "strings"
    "fmt"
    "encoding/json"

    "github.com/grafana/grafana-plugin-sdk-go/backend"
    "github.com/grafana/grafana-plugin-sdk-go/data"
    "github.com/grafana/grafana-plugin-sdk-go/backend/log"

    "github.com/samana-group/sammaws/pkg/samm"
)

func dereferenceDefault(ptr *string, default_value string) string {
    if ptr != nil {
        return *ptr
    } else {
        return default_value
    }
}

func dereferenceTimestamp(ptr *time.Time) string {
    if ptr != nil {
        return strconv.FormatInt(ptr.UnixMilli(), 10)
    }
    return "0"
}

func dereferenceArray(ptr []*string) string {
    var temp []string
    for _, v := range ptr {
        temp = append(temp, *v)
    }
    return strings.Join(temp, ",")
}

func fieldsToResponse(fields []string, fieldlist []string) backend.DataResponse {
    frame := data.NewFrame("")
    frame.Fields = append(frame.Fields, data.NewField(fieldlist[0], nil, fields))
    if len(fieldlist) > 1 {
        frame.Fields = append(frame.Fields, data.NewField(fieldlist[1], nil, fields))
    }
    return backend.DataResponse{Frames: []*data.Frame{frame}}
}

func CreateFrame(elements samm.SammElement, fieldList []string, Name string) (*data.Frame, error) {
    frame := data.NewFrame(Name)
    frame.Meta = &data.FrameMeta{
        PreferredVisualization: "table",
    }

    if len(fieldList) == 0 {
        fieldList = elements.DefaultFieldList()
    }

    for _, fieldName := range fieldList {
        fieldType, ok := elements.AttributeType(fieldName)
        if ! ok {
            return nil, fmt.Errorf("Invalid field %s.", fieldName)
        }
        f := data.NewField(fieldName, nil, fieldType)
        if f == nil {
            return nil, fmt.Errorf("Unable to create field for type %s.", fieldName)
        }
        frame.Fields = append(frame.Fields, f)
    }
    for itemIndex := 0; itemIndex < elements.Len(); itemIndex++ {
        for index, fieldName := range fieldList {
            elements.AppendData(itemIndex, frame.Fields[index], fieldName)
        }
    }
    log.DefaultLogger.Debug("Frame created.", "elementsInFrame", elements.Len())
    return frame, nil
}

/* The following section can be used to convert data to plain JSON */
type CascaderOption struct {
    Label string `json:"label"`
    Value interface{} `json:"value"`
}

func intToString(i interface{}) string {
    return strconv.Itoa(i.(int))
}
func float32ToString(f interface{}) string {
    return strconv.FormatFloat(f.(float64), 'g', -1, 32)
}
func float64ToString(f interface{}) string {
    return strconv.FormatFloat(f.(float64), 'g', -1, 64)
}
func stringToString(s interface{}) string {
    return s.(string)
}
func boolToString(b interface{}) string {
    return strconv.FormatBool(b.(bool))
}
func timeToString(t interface{}) string {
    return t.(time.Time).String()
}

func toString(labelData interface{}, dataType data.FieldType) (func (interface{}) (string), error) {
    if dataType.Nullable() {
        return nil, fmt.Errorf("Cannot use data of nullable type for label")
    }

    switch dataType {
    case data.FieldTypeUnknown:
        return nil, fmt.Errorf("Invalid type '%s' for label.", dataType.String())

    case data.FieldTypeInt8:
        return intToString, nil

    case data.FieldTypeInt16:
        return intToString, nil

    case data.FieldTypeInt32:
        return intToString, nil

    case data.FieldTypeInt64:
        return intToString, nil

    case data.FieldTypeUint8:
        return intToString, nil

    case data.FieldTypeUint16:
        return intToString, nil

    case data.FieldTypeUint32:
        return intToString, nil

    case data.FieldTypeUint64:
        return intToString, nil

    case data.FieldTypeFloat32:
        return float32ToString, nil

    case data.FieldTypeFloat64:
        return float64ToString, nil

    case data.FieldTypeString:
        return stringToString, nil

    case data.FieldTypeBool:
        return boolToString, nil

    case data.FieldTypeTime:
        return timeToString, nil

    case data.FieldTypeJSON:
        return nil, fmt.Errorf("Invalid type 'JSON' for label.")

    case data.FieldTypeEnum:
        return nil, fmt.Errorf("Invalid type 'Enum' for label.")
    }
    return nil, fmt.Errorf("Unable to identify type.")
}

func ToVariables(response backend.DataResponse) ([]byte, error) {
    var out []CascaderOption
    for _, frame := range response.Frames {
        var fieldLabel, fieldValue *data.Field
        if len(frame.Fields) == 1 {
            fieldLabel = frame.Fields[0]
            fieldValue = frame.Fields[0]
        } else if len(frame.Fields) == 2 {
            fieldLabel = frame.Fields[0]
            fieldValue = frame.Fields[1]
        } else {
            return []byte{}, fmt.Errorf("Invalid request. Can only have to fields")
        }
        for i := 0; i < fieldLabel.Len(); i++ {
            label, err := toString(fieldLabel.At(i), fieldLabel.Type())
            if err != nil {
                return []byte{}, err
            }
            out = append(out, CascaderOption{Label: label(fieldLabel.At(i)), Value: fieldValue.At(i) } )
        }
    }
    return json.Marshal(out)
}