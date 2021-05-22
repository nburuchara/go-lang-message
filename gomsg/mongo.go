package main




import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)
// DBName Database name.
const DBName = "gomsg"
var (
    conn       *mongo.Client
    ctx        = context.Background()
    connString string = "mongodb+srv://user_zhou:AB1Ck2Wg5MS3MfCq@cluster0.rktlj.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
)
// Database collections.
var (
    privateroomsCollection = "privaterooms"
    roomsCollection = "rooms"
	usersCollection = "users"
    msgsCollection = "msgs"
    privatemsgsCollection = "privatemsgs"
)
// createDBSession Create a new connection with the database.
func createDBSession() error {
    var err error
    conn, err = mongo.Connect(ctx, options.Client().
        ApplyURI(connString))
    if err != nil {
        return err
    }
    err = conn.Ping(ctx, nil)
    if err != nil {
        return err
    }
    return nil
}

