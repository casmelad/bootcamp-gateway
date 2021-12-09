package users

import (
	"context"

	"gopkg.in/go-playground/validator.v9"
)

//Service - the interface for the users logic
type Service interface {
	Create(context.Context, User) (int, error)
	GetByEmail(context.Context, string) (User, error)
	GetAll(context.Context) ([]User, error)
	Update(context.Context, User) error
	Delete(context.Context, int) error
}

//UserService - the implementation for the users logic
type UserService struct {
	repository Repository
}

//NewUserService - returns a UserService type pointer
func NewUserService(repo Repository) *UserService {
	return &UserService{repository: repo}
}

//Create - validates business rules and sends a user to the repository
func (us *UserService) Create(ctx context.Context, usr User) (int, error) {

	v := validator.New()

	if errVal := v.Struct(usr); errVal != nil {
		return 0, ErrInvalidData
	}

	dbUser, err := us.repository.GetByEmail(ctx, usr.Email)

	if err != nil {
		return 0, err
	}

	if dbUser.ID > 0 {

		return 0, ErrUserAlreadyExists
	}

	newID, errAdd := us.repository.Add(ctx, usr)

	if errAdd != nil {
		return 0, ErrInternalError
	}

	return newID, nil

}

//GetByEmail - retrieves the information of a user based on the email address
func (us *UserService) GetByEmail(ctx context.Context, email string) (User, error) {

	dbUser, err := us.repository.GetByEmail(ctx, email)

	if err != nil {
		return User{}, err
	}

	if dbUser.ID == 0 {
		return User{}, ErrNotFound
	}

	return dbUser, nil

}

//GetAll -  gets all the existing users
func (us *UserService) GetAll(ctx context.Context) ([]User, error) {

	users, err := us.repository.GetAll(ctx)

	if err != nil {
		return []User{}, ErrInternalError
	}

	return users, nil
}

//Update - validates the data and updates the user information
func (us *UserService) Update(ctx context.Context, usr User) error {

	/* v := validator.New()

	if errVal := v.Struct(usr); errVal != nil {
		return ErrInvalidData
	} */

	usrToUpdate, errU := us.repository.GetByID(ctx, usr.ID)

	if errU != nil {
		return errU
	}

	if usrToUpdate.ID == 0 {
		return ErrNotFound
	}

	usr.ID = usrToUpdate.ID

	if err := us.repository.Update(ctx, usr); err != nil {
		return ErrInternalError
	}

	return nil
}

//Delete - removes a user
func (us *UserService) Delete(ctx context.Context, usrID int) error {

	if usrID < 1 {
		return ErrInvalidData
	}

	usrToUpdate, err := us.repository.GetByID(ctx, usrID)

	if err != nil {
		return err
	}

	if usrToUpdate.ID == 0 {
		return ErrNotFound
	}

	if errD := us.repository.Delete(ctx, usrID); errD != nil {
		return ErrInternalError
	}

	return nil
}
