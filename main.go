package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/websocket"
)

var (
	es       *elasticsearch.Client
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

type ElasticsearchResponse struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Value string `json:"value"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type AggregationResult struct {
	Aggregations struct {
		MostFrequentValues struct {
			Buckets []struct {
				Key   string `json:"key"`
				Count int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"most_frequent_values"`
	} `json:"aggregations"`
}

type RequestBody struct {
	Word string `json:"word"`
}

func init() {
	var err error
	es, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Fatalf("error creating the Elasticsearch client: %s", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/populate", populateIndex)
	http.Handle("/send", withCORS(http.HandlerFunc(sendWord)))
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func populateIndex(w http.ResponseWriter, r *http.Request) {
	data := loadMockData()
	for _, value := range data {
		_, err := es.Index(
			"autocomplete",
			strings.NewReader(fmt.Sprintf(`{"value": "%s"}`, value)),
			es.Index.WithDocumentType("_doc"),
		)
		if err != nil {
			http.Error(w, "failed to index data", http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("data populated"))
}

func sendWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody RequestBody
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	_, err = es.Index(
		"autocomplete",
		strings.NewReader(fmt.Sprintf(`{"value": "%s"}`, requestBody.Word)),
		es.Index.WithDocumentType("_doc"),
	)

	if err != nil {
		http.Error(w, "failed to index data", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("data populated"))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("error while reading message:", err)
			break
		}

		query := string(msg)
		suggestions, err := getSuggestions(query)
		if err != nil {
			log.Println("error while getting suggestions:", err)
			continue
		}

		if query == "Word send" {
			time.Sleep(500 * time.Millisecond) // lol
		}

		getMostFrequent, err := getMostFrequentValues()
		if err != nil {
			log.Println("error while getting most frequent values:", err)
			continue
		}

		reponseWs := struct {
			FrequentValues []struct {
				Key   string `json:"key"`
				Count int    `json:"doc_count"`
			} `json:"frequent_values"`
			Suggestions []string `json:"suggestions"`
		}{
			getMostFrequent,
			suggestions,
		}

		err = conn.WriteJSON(reponseWs)
		if err != nil {
			log.Println("error while writing message:", err)
			break
		}
	}
}

func getSuggestions(query string) ([]string, error) {
	var suggestions []string

	resp, err := es.Search(
		es.Search.WithIndex("autocomplete"),
		es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
			"query": {
				"wildcard": {
					"value": "*%s*"
				}
			}
		}`, query))),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return suggestions, err
	}
	defer resp.Body.Close()

	var result ElasticsearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return suggestions, err
	}

	for _, hit := range result.Hits.Hits {
		suggestions = append(suggestions, hit.Source.Value)
	}

	return suggestions, nil
}

func getMostFrequentValues() ([]struct {
	Key   string `json:"key"`
	Count int    `json:"doc_count"`
}, error) {
	query := `{
		"size": 0,
		"aggs": {
			"most_frequent_values": {
				"terms": {
					"field": "value.keyword",
					"size": 10
				}
			}
		}
	}`

	resp, err := es.Search(
		es.Search.WithIndex("autocomplete"),
		es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result AggregationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Aggregations.MostFrequentValues.Buckets, nil

}

func loadMockData() []string {
	return []string{
		"Apple", "Banana", "Cherry", "Date", "Fig",
		"Grape", "Honeydew", "Kiwi", "Lemon", "Mango",
	}
}
