# Homer ingest: passing `[]byte` into DuckDB Appender

Homer builds `data_extra` JSON in pooled `[]byte` buffers. To avoid an extra
`string()` allocation at DuckLake flush time, use this fork's vector APIs and
pass bytes through `database/sql/driver.Value`.

## Bindings (this repo)

| API | Role |
|-----|------|
| `VectorAssignStringElementLen` | `[]byte` → VARCHAR with UTF-8 validation (no `C.CBytes`) |
| `UnsafeVectorAssignStringElementLen` | Same without UTF-8 check (caller guarantees UTF-8) |
| `VectorAssignUTF8Bytes` | Chooses safe vs unsafe assign |
| `UTF8Bytes` / `UnsafeUTF8Bytes` | Named slice types for documentation |

Tagged releases: `v0.10502.0-homer.gcopt.N` on [adubovikov/duckdb-go-bindings](https://github.com/adubovikov/duckdb-go-bindings).

## `go.mod` replace (Homer)

```go
replace github.com/duckdb/duckdb-go-bindings => github.com/adubovikov/duckdb-go-bindings v0.10502.0-homer.gcopt.3
```

## Homer `cellToDriverValue`

For JSON columns, `duckdb-go` `setJSON` accepts `json.RawMessage` (alias of
`[]byte`) without re-marshalling:

```go
import "encoding/json"

func cellToDriverValue(v interface{}) interface{} {
	if bp, ok := v.(*[]byte); ok {
		if bp == nil || len(*bp) == 0 {
			return json.RawMessage("{}")
		}
		return json.RawMessage(*bp)
	}
	return v
}
```

Plain `[]byte` works for `VARCHAR` columns via `setBytes`.

## Lifetime

`UnsafeVectorAssignStringElementLen` borrows the Go slice pointer until DuckDB
copies into vector storage. Keep pooled `[]byte` alive until the Appender batch
for that row is flushed (Homer: until `putRowSlice` after `AppendRow` loop).
