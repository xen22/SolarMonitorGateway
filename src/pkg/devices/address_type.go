package devices

import "strconv"

// AddressType type
type AddressType byte

// UnmarshalJSON converts string to AddressType.
func (a *AddressType) UnmarshalJSON(value []byte) (err error) {
	hex, err := strconv.ParseInt(string(value), 0, 8)
	if err != nil {
		return
	}
	*a = AddressType(hex)
	return
}
