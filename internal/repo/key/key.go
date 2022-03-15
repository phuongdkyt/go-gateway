package key_repo

//
//import (
//	"context"
//	"github.com/google/wire"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.uber.org/zap"
//	"phuongnd/gateway/internal/model/key"
//)
//
//var ProviderRepoSet = wire.NewSet(NewKeyRepo)
//
//type KeyRepo struct {
//	client *mongo.Client
//	ctx    context.Context
//	log    *zap.Logger
//}
//
//func (r *KeyRepo) Collection(collectionName string) *mongo.Collection {
//	//dbName := viper.GetString("MONGODB_DBNAME")
//	return r.client.Database("local2").Collection(collectionName)
//}
//
//func (srv *KeyRepo) Close() error {
//	if err := srv.client.Disconnect(srv.ctx); err != nil {
//		srv.log.Error("Failed to disconnect from mongodb", zap.Error(err))
//		return err
//	}
//
//	return nil
//}
//
//func NewKeyRepo(ctx context.Context, log *zap.Logger, client *mongo.Client) IRepo {
//	return &KeyRepo{
//		client: client,
//		ctx:    ctx,
//		log:    log,
//	}
//}
//
//func (srv *KeyRepo) Create(privateKey string, publicKey string) error {
//	keyCollection := srv.Collection("key")
//	_, err := keyCollection.InsertOne(srv.ctx, bson.M{
//		"private_key": privateKey,
//		"public_key":  publicKey,
//	})
//	if err != nil {
//		srv.log.Error("Error when insert key", zap.Error(err))
//		return err
//	}
//	return nil
//}
//
//func (srv *KeyRepo) FindPublicKey(publicKey string) (*key.Key, error) {
//	var keyResult *key.Key
//	keyCollection := srv.Collection("key")
//	err := keyCollection.FindOne(srv.ctx, bson.M{
//		"public_key": publicKey,
//	}).Decode(keyResult)
//	if err != nil {
//		srv.log.Error("Error when insert key", zap.Error(err))
//		return nil, err
//	}
//	return keyResult, nil
//}
//
//func (srv *KeyRepo) FindAll() ([]*key.Key, error) {
//	var lstKey []*key.Key
//	keyCollection := srv.Collection("key")
//	cursor, err := keyCollection.Find(srv.ctx, bson.M{})
//	cursor.Next(srv.ctx)
//	{
//		keyElem := &key.Key{}
//		err = cursor.Decode(keyElem)
//		if err != nil {
//			return nil, err
//		}
//		lstKey = append(lstKey, keyElem)
//
//	}
//	return lstKey, nil
//}
