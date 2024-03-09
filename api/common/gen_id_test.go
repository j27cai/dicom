package common

import (
    "testing"
)

func TestGenShortUUID(t *testing.T) {
    // Generate a short UUID
    uuid := genShortUUID()

    // Check if the generated UUID is not empty
    if uuid == "" {
        t.Errorf("Generated UUID is empty")
    }
}