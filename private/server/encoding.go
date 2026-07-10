package server

import (
	"net/http"
	"strings"
)

// clientAcceptsGzip reports whether the client has indicated it can decode
// gzip-compressed response frames. Per the gRPC spec:
//
//   - grpc-accept-encoding lists the compression algorithms the client can
//     decode in responses (comma-separated, case-insensitive). If it includes
//     "gzip" we should compress responses.
//
//   - If grpc-accept-encoding is absent but the client sent a gzip-encoded
//     request (grpc-encoding: gzip), it implicitly accepts gzip responses too,
//     so we mirror the encoding.
//
// Reference: https://grpc.github.io/grpc/core/md_doc_compression.html
func clientAcceptsGzip(r *http.Request) bool {
	if accept := r.Header.Get("grpc-accept-encoding"); accept != "" {
		// grpc-accept-encoding is the authoritative list of what the client can
		// decode. Only use the grpc-encoding fallback when it is absent entirely.
		for enc := range strings.SplitSeq(accept, ",") {
			if strings.EqualFold(strings.TrimSpace(enc), "gzip") {
				return true
			}
		}
		return false
	}
	// Fallback: if the client sent a gzip request it must be able to decode gzip.
	return strings.EqualFold(strings.TrimSpace(r.Header.Get("grpc-encoding")), "gzip")
}
