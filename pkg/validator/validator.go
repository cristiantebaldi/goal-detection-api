package validator

import (
	"encoding/json"
	"reflect"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var trans ut.Translator
var validate *validator.Validate

type ValidatorConfig struct {
	dicionario map[string]map[string]string
	traducao   map[string]map[string]string
}

type ValidatorOption func(c *ValidatorConfig)

func WithDicionario(dicionario map[string]map[string]string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.dicionario = dicionario
	}
}

func WithTraducoes(traducao map[string]map[string]string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.traducao = traducao
	}
}

type RequestError struct {
	Fields []Field `json:"fields"`
}

type Field struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Errs    string `json:"errs"`
}

func (r RequestError) Error() string {
	payload, _ := json.Marshal(r)
	return string(payload)
}

func GetDefaultConfig() (ValidatorOption, ValidatorOption) {
	dicionario := WithDicionario(map[string]map[string]string{
		"pt": {
			"timestamp":       "Timestamp",
			"positions":       "Positions",
			"anchor_id":       "Anchor ID",
			"distance_meters": "Distance meters",
		},
	})

	traducoes := WithTraducoes(map[string]map[string]string{
		"pt": {
			"required": "{0} é obrigatório",
			"gte":      "{0} deve ser maior ou igual a {1}",
			"min":      "{0} deve ter no mínimo {1} itens",
		},
	})

	return dicionario, traducoes
}

func translationFunc(t ut.Translator, fe validator.FieldError) string {
	field, err := t.T(fe.Field())
	if err != nil {
		field = fe.Field()
	}
	msg, err := t.T(fe.Tag(), field, fe.Param())
	if err != nil {
		return fe.Error()
	}
	return msg
}

func Initialize(opts ...ValidatorOption) {
	config := &ValidatorConfig{}
	for _, opt := range opts {
		opt(config)
	}

	v := validator.New()

	enLocale := en.New()
	ptLocale := pt_BR.New()
	utrans := ut.New(enLocale, enLocale, ptLocale)
	trans, _ = utrans.GetTranslator("pt")

	for locale, dict := range config.dicionario {
		engine, _ := utrans.FindTranslator(locale)
		for key, translation := range dict {
			_ = engine.Add(key, translation, false)
		}
	}

	for locale, translation := range config.traducao {
		engine, _ := utrans.FindTranslator(locale)
		for tag, msg := range translation {
			_ = v.RegisterTranslation(tag, engine, func(t ut.Translator) error {
				return t.Add(tag, msg, false)
			}, translationFunc)
		}
	}

	validate = v
}

func ValidateStruct(payload any) error {
	if err := validate.Struct(payload); err != nil {
		return handleValidatorFieldError(payload, err)
	}
	return nil
}

func handleValidatorFieldError(data any, err error) error {
	errs := err.(validator.ValidationErrors)
	requestError := RequestError{}
	for _, e := range errs {
		field := Field{
			Field:   getFormatedField(data, e.Field()),
			Message: e.Translate(trans),
			Errs:    e.Error(),
		}
		requestError.Fields = append(requestError.Fields, field)
	}
	return requestError
}

func getFormatedField(data any, field string) string {
	if f, ok := reflect.TypeOf(data).Elem().FieldByName(field); ok {
		return f.Tag.Get("json")
	}
	return ""
}
