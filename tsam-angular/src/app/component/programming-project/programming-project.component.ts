import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService, IProgrammingProject } from 'src/app/service/batch/batch.service';
import { UrlConstant } from 'src/app/service/constant';
import { GeneralService } from 'src/app/service/general/general.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { IResource, ResourceService } from 'src/app/service/resource/resource.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ITechnology, TechnologyService } from 'src/app/service/technology/technology.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-programming-project',
  templateUrl: './programming-project.component.html',
  styleUrls: ['./programming-project.component.css']
})
export class ProgrammingProjectComponent implements OnInit {

  // components
  resourceTypeList: any[]
  technologyList: ITechnology[]
  searchTechnologyList: ITechnology[]
  resourceList: IResource[]
  searchResourceList: IResource[]

  // tech component
  isTechLoading: boolean
  techLimit: number
  techOffset: number

  // batch-project
  programmingProjects: IProgrammingProject[]
  totalProgrammingProjects: number
  programmingProjectForm: FormGroup
  batchProjectTypes: string[]

  // search
  programmingProjectSearchForm: FormGroup
  isSearched: boolean
  searchFormValue: any

  // flags
  isOperationUpdate: boolean
  isViewMode: boolean
  isProgrammingProject: boolean

  // pagination
  limit: number
  currentPage: number
  offset: number
  paginationString: string

  // ck-editor
  ckConfig: any

  //project Document
  isDocumentUploadedToServer: boolean
  isFileUploading: boolean
  docStatus: string
  displayedFileName: string


  // modal
  modalRef: NgbModalRef
  @ViewChild('programmingProjectFormModal') programmingProjectFormModal: any
  @ViewChild('deleteModal') deleteModal: any
  @ViewChild('drawer') drawer: any

  // spinner

  isResourceLoading: boolean


  // permission
  permission: IPermission

  // constant
  readonly MAX_COMPLEXITY_LEVEL = 10
  private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]

  constructor(
    private formBuilder: FormBuilder,
    private utilService: UtilityService,
    private localService: LocalService,
    private urlConstant: UrlConstant,
    private batchService: BatchService,
    private generalService: GeneralService,
    private techService: TechnologyService,
    private resourceService: ResourceService,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private router: Router,
    private route: ActivatedRoute,
    private fileOperationService: FileOperationService,
  ) {
    this.editorConfig()
    this.initializeVariables()
    this.getAllComponents()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  editorConfig(): void {
    this.ckConfig = {
      extraPlugins: 'codeTag,kbdTag',
      removePlugins: "exportpdf",
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
        { name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
      ],
      toolbarGroups: [
        { name: 'styles' },
        { name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
        { name: 'document', groups: ['mode', 'document', 'doctools'] },
        { name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
        { name: 'links' },
      ],
      removeButtons: "",
      language: "en",
      resize_enabled: false,
      width: "100%",
      height: "80%",
      forcePasteAsPlainText: false,
    }
  }

  initializeVariables(): void {
    this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_PROJECT)

    this.limit = 5
    this.offset = 0
    this.totalProgrammingProjects = 0
    this.techLimit = 10
    this.techOffset = 0


    this.isOperationUpdate = false
    this.isViewMode = false
    this.isSearched = false
    this.isResourceLoading = false
    this.isProgrammingProject = true
    this.isTechLoading = false
    this.isDocumentUploadedToServer = false
    this.isFileUploading = false


    this.resourceList = []
    this.resourceTypeList = []
    this.technologyList = []
    this.programmingProjects = []
    this.searchResourceList = []
    this.searchTechnologyList = []
    this.batchProjectTypes = []
    this.displayedFileName = "Select file"

    this.createSearchForm()
  }

  getAllComponents(): void {
    this.getResourceType()
    this.getTechnologyList()
    this.getAllTechnologyList()
    this.searchOrGetProgrammingProject()
    this.getProjectType()
  }

  createSearchForm(): void {
    this.programmingProjectSearchForm = this.formBuilder.group({
      projectName: new FormControl(null, [Validators.maxLength(100)]),
      isActive: new FormControl(null),
      technologies: new FormControl(null),
      resourceType: new FormControl(null),
      resources: new FormControl(null),
      limit: new FormControl(this.limit),
      offset: new FormControl(this.offset),
    })
    this.programmingProjectSearchForm.get('resources').disable()
  }

  createProgrammingProjectForm(): void {
    const URLregex = /^(http(s)?:\/\/)[\w.-]+(\.[\w\.-]+)+[\w\-\._~:?#[\]@!\/$&'\(\)\*\+,;=.]+$/;

    this.programmingProjectForm = this.formBuilder.group({
      id: new FormControl(null),
      projectName: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
      description: new FormControl(null, [Validators.required, Validators.maxLength(250)]),
      code: new FormControl(null),
      isActive: new FormControl(true, [Validators.required]),
      complexityLevel: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(10)]),
      requiredHours: new FormControl(null, [Validators.required, Validators.min(1)]),
      sampleUrl: new FormControl(null, [Validators.maxLength(250), Validators.pattern(URLregex)]),
      resourceType: new FormControl(null),
      technologies: new FormControl(null, [Validators.required]),
      resources: new FormControl(null),
      document: new FormControl(null),
      projectType: new FormControl(null, [Validators.required]),
      score: new FormControl(null,[Validators.required, Validators.min(1), Validators.max(100)])
    })
    this.programmingProjectForm.get("resources").disable()
  }

  onAddClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createProgrammingProjectForm()
    this.openModal(this.programmingProjectFormModal, "xl")
  }

  onViewClick(courseProject: IProgrammingProject): void {
    this.isViewMode = true
    this.isOperationUpdate = false
    this.createProgrammingProjectForm()
    this.getResourceListByType(courseProject.resourceType)
    this.programmingProjectForm.disable()
    this.programmingProjectForm.patchValue(courseProject)
    if (courseProject.document) {
      this.displayedFileName = `<a href=${courseProject.document} target="_blank">Batch Project</a>`
    }
    this.openModal(this.programmingProjectFormModal, "xl")
  }

  onUpdateClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    if (this.programmingProjectForm.get("resourceType")?.value == null) {
      this.programmingProjectForm.get("resources").setValue(null)
    }
    this.programmingProjectForm.enable()
  }

  onDeleteClick(courseProjectID: string): void {
    this.openModal(this.deleteModal, 'md').result.then(() => {
      this.deleteProgrammingProject(courseProjectID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  setPaginationString() {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalProgrammingProjects < end) {
      end = this.totalProgrammingProjects
    }
    if (this.totalProgrammingProjects == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start}-${end}`
  }

  // Used to open modal.
  openModal(content: any, modalSize?: string): NgbModalRef {
    if (!modalSize) {
      modalSize = "lg"
    }
    this.resetDocumentUploadFields()
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: modalSize
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  resetSearchAndGetAll(): void {
    this.resetSearchForm()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([this.urlConstant.TRAINING_PROJECT])
  }

  resetSearchForm(): void {
    let limit = this.programmingProjectSearchForm.get("limit").value
    let offset = this.programmingProjectSearchForm.get("offset").value

    this.programmingProjectSearchForm.reset({
      limit: limit,
      offset: offset,
    })
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetProgrammingProject() {
    let queryParams = this.route.snapshot.queryParams
    if (!this.utilService.isObjectEmpty(queryParams)) {
      // this.getAllProgrammingProject()
      this.programmingProjectSearchForm.patchValue(queryParams)
    }
    console.log("searchOrGetProgrammingProject false");
    this.searchProgrammingProject()
  }

  searchAndCloseDrawer() {
    this.drawer.toggle()
    this.searchProgrammingProject()
  }

  searchProgrammingProject(): void {
    this.searchFormValue = { ...this.programmingProjectSearchForm?.value }
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: this.searchFormValue,
    })
    let flag: boolean = true

    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field];
      } else {
        if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
          this.isSearched = true
        }
        flag = false
      }
    }

    // No API call on empty search.
    if (flag) {
      return
    }
    this.getAllProgrammingProject()
  }

  // page change function
  changePage(pageNumber: number): void {
    // this.currentPage = pageNumber;
    // this.offset = this.currentPage - 1;
    // this.getAllProgrammingProject();

    this.programmingProjectSearchForm.get("offset").setValue(pageNumber - 1)
    this.searchProgrammingProject()
  }

  validateFormFields(): void {

    if (this.programmingProjectForm.get('resourceType')?.value != null) {
      this.programmingProjectForm.get('resources').setValidators([Validators.required])
    } else {
      this.resourceList = []
      this.programmingProjectForm.get('resources').clearValidators()
      this.programmingProjectForm.get('resources').disable()
    }
    this.programmingProjectForm.get('resources').updateValueAndValidity()
  }

  onSubmit(): void {
    console.log(this.programmingProjectForm.value);
    this.validateFormFields()

    if (this.programmingProjectForm.invalid) {
      this.programmingProjectForm.markAllAsTouched();
      return
    }
    if (this.isOperationUpdate) {
      this.updateProgrammingProject()
      return
    }
    this.addProgrammingProject()
  }

  // =============================================================== CRUD ===============================================================   

  getAllProgrammingProject(): void {
    this.spinnerService.loadingMessage = "Getting programming projects"

    this.totalProgrammingProjects = 0
    this.programmingProjects = []
    this.isProgrammingProject = true
    console.log("getAllProgrammingProject");
    this.batchService.getProgrammingProject(this.searchFormValue).subscribe((response: any) => {
      this.programmingProjects = response.body
      console.log(this.programmingProjects)
      this.totalProgrammingProjects = response.headers.get('X-Total-Count')
    }, (error) => {
      this.totalProgrammingProjects = 0
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error?.error)
    }).add(() => {

      if (this.totalProgrammingProjects == 0) {
        this.isProgrammingProject = false
      }
      this.setPaginationString()
    })
  }

  addProgrammingProject(): void {
    this.spinnerService.loadingMessage = "Adding programming project"
    this.batchService.addProgrammingProject(this.programmingProjectForm.value).subscribe((response: any) => {
      console.log(response)
      this.modalRef.close()
      alert("Project successfully added")
      this.getAllProgrammingProject()
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  updateProgrammingProject(): void {
    this.spinnerService.loadingMessage = "Updating programming project"
    console.log("this.programmingProjectForm.value",this.programmingProjectForm.value)
    this.batchService.updateProgrammingProject(this.programmingProjectForm.value).subscribe((response: any) => {
      console.log(response)
      this.modalRef.close()
      alert("Project successfully updated")
      this.getAllProgrammingProject()
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  deleteProgrammingProject(programmingProjectID: string): void {
    this.spinnerService.loadingMessage = "Deleting programming project"

    this.batchService.deleteProgrammingProject(programmingProjectID).subscribe((response: any) => {
      console.log(response)
      this.modalRef.close()
      alert("Project successfully deleted")
      this.getAllProgrammingProject()
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  // =============================================================== COMPONENTS ===============================================================   

  getResourceType(): void {
    this.generalService.getGeneralTypeByType("resource_type").subscribe((respond: any) => {
      this.resourceTypeList = respond;
    }, (err) => {
      console.error(err)
    })
  }

  getAllTechnologyList(): void {
    this.generalService.getTechnologies().subscribe((respond: any[]) => {
      this.searchTechnologyList = respond;
    }, (err) => {
      console.error(this.utilService.getErrorString(err))
    })
  }

  // Get Technology List.
  getTechnologyList(event?: any): void {
    let queryParams: any = {}
    if (event && event?.term != "") {
      queryParams.language = event?.term
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

  getResourceListByType(resourceType: string): void {
    this.isResourceLoading = true
    this.resourceList = []
    this.programmingProjectForm.get('resources').disable()

    let params: any = {
      resourceType: resourceType
    }

    this.resourceService.getResourcesList(params).subscribe((response: any) => {
      this.resourceList = response.body
      if (this.resourceList.length > 0 && !this.isViewMode) {
        this.programmingProjectForm.get('resources').enable()
      }
    }, (err: any) => {
      console.error(err)
    }).add(() => {
      this.isResourceLoading = false
    })
  }


  getSearchResourceListByType(resourceType: string): void {
    this.searchResourceList = []
    this.programmingProjectSearchForm.get('resources').setValue(null)
    this.programmingProjectSearchForm.get('resources').disable()

    let params: any = {
      resourceType: resourceType
    }

    this.resourceService.getResourcesList(params).subscribe((response: any) => {
      this.searchResourceList = response.body
      if (this.searchResourceList.length > 0) {
        this.programmingProjectSearchForm.get('resources').enable()
      }
    }, (err: any) => {
      console.error(err)
    }).add(() => {
      this.isResourceLoading = false
    })
  }

  getIntegerArray() {
    return Array.from(new Array(this.MAX_COMPLEXITY_LEVEL), (x, i) => i + 1)
  }

  //On uplaoding Document
  onResourceSelect(event: any) {
    this.docStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      let err = this.fileOperationService.isDocumentFileValid(file)
      if (err != null) {
        this.docStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // console.log(file)
      // Upload document if it is present.]
      this.isFileUploading = true
      this.fileOperationService.uploadBrochure(file,
        this.fileOperationService.BATCH_FOLDER + this.fileOperationService.PROJECT_FOLDER)
        .subscribe((data: any) => {
          this.programmingProjectForm.markAsDirty()
          this.programmingProjectForm.patchValue({
            document: data
          })
          if (file.name.toString().length > 25) {
            this.displayedFileName = file.name.toString().substr(0, 25) + "....."
          } else {
            this.displayedFileName = file.name
          }
          this.isFileUploading = false
          this.isDocumentUploadedToServer = true
          this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
        }, (error) => {
          this.isFileUploading = false
          this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
        })
    }
  }

  resetDocumentUploadFields(): void {
    this.isDocumentUploadedToServer = false
    this.isFileUploading = false
    this.displayedFileName = "Select file"
    this.docStatus = ""
  }

  dismissFormModal(modal: NgbModalRef) {
    if (this.isFileUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    if (this.isDocumentUploadedToServer) {
      if (!confirm("Uploaded document will be deleted.\nAre you sure you want to close?")) {
        return
      }
    }
    modal.dismiss()
    this.resetDocumentUploadFields()
  }
  // Get All Batch project type
  getProjectType(): void {
    this.generalService.getGeneralTypeByType("batch_project_type").subscribe((respond: any) => {
      // console.log(respond);
      this.batchProjectTypes = respond;

    }, (err) => {
      console.error(err.error.error)
    })
  }

}
