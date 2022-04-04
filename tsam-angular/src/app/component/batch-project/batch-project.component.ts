import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService, IBatchProject, IProgrammingProject } from 'src/app/service/batch/batch.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { IResource } from 'src/app/service/course/course.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { ResourceService } from 'src/app/service/resource/resource.service';
import { ITechnology, TechnologyService } from 'src/app/service/technology/technology.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-batch-project',
  templateUrl: './batch-project.component.html',
  styleUrls: ['./batch-project.component.css']
})
export class BatchProjectComponent implements OnInit {

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

  // batch
  batchID: string
  batchName: string

  // programming-project
  programmingProjects: IProgrammingProject[]
  totalProgrammingProjects: number
  programmingProjectForm: FormGroup

  // batch-project
  batchProjectForm: FormGroup
  batchProjectSearchForm: FormGroup
  multipleBatchProjectsForm: FormGroup
  programmingProjectIDList: string[]
  batchProjects: IBatchProject[]
  totalBatchProjects: number
  allBatchProjectIDs: string[]
  batchProjectTypes: string[]

  // search
  programmingProjectSearchForm: FormGroup
  isSearched: boolean
  searchFormValue: any
  // searchProject: Subject<IProgrammingProject>

  // flags
  isOperationUpdate: boolean
  isViewMode: boolean
  isProgrammingProject: boolean

  // pagination
  limit: number
  currentPage: number
  offset: number
  programmingPaginationString: string
  batchPaginationString: string

  // ck-editor
  ckConfig: any

  //project Document
  isDocumentUploadedToServer: boolean
  isFileUploading: boolean
  docStatus: string
  displayedFileName: string

  //Project Publish Form
  projectPublishForm: FormGroup
  publishProject: IBatchProject

  // modal
  modalRef: NgbModalRef
  @ViewChild('programmingProjectFormModal') programmingProjectFormModal: any
  @ViewChild('deleteModal') deleteModal: any
  @ViewChild('drawer') drawer: any
  @ViewChild('publishProjectFormModal') publishProjectFormModal: any

  isResourceLoading: boolean


  // permission
  permission: IPermission

  // Input
  @Input() isAssignClick: boolean
  @Input() isManageClick: boolean


  // constant
  readonly MAX_COMPLEXITY_LEVEL = 10
  private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]

  constructor(
    private formBuilder: FormBuilder,
    private utilService: UtilityService,
    private urlConstant: UrlConstant,
    private batchService: BatchService,
    private generalService: GeneralService,
    private techService: TechnologyService,
    private resourceService: ResourceService,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private route: ActivatedRoute,
    private fileOperationService: FileOperationService,
    private localService: LocalService,
		private role: Role,
    private el: ElementRef,
  ) {
    this.initializeVariables()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    this.extractQueryParams()
    this.editorConfig()
    this.getAllComponents()
  }

  extractQueryParams(): void {
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      this.batchName = params.get("batchName")
    }, err => {
      console.error(err);
    })
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
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN || this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER_SESSION_DETAILS)
		}
		if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
			this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH_SESSION_DETAILS)
		}

    this.limit = 10
    this.offset = 0
    this.totalProgrammingProjects = 0
    this.techLimit = 10
    this.techOffset = 0

    this.totalBatchProjects = 0

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
    this.batchProjects = []
    this.programmingProjectIDList = []
    this.allBatchProjectIDs = []
    this.batchProjectTypes = []


    this.displayedFileName = "Select file"

    // this.searchProject = new Subject<IProgrammingProject>()

    this.createProgrammingProjectSearchForm()
    this.createBatchProjectSearchForm()
    this.batchProjectForm = this.createBatchProjectForm()
  }

  getAllComponents(): void {
    if (this.isAssignClick) {
      this.getTechnologyList()
      this.getAllTechnologyList()
      this.searchOrGetProgrammingProject()
    }
    this.getResourceType()
    this.getProjectType()
    this.getBatchProjects()
  }

  createProgrammingProjectSearchForm(): void {
    this.programmingProjectSearchForm = this.formBuilder.group({
      projectName: new FormControl(null, [Validators.maxLength(100)]),
      limit: new FormControl(this.limit),
      offset: new FormControl(this.offset),
    })
  }

  createBatchProjectSearchForm(): void {
    this.batchProjectSearchForm = this.formBuilder.group({
      limit: new FormControl(this.limit),
      offset: new FormControl(this.offset),
    })
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

  createMultipleBatchProjectsForm(): void {
    this.multipleBatchProjectsForm = this.formBuilder.group({
      batchProjects: new FormArray([]),
    })
  }

  get batchProjectsControlArray(): FormArray {
    return this.multipleBatchProjectsForm.get("batchProjects") as FormArray
  }

  addBatchProjectsFormToControlArray(): void {
    this.batchProjectsControlArray.push(this.createBatchProjectForm())
  }

  createBatchProjectForm(): FormGroup {
    return this.formBuilder.group({
      id: new FormControl(null),
      batchID: new FormControl(null),
      programmingProjectID: new FormControl(null),
    })
  }

  onAddClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createProgrammingProjectForm()
    this.openModal(this.programmingProjectFormModal, "xl")
  }

  onBatchProjectDeleteClick(project: IBatchProject): void {
    // console.log(project);
    this.openModal(this.deleteModal, 'md').result.then(() => {
      this.deleteBatchProject(project.id)
    }, (err) => {
      console.error(err);
      return
    })
  }

  onViewClick(courseProject: IProgrammingProject): void {
    console.log("courseProject",courseProject)
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

  setProgrammingPaginationString(): void {
    this.programmingPaginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalProgrammingProjects < end) {
      end = this.totalProgrammingProjects
    }
    if (this.totalProgrammingProjects == 0) {
      this.programmingPaginationString = ''
      return
    }
    this.programmingPaginationString = `${start}-${end}`
  }

  setBatchPaginationString(): void {
    this.batchPaginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalBatchProjects < end) {
      end = this.totalBatchProjects
    }
    if (this.totalBatchProjects == 0) {
      this.batchPaginationString = ''
      return
    }
    this.batchPaginationString = `${start}-${end}`
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
    // this.router.navigate(['/batch/project'])
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
  searchOrGetProgrammingProject(): void {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllProgrammingProject()
      return
    }
    this.programmingProjectSearchForm.patchValue(queryParams)
    this.searchProgrammingProject()
  }

  searchProgrammingProject(): void {
    this.searchFormValue = { ...this.programmingProjectSearchForm?.value }

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

    // this.getAllProgrammingProject()
    this.getProgrammingProject()
  }

  searchBatchProject(): void {
    this.searchFormValue = { ...this.batchProjectSearchForm?.value }

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

    this.getBatchProjects()
  }

  // page change function
  changePage(pageNumber: number): void {
    this.programmingProjectSearchForm.get("offset").setValue(pageNumber - 1)
    this.searchProgrammingProject()
  }

  changeBatchProjectPage(pageNumber: number): void {
    this.batchProjectSearchForm.get("offset").setValue(pageNumber - 1)
    this.searchBatchProject()
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

  validateProgrammingProject(): void {
    // console.log(this.programmingProjectForm.controls);
    this.validateFormFields()

    if (this.programmingProjectForm.invalid) {
      this.programmingProjectForm.markAllAsTouched();
      // this.scrollToFirstInvalidControl();
      return
    }

    this.addProgrammingProject()
  }

  // =============================================================== CRUD ===============================================================   

  // getSearchProject(): void {
  //   this.searchProject.pipe(
  //     map((value) => {
  //       console.log(value);
  //     }),
  //     debounceTime(500),
  //     distinctUntilChanged(),
  //   ).subscribe
  // }

  getProgrammingProject(): void {
    this.totalProgrammingProjects = 0
    this.programmingProjects = []
    this.isProgrammingProject = true

    // .pipe(
    //   debounceTime(500),
    //   distinctUntilChanged()
    // )
    this.batchService.getProgrammingProject(this.searchFormValue).subscribe((response) => {
      this.programmingProjects = response.body
      // console.log("programmingProjects",this.programmingProjects)
      this.totalProgrammingProjects = response.headers.get('X-Total-Count')
    }, (err: any) => {
      this.totalProgrammingProjects = 0
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error?.error)
    }).add(() => {

      if (this.totalProgrammingProjects == 0) {
        this.isProgrammingProject = false
      }
      this.setProgrammingPaginationString()
    })
  }

  getAllProgrammingProject(): void {
    this.spinnerService.loadingMessage = "Getting programming projects"

    this.totalProgrammingProjects = 0
    this.programmingProjects = []
    this.isProgrammingProject = true
    this.batchService.getProgrammingProject(this.searchFormValue).subscribe((response: any) => {
      this.programmingProjects = response.body
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
      this.setProgrammingPaginationString()
    })
  }

  addProgrammingProject(): void {
    this.spinnerService.loadingMessage = "Adding programming project"

    this.batchService.addProgrammingProject(this.programmingProjectForm.value).subscribe((response: any) => {
      // console.log(response)
      this.modalRef.close()

      this.createNewBatchProject(response.body)
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

  createNewBatchProject(programmingProjectID: string): void {
    console.log("programmingProjectID",programmingProjectID)
    this.batchProjectForm = this.createBatchProjectForm()
    this.batchProjectForm.get("batchID").setValue(this.batchID)
    this.batchProjectForm.get("programmingProjectID").setValue(programmingProjectID)
    this.addBatchProject(this.batchProjectForm.value, true)
  }

  addProgrammingProjectID(): void {
    this.programmingProjectIDList = []
    this.allBatchProjectIDs = []

    for (let index = 0; index < this.batchProjects.length; index++) {
      this.programmingProjectIDList.push(this.batchProjects[index].programmingProject.id)
      this.allBatchProjectIDs.push(this.batchProjects[index].programmingProject.id)
      this.addBatchProjectToForm(index)
    }
    // console.log("this.batchProjectsControlArray.dirty -> ", this.batchProjectsControlArray.dirty);
    // console.log("addProgrammingProjectID -> ", this.batchProjectsControlArray.value);
  }

  addBatchProjectToForm(index: number): void {
    this.addBatchProjectsFormToControlArray()
    // this.batchProjectsControlArray.at(index).patchValue(this.batchProjects[index])
    this.batchProjectsControlArray.at(index).get("id").setValue(this.batchProjects[index].id)
    this.batchProjectsControlArray.at(index).get("batchID").setValue(this.batchID)
    this.batchProjectsControlArray.at(index).get("programmingProjectID").setValue(this.batchProjects[index].programmingProject.id)
  }

  toggleBatchProjects(programmingProjectID: string): void {
    // console.log(this.allBatchProjectIDs);
    this.batchProjectsControlArray.markAsDirty()

    if (this.programmingProjectIDList.includes(programmingProjectID)) {
      let index = this.programmingProjectIDList.indexOf(programmingProjectID)
      this.programmingProjectIDList.splice(index, 1)
      this.batchProjectsControlArray.removeAt(index)

      // console.log("toggleBatchProjects -> ", this.batchProjectsControlArray.value);
      return
    }

    // add
    if (!this.programmingProjectIDList.includes(programmingProjectID)) {
      this.programmingProjectIDList.push(programmingProjectID)
      let len = this.batchProjectsControlArray.length

      this.addBatchProjectsFormToControlArray()
      this.batchProjectsControlArray.at(len).get("batchID").setValue(this.batchID)
      this.batchProjectsControlArray.at(len).get("programmingProjectID").setValue(programmingProjectID)

      this.batchProjects.find((batchProject: IBatchProject) => {
        if (batchProject.programmingProject.id == programmingProjectID) {
          this.batchProjectsControlArray.at(len).get("id").setValue(batchProject.id)
        }
      })
    }
    // console.log("toggleBatchProjects -> ", this.batchProjectsControlArray.value);
  }

  validateBatchProject(): void {
    // console.log(this.multipleBatchProjectsForm.controls);

    if (this.multipleBatchProjectsForm.invalid) {
      this.multipleBatchProjectsForm.markAllAsTouched();
      return
    }

    this.submitBatchProject()
  }

  submitBatchProject(): void {
    let err: string[] = []
    let batchProjects: IBatchProject[] = this.multipleBatchProjectsForm.value.batchProjects
    console.log("batchProjects -> ", batchProjects);

    for (let index = 0; index < batchProjects.length; index++) {
      if (batchProjects[index].id) {
        this.updateBatchProject(batchProjects[index], err)

        // programming-project
        let j = this.programmingProjectIDList.indexOf(batchProjects[index].programmingProjectID)
        this.programmingProjectIDList.splice(j, 1)

        // batch-project
        j = this.allBatchProjectIDs.indexOf(batchProjects[index].programmingProjectID)
        this.allBatchProjectIDs.splice(j, 1)
        continue
      }
      this.addBatchProject(batchProjects[index], false, err)

      // programming-project
      let j = this.programmingProjectIDList.indexOf(batchProjects[index].programmingProjectID)
      this.programmingProjectIDList.splice(j, 1)

      // batch-project
      j = this.allBatchProjectIDs.indexOf(batchProjects[index].programmingProjectID)
      this.allBatchProjectIDs.splice(j, 1)
    }
  }

  getBatchProjects(): void {
    this.spinnerService.loadingMessage = "Getting batch projects"

    this.batchProjects = []
    this.totalBatchProjects = 0
    this.batchService.getBatchProject(this.batchID, this.searchFormValue).subscribe((response: any) => {
      // console.log(response);
      this.batchProjects = response.body
      console.log("batchProjects",this.batchProjects)
      this.totalBatchProjects = response.headers.get('X-Total-Count')
      this.createMultipleBatchProjectsForm()
      this.addProgrammingProjectID()
    }, (err: any) => {
      console.error(err);
      this.totalBatchProjects = 0
      if (err.error?.error) {
        alert(err.error?.error);
        return;
      }
      alert(err.statusText);
    }).add(() => {
      this.batchProjectSearchForm.get("limit").enable()
      if (this.totalBatchProjects == 0) {
        this.batchProjectSearchForm.get("limit").disable()
      }
      this.setBatchPaginationString()

    })
  }

  addBatchProject(batchProject: IBatchProject, isAutoAdd: boolean = false, errors: string[] = []): void {
    this.spinnerService.loadingMessage = "Updating batch project"

    console.log("batchProject",batchProject);
    this.batchService.addBatchProject(this.batchID, batchProject).subscribe((response: any) => {
      // console.log(response);
    }, (err: any) => {
      console.error(err);
      errors.push(err)
    }).add(() => {

      if (isAutoAdd) {
        this.getBatchProjects()
      }
      if (this.ongoingOperations == 0) {
        if (errors.length == 0) {
          if (this.allBatchProjectIDs.length > 0) {
            this.deleteBatchProjectFromList(errors)
            return
          }
          alert("Projects successfully updated in batch")
          this.getBatchProjects()
        } else {
          let errorString = ""
          for (let index = 0; index < errors.length; index++) {
            errorString += (index == 0 ? "" : "\n") + errors[index]
          }
          alert(errorString)
        }
      }
    })
  }

  updateBatchProject(batchProject: IBatchProject, errors: string[] = []): void {
    this.spinnerService.loadingMessage = "Updating batch project"
    console.log("updateBatchProject",batchProject)
    this.batchService.updateBatchProject(this.batchID, batchProject).subscribe((response: any) => {
      // console.log(response);
    }, (err: any) => {
      console.error(err);
      errors.push(err)
    }).add(() => {

      if (this.ongoingOperations == 0) {
        if (errors.length == 0) {
          if (this.allBatchProjectIDs.length > 0) {
            this.deleteBatchProjectFromList(errors)
            return
          }
          alert("Projects successfully updated in batch")
          this.getBatchProjects()
        } else {
          let errorString = ""
          for (let index = 0; index < errors.length; index++) {
            errorString += (index == 0 ? "" : "\n") + errors[index]
          }
          alert(errorString)
        }
      }
    })
  }

  deleteBatchProjectFromList(err: string[]): void {
    // console.log("deleteBatchProjectFromList -> ", this.allBatchProjectIDs);
    for (let index = 0; index < this.allBatchProjectIDs?.length; index++) {
      this.batchProjects.find((value: IBatchProject) => {
        if (value.programmingProject.id == this.allBatchProjectIDs[index]) {
          this.deleteBatchProject(value.id, err)
          // return
        }
      })
    }
  }

  deleteBatchProject(batchProjectID: string, err?: string[]): void {
    this.spinnerService.loadingMessage = "Updating batch project"

    this.batchService.deleteBatchProject(this.batchID, batchProjectID).subscribe((response: any) => {
      // console.log(response.body);
      if (err == null) {
        alert("Assignments successufully updated")
        this.getBatchProjects()
      }
    }, (err: any) => {
      console.error(err);
      if (err.error?.error) {
        alert(err.error?.error);
        return;
      }
      alert(err.statusText);
    }).add(() => {

      if (this.ongoingOperations == 0) {
        if (err?.length == 0) {
          alert("Assignments successufully updated")
          this.getBatchProjects()
          return
        }
        let errorString = ""
        for (let index = 0; index < err?.length; index++) {
          errorString += (index == 0 ? "" : "\n") + err[index]
        }
        alert(errorString)
      }
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
      // console.log("batchProjectTypes", respond);
      this.batchProjectTypes = respond;

    }, (err) => {
      console.error(err.error.error)
    })
  }

  // Publish Project
  onPublishClick( batchProject: IBatchProject): void{
    this.publishProject = batchProject
    this.createPublishProjectForm()
    this.openModal(this.publishProjectFormModal, "md")
  }

  //Create Publish Project Form
  createPublishProjectForm(): void{
    this.projectPublishForm = this.formBuilder.group({
      id : new FormControl(null),
      dueDate: new FormControl(null,[Validators.required])
    })
  }

  validatePublishProjectForm(): void{

    if (this.projectPublishForm.invalid) {
      this.projectPublishForm.markAllAsTouched()
      return
    }
    this.publishProject.dueDate = this.projectPublishForm.get("dueDate").value
    this.publishProject.batchID= this.batchID
    this.publishProject.programmingProjectID= this.publishProject.programmingProject.id

    console.log("this.publishProject",this.publishProject,this.projectPublishForm.get("dueDate").value)
    this.publishProjectDate(this.publishProject)

  }

  publishProjectDate(batchProject: IBatchProject): void{
    this.spinnerService.loadingMessage = "Updating batch project"

    this.batchService.updateBatchProject(this.batchID, batchProject).subscribe((response: any) => {
      this.modalRef.close()
			alert("Projects successfully updated in batch")
      this.getBatchProjects()
		}, (error) => {
			if (error.error?.error) {
				alert(error.error.error)
				return
			}
			alert("Check connection")
		}).add(() => {

		})

  }

  // private scrollToFirstInvalidControl() {
  //   const firstInvalidControl: HTMLElement = this.el.nativeElement.querySelector(
  //     "programmingProjectForm .ng-invalid"
  //   );
  //   console.log("firstInvalidControl",firstInvalidControl)
  //   firstInvalidControl.focus(); //without smooth behavior
  // }
  

}
