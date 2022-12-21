package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const RFC3339 = "2006-01-02T15:04:05Z07:00"

// service funcs

func getFileDescriptor() (*os.File, error) {
	var filename string
	flag.StringVar(&filename, "file", "", "path to operation file (json)")
	flag.Parse()
	if filename != "" {
		return os.Open(filename)
	}
	filename = os.Getenv("FILE")
	if filename != "" {
		return os.Open(filename)
	}
	return os.Stdin, nil
}

// OperationID struct

type OperationID struct {
	StrID string
	IntID int
}

func (i OperationID) String() string {
	if i.StrID != "" {
		return i.StrID
	} else {
		return strconv.Itoa(i.IntID)
	}
}

func (i *OperationID) UnmarshalJSON(data []byte) error {
	if strings.HasPrefix(string(data), "\"") {
		i.StrID = string(data[1 : len(data)-1])
	} else {
		if err := json.Unmarshal(data, &i.IntID); err != nil {
			return fmt.Errorf("operation id conversion error (data = %v): %w", data, err)
		}
	}
	return nil
}

func (i OperationID) MarshalJSON() ([]byte, error) {
	if i.StrID != "" {
		return json.Marshal(i.StrID)
	}
	return json.Marshal(i.IntID)
}

// struct to load JSON data

type OperationBody struct {
	Type      string      `json:"type"`
	Value     json.Number `json:"value"`
	ID        OperationID `json:"id"`
	CreatedAt string      `json:"created_at"`
}

type Operation struct {
	Company   string        `json:"company"`
	Operation OperationBody `json:"operation"`
	OperationBody
}

// struct to convert

type ConvertedOperation struct {
	Company   string
	Incorrect bool
	Value     int
	ID        OperationID
	CreatedAt time.Time
}

func (op Operation) convert() ConvertedOperation {
	// move fields from operation to root
	if op.Type == "" {
		op.Type = op.Operation.Type
	}
	if op.Value == "" {
		op.Value = op.Operation.Value
	}
	if op.ID.String() == "" {
		op.ID = op.Operation.ID
	}
	if op.CreatedAt == "" {
		op.CreatedAt = op.Operation.CreatedAt
	}

	// check & convert
	incorrect := (op.Company == "") || op.Type == "" || op.Value == "" ||
		op.ID.String() == "" || op.CreatedAt == "" ||
		(op.Type != "income" && op.Type != "outcome" && op.Type != "+" && op.Type != "-")
	
	var value int64
	var err error
	var createdAt time.Time
	
	if !incorrect {
		value, err = op.Value.Int64()
		incorrect = incorrect || (err != nil)
		if op.Type == "outcome" || op.Type == "-" {
			value = -value
		}
	}
	
	createdAt, err = time.Parse(RFC3339, op.CreatedAt)
	incorrect = incorrect || (err != nil)

	return ConvertedOperation{
		Company:   op.Company,
		Incorrect: incorrect,
		Value:     int(value),
		ID:        op.ID,
		CreatedAt: createdAt,
	}
}

// struct to result

type OperationIdTime struct {
	OperationID
	CreatedAt time.Time
}

type OperationResult struct {
	Company       string            `json:"company"`
	ValidOpsCount int               `json:"valid_operations_count"`
	Balance       int               `json:"balance"`
	InvalidOps    []OperationIdTime `json:"invalid_operations,omitempty"`
}

func processOperations(billingData []Operation) []OperationResult {
	resultMap := make(map[string]OperationResult)

	for _, op := range billingData {
		converted := op.convert()

		body := resultMap[converted.Company]
		body.Company = converted.Company
		
		if converted.Incorrect {
			body.InvalidOps = append(body.InvalidOps, OperationIdTime{
				OperationID: converted.ID,
				CreatedAt:   converted.CreatedAt,
			})
		} else {
			body.Balance += converted.Value
			body.ValidOpsCount++
		}

		resultMap[converted.Company] = body
	}

	resultMapKeys := make([]string, 0, len(resultMap))
	for k := range resultMap {
		resultMapKeys = append(resultMapKeys, k)
	}
	sort.Strings(resultMapKeys)

	result := make([]OperationResult, len(resultMap))
	for idx, k := range resultMapKeys {
		body := resultMap[k]
		sort.SliceStable(body.InvalidOps, func(i int, j int) bool {
			return body.InvalidOps[i].CreatedAt.Before(body.InvalidOps[j].CreatedAt)
		})
		result[idx] = body
	}

	return result
}

func main() {
	file, err := getFileDescriptor()
	if err != nil {
		log.Fatalf("error open file: %s", err.Error())
	}
	defer file.Close()

	var billingData []Operation
	err = json.NewDecoder(file).Decode(&billingData)
	if err != nil {
		log.Fatalf("error load data: %s", err.Error())
	}

	report := processOperations(billingData)

	data, err := json.MarshalIndent(report, "", "\t")
	if err != nil {
		log.Fatalf("error encode data: %s", err.Error())
	}

	err = os.WriteFile("out.json", data, 0644)
	if err != nil {
		log.Fatalf("error create file: %s", err.Error())
	}
}
