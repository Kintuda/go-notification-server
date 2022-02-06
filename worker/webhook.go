package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	db "github.com/Kintuda/notification-server/database"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v4"
)

type Event struct {
	Id       string                 `json:"id"`
	Url      string                 `json:"url"`
	Payload  map[string]interface{} `json:"payload"`
	Attempts int                    `json:"max_attempts"`
}

type Response struct {
	Status   int    `json:"status"`
	Response string `json:"body"`
}

type Row struct {
	Count int
}

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func FindAttempts(tx pgx.Tx, id string) (int, error) {
	count := 0
	rows, err := tx.Query(context.TODO(), "SELECT * FROM notification_attempts where notification_id = $1", id)

	if err != nil {
		return count, err
	}

	for rows.Next() {
		count += 1
	}

	defer rows.Close()

	return count, nil
}

func CreateWebhookTransaction(evt Event, database *db.DatabaseConnection) bool {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatal(err)
		return false
	}

	tx, err := database.Conn.BeginTx(context.TODO(), pgx.TxOptions{})

	if err != nil {
		logger.Error("error creating a transaction", zap.Error(err))
		return false
	}

	attempts, err := FindAttempts(tx, evt.Id)

	if err != nil {
		logger.Error("error finding the past attempts", zap.Error(err))
		return false
	}

	if evt.Attempts <= attempts {
		logger.Error("max attempt reached", zap.Int("max_attempts", evt.Attempts), zap.Int("attempts_count", attempts))
		return true
	}

	state := "failed"
	result, err := SendWebhook(evt)

	if err != nil {
		logger.Error("error sending postback", zap.Error(err))
		_, err = tx.Exec(context.TODO(), "INSERT INTO notification_attempts (state, notification_id, response_status, response_body) values ($1, $2, $3, $4)", state, evt.Id, result.Status, result.Response)

		if err != nil {
			logger.Error("error while inserting", zap.Error(err))
			return false
		}
		tx.Commit(context.TODO())
		return false
	}

	if result.Status > 200 && result.Status < 300 {
		state = "delivered"
	}

	_, err = tx.Exec(context.TODO(), "INSERT INTO notification_attempts (state, notification_id, response_status, response_body) values ($1, $2, $3, $4)", state, evt.Id, result.Status, result.Response)

	if err != nil {
		logger.Error("error while inserting", zap.Error(err))
		return false
	}

	tx.Commit(context.TODO())

	if state != "delivered" {
		return false
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	return true
}

func SendWebhook(evt Event) (*Response, error) {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	client := &http.Client{}
	json, err := json.Marshal(evt.Payload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, evt.Url, bytes.NewBuffer(json))

	if err != nil {
		return nil, err
	}

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	bodyString := string(body)

	logger.Info("request finished", zap.String("body", bodyString))

	if err != nil {
		return nil, err
	}

	return &Response{Response: bodyString, Status: response.StatusCode}, nil
}
