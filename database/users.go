package database

import (
	"database/sql"
	"fmt"
	"time"
)

// User operations

// CreateUser creates a new user
func (d *Database) CreateUser(user *User) error {
	query := `INSERT INTO users (username, email, password, first_name, last_name, role, is_active, company_id, intro_viewed, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	result, err := d.db.Exec(query, user.Username, user.Email, user.Password, user.FirstName, user.LastName, 
		user.Role, user.IsActive, user.CompanyID, user.IntroViewed, now, now)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting user ID: %v", err)
	}

	user.ID = int(id)
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// GetUserByID retrieves a user by ID
func (d *Database) GetUserByID(id int) (*User, error) {
	query := `SELECT id, username, email, password, first_name, last_name, role, is_active, company_id, intro_viewed, last_login, created_at, updated_at 
			  FROM users WHERE id = ?`
	
	user := &User{}
	var lastLogin sql.NullTime
	
	err := d.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.IsActive, &user.CompanyID, &user.IntroViewed, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (d *Database) GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, email, password, first_name, last_name, role, is_active, company_id, intro_viewed, last_login, created_at, updated_at 
			  FROM users WHERE username = ?`
	
	user := &User{}
	var lastLogin sql.NullTime
	
	err := d.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.IsActive, &user.CompanyID, &user.IntroViewed, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (d *Database) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, email, password, first_name, last_name, role, is_active, company_id, intro_viewed, last_login, created_at, updated_at 
			  FROM users WHERE email = ?`
	
	user := &User{}
	var lastLogin sql.NullTime
	
	err := d.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Role, &user.IsActive, &user.CompanyID, &user.IntroViewed, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return user, nil
}

// GetUsersByCompany retrieves all users for a specific company
func (d *Database) GetUsersByCompany(companyID int) ([]User, error) {
	query := `SELECT id, username, email, password, first_name, last_name, role, is_active, company_id, intro_viewed, last_login, created_at, updated_at 
			  FROM users WHERE company_id = ? ORDER BY created_at DESC`
	
	rows, err := d.db.Query(query, companyID)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		var lastLogin sql.NullTime
		
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName,
			&user.Role, &user.IsActive, &user.CompanyID, &user.IntroViewed, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %v", err)
		}

		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}

		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (d *Database) UpdateUser(user *User) error {
	query := `UPDATE users SET username = ?, email = ?, first_name = ?, last_name = ?, role = ?, is_active = ?, updated_at = ? 
			  WHERE id = ?`
	
	now := time.Now()
	_, err := d.db.Exec(query, user.Username, user.Email, user.FirstName, user.LastName, 
		user.Role, user.IsActive, now, user.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	user.UpdatedAt = now
	return nil
}

// UpdateUserPassword updates a user's password
func (d *Database) UpdateUserPassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password = ?, updated_at = ? WHERE id = ?`
	
	now := time.Now()
	_, err := d.db.Exec(query, hashedPassword, now, userID)
	if err != nil {
		return fmt.Errorf("error updating user password: %v", err)
	}

	return nil
}

// UpdateUserIntroViewed updates a user's intro viewed status
func (d *Database) UpdateUserIntroViewed(userID int, viewed bool) error {
	query := `UPDATE users SET intro_viewed = ?, updated_at = ? WHERE id = ?`
	
	now := time.Now()
	_, err := d.db.Exec(query, viewed, now, userID)
	if err != nil {
		return fmt.Errorf("error updating user intro viewed status: %v", err)
	}

	return nil
}

// UpdateUserLastLogin updates a user's last login time
func (d *Database) UpdateUserLastLogin(userID int) error {
	query := `UPDATE users SET last_login = ?, updated_at = ? WHERE id = ?`
	
	now := time.Now()
	_, err := d.db.Exec(query, now, now, userID)
	if err != nil {
		return fmt.Errorf("error updating user last login: %v", err)
	}

	return nil
}

// DeleteUser deletes a user (soft delete by setting is_active to false)
func (d *Database) DeleteUser(userID int) error {
	query := `UPDATE users SET is_active = 0, updated_at = ? WHERE id = ?`
	
	now := time.Now()
	_, err := d.db.Exec(query, now, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

// HardDeleteUser permanently deletes a user from the database
func (d *Database) HardDeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = ?`
	
	_, err := d.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error hard deleting user: %v", err)
	}

	return nil
}