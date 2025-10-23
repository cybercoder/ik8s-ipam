package netutils

import (
	"crypto/sha1"
	"fmt"
)

func GenerateStableMAC(iface, namespace, name string) string {
	// Combine the identifiers into a unique string
	key := fmt.Sprintf("%s-%s-%s", iface, namespace, name)

	// Hash the key (SHA1 gives us a stable 20-byte digest)
	hash := sha1.Sum([]byte(key))

	// Use the last 5 bytes of the hash for uniqueness
	// and set the first byte as a "locally administered unicast" MAC
	mac := []byte{
		0x02, // Locally administered MAC (not globally unique)
		hash[0],
		hash[1],
		hash[2],
		hash[3],
		hash[4],
	}

	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}
