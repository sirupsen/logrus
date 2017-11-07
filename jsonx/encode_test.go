package jsonx

import (
	"testing"
)

type aStruct struct {
	Channel chan string
	Name string
}

func Test_Client_CRUD(t *testing.T) {

	var myStruct = &aStruct{
		Channel: make(chan string),
		Name: "myles",
	}

	jsonBytes, err := MarshalWithOptions(myStruct, MarshalOptions{SkipUnserializableFields:true})

	if err != nil {
		t.Fatal()
	}

	jsonStr := string(jsonBytes)

	if jsonStr!="{\"Channel\":,\"Name\":\"myles\"}" {
		t.Fatal(jsonStr)
	}

}