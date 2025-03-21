package eventhandler

import (
	"context"
	"encoding/json"
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
	AddEventRecord(ctx context.Context, arg gen.AddEventRecordParams) error
	GetEventRecord(ctx context.Context, id int64) (int64, error)
}

type Handler struct {
	Producer Producer
	Store    Store
}

func (h *Handler) OnEvent(ctx context.Context, event *protos.Event) error {
	PrettyPrint(event)

	/**
	* 	TODO
	* 	- Implement the logic to handle and re-produce the event
	* 	- Return the event and an error if any (dlq store) - add option to create error
	* 	- Use the Store interface to interact with the database (use postgres?)
	 */

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

	PrettyPrint(&message)

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
	id, err := strconv.ParseInt(event.GetId(), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse event ID: %w", err)
	}

	dbRecord := gen.AddEventRecordParams{
		ID: int64(id),
		Tag: pgtype.Text{
			String: event.GetTag().String(),
		},
		Stat: pgtype.Text{
			String: event.GetStatus().String(),
		},
		Msg: pgtype.Text{
			String: event.GetMessage(),
		},
	}

	err = h.Store.AddEventRecord(ctx, dbRecord)
	if err != nil {
		return err
	}

	id, err = h.Store.GetEventRecord(ctx, int64(id))
	if err != nil {
		return err
	}

	fmt.Println("Event %s saved successfully!", id)

	return nil
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
