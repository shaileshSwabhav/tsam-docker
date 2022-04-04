import { Component, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IProject, ISearchFilterField, ProjectService } from 'src/app/service/project/project.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { UrlConstant } from 'src/app/service/constant';

@Component({
  selector: 'app-project',
  templateUrl: './project.component.html',
  styleUrls: ['./project.component.css']
})
export class ProjectComponent implements OnInit {

  // Flags.
  isSearched: boolean
  isOperationUpdate: boolean
  isViewMode: boolean
  showSubProjectsInForm: boolean

  // Project.
  projectList: IProject[]
  projectForm: FormGroup;

  // Pagination.
  limit: number
  currentPage: number
  totalProjects: number;
  offset: number
  paginationStart: number
  paginationEnd: number

  // Modal.
  modalRef: any
  @ViewChild('projectFormModal') projectFormModal: any
  @ViewChild('deleteProjectModal') deleteProjectModal: any

  // Spinner.


  // Search.
  projectSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  constructor(
    public utilityService: UtilityService,
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private projectService: ProjectService,
    private urlConstant: UrlConstant,
  ) {
    this.initializeVariables()
    this.searchOrGetProjects()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables
  initializeVariables(): void {

    // Components.
    this.projectList = [] as IProject[]

    // Flags.
    this.isOperationUpdate = false
    this.isViewMode = false
    this.isSearched = false
    this.showSubProjectsInForm = false

    // Pagination.
    this.limit = 5
    this.offset = 0
    this.currentPage = 0
    this.paginationStart = 0
    this.paginationEnd = 0

    // Initialize forms
    this.createProjectForm()
    this.createProjectSearchForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Projects"

    // Search.
    this.searchFilterFieldList = []
    this.searchFormValue = {}
  }

  // =============================================================CREATE FORMS==========================================================================
  // Create project form.
  createProjectForm() {
    this.projectForm = this.formBuilder.group({
      id: new FormControl(null),
      name: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z ]*$/), Validators.maxLength(100)]),
      isSubProjects: new FormControl(false),
      subProjects: this.formBuilder.array([]),
    });
  }

  // Add new project subProject to project form.
  addSubProject() {
    this.subProjectControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      projectID: new FormControl(null),
      name: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z ]*$/), Validators.maxLength(100)]),
    }));
  }

  // Create project search form.
  createProjectSearchForm() {
    this.projectSearchForm = this.formBuilder.group({
      projectName: new FormControl(null)
    })
  }

  // =============================================================PROJECT CRUD FUNCTIONS==========================================================================
  // On clicking add new project button.
  onAddNewProjectClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.showSubProjectsInForm = false
    this.createProjectForm()
    this.openModal(this.projectFormModal, 'md')
  }

  // Add new project.
  addProject(): void {
    this.spinnerService.loadingMessage = "Adding Project"

    this.projectService.addProject(this.projectForm.value).subscribe((response: any) => {
      this.modalRef.close()
      this.getAllProjects()
      alert(response)
    }, (error) => {
      console.error(error);

      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  // On clicking view project button.
  onViewProjectClick(project: IProject): void {
    this.isViewMode = true
    this.showSubProjectsInForm = false
    this.createProjectForm()
    this.projectForm.setControl("subProjects", this.formBuilder.array([]))
    if (project.subProjects == null) {
      project.subProjects = []
    }
    if (project.subProjects && project.subProjects.length > 0) {
      for (let i = 0; i < project.subProjects.length; i++) {
        this.addSubProject();
      }
      this.showSubProjectsInForm = true
      this.projectForm.get("isSubProjects").setValue(true)
    }
    this.projectForm.patchValue(project)
    this.projectForm.disable()
    this.openModal(this.projectFormModal, 'md')
  }

  // On cliking update form button in project form.
  onUpdateProjectClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.projectForm.enable()
  }

  // Update project.
  updateProject(): void {
    this.spinnerService.loadingMessage = "Updating Project"

    this.projectService.updateProject(this.projectForm.value).subscribe((response: any) => {

      this.modalRef.close()
      this.getAllProjects()
      alert(response)
    }, (error) => {
      console.error(error);

      if (error.error?.error) {
        alert(error.error.error);
        return;
      }
      alert("Check connection");
    })
  }

  // On clicking delete project button. 
  onDeleteProjectClick(projectID: string): void {
    this.openModal(this.deleteProjectModal, 'md').result.then(() => {
      this.deleteProject(projectID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete project after confirmation from user.
  deleteProject(projectID: string): void {
    this.spinnerService.loadingMessage = "Deleting Project"
    this.modalRef.close()

    this.projectService.deleteProject(projectID).subscribe((response: any) => {
      this.getAllProjects()
      alert(response)
    }, (error) => {
      console.error(error);

      if (error.error) {
        alert(error.error);
        return;
      }
      alert(error.statusText);
    })
  }

  // =============================================================PROJECT SEARCH FUNCTIONS==========================================================================
  // Reset search form and renaviagte page.
  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.projectSearchForm.reset()
    this.searchFormValue = {}
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([this.urlConstant.ADMIN_EMPLOYEE_PROJECT])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.projectSearchForm.reset()
  }

  // Search projects.
  searchProjects(): void {
    this.searchFormValue = { ...this.projectSearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Projects"
    this.changePage(1)
  }

  // Delete search criteria from project search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.projectSearchForm.get(searchName).setValue(null)
    this.searchProjects()
  }

  // ================================================OTHER FUNCTIONS FOR PROJECT===============================================

  // Page change.
  changePage(pageNumber: number): void {

    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getAllProjects()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetProjects() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilityService.isObjectEmpty(queryParams)) {
      this.getAllProjects()
      return
    }
    this.projectSearchForm.patchValue(queryParams)
    this.searchProjects()
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

  // On clicking sumbit button in project form.
  onFormSubmit(): void {
    if (this.projectForm.invalid) {
      this.projectForm.markAllAsTouched()
      return
    }
    if (this.isOperationUpdate) {
      this.updateProject()
      return
    }
    this.addProject()
  }

  // Set total feelings list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalProjects < this.paginationEnd) {
      this.paginationEnd = this.totalProjects
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

  // Get project subProjects array from project form.
  get subProjectControlArray(): FormArray {
    return this.projectForm.get("subProjects") as FormArray
  }

  // Delete subProject from project form.
  deleteSubProject(index: number) {
    if (confirm("Are you sure you want to delete sub project?")) {
      this.projectForm.markAsDirty();
      this.subProjectControlArray.removeAt(index);
    }
  }

  // On clicking isSubProjects checkbox, add or remove subProjects array.
  toggleSubProjectControls(event: any): void {
    if (this.showSubProjectsInForm) {
      if (confirm("Are you sure you want to delete all sub projects details?")) {
        this.showSubProjectsInForm = false;
        this.projectForm.setControl('subProjects', this.formBuilder.array([]));
        this.projectForm.get("isSubProjects").setValue(false)
        return;
      }
      event.target.checked = true
      return
    }
    this.showSubProjectsInForm = true;
    this.projectForm.get("isSubProjects").setValue(true)
    this.addSubProject();
  }

  // =============================================================GET FUNCTIONS==========================================================================
  // Get all projects by limit and offset.
  getAllProjects() {
    this.spinnerService.loadingMessage = "Getting All Projects"

    this.searchFormValue.limit = this.limit
    this.searchFormValue.offset = this.offset
    this.projectService.getAllProjects(this.searchFormValue).subscribe((response) => {

      this.totalProjects = response.headers.get('X-Total-Count')
      this.projectList = response.body
    }, error => {

      console.error(error);
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.setPaginationString()
    })
  }
}
