#ifndef YAML_EXTENSIONS_H
#define YAML_EXTENSIONS_H

#include "./yaml.h"

static const char * yaml_event_get_scalar_value(yaml_event_t *event) {
    return (const char *)event->data.scalar.value;
}

static yaml_scalar_style_t yaml_event_get_scalar_style(yaml_event_t *event) {
    return event->data.scalar.style;
}

#endif
