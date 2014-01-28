package gamble

// #include "./yaml.h"
import "C"

import (
	"unsafe"
)

func yamlString(input string) *C.yaml_char_t {
	return (*C.yaml_char_t)(unsafe.Pointer(C.CString(input)))
}

func yamlStringLength(input string) C.size_t {
	return (C.size_t)(len(input))
}
