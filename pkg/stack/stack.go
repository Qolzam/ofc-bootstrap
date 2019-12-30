package stack

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/Qolzam/ofc-bootstrap/pkg/types"
)

type gatewayConfig struct {
	Origin           string
	WebDomain        string
	ExternalDomain   string
	CookieRootDomain string
	Scheme           string
	Domain           string
	WebSocketURL     string
}

type serverWebConfig struct {
	PrettyURL bool
}

type authConfig struct {
	RootDomain           string
	ClientId             string
	CustomersURL         string
	Scheme               string
	OAuthProvider        string
	OAuthProviderBaseURL string
}

// Apply creates `templates/gateway_config.yml` to be referenced by stack.yml
func Apply(plan types.Plan) error {
	scheme := "http"
	if plan.TLS {
		scheme += "s"
	}

	gwConfigErr := generateTemplate("gateway_config", plan, gatewayConfig{
		Origin:           plan.Origin,
		WebDomain:        plan.WebDomain,
		ExternalDomain:   plan.ExternalDomain,
		CookieRootDomain: plan.CookieRootDomain,
		Scheme:           scheme,
		Domain:           plan.Domain,
		WebSocketURL:     plan.WebSocketURL,
	})

	if gwConfigErr != nil {
		return gwConfigErr
	}

	if slackConfigErr := generateTemplate("server_web_config", plan, serverWebConfig{
		PrettyURL: plan.PrettyURL,
	}); slackConfigErr != nil {
		return slackConfigErr
	}

	isGitHub := plan.SCM == "github"
	stackErr := generateTemplate("stack", plan, stackConfig{
		GitHub: isGitHub,
	})

	copyConfigErr := copyConfigFiles()
	if copyConfigErr != nil {
		return copyConfigErr
	}

	if stackErr != nil {
		return stackErr
	}

	return nil
}

type builderConfig struct {
	ECR bool
}

type stackConfig struct {
	GitHub bool
}

func generateTemplate(fileName string, plan types.Plan, templateType interface{}) error {

	generatedData, err := applyTemplate("templates/"+fileName+".yml", templateType)
	if err != nil {
		return err
	}

	tempFilePath := "tmp/generated-" + fileName + ".yml"
	file, fileErr := os.Create(tempFilePath)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	_, writeErr := file.Write(generatedData)
	file.Close()

	if writeErr != nil {
		return writeErr
	}

	return nil
}

func copyConfigFiles() error {
	err := copyDirectory("templates/config", "./tmp/config")
	return err
}

func applyTemplate(templateFileName string, templateType interface{}) ([]byte, error) {
	data, err := ioutil.ReadFile(templateFileName)
	if err != nil {
		return nil, err
	}
	t := template.Must(template.New(templateFileName).Parse(string(data)))

	buffer := new(bytes.Buffer)

	executeErr := t.Execute(buffer, templateType)

	return buffer.Bytes(), executeErr
}
