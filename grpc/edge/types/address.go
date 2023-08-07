package types

import (
	"fmt"
	"unicode"

	"github.com/0xPolygon/polygon-edge/helper/keccak"
	"github.com/0xPolygonHermez/zkevm-node/hex"
)

const AddressLength = 20

type Address [AddressLength]byte

func (a Address) Bytes() []byte {
	return a[:]
}

// checksumEncode returns the checksummed address with 0x prefix, as by EIP-55
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md
func (a Address) checksumEncode() string {
	addrBytes := a.Bytes() // 20 bytes

	// Encode to hex without the 0x prefix
	lowercaseHex := hex.EncodeToHex(addrBytes)[2:]
	hashedAddress := hex.EncodeToHex(keccak.Keccak256(nil, []byte(lowercaseHex)))[2:]

	result := make([]rune, len(lowercaseHex))
	// Iterate over each character in the lowercase hex address
	for idx, ch := range lowercaseHex {
		if ch >= '0' && ch <= '9' || hashedAddress[idx] >= '0' && hashedAddress[idx] <= '7' {
			// Numbers in range [0, 9] are ignored (as well as hashed values [0, 7]),
			// because they can't be uppercased
			result[idx] = ch
		} else {
			// The current character / hashed character is in the range [8, f]
			result[idx] = unicode.ToUpper(ch)
		}
	}

	return "0x" + string(result)
}

func (a Address) String() string {
	return a.checksumEncode()
}

func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Address) UnmarshalText(input []byte) error {
	buf := stringToBytes(string(input))
	if len(buf) != AddressLength {
		return fmt.Errorf("incorrect length")
	}

	*a = BytesToAddress(buf)

	return nil
}
