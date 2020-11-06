package user

import (
	"github.com/mixmaru/my_contracts/core/domain/validators"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCorporationNameValue_NewPresidentNameValue(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		name, err := NewCorporationNameValue("会社名")
		assert.NoError(t, err)
		assert.Equal(t, CorporationNameValue{"会社名"}, name)
	})

	t.Run("会社名が空文字だった時", func(t *testing.T) {
		name, err := NewCorporationNameValue("")
		assert.Error(t, err)
		assert.Equal(t, CorporationNameValue{}, name)
	})

	t.Run("会社名が50文字を超えていた時", func(t *testing.T) {
		name, err := NewCorporationNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Error(t, err)
		assert.Equal(t, CorporationNameValue{}, name)
	})

	t.Run("会社名が50文字だった時", func(t *testing.T) {
		name, err := NewCorporationNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.NoError(t, err)
		assert.Equal(t, CorporationNameValue{"0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789"}, name)
	})
}

func TestCorporationNameValue_Validate(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		validErrs, err := CorporationNameValue{}.Validate("会社名")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		validErrs, err := CorporationNameValue{}.Validate("")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.EmptyStringValidError, validErrs[0])
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		validErrs, err := CorporationNameValue{}.Validate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.OverLengthStringValidError, validErrs[0])
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		validErrs, err := CorporationNameValue{}.Validate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
	})
}
