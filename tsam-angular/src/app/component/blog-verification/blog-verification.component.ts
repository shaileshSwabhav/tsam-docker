import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalRef, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IBlogDTO, BlogService } from 'src/app/service/blog/blog.service';
import { ISearchFilterField, GeneralService } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-blog-verification',
  templateUrl: './blog-verification.component.html',
  styleUrls: ['./blog-verification.component.css']
})
export class BlogVerificationComponent implements OnInit {

  // Components.
  blogTopicList: any[]

  // Flags.
  isSearched: boolean

  // Blog.
  blogList: IBlogDTO[]
  blogForm: FormGroup

  // Pagination.
  limit: number
  currentPage: number
  totalBlogs: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Modal.
  modalRef: any
  @ViewChild('blogFormModal') blogFormModal: any
  @ViewChild('deleteBlogModal') deleteBlogModal: any
  @ViewChild("ckeditorContent") ckeditorContent: any

  // Spinner.



  // Search.
  blogSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  // Cke editor configuration.
  ckeEditorContentConfig: any

  // Banner image.
  selectedBlogBannerImage: string

  constructor(
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private utilService: UtilityService,
    private blogService: BlogService,
    private generalService: GeneralService,
    public utilityService: UtilityService,
    private localService: LocalService,
  ) {
    this.initializeVariables()
    this.getBlogTopicList()
    this.searchOrGetBlogs()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    // Components.
    this.blogTopicList = []

    // Blog.
    this.blogList = [] as IBlogDTO[]

    // Flags.
    this.isSearched = false

    // Pagination.
    this.limit = 5
    this.offset = 0
    this.currentPage = 0
    this.paginationStart = 0
    this.paginationEnd = 0

    // Initialize forms
    this.createBlogForm()
    this.createBlogSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Blogs"


    // Search.
    this.searchFilterFieldList = []

    // Cke editor configuration.
    this.ckeEditorContentCongiguration()

    // Initialize search form value.
    this.searchFormValue = {}
    this.searchFormValue.isVerified = "0"
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create blog form.
  createBlogForm(): void {
    this.blogForm = this.formBuilder.group({
      id: new FormControl(null),
      title: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
      content: new FormControl(null, [Validators.required, Validators.maxLength(5000)]),
      description: new FormControl(null, [Validators.required, Validators.maxLength(300)]),
      timeToRead: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(240)]),
      publishedDate: new FormControl(null),
      isVerified: new FormControl(false),
      isPublished: new FormControl(false),
      blogTopics: new FormControl(Array()),
      bannerImage: new FormControl(null),
    })
  }

  // Create blog search form.
  createBlogSearchForm(): void {
    this.blogSearchForm = this.formBuilder.group({
      title: new FormControl(null, [Validators.maxLength(200)]),
      publishedFromDate: new FormControl(null),
      publishedEndDate: new FormControl(null),
      timeToRead: new FormControl(null, [Validators.min(1), Validators.max(240)]),
      isVerified: new FormControl(null),
      isPublished: new FormControl(null),
    })
  }
  // =============================================================BLOG CRUD FUNCTIONS==========================================================================

  // On clicking view blog button.
  onViewBlogClick(blog: IBlogDTO): void {
    this.createBlogForm()
    this.selectedBlogBannerImage = blog.bannerImage
    this.blogForm.patchValue(blog)
    this.blogForm.disable()
    this.openModal(this.blogFormModal, 'xl')
  }

  // On clicking verify button.
  onVerifyButtonClik(isVerified: boolean, isPublished: boolean): void {
    if (confirm("Are you sure you want to verify this blog ?")) {
      this.spinnerService.loadingMessage = "Verifying Blog"


      let blog: any = {
        id: this.blogForm.get("id").value,
        isVerified: isVerified,
        isPublished: isPublished
      }
      this.blogService.updateBlogFlags(blog).subscribe((response: any) => {
        this.modalRef.close('success')
        this.getAllBlogs()
        alert(response)
      }, (error) => {
        alert("Blog could not be verified, please try again later")
        console.error(error)
      })
    }
  }

  // On clicking delete blog button. 
  onDeleteBlogClick(blogID: string): void {
    this.openModal(this.deleteBlogModal, 'md').result.then(() => {
      this.deleteBlog(blogID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete blog after confirmation from user.
  deleteBlog(blogID: string): void {
    this.spinnerService.loadingMessage = "Deleting Blog"


    this.blogService.deleteBlog(blogID).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllBlogs()
      alert(response)
    }, (error) => {
      alert("Blog could not be deleted, please try again later")
      console.error(error)
    })
  }

  // ============================================================BLOG SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.blogSearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate(['/admin/community/blog'])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.blogSearchForm.reset()
  }

  // Search blogs.
  searchBlogs(): void {
    this.searchFormValue = { ...this.blogSearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Blogs"
    this.changePage(1)
  }

  // Delete search criteria from blog search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.blogSearchForm.get(searchName).setValue(null)
    this.searchBlogs()
  }

  // ================================================OTHER FUNCTIONS FOR BLOG===============================================

  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllBlogs()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetBlogs() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllBlogs()
      return
    }
    this.blogSearchForm.patchValue(queryParams)
    this.searchBlogs()
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

  // Set total talents list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalBlogs < this.paginationEnd) {
      this.paginationEnd = this.totalBlogs
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

  // cke editor congiguration.
  ckeEditorContentCongiguration(): void {
    this.ckeEditorContentConfig = {
      extraPlugins: 'codeTag',
      // stylesSet: 'new_styles',
      toolbar: [
        { name: 'styles', items: ['Styles', 'Format'] },
        {
          name: 'basicstyles', groups: ['basicstyles', 'cleanup'],
          items: ['Bold', 'Italic', 'Underline', 'Strike', 'Subscript', 'Superscript', 'RemoveFormat', 'Code']
        },
        {
          name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'],
          items: ['NumberedList', 'BulletedList', '-', 'Outdent', 'Indent', '-', 'Blockquote']
        },
        { name: 'links', items: ['Link', 'Unlink'] }, //, 'Anchor' // Link
        { name: 'insert', items: ['Image'] }, //, 'Table', 'HorizontalRule' // Image
        { name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
        // { name: 'clipboard', groups: [ 'clipboard', 'undo' ], items: [ 'Cut', 'Copy', 'Paste', 'PasteText', 'PasteFromWord', '-', 'Undo', 'Redo' ] },
        // { name: 'editing', groups: [ 'find', 'selection', 'spellchecker' ], items: [ 'Scayt' ] },
        // { name: 'tools', items: [ 'Maximize' ] },
        // { name: 'others', items: [ '-' ] },
        // { name: 'about', items: [ 'About' ] }
      ],
      toolbarGroups: [
        { name: 'styles' },
        { name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
        { name: 'document', groups: ['mode', 'document', 'doctools'] },
        { name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
        { name: 'links' }, // Link
        { name: 'insert' }, // Image
        // '/',
        // { name: 'colors' },
        // { name: 'clipboard', groups: [ 'clipboard', 'undo' ] },
        // { name: 'editing', groups: [ 'find', 'selection', 'spellchecker' ] },
        // { name: 'forms' },
        // { name: 'tools' },
        // { name: 'others' },
      ],
      removeButtons: "",
      readOnly: true,
      language: 'en',
      forcePasteAsPlainText: false,
      ckfinder: {
        // Upload the images to the server using the CKFinder QuickUpload command.
        uploadUrl: '/ckfinder/core/connector/php/connector.php?command=QuickUpload&type=Files&responseType=json'
      },
      // Configure your file manager integration. This example uses CKFinder 3 for PHP.
      filebrowserBrowseUrl:
        'https://ckeditor.com/apps/ckfinder/3.4.5/ckfinder.html',
      filebrowserImageBrowseUrl:
        'https://ckeditor.com/apps/ckfinder/3.4.5/ckfinder.html?type=Images',
      filebrowserUploadUrl:
        'https://ckeditor.com/apps/ckfinder/3.4.5/core/connector/php/connector.php?command=QuickUpload&type=Files',
      filebrowserImageUploadUrl:
        'https://ckeditor.com/apps/ckfinder/3.4.5/core/connector/php/connector.php?command=QuickUpload&type=Images',
    }
  }

  // On changing tab get the component lists.
  onTabChange(event: any) {
    if (event == 1) {
      if (!this.searchFormValue) {
        this.searchFormValue = {}
      }
      this.searchFormValue.isVerified = "0"
      this.getAllBlogs()
    }
    if (event == 2) {
      if (!this.searchFormValue) {
        this.searchFormValue = {}
      }
      this.searchFormValue.isVerified = "1"
      this.getAllBlogs()
    }
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get all blogs.
  getAllBlogs() {
    this.spinnerService.loadingMessage = "Getting Blogs"


    this.blogService.getAllBlogs(this.limit, this.offset, this.searchFormValue).subscribe((response) => {
      this.totalBlogs = response.headers.get('X-Total-Count')
      this.blogList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }

  // Get blog topic list.
  getBlogTopicList(): void {
    this.generalService.getBlogTopicList().subscribe(response => {
      this.blogTopicList = response
    }, err => {
      console.error(this.utilityService.getErrorString(err))
    })
  }

}
