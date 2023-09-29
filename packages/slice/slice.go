package slice

import "fmt"

func Slice() {
	// Creating a slice
	numbers := []int{1, 2, 3, 4, 5}

	// Iterating over a slice
	for index, value := range numbers {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}

	// Modifying a slice
	numbers[0] = 10

	// Appending to a slice
	numbers = append(numbers, 6)

	// Slicing a slice
	subset := numbers[2:4]

	fmt.Println("Modified Slice:", numbers)
	fmt.Println("Subset Slice:", subset)
}
