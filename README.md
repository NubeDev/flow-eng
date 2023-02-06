# flow-eng

## node payload

using this api you can write a value to a node without having to download/deploy a new flow

`method` HTTP-POST

`url`
```
http://0.0.0.0:1665/api/nodes/payload/NODE-ID
```

body a float value

```json
{
  "float": 2
}
```

body a bool value

```json
{
  "boolean": true
}
```

body a string value

```json
{
  "string": "abc"
}
```
body as an interface
```json
{
  "any": [1,2,3,4]
}
```

## null values

```go
func (inst *Not) Process() {
in, null := inst.ReadPinBoolOk(node.Inp)
if null { // if input is null then set output to nil
inst.WritePinNull(node.Outp)
return
}
if in {
inst.WritePinFalse(node.Outp)
} else {
inst.WritePinTrue(node.Outp)
}
}
```

### reading inputs

```go
inst.ReadPinFloat(name InputName) (value float46, null bool) // will return the value as a float and if its `null/nil` the `boolean` `null` flag will be `true` 
inst.ReadPinBool(name InputName) (value bool, null bool) // same as above but value is a boolean
```

### writing outputs

```go
inst.WritePinFloat(node.Outp, f, inst.precision) // you can add number of decimal places 
inst.WritePinFloat(name OutputName, value float64) // set output to number value (float64)
inst.WritePinNull(node.Outp) // set output to nil
inst.WritePinTrue(node.Outp) // set output to true
inst.WritePinFalse(node.Outp) // set output to false
```

[interface](https://github.com/NubeDev/flow-eng/blob/f8778ee7402691a75516acdb9eef355038c8b17a/node/node.go#L7)

### check if input value has been updated and boolCOV will let you know if an input went from false to true

```go
InputUpdated(name InputName) (updated bool, boolCOV bool)
```

### get the flow loop count and if its the first loop

```go
count, firstLoop := Loop()
```

### check if an input has a connection

```go
if InputHasConnection(node.InBoolean) {
log.Infof("log: comment: %s boolean: %t", comment, inBool)
}
```

### blocking go-routine

when a blocking function is need to be executed eg:(do a database call) a `go-routine` can be used, see example below

[example](https://github.com/NubeDev/flow-eng/blob/620e14572a55b390c99e4efc2214d20a681423de/nodes/count/ramp.go#L49)