# flow-eng



### check if input value has been updated and boolCOV will let you know if an input went from false to true
```go
InputUpdated(name InputName) (updated bool, boolCOV bool)
```

### get the flow loop count and if its the first loop
```go
count, firstLoop := inst.Loop()
```

### check if an input has a connection
```go
if inst.InputHasConnection(node.InBoolean) {
    log.Infof("log: comment: %s bool: %t", comment, inBool)
}
```