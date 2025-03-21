package kafka

import (
	"errors"
	"hash/fnv"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func isPayloadTypeUnrecognizedError(err error) bool {
	return errors.Is(err, protoregistry.NotFound) || strings.Contains(err.Error(), "isn't linked in")
}

// Fnv1aHashFn is a simple hash adapter function that uses the FNV-1a [hash.Hash] to generate hashes
var Fnv1aHashFn = func(bytes []byte) uint32 {
	h := fnv.New32a()
	// this never returns an error as per https://pkg.go.dev/hash#Hash
	_, _ = h.Write(bytes)
	return h.Sum32()
}

func defaultNativeOpts() []kgo.Opt {
	return []kgo.Opt{
		kgo.ProducerBatchCompression(kgo.ZstdCompression()),
		kgo.AllowAutoTopicCreation(),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()),
		kgo.RecordPartitioner(kgo.StickyKeyPartitioner(kgo.KafkaHasher(Fnv1aHashFn))),
	}
}
