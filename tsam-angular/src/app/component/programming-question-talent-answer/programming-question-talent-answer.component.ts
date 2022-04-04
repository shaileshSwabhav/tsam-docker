import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant, Role } from 'src/app/service/constant';
import { ISearchFilterField } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IProgrammingQuestionTalentAnswerDTO, ProgrammingQuestionTalentAnswerService } from 'src/app/service/programming-question-talent-answer/programming-question-talent-answer.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { ProgrammingQuestionModalComponent } from '../programming-question-modal/programming-question-modal.component';

const PROBLEM_OF_THE_DAY = "Problem of the day"
const PRACTICE = "Practice"

@Component({
  selector: 'app-programming-question-talent-answer',
  templateUrl: './programming-question-talent-answer.component.html',
  styleUrls: ['./programming-question-talent-answer.component.css']
})
export class ProgrammingQuestionTalentAnswerComponent implements OnInit {

  // Programming question talent answer.
  answerList: IProgrammingQuestionTalentAnswerDTO[]
  multipleAnswerList: IProgrammingQuestionTalentAnswerDTO[]
  selectedAnswer: IProgrammingQuestionTalentAnswerDTO

  // Flags.
  isOperationUpdate: boolean
  isViewMode: boolean

  // Pagination.
  limit: number
  currentPage: number
  offset: number
  paginationStart: number
  paginationEnd: number
  totalAnswers: number

  // Modal.
  modalRef: any
  @ViewChild('deleteAnswerModal') deleteAnswerModal: any
  @ViewChild('viewAllAnswersModal') viewAllAnswersModal: any

  // Spinner.



  // Search.
  isSearched: boolean
  answerSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  // Permission.
  permission: IPermission
  roleName: string

  // Programming type.
  currentProgrammingType: string


  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private router: Router,
    private route: ActivatedRoute,
    private urlConstant: UrlConstant,
    private localService: LocalService,
    private role: Role,
    private answerService: ProgrammingQuestionTalentAnswerService,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables(): void {


    // Programming question talent answer.
    this.answerList = []
    this.multipleAnswerList = []

    // Flags.
    this.isOperationUpdate = false
    this.isViewMode = false

    // Search.
    this.isSearched = false
    this.searchFilterFieldList = []

    // Pagination.
    this.limit = 5
    this.offset = 0
    this.currentPage = 0

    // Initialize forms
    this.createAnswerSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Programming Question Talent Answers"


    // Programming type.
    this.currentProgrammingType = PROBLEM_OF_THE_DAY

    // Permision.
    // Get permissions from menus using utilityService function.
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.ADMIN_PROGRAMMING_QUESTION_TALENT_ANSWER)
    }
    else {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.PROGRAMMING_QUESTION_TALENT_ANSWER)
    }
    this.roleName = this.localService.getJsonValue("roleName")
  }

  // questionModal(): void{

  //   // console.log(this._programmingQuestionComponent.programmingModal)
  //   // // this._programmingQuestionComponent.onAddNewQuestionClick()
  //   // this.openModal(this._programmingQuestionComponent.programmingModal, 'md').result.then(() => {
  //   // }, (err) => {
  //   //   console.error(err)
  //   //   return
  //   // })
  //   this.modalService.open(ProgrammingQuestionModalComponent, { ariaLabelledBy: 'modal-basic-title', backdrop: 'static', size: 'xl' }).
  //   result.then(() => {
  //     alert("success")
  //   }).catch((err) => {
  //     console.log("Dismiss", err)
  //     return
  //   })
  // }

  // =============================================================CREATE FORMS==========================================================================

  // Create programming question talent answer search form.
  createAnswerSearchForm(): void {
    this.answerSearchForm = this.formBuilder.group({
      label: new FormControl(null),
    })
  }

  // =============================================================PROGRAMMING QUESTION TALENT ANSWER CRUD FUNCTIONS==========================================================================

  // On clicking delete programming question talent answer button. 
  onDeleteAnswerClick(answerID: string): void {
    this.openModal(this.deleteAnswerModal, 'md').result.then(() => {
      this.deleteAnswer(answerID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete programming question talent answer after confirmation from user.
  deleteAnswer(answerID: string): void {
    this.spinnerService.loadingMessage = "Deleting Programming Question Talent Answer"


    this.answerService.deleteAnswer(answerID).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllComponents()
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

  // =============================================================PROGRAMMING QUESTION TALENT ANSWER SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.answerSearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    if (this.roleName == this.role.ADMIN) {
      this.router.navigate(['/admin/coding-problems/answers'])
      return
    }
    this.router.navigate(['/coding-problems/answers'])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.answerSearchForm.reset()
  }

  // Search programming question talent answer.
  searchAnswers(): void {
    this.searchFormValue = { ...this.answerSearchForm?.value }
    if (this.searchFormValue.departmentID) {
      this.searchFormValue.departmentID = this.searchFormValue.departmentID.id
    }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field]
      } else {
        this.isSearched = true
      }
    }
    this.searchFilterFieldList = []
    for (var property in this.searchFormValue) {
      let text: string = property
      let result: string = text.replace(/([A-Z])/g, " $1");
      let finalResult: string = result.charAt(0).toUpperCase() + result.slice(1);
      let valueArray: any[] = []
      if (Array.isArray(this.searchFormValue[property])) {
        valueArray = this.searchFormValue[property]
      }
      else {
        valueArray.push(this.searchFormValue[property])
      }
      this.searchFilterFieldList.push(
        {
          propertyName: property,
          propertyNameText: finalResult,
          valueList: valueArray
        })
    }
    if (this.searchFilterFieldList.length == 0) {
      this.resetSearchAndGetAll()
    }
    if (!this.isSearched) {
      return
    }
    this.spinnerService.loadingMessage = "Searching Programming Question Talent Answers"
    this.changePage(1)
  }

  // Delete search criteria from programming question talent answer search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.answerSearchForm.get(searchName).setValue(null)
    this.searchAnswers()
  }

  // ================================================ OTHER FUNCTIONS FOR PROGRAMMING QUESTION TALENT ANSWERS ===============================================
  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllAnswers()
  }

  // Checks the url's query params and decides whether to call get or search.
  searchOrGetAnswers(): void {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllAnswers()
      return
    }
    this.answerSearchForm.patchValue(queryParams)
    this.getAllAnswers()
  }

  // Set total programming question talent answer list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalAnswers < this.paginationEnd) {
      this.paginationEnd = this.totalAnswers
    }
  }

  // Format fields of programming question talent answers.
  formatAnswerListFields(): void {

    for (let i = 0; i < this.answerList.length; i++) {

      // Set class for difficulty.
      if (this.answerList[i].programmingQuestion.level == 1) {
        this.answerList[i].programmingQuestion.levelName = "Easy"
        this.answerList[i].programmingQuestion.levelClass = "easy"
      }
      if (this.answerList[i].programmingQuestion.level == 2) {
        this.answerList[i].programmingQuestion.levelName = "Medium"
        this.answerList[i].programmingQuestion.levelClass = "card-sub-detail medium"
      }
      if (this.answerList[i].programmingQuestion.level == 3) {
        this.answerList[i].programmingQuestion.levelName = "Hard"
        this.answerList[i].programmingQuestion.levelClass = "card-sub-detail hard"
      }
    }
  }

  // Redirect to programming question talent answer details page.
  redirectToAnswerDetails(answerID: string): void {
    if (this.roleName == this.role.ADMIN) {
      this.router.navigate(['/admin/coding-problems/answers/answer-details'], {
        queryParams: {
          "answerID": answerID,
        }
      }).catch(err => {
        console.error(err)
      })
    } else {
      this.router.navigate(['/coding-problems/answers/answer-details'], {
        queryParams: {
          "answerID": answerID,
        }
      }).catch(err => {
        console.error(err)
      })
    }
  }

  // On changing tab get the component lists.
  onTabChange(event: any) {
    this.searchFormValue = {}
    this.searchFilterFieldList = []
    if (event == 1) {
      this.currentProgrammingType = PROBLEM_OF_THE_DAY
    }
    if (event == 2) {
      this.currentProgrammingType = PRACTICE
    }
    this.searchFormValue.programmingType = this.currentProgrammingType
    this.getAllAnswers()
  }

  // On clicking view all answers button.
  onViewAllAnswersButtonClick(answer: IProgrammingQuestionTalentAnswerDTO): void {
    this.selectedAnswer = answer
    this.openModal(this.viewAllAnswersModal, 'xl')
    this.getAllMultipleAnswers()
  }

  // Used to open modal.
  openModal(content: any, size?: string): NgbModalRef {
    if (!size) {
      size = 'lg'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: size
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get all components.
  getAllComponents(): void {
    this.searchFormValue = {}
    this.searchFormValue.programmingType = this.currentProgrammingType
    this.searchOrGetAnswers()
  }

  // Get all programming question talent answers.
  getAllAnswers(): void {
    this.spinnerService.loadingMessage = "Getting All Programming Question Talent Answer"


    this.answerService.getAllAnswers(this.limit, this.offset, this.searchFormValue).subscribe((response) => {
      this.totalAnswers = response.headers.get('X-Total-Count')
      this.answerList = response.body
      this.formatAnswerListFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

  // Get all multiple programming question talent answers.
  getAllMultipleAnswers(): void {
    this.spinnerService.loadingMessage = "Getting All Programming Question Talent Answer"


    let queryParams: any = {
      programmingType: this.currentProgrammingType,
      programmingQuestionID: this.selectedAnswer.programmingQuestion.id,
      programmingLanguageID: this.selectedAnswer.programmingLanguage.id,
      talentID: this.selectedAnswer.talent.id
    }
    this.answerService.getAllAnswers(-1, 0, queryParams).subscribe((response) => {
      this.multipleAnswerList = response.body
      this.formatAnswerListFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

}
