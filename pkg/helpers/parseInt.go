package helpers

import "strconv"

// Helper function to change id from string param to uint
// (according to gorm model)
func ParseUint(id string) (uint, error) {
	parsedUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsedUint), nil
}
