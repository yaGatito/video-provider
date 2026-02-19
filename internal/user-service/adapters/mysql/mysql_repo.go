package mysql

import (
	"database/sql"
	"log"
	"time"
	"video-provider/internal/user-service/domain"
	"video-provider/internal/user-service/ports"
)

type SQLUserRepository struct {
	DB *sql.DB
}

var _ ports.UserRepository = (*SQLUserRepository)(nil)

func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	log.Printf("Craeted SQLUserRepository\n")
	return &SQLUserRepository{DB: db}
}

func (r *SQLUserRepository) Create(user *domain.User, password string, passwordSalt string) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO users (name, lastname, email, passwordHash, passwordSalt, createdAt, status, isAdmin) VALUES(?, ?, ?, ?, ?, NOW(), ?, ?)`, user.Name, user.LastName, user.Email, password, passwordSalt, user.Status, user.IsAdmin)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SQLUserRepository) FindByID(id int64) (*domain.User, error) {
	var user domain.User
	var createdAt string
	err := r.DB.QueryRow(`SELECT name, lastname, email, createdAt FROM users WHERE id = ?`, id).Scan(&user.Name, &user.LastName, &user.Email, &createdAt)
	if err != nil {
		log.Fatal("Error querying user by ID:", err)
		return nil, err
	}
	if user.Email == "" {
		log.Printf("User with ID %d not found\n", id)
		return nil, sql.ErrNoRows
	}
	log.Printf("SQLUserRepository.FindByID: %v\n", err)

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		log.Printf("Error parsing createdAt: %v\n", err)
		return nil, err
	}
	return &user, nil
}
