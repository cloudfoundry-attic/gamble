package gamble

/*
#include "./yaml.h"
extern int writeHandler(void *data, unsigned char *buffer, size_t size);
*/
import "C"

import (
	"fmt"
	"unsafe"
	"bytes"
	"strconv"
)

type MarshalError struct {
	Node interface{}
}

func (err MarshalError) Error() string {
	return fmt.Sprintf("error marshaling unknown type %T", err.Node)
}

func Marshal(input interface{}) (result string, err error) {
	defer func() {
		panicd := recover()
		var ok bool
		err, ok = panicd.(MarshalError)
		if ok {
			return
		} else {
			panic(panicd)
		}
	}()

	var emitter C.yaml_emitter_t

	C.yaml_emitter_initialize(&emitter);
	C.yaml_emitter_set_unicode(&emitter, 1);
	C.yaml_emitter_set_indent(&emitter, 2);

	var buffer bytes.Buffer
	C.yaml_emitter_set_output(&emitter, (*[0]byte)(C.writeHandler), unsafe.Pointer(&buffer))

	startMarshaling(&emitter)
	marshalNode(&emitter, input)
	stopMarshaling(&emitter)

	C.yaml_emitter_delete(&emitter);
	result = buffer.String()
	return
}

func marshalNode(e *C.yaml_emitter_t, node interface{}) {
	if node == nil {
		marshalNull(e)
		return
	}

	switch value := node.(type) {
	case map[string]interface{}:
		marshalMap(e, value)
	case []interface{}:
		marshalSlice(e, value)
	case string:
		marshalString(e, value)
	case int:
		marshalInt(e, int64(value))
	case int8:
		marshalInt(e, int64(value))
	case int16:
		marshalInt(e, int64(value))
	case int32:
		marshalInt(e, int64(value))
	case int64:
		marshalInt(e, int64(value))
	case uint:
		marshalUint(e, uint64(value))
	case uint8:
		marshalUint(e, uint64(value))
	case uint16:
		marshalUint(e, uint64(value))
	case uint32:
		marshalUint(e, uint64(value))
	case uint64:
		marshalUint(e, uint64(value))
	case float32:
		marshalFloat(e, float64(value))
	case float64:
		marshalFloat(e, float64(value))
	default:
		panic(MarshalError{Node: node})
	}
}

func marshalMap(e *C.yaml_emitter_t, node map[string]interface{}) {
	var event C.yaml_event_t
	C.yaml_mapping_start_event_initialize(&event, nil, nil, 0, C.YAML_ANY_MAPPING_STYLE)
	emit(e, &event)

	for key, value := range node {
		marshalString(e, key)
		marshalNode(e, value)
	}

	C.yaml_mapping_end_event_initialize(&event)
	emit(e, &event)
}

func marshalSlice(e *C.yaml_emitter_t, node []interface{}) {
	var event C.yaml_event_t
	C.yaml_sequence_start_event_initialize(&event, nil, nil, 0, C.YAML_ANY_SEQUENCE_STYLE)
	emit(e, &event)

	for _, value := range node {
		marshalNode(e, value)
	}

	C.yaml_sequence_end_event_initialize(&event)
	emit(e, &event)
}

func marshalUint(e *C.yaml_emitter_t, value uint64) {
	marshalString(e, strconv.FormatUint(value, 10))
}

func marshalInt(e *C.yaml_emitter_t, value int64) {
	marshalString(e, strconv.FormatInt(value, 10))
}

func marshalFloat(e *C.yaml_emitter_t, value float64) {
	marshalString(e, strconv.FormatFloat(value, 'f', 2, 64))
}

func marshalString(e *C.yaml_emitter_t, node string) {
	var event C.yaml_event_t
	C.yaml_scalar_event_initialize(&event, nil, nil, yamlString(node), (C.int)(yamlStringLength(node)), 1, 1, C.YAML_PLAIN_SCALAR_STYLE)
	emit(e, &event)
}

func marshalNull(e *C.yaml_emitter_t) {
	var event C.yaml_event_t
	C.yaml_scalar_event_initialize(&event, nil, nil, yamlString("null"), 4, 1, 1, C.YAML_PLAIN_SCALAR_STYLE)
	emit(e, &event)
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
	C.yaml_document_end_event_initialize(&event, 0);
	emit(e, &event)

	C.yaml_stream_end_event_initialize(&event)
	emit(e, &event)

	code := C.yaml_emitter_flush(e)
	if code != C.int(1) {
		panic(fmt.Sprintf("YAML error flushing: %s", C.GoString(e.problem)))
	}
}

func emit(e *C.yaml_emitter_t, event *C.yaml_event_t) {
	code := C.yaml_emitter_emit(e, event)
	if code != C.int(1) {
		panic(fmt.Sprintf("YAML error emitting %s", C.GoString(e.problem)))
	}
}

//export writeHandler
func writeHandler(data unsafe.Pointer, buffer *C.uchar, size C.size_t) C.int {
	writer := *((*bytes.Buffer)(data))
	str := C.GoStringN((*C.char)(unsafe.Pointer(buffer)), C.int(size))
	writer.Write(([]byte)(str))
	return C.int(1)
}
