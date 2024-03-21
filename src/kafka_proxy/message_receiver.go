package kafka_proxy

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type MessageHandler struct {
	MessageQueue chan string
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (mh *MessageHandler) Setup(session sarama.ConsumerGroupSession) error {
	logrus.Infof("setup success,session context is %v ,claim is %v", session.Context(), session.Claims())
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (mh *MessageHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	logrus.Infof("cleanup success,session context is %v ,claim is %v", session.Context(), session.Claims())
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (mh *MessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		mh.MessageQueue <- string(message.Value)
		//logrus.Debugf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}
