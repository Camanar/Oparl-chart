package charts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Camanar/Oparl-neo4jConnector/reader"
)

// Define the Amendment struct to hold each amendment
type Data struct {
	X string `json:"x"`
	Y int    `json:"y"`
}

// Define the Response struct to hold the final response
type ResponseAm struct {
	Data  []Data `json:"data"`
	ID    string `json:"id"`
	Color string `json:"color"`
}

// Define the Filter struct to parse the JSON parameter
type Filter struct {
	ID     int      `json:"id"`
	Metric []Metric `json:"metric"`
	Symbol []Symbol `json:"symbol"`
	Value  []Value  `json:"value"`
}

type Metric struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Symbol struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Value struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func GetAmendments(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// Check for textId parameter
	textId := queryParams.Get("textId")
	amendmentId := queryParams.Get("amendmentId")
	fileId := queryParams.Get("uid")

	// Base query
	baseQuery := `
		MATCH (f:File)<-[:RELATED]-(stf:StateFile)<-[:RELATED]-(sf:StageFile)<-[:TEXT_VERSION]-(t:Text)<-[:AMENDS]-(a:Amendment)
		%s
		WITH DISTINCT a, t
		OPTIONAL MATCH (a)-[arta:AMENDS]->(art:Article)
		OPTIONAL MATCH (a)-[pa:AMENDS]->(p:Paragraph)
		MATCH (a)-[:WRITTEN_BY]->(d:Deputy)
		MATCH (d)-[:IS_IN_GROUP]->(g:Group)
		WITH a.lifecycle_date_publication AS x, COUNT(a) AS y
		ORDER BY x
		RETURN collect({x: x, y: y}) AS data
	`

	// Initialize a map to hold conditions grouped by field
	conditionsByField := make(map[string][]string)

	// Construct the WHERE clause
	var whereClauses []string
	for _, conditions := range conditionsByField {
		whereClauses = append(whereClauses, fmt.Sprintf("(%s)", strings.Join(conditions, " OR ")))
	}

	// Add textId condition if it exists
	if textId != "" {
		whereClauses = append(whereClauses, fmt.Sprintf(`t.text_id = "%s"`, textId))
	}

	// Add fileId condition if it exists
	if fileId != "" {
		whereClauses = append(whereClauses, fmt.Sprintf(`f.uid = "%s"`, fileId))
	}

	// Add amendmentId condition if it exists
	if amendmentId != "" {
		whereClauses = append(whereClauses, fmt.Sprintf(`a.amendment_id = "%s"`, amendmentId))
	}

	// Combine all conditions with AND
	var whereClause string
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Format the final query
	finalQuery := fmt.Sprintf(baseQuery, whereClause)

	log.Printf("Final query: %s", finalQuery)

	// Execute the query
	result, err := reader.ReadNeo4j(finalQuery, map[string]interface{}{}, "t")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//log.Printf("Query result: %v", result)

	// Initialize a slice to hold the amendments
	var data []Data

	if len(result) > 0 {
		// We only need to process the first record since the query uses collect()
		record := result[0]
		response := record.AsMap()

		jsonData := response["data"]

		// Convert jsonData into JSON bytes
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshal directly into amendments slice
		if err := json.Unmarshal(jsonBytes, &data); err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Create the final response structure
	finalResponse := ResponseAm{
		Data:  data,
		ID:    "Amendment",
		Color: "hsl(171, 70%, 50%)",
	}

	// Marshal the final response to JSON
	jsonResponse, err := json.Marshal(finalResponse)
	if err != nil {
		log.Printf("Error marshalling final response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Println(finalQuery)
}
