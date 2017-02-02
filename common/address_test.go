package common

import (
    "math/rand"
    "testing"
    "time"
)

func TestSortAddresses(test *testing.T) {
    rand.Seed(time.Now().UnixNano())
    n := rand.Intn(10)
    addresses, err := RandomAddresses(n)
    if (err != nil) {
        test.Fatal(err)
    }
    SortAddresses(addresses)
    for i := 0; i < n - 1; i++ {
        if (addresses[i].Hex() > addresses[i + 1].Hex()) {
            test.Fatal(addresses)
        }
    }
}
