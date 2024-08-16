package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	pb "github.com/exam-5/Car-Wash-Auth-Service/genproto/user"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (p *UserRepo) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var query string
	var args []interface{}
	if req.Email == "" && req.Id == "" {
        return nil, errors.New("either email or id must be provided")
    }
	if req.Email != "" {
		query = `
            SELECT id, username, email, full_name, role, created_at
            FROM users
            WHERE email = $1
        `
		args = append(args, req.Email)
	} else if req.Id != "" {
		query = `
            SELECT id, username, email, full_name, role, created_at
            FROM users
            WHERE id = $1
        `
		args = append(args, req.Id)
	} else {
		return nil, errors.New("neither email nor id provided")
	}

	var user pb.GetProfileResponse
	err := p.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Id, &user.Username, &user.Email, &user.FullName, &user.Role,&user.CreatedAt,
	)
	if err != nil {
		log.Printf("Error getting profile: %v\n", err)
		return nil, err
	}
	return &user, nil
}

func (p *UserRepo) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	query := `
		UPDATE users
		SET username = $1, email = $2, full_name = $3, role = $4, updated_at = NOW()
		WHERE id = $5
		`

	_, err := p.db.ExecContext(ctx, query, req.Username, req.Email, req.FullName, req.Role, req.Id)
	if err != nil {
		log.Printf("Error update profile: %v\n", err)
		return nil, err
	}
	return &pb.UpdateProfileResponse{}, nil

}

func (p *UserRepo) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := p.db.ExecContext(ctx, query, req.NewPassword, req.Id)
	if err != nil {
		log.Printf("Error change password: %v\n", err)
		return nil, err
	}
	return &pb.ChangePasswordResponse{}, nil
}

func (p *UserRepo) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	query := `
	SELECT id, username, email, full_name, role, created_at
    FROM users`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Error getting all users: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var users []*pb.GetProfileResponse
	for rows.Next() {
		var user pb.GetProfileResponse
		err := rows.Scan(
			&user.Id, &user.Username, &user.Email, &user.FullName, &user.Role, &user.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning user: %v\n", err)
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v\n", err)
		return nil, err
	}

	return &pb.GetAllUsersResponse{Users: users}, nil

}

