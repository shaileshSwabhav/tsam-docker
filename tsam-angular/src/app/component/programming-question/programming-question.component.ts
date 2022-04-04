import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IProgrammingQuestion, IProgrammingQuestionIsActive, IProgrammingQuestionSolution, IProgrammingQuestionSolutionDTO, IProgrammingQuestionTestCase, IProgrammingQuestionType, ProgrammingQuestionService } from 'src/app/service/programming-question/programming-question.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';


@Component({
  selector: 'app-programming-question',
  templateUrl: './programming-question.component.html',
  styleUrls: ['./programming-question.component.css']
})
export class ProgrammingQuestionComponent implements OnInit {

  // Components.
  programmingQuestionTypeList: IProgrammingQuestionType[]
  programmingQuestionLevel: any[]
  programmingLanguageList: any[]

  // Programming question.
  programmingQuestions: IProgrammingQuestion[]
  totalProgrammingQuestion: number
  selectedQuestion: IProgrammingQuestion

  // Form.
  programmingQuestionSearchForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  offset: number
  paginationString: string

  // Modal.
  modalRef: any

  // Access.
  permission: IPermission
  roleName: string

  // Flags.
  isViewClicked: boolean
  isUpdateClicked: boolean
  hasAnyTalentAnswered: boolean

  // Search.
  searchFormValue: any
  isSearched: boolean

  // Child elements.
  @ViewChild("drawer") drawer: any
  @ViewChild("deleteModal") deleteModal: any
  @ViewChild("programmingQuestionModal") programmingQuestionModal: any
  @ViewChild("ckeditorExample") ckeditorExample: any
  @ViewChild("ckeditorQuestion") ckeditorQuestion: any
  @ViewChild("ckeditorOption") ckeditorOption: any
  @ViewChild("ckeditorConstraints") ckeditorConstraints: any

  // Example and Configuration.
  exampleAndConfig: any

  // Programmimg question example.
  exampleValue: string

  // Cke editor configuration.
  ckeEditorCommonConfig: any
  ckeEditorExampleConfig: any
  ckeEditorQuestionConfig: any
  ckeEditorOptionConfig: any
  ckeEditorConstraintsConfig: any

  //****************************PROGRAMMING QUESTION SOLUTION*************************************** */

  // Flags.
  showSolutionForm: boolean
  isOperationSolutionUpdate: boolean

  // Programming Question Solution.
  solutionList: IProgrammingQuestionSolutionDTO[]
  solutionForm: FormGroup
  solutionCount: number

  // Modal.
  @ViewChild('solutionFormModal') solutionFormModal: any

  //****************************PROGRAMMING QUESTION TEST CASE*************************************** */

  // Flags.
  showTestCaseForm: boolean
  isOperationTestCaseUpdate: boolean

  // Programming question test case.
  testCaseList: IProgrammingQuestionTestCase[]
  testCaseForm: FormGroup
  testCaseCount: number

  // Modal.
  @ViewChild('testCaseFormModal') testCaseFormModal: any

  // Params.
  operation: string

  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private questionService: ProgrammingQuestionService,
    private generalService: GeneralService,
    private urlConstant: UrlConstant,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private route: ActivatedRoute,
    private router: Router,
    private localService: LocalService,
    private role: Role,
  ) {
    this.initializeVariables()
    this.createSearchProgrammingForm()
    this.getAllComponents()
  }

  // Initialize global variables.
  initializeVariables(): void {

    // Permision.
    // Get permissions from menus using utilityService function.
    this.roleName = this.localService.getJsonValue("roleName")
    if (this.roleName == this.role.ADMIN || this.roleName == this.role.SALES_PERSON) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.TRAINING_PROGRAMMING_QUESTION)
    }
    if (this.roleName == this.role.FACULTY) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.BANK_PROGRAMMING_QUESTION)
    }

    // Components.
    this.programmingQuestions = []
    this.programmingQuestionTypeList = []
    this.programmingQuestionLevel = []
    this.programmingLanguageList = []

    // Flags.
    this.isSearched = false
    this.isViewClicked = false
    this.isUpdateClicked = false

    this.limit = 5
    this.offset = 0
    this.totalProgrammingQuestion = 0

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Programming Questions"

    //****************************PROGRAMMING QUESTION SOLUTION*************************************** */

    // Flags.
    this.showSolutionForm = false
    this.isOperationSolutionUpdate = false
    this.hasAnyTalentAnswered = false

    // Programming Question Solution.
    this.solutionList = [] as IProgrammingQuestionSolutionDTO[]

    //****************************PROGRAMMING QUESTION TEST CASE*************************************** */

    // Programming Question Test case.
    this.testCaseList = []

    // Flags.
    this.showTestCaseForm = false
    this.isOperationTestCaseUpdate = false
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  ngAfterViewInit() {
    if (this.operation && this.operation == "add") {
      this.onAddNewQuestionClick()
      this.operation = null
      return
    }
  }

  // =============================================================CREATE FORMS==========================================================================

  // Create programming question search form.
  createSearchProgrammingForm(): void {
    this.programmingQuestionSearchForm = this.formBuilder.group({
      label: new FormControl(null),
      isActive: new FormControl(null),
      programmingType: new FormControl(null),
      programmingConcept: new FormControl(null),
    })
  }

  // Create new programming question test case form.
  createTestCaseForm(): void {
    this.testCaseForm = this.formBuilder.group({
      id: new FormControl(null),
      programmingQuestionID: new FormControl(null),
      input: new FormControl(null, [Validators.required, Validators.maxLength(500)]),
      output: new FormControl(null, [Validators.required, Validators.maxLength(500)]),
      explanation: new FormControl(null, [Validators.maxLength(500)]),
      isHidden: new FormControl(false, [Validators.required]),
      isActive: new FormControl(true, [Validators.required]),
    })
  }

  // Create new programming question solution form.
  createSolutionForm(): void {
    this.solutionForm = this.formBuilder.group({
      id: new FormControl(),
      solution: new FormControl(null, [Validators.required, Validators.maxLength(2000)]),
      programmingLanguage: new FormControl(null, [Validators.required]),
      programmingQuestionID: new FormControl(null)
    })
  }

  // =============================================================PROGRAMMING QUESTION CRUD FUNCTIONS==========================================================================

  // On clicking add new programming question button.
  onAddNewQuestionClick(): void {
    this.isViewClicked = false
    this.isUpdateClicked = false
    this.openModal(this.programmingQuestionModal, "xl")
  }

  // Programming question added by chil component.
  programmingQuestionAdded(programmingQuestionID: string): void{
    this.modalRef.close()
  }

  // On clicking view programming question button.
  onViewProgrammingQuestionClick(question: IProgrammingQuestion): void {
    this.isViewClicked = true
    this.isUpdateClicked = false
    this.selectedQuestion = question
    this.openModal(this.programmingQuestionModal, "xl")
  }

  // Update programming question.
  updateProgrammingQuestion(programmingQuestion: IProgrammingQuestion): void {
    this.spinnerService.loadingMessage = "Update Programming Question"
    this.questionService.updateProgrammingQuestion(programmingQuestion).subscribe((response: any) => {
      this.modalRef.close()
      this.getProgrammingQuestions()
      alert("Programming question successfully updated")
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  // Update is active field of programming question.
  updateIsActive(isActive: boolean, questionID: string): void {
    if (confirm("Are you sure you want to change the status of this programming question ?")) {
      this.spinnerService.loadingMessage = "Updating programming question status"
      let questionIsActive: IProgrammingQuestionIsActive = {
        id: questionID,
        isActive: isActive
      }
      this.questionService.updateProgrammingQuestionIsActive(questionIsActive).subscribe((response: any) => {
        alert(response)
        this.getProgrammingQuestions()
      }, (error) => {
        console.error(error)
        if (error.error?.error) {
          alert(error.error?.error)
          return
        }
        alert(error.statusText)
      })
    }
  }

  // On clicking delete programming question button.
  onDeleteProgrammingQuestionClick(questionID: string): void {
    this.openModal(this.deleteModal, 'md').result.then(() => {
      this.deleteProgrammingQuestion(questionID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete programming question.
  deleteProgrammingQuestion(questionID: string): void {
    this.spinnerService.loadingMessage = "Deleting Programming Question"

    this.questionService.deleteProgrammingQuestion(questionID).subscribe((response: any) => {
      this.modalRef.close()
      this.getProgrammingQuestions()
      alert("Programming question successfully deleted")
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  // =============================================================PROGRAMMING QUESTION SEARCH FUNCTIONS==========================================================================

  // Reset programming question search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.createSearchProgrammingForm()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    let url: string
    if (this.roleName == this.role.ADMIN || this.roleName == this.role.SALES_PERSON) {
      url = this.urlConstant.TRAINING_PROGRAMMING_QUESTION
    }
    if (this.roleName == this.role.FACULTY) {
      url = this.urlConstant.BANK_PROGRAMMING_QUESTION
    }
    this.router.navigate([url])
  }

  // Reset programming question search form.
  resetSearchForm(): void {
    this.programmingQuestionSearchForm.reset()
    this.programmingQuestionSearchForm.get('isActive').setValue("1")
  }

  // Search programming questions.
  searchProgrammingQuestion(): void {
    this.searchFormValue = { ...this.programmingQuestionSearchForm?.value }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
    let flag: boolean = true

    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field]
      } else {
        this.isSearched = true
        flag = false
      }
    }

    // No API call on empty search.
    if (flag) {
      return
    }
    this.changePage(1)
  }

  // On clicking search button og programming question search form.
  searchAndCloseDrawer(): void {
    this.drawer.toggle()
    this.searchProgrammingQuestion()
  }

  // Get searched programming questions.
  getSearchedProgrammingQuestions(): void {
    this.spinnerService.loadingMessage = "Searching Programming Questions"

    this.programmingQuestions = []
    this.totalProgrammingQuestion = 0
    this.questionService.getProgrammingQuestions(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      this.programmingQuestions = response.body
      this.totalProgrammingQuestion = response.headers.get('X-Total-Count')
    }, (err: any) => {
      this.totalProgrammingQuestion = 0
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.setPaginationString()

    })
  }

  // =======================================OTHER FUNCTIONS FOR PROGRAMMING QUESTION==========================================

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetProgrammingQuestion() {
    let queryParams = this.route.snapshot.queryParams
    if (queryParams.operation) {
      this.operation = queryParams.operation
      // queryParams = {}
    }
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getProgrammingQuestions()
      return
    }
    this.programmingQuestionSearchForm.patchValue(queryParams)
    this.searchProgrammingQuestion()
  }

  // Page change function.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getProgrammingQuestions()
  }

  // Set the pagination string.
  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalProgrammingQuestion < end) {
      end = this.totalProgrammingQuestion
    }
    if (this.totalProgrammingQuestion == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end}`
  }

  // Assign selected programming question.
  assignSelectedQuestion(id: string): void {
    let length = this.programmingQuestions.length
    for (let index = 0; index < length; index++) {
      if (this.programmingQuestions[index].id == id) {
        this.selectedQuestion = this.programmingQuestions[index]
      }
    }
  }

  // Dismiss programming question modal.
  dismissProgrammingQuestinModal(): void{
    this.modalRef.close()
  }

  //*********************************************CRUD FUNCTIONS FOR PROGRAMMING QUESTION SOLUTION************************************************************

  // On clicking programming question solution popup in question list.
  getSolutionsForSelectedquestion(questionID: any): void {
    this.isOperationSolutionUpdate = false
    this.solutionList = []
    this.assignSelectedQuestion(questionID)
    this.getAllSolutions()
    this.showSolutionForm = false
    this.openModal(this.solutionFormModal, 'xl')
  }

  // Get all programming question solutions by programming question id.
  getAllSolutions(): void {
    this.spinnerService.loadingMessage = "Getting All Programming Question Solutions"
    this.questionService.getSolutionsByQuestions(this.selectedQuestion.id).subscribe(response => {
      this.solutionList = response
      this.getProgrammingQuestions()
    }, err => {
      console.error(this.utilService.getErrorString(err))
    })
  }

  // On clicking add new programming question solution button.
  onAddNewSolutionButtonClick(): void {
    this.isOperationSolutionUpdate = false
    this.createSolutionForm()
    this.showSolutionForm = true
  }

  // Add programming question solution.
  addSolution(): void {
    this.spinnerService.loadingMessage = "Adding Programming Question Solution"
    let solution: IProgrammingQuestionSolution = this.solutionForm.value
    this.patchIDFromObjectsForSolution(solution)
    this.questionService.addSolution(this.solutionForm.value, this.selectedQuestion.id).subscribe((response: any) => {
      this.showSolutionForm = false
      this.getAllSolutions()
      alert(response)
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  // Update programming question solution form.
  onUpdateSolutionButtonClick(index: number): void {
    this.isOperationSolutionUpdate = true
    this.createSolutionForm()
    this.showSolutionForm = true
    this.solutionForm.patchValue(this.solutionList[index])
  }

  // Update programming question solution.
  updateSolution(): void {
    this.spinnerService.loadingMessage = "Updating Programming Question Solution"
    let solution: IProgrammingQuestionSolution = this.solutionForm.value
    this.patchIDFromObjectsForSolution(solution)
    this.questionService.updateSolution(this.solutionForm.value, this.selectedQuestion.id).subscribe((response: any) => {
      this.showSolutionForm = false
      this.getAllSolutions()
      alert(response)
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error.error)
        return
      }
      alert("Check connection")
    })
  }

  // Delete programming question solution.
  deleteSolution(solutionID: string): void {
    if (confirm("Are you sure you want to delete the programming question solution?")) {
      this.spinnerService.loadingMessage = "Deleting Programming Question Solution"
      this.questionService.deleteSolution(solutionID, this.selectedQuestion.id).subscribe((response: any) => {
        this.getAllSolutions()
        alert(response)
      }, (error) => {
        console.error(error)
        if (error.error) {
          alert(error.error)
          return
        }
        alert(error.statusText)
      })
    }
  }

  // Validate programming question solution form.
  validateSolutionForm(): void {
    if (this.solutionForm.invalid) {
      this.solutionForm.markAllAsTouched()
      return
    }
    if (this.isOperationSolutionUpdate) {
      this.updateSolution()
      return
    }
    this.addSolution()
  }

  // Extract ID from objects in programming question solution form.
  patchIDFromObjectsForSolution(solution: IProgrammingQuestionSolution): void {
    if (this.solutionForm.get('programmingLanguage').value) {
      solution.programmingLanguageID = this.solutionForm.get('programmingLanguage').value.id
      delete solution['programmingLanguage']
    }
  }

  //*********************************************CRUD FUNCTIONS FOR PROGRAMMING QUESTION TEST CASE************************************************************

  // On clicking programming question test case popup in question list.
  getTestCasesForSelectedquestion(questionID: any): void {
    this.isOperationTestCaseUpdate = false
    this.testCaseList = []
    this.assignSelectedQuestion(questionID)
    this.getAllTestCases()
    this.showTestCaseForm = false
    this.openModal(this.testCaseFormModal, 'xl')
  }

  // Get all programming question test cases by programming question id.
  getAllTestCases(): void {
    this.spinnerService.loadingMessage = "Getting All Programming Question Test Cases"


    this.questionService.getTestCasesByQuestions(this.selectedQuestion.id).subscribe(response => {
      this.testCaseList = response
      this.getProgrammingQuestions()
    }, err => {
      console.error(this.utilService.getErrorString(err))
    })
  }

  // On clicking add new programming question test case button.
  onAddNewTestCaseButtonClick(): void {
    this.isOperationTestCaseUpdate = false
    this.createTestCaseForm()
    this.showTestCaseForm = true
  }

  // Add programming question test case.
  addTestCase(): void {
    this.spinnerService.loadingMessage = "Adding Programming Question Test Case"
    this.questionService.addTestCase(this.testCaseForm.value, this.selectedQuestion.id).subscribe((response: any) => {
      this.showTestCaseForm = false
      this.getAllTestCases()
      alert(response)
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  // Update programming question test case form.
  onUpdateTestCaseButtonClick(index: number): void {
    this.isOperationTestCaseUpdate = true
    this.createTestCaseForm()
    this.showTestCaseForm = true
    this.testCaseForm.patchValue(this.testCaseList[index])
  }

  // Update programming question test case.
  updateTestCase(): void {
    this.spinnerService.loadingMessage = "Updating Programming Question Test Case"
    this.questionService.updateTestCase(this.testCaseForm.value, this.selectedQuestion.id).subscribe((response: any) => {
      this.showTestCaseForm = false
      this.getAllTestCases()
      alert(response)
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error.error)
        return
      }
      alert("Check connection")
    })
  }

  // Delete programming question test case.
  deleteTestCase(testCaseID: string): void {
    if (confirm("Are you sure you want to delete the programming question test case?")) {
      this.spinnerService.loadingMessage = "Deleting Programming Question Test Case"
      this.questionService.deleteTestCase(testCaseID, this.selectedQuestion.id).subscribe((response: any) => {
        this.getAllTestCases()
        alert(response)
      }, (error) => {
        console.error(error)
        if (error.error) {
          alert(error.error)
          return
        }
        alert(error.statusText)
      })
    }
  }

  // Validate programming question test case form.
  validateTestCaseForm(): void {
    if (this.testCaseForm.invalid) {
      this.testCaseForm.markAllAsTouched()
      return
    }
    if (this.isOperationTestCaseUpdate) {
      this.updateTestCase()
      return
    }
    this.addTestCase()
  }

  // =======================================OTHER FUNCTIONS========================================================

  // Compare Ob1 and Ob2.
  compareFn(optionOne: any, optionTwo: any): boolean {
    if (optionOne == null && optionTwo == null) {
      return true
    }
    return optionOne && optionTwo ? optionOne.id === optionTwo.id : optionOne === optionTwo
  }

  // Open modal with properties.
  openModal(modalContent: any, modalSize?: string): NgbModalRef {
    if (!modalSize) {
      modalSize = 'lg'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: modalSize
    }
    this.modalRef = this.modalService.open(modalContent, options)
    return this.modalRef
  }

  // =======================================GET FUNCTIONS========================================================

  // Get all components.
  getAllComponents(): void {
    this.getProgrammingQuestionTypeList()
    this.searchOrGetProgrammingQuestion()
    this.getProgrammingQuestionLevel()
    this.getProgrammingLanguageList()
  }

  // Get programming questions.
  getProgrammingQuestions(): void {

    if (this.isSearched) {
      this.getSearchedProgrammingQuestions()
      return
    }

    this.spinnerService.loadingMessage = "Getting Programming Questions"
    this.programmingQuestions = []
    this.totalProgrammingQuestion = 0
    this.questionService.getProgrammingQuestions(this.limit, this.offset).subscribe((response: any) => {
      this.programmingQuestions = response.body
      this.totalProgrammingQuestion = response.headers.get('X-Total-Count')
    }, (err: any) => {
      this.totalProgrammingQuestion = 0
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.setPaginationString()

    })
  }

  // Get programming question type list.
  getProgrammingQuestionTypeList(): void {
    this.generalService.getProgrammingQuestionTypeList().subscribe((response: any) => {
      this.programmingQuestionTypeList = response
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get programming question level list.
  getProgrammingQuestionLevel(): void {
    this.generalService.getGeneralTypeByType("programming_question_level").subscribe((response: any) => {
      this.programmingQuestionLevel = response
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get programming language list.
  getProgrammingLanguageList(): void {
    this.generalService.getProgrammingLanguageList().subscribe((response: any) => {
      this.programmingLanguageList = response
    }, (err: any) => {
      console.error(err)
    })
  }

}
