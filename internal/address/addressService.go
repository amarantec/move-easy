package address

import (
	"context"
	"github.com/amarantec/move-easy/internal"
	"errors"
	"unicode/utf8"
	"unicode"
)

type IAddressService interface {
	GetAddress(ctx context.Context, userID int64) (internal.Address, error)
	AddOrUpdateAddress(ctx context.Context, address internal.Address) (int64, error)
}

type addressService struct {
	addressRepository IAddressRepository
}

func NewAddressService(repository IAddressRepository) IAddressService {
	return &addressService{addressRepository: repository}
}

func (s *addressService) GetAddress(ctx context.Context, userID int64) (internal.Address, error) {
	return s.addressRepository.GetAddress(ctx, userID)
}

func (s *addressService) AddOrUpdateAddress(ctx context.Context, address internal.Address) (int64, error) {
	if valid, err := validateAddress(address); err != nil || !valid {
		return internal.ZERO, err
	}

	return s.addressRepository.AddOrUpdateAddress(ctx, address)
}

func validateAddress(a internal.Address) (bool, error) {
	if a.UserID <= internal.ZERO {
		return false, ErrAddressUserIDInvalid
	}

	if a.Street == internal.EMPTY {
		return false, ErrAddressStreetEmpty
	} else if utf8.RuneCountInString(a.Street) < 3 || utf8.RuneCountInString(a.Street) > 100 {
		return false, ErrAddressStreetInvalid
	}

	if a.Number == internal.EMPTY {
		return false, ErrAddressNumberEmpty
	} else if utf8.RuneCountInString(a.Number) < 1 {
		return false, ErrAddressNumberInvalid
	}

	for _, char := range a.Number {
		if !unicode.IsDigit(char) {
			return false, ErrAddressNumberInvalid
		}
	}

	if a.CEP == internal.EMPTY {
		return false, ErrAddressCEPEmpty
	} else if utf8.RuneCountInString(a.CEP) != 8 {
		return false, ErrAddressCEPInvalid
	}

	for _, char := range a.CEP {
		if !unicode.IsDigit(char) {
			return false, ErrAddressCEPInvalid
		}
	}
	
	if a.Neighborhood == internal.EMPTY {
		return false, ErrAddressNeighborhoodEmpty
	} else if utf8.RuneCountInString(a.Neighborhood) < 3 || utf8.RuneCountInString(a.Neighborhood) > 100 {
		return false, ErrAddressNeighborhoodInvalid
	}

	if a.City == internal.EMPTY {
		return false, ErrAddressCityEmpty
	} else if utf8.RuneCountInString(a.City) < 3 || utf8.RuneCountInString(a.City) > 100 {
		return false, ErrAddressCityInvalid
	}

	if a.State == internal.EMPTY {
		return false, ErrAddressStateEmpty
	} else if utf8.RuneCountInString(a.State) != 2 {
		return false, ErrAddressStateInvalid
	}

	return true, nil
}

var ErrAddressUserIDInvalid = errors.New("address user id is empty or negative")
var ErrAddressStreetEmpty = errors.New("address street is empty")
var ErrAddressStreetInvalid = errors.New("address street must be between 3-100 characters")
var ErrAddressNumberEmpty = errors.New("address number is empty")
var ErrAddressNumberInvalid = errors.New("address number must contain digits in range 0-9")
var ErrAddressCEPEmpty = errors.New("address cep is empty")
var ErrAddressCEPInvalid = errors.New("address cep must must contain digits in range 0-9")
var ErrAddressNeighborhoodEmpty = errors.New("address neighborhood is empty")
var ErrAddressNeighborhoodInvalid = errors.New("address neighborhood must be between 3-100 characters")
var ErrAddressCityEmpty = errors.New("address city is empty")
var ErrAddressCityInvalid = errors.New("address city must be between 3-100 characters")
var ErrAddressStateEmpty = errors.New("address state is empty")
var ErrAddressStateInvalid = errors.New("address state must contain only 2 characters, example: RS, SP, RJ")
