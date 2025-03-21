package eventhandler

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	protos "github.com/jparkkennaby/docker-kafka-go/gen/go/protos"

	"github.com/jparkkennaby/docker-kafka-go/store/gen"
	"google.golang.org/protobuf/proto"
)

const (
	AppName = "docker-kafka-go"
)

//go:generate mockery --name Producer
type Producer interface {
	ProduceMessage(ctx context.Context, message proto.Message, app, comment string) error
}

//go:generate mockery --name Store
type Store interface {
	AddEventRecord(ctx context.Context, request gen.AddEventRecordParams) (int64, error)
	GetEventRecord(ctx context.Context, id int64) (int64, error)
}

type Handler struct {
	Producer Producer
	Store    Store
}

func (h *Handler) OnEvent(ctx context.Context, event *protos.Event) error {
	// enabling this line will cause events to be sent to a DLQ topic
	// return fmt.Errorf("return an error so event are sent to the dlq topic...")

	// save the event to the db
	err := h.AddEventToStore(ctx, event)
	if err != nil {
		return err
	}

	// update the event
	message := &protos.Event{
		Id:      "123",
		Message: "Hello World!",
		Tag:     event.GetTag(),
		Status:  event.GetStatus(),
	}

	// publish the event
	err = h.Producer.ProduceMessage(ctx, message, AppName, "protos.Event")
	if err != nil {
		return err
	}

	slog.Info("message handled, updated and published", "id", message.GetId())
	fmt.Println("message handled, updated and published")

	return nil
}

func (h *Handler) AddEventToStore(ctx context.Context, event *protos.Event) error {
	eventID, err := strconv.ParseInt(event.GetId(), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse event ID: %w", err)
	}

	dbRecord := gen.AddEventRecordParams{
		EventID: pgtype.Int8{
			Int64: eventID,
		},
		EventTag: pgtype.Text{
			String: event.GetTag().String(),
		},
		EventStat: pgtype.Text{
			String: event.GetStatus().String(),
		},
		EventMsg: pgtype.Text{
			String: event.GetMessage(),
		},
	}

	id, err := h.Store.AddEventRecord(ctx, dbRecord)
	if err != nil {
		return err
	}

	id, err = h.Store.GetEventRecord(ctx, int64(id))
	if err != nil {
		return err
	}

	fmt.Println("Event saved successfully! Event ID: ", id)

	return nil
}
