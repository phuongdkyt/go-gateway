package language

import (
	"context"
	"golang.org/x/text/language"
	language2 "phuongnd/gateway/internal/domain/language"
)

var SupportedLanguageMatcher = language.NewMatcher([]language.Tag{
	language.Vietnamese, // The first language is used as fallback.
	language.English,
	language.AmericanEnglish,
	language.BritishEnglish,
})

type Translator interface {
	Translate(langTag language.Tag, key string) string
}

type remoteTranslator struct {
	langService language2.Service
}

func NewTranslator(langService language2.Service) (Translator, error) {
	if err := langService.Load(context.Background()); err != nil {
		return nil, err
	}
	return &remoteTranslator{langService: langService}, nil
}

func (r *remoteTranslator) Translate(langTag language.Tag, key string) string {
	switch langTag {
	case language.English, language.AmericanEnglish, language.BritishEnglish:
		return r.langService.Translate(context.Background(), key, language2.English)
	default:
		return r.langService.Translate(context.Background(), key, language2.Vietnamese)
	}
}
