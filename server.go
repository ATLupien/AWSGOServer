package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
)
type Response struct{
	Table string
	Count *int64
}
type status struct {
	HTTP int
	Time time.Time
}
type Item struct {
	FDate           interface{}
	FStatus         interface{}
	Fnumber         interface{}
	FIata           interface{}
	FIcao           interface{}
	ALIata          string
	ALIcao          interface{}
	Altitude        interface{}
	Direction       interface{}
	Latitude        interface{}
	Longitude       interface{}
	SpeedHorizontal interface{}
	SpeedVertical   interface{}
	Updated         interface{}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/alupien/status", Status).
		Methods("GET").
		Schemes("http", "https")
	router.HandleFunc("/alupien/all", Data).
			Methods("GET").
			Schemes("http", "https")
	router.HandleFunc("/alupien/search", Query).
				Methods("GET").
				Schemes("http", "https")
	router.HandleFunc("/*", HTTP405).
		Methods("HEAD", "POST", "PUT", "PATCH",
			"DELETE", "CONNECT", "OPTIONS", "TRACE").
		Schemes("http", "https")
  http.Handle("/alupien/status", router)
	http.Handle("/alupien/all", router)
	http.Handle("/alupien/search", router)
	http.ListenAndServe(":8080", nil)}

func HTTP405(w http.ResponseWriter, r *http.Request) {
	var tag string
	tag = "HTTP405"
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Println(tag, ip)
	Status := status{
		HTTP: http.StatusMethodNotAllowed,
		Time: time.Now(),	}
	msg, _ := json.Marshal(Status)
	fmt.Println(msg)
	w.Write(msg)}

func Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println("IN STATUS")
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	var tag string
	tag = "Status"
	fmt.Println(tag, ip)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),	})
	svc := dynamodb.New(sess)
	var tName = "alupien-Airplane"
	des := &dynamodb.DescribeTableInput{
		TableName: &tName,	}
	resp, err := svc.DescribeTable(des)
	if err != nil {
		fmt.Println("err:", err)	}
	var sResponse Response
	sResponse.Table = "alupien-Airplane"
	sResponse.Count = resp.Table.ItemCount
	json.NewEncoder(w).Encode(sResponse)}

func Data(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	var tag string
	tag = "Data"
	fmt.Println(tag, ip)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),	})
	svc := dynamodb.New(sess)
	var tableName = "alupien-Airplane"
	describeTables := &dynamodb.ScanInput{
		TableName: &tableName,	}
	resp, _ := svc.Scan(describeTables)
	messg := []Item{}
	dynamodbattribute.UnmarshalListOfMaps(resp.Items, &messg)
	//msg, _ := json.Marshal(messg)
fmt.Printf("err=%v", err)
//	w.Write(msg)
json.NewEncoder(w).Encode(messg)
}

func Query(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	var tag string
	tag = "Query"
	fmt.Println(tag, ip)
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),	})
	svc := dynamodb.New(sess)
	var tName = "alupien-Airplane"
	query := r.URL.Query()
	Fnumber, present := query["Fnumber"]
	g := 1
	if !present || len(Fnumber) == 0 {
		fmt.Println("No Fnumber's Present")
		g = 0			   }
	if len(Fnumber) > 1 {
		fmt.Println("Too Many Query Params")
		g = 0	}
	if g == 1 {
		fmt.Println("HERE WE GO!!!")
	param := Fnumber[0]
	found, err := regexp.MatchString("^.+$", param)
	fmt.Printf("found=%v, err=%v", found, err)
	que := &dynamodb.QueryInput{
		TableName: &tName,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":Fnumber": {
				S: aws.String(param),},},
		KeyConditionExpression: aws.String("Fnumber = :Fnumber"),	}
	fmt.Println(que)
	resp, err := svc.Query(que)
	if err != nil {fmt.Println("err:", err)}
	messg := []Item{}
	dynamodbattribute.UnmarshalListOfMaps(resp.Items, &messg)
	msg, _ := json.Marshal(messg)
	json.NewEncoder(w).Encode(messg)
	fmt.Println(msg)}}
