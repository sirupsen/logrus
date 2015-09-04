package logrus

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// HasFields is an interface that enables a type to control the manner in
// which it is formatted as a Logrus field value.
//
// For example, take the following type:
//
//     type Person struct {
//         Name  string
//         Alias string
//     }
//
//     p := &Person{"Bruce", "Batman"}
//
// If p were given to Logrus as a field value with an associated key of "hero"
// then it would be formatted as such (given a text formatter):
//
//     hero="&{Bruce Batman}"
//
// However, the Person type may know how it wants to be formatted when it comes
// to be logged as a field value, and this interface will provide that
// advantage.
//
//     func (p *Person) Fields() map[string]interface{} {
//         return map[string]interface{} {"name" : p.Name, "alias" : p.Alias}
//     }
//
// Now when a Person instance is provided to Logrus as a field value with an
// associated key of hero it would be formatted as such (given a text
// formatter):
//
//     hero.name=Bruce hero.alias=Batman...
//
// The JSON formatter will not nest field names as there is no precedent for
// such behavior and consumers of the data expect a direct serialization of a
// series of key/value pairs without additional object structure. Given then
// the same example above, the JSON formatter would produce the following
// output:
//
//     { "hero.name" : "bruce", "hero.alias" : "batman" ... }
//
// The values in the map returned by the Fields function are also considered
// with regards to whether they implement the HasFields interface, enabling
// multiple levels of field data to find its way into the log statement in an
// ordered fashion.
//
// There may be cases where it is not desireable to have all keys in a possible
// key path emitted along with a value. To that end this interface defines the
// Flatten function. If Flatten returns a true value, it indicates that the
// top-level key should be removed when emitting the field data. For example:
//
//     func (p *Person) Flatten() {
//         return true
//     }
//
// Because the Flatten function returns true, if the Person instance p
// is logged as a field value with an associated key of "hero" this is
// what will actually be emitted by the text formatter:
//
//     name=Bruce alias=Batman...
//
// The key supplied along with a value is stripped from the emitted field
// names when the Flatten function returns true.
//
// Additionally, it's possible to bypass the explicit naming of fields
// altogether. The HasFields interface requires types to implement a function
// called UseTypeFields that returns a flag indicating whether or not a type's
// public fields should be used instead of any explicit key/value pairs
// returned by the Fields function.
//
//
//     func (p *Person) Fields() map[string]interface{} {
//         return nil
//     }
//
//     func (p *Person) UseTypeFields() bool {
//         return true
//     }
//
// With the above redefinition of the Fields function for the Person type, and
// the UseTypeFields function returning true, the text and JSON formatters
// will produce the same output as the earlier examples without having to
// specify any explicit fields.
//
// Finally, the HasFields interface defines one more function, ExceptFields.
// This function returns an array of strings -- a list of field names to not
// allow to be emitted as part of the log message.
//
//     func (p *Person) ExceptFields() []string {
//         return []string {"name"}
//     }
//
// Because the ExceptFields function returns an array containing the field
// called "name," this is what will be emitted by the text formatter:
//
//     alias=Batman...
//
// It doesn't matter whether the field is discovered because UseTypeFields
// returns true or if the field is explicitly returned via the Fields function.
// If the name of the field (case insensitive) is in the array returned by the
// ExceptFields function  the field's value will not be emitted.
type HasFields interface {
	// UseTypeFields returns a flag that indicates whether or not a type's
	// public fields should be used instead of the key/value pairs returned by
	// the Fields function.
	UseTypeFields() bool

	// ExceptFields returns a list of a type's field names to omit even
	// if the UseTypeFields funtion returns a true value.
	ExceptFields() []string

	// Fields returns the data to format as key/value pairs.
	Fields() map[string]interface{}

	// Flatten returns a flag indicating whether or not to keep the top-level
	// key when emitting the field data.
	Flatten() bool
}

func parseFields(key string, val HasFields, fields map[string]interface{}) {
	var vf map[string]interface{}

	if val.UseTypeFields() {

		vf = map[string]interface{}{}

		valType := reflect.ValueOf(val).Elem()
		elmType := valType.Type()

		for x := 0; x < elmType.NumField(); x++ {
			f := valType.Field(x)
			if !f.CanInterface() {
				continue
			}

			fn := strings.ToLower(elmType.Field(x).Name)
			vf[fn] = f.Interface()
		}

	} else {
		vf = val.Fields()
	}

	if val.ExceptFields() != nil {
		except := strings.Join(val.ExceptFields(), " ")
		dk := []string{}
		for k, _ := range vf {
			namePatt := fmt.Sprintf("(?i)\\b%s\\b", k)

			isMatch, mErr := regexp.MatchString(namePatt, except)
			if mErr != nil {
				panic(mErr)
			}

			if isMatch {
				dk = append(dk, k)
			}
		}
		for _, k := range dk {
			delete(vf, k)
		}
	}

	for k, v := range vf {
		if !val.Flatten() {
			k = fmt.Sprintf("%s.%s", key, k)
		}
		switch vt := v.(type) {
		case HasFields:
			parseFields(k, vt, fields)
		default:
			fields[k] = vt
		}
	}
}
