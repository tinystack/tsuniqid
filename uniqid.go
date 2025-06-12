// Package tsuniqid provides a high-performance unique ID generator
// that generates both string and uint64 type unique identifiers.
//
// The generated IDs are composed of:
// - Machine ID (8 bits): Unique identifier for the machine/process
// - Timestamp (42 bits): Current Unix timestamp in milliseconds
// - Counter (14 bits): Atomic counter to ensure uniqueness within the same millisecond
// - Random suffix (for string IDs): Additional randomness for string-based IDs
package tsuniqid

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Bit allocation constants for the unique ID generation
const (
	// MaxMachineID represents the maximum machine ID value (4 bits)
	MaxMachineID = 0xf

	// MaxInstanceID represents the maximum instance ID value (4 bits)
	MaxInstanceID = 0xf

	// MaxCounter represents the maximum counter value (14 bits)
	MaxCounter = 0x3fff

	// MaxTimestamp represents the maximum timestamp value (42 bits)
	MaxTimestamp = 0x3ffffffffff

	// RandomSuffixLength is the length of random suffix for string IDs
	RandomSuffixLength = 8

	// CharSet contains characters used for random string generation
	CharSet = "0123456789abcdefghijklmnopqrstuvwxyz"

	// TimestampShift is the number of bits to shift timestamp
	TimestampShift = 14

	// InstanceIDShift is the number of bits to shift instance ID
	InstanceIDShift = 56

	// MachineIDShift is the number of bits to shift machine ID
	MachineIDShift = 60
)

// globalInstanceCounter is used to assign unique instance IDs to each generator
var globalInstanceCounter uint64

// Generator is the default global generator instance
var Generator = NewGenerator()

// UniqID generates a unique string ID using the default generator.
// The string ID consists of a hex-encoded uint64 ID plus a random suffix.
//
// Returns: A unique string identifier
func UniqID() string {
	return Generator.GenerateStringID()
}

// UniqUID generates a unique uint64 ID using the default generator.
// The uint64 ID is composed of machine ID, timestamp, and counter.
//
// Returns: A unique uint64 identifier
func UniqUID() uint64 {
	return Generator.GenerateUint64ID()
}

// IDGenerator is responsible for generating unique identifiers.
// It maintains machine ID, instance ID and an atomic counter to ensure uniqueness.
type IDGenerator struct {
	machineID  uint64     // 4-bit machine identifier
	instanceID uint64     // 4-bit instance identifier for distinguishing multiple generators
	counter    uint64     // atomic counter for uniqueness within the same millisecond
	rng        *rand.Rand // local random number generator for better performance
	mu         sync.Mutex // mutex to protect rng from concurrent access
}

// NewGenerator creates a new IDGenerator instance with initialized machine ID and unique instance ID.
//
// Returns: A new IDGenerator instance
func NewGenerator() *IDGenerator {
	// Initialize with current time as seed for better randomness
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Assign a unique instance ID to this generator
	instanceID := atomic.AddUint64(&globalInstanceCounter, 1) & MaxInstanceID

	return &IDGenerator{
		machineID:  generateMachineID() & MaxMachineID, // Ensure within 6-bit range
		instanceID: instanceID,                         // Ensure within 2-bit range
		counter:    0,
		rng:        rng,
	}
}

// GenerateStringID creates a unique string identifier.
// Format: hex(uint64_id) + random_suffix
//
// Returns: A unique string identifier
func (g *IDGenerator) GenerateStringID() string {
	id := g.GenerateUint64ID()
	suffix := g.generateRandomSuffix(RandomSuffixLength)
	return fmt.Sprintf("%s%s", strconv.FormatUint(id, 16), suffix)
}

// GenerateUint64ID creates a unique uint64 identifier.
//
// Bit layout (64 bits total):
// - Bits 63-60 (4 bits): Machine ID
// - Bits 59-56 (4 bits): Instance ID
// - Bits 55-14 (42 bits): Timestamp (milliseconds since Unix epoch)
// - Bits 13-0 (14 bits): Counter
//
// Returns: A unique uint64 identifier
func (g *IDGenerator) GenerateUint64ID() uint64 {
	counter := g.nextCounter()
	timestamp := uint64(time.Now().UnixMilli())

	// Combine components with bit shifting
	id := (g.machineID << MachineIDShift) |
		(g.instanceID << InstanceIDShift) |
		((timestamp & MaxTimestamp) << TimestampShift) |
		(counter & MaxCounter)

	return id
}

// nextCounter atomically increments and returns the next counter value.
//
// Returns: The next counter value
func (g *IDGenerator) nextCounter() uint64 {
	return atomic.AddUint64(&g.counter, 1)
}

// generateRandomSuffix creates a random string of specified length.
// Uses a more efficient approach than crypto/rand for non-cryptographic purposes.
// This method is thread-safe.
//
// Parameters:
//   - length: The desired length of the random string
//
// Returns: A random string of the specified length
func (g *IDGenerator) generateRandomSuffix(length int) string {
	if length <= 0 {
		return ""
	}

	result := make([]byte, length)
	charSetLen := len(CharSet)

	// Lock to ensure thread-safe access to the random number generator
	g.mu.Lock()
	for i := 0; i < length; i++ {
		result[i] = CharSet[g.rng.Intn(charSetLen)]
	}
	g.mu.Unlock()

	return string(result)
}

// generateMachineID creates a unique machine identifier based on hostname and local IP.
// If hostname or IP cannot be obtained, it falls back to random generation.
//
// Returns: A machine-specific identifier
func generateMachineID() uint64 {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = generateFallbackString(10)
	}

	// Get local IP
	localIP, err := getLocalIP()
	var ipStr string
	if err != nil {
		ipStr = generateFallbackString(10)
	} else {
		ipStr = localIP.String()
	}

	// Create machine ID from hostname and IP
	return hashToUint64(hostname + ipStr)
}

// hashToUint64 converts a string to uint64 using SHA1 hash.
//
// Parameters:
//   - input: The string to hash
//
// Returns: A uint64 representation of the hash
func hashToUint64(input string) uint64 {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)

	// Use the last 8 bytes of the hash for uint64 conversion
	if len(hashBytes) >= 8 {
		return binary.BigEndian.Uint64(hashBytes[len(hashBytes)-8:])
	}

	// Fallback: pad with zeros if hash is somehow shorter
	padded := make([]byte, 8)
	copy(padded[8-len(hashBytes):], hashBytes)
	return binary.BigEndian.Uint64(padded)
}

// generateFallbackString creates a random string for fallback purposes.
//
// Parameters:
//   - length: The desired length of the random string
//
// Returns: A random string of the specified length
func generateFallbackString(length int) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = CharSet[rng.Intn(len(CharSet))]
	}

	return string(result)
}
