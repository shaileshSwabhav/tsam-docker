package contact

import "strings"

type ContactInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	// BCC []string
	// CC []string
}

func (email *ContactInfo) CreateMessageBody() (string, error) {
	// var messageBody string
	builder := strings.Builder{}
	err := stringWriter("Hey, You got a new Enquiry !", &builder)
	if err != nil {
		return "", err
	}
	err = stringWriter("Name :"+email.Name, &builder)
	if err != nil {
		return "", err
	}
	err = stringWriter("Email :"+email.Email, &builder)
	if err != nil {
		return "", err
	}
	err = stringWriter("Phone :"+email.Phone, &builder)
	if err != nil {
		return "", err
	}
	err = stringWriter("Subject :"+email.Subject, &builder)
	if err != nil {
		return "", err
	}
	err = stringWriter("Message :"+email.Message, &builder)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}
func stringWriter(msg string, builder *strings.Builder) error {
	// builder:= strings.Builder{}
	var err error
	_, err = builder.WriteString(msg)
	if err != nil {
		return err
	}
	_, err = builder.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}
