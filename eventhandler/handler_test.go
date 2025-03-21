package eventhandler

import (
	"context"
	"database/sql"
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
		expectedEvent := &protos.Event{
			Id:     eventID,
			Tag:    protos.Tag_EVENT_TAG_TECHNOLOGY,
			Status: protos.Status_EVENT_STATUS_SUCCESSFUL,
		}

		Producer := &mocks.Producer{}
		Producer.On(
			"ProduceMessage",
			mock.Anything,
			expectedEvent,
			AppName,
			eventID,
		).Once().Return(nil)

		mockDB := &mocks.Store{}
		mockDB.On("AddEventRecord", mock.Anything, mock.Anything).Once().Return("", sql.ErrNoRows)
		mockDB.On("GetEventRecord", mock.Anything, mock.Anything).Once().Return(nil)

		handler := &Handler{
			Producer: Producer,
			Store:    mockDB,
		}

		// act
		err := handler.OnEvent(context.Background(), &protos.Event{
			Id:     eventID,
			Tag:    protos.Tag_EVENT_TAG_TECHNOLOGY,
			Status: protos.Status_EVENT_STATUS_SUCCESSFUL,
		})

		// assert
		assert.NoError(t, err)
		Producer.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})

}
