import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalRef, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BlogTopicService, IBlogTopic } from 'src/app/service/blog-topic/blog-topic.service';
import { ISearchFilterField } from 'src/app/service/general/general.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-blog-topic',
  templateUrl: './blog-topic.component.html',
  styleUrls: ['./blog-topic.component.css']
})
export class BlogTopicComponent implements OnInit {

  // Flags.
  isSearched: boolean
  isOperationUpdate: boolean
  isViewMode: boolean

  // Blog topic.
  blogTopicList: IBlogTopic[]
  blogTopicForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  totalBlogTopics: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Modal.
  modalRef: any
  @ViewChild('blogTopicFormModal') blogTopicFormModal: any
  @ViewChild('deleteBlogTopicModal') deleteBlogTopicModal: any

  // Spinner.



  // Search.
  blogTopicSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  constructor(
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private utilService: UtilityService,
    private blogTopicService: BlogTopicService,
  ) {
    this.initializeVariables()
    this.searchOrGetBlogTopics()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    // Components.
    this.blogTopicList = [] as IBlogTopic[]

    // Flags.
    this.isOperationUpdate = false
    this.isViewMode = false
    this.isSearched = false

    // Pagination.
    this.limit = 5
    this.offset = 0
    this.currentPage = 0
    this.paginationStart = 0
    this.paginationEnd = 0

    // Initialize forms
    this.createBlogTopicForm()
    this.createBlogTopicSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Blog Topics"


    // Search.
    this.searchFilterFieldList = []
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create blog topic form.
  createBlogTopicForm(): void {
    this.blogTopicForm = this.formBuilder.group({
      id: new FormControl(null),
      name: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
    })
  }

  // Create blog topic search form.
  createBlogTopicSearchForm(): void {
    this.blogTopicSearchForm = this.formBuilder.group({
      name: new FormControl(null),
    })
  }
  // =============================================================BLOG TOPIC CRUD FUNCTIONS==========================================================================
  // On clicking add new blog topic button.
  onAddNewBlogTopicClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createBlogTopicForm()
    this.openModal(this.blogTopicFormModal, 'sm')
  }

  // Add new blog topic.
  addBlogTopic(): void {
    this.spinnerService.loadingMessage = "Adding Blog Topic"


    this.blogTopicService.addBlogTopic(this.blogTopicForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllBlogTopics()
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

  // On clicking view blog topic button.
  onViewBlogTopicClick(blogTopic: IBlogTopic): void {
    this.isViewMode = true
    this.createBlogTopicForm()
    this.blogTopicForm.patchValue(blogTopic)
    this.blogTopicForm.disable()
    this.openModal(this.blogTopicFormModal, 'sm')
  }

  // On cliking update form button in blog topic form.
  onUpdateBlogTopicClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.blogTopicForm.enable()
  }

  // Update blog topic.
  updateBlogTopic(): void {
    this.spinnerService.loadingMessage = "Updating Blog Topic"


    this.blogTopicService.updateBlogTopic(this.blogTopicForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllBlogTopics()
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

  // On clicking delete blog topic button. 
  onDeleteBlogTopicClick(blogTopicID: string): void {
    this.openModal(this.deleteBlogTopicModal, 'md').result.then(() => {
      this.deleteBlogTopic(blogTopicID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete blog topic after confirmation from user.
  deleteBlogTopic(blogTopicID: string): void {
    this.spinnerService.loadingMessage = "Deleting Blog Topic"


    this.blogTopicService.deleteBlogTopic(blogTopicID).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllBlogTopics()
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

  // ============================================================BLOG TOPIC SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.blogTopicSearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate(['/admin/models/blog-topic'])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.blogTopicSearchForm.reset()
  }

  // Search blog topics.
  searchBlogTopics(): void {
    this.searchFormValue = { ...this.blogTopicSearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Blog Topics"
    this.changePage(1)
  }

  // Delete search criteria from blog topic search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.blogTopicSearchForm.get(searchName).setValue(null)
    this.searchBlogTopics()
  }

  // ================================================OTHER FUNCTIONS FOR BLOG TOPIC===============================================
  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllBlogTopics()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetBlogTopics() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllBlogTopics()
      return
    }
    this.blogTopicSearchForm.patchValue(queryParams)
    this.searchBlogTopics()
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

  // On clicking sumbit button in blog topic form.
  onFormSubmit(): void {
    if (this.blogTopicForm.invalid) {
      this.blogTopicForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateBlogTopic()
      return
    }
    this.addBlogTopic()
  }

  // Set total talents list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalBlogTopics < this.paginationEnd) {
      this.paginationEnd = this.totalBlogTopics
    }
  }

  // Compare for select option field.
  compareFn(optionOne: any, optionTwo: any): boolean {
    if (optionOne == null && optionTwo == null) {
      return true
    }
    if (optionTwo != undefined && optionOne != undefined) {
      return optionOne.id === optionTwo.id
    }
    return false
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get all blog topics.
  getAllBlogTopics() {
    this.spinnerService.loadingMessage = "Getting Blog Topics"


    this.blogTopicService.getAllBlogTopics(this.limit, this.offset, this.searchFormValue).subscribe((response) => {
      this.totalBlogTopics = response.headers.get('X-Total-Count')
      this.blogTopicList = response.body
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
