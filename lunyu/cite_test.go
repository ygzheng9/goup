package lunyu

import (
	"fmt"
	"testing"
)

func TestIntersect(t *testing.T) {
	old := []string{"a", "b", "c"}
	curr := []string{"b", "d"}

	added, removed := Intersect(old, curr)

	fmt.Printf("added: %s\n", added)
	fmt.Printf("removed: %s\n", removed)
}

func TestMinus(t *testing.T) {
	old := []string{"a", "b", "c"}
	curr := []string{"b", "d"}

	removed := Minus(old, curr)

	fmt.Printf("removed: %s\n", removed)
}
