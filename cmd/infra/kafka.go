package infra

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var producer *kafka.Writer

func ConnectKafka() *kafka.Writer {
	KAFKA_CA_CERT := os.Getenv("KAFKA_CA_CERT")
	KAFKA_TOPIC := os.Getenv("KAFKA_TOPIC")
	KAFKA_USERNAME := os.Getenv("KAFKA_USERNAME")
	KAFKA_PASSWORD := os.Getenv("KAFKA_PASSWORD")
	KAFKA_BROKER := os.Getenv("KAFKA_BROKER")

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM([]byte(KAFKA_CA_CERT))
	if !ok {
		log.Fatalf("Falha ao carregar certificado CA do Kafka")
	}

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	scramMechanism, err := scram.Mechanism(scram.SHA512, KAFKA_USERNAME, KAFKA_PASSWORD)
	if err != nil {
		log.Fatalf("Falha ao criar mecanismo scram: %s", err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		TLS:           tlsConfig,
		SASLMechanism: scramMechanism,
	}

	producer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{KAFKA_BROKER},
		Topic:    KAFKA_TOPIC,
		Balancer: &kafka.Hash{},
		Dialer:   dialer,
	})

	return producer
}

func SendMessage(data []byte) error {
	message := kafka.Message{Value: data}

	err := producer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("Falha ao enviar mensagem: %s", err)
		return err
	}

	log.Printf("Mensagem enviada para o tópico do Kafka")
	return nil
}

func SendMessageJSON(data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Falha ao serializar dados para JSON: %s", err)
		return err
	}

	return SendMessage(jsonData)
}
