package golib

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (g *Golib) structToBson(elements reflect.Value, tag string) bson.D {
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
				o := g.structToBson(fieldValue, tag)
				out = append(out, primitive.E{Key: fieldTag, Value: o})
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
						o := g.structToBson(val, tag)
						out = append(out, primitive.E{Key: fieldTag, Value: o})
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

func (g *Golib) StructToBson(from interface{}, tags ...string) bson.D {
	tag := "bson"
	if len(tags) > 0 {
		tag = tags[0]
	}
	elems := reflect.ValueOf(from).Elem()
	out := g.structToBson(elems, tag)
	return out
}
