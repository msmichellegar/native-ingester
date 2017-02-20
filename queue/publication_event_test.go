package queue

import (
	"testing"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/stretchr/testify/assert"
)

const expectedTID = "tid_test"
const expectedOriginSystemID = "http://cmdb.ft.com/systems/methode-web-pub"
const expectedTimestamp = "2017-02-16T12:56:16Z"

var someMsgHeaders = map[string]string{
	"X-Request-Id":      expectedTID,
	"Origin-System-Id":  expectedOriginSystemID,
	"Message-Timestamp": expectedTimestamp,
}

var aMsg = consumer.Message{
	Headers: someMsgHeaders,
	Body:    `{"foo":"bar"}`,
}

var aMsgWithBadBody = consumer.Message{
	Headers: someMsgHeaders,
	Body:    `I'm not JSON`,
}

var aMsgWithoutTimestamp = consumer.Message{
	Headers: map[string]string{},
	Body:    `{"foo":"bar"}`,
}

func TestGetTransactionID(t *testing.T) {
	pe := publicationEvent{aMsg}
	actualTID := pe.transactionID()
	assert.Equal(t, expectedTID, actualTID, "The transaction ID shoud be the same of a consumer message")
}

func TestGetOriginSystemID(t *testing.T) {
	pe := publicationEvent{aMsg}
	actualOriginSystemID := pe.originSystemID()
	assert.Equal(t, expectedOriginSystemID, actualOriginSystemID, "The Origin-System-Id shoud be the same of a consumer message")
}

func TestGetContentBodySuccessfully(t *testing.T) {
	pe := publicationEvent{aMsg}
	body, err := pe.contentBody()

	assert.Nil(t, err, "It should not return an error")
	assert.Equal(t, "bar", body["foo"], "The body should contain the original data")
	assert.Equal(t, expectedTimestamp, body["lastModified"], "The body should contain a lastModified attribute equal to the timestamp message header")
	assert.Equal(t, expectedTID, body["publishReference"], "The body should contain the publishReference attribute equal to the message transaction ID")
}

func TestGetContentBodyFailBecauseBadBody(t *testing.T) {
	pe := publicationEvent{aMsgWithBadBody}
	_, err := pe.contentBody()

	assert.EqualError(t, err, "invalid character 'I' looking for beginning of value", "It should return an error")
}

func TestGetContentBodyFailBecauseMissingTimstamp(t *testing.T) {
	pe := publicationEvent{aMsgWithoutTimestamp}
	_, err := pe.contentBody()

	assert.EqualError(t, err, "Publish event does not contain timestamp", "It should return an error")
}

func TestGetProducerMessage(t *testing.T) {
	pe := publicationEvent{aMsg}
	actualProducerMsg := pe.producerMsg()

	assert.Equal(t, aMsg.Body, actualProducerMsg.Body, "It should have the same body of the consumer message")
	assert.Equal(t, aMsg.Headers, actualProducerMsg.Headers, "It should have the same headers of the consumer message")
}