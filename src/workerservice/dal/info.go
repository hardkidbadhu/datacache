package dal

import (
	"context"
	"datamodels"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ UserInfo = &userInfDal{}
type UserInfo interface {
	GetData(int,int) ([]*datamodels.UserInfo, error)
}

type userInfDal struct {
	col *mongo.Collection
}

func NewUserInfDal(col *mongo.Collection) *userInfDal {
	return &userInfDal{
		col: col,
	}
}

func (u userInfDal) GetData(pageNo, limit int) (uList []*datamodels.UserInfo, err error) {

	var (
		findOpt *options.FindOptions
		cur     *mongo.Cursor
	)

	//calculate the paginate
	skip := (pageNo - 1) * limit

	findOpt = new(options.FindOptions)
	findOpt.SetSkip(int64(skip))
	findOpt.SetLimit(int64(limit))

	log.Println(int64(skip), int64(limit))

	cur, err = u.col.Find(context.Background(), bson.D{}, findOpt)
	defer cur.Close(context.Background())
	if err != nil {
		log.Printf("Error: Datastore - GetData - %s", err.Error())
		return
	}

	if err = cur.All(context.Background(), &uList); err != nil {
		log.Printf("Error: Datastore - GetData - %s", err.Error())
		return
	}

	return
}


