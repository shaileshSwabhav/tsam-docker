import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { GeneralService, IRole, ISearchFilterField } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-management',
  templateUrl: './admin-management.component.html',
  styleUrls: ['./admin-management.component.css']
})
export class AdminManagementComponent implements OnInit {

  // components
  allRoles: IRole[]
  allEmployees: IEmployee[]
  employeeList: IEmployee[]

  // flags

  // supervisor
  supervisorForm: FormGroup

  // search
  employeeSearchForm: FormGroup
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  // pagination
  limit: number;
  currentPage: number;
  offset: number;
  totalEmployees: number;
  paginationStart: number
  paginationEnd: number

  // modal
  modalRef: any;
  @ViewChild('supervisorModal') supervisorModal: any
  @ViewChild('deleteSupervisorModal') deleteSupervisorModal: any

  // spinner

  isSearched: boolean

  currentDate: Date

  // beta
  permission: IPermission

  constructor(
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private generalService: GeneralService,
    private localService: LocalService,
    private utilService: UtilityService,
    private urlConstant: UrlConstant,
    private router: Router,
    private route: ActivatedRoute,
    private roleConstant: Role
  ) {
    this.initializeVariables()
    this.getAllComponents()

  }

  initializeVariables() {
    this.allRoles = [] as IRole[]
    this.searchFormValue = []
    this.currentDate = new Date()

    // flags
    this.isSearched = false

    // employees limit offset
    this.limit = 5
    this.offset = 0

    // Search.
    this.searchFilterFieldList = []

    // initialize forms
    this.createSupervisorForm()
    this.createEmployeeSearchForm()

    this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'),
      this.urlConstant.ADMIN_ALL_EMPLOYEE)
  }

  getAllComponents() {
    this.getAllRoles()
    this.getEmployeeList()
    this.searchOrGetEmployees()

  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit() { }

  // Create form.
  createSupervisorForm(): void {
    this.supervisorForm = this.formBuilder.group({
      employeeCredentialID: new FormControl(null, [Validators.required]),
      supervisorCredentialID: new FormControl(null, [Validators.required])
    })
  }

  // Create search form.
  createEmployeeSearchForm(): void {
    this.employeeSearchForm = this.formBuilder.group({
      roleID: new FormControl(null),
      firstName: new FormControl(null, [Validators.maxLength(100)]),
      lastName: new FormControl(null, [Validators.maxLength(100)]),
      email: new FormControl(null, [Validators.maxLength(100)])
    });
  }

  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.employeeSearchForm.reset()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([this.urlConstant.ADMIN_ALL_EMPLOYEE])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.employeeSearchForm.reset()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetEmployees() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllEmployees()
      return
    }
    this.employeeSearchForm.patchValue(queryParams)
    this.searchEmployees()
  }

  searchEmployees(): void {
    this.searchFormValue = { ...this.employeeSearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Employees"
    this.changePage(1)
  }

  // Delete search criteria from employee search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.employeeSearchForm.get(searchName).setValue(null)
    this.searchEmployees()
  }

  getAllRoles(): void {
    this.generalService.getAllRoles().subscribe((data: any) => {
      this.allRoles = data
    }, (err) => {
      console.error(err)
    })
  }

  getEmployeeList(): void {
    this.generalService.getAllEmployeeList().subscribe((data: any) => {
      this.employeeList = data.body
    }, (err) => {
      console.error(err)
    })
  }

  // page change function
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber;
    this.offset = this.currentPage - 1;
    this.getAllEmployees();
  }

  onAddClick(employeeCredentialID: string): void {
    this.createSupervisorForm()
    this.supervisorForm.get('employeeCredentialID').setValue(employeeCredentialID)
    this.openModal(this.supervisorModal, 'mds')
  }

  onDeleteSupervisorClick(supervisorID: string, employeeID: string) {

    this.openModal(this.deleteSupervisorModal, 'md').result.then(() => {
      this.deleteSupervisor(supervisorID, employeeID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  // load all employees
  getAllEmployees(): void {
    this.spinnerService.loadingMessage = "Getting Employees"


    this.generalService.getAllEmployeeCredentials(this.limit, this.offset, this.searchFormValue).
      subscribe((data) => {
        this.totalEmployees = data.headers.get('X-Total-Count');
        this.allEmployees = data.body;
      }, (error) => {
        console.error(error);
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
        }
      }).add(() => {
        // Will work like finally.

        this.setPaginationString()
      })
  }


  // Add supervisor.
  addSupervisor(): void {
    this.spinnerService.loadingMessage = "Adding supervisor"

    this.generalService.addSupervisor(this.supervisorForm.value).subscribe(() => {
      this.modalRef.close();
      this.getAllEmployees();
      alert("Supervisor successfully added");
    }, (error) => {
      console.error(error);

      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  deleteSupervisor(supervisorID: string, employeeID: string): void {
    this.spinnerService.loadingMessage = "Deleting supervisor"

    this.generalService.deleteSupervisor(supervisorID, employeeID).
      subscribe(() => {
        this.modalRef.close();
        this.getAllEmployees();
        alert("Supervisor successfully deleted");
      }, (error) => {
        console.error(error);
        if (error.error?.error) {
          alert(error.error?.error);
          return;
        }
        alert(error.statusText);
      }).add(() => {
        // Will work like finally.

        this.setPaginationString()
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


  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef) {
    modal.dismiss()
  }

  onSubmit(): void {
    if (this.supervisorForm.invalid) {
      this.supervisorForm.markAllAsTouched();
      return
    }
    if (this.supervisorForm.get("employeeCredentialID").value === this.supervisorForm.get("supervisorCredentialID").value) {
      this.supervisorForm.get("supervisorCredentialID").setErrors({ "sameSupervisor": true })
      return
    }
    // for (let employee of )
    this.addSupervisor()
  }

  // Set total employees list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalEmployees < this.paginationEnd) {
      this.paginationEnd = this.totalEmployees
    }
  }
}

interface IEmployee {
  id?: string
  firstName: string
  lastName: string
  email?: string
  contact?: string
  role?: IRole
  supervisor?: IEmployee[]
}

interface ISupervisor {
  employeeCredentialID: string
  supervisorCredentialID: string
}