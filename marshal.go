package gamble

// #include "./yaml.h"
// #include "./yaml_extensions.h"
import "C"

import (
	"unsafe"
	"fmt"
)

func Marshal(input interface{}) (result string, err error) {
	var emitter C.yaml_emitter_t
	var output_length C.size_t
	var output_capacity = C.size_t(10000)
	var output *C.uchar = (*C.uchar)(C.malloc(output_capacity))
	C.memset(unsafe.Pointer(output), 0, output_capacity);

	C.yaml_emitter_initialize(&emitter);
	C.yaml_emitter_set_unicode(&emitter, 1);
	C.yaml_emitter_set_indent(&emitter, 2);
	C.yaml_emitter_set_output_string(&emitter, output, output_capacity, &output_length)

	startMarshaling(&emitter)
	marshalNode(&emitter, input)
	stopMarshaling(&emitter)

	C.yaml_emitter_delete(&emitter);
	result = C.GoStringN((*C.char)(unsafe.Pointer(output)), C.int(output_length))
	return
}

func marshalNode(e *C.yaml_emitter_t, node interface{}) {
	switch node := node.(type) {
	case map[string]interface{}:
		marshalMap(e, node)
	case string:
		marshalString(e, node)
	}
}

func marshalMap(e *C.yaml_emitter_t, node map[string]interface{}) {
	var event C.yaml_event_t
	C.yaml_mapping_start_event_initialize(&event, nil, nil, 0, C.YAML_ANY_MAPPING_STYLE)
	emit(e, &event)

	for key, value := range node {
		marshalNode(e, key)
		marshalNode(e, value)
	}

	C.yaml_mapping_end_event_initialize(&event)
	emit(e, &event)
}

func marshalString(e *C.yaml_emitter_t, node string) {
	var event C.yaml_event_t
	C.yaml_scalar_event_initialize(&event, nil, nil, yamlString(node), (C.int)(yamlStringLength(node)), 0, 0, C.YAML_PLAIN_SCALAR_STYLE)
	emit(e, &event)
}

func emit(e *C.yaml_emitter_t, event *C.yaml_event_t) {
	code := C.yaml_emitter_emit(e, event)
	if code != C.int(1) {
		panic("ZOMG FREAKOUT" + C.GoString(e.problem))
	}
}

func startMarshaling(e *C.yaml_emitter_t) {
	var event C.yaml_event_t

	C.yaml_stream_start_event_initialize(&event, C.YAML_UTF8_ENCODING)
	emit(e, &event)

	var version_directive C.yaml_version_directive_t
	version_directive.major = C.int(1)
	version_directive.minor = C.int(2)

	C.yaml_document_start_event_initialize(&event, nil, nil, nil, 0)
	emit(e, &event)
}

func stopMarshaling(e *C.yaml_emitter_t) {
	var event C.yaml_event_t
	C.yaml_stream_end_event_initialize(&event)
	emit(e, &event)

	C.yaml_document_end_event_initialize(&event, 0);
	emit(e, &event)

	code := C.yaml_emitter_flush(e)
	if code != C.int(1) {
		fmt.Println("WTF!")
	}
}
