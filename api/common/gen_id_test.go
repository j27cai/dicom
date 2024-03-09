package common

import (
    "testing"
)

func TestGenShortUUID(t *testing.T) {
    uuid := GenShortUUID()

    if uuid == "" {
        t.Errorf("Generated UUID is empty")
    }
}