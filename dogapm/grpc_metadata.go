package dogapm

import "google.golang.org/grpc/metadata"
type metadataSupplier struct {
	metadata *metadata.MD
}

func (s *metadataSupplier) Get(key string) string {
	if s.metadata == nil {
		return ""
	}
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key, value string) {
	if s.metadata == nil {
		s.metadata = &metadata.MD{}
	}
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	if s.metadata == nil {
		return nil
	}
	keys := make([]string, 0, len(*s.metadata))
	for k := range *s.metadata {
		keys = append(keys, k)
	}
	return keys
}

