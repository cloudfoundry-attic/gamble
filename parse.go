package gamble

// #include "./yaml.h"
// #include "./yaml_extensions.h"
import "C"

import (
	"errors"
)

type Node interface{}

func Parse(input string) (result Node, err error) {
	var parser C.yaml_parser_t
	C.yaml_parser_initialize(&parser)
	C.yaml_parser_set_input_string(&parser, (*C.uchar)(yamlString(input)), yamlStringLength(input))

	defer func() {
		if e, ok := recover().(error); ok {
			err = e
		}
	}()

	result = getNode(&parser, C.YAML_DOCUMENT_END_EVENT)
	return
}

func getNode(p *C.yaml_parser_t, stopEvent C.yaml_event_type_t) Node {
	for {
		var event C.yaml_event_t
		if int(C.yaml_parser_parse(p, &event)) == 0 {
			panic(errors.New("Error parsing YAML."))
		}

		eventType := event._type
		if eventType == stopEvent || eventType == 0 {
			return nil
		}

		switch eventType {
		case C.YAML_SCALAR_EVENT:
			str := C.GoString(C.yaml_event_get_scalar_value(&event))
			isPlainScalar := C.yaml_event_get_scalar_style(&event) == C.YAML_PLAIN_SCALAR_STYLE
			if isPlainScalar && (str == "null" || str == "") {
				return nil
			}
			return str
		case C.YAML_SEQUENCE_START_EVENT:
			sequenceNode := make([]interface{}, 0)
			for {
				item := getNode(p, C.YAML_SEQUENCE_END_EVENT)
				if item != nil {
					sequenceNode = append(sequenceNode, item)
				} else {
					break
				}
			}
			return sequenceNode
		case C.YAML_MAPPING_START_EVENT:
			mappingNode := make(map[string]interface{})
			for {
				key := getNode(p, C.YAML_MAPPING_END_EVENT)
				if key != nil {
					value := getNode(p, C.YAML_MAPPING_END_EVENT)
					mappingNode[key.(string)] = value
				} else {
					break
				}
			}
			return mappingNode
		}
	}
}
