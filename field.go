package param

import (
	"time"
)

// Helper functions for common types.
func String(s string) Opt[string]     { return From(s) }
func Int(i int) Opt[int]              { return From(i) }
func Int64(i int64) Opt[int64]        { return From(i) }
func Bool(b bool) Opt[bool]           { return From(b) }
func Float(f float64) Opt[float64]    { return From(f) }
func Float32(f float32) Opt[float32]  { return From(f) }
func Time(t time.Time) Opt[time.Time] { return From(t) }

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T { return &v }
