import { DatePipe } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IRole, AdminService } from 'src/app/service/admin/admin.service';
import { UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { IEmployee, GeneralService, ITechnologies, ISearchFilterField } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ICountry, IState } from 'src/app/service/talent/talent.service';
import { TechnologyService } from 'src/app/service/technology/technology.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-employee',
  templateUrl: './admin-employee.component.html',
  styleUrls: ['./admin-employee.component.css']
})
export class AdminEmployeeComponent implements OnInit {

  // components
  countryList: ICountry[]
  stateList: IState[]
  technologyList: ITechnologies[]

  // tech component
  isTechLoading: boolean
  techLimit: number
  techOffset: number

  //modal
  modalHeader: string;
  modalButton: string;
  modalRef: any;

  // form
  employeeForm: FormGroup
  employeeSearchForm: FormGroup

  // employee
  employees: IEmployee[]
  selectedEmployee: IEmployee
  totalEmployees: number
  employeeID: string

  // role list
  roleList: IRole[]

  //pagination
  limit: number;
  currentPage: number;
  offset: number;
  paginationStart: number
  paginationEnd: number

  // spinner


  // search
  searched: boolean
  searchFormValue: any
  searchFilterFieldList: ISearchFilterField[]

  // date
  currentDate: Date

  // resume
  isResumeUploadedToServer: boolean
  isFileUploading: boolean
  docStatus: string
  displayedFileName: string

  // access
  permission: IPermission

  // add/update/view
  isViewClicked: boolean
  isUpdateClicked: boolean

  // api buttons
  disableButton: boolean

  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private adminService: AdminService,
    private generalService: GeneralService,
    private techService: TechnologyService,
    private fileOpsService: FileOperationService,
    private urlConstant: UrlConstant,
    private modalService: NgbModal,
    private spinnerService: SpinnerService,
    private localService: LocalService,
    private datePipe: DatePipe,
    private activatedRoute: ActivatedRoute,
    private router: Router,
  ) {
    this.initialize()
    this.createForms()
    this.getAllComponents()
  }

  initialize(): void {
    this.permission = this.utilService.getPermission(this.urlConstant.ADMIN_OTHER_EMPLOYEE)

    this.spinnerService.loadingMessage = "Getting all employees"
    this.displayedFileName = "Select file"

    this.limit = 5
    this.offset = 0
    this.totalEmployees = 0
    this.techLimit = 10
    this.techOffset = 0

    this.isViewClicked = false
    this.isUpdateClicked = false
    this.disableButton = false
    this.isFileUploading = false
    this.searched = false
    this.isTechLoading = false

    this.employees = []
    this.countryList = []
    this.stateList = []

    // Search.
    this.searchFilterFieldList = []

    this.currentDate = new Date()
  }

  getAllComponents(): void {
    // this.getStateList()
    // this.getAllEmployees()
    this.getCountryList()
    this.getTechnologyList()
    this.getRoleList()
    this.getQueryParams()
  }

  createForms(): void {
    this.createEmployeeSearchForm()
    this.createEmployeeForm()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  createEmployeeSearchForm(): void {
    this.employeeSearchForm = this.formBuilder.group({
      firstName: new FormControl(null, [Validators.pattern("^[a-zA-Z]+$")]),
      lastName: new FormControl(null, [Validators.pattern("^[a-zA-Z]+$")]),
      email: new FormControl(null),
      contact: new FormControl(null),
      isActive: new FormControl('1'),
    })
  }

  createEmployeeForm(): void {
    this.employeeForm = this.formBuilder.group({
      id: new FormControl(null),
      code: new FormControl(null),
      firstName: new FormControl(null, [Validators.required, Validators.pattern("^[a-zA-Z]+$")]),
      lastName: new FormControl(null, [Validators.required, Validators.pattern("^[a-zA-Z]+$")]),
      contact: new FormControl(null, [Validators.required, Validators.pattern(/^[6789]\d{9}$/)]),
      dateOfBirth: new FormControl(null),
      dateOfJoining: new FormControl(null),
      email: new FormControl(null, [Validators.required, Validators.pattern(/^([a-zA-Z0-9._%+-]+@[a-zA-Z0-9]+\.[a-zA-z]{2,3})/)]),
      technologies: new FormControl(Array()),
      resume: new FormControl(null),
      isActive: new FormControl(null, [Validators.required]),
      type: new FormControl(null),
      role: new FormControl(null, [Validators.required]),
      address: new FormControl(null, [Validators.required, Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
      state: new FormControl(null, Validators.required),
      city: new FormControl(null, [Validators.required]), //, Validators.pattern(/^[a-zA-Z\s*]+[a-zA-Z]+$/)
      country: new FormControl(null, Validators.required),
      pinCode: new FormControl(null, [Validators.required, Validators.pattern(/^[0-9]{6}$/)]),
    })
    this.employeeForm.get('state')?.disable()
  }

  onAddEmployeeClick(modalContent: any): void {
    this.isViewClicked = false
    this.isUpdateClicked = false
    this.isFileUploading = false
    this.createEmployeeForm()
    this.modalHeader = "Add New Employee"
    this.modalButton = "Add Employee"
    this.openModal(modalContent, "xl")
  }

  onViewEmployeeClick(employee: IEmployee, modalContent: any): void {
    this.isViewClicked = true
    this.isUpdateClicked = false
    this.isFileUploading = false
    this.modalHeader = "Employee Details"
    this.modalButton = ""
    this.selectedEmployee = employee
    this.updateEmployeeForm()
    this.openModal(modalContent, "xl")
    this.employeeForm.disable()
  }

  onUpdateEmployeeClick(): void {
    this.isViewClicked = false
    this.isUpdateClicked = true
    this.isFileUploading = false
    this.modalHeader = "Update Employee"
    this.modalButton = "Update Employee"
    this.employeeForm.enable()
    this.employeeForm.get('code').disable()
  }

  onDeleteEmployeeClick(employeeID: string, modalContent: any): void {
    this.employeeID = employeeID
    this.openModal(modalContent, "md")
  }

  updateEmployeeForm(): void {
    this.createEmployeeForm()
    this.displayedFileName = "No resume uploaded"
    if (this.selectedEmployee.resume) {
      this.displayedFileName = `<a href=${this.selectedEmployee.resume} target="_blank">Resume present</a>`
    }
    this.employeeForm.patchValue(this.selectedEmployee)
  }

  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.searched = false
    this.resetSearchForm()
    this.changePage(1)
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.employeeSearchForm.reset()
    this.employeeSearchForm.get("isActive").setValue("1")
  }

  changePage($event: any): void {

    this.offset = $event - 1
    this.currentPage = $event;
    this.getAllEmployees()
  }

  // Compare Ob1 and Ob2
  compareFn(ob1: any, ob2: any): boolean {
    if (ob1 == null && ob2 == null) {
      return true;
    }
    return ob1 && ob2 ? ob1.id === ob2.id : ob1 === ob2
  }

  // Set total users list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalEmployees < this.paginationEnd) {
      this.paginationEnd = this.totalEmployees
    }
  }

  showSpecificStates(employee: any, state: IState): boolean {
    if (employee.get('country').value == null) {
      return false
    }
    if (employee.get('country').value.id === state.countryID) {
      return true
    }
    return false
  }

  searchEmployee(): void {
    this.spinnerService.loadingMessage = "Searching employee";
    this.searchFormValue = { ...this.employeeSearchForm.value }
    let flag: boolean = true
    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field];
      } else {
        this.searched = true
        flag = false
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
    // No API call on empty search.
    if (flag) {
      return
    }
    this.router.navigate([], {
      relativeTo: this.activatedRoute,
      queryParams: this.searchFormValue,
    })
    console.log(this.searchFormValue)
    this.searched = true
    this.changePage(1)
  }

  // Delete search criteria from designation search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.employeeSearchForm.get(searchName).setValue(null)
    this.searchEmployee()
  }

  openModal(modalContent: any, modalSize?: string): void {
    if (modalSize == undefined) {
      modalSize = 'lg'
    }
    this.modalRef = this.modalService.open(modalContent,
      {
        ariaLabelledBy: 'modal-basic-title',
        backdrop: 'static', size: modalSize,
        keyboard: false
      }
    );
    /*this.modalRef.result.subscribe((result) => {
    }, (reason) => {

    });*/
  }

  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef) {
    if (this.isFileUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    if (this.isResumeUploadedToServer) {
      if (!confirm("Uploaded resume will be deleted.\nAre you sure you want to close?")) {
        return
      }
      // this.deleteResume()
    }
    modal.dismiss()
    this.resetUploadFields()
    // this.isResumeUploadedToServer = false
    // this.displayedFileName = "Select file"
    // this.docStatus = ""
  }

  getQueryParams(): void {
    this.activatedRoute.queryParams.subscribe((param) => {
      // console.log(param)
      this.searchFormValue = { ...param }
      for (let field in this.searchFormValue) {
        if (!this.searchFormValue[field]) {
          delete this.searchFormValue[field];
        } else {
          this.searched = true
        }
      }
      this.employeeSearchForm.patchValue(this.searchFormValue)
      this.changePage(1)
    })
  }

  // =============================================================CRUD=============================================================

  getAllEmployees(): void {
    if (this.searched) {
      this.getSearchedEmployees()
      return
    }
    this.spinnerService.loadingMessage = "Getting employees"

    this.totalEmployees = 1
    this.employees = []
    this.generalService.getAllEmployees(this.limit, this.offset).subscribe((response: any) => {
      this.employees = response.body
      this.totalEmployees = response.headers.get("X-Total-Count")
      this.setPaginationString()

      for (let index in this.employees) {
        let dateOfBirth = this.employees[index].dateOfBirth
        let dateOfJoining = this.employees[index].dateOfJoining
        if (!dateOfBirth && !dateOfJoining) {
          continue
        }
        if (dateOfBirth) {
          this.employees[index].dateOfBirth = this.datePipe.transform(dateOfBirth, 'yyyy-MM-dd')
        }
        if (dateOfJoining) {
          this.employees[index].dateOfJoining = this.datePipe.transform(dateOfJoining, 'yyyy-MM-dd')
        }
      }


    }, (error: any) => {
      console.error(error);
      this.totalEmployees = 0
      this.setPaginationString()
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    })
  }

  getSearchedEmployees(): void {
    this.spinnerService.loadingMessage = "Getting searched employees"

    this.totalEmployees = 1
    this.employees = []
    this.generalService.getAllEmployees(this.limit, this.offset, this.searchFormValue).subscribe((response: any) => {
      this.employees = response.body
      this.totalEmployees = response.headers.get("X-Total-Count")
      this.setPaginationString()

      for (let index in this.employees) {
        let dateOfBirth = this.employees[index].dateOfBirth
        let dateOfJoining = this.employees[index].dateOfJoining
        if (!dateOfBirth && !dateOfJoining) {
          continue
        }
        if (dateOfBirth) {
          this.employees[index].dateOfBirth = this.datePipe.transform(dateOfBirth, 'yyyy-MM-dd')
        }
        if (dateOfJoining) {
          this.employees[index].dateOfJoining = this.datePipe.transform(dateOfJoining, 'yyyy-MM-dd')
        }
      }


    }, (error: any) => {
      console.error(error);
      this.totalEmployees = 0
      this.setPaginationString()
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    })
  }

  addEmployee(): void {
    this.spinnerService.loadingMessage = "Adding Employee"

    this.generalService.addEmployee(this.employeeForm.value).subscribe((response: any) => {
      alert("Employee successfully added")
      this.resetUploadFields()
      this.modalRef.close()

      this.getAllEmployees()
    }, (error: any) => {
      console.error(error);

      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    })
  }

  updateEmployee(): void {
    this.spinnerService.loadingMessage = "Updating Employee"

    this.generalService.updateEmployee(this.employeeForm.value).subscribe((response: any) => {
      alert("Employee successfully updated")
      this.resetUploadFields()
      this.modalRef.close()

      this.getAllEmployees()
    }, (error: any) => {
      console.error(error);

      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    })
  }

  deleteEmployee(): void {
    this.spinnerService.loadingMessage = "Deleting Employee"

    this.generalService.deleteEmployee(this.employeeID).subscribe((response: any) => {
      this.modalRef.close()
      alert("Employee successfully deleted")

      this.getAllEmployees()
    }, (error: any) => {
      console.error(error);

      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    })
  }

  validate(): void {
    console.log(this.employeeForm.controls);

    if (this.employeeForm.invalid) {
      this.employeeForm.markAllAsTouched()
      return
    }

    if (this.isUpdateClicked) {
      this.updateEmployee()
      return
    }

    this.addEmployee()
  }

  //On uplaoding resume
  onResourceSelect(event: any) {
    this.docStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      // Upload resume if it is present.]
      this.isFileUploading = true
      this.fileOpsService.uploadResume(file).subscribe((data: any) => {
        this.employeeForm.markAsDirty()
        this.employeeForm.patchValue({
          resume: data
        })
        this.displayedFileName = file.name
        this.isFileUploading = false
        this.isResumeUploadedToServer = true
        this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
      }, (error) => {
        this.isFileUploading = false
        this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
      })

    }
  }

  // =============================================================COMPONENTS=============================================================

  // Get Country List
  getCountryList(): void {
    this.generalService.getCountries().subscribe((respond: any[]) => {
      this.countryList = respond;
    }, (err) => {
      console.error(this.utilService.getErrorString(err))
    })
  }

  // getStateList(): void {
  //   this.generalService.getStates().subscribe((data: IState[]) => {
  //         this.stateList = data.sort();
  //   }, (err) => {
  //         console.error(err)
  //   })
  // }

  getStateList(countryid: any): void {
    this.employeeForm.get('state')?.reset()
    this.employeeForm.get('state')?.disable()
    this.stateList = [] as IState[]
    this.generalService.getStatesByCountryID(countryid).subscribe((respond: any[]) => {
      this.stateList = respond
      if (this.stateList.length > 0 && !this.isViewClicked) {
        this.employeeForm.get('state')?.enable()
      }
    }, (err) => {
      console.error(this.utilService.getErrorString(err))
    })
  }

  getRoleList(): void {
    this.adminService.getAllRole().subscribe((response: any) => {
      this.roleList = response.body
    }, (err: any) => {
      console.error(err);
    })
  }

  // Get Technology List.
  getTechnologyList(event?: any): void {
    // this.generalService.getTechnologies().subscribe((respond: any[]) => {
    //   this.technologyList = respond;
    // }, (err) => {
    //   console.error(this.utilService.getErrorString(err))
    // })
    // console.log(event);
    let queryParams: any = {}
    if (event && event?.term != "") {
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

  resetUploadFields(): void {
    this.displayedFileName = ""
    this.docStatus = ""
    this.isResumeUploadedToServer = false
    this.isFileUploading = false
  }

}
