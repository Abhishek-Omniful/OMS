package models

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockCursor implements MongoCursor-like behavior
type MockCursor struct {
	mock.Mock
}

// MockCollection implements basic InsertOne and Find
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(_ interface{}, doc interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(doc)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockCollection) Find(_ interface{}, filter interface{}) (MongoCursor, error) {
	args := m.Called(filter)
	return args.Get(0).(MongoCursor), args.Error(1)
}

func TestCreateWebhook_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		input       *Webhook
		mockSetup   func(*MockCollection)
		expectedErr bool
	}{
		{
			name: "success",
			input: &Webhook{
				URL: "https://test.com", TenantID: 1,
			},
			mockSetup: func(m *MockCollection) {
				m.On("InsertOne", mock.Anything).Return(&mongo.InsertOneResult{}, nil)
			},
			expectedErr: false,
		},
		{
			name: "invalid input",
			input: &Webhook{
				URL: "", TenantID: 0,
			},
			mockSetup:   func(_ *MockCollection) {},
			expectedErr: true,
		},
		{
			name: "insert error",
			input: &Webhook{
				URL: "https://test.com", TenantID: 1,
			},
			mockSetup: func(m *MockCollection) {
				m.On("InsertOne", mock.Anything).Return(&mongo.InsertOneResult{}, errors.New("db error"))
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockColl := new(MockCollection)
			tc.mockSetup(mockColl)
			ctx = context.Background()
			webhookCollection = mockColl

			err := CreateWebhook(tc.input)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockColl.AssertExpectations(t)
		})
	}
}

func TestListWebhooks_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*MockCollection, *MockCursor)
		expectedErr bool
	}{
		{
			name: "success",
			mockSetup: func(mc *MockCollection, cursor *MockCursor) {
				cursor.On("All").Return(nil)
				cursor.On("Close").Return(nil)
				mc.On("Find", mock.Anything).Return(cursor, nil)
			},
			expectedErr: false,
		},
		{
			name: "find error",
			mockSetup: func(mc *MockCollection, _ *MockCursor) {
				mc.On("Find", mock.Anything).Return((*MockCursor)(nil), errors.New("find error"))
			},
			expectedErr: true,
		},

		{
			name: "decode error",
			mockSetup: func(mc *MockCollection, cursor *MockCursor) {
				cursor.On("All").Return(errors.New("decode error"))
				cursor.On("Close").Return(nil)
				mc.On("Find", mock.Anything).Return(cursor, nil)
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockColl := new(MockCollection)
			mockCursor := new(MockCursor)
			ctx = context.Background()
			tc.mockSetup(mockColl, mockCursor)
			webhookCollection = mockColl

			result, err := ListWebhooks()

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
			}
			mockColl.AssertExpectations(t)
			mockCursor.AssertExpectations(t)
		})
	}
}
