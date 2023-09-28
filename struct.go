package golib

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func regStructToBson(elements reflect.Value, parentFieldTag string, result *bson.D, isUpdate bool) {
	tag := "bson"
	total := elements.NumField()
	elemTypes := elements.Type()
	var out bson.D

	if isUpdate {
		for i := 0; i < total; i++ {
			field := elemTypes.Field(i)
			if !field.IsExported() {
				continue
			}
			fieldTag := strings.Split(field.Tag.Get(tag), ",")[0]
			fieldValue := elements.Field(i)
			fieldZero := reflect.Zero(fieldValue.Type())
			// check for not zero
			if reflect.DeepEqual(fieldZero.Interface(), fieldValue.Interface()) {
				continue
			}
			if isUpdate {
				key := fmt.Sprintf("%s.%s", parentFieldTag, fieldTag)
				out = append(out, primitive.E{Key: key, Value: fieldValue.Interface()})
				continue
			}
		}
	} else {
		o := structToBson(elements, isUpdate, tag)
		out = append(out, primitive.E{Key: parentFieldTag, Value: o})
	}

	*result = append(*result, out...)
}

func structToBson(elements reflect.Value, isUpdate bool, tag string) bson.D {
	total := elements.NumField()
	elemTypes := elements.Type()
	var out bson.D

	for i := 0; i < total; i++ {
		field := elemTypes.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldTag := strings.Split(field.Tag.Get(tag), ",")[0]
		fieldValue := elements.Field(i)
		fieldZero := reflect.Zero(fieldValue.Type())

		// check for not zero
		if !reflect.DeepEqual(fieldZero.Interface(), fieldValue.Interface()) {
			// check for struct
			if fieldValue.Kind() == reflect.Struct {
				regStructToBson(fieldValue, fieldTag, &out, isUpdate)
				continue
			}

			val := reflect.ValueOf(fieldValue.Interface())
			valType := val.Kind()
			// check if field is pointer
			if valType == reflect.Pointer {
				val = val.Elem()

				// handle val to not zero value
				if val.Interface() != "" {
					actualType := val.Kind()
					// check if val is struct
					if actualType == reflect.Struct {
						regStructToBson(fieldValue, fieldTag, &out, isUpdate)
						continue
					}
					// if actual type is not struct
					out = append(out, primitive.E{Key: fieldTag, Value: val.Interface()})
				}
				continue
			}
			// regular type
			out = append(out, primitive.E{Key: fieldTag, Value: val.Interface()})
		}
	}

	return out
}

func StructToBson(from interface{}, isUpdate bool, tags ...string) bson.D {
	tag := "bson"
	if len(tags) > 0 {
		tag = tags[0]
	}
	elems := reflect.ValueOf(from).Elem()
	out := structToBson(elems, isUpdate, tag)
	return out
}
