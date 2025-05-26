package tsuniqid

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	simpleRand "math/rand"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	generator = NewUniqIDGenerator()
)

func UniqID() string {
	return generator.uniqID()
}

func UniqUID() uint64 {
	return generator.nextID()
}

const (
	maxMachineNum = 0xff
	maxCounterNum = 0x3fff
	maxTimestamp  = 0x3ffffffffff
)

type UniqIDGenerator struct {
	machineID uint64
	counter   uint64
}

func NewUniqIDGenerator() *UniqIDGenerator {
	return &UniqIDGenerator{
		machineID: initMachineID() % maxMachineNum,
	}
}

func (i *UniqIDGenerator) uniqID() string {
	return fmt.Sprintf("%s%s", strconv.FormatUint(i.nextID(), 16), randString(8))
}

func (i *UniqIDGenerator) nextID() uint64 {
	nextCounter := i.nextCounter()
	timestamp := uint64(time.Now().UnixMilli())
	id := (i.machineID << 56) | (timestamp % maxTimestamp << 14) | (nextCounter % maxCounterNum)
	return id
}

func (i *UniqIDGenerator) nextCounter() uint64 {
	return atomic.AddUint64(&i.counter, 1)
}

func initMachineID() uint64 {
	var (
		machineID uint64
		host      string
		ip        string
		err       error
	)

	host, err = os.Hostname()
	if err != nil || host == "" {
		host = randString(10)
	}

	locIP, err := getLocalIP()
	if err != nil {
		ip = randString(10)
	} else {
		ip = locIP.String()
	}

	machineBytes := []byte(hashSha1(host + ip))
	machineBytes = machineBytes[len(machineBytes)-8:]
	_ = binary.Read(bytes.NewBuffer(machineBytes), binary.BigEndian, &machineID)

	if machineID == 0 {
		machineID = randUint64()
	}
	return machineID
}

var (
	newSimpleRand = simpleRand.New(simpleRand.NewSource(time.Now().UnixNano()))
	letters       = "1234567890abcde"
	lettersLen    = 15
)

func randString(len int) string {
	var s []string
	b := new(big.Int).SetInt64(int64(lettersLen))
	for i := 0; i < len; i++ {
		if i, err := rand.Int(rand.Reader, b); err == nil {
			s = append(s, string(letters[i.Int64()]))
		}
	}
	return strings.Join(s, "")
}

func randUint64() uint64 {
	r, err := rand.Int(rand.Reader, new(big.Int).SetUint64(math.MaxUint64))
	if err != nil {
		return newSimpleRand.Uint64()
	}
	return r.Uint64()
}

func hashSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
