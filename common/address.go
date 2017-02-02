package common

import (
    "crypto/rand"
    "github.com/ethereum/go-ethereum/common"
)

// RandomAddress - Generate a random address.
func RandomAddress() (*common.Address, error) {
    bytes := make([]byte, common.AddressLength)
    _, err := rand.Read(bytes)
    if err != nil {
        return nil, err
    }
    address := common.BytesToAddress(bytes)
    return &address, nil
}

// RandomAddresses - Generate a list of random address.
func RandomAddresses(n int) ([]common.Address, error) {
    addresses := make([]common.Address, n)
    for i := range addresses {
        ptr, err := RandomAddress()
        addresses[i] = *ptr
        if err != nil {
            return nil, err
        }
    }
    return addresses, nil
}

func sortByHex(addresses []common.Address, l int, r int) {
    if l < r {
        pivot := addresses[(l + r) / 2].Hex()
        i := l
        j := r
        var tmp common.Address
        for i <= j {
            for addresses[i].Hex() < pivot { i++ }
            for addresses[j].Hex() > pivot { j-- }
            if i <= j {
                tmp = addresses[i]
                addresses[i] = addresses[j]
                addresses[j] = tmp
                i++
                j--
            }
        }
        if l < j {
            sortByHex(addresses, l, j)
        }
        if i < r {
            sortByHex(addresses, i, r)
        }
    }
}

// SortAddresses - Sort a list of address.
func SortAddresses(addresses []common.Address) {
    n := len(addresses)
    sortByHex(addresses, 0, n - 1)
}
