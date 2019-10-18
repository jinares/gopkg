package xerr

type (
	XCode int
)

const (
	OK              XCode = 0
	Canceled        XCode = 1
	Unknown         XCode = 2
	InvalidArgument XCode = 3

	DeadlineExceeded XCode = 4

	NotFound XCode = 5

	AlreadyExists XCode = 6

	PermissionDenied XCode = 7

	ResourceExhausted XCode = 8

	FailedPrecondition XCode = 9

	Aborted XCode = 10

	OutOfRange XCode = 11

	Unimplemented XCode = 12

	Internal XCode = 13

	Unavailable XCode = 14

	DataLoss XCode = 15

	Unauthenticated XCode = 16
)
