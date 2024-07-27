package v1

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	URI                  string               `yaml:"uri"`
	Username             string               `yaml:"username"`
	Password             string               `yaml:"password"`
	ConnectionPoolDetail ConnectionPoolDetail `yaml:"connectionPoolDetail"`
}

type ConnectionPoolDetail struct {
	MaxPoolSize     uint64 `yaml:"maxPoolSize"`
	MinPoolSize     uint64 `yaml:"minPoolSize"`
	MaxIdleTime     int64  `yaml:"maxIdleTime"`
	MaxConnIdleTime int64  `yaml:"maxConnIdleTime"`
	ConnectTimeout  int64  `yaml:"connectTimeout"`
}

func Connect(clientOptions *options.ClientOptions, poolOptions *options.ClientOptions) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, poolOptions, clientOptions)
	if err != nil {
		panic(err)
	}
	return client, ctx, cancel, err
}

func Ping(client *mongo.Client, ctx context.Context) error {

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	log.Println("Ping to mongo server successfully")
	return nil
}

func (mongo *Mongo) InitializeMongoConnection() (*mongo.Client, error) {
	log.Println("Initializing mongo db connection...")

	url := mongo.URI
	credential := options.Credential{
		Username: mongo.Username,
		Password: mongo.Password,
	}

	mongoClientOptions := options.Client().ApplyURI(url)

	// Creating connection pool options
	poolOptions := options.Client().SetMaxPoolSize(mongo.ConnectionPoolDetail.MaxPoolSize)
	poolOptions.SetMinPoolSize(mongo.ConnectionPoolDetail.MinPoolSize)
	poolOptions.SetMaxConnIdleTime(time.Minute * time.Duration(mongo.ConnectionPoolDetail.MaxConnIdleTime))
	poolOptions.SetConnectTimeout(time.Second * time.Duration(mongo.ConnectionPoolDetail.ConnectTimeout))

	if mongo.Username != "" && mongo.Password != "" {
		mongoClientOptions.SetAuth(credential)
	}

	client, ctx, _, err := Connect(mongoClientOptions, poolOptions)
	if err != nil {
		log.Println("Faield to connect to mongo server")
		panic(err)
	}

	err = Ping(client, ctx)
	if err != nil {
		log.Println("Faield to ping mongo server")
		panic(err)
	}

	log.Println("Connected to MongoDB")

	return client, nil
}
