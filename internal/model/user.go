package model

import (
	"context"
	"time"

	"github.com/isaiahwong/auth-go/internal/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Auth type
type Auth struct {
	Email                    string    `bson:"email" json:"email"`
	FirstName                string    `bson:"first_name" json:"first_name"`
	LastName                 string    `bson:"last_name" json:"last_name"`
	Name                     string    `bson:"name" json:"name"`
	Picture                  string    `bson:"picture" json:"picture"`
	Password                 string    `bson:"password" json:"password"`
	PasswordHashMethod       string    `bson:"password_hash_method" json:"password_hash_method"`
	PasswordResetID          string    `bson:"password_reset_id" json:"password_reset_id"`
	PasswordResetToken       string    `bson:"password_reset_token" json:"password_reset_token"`
	PasswordResetExpires     time.Time `bson:"password_reset_expires" json:"password_reset_expires"`
	PasswordModified         time.Time `bson:"password_modified" json:"password_modified"`
	Verified                 bool      `bson:"verified" json:"verified"`
	VerifiedDate             time.Time `bson:"verified_date" json:"verified_date"`
	VerificationToken        string    `bson:"verification_token" json:"verification_token"`
	VerificationTokenExpires time.Time `bson:"verification_token_expires" json:"verification_token_expires"`
	CaptchaChallengeTS       time.Time `bson:"captcha_challenge_ts" json:"captcha_challenge_ts"`
}

// User type
type User struct {
	ID         primitive.ObjectID `bson:"id,omitempty" json:"id"`
	Object     string             `bson:"object" json:"object" validate:" eq=user,required" `
	Auth       Auth               `bson:"auth" json:"auth"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	LoggedIn   time.Time          `bson:"logged_in" json:"logged_in"`
	LoggedInIP string             `bson:"logged_in_ip" json:"logged_in_ip"`
}

// Save persists the document. If context is nil, a default context with timeout will be used
// based on the store's Timeout variable
func (u *User) Save(ctx context.Context, s *store.MongoStore) (*primitive.ObjectID, error) {
	var cancel context.CancelFunc
	if s == nil {
		return nil, &store.StoreEmpty{}
	}
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), s.Timeout)
		defer cancel()
	}
	coll := s.Client.Database(s.Database).Collection("user")
	u.ID = primitive.NewObjectID()
	u.UpdatedAt = time.Now()
	u.CreatedAt = time.Now()

	// TODO: Validate
	// perr := p.val.Struct(u)
	// if perr != nil {
	// 	for _, err := range p.valcast(perr) {

	// 		fmt.Println("\nFIELD")
	// 		fmt.Println(err.Namespace())
	// 		fmt.Println(err.Field())
	// 		fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
	// 		fmt.Println(err.StructField())     // by passing alt name to ReportError like below
	// 		fmt.Println(err.Tag())
	// 		fmt.Println(err.ActualTag())
	// 		fmt.Println(err.Kind())
	// 		fmt.Println(err.Type())
	// 		fmt.Println(err.StructField(), err.Value(), err.Param())
	// 	}
	// }

	res, err := coll.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, &store.OIDTypeError{}
	}
	return &oid, nil
}
