// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package gen

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type KafkaEvent struct {
	ID          int64
	EventID     pgtype.Int8
	EventTag    pgtype.Text
	EventStat   pgtype.Text
	EventMsg    pgtype.Text
	RequestDate pgtype.Timestamptz
}
