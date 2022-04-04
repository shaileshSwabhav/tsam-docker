import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalRef, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { GeneralService } from 'src/app/service/general/general.service';
import { FeedbackGroupService, IFeedbackQuestionGroup } from 'src/app/service/feedback-group/feedbak-group.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UrlConstant } from 'src/app/service/constant';

@Component({
  selector: 'app-admin-group',
  templateUrl: './admin-group.component.html',
  styleUrls: ['./admin-group.component.css']
})
export class AdminGroupComponent implements OnInit {

  // Flags.
  isSearched: boolean
  isOperationUpdate: boolean
  isViewMode: boolean
  showSearch: boolean

  // Feedback Question Group.
  feedbackQuestionGroupList: IFeedbackQuestionGroup[]
  feedbackQuestionGroupForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  totalFeedbackQuestionGroups: number
  offset: number
  paginationString: string

  // Modal.
  modalRef: any
  @ViewChild('feedbackQuestionGroupFormModal') feedbackQuestionGroupFormModal: any
  @ViewChild('deleteFeedbackQuestionGroupModal') deleteFeedbackQuestionGroupModal: any
  @ViewChild("drawer") drawer: any

  // Spinner.



  // Search.
  feedbackQuestionGroupSearchForm: FormGroup
  searchFormValue: any

  // component
  feedbackType: any[]

  // permission
  permission: IPermission

  constructor(
    private formBuilder: FormBuilder,
    private groupService: FeedbackGroupService,
    private utilService: UtilityService,
    private generalService: GeneralService,
    private localService: LocalService,
    private urlConstant: UrlConstant,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.TRAINING_GROUP)

    // Components.
    this.feedbackQuestionGroupList = [] as IFeedbackQuestionGroup[]
    this.feedbackType = []

    // Flags.
    this.isOperationUpdate = false
    this.isViewMode = false
    this.isSearched = false
    this.showSearch = false

    // Pagination.
    this.limit = 5
    this.offset = 0
    this.currentPage = 0

    // Initialize forms
    this.createFeedbackQuestionGroupForm()
    this.createFeedbackQuestionGroupSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Feedback Question Groups"

  }

  getAllComponents() {
    this.getFeedbackType()
    this.searchOrGetFeedbackQuestionGroups()
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create feedbackQuestionGroup form.
  createFeedbackQuestionGroupForm(): void {
    this.feedbackQuestionGroupForm = this.formBuilder.group({
      id: new FormControl(null),
      groupName: new FormControl(null, [Validators.required, Validators.maxLength(50)]),
      groupDescription: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
      order: new FormControl(null, [Validators.required, Validators.min(1)]),
      type: new FormControl(null, [Validators.required]),
    })
  }

  // Create feedbackQuestionGroup search form.
  createFeedbackQuestionGroupSearchForm(): void {
    this.feedbackQuestionGroupSearchForm = this.formBuilder.group({
      groupName: new FormControl(null),
      groupDescription: new FormControl(null),
      order: new FormControl(null),
      type: new FormControl(null),
    })
  }
  // =============================================================FEEDBACK QUESTION GROUP CRUD FUNCTIONS==========================================================================
  // On clicking add new feedbackQuestionGroup button.
  onAddNewFeedbackQuestionGroupClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createFeedbackQuestionGroupForm()
    this.openModal(this.feedbackQuestionGroupFormModal, 'xl')
  }

  // Add new feedbackQuestionGroup.
  addFeedbackQuestionGroup(): void {
    this.spinnerService.loadingMessage = "Adding Feedback Question Group"


    this.groupService.addFeedbackQuestionGroup(this.feedbackQuestionGroupForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllFeedbackQuestionGroups()
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

  // On clicking view feedbackQuestionGroup button.
  onViewFeedbackQuestionGroupClick(feedbackQuestionGroup: IFeedbackQuestionGroup): void {
    this.isViewMode = true
    this.createFeedbackQuestionGroupForm()
    this.feedbackQuestionGroupForm.patchValue(feedbackQuestionGroup)
    this.feedbackQuestionGroupForm.disable()
    this.openModal(this.feedbackQuestionGroupFormModal, 'xl')
  }

  // On cliking update form button in feedbackQuestionGroup form.
  onUpdateFeedbackQuestionGroupClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.feedbackQuestionGroupForm.enable()
  }

  // Update feedbackQuestionGroup.
  updateFeedbackQuestionGroup(): void {
    this.spinnerService.loadingMessage = "Updating Feedback Question Group"


    this.groupService.updateFeedbackQuestionGroup(this.feedbackQuestionGroupForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllFeedbackQuestionGroups()
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

  // On clicking delete feedbackQuestionGroup button. 
  onDeleteFeedbackQuestionGroupClick(feedbackQuestionGroupID: string): void {
    this.openModal(this.deleteFeedbackQuestionGroupModal, 'md').result.then(() => {
      this.deleteFeedbackQuestionGroup(feedbackQuestionGroupID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete feedbackQuestionGroup after confirmation from user.
  deleteFeedbackQuestionGroup(feedbackQuestionGroupID: string): void {
    this.spinnerService.loadingMessage = "Deleting Feedback Question Group"


    this.groupService.deleteFeedbackQuestionGroup(feedbackQuestionGroupID).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllFeedbackQuestionGroups()
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

  // =============================================================FEEDBACK QUESTION GROUP SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.feedbackQuestionGroupSearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.showSearch = false
    this.router.navigate([this.urlConstant.TRAINING_GROUP])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.feedbackQuestionGroupSearchForm.reset()
  }

  searchAndCloseDrawer(): void {
    this.drawer.toggle()
    this.searchFeedbackQuestionGroups()
  }

  // Search feedback question groups.
  searchFeedbackQuestionGroups(): void {
    this.searchFormValue = { ...this.feedbackQuestionGroupSearchForm?.value }
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
    if (!this.isSearched) {
      return
    }
    this.spinnerService.loadingMessage = "Searching Feedback Question Groups"
    this.changePage(1)
  }

  // ================================================OTHER FUNCTIONS FOR FEEDBACK QUESTION GROUP===============================================
  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllFeedbackQuestionGroups()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetFeedbackQuestionGroups() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllFeedbackQuestionGroups()
      return
    }
    this.feedbackQuestionGroupSearchForm.patchValue(queryParams)
    this.searchFeedbackQuestionGroups()
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

  // On clicking sumbit button in feedbackQuestionGroup form.
  onFormSubmit(): void {
    if (this.feedbackQuestionGroupForm.invalid) {
      this.feedbackQuestionGroupForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateFeedbackQuestionGroup()
      return
    }
    this.addFeedbackQuestionGroup()
  }

  // Set total feedback question group list on current page.
  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalFeedbackQuestionGroups < end) {
      end = this.totalFeedbackQuestionGroups
    }
    if (this.totalFeedbackQuestionGroups == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start}-${end}`
  }

  // Compare for select option field.
  // compareFn(optionOne: any, optionTwo: any): boolean {
  //   if (optionOne == null && optionTwo == null) {
  //     return true
  //   }
  //   if (optionTwo != undefined && optionOne != undefined) {
  //     return optionOne.value === optionTwo.value
  //   }
  //   return false
  // }

  // =============================================================GET FUNCTIONS==========================================================================
  // Get all feedback question groups.
  getAllFeedbackQuestionGroups() {
    this.spinnerService.loadingMessage = "Getting All Feedback Question Groups"


    this.groupService.getAllFeedbackQuestionGroups(this.limit, this.offset, this.searchFormValue).subscribe((response) => {
      this.totalFeedbackQuestionGroups = response.headers.get('X-Total-Count')
      this.feedbackQuestionGroupList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

  getFeedbackType(): void {
    this.generalService.getGeneralTypeByType("feedback_type").subscribe((response: any) => {
      this.feedbackType = response
    }, (err: any) => {
      console.error(err)
    })
  }

}
