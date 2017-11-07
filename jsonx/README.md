### JSONX

A copy of the json package shipped with Go 1.6.2 with the addition of ```MarshalOptions```.  See ```MarshalWithOptions``` function in encode file.

Currently there is a single option: ```SkipUnserializableFields``` which will do just that.  See [encode_test.go](encode_test.go) for details.

