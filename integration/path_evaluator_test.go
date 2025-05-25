package integration

import (
	"fmt"
	"testing"
)

func TestPathEvaluator(t *testing.T) {
	fmt.Println(EvaluatePath("./scripts/script1/"))
	fmt.Println(EvaluatePath("./scripts/script1"))
	fmt.Println(EvaluatePath("./scripts/script1.go"))
	fmt.Println(EvaluatePath("./"))
	fmt.Println(EvaluatePath("./scripts/script1/script7"))
	fmt.Println(EvaluatePath("./scripts/script1/script7/test.go"))
}
