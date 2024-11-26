package auth

import (
	"AABBCCDD/app/db"
	"AABBCCDD/app/db/sqlc"
	"fmt"

	"github.com/ignoxx/ohmyskit/kit"
	v "github.com/ignoxx/ohmyskit/validate"
)

var profileSchema = v.Schema{
	"firstName": v.Rules(v.Min(3), v.Max(50)),
	"lastName":  v.Rules(v.Min(3), v.Max(50)),
}

type ProfileFormValues struct {
	ID        uint   `form:"id"`
	FirstName string `form:"firstName"`
	LastName  string `form:"lastName"`
	Email     string
	Success   string
}

func HandleProfileShow(kit *kit.Kit) error {
	auth := kit.Auth().(Auth)

	user, err := db.Get().FindUserByID(kit.Request.Context(), int64(auth.UserID))
	if err != nil {
		return err
	}

	formValues := ProfileFormValues{
		ID:        uint(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return kit.Render(ProfileShow(formValues))
}

func HandleProfileUpdate(kit *kit.Kit) error {
	var values ProfileFormValues
	errors, ok := v.Request(kit.Request, &values, profileSchema)
	if !ok {
		return kit.Render(ProfileForm(values, errors))
	}

	auth := kit.Auth().(Auth)
	if auth.UserID != values.ID {
		return fmt.Errorf("unauthorized request for profile %d", values.ID)
	}
	err := db.Get().UpdateUserFirstLastName(kit.Request.Context(), sqlc.UpdateUserFirstLastNameParams{
		ID:        int64(auth.UserID),
		FirstName: values.FirstName,
		LastName:  values.LastName,
	})

	if err != nil {
		return err
	}

	values.Success = "Profile successfully updated!"
	values.Email = auth.Email

	return kit.Render(ProfileForm(values, v.Errors{}))
}
