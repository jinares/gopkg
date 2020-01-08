package xerr

import (
	"fmt"
	"testing"

	"github.com/jinares/gopkg/xtools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestXErr(t *testing.T) {
	xerr := XErr(OK, "ok", true)
	fmt.Println(xerr.Error())

}
func TestFromXErr(t *testing.T) {
	xerr := FromXErr(status.New(codes.InvalidArgument, "InvalidArgument").Err())

	fmt.Println(xerr.Error())
	st := status.New(codes.InvalidArgument, "InvalidArgument")
	fmt.Println("===================", st.Err(), xtools.JSONToStr(st))

}
