package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/repository"
)

// LanguageCompilerService provide methods to communicate with external language compiler API.
type LanguageCompilerService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewLanguageCompilerService creates new instance of LanguageCompilerService.
func NewLanguageCompilerService(db *gorm.DB, repository repository.Repository) *LanguageCompilerService {
	return &LanguageCompilerService{
		DB:         db,
		Repository: repository,
	}
}

// CompilerSphere is for sending information to extenmal compiler API (Sphere).
type CompilerSphere struct {
	CompilerID                 int   `json:"compilerId"` 
	Source string `json:"source"` 
	Input string `json:"input"` 
}

// CompilerRemoteCode is for sending information to extenmal compiler API (Remote Code).
type CompilerRemoteCode struct {
	SourceCode string `json:"sourceCode"` 
}

// OutputIDSphere is for storing result sent back from external API (Sphere).
type OutputIDSphere struct {
	ID                 int   `json:"id"` 
}

// ResultURLSphere is for getting result url sent from front end.
type ResultURLSphere struct {
	URL                string   `json:"url"` 
}

// JudgeSphere is add/update of judge (Sphere).
type JudgeSphere struct {
	CompilerID                 int   `json:"compilerId"` 
	Name string `json:"name"` 
	Source string `json:"source"` 
	TypeID int `json:"typeId"` 
}

// ProblemSphere is add/update of problem (Sphere).
type ProblemSphere struct {
	Name string `json:"name"` 
	MasterjudgeID int `json:"masterjudgeId"` 
	TypeID int `json:"typeId"` 
}

// TestCaseSphere is add/update of test case (Sphere).
type TestCaseSphere struct {
	Input string `json:"input"` 
	Output string `json:"output"` 
	JudgeID int `json:"judgeId"` 
	ProblemID int `json:"problemID"` 
}

// SubmissionSphere is add/update of submission (Sphere).
type SubmissionSphere struct {
	ProblemID int `json:"problemId"` 
	Source string `json:"source"` 
	Tests string `json:"tests"`
	CompilerID int `json:"compilerId"` 
}

// // ResultRemoteCode is for storing result sent back from external API (Remote Code).
// type ResultRemoteCode struct {
// 	// Date string `json:"date"` 
// 	ExpectedOutput string `json:"expectedOutput"` 
// 	Output string `json:"output"` 
// 	Status string `json:"status"` 
// }

//==================================== Compiler Sphere API ====================================================

// SendCodeSphere sends code and compiler language id to external API (Sphere).
func (service *LanguageCompilerService) SendCodeSphere(compiler *CompilerSphere, result *OutputIDSphere) error {

	// URL for sending source code and other info.
	url := "https://user_0b71f8472c.compilers.sphere-engine.com/api/v4/submissions?access_token=664b16978717b76aabe8bdcc95c53759"

	// Convert info to json.
	response, error := json.Marshal(compiler)
	if error != nil{
		return error
	}

	// Create the external API request.
    request, err := http.NewRequest("POST", url, bytes.NewBuffer(response))

	// Set the header for the request.
    request.Header.Set("Content-Type", "application/json")

	// Create http client.
    client := &http.Client{}

	// Call the external API.
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }

	// Close the response body.
    defer resp.Body.Close()

	// Read the response body.
    body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.NewHTTPError(errors.ErrorCodeReadWriteFailure, http.StatusBadRequest)
	}

	// If body is empty.
	if len(body) == 0 {
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	// Convert the body from json to struct format.
	err = json.Unmarshal(body, result)
	if err != nil {
		return errors.NewHTTPError(errors.ErrorCodeInvalidJSON, http.StatusBadRequest)
	}

	return nil
}

// GetOutputSphere gets the output from the external API (Sphere).
func (service *LanguageCompilerService) GetOutputSphere(outputID string, resultBytes *[]byte) error {

	// URL for sending output id.
	url := "https://user_0b71f8472c.compilers.sphere-engine.com/api/v4/submissions/" + outputID + "?access_token=664b16978717b76aabe8bdcc95c53759"

	// Call the get output external API.
	response, err := http.Get(url)
   	if err != nil {
		return errors.NewHTTPError("Output not recevied", http.StatusBadRequest)
   	}

	// Read the response body.
    tempResultBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.NewHTTPError(errors.ErrorCodeReadWriteFailure, http.StatusBadRequest)
	}

	// If body is empty.
	if len(tempResultBytes) == 0 {
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	*resultBytes = tempResultBytes

	return nil
}

// GetResultSphere gets the result from the external API (Sphere).
func (service *LanguageCompilerService) GetResultSphere(resultURL string, resultBytes *[]byte) error {

	// Call the get output external API.
	response, err := http.Get(resultURL)
   	if err != nil {
		return errors.NewHTTPError("Output not recevied", http.StatusBadRequest)
   	}

	// Read the response body.
    tempResultBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.NewHTTPError(errors.ErrorCodeReadWriteFailure, http.StatusBadRequest)
	}

	// If body is empty.
	if len(tempResultBytes) == 0 {
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	*resultBytes = tempResultBytes

	return nil
}

//==================================== Problems Sphere API ====================================================

// AddJudgeSphere adds testcase/master judge (Sphere).
func (service *LanguageCompilerService) AddJudgeSphere(judge *JudgeSphere, resultBytes *[]byte) error {

	// URL for sending source code and other info.
	url := "https://user_0b71f8472c.problems.sphere-engine.com/api/v4/judges?access_token=078a0bed4ac6687f5b92453621f5dc24"

	// Convert info to json.
	response, error := json.Marshal(judge)
	if error != nil{
		return error
	}

	// Create the external API request.
    request, err := http.NewRequest("POST", url, bytes.NewBuffer(response))

	// Set the header for the request.
    request.Header.Set("Content-Type", "application/json")

	// Create http client.
    client := &http.Client{}

	// Call the external API.
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }

	// Close the response body.
    defer resp.Body.Close()

	// Read the response body.
    tempResultBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.NewHTTPError(errors.ErrorCodeReadWriteFailure, http.StatusBadRequest)
	}

	// If body is empty.
	if len(tempResultBytes) == 0 {
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	*resultBytes = tempResultBytes

	return nil
}

// AddProblemSphere adds problem (Sphere).
func (service *LanguageCompilerService) AddProblemSphere(problem *ProblemSphere, resultBytes *[]byte) error {

	// URL for sending source code and other info.
	url := "https://user_0b71f8472c.problems.sphere-engine.com/api/v4/problems?access_token=078a0bed4ac6687f5b92453621f5dc24"

	// Convert info to json.
	response, error := json.Marshal(problem)
	if error != nil{
		return error
	}

	// Create the external API request.
    request, err := http.NewRequest("POST", url, bytes.NewBuffer(response))

	// Set the header for the request.
    request.Header.Set("Content-Type", "application/json")

	// Create http client.
    client := &http.Client{}

	// Call the external API.
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }

	// Close the response body.
    defer resp.Body.Close()

	// Read the response body.
    tempResultBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.NewHTTPError(errors.ErrorCodeReadWriteFailure, http.StatusBadRequest)
	}

	// If body is empty.
	if len(tempResultBytes) == 0 {
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	*resultBytes = tempResultBytes

	return nil
}

// AddTestCaseSphere adds test case (Sphere).
func (service *LanguageCompilerService) AddTestCaseSphere(testCase *TestCaseSphere, resultBytes *[]byte) error {

	// URL for sending source code and other info.
	url := "https://user_0b71f8472c.problems.sphere-engine.com/api/v4/problems/" + strconv.Itoa(testCase.ProblemID) + "/testcases?access_token=078a0bed4ac6687f5b92453621f5dc24"

	// Convert info to json.
	response, error := json.Marshal(testCase)
	if error != nil{
		return error
	}

	// Create the external API request.
    request, err := http.NewRequest("POST", url, bytes.NewBuffer(response))

	// Set the header for the request.
    request.Header.Set("Content-Type", "application/json")

	// Create http client.
    client := &http.Client{}

	// Call the external API.
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }

	// Close the response body.
    defer resp.Body.Close()

	// Read the response body.
    tempResultBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.NewHTTPError(errors.ErrorCodeReadWriteFailure, http.StatusBadRequest)
	}

	// If body is empty.
	if len(tempResultBytes) == 0 {
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	*resultBytes = tempResultBytes

	return nil
}

// AddSubmissionSphere adds submission (Sphere).
func (service *LanguageCompilerService) AddSubmissionSphere(submission *SubmissionSphere, resultBytes *[]byte) error {

	// URL for sending source code and other info.
	url := "https://user_0b71f8472c.problems.sphere-engine.com/api/v4/submissions?access_token=078a0bed4ac6687f5b92453621f5dc24"

	// Convert info to json.
	response, error := json.Marshal(submission)
	if error != nil{
		return error
	}

	// Create the external API request.
    request, err := http.NewRequest("POST", url, bytes.NewBuffer(response))

	// Set the header for the request.
    request.Header.Set("Content-Type", "application/json")

	// Create http client.
    client := &http.Client{}

	// Call the external API.
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }

	// Close the response body.
    defer resp.Body.Close()

	// Read the response body.
    tempResultBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.NewHTTPError(errors.ErrorCodeReadWriteFailure, http.StatusBadRequest)
	}

	// If body is empty.
	if len(tempResultBytes) == 0 {
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	*resultBytes = tempResultBytes

	return nil
}

//==================================== Remote Code API ====================================================

// // SendCodeRemoteCode sends code and compiler language id to external API (Remote Code).
// func (service *LanguageCompilerService) SendCodeRemoteCode(compiler *CompilerRemoteCode, result *ResultRemoteCode) error {

// 	// URL for send source code and other info.
// 	url := "http://localhost:8081/compiler/java?sourceCode=" + compiler.SourceCode 

// 	// Convert info to json.
// 	response, error := json.Marshal(compiler)
// 	if error != nil{
// 		return error
// 	}

// 	// Create the external API request.
//     request, err := http.NewRequest("POST", url, bytes.NewBuffer(response))

// 	// Set the header for the request.
//     request.Header.Set("Content-Type", "application/json")

// 	// Create http client.
//     client := &http.Client{}

// 	// Call the external API.
//     resp, err := client.Do(request)
//     if err != nil {
//         panic(err)
//     }

// 	// Close the response body.
//     defer resp.Body.Close()

// 	// Read the response body.
//     body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return errors.NewHTTPError("Unable to read response", http.StatusBadRequest)
// 	}

// 	// If body is empty.
// 	if len(body) == 0 {
// 		return errors.NewHTTPError("Response is empty", http.StatusBadRequest)
// 	}

// 	var dat map[string]interface{}

// 	// Convert the body from json to struct format.
// 	err = json.Unmarshal(body, &dat)
// 	if err != nil {
// 		return errors.NewHTTPError("Cannot Unmarshal JSON", http.StatusBadRequest)
// 	}

// 	fmt.Println("Result@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
// 	fmt.Println(dat)

// 	return nil
// }


