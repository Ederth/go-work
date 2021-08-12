package error

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	kHttp "github.com/go-kratos/kratos/v2/transport/http"
)

type ErrResp struct {
	Code      int32             `json:"code"`
	Reason    string            `json:"reason"`
	Message   string            `json:"message"`
	Metadata  map[string]string `json:"metadata"`
	Timestamp int64             `json:"timestamp"`
}

func Encoder(w http.ResponseWriter, r *http.Request, se error) {
	e := se.(*errors.Error)

	er := ErrResp{
		Code:      e.Code,
		Reason:    e.Reason,
		Message:   e.Message,
		Metadata:  e.Metadata,
		Timestamp: time.Now().Unix(),
	}

	codec, _ := kHttp.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(er)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", strings.Join([]string{"application", codec.Name()}, "/"))
	if sc, ok := se.(interface {
		StatusCode() int
	}); ok {
		w.WriteHeader(sc.StatusCode())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(body)
}
