package grdatabase

/*
import (
	"encoding/json"
)


type jsonStruct struct {
}

func (e jsonStruct) Value() jsonStruct {
	return e
}

func (e *jsonStruct) Set(d []string) {
	//*e = jsonStruct(d)
}

func (e *jsonStruct) Add(v string) {
	//*e = append(*e, v)
}

func (e *jsonStruct) String() string {

	bytes, _ := json.Marshal(*e)

	return string(bytes)
}

func (e *jsonStruct) FieldType() int {
	return TypeCharField
}

func (e *jsonStruct) SetRaw(value interface{}) error {
	switch d := value.(type) {
	case []string:
		e.Set(d)
	case string:
		if len(d) > 0 {
			parts := strings.Split(d, ",")
			v := make([]string, 0, len(parts))
			for _, p := range parts {
				v = append(v, strings.TrimSpace(p))
			}
			e.Set(v)
		}
	default:
		return fmt.Errorf("<jsonStruct.SetRaw> unknown value `%v`", value)
	}
	return nil
}

func (e *jsonStruct) RawValue() interface{} {
	return e.String()
}
*/
