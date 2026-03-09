package errors

import "fmt"

// NegativeOffsetError is returned when offset is negative.
type NegativeOffsetError struct {
	Offset int
}

func (e *NegativeOffsetError) Error() string {
	return fmt.Sprintf("Offset `%d` cannot be negative.", e.Offset)
}

// PositionOutOfBoundsError is returned when position is out of bounds.
type PositionOutOfBoundsError struct {
	Length   int
	Position int
}

func (e *PositionOutOfBoundsError) Error() string {
	return fmt.Sprintf("Position `%d` is out of bounds (`0 < position < %d`).", e.Position, e.Length)
}

// RecursiveReadLimitExceededError is returned when recursive read limit is exceeded.
type RecursiveReadLimitExceededError struct {
	Count int
	Limit int
}

func (e *RecursiveReadLimitExceededError) Error() string {
	return fmt.Sprintf("Recursive read limit of `%d` exceeded (recursive read count: `%d`).", e.Limit, e.Count)
}
