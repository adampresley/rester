package contenttype

type ContentTypeHandler func(body []byte, result any) error
