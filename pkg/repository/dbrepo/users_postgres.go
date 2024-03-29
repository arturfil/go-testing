package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"
	"webapp/pkg/data"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

type PostgresDBRepo struct {
    DB *sql.DB
}

func (m *PostgresDBRepo) Connection() *sql.DB {
    return m.DB
}

// AllUsers returns all users as a slice of *data.User
func (m *PostgresDBRepo) AllUsers() ([]*data.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

    defer cancel()

    query := `
        SELECT id, email, first_name, last_name, password, is_admin, created_at, updated_at
        FROM users
        ORDER BY last_name
    `

    rows, err := m.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*data.User

    for rows.Next() {
        var user data.User
        err := rows.Scan(
            &user.ID,
            &user.Email,
            &user.FirstName,
            &user.LastName,
            &user.Password,
            &user.IsAdmin,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            log.Println("Error scanning", err)
            return nil, err
        }

        users = append(users, &user)
    }

    return users, nil
}

// GetUser - when providing an id, you return the user
func (m *PostgresDBRepo) GetUser(id int) (*data.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    query := `
        SELECT id, email, first_name, last_name, password, is_admin, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    var user data.User
    row := m.DB.QueryRowContext(ctx, query, id)

    err := row.Scan(
        &user.ID,
        &user.Email,
        &user.FirstName,
        &user.LastName,
        &user.Password,
        &user.IsAdmin,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

// GetUserByEmail - given an email, this method will return a user
func (m *PostgresDBRepo) GetUserByEmail(email string) (*data.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    query := `
        SELECT id, email, first_name, last_name, password, is_admin, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    var user data.User
    row := m.DB.QueryRowContext(ctx, query, email)
    err := row.Scan(
        &user.ID,
        &user.Email,
        &user.FirstName,
        &user.LastName,
        &user.Password,
        &user.IsAdmin,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
}

//  InserUser - Create user, provided the necessary information
func (m *PostgresDBRepo) InsertUser(user data.User) (int, error) {
        ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
        defer cancel()

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
        if err != nil {
            return 0, err
        }

        var newID int
        query := `
            INSERT INTO users (email, first_name, last_name, password, is_admin, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7) returning id
        `
         err = m.DB.QueryRowContext(ctx, query, 
            user.Email,
            user.FirstName,
            user.LastName,
            hashedPassword,
            user.IsAdmin,
            time.Now(),
            time.Now(),
         ).Scan(&newID)
        
         if err != nil {
             return 0, err
         }

         return newID, nil
}

// UpdateUser - change a users info, provided the data and the ID included
func(m *PostgresDBRepo) UpdateUser(u data.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    query := `
        UPDATE users 
        SET email = $1, first_name = $2, last_name = $3, is_admin = $4, updated_at = $5
        WHERE id = $6
    `

    _, err := m.DB.ExecContext(ctx, query, 
       u.Email,
       u.FirstName,
       u.LastName,
       u.IsAdmin,
       time.Now(),
       u.ID,
    )

    if err != nil {
        return err
    }

    return nil
}

// DeleteUser - delete user once an id is providedf
func (m *PostgresDBRepo) DeleteUser(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    query := `DELETE FROM users where id = $1`
    _, err := m.DB.ExecContext(ctx, query, id)

    if err != nil {
        return err
    }
    return nil
}

// ResetPassword - This method will change the password
func (m *PostgresDBRepo) ResetPassword(id int, password string) error {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }

    query := `UPDATE users SET password = $1 WHERE id = $2`
    _, err = m.DB.ExecContext(ctx, query, hashedPassword, id)
    if err != nil {
        return err
    }

    return nil
}

// Insert User Image - This method will associalte an image to a user
func (m *PostgresDBRepo) InsertUserImage(i data.UserImage) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    var newID int
    stmt := `insert into user_images (user_id, file_name, created_at, updated_at)
        values ($1, $2, $3, $4) returning id`

    err := m.DB.QueryRowContext(ctx, stmt,
        i.UserID,
        i.FileName,
        time.Now(),
        time.Now(),
    ).Scan(&newID)

    if err != nil {
        return 0, err
    }

    return newID, nil
}
