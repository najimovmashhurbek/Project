package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	pb "github.com/najimovmashhurbek/project-api/user-service.ozim/genproto"
	//"google.golang.org/grpc/internal/status"
)

type userRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *pb.User) (*pb.CreatePostRes, error) {
	query := "insert into users (id,name,firstName,lastName,bio,phoneNumbers,createdAt,status) values ($1,$2,$3,$4,$5,$6,$7,$8)"
	time := time.Now()
	//id := uuid.New()
	_, err := r.db.Exec(query, user.Id, user.Name, user.FirstName, user.LastName, user.Bio, pq.Array(user.PhoneNumbers), time, user.Status)
	if err != nil {
		return nil, err
	}
	defer r.db.Close()
	for _, adres := range user.Adress {
		adid := uuid.New()
		query1 := "insert into adress (id,users_id,country,city,district,postalCodes) values ($1,$2,$3,$4,$5,$6)"
		_, err := r.db.Exec(query1, adid, user.Id, adres.Country, adres.City, adres.District, adres.PostalCodes)
		if err != nil {
			return nil, err
		}
	}
	return &pb.CreatePostRes{Status: true}, nil
}

func (r *userRepo) DeleteUser(delete *pb.DeleteById) (*pb.DeleteUserRes, error) {
	query := "DELETE FROM users WHERE id=$1"
	_, err := r.db.Exec(query, delete.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserRes{Status: true}, nil
}
func (r *userRepo) UpdateUser(user *pb.User) (*pb.UpdateUserRes, error) {
	query := "UPDATE users set name=$1,firstname=$2,lastname=$3,bio=$4,updateat=$5,status=$6,phonenumbers=$7 where id=$8"
	query1 := "UPDATE adress set country=$1,city=$2,district=$3,postalcodes=$4 where users_id=$5"
	time := time.Now()
	_, err := r.db.Exec(query, user.Name, user.FirstName, user.LastName, user.Bio, time, user.Status, pq.Array(user.PhoneNumbers), user.Id)
	if err != nil {
		return nil, err
	}
	for _, adres := range user.Adress {
		_, err := r.db.Exec(query1, adres.Country, adres.City, adres.District, adres.PostalCodes, adres.UserId)
		if err != nil {
			return nil, err
		}
	}
	return &pb.UpdateUserRes{Status: true}, nil

}
func (r *userRepo) GetAllUser(get *pb.GetAllById) (*pb.User, error) {
	var res pb.User
	//var resp []*pb.User
	query1 := "select id,country,city,district,postalcodes from adress where users_id=$1"
	query := "select id,name,firstname,lastname,bio,createdat,status from users where id=$1"
	rows, err := r.db.Query(query, get.Id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(
			&res.Id,
			&res.Name,
			&res.FirstName,
			&res.LastName,
			&res.Bio,
			&res.CreatedAt,
			&res.Status,
		)
		if err != nil {
			return nil, err
		}

	}
	row, err := r.db.Query(query1, get.Id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var adres pb.Adress
		err = row.Scan(
			&adres.Id,
			&adres.Country,
			&adres.City,
			&adres.District,
			&adres.PostalCodes,
		)
		if err != nil {
			return nil, err
		}
		res.Adress = append(res.Adress, &adres)
	}
	return &res, nil
}

func (r *userRepo) ListUsers(limit, page int64) ([]*pb.User, int64, error) {
	var (
		users []*pb.User
		count int64
	)
	offset := (page - 1) * limit
	query := "select id,firstname,lastname,bio,status,name,phonenumbers from users order by firstname OFFSET $1 LIMIT $2"
	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		var user pb.User
		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Bio,
			&user.Status,
			&user.Name,
			pq.Array(&user.PhoneNumbers),
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}
	countquery := "SELECT count(*) FROM users"
	err = r.db.QueryRow(countquery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}
