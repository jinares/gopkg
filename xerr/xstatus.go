package xerr

import (
	"fmt"

	//"github.com/golang/protobuf/proto"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// statusError is an alias of a status proto.  It implements error and Status,
// and a nil statusError should never be returned by this package.
type XStatus spb.Status

func (se *XStatus) Error() string {
	p := (*spb.Status)(se)
	return fmt.Sprintf("code = %s desc = %s", codes.Code(p.GetCode()), p.GetMessage())
}

func (se *XStatus) GRPCStatus() *XStatus {
	return se
}
