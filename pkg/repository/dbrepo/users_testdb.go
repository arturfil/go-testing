package dbrepo

import (
	"database/sql"
	"errors"
	"time"
	"webapp/pkg/data"
)

type TestDBRepo struct {
    
}

func (m *TestDBRepo) Connection() *sql.DB {
    return nil
}

// AllUsers returns all users as a slice of *data.User
func (m *TestDBRepo) AllUsers() ([]*data.User, error) {
    var users []*data.User 
    return users, nil
}

// GetUser - when providing an id, you return the user
func (m *TestDBRepo) GetUser(id int) (*data.User, error) {
    var user = data.User{ID: 1}
    return &user, nil
}

// GetUserByEmail - given an email, this method will return a user
func (m *TestDBRepo) GetUserByEmail(email string) (*data.User, error) {
    if email == "admin@example.com" {
        user := data.User {
            ID: 1,
            FirstName: "Admin",
            LastName: "User",
            Email: "admin@example.com",
            Password: "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
            IsAdmin: 1,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        return &user, nil
    }
    return nil, errors.New("not found")
}

//  InserUser - Create user, provided the necessary information
func (m *TestDBRepo) InsertUser(user data.User) (int, error) {
     return  1, nil
}

// UpdateUser - change a users info, provided the data and the ID included
func(m *TestDBRepo) UpdateUser(u data.User) error {
    return nil
}

// DeleteUser - delete user once an id is providedf
func (m *TestDBRepo) DeleteUser(id int) error {
    return nil
}

// ResetPassword - This method will change the password
func (m *TestDBRepo) ResetPassword(id int, password string) error {
    return nil
}

// Insert User Image - This method will associalte an image to a user
func (m *TestDBRepo) InsertUserImage(i data.UserImage) (int, error) {
    return 1, nil
}
