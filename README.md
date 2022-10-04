# flow-eng


## node programming

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
    log.Infof("log: comment: %s bool: %t", comment, inBool)
}
```

### blocking go-routine
when a blocking function is need to be executed eg:(do a database call) a `go-routine` can be used, see example blow

[example](https://github.com/NubeDev/flow-eng/blob/620e14572a55b390c99e4efc2214d20a681423de/nodes/count/ramp.go#L49)