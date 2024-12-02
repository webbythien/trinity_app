package task

import "fmt"

func HealthCheck(input int64) error {
	fmt.Printf("HealthCheck Number: %d\n", input)
	return nil
}
