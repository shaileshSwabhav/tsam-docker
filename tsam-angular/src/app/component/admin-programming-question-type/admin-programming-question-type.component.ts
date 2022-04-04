import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UrlConstant } from 'src/app/service/constant';
import { GeneralService, ISearchFilterField } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IProgrammingQuestionType, ProgrammingQuestionService } from 'src/app/service/programming-question/programming-question.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-programming-question-type',
  templateUrl: './admin-programming-question-type.component.html',
  styleUrls: ['./admin-programming-question-type.component.css']
})
export class AdminProgrammingQuestionTypeComponent implements OnInit {

  // Flags.
  isSearched: boolean
  isViewClicked: boolean
  isUpdateClicked: boolean

  // Programming Question Type.
  programmingQuestionType: IProgrammingQuestionType[]
  programmingQuestionTypeForm: FormGroup
  totalProgrammingQuestionTypes: number

  // Pagination.
  limit: number
  currentPage: number
  offset: number
  paginationString: string

  // Modal.
  modalRef: any
  @ViewChild('programmingQuestionTypeModal') programmingQuestionTypeModal: any
  @ViewChild('deleteModal') deleteModal: any

  // Spinner.



  // Search.
  programmingQuestionTypeSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  // access
  permission: IPermission

  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private questionService: ProgrammingQuestionService,
    private urlConstant: UrlConstant,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.initalizeVariables()
    this.createProgrammingQuestionTypeSearchForm()
    this.searchOrGetProgrammingQuestionType()
  }

  initalizeVariables(): void {
    this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_PROGRAMMING_QUESTION_TYPE)

    this.programmingQuestionType = []

    this.isSearched = false
    this.isViewClicked = false
    this.isUpdateClicked = false

    this.limit = 5
    this.offset = 0
    this.totalProgrammingQuestionTypes = 0
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  createProgrammingQuestionTypeForm(): void {
    this.programmingQuestionTypeForm = this.formBuilder.group({
      id: new FormControl(null),
      programmingType: new FormControl(null),
    })
  }

  createProgrammingQuestionTypeSearchForm(): void {
    this.programmingQuestionTypeSearchForm = this.formBuilder.group({
      programmingType: new FormControl(null),
    })
  }

  // page change function
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber;
    this.offset = this.currentPage - 1;
    this.getProgrammingQuestionTypes();
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetProgrammingQuestionType() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getProgrammingQuestionTypes()
      return
    }
    this.programmingQuestionTypeSearchForm.patchValue(queryParams)
    this.searchProgrammingQuestion()
  }

  searchProgrammingQuestion(): void {
    // console.log(this.searchBatchForm.value)
    this.searchFormValue = { ...this.programmingQuestionTypeSearchForm?.value }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
    let flag: boolean = true

    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
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
    this.changePage(1)
  }

  onAddClick(): void {
    this.isViewClicked = false
    this.isUpdateClicked = false
    this.createProgrammingQuestionTypeForm()
    this.openModal(this.programmingQuestionTypeModal, "md")
  }

  onViewClick(programmingQuestionType: IProgrammingQuestionType): void {
    this.isViewClicked = true
    this.isUpdateClicked = false
    this.createProgrammingQuestionTypeForm()
    this.programmingQuestionTypeForm.patchValue(programmingQuestionType)
    this.programmingQuestionTypeForm.disable()
    this.openModal(this.programmingQuestionTypeModal, "md")
  }

  onUpdateClick(): void {
    this.isViewClicked = false
    this.isUpdateClicked = true
    this.programmingQuestionTypeForm.enable()
  }

  getProgrammingQuestionTypes(): void {

    if (this.isSearched) {
      this.getSearchedProgrammingQuestionTypes()
      return
    }

    this.spinnerService.loadingMessage = "Getting programming question types"

    this.questionService.getProgrammingQuestionTypes(this.limit, this.offset).subscribe((response: any) => {
      this.programmingQuestionType = response.body
      this.totalProgrammingQuestionTypes = response.headers.get('X-Total-Count')
    }, (err: any) => {
      this.totalProgrammingQuestionTypes = 0
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

  getSearchedProgrammingQuestionTypes(): void {
    this.spinnerService.loadingMessage = "Getting programming question types"

    this.questionService.getProgrammingQuestionTypes(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      this.programmingQuestionType = response.body
      this.totalProgrammingQuestionTypes = response.headers.get('X-Total-Count')
    }, (err: any) => {
      this.totalProgrammingQuestionTypes = 0
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

  addProgrammingQuestionType(): void {
    this.spinnerService.loadingMessage = "Adding programming question type"

    this.questionService.addProgrammingQuestionType(this.programmingQuestionTypeForm.value).subscribe((response: any) => {
      this.modalRef.close()
      alert("Programming question type successfully added")
      this.getProgrammingQuestionTypes()
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  updateProgrammingQuestionType(): void {
    this.spinnerService.loadingMessage = "Updating programming question type"

    this.questionService.updateProgrammingQuestionType(this.programmingQuestionTypeForm.value).subscribe((response: any) => {
      this.modalRef.close()
      alert("Programming question type successfully updated")
      this.getProgrammingQuestionTypes()
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  onDeleteClick(programmingTypeID: string): void {
    this.openModal(this.deleteModal, 'md').result.then(() => {
      this.deleteProgrammingQuestionType(programmingTypeID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  deleteProgrammingQuestionType(programmingTypeID: string): void {
    this.spinnerService.loadingMessage = "Deleting programming question type"

    this.questionService.deleteProgrammingQuestionType(programmingTypeID).subscribe((response: any) => {
      this.modalRef.close()
      alert("Programming question type successfully deleted")
      this.getProgrammingQuestionTypes()
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  openModal(modalContent: any, modalSize?: string): NgbModalRef {
    if (!modalSize) {
      modalSize = 'lg'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: modalSize, centered: true,
    }
    this.modalRef = this.modalService.open(modalContent, options)
    return this.modalRef
  }

  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalProgrammingQuestionTypes < end) {
      end = this.totalProgrammingQuestionTypes
    }
    if (this.totalProgrammingQuestionTypes == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end}`
  }

  //validate add/update form
  validate(): void {
    // console.log(this.programmingQuestionTypeForm.controls);

    if (this.programmingQuestionTypeForm.invalid) {
      this.programmingQuestionTypeForm.markAllAsTouched();
      return
    }

    if (this.isUpdateClicked) {
      this.updateProgrammingQuestionType()
      return
    }
    this.addProgrammingQuestionType()
  }

  resetSearchAndGetAll(): void {
    this.createProgrammingQuestionTypeSearchForm()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([this.urlConstant.TRAINING_PROGRAMMING_QUESTION_TYPE])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.programmingQuestionTypeSearchForm.reset()
  }
}
