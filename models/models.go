package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	Name           *string            `json:"name" validate:"required,min=2,max=30"`
	Password       *string            `json:"password" validate:"required"`
	Email          *string            `json:"email" validate:"required,email"`
	Phone          *string            `json:"phone" validate:"required"`
	Token          *string            `json:"token"`
	RefreshToken   *string            `json:"refreshToken"`
	CreatedTime    time.Time          `json:"createdTime"`
	UpdatedTime    time.Time          `json:"updatedTime"`
	UserId         string             `json:"useId"`
	UserCart       []ProductUser      `json:"usercart" bson:"usercart"`
	AddressDetails []Address          `json:"address" bson:"address"`
	OrderStatus    []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	ProductId   primitive.ObjectID `json:"_id" bson:"_id"`
	ProductName *string            `json:"product_name" bson:"product_name"`
	Price       *uint64            `json:"price"`
	Rating      *uint8             `json:"rating"`
	Image       *string            `json:"image"`
}

type ProductUser struct {
	ProductId   primitive.ObjectID `json:"_id" bson:"_id"`
	ProductName *string            `json:"product_name" bson:"product_name"`
	Price       *uint64            `json:"price" bson:"price"`
	Rating      *uint8             `json:"rating" bson:"rating"`
	Image       *string            `json:"image" bson:"image"`
}

type Address struct {
	AddressId primitive.ObjectID `json:"_id" bson:"_id"`
	House     *string            `json:"house_name" bson:"house_name"`
	Street    *string            `json:"street_name" bson:"street_name"`
	City      *string            `json:"city_name" bson:"city_name"`
}
type Order struct {
	OrderId       primitive.ObjectID `json:"_id" bson:"_id"`
	OrderCart     []ProductUser      `json:"order_list" bson:"order_list"`
	OrderTime     time.Time          `json:"order_time" bson:"order_time"`
	Price         int                `json:"total_price" bson:"total_price"`
	Discount      int                `json:"discount" bson:"discount"`
	PaymentMethod Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool
	COD     bool
}
