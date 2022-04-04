import { Component, OnInit, SecurityContext, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { DomSanitizer } from '@angular/platform-browser';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalRef, NgbModalOptions } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BlogService, IBlog, IBlogDTO } from 'src/app/service/blog/blog.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService, ISearchFilterField } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-blog',
  templateUrl: './blog.component.html',
  styleUrls: ['./blog.component.css']
})
export class BlogComponent implements OnInit {

  // Components.
  blogTopicList: any[]

  // Flags.
  isSearched: boolean
  isOperationUpdate: boolean
  isViewMode: boolean

  // Blog.
  blogList: IBlogDTO[]
  blogForm: FormGroup
  selectedBlog: any

  // Pagination.
  limit: number
  currentPage: number
  totalBlogs: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Login details.
  rolName: string
  credentialID: string

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
  blogImagesFolderPath: string
  imageUploadURL: string = environment.UPLOAD_API_PATH
  fileUploadLocation: string = environment.FILE_UPLOAD_LOACTION

  // Banner Image.
  isBannerImageUploadedToServer: boolean
  isBannerImageUploading: boolean
  bannerImageDocStatus: string
  bannerImageDisplayedFileName: string
  defaultBannerImage: string

  constructor(
    public sanitizer: DomSanitizer,
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
    private fileOps: FileOperationService,
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
    this.createBlogForm()
    this.createBlogSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Blogs"


    // Search.
    this.searchFilterFieldList = []

    // Cke editor configuration.
    this.blogImagesFolderPath = this.fileOps.BLOG_IMAGES
    this.ckeEditorContentCongiguration()

    // Login details.
    this.rolName = this.localService.getJsonValue("roleName")
    this.credentialID = this.localService.getJsonValue("credentialID")

    // Image.
    this.bannerImageDocStatus = ""
    this.bannerImageDisplayedFileName = "Select file"
    this.isBannerImageUploadedToServer = false
    this.isBannerImageUploading = false
    this.defaultBannerImage = "assets/images/blog-banner-image-default.jpg"
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

  // On clicking add new blog button.
  onAddNewBlogClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createBlogForm()
    this.openModal(this.blogFormModal, 'xl')
  }

  // Add new blog.
  addBlog(): void {
    this.spinnerService.loadingMessage = "Adding Blog"


    this.blogService.addBlog(this.blogForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllBlogs()
      alert("Blog added successfully, please wait for verification from Swabhav")
    }, (error) => {
      alert("Blog could not be added, please try again later")
      console.error(error)
    })
  }

  // On clicking view blog button.
  onViewBlogClick(blog: IBlogDTO): void {
    this.isViewMode = true
    this.createBlogForm()
    this.selectedBlog = blog
    // this.sanitizer.sanitize(SecurityContext.HTML, blog.content)
    // this.sanitizer.bypassSecurityTrustHtml(this.selectedBlog.content)
    // console.log(this.selectedBlog);

    //  Banner image.
    this.bannerImageDisplayedFileName = "No image uploaded"
    if (this.selectedBlog.bannerImage) {
      this.bannerImageDisplayedFileName = `<a href=${this.selectedBlog.bannerImage} target="_blank">Banner image present</a>`
    }

    this.blogForm.patchValue(this.selectedBlog)
    // this.blogForm.get("content").setValue(this.sanitizer.bypassSecurityTrustHtml(this.selectedBlog.content))
    this.blogForm.disable()
    this.openModal(this.blogFormModal, 'xl')
  }

  // On cliking update form button in blog form.
  onUpdateBlogClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.blogForm.enable()
  }

  // Update blog.
  updateBlog(): void {
    this.spinnerService.loadingMessage = "Updating Blog"


    this.blogForm.get("isVerified").setValue(false)
    this.blogForm.get("isPublished").setValue(false)
    this.blogForm.get("publishedDate").setValue(null)
    this.blogService.updateBlog(this.blogForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllBlogs()
      alert("Blog updated successfully, please wait for verification from Swabhav")
    }, (error) => {
      alert("Blog could not be updated, please try again later")
      console.error(error)
    })
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

  // Delete banner image.
  deleteBannerImage(): void {
    this.fileOps.deleteUploadedFile().subscribe((data: any) => {
    }, (error) => {
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
    this.router.navigate(['/community/blog/my-blogs'])
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

    this.isBannerImageUploadedToServer = false
    this.bannerImageDisplayedFileName = "Select file"
    this.bannerImageDocStatus = ""

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

  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef): void {
    if (this.isBannerImageUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    if (this.isBannerImageUploadedToServer) {
      if (!confirm("Uploaded image will be deleted.\nAre you sure you want to close?")) {
        return
      }
      this.deleteBannerImage()
    }
    modal.dismiss()
    this.isBannerImageUploadedToServer = false
    this.bannerImageDisplayedFileName = "Select file"
    this.bannerImageDocStatus = ""
  }

  // On clicking sumbit button in blog form.
  onFormSubmit(): void {
    if (this.blogForm.invalid) {
      this.blogForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateBlog()
      return
    }
    this.addBlog()
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
      extraPlugins: 'codeTag,kbdTag,simage',
      removePlugins: "exportpdf",
      // stylesSet: 'new_styles',
      toolbar: [
        { name: 'styles', items: ['Styles', 'Format'] },
        {
          name: 'basicstyles', groups: ['basicstyles', 'cleanup'],
          items: ['Bold', 'Italic', 'Underline', 'Strike', 'Subscript', 'Superscript', 'RemoveFormat', 'Code', 'Kbd']
        },
        {
          name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'],
          items: ['NumberedList', 'BulletedList', '-', 'Outdent', 'Indent', '-', 'Blockquote']
        },
        { name: 'links', items: ['Link', 'Unlink'] }, //, 'Anchor' // Link
        { name: 'insert', items: ['SImage'] }, //, 'Table', 'HorizontalRule' // Image
        { name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
      ],
      toolbarGroups: [
        { name: 'styles' },
        { name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
        { name: 'document', groups: ['mode', 'document', 'doctools'] },
        { name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
        { name: 'links' }, // Link
        { name: 'insert' }, // Image
      ],
      removeButtons: "",
      language: 'en',
      resize_enabled: false,
      width: "100%", height: "80%",
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
      folderPath: this.blogImagesFolderPath,
      imageUploadURL: this.imageUploadURL,
      fileUploadLocation: this.fileUploadLocation,
      allowedContent: true,
      extraAllowedContent: 'img'
    }
  }

  // On clicking publish button.
  onPublishButtonClik(blogID: string, isVerified: boolean, isPublished: boolean): void {
    if (confirm("Are you sure you want to publish your blog ?")) {
      this.spinnerService.loadingMessage = "Publishing Blog"


      let blog: any = {
        id: blogID,
        isVerified: isVerified,
        isPublished: isPublished
      }
      this.blogService.updateBlogFlags(blog).subscribe((response: any) => {
        this.getAllBlogs()
        alert("Your blog has been published successfully")
      }, (error) => {
        alert("Blog could not be published, please try again later")
        console.error(error)
      })
    }
  }

  // On uplaoding banner image.
  onBannerImageSelect(event: any): void {
    this.bannerImageDocStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]

      // Upload banner image if it is present.
      this.isBannerImageUploading = true
      this.fileOps.uploadBlogBannerImage(file).subscribe((response: any) => {
        this.blogForm.markAsDirty()
        this.blogForm.patchValue({
          bannerImage: response
        })
        this.bannerImageDisplayedFileName = file.name
        this.isBannerImageUploading = false
        this.isBannerImageUploadedToServer = true
        this.bannerImageDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
      }, (error) => {
        this.isBannerImageUploading = false
        this.bannerImageDocStatus = `<p><span>&#10060;</span> ${error}</p>`
      })
    }
  }

  // Format fields of blog list.
  formatBlogListFields(): void {

    // Give default banner image to all blogs 
    for (let i = 0; i < this.blogList.length; i++) {
      if (!this.blogList[i].bannerImage) {
        this.blogList[i].bannerImage = this.defaultBannerImage
      }
    }
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get all blogs.
  getAllBlogs() {
    this.spinnerService.loadingMessage = "Getting Blogs"


    if (!this.searchFormValue) {
      this.searchFormValue = {}
    }
    this.searchFormValue.authorID = this.credentialID
    this.blogService.getAllBlogs(this.limit, this.offset, this.searchFormValue).subscribe((response) => {
      this.totalBlogs = response.headers.get('X-Total-Count')
      this.blogList = response.body
      this.formatBlogListFields()
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
