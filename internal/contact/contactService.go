package contact

import (
	"context"
	"errors"
	"unicode"
	"unicode/utf8"

	"github.com/amarantec/move-easy/internal"
)

type IContactService interface {
	SaveContact(ctx context.Context, contact internal.Contact) (int64, error)
	GetContact(ctx context.Context, userID, contactID int64) (internal.Contact, error)
	ListContacts(ctx context.Context, userID int64) ([]internal.Contact, error)
	UpdateContact(ctx context.Context, contact internal.Contact) (bool, error)
	DeleteContact(ctx context.Context, userID, contactID int64) (bool, error)
}

type contactService struct {
	contactRepository IContactRepository
}

func NewContactService(repository IContactRepository) IContactService {
	return &contactService{contactRepository: repository}
}

func (s *contactService) SaveContact(ctx context.Context, contact internal.Contact) (int64, error) {
	if valid, err := validateContact(contact); err != nil || !valid {
		return internal.ZERO, err
	}
	return s.contactRepository.SaveContact(ctx, contact)
}

func (s *contactService) GetContact(ctx context.Context, userID, contactID int64) (internal.Contact, error) {
	if userID <= internal.ZERO {
		return internal.Contact{}, ErrContactUserIDInvalid
	}
	return s.contactRepository.GetContact(ctx, userID, contactID)
}

func (s *contactService) ListContacts(ctx context.Context, userID int64) ([]internal.Contact, error) {
	if userID <= internal.ZERO {
		return []internal.Contact{}, ErrContactUserIDInvalid
	}
	return s.contactRepository.ListContacts(ctx, userID)
}

func (s *contactService) UpdateContact(ctx context.Context, contact internal.Contact) (bool, error) {
	if valid, err := validateContact(contact); err != nil || !valid {
		return false, err
	}

	return s.contactRepository.UpdateContact(ctx, contact)
}

func (s *contactService) DeleteContact(ctx context.Context, userID, contactID int64) (bool, error) {
	if userID <= internal.ZERO || contactID <= internal.ZERO {
		return false, ErrContactIDInvalid
	}
	return s.contactRepository.DeleteContact(ctx, userID, contactID)
}

func validateContact(c internal.Contact) (bool, error) {
	if c.UserID <= internal.ZERO {
		return false, ErrContactUserIDInvalid
	}

	if c.Name == internal.EMPTY {
		return false, ErrContactNameEmpty
	} else if utf8.RuneCountInString(c.Name) < 3 || utf8.RuneCountInString(c.Name) > 100 {
		return false, ErrContactNameInvalid
	}

	if c.DDI == internal.EMPTY {
		return false, ErrContactDDIEmpty
	} else if utf8.RuneCountInString(c.DDI) != 3 {
		return false, ErrContactDDIInvalid
	}

	for _, char := range c.DDI {
		if !unicode.IsDigit(char) {
			return false, ErrContactDDIInvalid
		}
	}

	if c.DDD == internal.EMPTY {
		return false, ErrContactDDDEmpty
	} else if utf8.RuneCountInString(c.DDD) != 3 {
		return false, ErrContactDDDInvalid
	}

	for _, char := range c.DDD {
		if !unicode.IsDigit(char) {
			return false, ErrContactDDDInvalid
		}
	}

	if c.PhoneNumber == internal.EMPTY {
		return false, ErrContactPhoneNumberEmpty
	} else if utf8.RuneCountInString(c.PhoneNumber) != 9 {
		return false, ErrContactPhoneNumberInvalid
	}

	for _, char := range c.PhoneNumber {
		if !unicode.IsDigit(char) {
			return false, ErrContactPhoneNumberInvalid
		}
	}

	return true, nil
}

var (
	ErrContactUserIDInvalid      = errors.New("contact user id is empty or negative")
	ErrContactIDInvalid          = errors.New("contact id is empty or negative")
	ErrContactNameEmpty          = errors.New("contact name is empty")
	ErrContactNameInvalid        = errors.New("contact name must be between 3-100 characters")
	ErrContactDDIEmpty           = errors.New("contact ddi is empty")
	ErrContactDDIInvalid         = errors.New("contact ddi must have only 3 digits")
	ErrContactDDDEmpty           = errors.New("contact ddd is empty")
	ErrContactDDDInvalid         = errors.New("contact ddd must have only 3 digits")
	ErrContactPhoneNumberEmpty   = errors.New("contact phone number is empty")
	ErrContactPhoneNumberInvalid = errors.New("contact phone number must have 9 digits in range 0-9")
)
