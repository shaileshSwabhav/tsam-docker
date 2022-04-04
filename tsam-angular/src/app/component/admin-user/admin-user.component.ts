import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService, ICountry, IState, IUser, IRole, ISearchFilterField } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-admin-user',
  templateUrl: './admin-user.component.html',
  styleUrls: ['./admin-user.component.css']
})
export class AdminUserComponent implements OnInit {

  // components
  countryList: ICountry[]
  stateList: IState[]
  allRoles: IRole[]

  // flags
  isOperationUpdate: boolean
  isViewMode: boolean

  // user
  salesPersonList: IUser[]
  allUsers: IUser[]
  userForm: FormGroup

  // search
  userSearchForm: FormGroup
  searchFormValue: any
  isSearched: boolean
  searchFilterFieldList: ISearchFilterField[]

  // pagination
  limit: number;
  currentPage: number;
  offset: number;
  totalUsers: number;
  paginationStart: number
  paginationEnd: number

  // modal
  modalRef: NgbModalRef;
  @ViewChild('userFormModal') userFormModal: any
  @ViewChild('deleteConfirmationModal') deleteConfirmationModal: any

  // resume
  isResumeUploadedToServer: boolean
  isFileUploading: boolean
  docStatus: string
  displayedFileName: string


  // spinner

  currentDate: Date

  // beta
  permission: IPermission
  readonly OPERATIONS = { 'add': "Add User", 'update': "Update User", 'view': "User Details" }


  constructor(
    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private route: ActivatedRoute,
    private generalService: GeneralService,
    private localService: LocalService,
    private utilService: UtilityService,
    private fileOperationService: FileOperationService,
    private datePipe: DatePipe,
    private urlConstant: UrlConstant,
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
    this.isOperationUpdate = false
    this.isSearched = false
    this.isViewMode = false

    // user limit offset
    this.limit = 5
    this.offset = 0

    // initialize forms
    this.createUserForm()
    this.createUserSearchForm()

    // Search.
    this.searchFilterFieldList = []

    this.permission = this.utilService.getMenu(this.localService.getJsonValue('menus'),
      this.urlConstant.ADMIN_USER)
  }

  getAllComponents() {
    this.getAllCountries()
    this.getAllRoles()
    this.searchOrGetUsers()

  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit() { }

  // Create form.
  createUserForm(): void {
    this.isResumeUploadedToServer = false
    this.displayedFileName = "Select file"
    this.docStatus = ""
    // Create form
    this.userForm = this.formBuilder.group({
      id: new FormControl(null),
      code: new FormControl(null),
      role: new FormControl(null, [Validators.required]),

      firstName: new FormControl(null,
        [Validators.required, Validators.maxLength(100), Validators.pattern(/^[a-zA-Z]+$/)]),
      lastName: new FormControl(null,
        [Validators.required, Validators.maxLength(100), Validators.pattern(/^[a-zA-Z]+$/)]),
      contact: new FormControl(null,
        [Validators.required, Validators.maxLength(100), Validators.pattern(/^[6789]\d{9}$/)]),
      dateOfBirth: new FormControl(null),
      dateOfJoining: new FormControl(null),
      email: new FormControl(null,
        [Validators.required, Validators.maxLength(100), Validators.pattern(/^([a-zA-Z0-9._%+-]+@[a-zA-Z0-9]+\.[a-zA-z]{2,3})/)]),
      resume: new FormControl(null),
      isActive: new FormControl(null),
      address: this.formBuilder.group({
        state: new FormControl(null, [Validators.required]),
        country: new FormControl(null, [Validators.required]),
        address: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
        city: new FormControl(null, [Validators.required, Validators.maxLength(50), Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/)]),
        pinCode: new FormControl(null, [Validators.required, Validators.pattern(/^[1-9][0-9]{5}$/)])
      })
    })
    this.userForm.get('address')?.get('state')?.disable()
  }

  // Create search form.
  createUserSearchForm(): void {
    this.userSearchForm = this.formBuilder.group({
      roleID: new FormControl(null),
      firstName: new FormControl(null, [Validators.maxLength(100)]),
      lastName: new FormControl(null, [Validators.maxLength(100)]),
      email: new FormControl(null, [Validators.maxLength(100)]),
      isActive: new FormControl('1'),
    });
  }

  resetSearchAndGetAll(): void {
    this.searchFilterFieldList = []
    this.createUserSearchForm()
    this.searchFormValue = null
    this.changePage(1)
    this.isSearched = false
    this.router.navigate([this.urlConstant.ADMIN_USER])
  }

  // Reset search form.
  resetSearchForm(): void {
    this.searchFilterFieldList = []
    this.userSearchForm.reset()
    this.userSearchForm.get('isActive').setValue("1")
  }


  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetUsers() {
    let queryParams = this.route.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getAllUsers()
      return
    }
    this.userSearchForm.patchValue(queryParams)
    this.searchUsers()
  }

  searchUsers(): void {
    this.searchFormValue = { ...this.userSearchForm?.value }
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
    this.spinnerService.loadingMessage = "Searching Users"
    this.changePage(1)
  }

  // Delete search criteria from designation search form by search name.
  deleteSearchCriteria(searchName: string): void {
    this.userSearchForm.get(searchName).setValue(null)
    this.searchUsers()
  }

  getAllCountries(): void {
    this.generalService.getCountries().subscribe((data: any) => {
      this.countryList = data
    }, (err) => {
      console.error(err)
    })
  }

  // shouldReset flag decides wheter to reset state control or not.
  getStateList(countryID: string, shouldReset?: boolean) {
    let stateControl = this.userForm.get('address').get('state')
    if (shouldReset) {
      stateControl.reset()
    }
    // Need to reset first incase of valid country with no states.
    stateControl.disable()
    this.stateList = [] as IState[]
    if (!countryID) {
      return
    }

    this.generalService.getStatesByCountryID(countryID).subscribe((data: IState[]) => {
      this.stateList = data
      // View mode needs to be checked as it is an async call, it will enable the field.
      if (this.stateList.length > 0 && !this.isViewMode) {
        stateControl.enable()
        return
      }
    }, (err) => {
      console.error(err)
    })
  }


  getAllRoles(): void {
    this.allRoles = []
    this.generalService.getAllRoles().subscribe((data: any) => {
      for (let role of data) {
        if (role.roleName === this.roleConstant.ADMIN || role.roleName === this.roleConstant.SALES_PERSON) {
          this.allRoles.push(role)
        }
      }
    }, (err) => {
      console.error(err)
    })
  }

  // page change function
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber;
    this.offset = this.currentPage - 1;
    this.getAllUsers();
  }

  onAddClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    this.createUserForm()
    this.openModal(this.userFormModal, 'xl')
  }

  // Update Form(prepopulate form)
  onViewClick(user: IUser): void {
    this.isViewMode = true
    this.isOperationUpdate = false
    this.createUserForm();

    // The date transform is handled here so that the form patches properly.
    user.dateOfBirth = this.datePipe.transform(user.dateOfBirth, 'yyyy-MM-dd')
    user.dateOfJoining = this.datePipe.transform(user.dateOfJoining, 'yyyy-MM-dd')
    this.userForm.patchValue(user);

    this.getStateList(user.address?.country?.id)
    this.userForm.disable()

    this.displayedFileName = "No resume uploaded"
    if (user.resume) {
      this.displayedFileName = `<a href=${user.resume} target="_blank">Resume present</a>`
    }
    this.openModal(this.userFormModal, 'xl')
  }

  onUpdateClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = true
    this.userForm.enable()
    this.userForm.get('code').disable()
  }

  onDeleteClick(userID: string): void {
    this.openModal(this.deleteConfirmationModal, 'md').result.then(() => {
      this.deleteUser(userID)
    }, (err) => {
      console.error(err);
      return
    })
  }



  // load all users in user list
  getAllUsers(): void {
    this.spinnerService.loadingMessage = "Getting Users"


    this.generalService.getSpecificUsers(this.limit, this.offset, this.searchFormValue).
      subscribe((data) => {
        this.totalUsers = data.headers.get('X-Total-Count');
        this.allUsers = data.body;
      }, (error) => {
        console.error(error);
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
        }
      }).add(() => {
        // Will be executed when subscription ends.
        this.setPaginationString()
      })
  }

  // Update User.
  updateUser(): void {
    this.spinnerService.loadingMessage = "Updating user";

    this.generalService.updateUser(this.userForm.value).subscribe((data) => {

      this.modalRef.close();
      this.getAllUsers();
      alert(data);
      this.userForm.reset();
    }, (error) => {
      console.error(error);

      if (error.error?.error) {
        alert(error.error.error);
        return;
      }
      alert("Check connection");
    });
  }

  // Add User.
  addUser(): void {
    this.spinnerService.loadingMessage = "Adding user"

    this.generalService.addUser(this.userForm.value).subscribe(() => {
      this.modalRef.close();
      this.getAllUsers();
      alert("User successfully added");
      this.userForm.reset();
    }, (error) => {
      console.error(error);

      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  //delete user after confirmation from user
  deleteUser(userID: string): void {
    this.spinnerService.loadingMessage = "Deleting user";
    this.modalRef.close();

    this.generalService.deleteUser(userID).subscribe((data) => {
      this.getAllUsers();
      alert("User deleted");
      this.userForm.reset();
    }, (error) => {
      console.error(error);

      if (error.error) {
        alert(error.error);
        return;
      }
      alert(error.statusText);
    });
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
    if (this.isFileUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    if (this.isResumeUploadedToServer) {
      if (!confirm("Uploaded resume will be deleted.\nAre you sure you want to close?")) {
        return
      }
      // #niranjan
      // this.deleteResume()
    }
    modal.dismiss()
  }

  onSubmit(): void {
    if (this.userForm.invalid) {
      this.userForm.markAllAsTouched();
      return
    }
    if (this.isOperationUpdate) {
      this.updateUser()
      return
    }
    this.addUser()
  }


  //On uplaoding resume
  onResourceSelect(event: any) {
    this.docStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      // Upload resume if it is present.]
      this.isFileUploading = true
      this.fileOperationService.uploadResume(file).subscribe((data: any) => {
        this.userForm.markAsDirty()
        this.userForm.patchValue({
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

  // Set total users list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalUsers < this.paginationEnd) {
      this.paginationEnd = this.totalUsers
    }
  }

  // Compare for select option field.
  compareFn(optionOne: any, optionTwo: any): boolean {
    if (optionOne == null && optionTwo == null) {
      return true;
    }
    if (optionTwo != undefined && optionOne != undefined) {
      return optionOne.id === optionTwo.id;
    }
    return false
  }
}
