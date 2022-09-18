package modbusclient

import (
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type InputAddr struct {
	bacnetAddr int
	deviceAddr int
	ioPin      int
	temp       int
	volt       int
	current    int
}

func (inst *Modbus) BuildOutput(ioType bstore.IoType, id bstore.ObjectID) (OutputAddr, model.ObjectType) {
	_, out := outputAddress(0, int(id))

	return out, typeSelect(ioType, true)
}

func (inst *Modbus) BuildInput(ioType bstore.IoType, id bstore.ObjectID) (*InputAddr, model.ObjectType) {
	_, out := InputAddress(0, int(id))
	return out, typeSelect(ioType, true)
}

func typeSelect(objectType bstore.IoType, write bool) model.ObjectType {
	if objectType == bstore.IoTypeVolts {
		if write {
			return model.ObjTypeWriteHolding
		}
		return model.ObjTypeReadRegister
	}
	if objectType == bstore.IoTypeDigital {
		if write {
			return model.ObjTypeWriteCoil
		}
		return model.ObjTypeReadCoil
	}
	if objectType == bstore.IoTypeTemp {
		return model.ObjTypeReadRegister
	}
	return ""

}

func InputAddress(deviceCount int, filterByBacnet int) ([]*InputAddr, *InputAddr) {
	var ioNumber = 1
	var temp = 1
	var volt = 250
	var current = 301
	var count = 0
	var innerCount = 0
	if deviceCount == 0 {
		deviceCount = 4
	}
	ioCount := []int{1, 2, 3, 4, 5, 6, 7, 8}
	ioList := map[int][]int{}
	sum := 0
	for i := 1; i < deviceCount+1; i++ {
		sum++
		ioList[i] = ioCount
	}
	var addresses []*InputAddr
	address := &InputAddr{}
	for _, ints := range ioList {
		count++
		for i := range ints {
			innerCount++
			address.deviceAddr = count
			address.ioPin = i + ioNumber
			address.bacnetAddr = innerCount
			address.temp = i + temp
			address.volt = i + volt
			address.current = i + current
			//fmt.Println("device-addr", count, "bacnet address", innerCount, "io-number", i+ioNumber, "point-tmp", i+temp, "point-volt", i+volt)
			addresses = append(addresses, address)
		}
	}
	filtered := &InputAddr{}
	if filterByBacnet != 0 {
		for _, addr := range addresses {
			if addr.bacnetAddr == filterByBacnet {
				filtered = addr
			}
		}
	}
	return addresses, filtered
}

type OutputAddr struct {
	bacnetAddr int
	deviceAddr int
	ioPin      int
	relay      int
	volt       int
	modbusType model.ObjectType
}

func outputAddress(deviceCount int, filterByBacnet int) ([]OutputAddr, OutputAddr) {
	var ioNumber = 1
	var relay = 1
	var volt = 250
	var count = 0
	var innerCount = 0
	if deviceCount == 0 {
		deviceCount = 4
	}
	ioCount := []int{1, 2, 3, 4, 5, 6, 7, 8}
	ioList := map[int][]int{}
	sum := 0
	for i := 1; i < deviceCount+1; i++ {
		sum++
		ioList[i] = ioCount
	}
	var addresses []OutputAddr
	address := OutputAddr{}
	for _, ints := range ioList {
		count++
		for i := range ints {
			innerCount++
			address.deviceAddr = count
			address.ioPin = i + ioNumber
			address.bacnetAddr = innerCount
			address.relay = i + relay
			address.volt = i + volt
			addresses = append(addresses, address)
		}
	}
	filtered := OutputAddr{}
	if filterByBacnet != 0 {
		for _, addr := range addresses {
			if addr.bacnetAddr == filterByBacnet {
				filtered = addr
			}
		}
	}
	return addresses, filtered
}
