package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	service "github.com/techlabs/swabhav/tsam/college/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	college "github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// StudentController provides method to update, delete, add, get all, get one for student.
type StudentController struct {
	StudentService *service.StudentService
}

// NewStudentController creates new instance of StudentController.
func NewStudentController(studentService *service.StudentService) *StudentController {
	return &StudentController{
		StudentService: studentService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *StudentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all students by seminar id.
	router.HandleFunc("/tenant/{tenantID}/student/seminar/{seminarID}/limit/{limit}/offset/{offset}",
		controller.GetStudents).Methods(http.MethodGet)

	// Add one student.
	router.HandleFunc("/tenant/{tenantID}/student/seminar/{seminarID}/credential/{credentialID}",
		controller.AddStudent).Methods(http.MethodPost)

	// Add one student from seminar registration form.
	addStudent := router.HandleFunc("/tenant/{tenantID}/student-reg-form/seminar/{seminarID}",
		controller.AddStudentFromRegForm).Methods(http.MethodPost)

	// Update one student.
	router.HandleFunc("/tenant/{tenantID}/student/{studentID}/seminar/{seminarID}/credential/{credentialID}",
		controller.UpdateStudent).Methods(http.MethodPut)

	// Delete one student.
	router.HandleFunc("/tenant/{tenantID}/student/{studentID}/seminar/{seminarID}/credential/{credentialID}",
		controller.DeleteStudent).Methods(http.MethodDelete)

	// Update all students seminar talent registration field.
	router.HandleFunc("/tenant/{tenantID}/student/seminar/{seminarID}/credential/{credentialID}",
		controller.UpdateMultipleStudent).Methods(http.MethodPut)

	// Exculde routes.
	*exclude = append(*exclude, addStudent)

	log.NewLogger().Info("Student Routes Registered")
}

// AddStudent adds one student.
func (controller *StudentController) AddStudent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddStudent called=======================================")

	// Create bucket.
	student := college.Student{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the student variable with given data.
	if err := web.UnmarshalJSON(r, &student); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Declare err.
	var err error

	// Parse and set tenant ID.
	student.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of student.
	student.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set seminar ID.
	student.SeminarID, err = util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.StudentService.AddStudentWithCredential(&student); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Student added successfully")
}

// AddStudentFromRegForm adds one student from seminar registration form.
func (controller *StudentController) AddStudentFromRegForm(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddStudentFromRegForm called=======================================")

	// Create bucket.
	student := college.Student{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the student variable with given data.
	if err := web.UnmarshalJSON(r, &student); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Declare err.
	var err error

	// Parse and set tenant ID.
	student.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set seminar ID.
	student.SeminarID, err = util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.StudentService.AddStudent(&student); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Student added successfully")
}

// GetStudents gets all students.
func (controller *StudentController) GetStudents(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetStudents called=======================================")

	// Create bucket.
	students := []college.StudentDTO{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the r.Form.
	r.ParseForm()

	// Create bucket for total seminar count.
	var totalCount int

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting seminar id from param and parsing it to uuid.
	seminarID, err := util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Get limit and offset from param and convert it to int.
	limit, offset := web.GetLimitAndOffset(r)

	// Call get students method.
	if err := controller.StudentService.GetStudents(&students, tenantID, seminarID, limit, offset, &totalCount, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, students)
}

// UpdateStudent updates student.
func (controller *StudentController) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateStudent called=======================================")

	// Create bucket.
	student := college.Student{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the student variable with given data.
	err := web.UnmarshalJSON(r, &student)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := student.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set student ID to student.
	student.SeminarTalentRegistrationID, err = util.ParseUUID(params["studentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse student id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to student.
	student.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of student.
	student.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set student ID to student.
	student.SeminarID, err = util.ParseUUID(params["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.StudentService.UpdateStudent(&student); err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Student updated successfully")
}

// DeleteStudent deletes one student.
func (controller *StudentController) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteStudent called=======================================")

	// Create bucket
	student := college.Student{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set student ID.
	student.SeminarTalentRegistrationID, err = util.ParseUUID(params["studentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse student id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	student.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to student's DeletedBy field.
	student.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set student ID to student.
	student.SeminarID, err = util.ParseUUID(mux.Vars(r)["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.StudentService.DeleteStudent(&student); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Student deleted successfully")
}

// UpdateMultipleStudent updates multiple students' multiple fields.
func (controller *StudentController) UpdateMultipleStudent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateMultipleStudent called=======================================")

	// Create bucket.
	updateMultipleStudent := college.UpdateMultipleStudent{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the updateMultipleStudent variable with given data.
	err := web.UnmarshalJSON(r, &updateMultipleStudent)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set student ID to student.
	updateMultipleStudent.SeminarID, err = util.ParseUUID(mux.Vars(r)["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	updateMultipleStudent.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	updateMultipleStudent.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.StudentService.UpdateMultipleStudent(&updateMultipleStudent)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Students updated successfully")
}
