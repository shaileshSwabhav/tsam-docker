import { DatePipe } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { concatMap } from 'rxjs/operators';
import { BatchTalentService, IBatchTalent } from 'src/app/service/batch-talent/batch-talent.service';
import { BatchService, ITalentBatchDTO } from 'src/app/service/batch/batch.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { GeneralService, IFeedbackOptions, IFeedbackQuestion } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ITalent, TalentService } from 'src/app/service/talent/talent.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-batch-talent',
  templateUrl: './batch-talent.component.html',
  styleUrls: ['./batch-talent.component.css']
})
export class BatchTalentComponent implements OnInit {

  // Batch talent.
  batchTalentList: any[]
  totalBatchTalents: number

  // Batch.
  @Input() batchID: string

  // Form.
  searchedTalentForm: FormGroup

  // Spinner.


  @Output() changeLoadingMessage: EventEmitter<string>

  // Flags.
  isTalentSearched: boolean

  // Access.
  permission: IPermission
  loginID: string

  // Pagination.
  totalSelectedTalents: number
  selectedTalentLimit: number
  selectedTalentOffset: number
  selectedTalentPaginationStart: number
  selectedTalentPaginationEnd: number
  currentSelectedTalentPage: number

  totalSearchedTalents: number
  searchedTalentLimit: number
  searchedTalentOffset: number
  searchedTalentPaginationStart: number
  searchedTalentPaginationEnd: number
  currentSearchedTalentPage: number

  // Talent.
  searchedTalentList: any[]
  selectedTalentList: any[]
  addSelectedTalents: any[]

  // Modal.
  modalRef: any
  @ViewChild('talentModal') talentModal: any

  constructor(
    private formBuilder: FormBuilder,
    private batchTalentService: BatchTalentService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private route: ActivatedRoute,
    private utilService: UtilityService,
    private batchService: BatchService,
    private datePipe: DatePipe,
    private urlConstant: UrlConstant,
    private localService: LocalService,
    private talentService: TalentService,
		private role: Role,
  ) {
    this.initializeVariables()
    // this.extractID()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    // Components.
    this.batchTalentList = []

    // Talent
    this.searchedTalentList = []

    // Spinner.
    this.spinnerService.loadingMessage = "Getting student list"
    this.changeLoadingMessage = new EventEmitter()


    // Access.
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN || this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER_SESSION_DETAILS)
		}
		if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
			this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH_SESSION_DETAILS)
		}
    this.loginID = this.localService.getJsonValue("loginID")

    // Flags.
    this.isTalentSearched = false

    // Pagination.
    this.searchedTalentOffset = 0
    this.searchedTalentLimit = 5
    this.selectedTalentOffset = 0
    this.selectedTalentLimit = 5

    // Batch id.
    this.batchID = this.route.snapshot.queryParamMap.get("batchID")
    if (this.batchID) {
      this.getBatchTalentDetails()
    }
  }

  // =============================================================CREATE FORMS==========================================================================

  // Create search talent form.
  createSearchTalentForm() {
    this.searchedTalentForm = this.formBuilder.group({
      firstName: new FormControl(null),
      lastName: new FormControl(null),
      email: new FormControl(null),
    })
  }

  // ============================================================= BATCH TALENT CRUD FUNCTIONS ==========================================================================

  // Update suspension date of one batch talent.
  updateSuspensionDateBatchTalent(batchTalentID: string, questionString: string, isSuspend: boolean): void {
    if (confirm("Are you sure you want to " + questionString + " the batch talent?")) {
      this.spinnerService.loadingMessage = "Updating batch talent"
      this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

      let batchTalent: any = {
        id: batchTalentID,
        isSuspend: isSuspend
      }
      this.batchTalentService.updateSuspensionDateBatchTalent(batchTalent).subscribe((response: any) => {
        this.getBatchTalentDetails()
      }, (err: any) => {
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(err.error.error)
        console.error(err);
      })
    }
  }

  // Update is active of batch talent (delete or restore batch talent).
  updateIsActiveBatchTalent(tempBatchTalent: any, questionString: string, isActive: boolean): void {
    if (confirm("Are you sure you want to " + questionString + " the batch talent?")) {
      this.spinnerService.loadingMessage = "Updating batch talent"
      this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

      let batchTalent: any = {
        id: tempBatchTalent.id,
        batchID: tempBatchTalent.batch.id,
        isActive: isActive,
      }
      this.batchTalentService.updateIsActiveBatchTalent(batchTalent).subscribe((response: any) => {
        this.getBatchTalentDetails()
      }, (err: any) => {
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(err.error.error)
        console.error(err);
      })
    }
  }

  // ============================================================= TALENT ASSIGN FUNCTIONS ==========================================================================

  // On clicking add talents to batch.
  assignTalentsToBatch() {
    this.addSelectedTalents = []
    this.createSearchTalentForm()
    this.openModal(this.talentModal, "xl")
    this.getSelectedTalentsForBatch()
    this.getAllTalents()
  }

  // Reset talent search form.
  resetTaleSearchForm(): void {
    this.searchedTalentForm.reset()
  }

  // Reset talent search form and get all talents.
  resetTalentSearchAndGetAll(): void {
    this.searchedTalentForm.reset()
    this.isTalentSearched = false
    this.changeSearchedTalentsPage(1)
  }

  // Search talents.
  searchTalents(): void {
    this.isTalentSearched = true
    this.spinnerService.loadingMessage = "Searching Talents"
    this.changeSearchedTalentsPage(1)
  }

  // Set isChecked field of all selected talents.
  setSelectAllTalents(isSelectedAll: boolean): void {
    for (let i = 0; i < this.searchedTalentList.length; i++) {
      this.addTalentToList(isSelectedAll, this.searchedTalentList[i])
    }
  }

  // Check if all talents in selected talents are added in multiple select or not.
  checkTalentsAdded(): boolean {
    let count: number = 0

    for (let i = 0; i < this.searchedTalentList.length; i++) {
      if (this.addSelectedTalents.includes(this.searchedTalentList[i].id))
        count = count + 1
    }
    return (count == this.searchedTalentList.length)
  }

  // Check if talent is added in multiple select or not.
  checkTalentAdded(talentID): boolean {
    return this.addSelectedTalents.includes(talentID)
  }

  // Takes a list called selectedTalent and adds all the checked talents to list, also does not contain duplicate values.
  addTalentToList(isChecked: boolean, talent: any): void {
    if (isChecked) {
      if (!this.addSelectedTalents.includes(talent.id)) {
        this.addSelectedTalents.push(talent.id)
      }
      return
    }
    if (this.addSelectedTalents.includes(talent.id)) {
      let index = this.addSelectedTalents.indexOf(talent.id)
      this.addSelectedTalents.splice(index, 1)
    }
  }

  // Allocate talent to batch.
  allocateTalentsToBatch(): void {
    this.spinnerService.loadingMessage = "Talents are getting allocated to batch"
    let batchTalents = []
    for (let index = 0; index < this.addSelectedTalents.length; index++) {
      batchTalents.push({
        "talentID": this.addSelectedTalents[index],
        "batchID": this.batchID,
        "isActive": true
      })
    }
    if (batchTalents.length == 0) {
      alert("Please select talents")
      return
    }
    this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.batchTalentService.addBatchTalents(batchTalents, this.batchID).subscribe(() => {
      this.getSelectedTalentsForBatch()
      this.getBatchTalentDetails()
      alert("Talent allocated to batch successfully")
      this.addSelectedTalents = []
    }, (error) => {
      console.error(error)
      if (typeof error.error == 'object' && error) {
        alert(this.utilService.getErrorString(error))
        return
      }
      if (error.error == undefined) {
        alert('Talent could not be allocated, try again')
      }
      alert(error.statusText)
    })
  }

  // ============================================================= OTHER FUNCTIONS ==========================================================================



  // Extract id from params.
  extractID(): void {
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      if (this.batchID) {
        this.getBatchTalentDetails()
      }
    }, err => {
      console.error(err);
    })
  }

  // Set pagination of selected talents.
  setSelectedTalentPaginationString(): void {
    this.selectedTalentPaginationStart = this.selectedTalentLimit * this.selectedTalentOffset + 1
    this.selectedTalentPaginationEnd = +this.selectedTalentLimit + this.selectedTalentLimit * this.selectedTalentOffset
    if (this.totalSelectedTalents < this.selectedTalentPaginationEnd) {
      this.selectedTalentPaginationEnd = this.totalSelectedTalents
    }
  }

  // Set pagination of searched talents.
  setSearchedTalentPaginationString(): void {
    this.searchedTalentPaginationStart = this.searchedTalentLimit * this.searchedTalentOffset + 1
    this.searchedTalentPaginationEnd = +this.searchedTalentLimit + this.searchedTalentLimit * this.searchedTalentOffset
    if (this.totalSearchedTalents < this.searchedTalentPaginationEnd) {
      this.searchedTalentPaginationEnd = this.totalSearchedTalents
    }
  }

  // Page change for selected talents.
  changeSelectedTalentsPage(pageNumber: number): void {
    this.currentSelectedTalentPage = pageNumber
    this.selectedTalentOffset = this.currentSelectedTalentPage - 1
    this.getSelectedTalentsForBatch()
  }

  // Page change for searched talents.
  changeSearchedTalentsPage(pageNumber: number): void {
    this.currentSearchedTalentPage = pageNumber
    this.searchedTalentOffset = this.currentSearchedTalentPage - 1
    if (this.isTalentSearched) {
      this.getAllSearchedTalents()
      return
    }
    this.getAllTalents()
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

  // ============================================================= GET FUNCTIONS ==========================================================================

  // Get all batch talent details by batch id.
  getBatchTalentDetails(): void {
    this.spinnerService.loadingMessage = "Getting student list"
    this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.batchService.getBatchTalentDetails(this.batchID).subscribe((response: any) => {
      this.batchTalentList = response.body
      this.totalBatchTalents = this.batchTalentList.length
    }, (err: any) => {
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
      console.error(err);
    })
  }

  // Get all talents.
  getAllTalents(): void {
    if (this.isTalentSearched) {
      this.getAllSearchedTalents()
      return
    }
    let roleNameAndLogin: any = {
      roleName: this.localService.getJsonValue("roleName"),
      loginID: this.localService.getJsonValue("loginID")
    }
    this.spinnerService.loadingMessage = "Getting All talents"
    this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.talentService.getTalents(this.searchedTalentLimit, this.searchedTalentOffset, roleNameAndLogin).subscribe(response => {
      this.searchedTalentList = response.body
      this.totalSearchedTalents = parseInt(response.headers.get("X-Total-Count"))
    }, err => {
      console.error(this.utilService.getErrorString(err))
    }).add(() => {

      this.setSearchedTalentPaginationString()
    })
  }

  // Get all searched talents.
  getAllSearchedTalents(): void {
    let roleNameAndLogin: any = {
      roleName: this.localService.getJsonValue("roleName"),
      loginID: this.localService.getJsonValue("loginID")
    }
    let searchFormValue = this.searchedTalentForm.value
    this.utilService.deleteNullValuePropertyFromObject(searchFormValue)
    this.spinnerService.loadingMessage = "Getting All Searched Talents"
    this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.talentService.getAllSearchedTalents(searchFormValue, this.searchedTalentLimit, this.searchedTalentOffset, roleNameAndLogin).subscribe(response => {
      this.searchedTalentList = response.body
      this.totalSearchedTalents = parseInt(response.headers.get('X-Total-Count'))
    }, (error) => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {

      this.setSearchedTalentPaginationString()
    })
  }

  // Get all selected talents for batch.
  getSelectedTalentsForBatch() {
    this.spinnerService.loadingMessage = "Getting Talents For Batch"
    this.changeLoadingMessage.emit(this.spinnerService.loadingMessage)

    this.batchService.getTalentsForBatch(this.batchID, this.selectedTalentLimit, this.selectedTalentOffset).
      subscribe((response: any) => {
        this.selectedTalentList = response.body
        this.totalSelectedTalents = response.headers.get('X-Total-Count')
      }, (err: any) => {
        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        console.error(err)
      }).add(() => {

        this.setSelectedTalentPaginationString()
      })
  }

}

