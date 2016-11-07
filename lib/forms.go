package lib

import (
	"bytes"
	"errors"
	"fmt"
	. "github.com/zballs/comit/util"
	"time"
)

var Fmt = fmt.Sprintf

const (
	ErrDecodeForm          = 100
	ErrFindForm            = 101
	ErrFormAlreadyResolved = 102
	ErrDecodingFormID      = 103
	field                  = "<strong style='opacity:0.8;'>%s</strong> <small>%s</small><br>"
	miniField              = "<strong style='opacity:0.8;'>%s</strong> <really-small>%s</really-small><br>"
)

type Form struct {
	SubmittedAt string
	Issue       string
	Location    string
	Description string
	Submitter   string //optional

	//==============//

	Status     string
	ResolvedBy string
	ResolvedAt string
}

// Functional Options

type Item func(*Form) error

func newForm(items ...Item) (*Form, error) {
	form := &Form{}
	for _, item := range items {
		err := item(form)
		if err != nil {
			return nil, err
		}
	}
	return form, nil
}

func setSubmittedAt(submittedAt string) Item {
	return func(form *Form) error {
		form.SubmittedAt = submittedAt
		return nil
	}
}

func setIssue(issue string) Item {
	return func(form *Form) error {
		form.Issue = issue
		return nil
	}
}

func setLocation(location string) Item {
	return func(form *Form) error {
		form.Location = location
		return nil
	}
}

func setDescription(description string) Item {
	return func(form *Form) error {
		// TODO: field validation
		form.Description = description
		return nil
	}
}

func setSubmitter(submitter string) Item {
	return func(form *Form) error {
		form.Submitter = submitter
		return nil
	}
}

func MakeAnonymousForm(issue, location, description string) (*Form, error) {
	submittedAt := time.Now().Local().String()
	form, err := newForm(
		setSubmittedAt(submittedAt),
		setIssue(issue),
		setLocation(location),
		setDescription(description))
	if err != nil {
		return nil, err
	}
	form.Status = "unresolved"
	return form, nil
}

func MakeForm(issue, location, description, submitter string) (*Form, error) {
	form, err := MakeAnonymousForm(issue, location, description)
	if err != nil {
		return nil, err
	}
	err = setSubmitter(submitter)(form)
	if err != nil {
		return nil, err
	}
	return form, nil
}

func (form *Form) Resolved() bool {
	return form.Status == "resolved"
}

func (form *Form) Resolve(timestr, addr string) error {
	if form.Resolved() {
		return errors.New("form already resolved")
	}
	form.Status = "resolved"
	form.ResolvedAt = timestr
	form.ResolvedBy = addr
	return nil
}

func XOR(bytes []byte, items ...string) []byte {
	for _, item := range items {
		for idx, _ := range bytes {
			if idx < len(item) {
				bytes[idx] ^= byte(item[idx])
			} else {
				break
			}
		}
	}
	return bytes
}

// TODO: add location
func (form *Form) ID() []byte {
	bytes := make([]byte, 16)
	daystr := ToTheMinute(form.SubmittedAt)
	bytes = XOR(bytes, daystr, form.Issue)
	return bytes
}

func (form *Form) Summary() string {
	status := "unresolved"
	if form.Resolved() {
		status = Fmt(
			"resolved at %v by <really-small>%v</really-small>",
			form.ResolvedAt,
			form.ResolvedBy)
	}
	submittedAt := ToTheMinute(form.SubmittedAt)
	var summary bytes.Buffer
	summary.WriteString(Fmt(field, "submitted", submittedAt))
	summary.WriteString(Fmt(field, "issue", form.Issue))
	summary.WriteString(Fmt(field, "location", form.Location))
	summary.WriteString(Fmt(field, "description", form.Description))
	summary.WriteString(Fmt(field, "status", status))
	if len(form.Submitter) == 64 {
		summary.WriteString(Fmt(miniField, "submitter", form.Submitter))
	}
	return summary.String()
}
