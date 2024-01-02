package points

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/datatype"
)

type ModbusInputAddr struct {
	BacnetAddr int `json:"bacnetAddr"`
	// DeviceAddr int `json:"modbusAddr"`
	IoPin   int `json:"ioPin"`
	Temp    int `json:"temp"`
	Volt    int `json:"volt"`
	Current int `json:"current"`
}

func ModbusBuildOutput(ioType IoType, id ObjectID) (OutputAddr, datatype.ObjectType) {
	_, out := outputAddress(0, int(id))
	return out, typeSelect(ioType, true)
}

func ModbusBuildInput(ioType IoType, id ObjectID) (ModbusInputAddr, datatype.ObjectType) {
	_, out := InputAddress(0, int(id))
	return out, typeSelect(ioType, false)
}

func typeSelect(objectType IoType, write bool) datatype.ObjectType {
	if objectType == IoTypeVolts {
		if write {
			return datatype.ObjTypeWriteHolding
		}
		return datatype.ObjTypeReadRegister
	}
	if objectType == IoTypeDigital {
		if write {
			return datatype.ObjTypeWriteCoil
		}
		return datatype.ObjTypeReadCoil
	}
	if objectType == IoTypeTemp {
		return datatype.ObjTypeReadRegister
	}
	if objectType == IoTypePulseOnFall {
		return datatype.ObjTypeReadRegister
	}
	if objectType == IoTypePulseOnRise {
		return datatype.ObjTypeReadRegister
	}
	return ""

}

func InputAddress(deviceCount int, filterByBacnet int) ([]ModbusInputAddr, ModbusInputAddr) {
	var ioNumber = 1
	var temp = 1
	var volt = 201
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
	var addresses []ModbusInputAddr
	address := ModbusInputAddr{}
	for _, ints := range ioList {
		count++
		for i := range ints {
			innerCount++
			// address.DeviceAddr = count
			address.IoPin = i + ioNumber
			address.BacnetAddr = innerCount
			address.Temp = i + temp
			address.Volt = i + volt
			address.Current = i + current
			// fmt.Println("device-addr", count, "bacnet address", innerCount, "io-number", i+ioNumber, "point-tmp", i+temp, "point-volt", i+volt)
			addresses = append(addresses, address)
		}
	}
	filtered := ModbusInputAddr{}
	if filterByBacnet != 0 {
		for _, addr := range addresses {
			if addr.BacnetAddr == filterByBacnet {
				filtered = addr
			}
		}
	}
	return addresses, filtered
}

type OutputAddr struct {
	BacnetAddr int `json:"bacnetAddr"`
	// DeviceAddr int `json:"modbusAddr"`
	IoPin int `json:"ioPin"`
	Relay int `json:"relay"`
	Volt  int `json:"volt"`
}

func outputAddress(deviceCount int, filterByBacnet int) ([]OutputAddr, OutputAddr) {
	var ioNumber = 1
	var relay = 1
	var volt = 0
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
			// address.DeviceAddr = count
			address.IoPin = i + ioNumber
			address.BacnetAddr = innerCount
			address.Relay = i + relay
			address.Volt = i + volt
			addresses = append(addresses, address)
		}
	}
	filtered := OutputAddr{}
	if filterByBacnet != 0 {
		for _, addr := range addresses {
			if addr.BacnetAddr == filterByBacnet {
				filtered = addr
			}
		}
	}
	return addresses, filtered
}
