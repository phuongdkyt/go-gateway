package options

import (
	"context"
	"encoding/json"
	"github.com/duclm2609/status"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"net/http"
	"phuongnd/gateway/api"
	lang "phuongnd/gateway/internal/language"
)

type ProtoErrorHandler struct {
	MessageTranslator lang.Translator
}

func (peh ProtoErrorHandler) Handle(ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	request *http.Request,
	err error) {

	const fallback = `{"message": "Unexpected error happens :("}`

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(err)))

	errResponse := status.Convert(err)

	lg, _ := request.Cookie("lang")
	accept := request.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(lang.SupportedLanguageMatcher, lg.String(), accept)

	jErr := json.NewEncoder(w).Encode(api.ErrorResponse{
		Message: peh.MessageTranslator.Translate(tag, errResponse.Message()),
		Details: errResponse.Details(),
	})

	if jErr != nil {
		w.WriteHeader(runtime.HTTPStatusFromCode(codes.Internal))
		_, _ = w.Write([]byte(fallback))
	}
}
