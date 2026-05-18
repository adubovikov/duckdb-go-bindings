package duckdb_go_bindings

// UTF8Bytes is a UTF-8 byte slice for VARCHAR / JSON appender columns.
// Use with duckdb-go AppendRow by passing UTF8Bytes(...) as driver.Value, or call
// VectorAssignUTF8Bytes directly on a writable vector.
type UTF8Bytes []byte

// UnsafeUTF8Bytes is UTF-8 payload for duckdb_unsafe_vector_assign_string_element_len.
// The backing slice must stay valid until DuckDB has copied the value into vector
// storage (for Homer: until the Appender row batch flush for that chunk completes).
type UnsafeUTF8Bytes []byte

// Bytes returns the underlying slice without allocation.
func (b UTF8Bytes) Bytes() []byte { return []byte(b) }

// Bytes returns the underlying slice without allocation.
func (b UnsafeUTF8Bytes) Bytes() []byte { return []byte(b) }

// VectorAssignUTF8Bytes writes blob into a VARCHAR/BLOB vector element.
// When validateUTF8 is true, DuckDB validates UTF-8 (VARCHAR); when false, the
// caller must ensure valid UTF-8 (typical for pre-serialized JSON from a pool).
func VectorAssignUTF8Bytes(vec Vector, index IdxT, blob []byte, validateUTF8 bool) {
	if validateUTF8 {
		VectorAssignStringElementLen(vec, index, blob)
		return
	}
	UnsafeVectorAssignStringElementLen(vec, index, blob)
}
