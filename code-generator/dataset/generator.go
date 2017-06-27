package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/chuckpreslar/inflect"
	"github.com/fatih/camelcase"
	"github.com/reivaj05/GoConfig"
)

type generateOptions struct {
	path          string
	serviceName   string
	fileExtension string
	fileTemplate  string
	data          interface{}
}

func (op *generateOptions) getFilePath() string {
	return op.path + op.serviceName + op.fileExtension
}

type goAPITemplateData struct {
	ServiceName      string
	UpperServiceName string
}

type protoAPITemplateData struct {
	ServiceName      string
	ResourcePath     string
	UpperServiceName string
}

type EndpointsData struct {
	Services []string
}

// TODO: Add tests
// TODO: Refactor code when done

func Generate(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Service name not provided")
	}
	createServices(args...)
	updateServerEndpoints()
	return runProtoGenScript()
}

func createServices(args ...string) {
	// TODO: Create properly template files (protos and go files for services)
	// TODO: Implement rest of template files
	basePath := joinPath()
	for _, serviceName := range args {
		if err := generateFiles(basePath, serviceName); err != nil {
			fmt.Errorf("Service not created: " + err.Error())
			rollback(serviceName)
		}
	}

}

func joinPath() string {
	const relativePath = "/src/github.com/reivaj05/apigateway"
	goPath := os.Getenv("GOPATH")
	return goPath + relativePath
}

func generateFiles(path, serviceName string) error {
	if err := generateAPIFile(path, serviceName); err != nil {
		return err
	}
	if err := generateServiceFile(path, serviceName); err != nil {
		return err
	}
	if err := generateProtoFiles(path, serviceName); err != nil {
		return err
	}
	return nil
}

func generateAPIFile(path, serviceName string) error {
	path += "/api/" + serviceName + "/"
	return _generateFile(&generateOptions{
		path:          path,
		serviceName:   serviceName,
		fileExtension: ".go",
		fileTemplate:  "goAPI.txt",
		data: &goAPITemplateData{
			ServiceName:      serviceName,
			UpperServiceName: inflect.Titleize(serviceName),
		},
	})

}

func generateServiceFile(path, serviceName string) error {
	path += "/services/" + serviceName + "/"
	return _generateFile(&generateOptions{
		path:          path,
		serviceName:   serviceName,
		fileExtension: ".go",
		fileTemplate:  "goService.txt",
		data: &goAPITemplateData{
			ServiceName:      serviceName,
			UpperServiceName: inflect.Titleize(serviceName),
		},
	})
}

func generateProtoFiles(path, serviceName string) error {
	sp := camelcase.Split(serviceName)[0]
	if err := _generateFile(&generateOptions{
		path:          path + "/protos/api/",
		serviceName:   serviceName,
		fileExtension: ".proto",
		fileTemplate:  "protoAPI.txt",
		data: &protoAPITemplateData{
			ServiceName:      serviceName,
			ResourcePath:     sp,
			UpperServiceName: inflect.Titleize(serviceName),
		},
	}); err != nil {
		return err
	}
	return _generateFile(&generateOptions{
		path:          path + "/protos/services/",
		serviceName:   serviceName,
		fileExtension: ".proto",
		fileTemplate:  "protoService.txt",
		data: struct {
			ServiceName      string
			UpperServiceName string
		}{
			ServiceName:      serviceName,
			UpperServiceName: inflect.Titleize(serviceName),
		},
	})
}

func _generateFile(options *generateOptions) error {
	file, err := _createFile(options)
	if err != nil {
		return err
	}
	return _writeTemplateContent(file, options)
}

func _createFile(options *generateOptions) (*os.File, error) {
	err := os.MkdirAll(options.path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return os.Create(options.getFilePath())
}

func _writeTemplateContent(file *os.File, options *generateOptions) error {
	defer file.Close()
	tmpl := template.Must(template.ParseFiles(
		GoConfig.GetConfigStringValue("templatesPath") + options.fileTemplate),
	)
	return tmpl.Execute(file, options.data)
}

func rollback(serviceName string) {
	path := joinPath()
	os.RemoveAll(path + "/api/" + serviceName)
	os.RemoveAll(path + "/services/" + serviceName)
	os.Remove(path + "/protos/api/" + serviceName + ".proto")
	os.Remove(path + "/protos/services/" + serviceName + ".proto")
}

func updateServerEndpoints() {
	services := getServicesNames()
	if err := updateServerFiles(&EndpointsData{Services: services}); err != nil {
		// TODO: Rollback
		// rollback(serviceName)
	}

}

func getServicesNames() (services []string) {
	files, _ := ioutil.ReadDir(GoConfig.GetConfigStringValue("protosPath"))
	for _, file := range files {
		services = append(services, strings.TrimSuffix(file.Name(), ".proto"))
	}
	return services
}

func updateServerFiles(data *EndpointsData) error {
	if err := _updateFile(data, "registeredGRPCEndpoints.go", "grpcEndpoints.txt"); err != nil {
		return err
	}
	return _updateFile(data, "registeredHTTPEndpoints.go", "httpEndpoints.txt")
}

func _updateFile(data *EndpointsData, endpointsFile, fileTemplate string) error {
	serverPath := GoConfig.GetConfigStringValue("serverPath")
	file, err := os.OpenFile(serverPath+endpointsFile, os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	return _writeTemplateContent(file, &generateOptions{
		data: data, fileTemplate: fileTemplate})
}

func runProtoGenScript() error {
	cmd := exec.Command("/bin/sh", GoConfig.GetConfigStringValue("protoGenPath"))
	return cmd.Run()
}
