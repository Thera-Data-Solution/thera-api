package utils

import "github.com/lucsky/cuid"

// GenerateID menghasilkan ID unik seperti "cmh0pr9j40000onyr2jgzlj3x"
func GenerateID() string {
	return cuid.New()
}
