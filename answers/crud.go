package answers

import (
	"createanswers/cachedb"
	"createanswers/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var answersKey = "ANSWER_KEY_"
var historyKey = "HISTORY_KEY_"

func CreateAnswer(w http.ResponseWriter, r *http.Request) {
	var newAnswer models.Answer
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Kindly enter data correctly")
		return
	}

	json.Unmarshal(reqBody, &newAnswer)
	if cmd := cachedb.RedisClient.Get(answersKey + newAnswer.Key); cmd.Err() != nil && cmd.Err() != redis.Nil {
		log.Printf("CreateAnswer: failed to get key from cache db. db error %v", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	} else if cmd.Err() == nil {
		log.Printf("CreateAnswer: key already exist in db key: %v", newAnswer.Key)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "key already exists")
		return
	}
	if cmd := cachedb.RedisClient.Set(answersKey+newAnswer.Key, newAnswer.Value, 0); cmd.Err() != nil {
		log.Printf("CreateAnswer: failed to set key in cache db %v", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}
	var event models.Event
	event.Type = "CREATE"
	event.Answer = newAnswer
	jsonData, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	data := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: jsonData,
	}
	if cmd := cachedb.RedisClient.ZAdd(historyKey+newAnswer.Key, data); cmd.Err() != nil {
		log.Printf("CreateAnswer: failed to set answer history value. %v", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAnswer)
}

func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	var newAnswer models.Answer
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Kindly enter data correctly")
	}
	json.Unmarshal(reqBody, &newAnswer)
	if cmd := cachedb.RedisClient.Get(answersKey + newAnswer.Key); cmd.Err() != nil && cmd.Err() != redis.Nil {
		log.Printf("UpdateAnswer: failed to get key from cache db. db error %v", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	} else if cmd.Err() != nil {
		log.Printf("UpdateAnswer: key does not exist db err:(%v), key: %v", cmd.Err(), newAnswer.Key)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "key does not exist")
		return
	}

	if cmd := cachedb.RedisClient.Set(answersKey+newAnswer.Key, newAnswer.Value, 0); cmd.Err() != nil {
		log.Printf("UpdateAnswer: failed to set key in cache db %v", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}
	var event models.Event
	event.Type = "UPDATE"
	event.Answer = newAnswer
	jsonData, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	data := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: jsonData,
	}
	if cmd := cachedb.RedisClient.ZAdd(historyKey+newAnswer.Key, data); cmd.Err() != nil {
		log.Printf("UpdateAnswer: failed to set answer history value. %v", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newAnswer)
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	cmd := cachedb.RedisClient.Get(answersKey + key)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		log.Printf("GetHistory: failed to get key from cache db %v", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	} else if cmd.Err() != nil {
		log.Printf("GetHistory: key entered doesn't exist in cache db %v", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "key doesn't exists")
		return
	}
	data := cachedb.RedisClient.ZRangeByScore(historyKey+key, redis.ZRangeBy{Max: "+inf", Min: "-inf"})
	if data.Err() != nil {
		log.Printf("GetHistory: failed to get the key history %v", data.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}
	var events []models.Event
	for _, v := range data.Val() {
		//str := v
		var event models.Event
		json.Unmarshal([]byte(v), &event)
		//json.Unmarshal(v.Member, &event)
		events = append(events, event)
	}

	/*var events []models.Event
	for _, v := range data.Val() {
		//str := v
		fmt.Println(v)
		var event models.Event
		//json.Unmarshal(v.Member, &event)
		events = append(events, v)
	}*/

	/*var event models.Event
	var answer models.Answer
	event.Type = "GET"
	answer.Key = key
	answer.Value = cmd.Val()
	event.Answer = answer
	events = append(events, event)*/
	json.NewEncoder(w).Encode(events)
}

func GetAnswerDetails(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	cmd := cachedb.RedisClient.Get(answersKey + key)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		log.Printf("GetAnswerDetails: failed to get key from cache db error (%v)", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	} else if cmd.Err() != nil {
		log.Printf("GetAnswerDetails: key entered doesn't exist in cache db %v", key)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "key doesn't exists")
		return
	}
	var answer models.Answer
	answer.Key = key
	answer.Value = cmd.Val()

	json.NewEncoder(w).Encode(answer)
}

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	cmd := cachedb.RedisClient.Get(answersKey + key)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		log.Printf("DeleteAnswer: failed to get key from cache db error (%v)", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	} else if cmd.Err() == nil {
		log.Printf("DeleteAnswer: key entered doesn't exist in cache db %v", key)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "key doesn't exists")
		return
	}
	value := cmd.Val()
	cmd2 := cachedb.RedisClient.Del(answersKey + key)
	if cmd2.Err() != nil {
		log.Printf("DeleteAnswer: failed to get key from cache db error (%v)", cmd2.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
		return
	}
	var event models.Event
	event.Type = "DELETE"
	event.Answer = models.Answer{Key: key, Value: value}
	jsonData, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	data := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: jsonData,
	}
	if cmd := cachedb.RedisClient.ZAdd(historyKey+key, data); cmd.Err() != nil {
		log.Printf("CreateAnswer: failed to set answer history value. %v", cmd.Err())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error")
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "The answer with key(%v) deleted successfully", key)
}
