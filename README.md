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