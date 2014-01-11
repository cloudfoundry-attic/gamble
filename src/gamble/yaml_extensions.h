yaml_event_type_t yaml_event_get_type(yaml_event_t *event) {
    return event->type;
}

const char * yaml_event_get_scalar_value(yaml_event_t *event) {
    return (const char *)event->data.scalar.value;
}
