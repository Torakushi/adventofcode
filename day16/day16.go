package day16

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const input = ""

var hexMap = map[rune]string{
	'0': "0000",
	'1': "0001",
	'2': "0010",
	'3': "0011",
	'4': "0100",
	'5': "0101",
	'6': "0110",
	'7': "0111",
	'8': "1000",
	'9': "1001",
	'A': "1010",
	'B': "1011",
	'C': "1100",
	'D': "1101",
	'E': "1110",
	'F': "1111",
}

func Day16() {
	fmt.Println("DAY16")

	if err := process(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}

func process() error {
	b, err := ioutil.ReadFile("day16/data.txt")
	if err != nil {
		return err
	}
	input := string(b)

	h := getHexBinaryString(input)
	packetObject, err := parsePacketObject(h)
	if err != nil {
		return err
	}

	fmt.Printf("First Part: total versions sum is: %d\n", getSumPacketVersion(packetObject))
	fmt.Printf("Second Part: Final value after operations: %d\n", packetObject.getValue())

	return nil
}

func getSumPacketVersion(p packetsInterface) int64 {
	sum := p.getPacketVersion()
	for _, c := range p.getSubItems() {
		sum += getSumPacketVersion(c)
	}
	return sum
}

func getHexBinaryString(s string) string {
	var sb strings.Builder
	for _, r := range s {
		sb.WriteString(hexMap[r])
	}
	return sb.String()
}

func getPacketVersion(s string) (int64, error) {
	return strconv.ParseInt(s[:3], 2, 64)
}

func getPacketTypeID(s string) (int64, error) {
	return strconv.ParseInt(s[3:6], 2, 64)
}

func getLengthTypeID(s string) (int64, error) {
	return strconv.ParseInt(string(s[6]), 2, 64)
}

func getSubPacketsBitLength(s string) (int64, error) {
	return strconv.ParseInt(string(s[7:22]), 2, 64)
}

func getSubPacketsNumber(s string) (int64, error) {
	return strconv.ParseInt(string(s[7:18]), 2, 64)
}

func binaryAsInt(s string) (int64, error) {
	return strconv.ParseInt(s, 2, 64)
}

func parsePacketObject(s string) (packetsInterface, error) {
	t, err := getPacketTypeID(s)
	if err != nil {
		return nil, err
	}

	if t == 4 {
		return parseSubpacket(s)
	}

	lt, err := getLengthTypeID(s)
	if err != nil {
		return nil, err
	}

	if lt == 1 {
		return parseOperationLengthTypeOne(s)
	}

	return parseOperationLengthTypeZero(s)
}

func parseSubpacket(s string) (*subPacket, error) {
	sp := &subPacket{}
	var err error
	sp.packetVersion, err = getPacketVersion(s)
	if err != nil {
		return nil, err
	}

	sp.packetTypeID, err = getPacketTypeID(s)
	if err != nil {
		return nil, err
	}

	lastLabelBit := "1"
	index := 6
	var sb strings.Builder
	for lastLabelBit != "0" {
		lastLabelBit = string(s[index])
		sb.WriteString(s[index+1 : index+5])
		index += 5
	}

	sp.bitLength = int64(index)

	lv, err := binaryAsInt(sb.String())
	if err != nil {
		return nil, err
	}
	sp.literalValue = lv
	return sp, nil
}

func parseOperationLengthTypeZero(s string) (*operationLengthTypeZero, error) {
	o := &operationLengthTypeZero{}
	var err error
	o.packetVersion, err = getPacketVersion(s)
	if err != nil {
		return nil, err
	}

	o.packetTypeID, err = getPacketTypeID(s)
	if err != nil {
		return nil, err
	}

	o.f = getFuncById(int(o.packetTypeID))

	o.lengthTypeID, err = getLengthTypeID(s)
	if err != nil {
		return nil, err
	}

	o.subPacketsBitLength, err = getSubPacketsBitLength(s)
	if err != nil {
		return nil, err
	}

	index := 22
	arr := []packetsInterface{}
	for index < 22+int(o.subPacketsBitLength) {
		p, err := parsePacketObject(s[index:])
		if err != nil {
			return nil, err
		}
		arr = append(arr, p)
		index += int(p.getLength())
	}

	o.subItems = arr
	o.bitLength = int64(index)

	iArr := []int64{}
	for i := 0; i < len(arr); i++ {
		iArr = append(iArr, arr[i].getValue())
	}
	o.opValue = o.f(iArr)

	return o, nil
}

func parseOperationLengthTypeOne(s string) (*operationLengthTypeOne, error) {
	o := &operationLengthTypeOne{}
	var err error
	o.packetVersion, err = getPacketVersion(s)
	if err != nil {
		return nil, err
	}

	o.packetTypeID, err = getPacketTypeID(s)
	if err != nil {
		return nil, err
	}

	o.f = getFuncById(int(o.packetTypeID))

	o.lengthTypeID, err = getLengthTypeID(s)
	if err != nil {
		return nil, err
	}

	o.subPacketsNumber, err = getSubPacketsNumber(s)
	if err != nil {
		return nil, err
	}

	index := 18
	arr := []packetsInterface{}
	i := 0
	for i < int(o.subPacketsNumber) {
		p, err := parsePacketObject(s[index:])
		if err != nil {
			return nil, err
		}
		arr = append(arr, p)
		index += int(p.getLength())
		i++
	}

	o.subItems = arr
	o.bitLength = int64(index)

	iArr := []int64{}
	for i := 0; i < len(arr); i++ {
		iArr = append(iArr, arr[i].getValue())
	}
	o.opValue = o.f(iArr)

	return o, nil
}

func getFuncById(i int) func([]int64) int64 {
	switch i {
	case 0:
		return func(a []int64) int64 {
			s := int64(0)
			for _, v := range a {
				s += v
			}
			return s
		}
	case 1:
		return func(a []int64) int64 {
			p := int64(1)
			for _, v := range a {
				p *= v
			}
			return p
		}
	case 2:
		return func(a []int64) int64 {
			min := a[0]
			for _, v := range a {
				if v < min {
					min = v
				}
			}
			return min
		}
	case 3:
		return func(a []int64) int64 {
			max := a[0]
			for _, v := range a {
				if v > max {
					max = v
				}
			}
			return max
		}

	case 5:
		return func(a []int64) int64 {
			if a[0] > a[1] {
				return 1
			}
			return 0
		}

	case 6:
		return func(a []int64) int64 {
			if a[0] < a[1] {
				return 1
			}
			return 0
		}
	case 7:
		return func(a []int64) int64 {
			if a[0] == a[1] {
				return 1
			}
			return 0
		}
	default:
		panic("ID FUNC > 7")
	}
}

type packetsInterface interface {
	getPacketVersion() int64
	getLength() int64
	getValue() int64
	getSubItems() []packetsInterface
}

type operationLengthTypeOne struct {
	packetVersion    int64
	packetTypeID     int64
	lengthTypeID     int64
	subPacketsNumber int64
	bitLength        int64
	subItems         []packetsInterface
	opValue          int64
	f                func([]int64) int64
}

type operationLengthTypeZero struct {
	packetVersion       int64
	packetTypeID        int64
	lengthTypeID        int64
	subPacketsBitLength int64
	bitLength           int64
	subItems            []packetsInterface
	opValue             int64
	f                   func([]int64) int64
}

type subPacket struct {
	packetVersion int64
	packetTypeID  int64
	literalValue  int64
	bitLength     int64
}

func (sp *subPacket) getLength() int64 {
	return sp.bitLength
}

func (o *operationLengthTypeZero) getLength() int64 {
	return o.bitLength
}

func (o *operationLengthTypeOne) getLength() int64 {
	return o.bitLength
}
func (sp *subPacket) getValue() int64 {
	return sp.literalValue
}

func (o *operationLengthTypeZero) getValue() int64 {
	return o.opValue
}

func (o *operationLengthTypeOne) getValue() int64 {
	return o.opValue
}

func (sp *subPacket) getPacketVersion() int64 {
	return sp.packetVersion
}

func (o *operationLengthTypeZero) getPacketVersion() int64 {
	return o.packetVersion
}

func (o *operationLengthTypeOne) getPacketVersion() int64 {
	return o.packetVersion
}

func (sp *subPacket) getSubItems() []packetsInterface {
	return nil
}

func (o *operationLengthTypeZero) getSubItems() []packetsInterface {
	return o.subItems
}

func (o *operationLengthTypeOne) getSubItems() []packetsInterface {
	return o.subItems
}
