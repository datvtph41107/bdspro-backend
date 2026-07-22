// pkg/optional/apply.go
package optional

// ApplyString áp dụng String vào fields map
func ApplyString(fields map[string]interface{}, col string, opt String) {
	switch opt.State() {
	case StatePresent:
		fields[col] = opt.Value()
	case StateNull:
		fields[col] = nil
		// StateAbsent: không làm gì
	}
}

// ApplyUint64 áp dụng Uint64 vào fields map
func ApplyUint64(fields map[string]interface{}, col string, opt Uint64) {
	switch opt.State() {
	case StatePresent:
		fields[col] = opt.Value()
	case StateNull:
		fields[col] = nil
	}
}

// ApplyUint32 áp dụng Uint32 vào fields map
func ApplyUint32(fields map[string]interface{}, col string, opt Uint32) {
	switch opt.State() {
	case StatePresent:
		fields[col] = opt.Value()
	case StateNull:
		fields[col] = nil
	}
}

// ApplyInt32 áp dụng Int32 vào fields map
func ApplyInt32(fields map[string]interface{}, col string, opt Int32) {
	switch opt.State() {
	case StatePresent:
		fields[col] = opt.Value()
	case StateNull:
		fields[col] = nil
	}
}

// ApplyFloat64 áp dụng Float64 vào fields map
func ApplyFloat64(fields map[string]interface{}, col string, opt Float64) {
	switch opt.State() {
	case StatePresent:
		fields[col] = opt.Value()
	case StateNull:
		fields[col] = nil
	}
}

// ApplyBool áp dụng Bool vào fields map
func ApplyBool(fields map[string]interface{}, col string, opt Bool) {
	switch opt.State() {
	case StatePresent:
		fields[col] = opt.Value()
	case StateNull:
		fields[col] = nil
	}
}

// ApplySlice áp dụng Slice vào fields map
func ApplySlice[T any](fields map[string]interface{}, col string, opt Slice[T]) {
	switch opt.State() {
	case StatePresent:
		fields[col] = opt.Value()
	case StateNull:
		fields[col] = nil
	}
}
