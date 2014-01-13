const char * yaml_event_get_scalar_value(yaml_event_t *event) {
    return (const char *)event->data.scalar.value;
}
