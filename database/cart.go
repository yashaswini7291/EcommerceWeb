package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/yashaswini7291/ecommerceWeb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ProductNotFound       = errors.New("Produt Not Found:(")
	ErrCantDecodeProducts = errors.New("Produt Not Found:(")
	ErrUserIdIsInValid    = errors.New("Invalid User:(")
	ErrCantupdateUser     = errors.New("Cannot Add This Product To Cart:(")
	ErrCantRemoveItemCart = errors.New("Item Cant Be Removed From Cart:(")
	ErrCantGetItem        = errors.New("Unable To Get Item From Cart:(")
	ErrCantBuyCartItem    = errors.New("Cannot Update The Purchase:(")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDB, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ProductNotFound
	}
	var productCart []models.ProductUser

	err = searchFromDB.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsInValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantupdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsInValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsInValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.OrderId = primitive.NewObjectID()
	orderCart.OrderTime = time.Now()
	orderCart.OrderCart = make([]models.ProductUser, 0)
	orderCart.PaymentMethod.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$userCart"}}}}
	grouping := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "total", Value: bson.D{
				{Key: "$sum", Value: "$usercart.price"},
			}},
		}},
	}
	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getUserCart []bson.M
	if err = currentResults.All(ctx, &getUserCart); err != nil {
		panic(err)
	}
	var totalPrice int32

	for _, userItem := range getUserCart {
		price := userItem["total"]
		totalPrice = price.(int32)
	}

	orderCart.Price = int(totalPrice)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "orders", Value: orderCart},
		}},
	}

	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		log.Println(err)
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)

	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].orderList": bson.M{"$each": getCartItems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(ErrCantupdateUser)
	}

	userCartEmpty := make([]models.ProductUser, 0)

	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: userCartEmpty}}}}
	_, err = userCollection.UpdateMany(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsInValid
	}

	var productDetails models.ProductUser
	var orderDetails models.Order

	orderDetails.OrderId = primitive.NewObjectID()
	orderDetails.OrderTime = time.Now()
	orderDetails.OrderCart = make([]models.ProductUser, 0)
	orderDetails.PaymentMethod.COD = true
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
	}
	orderDetails.Price = int(*productDetails.Price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "orders", Value: orderDetails},
		}},
	}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].orderList": productDetails}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}
