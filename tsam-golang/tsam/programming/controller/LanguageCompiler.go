package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/web"
)

// LanguageCompilerController provide methods to communicate with external language compiler API.
type LanguageCompilerController struct {
	LanguageCompilerService *service.LanguageCompilerService
}

// NewLanguageCompilerController creates new instance of LanguageCompilerService.
func NewLanguageCompilerController(languageCompilerService *service.LanguageCompilerService) *LanguageCompilerController {
	return &LanguageCompilerController{
		LanguageCompilerService: languageCompilerService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *LanguageCompilerController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	//============================== Compiler Sphere API ==============================

	// Send code to external compiler (Sphere).
	router.HandleFunc("/tenant/{tenantID}/compiler-send-code-sphere",
		controller.SendCodeSphere).Methods(http.MethodPost)

	// Get the output from external compiler (Sphere).
	router.HandleFunc("/tenant/{tenantID}/compiler-send-code-sphere-output/{outputID}",
		controller.GetOutputSphere).Methods(http.MethodGet)

	// Get the result from external compiler (Sphere).
	router.HandleFunc("/tenant/{tenantID}/compiler-send-code-sphere-result",
		controller.GetResultSphere).Methods(http.MethodPost)

	//============================== Problems Sphere API ==============================

	// Add judge (Sphere).
	router.HandleFunc("/tenant/{tenantID}/add-judge-sphere",
		controller.AddJudgeSphere).Methods(http.MethodPost)

	// Add problem (Sphere).
	router.HandleFunc("/tenant/{tenantID}/add-problem-sphere",
		controller.AddProblemSphere).Methods(http.MethodPost)

	// Add test case (Sphere).
	router.HandleFunc("/tenant/{tenantID}/add-test-case-sphere",
		controller.AddTestCaseSphere).Methods(http.MethodPost)

	// Add submissiom (Sphere).
	router.HandleFunc("/tenant/{tenantID}/add-submission-sphere",
		controller.AddSubmissionSphere).Methods(http.MethodPost)

	//============================== Remote Code API ==============================

	// // Send code to external compiler (Remote Code).
	// router.HandleFunc("/tenant/{tenantID}/compiler-send-code-remote-code",
	// 	controller.SendCodeRemoteCode).Methods(http.MethodPost)

	log.NewLogger().Info("Langauge Compiler Routes Registered")
}

//==================================== Compiler Sphere API ====================================================

// SendCodeSphere sends code and compiler language id to external API (Sphere).
func (controller *LanguageCompilerController) SendCodeSphere(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================SendCodeSphere called==============================")

	// Create buckets.
	compiler := service.CompilerSphere{}
	result := service.OutputIDSphere{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &compiler)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Call send code service method.
	err = controller.LanguageCompilerService.SendCodeSphere(&compiler, &result)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, result)
}

// GetOutputSphere gets the output from the external API (Sphere).
func (controller *LanguageCompilerController) GetOutputSphere(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetOutputSphere called==============================")

	// Getting tenant id from param and parsing it to uuid.
	outputID := mux.Vars(r)["outputID"]

	// Create bucket for bytes.
	resultBytes:= []byte{}

	// Call get output service method.
	err := controller.LanguageCompilerService.GetOutputSphere(string(outputID), &resultBytes)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultBytes)
}

// GetResultSphere gets the result from the external API (Sphere).
func (controller *LanguageCompilerController) GetResultSphere(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetResultSphere called==============================")

	// Create bucket for result URL.
	resultURL := service.ResultURLSphere{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &resultURL)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Create bucket for bytes.
	resultBytes:= []byte{}

	// Call get output service method.
	err = controller.LanguageCompilerService.GetResultSphere(resultURL.URL, &resultBytes)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	
	// Writing Response with OK Status to ResponseWriter.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(resultBytes)
}

//==================================== Problems Sphere API ====================================================

// AddJudgeSphere adds testcase/master judge (Sphere).
func (controller *LanguageCompilerController) AddJudgeSphere(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddJudgeSphere called==============================")

	// Create buckets.
	judge := service.JudgeSphere{}
	resultBytes:= []byte{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &judge)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Call send code service method.
	err = controller.LanguageCompilerService.AddJudgeSphere(&judge, &resultBytes)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultBytes)
}

// AddProblemSphere adds problem (Sphere).
func (controller *LanguageCompilerController) AddProblemSphere(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddProblemSphere called==============================")

	// Create buckets.
	problem := service.ProblemSphere{}
	resultBytes:= []byte{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &problem)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Call send code service method.
	err = controller.LanguageCompilerService.AddProblemSphere(&problem, &resultBytes)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultBytes)
}

// AddTestCaseSphere adds test case (Sphere).
func (controller *LanguageCompilerController) AddTestCaseSphere(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddTestCaseSphere called==============================")

	// Create buckets.
	testCase := service.TestCaseSphere{}
	resultBytes:= []byte{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &testCase)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Call send code service method.
	err = controller.LanguageCompilerService.AddTestCaseSphere(&testCase, &resultBytes)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultBytes)
}


// AddSubmissionSphere adds submission (Sphere).
func (controller *LanguageCompilerController) AddSubmissionSphere(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddSubmissionSphere called==============================")

	// Create buckets.
	submission := service.SubmissionSphere{}
	resultBytes:= []byte{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &submission)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Call send code service method.
	err = controller.LanguageCompilerService.AddSubmissionSphere(&submission, &resultBytes)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultBytes)
}

//==================================== Remote Code API ====================================================

// // SendCodeRemoteCode sends code and compiler language id to external API (Remote Code).
// func (controller *LanguageCompilerController) SendCodeRemoteCode(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================SendCodeRemoteCode called==============================")

// 	// Create buckets.
// 	compiler := service.CompilerRemoteCode{}
// 	result := service.ResultRemoteCode{}

// 	// Unmarshal json.
// 	err := web.UnmarshalJSON(r, &compiler)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
// 		return
// 	}

// 	// Call send code service method.
// 	err = controller.LanguageCompilerService.SendCodeRemoteCode(&compiler, &result)
// 	if err != nil {
// 		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
// 		return
// 	}

// 	// Writing Response with OK Status to ResponseWriter.
// 	web.RespondJSON(w, http.StatusOK, result)
// }