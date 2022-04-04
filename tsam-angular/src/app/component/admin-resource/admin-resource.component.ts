import { HttpParams } from '@angular/common/http';
import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { startWith, pairwise } from 'rxjs/operators';
import { Role, UploadConstant, UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService, IFileTypeList } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IResource, IResourceDownload, IResourceLike, ResourceService } from 'src/app/service/resource/resource.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ITechnology, TechnologyService } from 'src/app/service/technology/technology.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-resource',
  templateUrl: './admin-resource.component.html',
  styleUrls: ['./admin-resource.component.css']
})
export class AdminResourceComponent implements OnInit {

  // forms
  resourceSearchForm: FormGroup
  resourceForm: FormGroup

  // resources
  resources: IResource[]
  isResourceAvailable: boolean

  // resource-download
  isResourceDownload: boolean
  resourceDownload: IResourceDownload[]
  resourceDownloadForm: FormGroup

  // resource-like
  isResourceLike: boolean
  resourceLike: IResourceLike[]
  resourceLikeForm: FormGroup

  // components
  resourceTypeList: any[]
  fileTypeList: IFileTypeList[]
  technologyList: ITechnology[]
  resourceSubTypeList: any[]

  // pagination
  limit: number
  offset: number
  currentPage: number
  totalResources: number
  paginationString: string
  techLimit: number
  techOffset: number

  // flags
  isOperationUpdate: boolean
  isViewMode: boolean
  isTechLoading: boolean

  // spinner

  loadingLikeContent: string

  // access
  permission: IPermission
  isAdmin: boolean
  isFaculty: boolean
  loginID: string

  // modal
  modalRef: any

  // resource
  isResourceUploadedToServer: boolean
  isFileUploading: boolean
  docStatus: string
  displayedFileName: string
  previousFileType: string
  allowedExtension: string

  // preview image
  isPreviewUploadedToServer: boolean
  isPreviewFileUploading: boolean
  previewImageDocStatus: string
  previewDisplayedFileName: string

  // search
  searchFormValue: any
  showSearch: boolean
  isSearched: boolean

  // tab
  activeTab: number

  // fileType count
  fileTypeCount: any
  fileTypeMap: any

  // private readonly RESOURCE = "resources"
  readonly SWABHAV_DOCS = "Swabhav_Docs"

  readonly IMAGE = "Image"
  readonly AUDIO = "Audio"
  readonly VIDEO = "Video"
  readonly URL = "URL"
  readonly DOCUMENT = "Document"
  readonly BOOK = "Book"
  readonly PREVIEW = "Preview"
  private readonly URLregex = /^(http(s)?:\/\/)[\w.-]+(\.[\w\.-]+)+[\w\-\._~:?#[\]@!\/$&'\(\)\*\+,;=.]+$/

  // modal
  @ViewChild('resourceFormModal') resourceFormModal: ElementRef
  @ViewChild('resourceDownloadModal') resourceDownloadModal: ElementRef
  @ViewChild('resourceLikeModal') resourceLikeModal: ElementRef
  @ViewChild('deleteConfirmationModal') deleteConfirmationModal: ElementRef

  constructor(
    private formBuilder: FormBuilder,
    private generalService: GeneralService,
    private techService: TechnologyService,
    private resourceService: ResourceService,
    private urlConstant: UrlConstant,
    private localService: LocalService,
    private utilService: UtilityService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private fileOperationService: FileOperationService,
    private uploadConstant: UploadConstant,
    private role: Role,
  ) {
    this.inititalizeVariables()
    this.getAllComponents()
  }

  inititalizeVariables(): void {
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.TRAINING_RESOURCE)
    }
    if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.BANK_RESOURCE)
    }
    this.loginID = this.localService.getJsonValue("credentialID")

    this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
    this.isAdmin = (this.localService.getJsonValue("roleName") == this.role.ADMIN ? true : false)

    if (!this.isFaculty && !this.isAdmin) {
      this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'), this.urlConstant.RESOURCE)
    }

    // console.log(this.localService.getJsonValue("roleName"));

    this.allowedExtension = ""

    this.limit = 9
    this.offset = 0
    this.techLimit = 10
    this.techOffset = 0

    this.isOperationUpdate = false
    this.isSearched = false
    this.showSearch = false
    this.isResourceAvailable = true
    this.isTechLoading = false

    this.isResourceDownload = true
    this.isResourceLike = true

    this.resources = [] as IResource[]
    this.resourceDownload = [] as IResourceDownload[]
    this.resourceLike = [] as IResourceLike[]

    this.resourceTypeList = []
    this.fileTypeList = []
    this.resourceSubTypeList = []


    this.previousFileType = null
    this.activeTab = 1

    this.fileTypeMap = new Map()
    this.fileTypeCount = new Map()

    this.searchFormValue = {}

    this.createResourceSearchForm()
  }

  getAllComponents(): void {
    this.getFileTypeCount()
    this.getResourceSubType()
    this.getResourceType()
    this.getTechnologyList()
    this.getFileType()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  operation: string

  onRedirect(): void {
    // console.log("resourceSearchForm.value -> ", this.resourceSearchForm.value);
    if (this.resourceSearchForm.get("operation")?.value === "add") {
      this.onAddClick()
      return
    }
  }

  createResourceForm() {
    this.resetFileUploadFields()

    this.resourceForm = this.formBuilder.group({
      id: new FormControl(null),
      resourceName: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
      isExistingResource: new FormControl(false),
      resourceType: new FormControl(null, [Validators.required]),
      resourceSubType: new FormControl(null),
      fileType: new FormControl(null, [Validators.required]),
      resourceURL: new FormControl(null, [Validators.required]),
      previewURL: new FormControl(null),
      description: new FormControl(null, [Validators.maxLength(250)]),
      technologyID: new FormControl(null),
      technology: new FormControl(null),

      // control for book.
      isBook: new FormControl(false, [Validators.required]),
      author: new FormControl(null),
      publication: new FormControl(null),
    })

    this.emitPreviousFileType()
  }

  // emitPreviousFileType will fetch previously selected file type.
  emitPreviousFileType(): void {

    // Fill buffer with initial value.  Will emit immediately.
    this.resourceForm.get('fileType')
      .valueChanges.pipe(
        startWith(null as string),
        pairwise()).subscribe(([prev, next]: [any, any]) => {
          this.previousFileType = prev
        })
  }

  createResourceSearchForm() {
    this.resourceSearchForm = this.formBuilder.group({
      fileType: new FormControl(this.DOCUMENT),
      resourceType: new FormControl(((!this.isAdmin && !this.isFaculty) ? this.SWABHAV_DOCS : null)),
      resourceSubType: new FormControl(null),
      resourceName: new FormControl(null, [Validators.maxLength(100)]),
      credentialID: new FormControl(this.loginID),
      technologyID: new FormControl(null),
      operation: new FormControl(null),
    })
  }

  onAddClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createResourceForm()
    this.openModal(this.resourceFormModal, 'xl')
  }

  onUpdateClick(resource: IResource): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.createResourceForm()

    this.resourceForm.patchValue(resource)

    this.displayedFileName = "No resource uploaded"
    if (resource.resourceURL) {
      this.displayedFileName = `<a href=${resource.resourceURL} target="_blank">${resource.resourceName}</a>`
    }
    if (resource.previewURL) {
      this.previewDisplayedFileName = `<a href=${resource.previewURL} target="_blank">Preview Present</a>`
    }
    this.openModal(this.resourceFormModal, 'xl')
  }

  onDeleteClick(resourceID: string): void {
    this.openModal(this.deleteConfirmationModal, 'md').result.then(() => {
      this.deleteResource(resourceID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  onIsBookChange(): void {
    if (!this.resourceForm.get("isBook")?.value) {
      this.resourceForm.get("author")?.setValue(null)
      this.resourceForm.get("publication")?.setValue(null)
    }
  }

  checkFileType() {
    this.allowedExtension = ""
    if (this.resourceForm.get("fileType") != null) {
      this.setAllowedExtensions()
    }

    if (this.resourceForm.get('resourceURL').value && !this.onFileTpeChange()) {
      return
    }

    this.resetFileUploadFields()
    this.resetResourceURLField()
  }

  setAllowedExtensions(): void {
    if (this.resourceForm.get("fileType").value == this.AUDIO) {
      this.allowedExtension = this.uploadConstant.ALLOWED_AUDIO_EXTENSION
    } else if (this.resourceForm.get("fileType").value == this.IMAGE) {
      this.allowedExtension = this.uploadConstant.ALLOWED_IMAGE_EXTENSION
    } else {
      this.allowedExtension = this.uploadConstant.ALLOWED_DOCUMENT_EXTENSION
    }
  }

  onFileTpeChange(): boolean {
    if (confirm("Uploaded resource will be deleted.\nAre you sure you want to continue?")) {
      this.resourceForm.get('resourceURL').reset()
      return true
    }
    this.resourceForm.get('fileType').setValue(this.previousFileType)
    return false
  }

  resetFileUploadFields(): void {
    this.displayedFileName = "Select File"
    this.docStatus = ""
    this.isResourceUploadedToServer = false

    this.previewDisplayedFileName = "Select Image"
    this.previewImageDocStatus = ""
    this.isPreviewUploadedToServer = false
  }

  resetResourceURLField(): void {

    if (this.resourceForm.get('fileType').value == this.VIDEO ||
      this.resourceForm.get('fileType').value == this.URL) {
      this.resourceForm.get('resourceURL').setValidators([Validators.required, Validators.pattern(this.URLregex)])
      this.resourceForm.get('resourceURL').setValue(null)
    } else {
      this.resourceForm.get('resourceURL').setValidators([Validators.required])
    }
    this.resourceForm.get("resourceURL").updateValueAndValidity()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetResources() {
    let queryParams = this.route.snapshot.queryParams

    if (this.utilService.isObjectEmpty(queryParams)) {
      this.spinnerService.loadingMessage = "Getting " + this.resourceSearchForm.get("fileType").value.toLowerCase()
      this.showSearch = false
      this.searchFormValue.fileType = this.resourceSearchForm.get("fileType").value

      if (!this.isAdmin && !this.isFaculty) {
        this.searchFormValue.resourceType = this.SWABHAV_DOCS
      }
      this.searchFormValue.credentialID = this.loginID

      this.addQueryParams()
      this.getResources()
      return
    }

    // if (queryParams.operation) {
    // 	this.operation = queryParams.operation
    //   this.onRedirect()
    //   // return
    // }

    this.resourceSearchForm.patchValue(queryParams)
    this.setActiveTab(this.resourceSearchForm.get("fileType").value)
    this.spinnerService.loadingMessage = "Getting " + this.resourceSearchForm.get("fileType").value.toLowerCase()
    this.searchResources()

  }

  setActiveTab(fileType: string): void {
    this.activeTab = this.fileTypeMap.get(fileType)
  }

  resetSearchAndGetAll(): void {
    this.resetSearchForm()
    this.searchFormValue = {}
    this.searchFormValue.fileType = this.resourceSearchForm.get("fileType").value
    if (!this.isAdmin && !this.isFaculty) {
      this.searchFormValue.resourceType = this.SWABHAV_DOCS
    }
    this.searchFormValue.credentialID = this.loginID
    this.getFileTypeCount()
    this.addQueryParams()
    this.changePage(1)

    this.isSearched = false
    this.showSearch = false
  }

  addQueryParams(): void {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
  }

  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalResources < end) {
      end = this.totalResources
    }
    if (this.totalResources == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start}-${end} of ${this.totalResources}`
  }

  onTabChange(event: any) {
    // console.log(event)
    if (this.resourceSearchForm.get("fileType").value == event) {
      return
    }
    this.resourceSearchForm.get('fileType').setValue(event)
    this.spinnerService.loadingMessage = "Getting " + event
    this.setActiveTab(this.resourceSearchForm.get("fileType").value)
    this.searchResources(true)
    // this.changePage(1)
  }

  searchResources(isTabChanged?: boolean): void {
    this.searchFormValue = { ...this.resourceSearchForm?.value }

    let flag: boolean = true

    if (!this.isAdmin && !this.isFaculty) {
      this.searchFormValue.resourceType = this.SWABHAV_DOCS
    }

    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field]
      } else {
        if (!this.isAdmin && !this.isFaculty) {
          // if (field != "fileType" || (field != "fileType" && field != "resourceType")) {
          if (field != "fileType" && field != "resourceType") {
            this.isSearched = true
            this.showSearch = true
          }
        } else if (field != "fileType") {
          this.showSearch = true
          this.isSearched = true
        }
        flag = false
      }
    }

    if (flag) {
      return
    }

    if (isTabChanged == undefined || isTabChanged == null || !isTabChanged) {
      this.getFileTypeCount()
    }
    this.addQueryParams()
    this.changePage(1)
  }

  // page change function
  changePage(event: any): void {
    // console.log(event)
    this.currentPage = event;
    this.offset = event - 1;
    this.getResources();
  }

  getResources(): void {
    this.totalResources = 0
    this.resources = []
    this.isResourceAvailable = true
    // console.log(this.searchFormValue);
    this.resourceService.getAllResources(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      this.resources = response.body
      console.log(this.resources)
      this.totalResources = response.headers.get('X-Total-Count')
      this.onRedirect()
      this.resourceSearchForm.get("operation").setValue(null)
    }, (error: any) => {
      this.totalResources = 0
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error?.statusText)
    }).add(() => {
      this.setPaginationString()
      if (this.totalResources == 0) {
        this.isResourceAvailable = false
      }
    })
  }

  addResource(): void {
    this.spinnerService.loadingMessage = "Adding resource"

    this.resourceService.addResource(this.resourceForm.value).subscribe(() => {
      this.modalRef.close()
      alert("Resource successfully added")
      this.getFileTypeCount()
      this.getResources()
      this.resourceForm.reset()
    }, (error) => {
      console.error(error)
      // 
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  updateResource(): void {
    this.spinnerService.loadingMessage = "Updating resource"

    this.resourceService.updateResource(this.resourceForm.value).subscribe(() => {
      this.modalRef.close()
      alert("Resource successfully updated")
      this.getFileTypeCount()
      this.getResources()
      this.resourceForm.reset()
    }, (error) => {
      console.error(error)
      // 
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  deleteResource(resourceID: string): void {
    this.spinnerService.loadingMessage = "Deleting resource"

    this.resourceService.deleteResource(resourceID).subscribe(() => {
      this.modalRef.close()
      alert("Resource successfully deleted")
      this.getFileTypeCount()
      this.getResources()
      // this.resourceForm.reset()
    }, (error) => {
      console.error(error)
      // 
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
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

  // Used to open modal.
  openScrollableModal(content: any, size?: string): NgbModalRef {
    if (!size) {
      size = 'lg'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: size, scrollable: true, centered: true,
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef) {
    if (this.isFileUploading || this.isPreviewFileUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    if (this.isResourceUploadedToServer || this.isPreviewUploadedToServer) {
      if (!confirm("Uploaded resource will be deleted.\nAre you sure you want to close?")) {
        return
      }
      // #niranjan
      // this.deleteResume()
      // let file = this.resourceForm.get("resourceURL").value.split(this.RESOURCE)
      // this.deleteUploadedResouce(this.RESOURCE+file[1])
    }

    modal.dismiss()
    this.resetFileUploadFields()
  }


  //Delete uploaded resource
  // deleteUploadedResouce(fileName?: string) {
  //   this.fileOperationService.deleteUploadedFile(fileName).subscribe((data: any) => {
  //     // console.log(data)
  //   }, (error) => {
  //     console.error(error);
  //   })
  // }

  patchIDFromObjects(): void {
    if (this.resourceForm.get("technology")?.value) {
      let technologyID = this.resourceForm.get("technology").value.id
      this.resourceForm.get("technologyID").setValue(technologyID)
    }
  }

  doesURLExist(): boolean {
    return this.resourceForm.get("resourceURL").value
  }

  onSubmit(): void {
    this.patchIDFromObjects()

    if (!this.doesURLExist()) {
      alert("Either file should be uploaded or link should be specified.")
      return
    }

    // console.log(this.resourceForm.controls)

    if (this.resourceForm.invalid) {
      this.resourceForm.markAllAsTouched();
      return
    }

    if (this.isFileUploading) {
      alert("Please wait till file is being uploaded")
      return
    }

    if (this.isOperationUpdate) {
      this.updateResource()
      return
    }
    this.addResource()
  }

  onResourceSelect(event: any, fileType: string): void {
    let files = event.target.files
    if (this.PREVIEW === fileType) {
      this.previewFileSelect(files)
      return
    }
    this.docFileSelect(files, fileType)
  }

  previewFileSelect(files: any): void {
    this.previewImageDocStatus = ""
    if (files && files.length) {
      let file = files[0]
      // check if file extension is valid
      let err = this.fileOperationService.isImageFileValid(file)
      if (err != null) {
        this.previewImageDocStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // Upload document if it is present.
      this.uploadPreview(file)
    }
  }

  docFileSelect(files: any, fileType: string): void {
    this.docStatus = ""
    if (files && files.length) {
      let file = files[0]
      // check if file extension is valid
      let err = this.checkFileExtension(file, fileType)
      if (err != null) {
        this.docStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // Upload document if it is present.
      this.uploadFile(file)
    }
  }

  checkFileExtension(file: any, fileType: string): Error {

    switch (fileType) {
      // case this.DOCUMENT || this.BOOK:
      //   return this.fileOperationService.isDocumentFileValid(file)
      case this.AUDIO:
        return this.fileOperationService.isAudioFileValid(file)
      case this.IMAGE:
        return this.fileOperationService.isImageFileValid(file)
      case this.PREVIEW:
        return this.fileOperationService.isImageFileValid(file)

      default:
        return this.fileOperationService.isDocumentFileValid(file)
    }
  }

  // uploads file to the server
  uploadFile(file: any): void {
    this.isFileUploading = true
    this.fileOperationService.uploadResource(file).subscribe((data: any) => {
      // console.log(data);
      this.resourceForm.markAsDirty()
      this.resourceForm.patchValue({
        resourceURL: data
      })
      this.displayedFileName = file.name
      this.isFileUploading = false
      this.isResourceUploadedToServer = true
      this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
    }, (error: any) => {
      this.isFileUploading = false
      this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
    })
  }

  // uploads file to the server
  uploadPreview(file: any): void {
    this.isPreviewFileUploading = true
    this.fileOperationService.uploadResource(file, this.PREVIEW).subscribe((data: any) => {
      // console.log(data);
      this.resourceForm.markAsDirty()
      this.resourceForm.patchValue({
        previewURL: data
      })
      this.previewDisplayedFileName = file.name
      this.isPreviewFileUploading = false
      this.isPreviewUploadedToServer = true
      this.previewImageDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
    }, (error: any) => {
      this.isPreviewFileUploading = false
      this.previewImageDocStatus = `<p><span>&#10060;</span> ${error}</p>`
    })
  }

  deletePreviousResource(): void {
    // let file = this.resourceForm.get("resourceURL").value.split(this.RESOURCE)
    // this.deleteUploadedResouce(this.RESOURCE+file[1])
  }

  getFileType(): void {
    this.spinnerService.loadingMessage = "Getting resource"


    this.generalService.getGeneralTypeByType("file_type").subscribe((respond: any) => {
      this.fileTypeList = respond
      // this.fileType = this.fileTypeList[0].value
      for (let index = 0; index < this.fileTypeList.length; index++) {
        this.fileTypeMap.set(this.fileTypeList[index].value, this.fileTypeList[index].key)
      }
      this.searchOrGetResources()
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  getResourceType(): void {
    this.generalService.getGeneralTypeByType("resource_type").subscribe((respond: any) => {
      this.resourceTypeList = respond;
    }, (err) => {
      console.error(err)
    })
  }

  getResourceSubType(): void {
    this.generalService.getGeneralTypeByType("resource_sub_type").subscribe((respond: any) => {
      this.resourceSubTypeList = respond;
    }, (err) => {
      console.error(err)
    })
  }

  // Get Technology List.
  getTechnologyList(event?: any): void {
    // this.generalService.getTechnologies().subscribe((respond: any[]) => {
    //   this.technologyList = respond;
    // }, (err) => {
    //   console.error(this.utilService.getErrorString(err))
    // })

    let queryParams: any = {}
    if (event && event.term != "") {
      queryParams.language = event.term
    }
    this.isTechLoading = true
    this.techService.getAllTechnologies(this.techLimit, this.techOffset, queryParams).subscribe((response) => {
      // console.log("getTechnology -> ", response);
      this.technologyList = []
      this.technologyList = this.technologyList.concat(response.body)
    }, (err) => {
      console.error(err)
    }).add(() => {
      this.isTechLoading = false
    })
  }

  getFileTypeCount(): void {
    let params: any = {}

    if (this.searchFormValue.resourceName != null) {
      params.resourceName = this.searchFormValue.resourceName
    }

    if (this.searchFormValue.technologyID != null) {
      params.technologyID = this.searchFormValue.technologyID
    }

    if (this.searchFormValue.resourceSubType != null) {
      params.resourceSubType = this.searchFormValue.resourceSubType
    }

    if (!this.isAdmin && !this.isFaculty) {
      params.resourceType = this.SWABHAV_DOCS
    } else if (this.searchFormValue.resourceType != null) {
      params.resourceType = this.searchFormValue.resourceType
    }

    this.resourceService.getResourceCount(params).subscribe((response: any) => {
      for (let index = 0; index < response.body.length; index++) {
        this.fileTypeCount.set(response.body[index].fileType, response.body[index].totalCount)
      }
    }, (err) => {
      console.error(err)
    })
  }

  resetSearchForm(): void {
    let fileType = this.resourceSearchForm.get("fileType").value

    if (this.isAdmin || this.isFaculty) {
      this.resourceSearchForm.reset({
        fileType: fileType,
        credentialID: this.loginID,
      })
      // this.resourceSearchForm.get("fileType").setValue(fileType)
      // this.resourceSearchForm.get("credentialID").setValue(this.loginID)
      return
    }
    this.resourceSearchForm.get('resourceName').reset()
    this.resourceSearchForm.get("fileType").setValue(fileType)
    this.resourceSearchForm.get("credentialID").setValue(this.loginID)
  }

  // getThumbnail(videoURL: string): string {
  //   const VID_REGEX = /(?:youtube(?:-nocookie)?\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})/
  //   return videoURL.match(VID_REGEX) === null ? null : 'http://img.youtube.com/vi/' + videoURL.match(VID_REGEX)[1] + '/mqdefault.jpg'
  // }

  // getSanitizedURL(videoURL: string): SafeUrl {
  //   return this.sanitizer.bypassSecurityTrustResourceUrl(videoURL);
  // }

  // ======================================================== RESOURCE DOWNLOAD ========================================================

  createResourceDownload(): void {
    this.resourceDownloadForm = this.formBuilder.group({
      id: new FormControl(null),
      resourceID: new FormControl(null, [Validators.required]),
      credentialID: new FormControl(null, [Validators.required]),
    })
  }

  onDownloadClick(resource: IResource): void {

    this.createResourceDownload()
    this.resourceDownloadForm.get('resourceID').setValue(resource.id)
    this.resourceDownloadForm.get('credentialID').setValue(this.loginID)
    this.addResourceDownload(resource)
  }

  addResourceDownload(resource: IResource): void {

    this.resourceService.addResourceDownload(this.resourceDownloadForm.value).subscribe(response => {
      // console.log(response)
      this.getResourceDownloadCount(resource)
    }, (error: any) => {
      console.error(error);
    })
  }

  getResourceDownloadCount(resource: IResource): void {
    this.resourceService.getResourceDownloadCount(resource.id).subscribe((response: any) => {
      resource.totalDownload = response.body.totalCount
    }, (error: any) => {
      console.error(error);
    })
  }

  getResourceDownload(resource: IResource): void {

    this.spinnerService.loadingMessage = "Getting downloaded by details"

    this.resourceDownload = []
    this.isResourceDownload = true
    this.openScrollableModal(this.resourceDownloadModal, "md")
    this.resourceService.getResourceDownload(resource.id).subscribe((response: any) => {
      // console.log(response.body);
      this.resourceDownload = response.body
      resource.totalDownload = this.resourceDownload.length
    }, (error: any) => {
      this.resourceDownload = []
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    }).add(() => {
      if (this.resourceDownload.length == 0) {
        this.isResourceDownload = false
      }
    })
  }

  // ======================================================== RESOURCE LIKE ========================================================

  createResourceLikeForm(): void {
    this.resourceLikeForm = this.formBuilder.group({
      id: new FormControl(null),
      resourceID: new FormControl(null, [Validators.required]),
      credentialID: new FormControl(null, [Validators.required]),
      isLiked: new FormControl(null),
    })
  }

  onLikeClick(resource: IResource): void {
    // console.log(resource)

    this.createResourceLikeForm()
    if (resource.resourceLike != null) {
      resource.resourceLike.isLiked = !resource.resourceLike.isLiked
      this.resourceLikeForm.patchValue(resource.resourceLike)
      this.updateResourceLike(resource)
      return
    }

    this.resourceLikeForm.get('resourceID').setValue(resource.id)
    this.resourceLikeForm.get('credentialID').setValue(this.loginID)
    this.resourceLikeForm.get('isLiked').setValue(true)
    this.addResourceLike(resource)
  }

  addResourceLike(resource: IResource): void {

    this.resourceService.addResourceLike(this.resourceLikeForm.value).subscribe((response: any) => {
      // console.log(response);
      this.getResourceLike(resource)
    }, (error: any) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  updateResourceLike(resource: IResource): void {

    this.resourceService.updateResourceLike(this.resourceLikeForm.value).subscribe((response: any) => {
      // console.log(response);
      this.getResourceLike(resource)
    }, (error: any) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  getResourceLike(resource: IResource): void {

    this.resourceService.getResourceLike(resource.id).subscribe((response: any) => {
      resource.resourceLike = response.body
      resource.totalLike = response.body.totalCount
      // console.log(resource)
    }, (error: any) => {
      console.error(error);
      resource.resourceLike.isLiked = !resource.resourceLike.isLiked
    })
  }

  getAllResourceLike(resource: IResource): void {

    this.spinnerService.loadingMessage = "Getting liked by details"

    this.resourceLike = []
    this.isResourceLike = true
    this.openScrollableModal(this.resourceLikeModal, "md")
    this.resourceService.getAllResourceLike(resource.id).subscribe((response: any) => {
      this.resourceLike = response.body
      resource.totalLike = this.resourceLike.length
    }, (error: any) => {
      this.resourceLike = []
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    }).add(() => {
      if (this.resourceLike.length == 0) {
        this.isResourceLike = false
      }
    })
  }


  // getResourcesByFileType(): void {
  //   if (this.isSearched) {
  //     this.getSearchedResources()
  //     return
  //   }
  //   // this.spinnerService.loadingMessage = "Getting all resources"
  //   
  //   this.totalResources = 0
  //   this.resources = []
  //   this.isResourceAvailable = true
  //   this.generalService.getResourcesByFileType(this.limit, this.offset, this.fileType).subscribe((response: any) => {
  //     this.resources = response.body
  //     this.totalResources = response.headers.get('X-Total-Count')
  //     if (this.totalResources == 0) {
  //       this.isResourceAvailable = false
  //     }
  //     
  //   }, (error: any) => {
  //     this.totalResources = 0
  //     this.isResourceAvailable = false
  //     console.error(error);
  //     if (error.statusText.includes('Unknown')) {
  //       alert("No connection to server. Check internet.")
  //     }
  //   }).add(() => {
  //     ;
  //     this.setPaginationString()
  //   })
  // }


  // onAudioSelect(event: any) {

  //   // if (this.resourceForm.get("resourceURL").value) {
  //   //   this.deletePreviousResource()
  //   // }

  //   this.docStatus = ""
  //   let files = event.target.files
  //   if (files && files.length) {
  //     let file = files[0]
  //     // check if file extension is valid
  //     let err = this.fileOperationService.isAudioFileValid(file)
  //     if (err != null) {
  //       this.docStatus = `<p><span>&#10060;</span> ${err}</p>`
  //       return
  //     }
  //     // Upload document if it is present.]
  //     this.uploadFile(file)
  //   }
  // }

  //On uplaoding document
  // onDocumentSelect(event: any) {

  //   // if (this.resourceForm.get("resourceURL").value) {
  //   //   this.deletePreviousResource()
  //   // }

  //   this.docStatus = ""
  //   let files = event.target.files
  //   if (files && files.length) {
  //     let file = files[0]
  //     // check if file extension is valid
  //     let err = this.fileOperationService.isDocumentFileValid(file)
  //     if (err != null) {
  //       this.docStatus = `<p><span>&#10060;</span> ${err}</p>`
  //       return
  //     }
  //     // Upload document if it is present.]
  //     this.uploadFile(file)
  //   }
  // }

  // onImageSelect(event: any) {

  //   // if (this.resourceForm.get("resourceURL").value) {
  //   //   this.deletePreviousResource()
  //   // }

  //   this.docStatus = ""
  //   let files = event.target.files
  //   if (files && files.length) {
  //     let file = files[0]
  //     // check if file extension is valid
  //     let err = this.fileOperationService.isImageFileValid(file)
  //     if (err != null) {
  //       this.docStatus = `<p><span>&#10060;</span> ${err}</p>`
  //       return
  //     }
  //     // Upload document if it is present.]
  //     this.uploadFile(file)
  //   }
  // }


}