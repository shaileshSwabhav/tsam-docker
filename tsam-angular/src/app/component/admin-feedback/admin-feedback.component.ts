import { Component, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { FeedbackType, UrlConstant } from 'src/app/service/constant';
import { FeedbackGroupService, IFeedbackQuestionGroup } from 'src/app/service/feedback-group/feedbak-group.service';
import { FeedbackService, IFeedbackQuestion } from 'src/app/service/feedback/feedback.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-feedback',
  templateUrl: './admin-feedback.component.html',
  styleUrls: ['./admin-feedback.component.css']
})
export class AdminFeedbackComponent implements OnInit {

  //modal
  modalHeader: string;
  modalButton: string;
  modelAction: () => void
  formHandler: (param?) => void;
  modalRef: any;

  // feedback form
  feedbackQuestionForm: FormGroup
  feedbackQuestionSearchForm: FormGroup

  // feedback questions
  feedbackType: any[]
  feedbackQuestion: IFeedbackQuestion
  feedbackQuestionsList: IFeedbackQuestion[]
  totalFeedbackQuestions: number
  feedbackQuestionGroupList: IFeedbackQuestionGroup[]
  searchFeedbackQuestionGroupList: IFeedbackQuestionGroup[]

  //pagination
  limit: number;
  currentPage: number;
  offset: number;
  paginationString: string

  // spinner


  // access
  permission: IPermission

  isViewClicked: boolean
  isUpdateClicked: boolean

  // search
  searchFormValue: any
  showSearch: boolean
  isSearched: boolean

  isGroupRequired: boolean
  isOptionsAdded: boolean
  isSearchGroupLoading: boolean

  private MAX_RATING = 10

  @ViewChild("drawer") drawer: any
  @ViewChild("feedbackQuestionFormModal") feedbackQuestionFormModal: any
  @ViewChild("deleteFeedbackQuestionModal") deleteFeedbackQuestionModal: any

  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private generalService: GeneralService,
    private urlConstant: UrlConstant,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private route: ActivatedRoute,
    private router: Router,
    private groupService: FeedbackGroupService,
    private feedbackService: FeedbackService,
    private feedbackTypeConstant: FeedbackType,
  ) {
    this.initialVariables()
    this.createForms()
    this.getAllComponents()
  }

  initialVariables(): void {
    this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_FEEDBACK)

    this.feedbackQuestionsList = []
    this.feedbackType = []
    this.feedbackQuestionGroupList = []

    this.limit = 5
    this.offset = 0
    this.totalFeedbackQuestions = 0

    this.isViewClicked = false
    this.showSearch = false
    this.isSearched = false
    this.isGroupRequired = false
    this.isOptionsAdded = false
    this.isUpdateClicked = false

    this.spinnerService.loadingMessage = "Getting all feedback questions"
  }

  getAllComponents(): void {
    this.getFeedbackType()
    this.getQueryParams()
  }

  createForms(): void {
    this.createFeedbackQuestionSearchForm()
    this.createFeedbackQuestionForm()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  createFeedbackQuestionSearchForm(): void {
    this.feedbackQuestionSearchForm = this.formBuilder.group({
      keyword: new FormControl(null),
      questionType: new FormControl(null),
      isActive: new FormControl('1'),
      groupID: new FormControl(null),
    })
  }

  createFeedbackQuestionForm(): void {
    this.feedbackQuestionForm = this.formBuilder.group({
      id: new FormControl(null),
      type: new FormControl(null, [Validators.required]),
      question: new FormControl(null, [Validators.required, Validators.maxLength(250)]),
      order: new FormControl(null, [Validators.required, Validators.min(1)]),
      hasOptions: new FormControl(false, [Validators.required]),
      keyword: new FormControl(null),
      feedbackQuestionGroup: new FormControl(null),
      isActive: new FormControl(true, [Validators.required]),
      options: new FormArray([]),
    })
  }

  get feedbackOptions() {
    return this.feedbackQuestionForm.get('options') as FormArray
  }

  addFeedbackOptions(): void {
    this.feedbackOptions.push(this.formBuilder.group({
      id: new FormControl(null),
      order: new FormControl(null, [Validators.required, Validators.min(1)]),
      key: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(10)]),
      value: new FormControl(null, [Validators.required, Validators.maxLength(250)]),
    }))
  }

  onAddRatingTypeOptionClick(event: any): void {

    if (event.target.checked) {
      // Remove previously added options
      if (this.feedbackOptions.controls.length > 0) {
        if (!confirm("Selecting this option will delete all the previouusly added options")) {
          event.target.checked = false
          return
        }
      }
      this.feedbackOptions.clear()
      this.setKeywordField()
      this.createMaxRatingFeedbackOptions()
      // this.isOptionsAdded = true
      event.target.checked = true
      return
    }

    // Remove previously added options
    if (this.feedbackOptions.controls.length > 0) {
      if (!confirm("Are you sure you want to delete all options?")) {
        event.target.checked = true
        return
      }
      this.feedbackOptions.clear()
      // this.isOptionsAdded = false
      event.target.checked = false
    }
  }

  createMaxRatingFeedbackOptions(): void {

    for (let index = 0; index < this.MAX_RATING; index++) {
      // console.log(index);
      this.addFeedbackOptions()
      this.feedbackOptions.at(index).get('key').setValue(index + 1)
      this.feedbackOptions.at(index).get('order').setValue(index + 1)
      this.feedbackOptions.at(index).get('value').setValue((index + 1).toString())
    }
    this.feedbackQuestionForm.markAsDirty()
  }



  onAddFeedbackQuestionClick(): void {
    this.isViewClicked = false
    this.isUpdateClicked = false
    this.isOptionsAdded = false
    this.updateVariable(this.createFeedbackQuestionForm, "Add Feedback Question", "Add", this.addFeedbackQuestion)
    this.formHandler()
    this.openModal(this.feedbackQuestionFormModal, "xl")
  }

  onViewFeedbackQuestionClick(feedbackQuestion: IFeedbackQuestion): void {
    this.feedbackQuestion = feedbackQuestion
    this.isUpdateClicked = false
    // console.log(this.feedbackQuestion);
    this.isViewClicked = true
    if (this.feedbackQuestion.options?.length == this.MAX_RATING) {
      this.isOptionsAdded = true
    }
    this.updateVariable(this.createFeedbackQuestionForm, "Feedback Question Details", "Add", this.addFeedbackQuestion)
    this.formHandler()
    this.updateForm()
    this.feedbackQuestionForm.disable()
    this.openModal(this.feedbackQuestionFormModal, "xl")
  }

  onUpdateFeedbackQuestionClick(): void {
    this.getFeedbackQuestionGroupByType(this.feedbackQuestionForm.get('type').value)
    this.isViewClicked = false
    this.isUpdateClicked = true
    this.feedbackQuestionForm.enable()
    this.modalHeader = "Update Feedback Question"
    this.modalButton = "Update"
    this.modelAction = this.updateFeedbackQuestion
  }

  updateForm(): void {

    if (!this.feedbackQuestion.hasOptions) {
      this.feedbackQuestion.options = []
    } else {
      for (let index = 0; index < this.feedbackQuestion.options.length; index++) {
        this.addFeedbackOptions()
      }
    }
    this.feedbackQuestionForm.patchValue(this.feedbackQuestion)
  }

  onDeleteFeedbackQuestionClick(feedbackQuestion: IFeedbackQuestion): void {
    this.feedbackQuestion = feedbackQuestion
    this.openModal(this.deleteFeedbackQuestionModal, 'md')
  }

  updateVariable(formaction: any, modelheader: string, modelbutton: string, modelaction: any): void {
    this.formHandler = formaction;
    this.modalHeader = modelheader;
    this.modalButton = modelbutton;
    this.modelAction = modelaction;
  }

  // Compare two Object.
  compareFn(ob1: any, ob2: any): any {
    if (ob1 == null && ob2 == null) {
      return true
    }
    return ob1 && ob2 ? ob1.id == ob2.id : ob1 == ob2;
  }

  //page change function
  changePage(pageNumber: number): void {

    this.currentPage = pageNumber;
    this.offset = this.currentPage - 1;
    this.getAllFeedbackQuestions();
  }

  resetSearchAndGetAll(): void {
    this.isSearched = false
    this.showSearch = false
    this.feedbackQuestionSearchForm.reset()
    this.createQueryParams()
    this.changePage(1)
  }

  addFeedbackOptionsToForm(hasOptions: boolean) {
    this.setKeywordField()

    // remove all obejcts from feedbackOptions formarray
    this.feedbackOptions.clear()
  }

  setKeywordField(): void {
    if (this.feedbackQuestionForm.controls["hasOptions"]?.value) {
      this.feedbackQuestionForm.controls["keyword"].setValidators([Validators.required, Validators.maxLength(30)])
    } else {
      this.feedbackQuestionForm.controls["keyword"].setValidators(null)
    }
    this.feedbackQuestionForm.controls["keyword"].updateValueAndValidity()
  }

  deleteFeedbackOption(index: number) {
    this.feedbackOptions.removeAt(index)
    this.feedbackQuestionForm.markAsDirty()
  }

  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalFeedbackQuestions < end) {
      end = this.totalFeedbackQuestions
    }
    if (this.totalFeedbackQuestions == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end}`
  }

  openModal(contentModal: any, modalSize?: string): void {
    if (modalSize == undefined || modalSize == "") {
      modalSize = 'xl'
    }
    this.modalRef = this.modalService.open(contentModal, {
      ariaLabelledBy: 'modal-basic-title',
      backdrop: 'static',
      size: modalSize
    })
  }

  assignNecessaryValidators(): void {

    if (this.feedbackQuestionForm.get('type').value != this.feedbackTypeConstant.AHA_MOMENT_FEEDBACK &&
      this.feedbackQuestionForm.get('type').value != this.feedbackTypeConstant.FACULTY_SESSION_FEEDBACK &&
      this.feedbackQuestionForm.get('type').value != this.feedbackTypeConstant.FACULTY_BATCH_FEEDBACK) {
      this.feedbackQuestionForm.get('feedbackQuestionGroup').clearValidators()
      this.isGroupRequired = false
    } else {
      this.feedbackQuestionForm.get('feedbackQuestionGroup').setValidators([Validators.required])
      this.isGroupRequired = true
    }

    this.utilService.updateValueAndValiditors(this.feedbackQuestionForm)
  }

  validateFeedbackQuestion(): void {

    console.log(this.feedbackQuestionForm.controls)

    if (this.feedbackQuestionForm.invalid) {
      this.feedbackQuestionForm.markAllAsTouched()
      return
    }
    this.modelAction()
  }

  searchAndCloseDrawer(): void {
    this.drawer.toggle()
    this.searchFeedbackQuestion()
  }

  searchFeedbackQuestion(): void {
    this.spinnerService.loadingMessage = "Searching feedback questions";
    this.searchFormValue = { ...this.feedbackQuestionSearchForm.value }
    let flag: boolean = true

    for (let field in this.searchFormValue) {
      if (!this.searchFormValue[field]) {
        delete this.searchFormValue[field];
      } else {
        this.isSearched = true
        flag = false
      }
    }
    // No API call on empty search.
    if (flag) {
      return
    }
    this.createQueryParams()
    this.changePage(1)
  }

  createQueryParams() {
    if (this.isSearched) {
      this.router.navigate([], {
        queryParams: this.searchFormValue,
        // queryParamsHandling: 'merge',
      })
      return
    }
    this.router.navigate([this.urlConstant.TRAINING_FEEDBACK])
  }

  getQueryParams(): void {
    this.route.queryParams.subscribe((param) => {
      // console.log(param)
      this.searchFormValue = { ...param }
      for (let field in this.searchFormValue) {
        if (!this.searchFormValue[field]) {
          delete this.searchFormValue[field];
        } else {
          this.isSearched = true
          this.showSearch = true
        }
      }

      if (this.searchFormValue.questionType) {
        this.getSearchFeedbackQuestionGroupByType(this.searchFormValue.questionType)
      }
      this.feedbackQuestionSearchForm.patchValue(this.searchFormValue)
      this.changePage(1)
    })
  }

  // =============================================================CRUD=============================================================

  getAllFeedbackQuestions(): void {
    if (this.isSearched) {
      this.getSearchedFeedbackQuestion()
      return
    }
    this.spinnerService.loadingMessage = "Getting all feedback questions"

    this.feedbackService.getFeedbackQuestions(this.limit, this.offset).subscribe((response: any) => {
      // console.log(response.body);
      this.feedbackQuestionsList = response.body
      this.totalFeedbackQuestions = response.headers.get('X-Total-Count')
      this.setPaginationString()

    }, (err: any) => {
      this.totalFeedbackQuestions = 0
      this.setPaginationString()
      console.error(err)

      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getSearchedFeedbackQuestion(): void {
    this.spinnerService.loadingMessage = "Searching feedback questions"

    console.log(this.searchFormValue)
    this.feedbackService.getFeedbackQuestions(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      // console.log(response.body);
      this.feedbackQuestionsList = response.body
      this.totalFeedbackQuestions = response.headers.get('X-Total-Count')
      this.setPaginationString()

    }, (err: any) => {
      console.error(err)
      this.totalFeedbackQuestions = 0
      this.setPaginationString()

      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  addFeedbackQuestion(): void {
    // console.log(this.feedbackQuestionForm.value);
    this.spinnerService.loadingMessage = "Adding feedback question"

    this.feedbackService.addFeedbackQuestion(this.feedbackQuestionForm.value).subscribe((response: any) => {
      // console.log(response);
      alert("Feedback question successfully added")
      this.modalRef.close()
      this.changePage(1)

    }, (err: any) => {
      console.error(err)

      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  updateFeedbackQuestion(): void {
    // console.log(this.feedbackQuestionForm.value);
    this.spinnerService.loadingMessage = "Updating feedback question"

    this.feedbackService.updateFeedbackQuestion(this.feedbackQuestionForm.value).subscribe((response: any) => {
      // console.log(response);
      alert("Feedback question successfully updated")
      this.modalRef.close()
      this.changePage(this.currentPage)

    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)

    })
  }

  updateFeedbackQuestionStatus(feedbackQuestion: IFeedbackQuestion): void {
    if (!confirm("Are you sure you want to update status of the feedback question?")) {
      return
    }
    feedbackQuestion.isActive = !feedbackQuestion.isActive
    console.log(feedbackQuestion);

    this.spinnerService.loadingMessage = "Updating feedback question status"

    this.feedbackService.updateFeedbackQuestionStatus(feedbackQuestion).subscribe((response: any) => {
      console.log(response)
      alert("Feedback question successfully updated")

      this.changePage(this.currentPage)
    }, (err: any) => {
      console.error(err)
      feedbackQuestion.isActive = !feedbackQuestion.isActive
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)

    })
  }

  deleteFeedbackQuestion(): void {
    this.spinnerService.loadingMessage = "Deleting feedback question"

    this.feedbackService.deleteFeedbackQuestion(this.feedbackQuestion.id).subscribe((response: any) => {
      // console.log(response);
      alert("Feedback question successfully deleted")
      this.modalRef.close()
      // this.getAllFeedbackQuestions()
      this.changePage(this.currentPage)

    }, (err: any) => {
      console.error(err)

      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getFeedbackType(): void {
    this.generalService.getGeneralTypeByType("feedback_type").subscribe((response: any) => {
      this.feedbackType = response
    }, (err: any) => {
      console.error(err)
    })
  }

  getFeedbackQuestionGroupByType(type: string): void {
    this.feedbackQuestionGroupList = []
    // this.assignNecessaryValidators()
    // this.spinnerService.loadingMessage = "Getting data..."
    // 
    this.groupService.getFeedbackQuestionGroupList(type).subscribe((response: any) => {
      this.feedbackQuestionGroupList = response.body
      // 
    }, (err: any) => {
      console.error(err)
      // 
    })
  }

  getSearchFeedbackQuestionGroupByType(type: string): void {
    this.searchFeedbackQuestionGroupList = []
    // this.spinnerService.loadingMessage = "Getting data..."
    // 
    this.isSearchGroupLoading = true
    this.groupService.getFeedbackQuestionGroupList(type).subscribe((response: any) => {
      this.searchFeedbackQuestionGroupList = response.body
      // 
    }, (err: any) => {
      console.error(err)
      // 
    }).add(() => {
      this.isSearchGroupLoading = false
    })
  }

  resetSearchForm(): void {
    this.searchFeedbackQuestionGroupList = []
    this.feedbackQuestionSearchForm.reset()
  }

}
