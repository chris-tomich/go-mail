package params

type ContentDisposition string

const (
	DispInline ContentDisposition = "inline"
	DispAttachment ContentDisposition = "attachment"
)
