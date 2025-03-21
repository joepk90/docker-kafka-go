package eventhandler

import (
	"context"
	"testing"

	"github.com/jparkkennaby/docker-kafka-go/eventhandler/mocks"
	protos "github.com/jparkkennaby/docker-kafka-go/gen/go/protos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_OnEvent(t *testing.T) {
	eventID := "123"
	t.Run("success", func(t *testing.T) {
		// arrange
		databaseRecordID := int64(1)
		expectedEvent := &protos.Event{
			Id:      eventID,
			Message: "Hello World!",
			Tag:     protos.Tag_EVENT_TAG_TECHNOLOGY,
			Status:  protos.Status_EVENT_STATUS_SUCCESSFUL,
		}

		Producer := &mocks.Producer{}
		Producer.On(
			"ProduceMessage",
			mock.Anything,
			expectedEvent,
			AppName,
			"protos.Event",
		).Once().Return(nil)

		mockDB := &mocks.Store{}
		mockDB.On("AddEventRecord", mock.Anything, mock.Anything).Once().Return(databaseRecordID, nil)
		mockDB.On("GetEventRecord", mock.Anything, mock.Anything).Once().Return(databaseRecordID, nil)

		handler := &Handler{
			Producer: Producer,
			Store:    mockDB,
		}

		// act
		err := handler.OnEvent(context.Background(), expectedEvent)

		// assert
		assert.NoError(t, err)
		Producer.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})
}
