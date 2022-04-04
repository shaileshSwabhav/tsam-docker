import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { AbstractControl, FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { AdminService, ITimesheet, ITimesheetActivity } from 'src/app/service/admin/admin.service';
import { BatchTopicService } from 'src/app/service/batch-topic/batch-topic.service';
import { BatchService, IBatch, IBatchSession, IBatchTopic } from 'src/app/service/batch/batch.service';
import { Role, UrlConstant } from 'src/app/service/constant';
import { GeneralService, IProject, IRole } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-timesheet',
  templateUrl: './timesheet.component.html',
  styleUrls: ['./timesheet.component.css'],
  // providers: [AdminService, GeneralService]
})
export class TimesheetComponent implements OnInit {

  // user details
  credentialID: string
  departmentID: string
  roleName: string

  // selected
  selectedCredentialID: string

  // components
  directReportList: ICredential[]
  batchList: IBatch[]
  // batchSessionList: IBatchSession[][]
  batchTopicList: IBatchTopic[][]
  projectList: IProject[]
  subProjectList: IProject[][]

  // flags
  isOperationUpdate: boolean
  isOperationAdd: boolean
  isViewMode: boolean

  areBatchDetailsVisible: boolean

  // timesheet
  salesPersonList: ITimesheet[]
  allTimesheets: ITimesheet[]
  timesheetForm: FormGroup
  multipleTimesheetForm: FormGroup

  // new timesheet
  newTimesheetForm: FormGroup
  timesheetActivityForm: FormGroup

  // search
  timesheetSearchForm: FormGroup
  searchFormValue: any
  showSearch: boolean
  isSearched: boolean
  isCredentialSearched: boolean
  isWeeklySearch: boolean

  // pagination
  limit: number;
  currentPage: number;
  offset: number;
  totalTimesheets: number;
  totalHours: number
  totalFreeHours: number
  paginationString: string
  weekIndex: number

  // modal
  modalRef: NgbModalRef;
  @ViewChild('timesheetFormModal') timesheetFormModal: any
  @ViewChild('deleteConfirmationModal') deleteConfirmationModal: any
  @ViewChild('newTimesheetModal') newTimesheetModal: any
  @ViewChild('timesheetActivityModal') timesheetActivityModal: any
  @ViewChild('previewModal') previewModal: any


  // validation
  readonly numValidator = /(^100([.]0{1,2})?)$|(^\d{1,2}([.]\d{1,2})?)$/i

  // spinner

  // #
  currentDate: Date

  // ng-select spinners
  subProjectLoadingSignals: boolean[]
  sessionLoadingSignals: boolean[]

  permission: IPermission
  // beta
  columns: IColumn[]
  subColumns: IColumn[]
  selectedOperation: IOperation
  readonly OPERATIONS: IOperation[] =
    [{ operation: "add", tag: "Add Timesheet" },
    { operation: "update", tag: "Update Timesheet" },
    { operation: "view", tag: "Timesheet Details" }]


  constructor(
    public readonly roleConstant: Role,

    private formBuilder: FormBuilder,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private router: Router,
    private activatedRoute: ActivatedRoute,
    private generalService: GeneralService,
    private batchService: BatchService,
    private batchTopicService: BatchTopicService,
    private localService: LocalService,
    private utilService: UtilityService,
    private datePipe: DatePipe,
    private urlConstant: UrlConstant,
    private adminService: AdminService,
  ) {
    this.initializeVariables()
    this.getAllComponents()

  }

  initializeVariables(): void {

    // User information
    this.roleName = this.localService.getJsonValue('roleName')
    this.credentialID = this.localService.getJsonValue('credentialID')
    this.departmentID = this.localService.getJsonValue('departmentID')
    this.selectedCredentialID = this.credentialID

    this.searchFormValue = []
    this.currentDate = new Date()

    // ng-select spinner
    this.subProjectLoadingSignals = []
    this.sessionLoadingSignals = []

    // flags
    this.isOperationUpdate = false
    this.isSearched = false
    this.isCredentialSearched = false
    this.isViewMode = false
    this.isWeeklySearch = true
    this.isOperationAdd = false
    this.areBatchDetailsVisible = false

    if (this.roleName === this.roleConstant.FACULTY) {
      this.areBatchDetailsVisible = true
    }

    // timesheet limit offset
    this.limit = 10
    this.offset = 0
    this.weekIndex = 0
    this.totalTimesheets = 0
    this.totalHours = 0
    this.totalFreeHours = 0

    // Lists
    this.projectList = []
    this.subProjectList = []
    this.batchList = []
    // this.batchSessionList = []
    this.batchTopicList = []
    this.directReportList = [{
      id: this.credentialID, firstName: "Self", lastName: "",
      role: {
        roleName: this.roleName
      }
    }]
    // if loadingMessageList is empty, we can close the spinner.


    // initialize forms
    // this.createTimesheetForm()
    this.createNewTimesheetForm()
    this.createTimesheetSearchForm()

    if (this.localService.getJsonValue("roleName") == this.roleConstant.FACULTY) {
      this.permission = this.utilService.getPermission(this.urlConstant.MY_TIMESHEET)
    }
    else {
      this.permission = this.utilService.getPermission(this.urlConstant.TIMESHEET)
    }
    this.columns = [
      { name: "Date", class: "col-2 text-break" },
      { name: "On Leave", class: "col-1 text-break" },
      { name: "Activities", class: "col-8 text-break" }]
  }

  getAllComponents(): void {
    this.getBatchList()
    this.getListOfProjects()
    this.searchOrGetTimesheets()
    this.getDirectReportsOrAllEmployees()
  }

  // ========================================testing================================================

  getColumnClass(name: string): string {
    return this.columns.find((val) => {
      return val.name.toLowerCase() === name.toLowerCase()
    }).class
  }

  onAddMultipleClick(): void {
    this.isOperationAdd = true
    this.showSearch = false
    this.createMultipleTimesheetForm()
  }

  cancelAdd(): void {
    this.isOperationAdd = false
    this.multipleTimesheetForm.reset()
    // this.batchSessionList = []
    this.batchTopicList = []
    this.sessionLoadingSignals = []
    this.subProjectList = []
    this.subProjectLoadingSignals = []
  }

  createMultipleTimesheetForm(): void {
    this.multipleTimesheetForm = this.formBuilder.group({
      timesheetArray: this.formBuilder.array([
        // this.createTimesheetForm()
        this.createNewTimesheetForm()
      ])
    })
    // this.batchSessionList = []
    this.batchTopicList = []
    this.sessionLoadingSignals = []
    this.subProjectList = []
    this.subProjectLoadingSignals = []
  }

  get timesheetArray(): FormArray {
    return this.multipleTimesheetForm.get('timesheetArray') as FormArray;
  }

  addNewActivity(timesheetForm: FormGroup): void {
    (timesheetForm.get('activities') as FormArray).push(this.createTimesheetActivityForm())
  }

  addNewTimesheetRow(): void {
    // this.timesheetArray.push(this.createTimesheetForm());
    this.timesheetArray.push(this.createNewTimesheetForm());
    this.subProjectList.push([])
    this.subProjectLoadingSignals.push(false)
    // this.batchSessionList.push([])
    this.batchTopicList.push([])
    this.sessionLoadingSignals.push(false)
  }

  duplicateTimesheetRow(timesheetForm: FormGroup): void {
    this.assignValidatorsForMultipleTimesheet(timesheetForm)
    if (this.multipleTimesheetForm.invalid) {
      this.multipleTimesheetForm.markAllAsTouched()
      return
    }
    let index: number = this.timesheetArray.length

    // this.timesheetArray.push(this.createTimesheetForm())
    this.timesheetArray.push(this.createNewTimesheetForm())

    // let newTimesheetForm = this.timesheetArray.at(index)
    // newTimesheetForm.get("date").setValue(timesheetForm.get("date").value)
    // newTimesheetForm.get("isOnLeave").setValue(timesheetForm.get("isOnLeave").value)

    // // Add project & subproject
    // let projectID = timesheetForm.get("projectID").value
    // newTimesheetForm.get("projectID").setValue(projectID)
    // let subProjectID = timesheetForm.get("subProjectID").value
    // this.subProjectList.push([])
    // this.subProjectLoadingSignals.push(false)
    // this.getListOfSubProjects(projectID, newTimesheetForm.get('subProjectID'), true, index)
    // newTimesheetForm.get("subProjectID").setValue(subProjectID)

    // // Add batch and session
    // let batchID = timesheetForm.get("batchID").value
    // newTimesheetForm.get("batchID").setValue(batchID)
    // let batchSessionID = timesheetForm.get("batchSessionID").value
    // this.batchSessionList.push([])
    // this.sessionLoadingSignals.push(false)
    // this.getBatchSessionList(batchID, newTimesheetForm.get('batchSessionID'), true, index)
    // newTimesheetForm.get("batchSessionID").setValue(batchSessionID)

    // newTimesheetForm.get("activity").setValue(timesheetForm.get("activity").value)
    // newTimesheetForm.get("hoursNeeded").setValue(timesheetForm.get("hoursNeeded").value)
  }

  deleteTimesheetRow(index: number): void {
    this.timesheetArray.removeAt(index);
    // this.batchSessionList.splice(index)
    this.batchTopicList.splice(index)
    this.sessionLoadingSignals.splice(index)
  }

  deleteActivityRowForTimesheetArray(rowIndex: number, activityIndex: number): void {
    let activityForm = this.timesheetArray.at(rowIndex).get('activities') as FormArray
    activityForm.removeAt(activityIndex)
  }

  deleteActivityRowForTimesheet(activityIndex: number): void {
    let activityForm = this.timesheetForm.get('activities') as FormArray
    activityForm.removeAt(activityIndex)
    this.timesheetForm.markAsDirty()
  }

  previewAddedTimesheets(): void {
    for (let i = 0; i < this.timesheetArray.length; i++) {
      this.assignValidatorsForMultipleTimesheet(this.timesheetArray.at(i))
    }
    if (this.multipleTimesheetForm.invalid) {
      this.multipleTimesheetForm.markAllAsTouched()
      return
    }
    this.openModal(this.previewModal, 'xl')
  }

  getEntityByID(entityName: string, entityID: string, list?: any[]) {
    if (entityName === "project") {
      return this.projectList.find((val) => {
        return val.id === entityID
      }).name
    }
    if (entityName === "subProject") {
      return list?.find((val: any) => {
        return val.id === entityID
      })?.name
    }
    if (entityName === "batch") {
      return this.batchList.find((val) => {
        return val.id === entityID
      }).batchName
    }
    if (entityName === "batchTopic") {
      return list?.find((val) => {
        return val.id === entityID
      })?.name
    }
    // if (entityName === "batchSession") {
    //   return list?.find((val) => {
    //     return val.id === entityID
    //   })?.name
    // }
  }

  // Add multiple Timesheet.
  addMultipleTimesheets(): void {
    for (let i = 0; i < this.timesheetArray.length; i++) {
      this.assignValidatorsForMultipleTimesheet(this.timesheetArray.at(i))
    }
    if (this.multipleTimesheetForm.invalid) {
      this.multipleTimesheetForm.markAllAsTouched()
      return
    }
    this.spinnerService.loadingMessage = "Adding timesheets"



    this.adminService.addMultipleTimesheets(this.timesheetArray.value).subscribe(() => {
      this.isOperationAdd = false
      this.searchOrGetTimesheets();
      alert("Timesheet successfully added");
      this.multipleTimesheetForm.reset()
      this.modalRef.close()
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }


  // For multiple timesheet
  assignValidatorsForMultipleTimesheet(timesheetRow: FormGroup | AbstractControl): void {
    let activityControl = (timesheetRow.get("activities") as FormArray).controls

    if (timesheetRow.get('isOnLeave').value === true) {

      for (let index = 0; index < activityControl.length; index++) {
        activityControl[index].get('project').setValidators(null)
        activityControl[index].get('projectID').setValidators(null)
        activityControl[index].get('subProject').setValidators(null)
        activityControl[index].get('subProjectID').setValidators(null)
        activityControl[index].get('activity').setValidators(null)
        activityControl[index].get('hoursNeeded').setValidators([Validators.pattern(this.numValidator), Validators.min(0.1)])

        // if (activityControl[index].get('isCompleted')?.value === false) {
        //   activityControl[index].get('nextEstimatedDate').setValidators([Validators.required])
        // } else {
        //   activityControl[index].get('nextEstimatedDate').clearValidators()
        // }

        this.utilService.updateValueAndValiditors(activityControl[index] as FormGroup)
      }
    } else {
      timesheetRow.get('date').setValidators([Validators.required])
      // form.get('project').setValidators([Validators.required])
      // form.get('subProject').setValidators([Validators.required])
      for (let index = 0; index < activityControl.length; index++) {
        activityControl[index].get('projectID').setValidators([Validators.required])
        activityControl[index].get('subProjectID').setValidators([Validators.required])
        activityControl[index].get('activity').setValidators([Validators.required])
        activityControl[index].get('hoursNeeded').setValidators([Validators.required,
        Validators.pattern(this.numValidator), Validators.min(0.1)])

        // if (activityControl[index].get('isCompleted')?.value == false) {
        //   activityControl[index].get('nextEstimatedDate').setValidators([Validators.required])
        // } else {
        //   activityControl[index].get('nextEstimatedDate').clearValidators()
        // }

        this.utilService.updateValueAndValiditors(activityControl[index] as FormGroup)
      }

    }
    this.utilService.updateValueAndValiditors(timesheetRow as FormGroup)
  }


  // ========================================testing================================================


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit() { }

  // Create form.
  createTimesheetForm(): void {
    // Create form
    // return this.timesheetForm = this.formBuilder.group({
    //   id: new FormControl(null),
    //   tenantID: new FormControl(null),
    //   date: new FormControl(null, [Validators.required]),
    //   departmentID: new FormControl(this.departmentID, [Validators.required]),
    //   credentialID: new FormControl(this.credentialID, [Validators.required]),
    //   isOnLeave: new FormControl(false, [Validators.required]),
    //   project: new FormControl(null),
    //   projectID: new FormControl(null),
    //   subProject: new FormControl({ value: null, disabled: true }),
    //   subProjectID: new FormControl({ value: null, disabled: true }),
    //   activity: new FormControl(null),
    //   hoursNeeded: new FormControl(null, [Validators.pattern(this.numValidator), Validators.min(0.1)]),
    //   batchID: new FormControl(null),
    //   batch: new FormControl(null),
    //   batchSessionID: new FormControl({ value: null, disabled: true }),
    //   batchSession: new FormControl({ value: null, disabled: true }),
    //   // for edit only
    //   isBillable: new FormControl(null),
    //   isCompleted: new FormControl(null),
    //   workDone: new FormControl(null),
    //   nextEstimatedDate: new FormControl(null),
    // })
  }

  // For single timesheet
  assignNecessaryValidators(): void {
    let activityControl = (this.timesheetForm.get("activities") as FormArray).controls

    if (this.timesheetForm.get('isOnLeave').value === true) {

      for (let index = 0; index < activityControl.length; index++) {
        if (activityControl[index].get('isCompleted').value === false) {
          activityControl[index].get('nextEstimatedDate').setValidators([Validators.required])
        } else {
          activityControl[index].get('nextEstimatedDate').setValidators(null)
        }
        activityControl[index].get('project').setValidators(null)
        activityControl[index].get('projectID').setValidators(null)
        activityControl[index].get('subProject').setValidators(null)
        activityControl[index].get('subProjectID').setValidators(null)
        activityControl[index].get('activity').setValidators(null)
        activityControl[index].get('hoursNeeded').setValidators([Validators.pattern(this.numValidator), Validators.min(0.1)])

        this.utilService.updateValueAndValiditors(activityControl[index] as FormGroup)
      }

      // this.timesheetForm.get('project').setValidators(null)
      // this.timesheetForm.get('projectID').setValidators(null)
      // this.timesheetForm.get('subProject').setValidators(null)
      // this.timesheetForm.get('subProjectID').setValidators(null)
      // this.timesheetForm.get('activity').setValidators(null)
      // this.timesheetForm.get('hoursNeeded').setValidators([Validators.pattern(this.numValidator), Validators.min(0.1)])
    } else {
      this.timesheetForm.get('date').setValidators([Validators.required])
      // form.get('project').setValidators([Validators.required])
      // form.get('subProject').setValidators([Validators.required])
      for (let index = 0; index < activityControl.length; index++) {
        activityControl[index].get('projectID').setValidators([Validators.required])
        activityControl[index].get('project').setValidators([Validators.required])
        activityControl[index].get('subProjectID').setValidators([Validators.required])
        activityControl[index].get('subProject').setValidators([Validators.required])
        activityControl[index].get('activity').setValidators([Validators.required])
        activityControl[index].get('hoursNeeded').setValidators([Validators.required,
        Validators.pattern(this.numValidator), Validators.min(0.1)])

        if (activityControl[index].get('isCompleted')?.value === false) {
          activityControl[index].get('nextEstimatedDate').setValidators([Validators.required])
        } else {
          activityControl[index].get('nextEstimatedDate').clearValidators()
        }

        this.utilService.updateValueAndValiditors(activityControl[index] as FormGroup)
      }
      // this.timesheetForm.get('project').setValidators([Validators.required])
      // this.timesheetForm.get('projectID').setValidators([Validators.required])
      // this.timesheetForm.get('subProject').setValidators([Validators.required])
      // this.timesheetForm.get('subProjectID').setValidators([Validators.required])
      // this.timesheetForm.get('activity').setValidators([Validators.required])
      // this.timesheetForm.get('hoursNeeded').setValidators([Validators.required, Validators.pattern(this.numValidator), Validators.min(0.1)])
    }
    this.utilService.updateValueAndValiditors(this.timesheetForm)
  }

  // Create search form.
  createTimesheetSearchForm(): void {
    this.timesheetSearchForm = this.formBuilder.group({
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
      departmentID: new FormControl(null),
      credentialID: new FormControl(this.credentialID),
      projectID: new FormControl(null),
      batchID: new FormControl(null),
      // batchSessionID: new FormControl({ value: null, disabled: true }),
      batchTopicID: new FormControl({ value: null, disabled: true }),
      isOnLeave: new FormControl(null),
      isCompleted: new FormControl(null),
    });
  }


  resetSearchAndGetAll(): void {
    this.searchFormValue = []
    this.changePage(1)
    this.isSearched = false
    this.isCredentialSearched = false
    this.isWeeklySearch = false
    this.weekIndex = null
    this.showSearch = false
    this.onResetClick()
  }

  onResetClick(): void {
    this.router.navigate(
      ['.'],
      { relativeTo: this.activatedRoute, queryParams: { "credentialID": this.credentialID } }
    ).catch((err: any) => {
      console.error(err);
    })
    // this.timesheetSearchForm.get('batchSessionID').setValue(null)
    this.timesheetSearchForm.get('batchTopicID').setValue(null)
    // doesn't work if it is after form creation for some reason #Niranjan.
    // this.timesheetSearchForm.get('batchSessionID').disable()
    this.timesheetSearchForm.get('batchTopicID').disable()
    this.createTimesheetSearchForm()
  }

  // Checks the url's query params and decides to call whether to call get or search.
  searchOrGetTimesheets(): void {
    let queryParams = this.activatedRoute.snapshot.queryParams
    if (this.utilService.isObjectEmpty(queryParams)) {
      this.getWeeklyTimesheet()
      return
    }
    this.timesheetSearchForm.patchValue(queryParams)
    let batchID = this.timesheetSearchForm.get("batchID").value
    if (batchID) {
      // this.getBatchSessionList(batchID, this.timesheetSearchForm.get("batchSessionID"))
      // this.getBatchTopicList(batchID, this.timesheetSearchForm.get("batchTopicID"))
    }
    this.searchTimesheets()
  }

  searchTimesheets(isDirectReportSearch?: boolean, isWeeklySearch?: boolean): void {
    if (!isWeeklySearch && !isDirectReportSearch) {
      this.isWeeklySearch = false
      this.weekIndex = null
    }
    this.searchFormValue = { ...this.timesheetSearchForm?.getRawValue() }
    this.router.navigate([], {
      relativeTo: this.activatedRoute,
      queryParams: this.searchFormValue,
    })
    for (let field in this.searchFormValue) {
      if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
        delete this.searchFormValue[field]
      } else if (field === "credentialID") {
        this.isCredentialSearched = true
      } else {
        this.isSearched = true
        if (!isWeeklySearch && !isDirectReportSearch) {
          this.showSearch = true
        }
      }
    }
    if (!this.isSearched && !this.isCredentialSearched && !this.isWeeklySearch) {
      return
    }
    this.spinnerService.loadingMessage = "Searching Timesheets"
    this.changePage(1)
  }

  getDirectReportsOrAllEmployees(): void {


    if (this.roleName === this.roleConstant.ADMIN) {
      this.generalService.getAllEmployeeList().subscribe((data: any) => {
        this.directReportList = data.body
      }, (err) => {
        console.error(err)
      })
      return
    }
    this.generalService.getDirectReports(this.credentialID).subscribe((data: any) => {
      this.directReportList = this.directReportList.concat(data.body)
    }, (err) => {
      console.error(err)
    })
  }

  getBatchList(): void {
    let queryParams: any = {
      isActive: '1'
    }



    this.batchService.getBatchList(queryParams).subscribe((data: any) => {
      this.batchList = data.body
    }, (err) => {
      console.error(err)
    })
  }

  getBatchTopicList(batchID: string, activityForm: FormGroup | AbstractControl,
    shouldReset?: boolean, index?: number): void {
    console.log(activityForm);

    if (shouldReset) {
      activityForm.get('batchTopicID')?.reset()
      activityForm.get('batchTopic')?.reset()
    }
    if (!index) {
      index = 0
    }

    // if it is called on clear, we need to clear the list before returning.
    this.batchTopicList[index] = []
    activityForm.get('batchTopicList')?.setValue(null)
    activityForm.get('batchTopicID')?.disable()
    activityForm.get('batchTopic')?.disable()
    if (!batchID) {
      return
    }
    this.sessionLoadingSignals[index] = true
    activityForm.get('isBatchTopicLoading')?.setValue(true)

    let queryParams: any = {
      "filter": "list"
    }

    this.batchTopicService.getBatchTopic(batchID, queryParams).subscribe((data: any) => {
      this.batchTopicList[index] = data.body
      activityForm.get('batchTopicList')?.setValue(data.body)

      // View mode needs to be checked as it is an async call, it will enable the field.
      if ((activityForm.get('batchTopicList')?.value.length > 0 || this.batchTopicList[index].length) > 0
        && !this.isViewMode) {
        activityForm.get('batchTopicID')?.enable()
        activityForm.get('batchTopic')?.enable()
      }
    }, (err) => {
      console.error(err)
    }).add(() => {
      this.sessionLoadingSignals[index] = false
      activityForm.get('isBatchTopicLoading')?.setValue(false)
    })
  }

  // getBatchSessionList(batchID: string, activityForm: FormGroup | AbstractControl,
  //   shouldReset?: boolean, index?: number): void {
  //   if (shouldReset) {
  //     activityForm.get('batchSessionID').reset()
  //     activityForm.get('batchSession').reset()
  //   }
  //   if (!index) {
  //     index = 0
  //   }
  //   // if it is called on clear, we need to clear the list before returning.
  //   this.batchSessionList[index] = []
  //   activityForm.get('batchSessionList').setValue(null)
  //   activityForm.get('batchSessionID').disable()
  //   activityForm.get('batchSession').disable()
  //   if (!batchID) {
  //     return
  //   }
  //   this.sessionLoadingSignals[index] = true
  //   activityForm.get('isBatchSessionLoading').setValue(true)
  //   this.batchService.getBatchSessionList(batchID).subscribe((data: any) => {
  //     this.batchSessionList[index] = data.body
  //     activityForm.get('batchSessionList').setValue(data.body)
  //     // View mode needs to be checked as it is an async call, it will enable the field.
  //     if (activityForm.get('batchSessionList').value.length > 0 && this.batchSessionList[index].length > 0
  //       && !this.isViewMode) {
  //       activityForm.get('batchSessionID').enable()
  //       activityForm.get('batchSession').enable()
  //     }
  //   }, (err) => {
  //     console.error(err)
  //   }).add(() => {
  //     this.sessionLoadingSignals[index] = false
  //     activityForm.get('isBatchSessionLoading').setValue(false)
  //   })
  // }

  getListOfProjects(): void {


    this.generalService.getListOfProjects().subscribe((data: any) => {
      this.projectList = data.body
    }, (err) => {
      console.error(err)
    })
  }

  getListOfSubProjects(projectID: string, activityForm: FormGroup | AbstractControl,
    shouldReset?: boolean, index?: number): void {

    if (shouldReset) {
      activityForm.get('subProjectID').reset()
      activityForm.get('subProject').reset()
    }
    if (!index) {
      index = 0
    }
    // this.subProjectList[index] = []
    activityForm.get('subProjectList').setValue(null)
    activityForm.get('subProjectID').disable()
    activityForm.get('subProject').disable()
    if (!projectID) {
      return
    }
    this.subProjectLoadingSignals[index] = true
    activityForm.get('isSubProjectLoading').setValue(true)

    this.generalService.getListOfSubProjects(projectID).subscribe((data: any) => {
      this.subProjectList[index] = data.body
      activityForm.get('subProjectList').setValue(data.body)

      // View mode needs to be checked as it is an async call, it will enable the field.
      if (activityForm.get('subProjectList').value.length > 0 && this.subProjectList.length > 0
        && !this.isViewMode) {
        activityForm.get('subProject').enable()
        activityForm.get('subProjectID').enable()
      }
    }, (err) => {
      console.error(err)
    }).add(() => {
      this.subProjectLoadingSignals[index] = false
      activityForm.get('isSubProjectLoading').setValue(false)
    })
  }

  // page change function
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber;
    this.offset = this.currentPage - 1;
    this.getAllTimesheets();
  }

  onAddClick(): void {
    this.isViewMode = false
    this.isOperationUpdate = false
    // this.createTimesheetForm()
    this.createNewTimesheetForm()
    this.openModal(this.timesheetFormModal, 'xl')
  }

  // Update Form(prepopulate form)
  onViewClick(timesheet: ITimesheet): void {
    this.isViewMode = true
    this.isOperationUpdate = false

    // this.createTimesheetForm();
    this.createNewTimesheetForm();

    for (let index = 0; index < timesheet.activities.length - 1; index++) {
      (this.timesheetForm.get('activities') as FormArray).push(this.createTimesheetActivityForm())
    }
    this.patchFormattedValueToForm(timesheet)
    // console.log("timesheetForm",this.timesheetForm.value)
    this.timesheetForm.disable()
    this.openNewTimesheetModal(this.newTimesheetModal, 'xl')
  }

  patchFormattedValueToForm(timesheet: ITimesheet) {

    // The date transform is handled here so that the form patches properly.
    timesheet.date = this.datePipe.transform(timesheet.date, 'yyyy-MM-dd')

    for (let index = 0; index < timesheet.activities.length; index++) {
      timesheet.activities[index].nextEstimatedDate =
        this.datePipe.transform(timesheet.activities[index].nextEstimatedDate, 'yyyy-MM-dd')
    }
    this.timesheetForm.patchValue(timesheet)
  }

  onUpdateClick(): void {
    if (this.isOperationAdd) {
      if (!confirm("Proceeding will cancel the add operation. Do you wish to proceed?")) {
        return
      }
      this.cancelAdd()
    }
    this.isViewMode = false
    let activityControl = (this.timesheetForm.get("activities") as FormArray).controls
    this.isOperationUpdate = true
    this.timesheetForm.enable()

    for (let index = 0; index < activityControl.length; index++) {
      if (activityControl[index].get("nextEstimatedDate").value) {
        activityControl[index].disable()
        continue
      }
      this.getListOfSubProjects(activityControl[index].get("project").value?.id, activityControl[index])
      // this.getBatchSessionList(activityControl[index].get("batch").value?.id, activityControl[index])
      // this.getBatchTopicList(activityControl[index].get("batch").value?.id, activityControl[index])
    }
  }

  onMarkCompleteClick(timesheetActivity: ITimesheetActivity): void {
    if (this.isOperationAdd) {
      if (!confirm("Proceeding will cancel the add operation. Do you wish to proceed?")) {
        return
      }
      this.cancelAdd()
    }

    if (!confirm("Are you sure you want to mark this acitivity as complete?")) {
      return
    }

    timesheetActivity.isCompleted = true
    this.createTimesheetActivityForm()
    this.timesheetActivityForm.patchValue(timesheetActivity)

    if (timesheetActivity.subProject) {
      this.timesheetActivityForm.get('subProjectID').enable()
    }

    // if (timesheetActivity.batchSession) {
    //   this.timesheetActivityForm.get('batchSessionID').enable()
    // }
    if (timesheetActivity.batchTopic) {
      this.timesheetActivityForm.get('batchTopicID').enable()
    }

    // patch id
    this.patchTimesheetActivityID()
    // console.log(this.timesheetActivityForm.controls)
    this.updateActivityComplete(timesheetActivity)
  }

  // Update Timesheet Activity.
  updateActivityComplete(timesheetActivity: ITimesheetActivity): void {
    this.spinnerService.loadingMessage = "Updating timesheet activity";


    this.adminService.updateTimesheetActivity(this.timesheetActivityForm.value).subscribe((data) => {
      // this.getAllTimesheets();
      alert(data);
      this.timesheetActivityForm.reset();
    }, (error) => {
      timesheetActivity.isCompleted = false
      console.error(error);
      if (error.error?.error) {
        alert(error.error.error);
        return;
      }
      alert("Check connection");
    });
  }


  onUpdateTimesheetActivityClick(timesheetActivity: ITimesheetActivity): void {
    if (this.isOperationAdd) {
      if (!confirm("Proceeding will cancel the add operation. Do you wish to proceed?")) {
        return
      }
      this.cancelAdd()
    }

    this.createTimesheetActivityForm()

    if (timesheetActivity.nextEstimatedDate) {
      timesheetActivity.nextEstimatedDate = this.datePipe.transform(timesheetActivity.nextEstimatedDate, 'yyyy-MM-dd')
      this.timesheetActivityForm.disable()
    } else {
      this.getListOfSubProjects(timesheetActivity.project.id, this.timesheetActivityForm)
      if (timesheetActivity.batch) {
        // this.getBatchSessionList(timesheetActivity.batch.id, this.timesheetActivityForm)
        // this.getBatchTopicList(timesheetActivity.batch.id, this.timesheetActivityForm)
      }
    }
    this.timesheetActivityForm.patchValue(timesheetActivity)
    this.openNewTimesheetModal(this.timesheetActivityModal, 'xl')
  }

  onDeleteClick(timesheetID: string): void {
    this.openModal(this.deleteConfirmationModal, 'md').result.then(() => {
      this.deleteTimesheet(timesheetID)
    }, (err) => {
      console.error(err);
      return
    })
  }

  // load all timesheets in timesheet list
  getAllTimesheets(): void {
    this.spinnerService.loadingMessage = "Getting Timesheets"


    this.setCrendential()

    // if (this.isWeeklySearch) {
    //   this.limit = -1
    // } else if (this.limit === -1) {
    //   this.limit = 5
    // }
    this.allTimesheets = []

    this.adminService.getTimesheets(this.limit, this.offset, this.searchFormValue).
      subscribe((data) => {
        this.totalTimesheets = data.headers.get('X-Total-Count');
        this.totalHours = data.headers.get('Total-Hours');
        this.totalFreeHours = data.headers.get('Free-Hours');
        this.allTimesheets = data.body;
      }, (error) => {
        console.error(error);
        // #niranjan untested
        this.initializeVariables()
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(error.statusText)
      }).add(() => {
        // Will work like finally.
        this.setPaginationString()
        this.isWeeklySearch = false
      })
  }

  setCrendential(): void {
    if (!this.searchFormValue["credentialID"]) {
      this.timesheetSearchForm.get('credentialID').setValue(this.credentialID)
      this.searchFormValue["credentialID"] = this.credentialID
    }
    this.selectedCredentialID = this.searchFormValue["credentialID"]
  }

  // Update Timesheet.
  updateTimesheet(): void {
    this.spinnerService.loadingMessage = "Updating timesheet";

    // return



    this.adminService.updateTimesheet(this.timesheetForm.value).subscribe((data) => {
      this.modalRef.close();
      this.getAllTimesheets();
      alert(data);
      this.timesheetForm.reset();
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error.error);
        return;
      }
      alert("Check connection");
    });
  }

  // Update Timesheet Activity.
  updateTimesheetActivity(): void {
    this.spinnerService.loadingMessage = "Updating timesheet activity";


    this.adminService.updateTimesheetActivity(this.timesheetActivityForm.value).subscribe((data) => {
      this.modalRef.close();
      this.getAllTimesheets();
      alert(data);
      this.timesheetActivityForm.reset();
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error.error);
        return;
      }
      alert("Check connection");
    });
  }

  // Add Timesheet.
  addTimesheet(): void {
    this.spinnerService.loadingMessage = "Adding timesheet"



    this.adminService.addTimesheet(this.timesheetForm.value).subscribe(() => {
      this.modalRef.close();
      this.getAllTimesheets();
      alert("Timesheet successfully added");
      this.timesheetForm.reset();
    }, (error) => {
      console.error(error);
      if (error.error?.error) {
        alert(error.error?.error);
        return;
      }
      alert(error.statusText);
    })
  }

  //delete timesheet after confirmation from timesheet
  deleteTimesheet(timesheetID: string): void {
    this.spinnerService.loadingMessage = "Deleting timesheet";
    this.modalRef.close();



    this.adminService.deleteTimesheet(timesheetID).subscribe((data) => {
      this.getAllTimesheets();
      alert("Timesheet deleted");
      this.timesheetForm.reset();
    }, (error) => {
      console.error(error);
      if (error.error) {
        alert(error.error);
        return;
      }
      alert(error.statusText);
    })
  }

  onIsOnLeaveChange(form: FormGroup | AbstractControl) {

    if (form.get('isOnLeave').value === true) {
      form.reset({
        id: form.get('id').value,
        isOnLeave: true,
        date: form.get('date').value,
        credentialID: form.get('credentialID').value,
        departmentID: form.get('departmentID').value
      })
      form.markAsDirty()
      let formArray = form.get('activities') as FormArray
      let firstControl = formArray.at(0)
      formArray.clear()
      formArray.push(firstControl)
      // while (formArray.length > 1) {
      //   formArray.removeAt(1)
      // }
      return
    }
  }

  // # take event and check role to hide show batch/sessions
  onDirectReportChange(credential: ICredential): void {
    this.showOrHideBatchDetails(credential)
    this.timesheetSearchForm.get('credentialID').setValue(this.selectedCredentialID)
    this.searchTimesheets(true)
  }

  showOrHideBatchDetails(credential: ICredential) {
    if (credential?.role?.roleName === this.roleConstant.FACULTY) {
      this.areBatchDetailsVisible = true
      return
    }
    this.areBatchDetailsVisible = false
  }

  isBatchAndSessionsVisible(activityForm: FormGroup | AbstractControl): boolean {
    if ((this.timesheetForm.get('isOnLeave').value === false) || (this.roleName === this.roleConstant.FACULTY ||
      (this.roleName === this.roleConstant.ADMIN && activityForm.get('batch').value))) {
      return true
    }
    return false
  }

  getWeeklyTimesheet(): void {
    if (this.weekIndex === null) {
      this.isWeeklySearch = false
      this.getAllTimesheets()
    }
    this.isWeeklySearch = true
    let fromDate: Date = this.getPreviousMonday()
    let toDate: Date = new Date(fromDate)
    toDate.setDate(toDate.getDate() + 6)

    this.timesheetSearchForm.get('toDate').setValue(this.datePipe.transform(toDate, 'yyyy-MM-dd'))
    this.timesheetSearchForm.get('fromDate').setValue(this.datePipe.transform(fromDate, 'yyyy-MM-dd'))
    this.searchTimesheets(false, true)
  }

  /**
  * getPreviousMonday
  * @returns date depending on the weekIndex value, 0 is for this week,
  *  1 is for previous week and so on.
  */
  getPreviousMonday(): Date {
    let date = new Date();
    let day = date.getDay()
    let prevMonday = new Date()
    let index = this.weekIndex * 7
    // if (day == 1) {
    //   prevMonday.setDate(date.getDate() - (7 + index))
    //   return prevMonday
    // }
    let subDate = ((day + 6) % 7) + index
    prevMonday.setDate(date.getDate() - Math.abs(subDate))
    return prevMonday
  }
  // getPreviousMonday(): Date {
  // let date = new Date();
  // date.setDate(date.getDate() + 3)

  // let day = date.getDay();
  // let prevMonday = new Date();
  // if (day == 1 || day == 0) {
  //   prevMonday.setDate(date.getDate() - (day + 6));
  // }
  // else {
  //   prevMonday.setDate(date.getDate() - (day - 1));
  // }
  //   return prevMonday
  // }

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
  openNewTimesheetModal(content: any, size?: string): NgbModalRef {
    if (!size) {
      size = 'lg'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: size, scrollable: true,
      centered: true,
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef): void {
    modal.dismiss()
  }

  assignTimesheetActivityValidators(): void {
    this.timesheetActivityForm.get('projectID').setValidators([Validators.required])
    this.timesheetActivityForm.get('project').setValidators([Validators.required])
    this.timesheetActivityForm.get('subProjectID').setValidators([Validators.required])
    this.timesheetActivityForm.get('subProject').setValidators([Validators.required])
    this.timesheetActivityForm.get('activity').setValidators([Validators.required])
    this.timesheetActivityForm.get('hoursNeeded').setValidators([Validators.required,
    Validators.pattern(this.numValidator), Validators.min(0.1)])

    if (this.timesheetActivityForm.get('isCompleted')?.value === false) {
      this.timesheetActivityForm.get('nextEstimatedDate').setValidators([Validators.required])
    } else {
      this.timesheetActivityForm.get('nextEstimatedDate').clearValidators()
    }

    this.utilService.updateValueAndValiditors(this.timesheetActivityForm)
  }

  patchTimesheetActivityID(): void {

    // patch id's
    let project = this.timesheetActivityForm.get('project').value
    if (project) {
      this.timesheetActivityForm.get('projectID').setValue(project.id)
    }
    let subProject = this.timesheetActivityForm.get('subProject').value
    if (subProject) {
      this.timesheetActivityForm.get('subProjectID').setValue(subProject.id)
    }
    let batch = this.timesheetActivityForm.get('batch').value
    if (batch) {
      this.timesheetActivityForm.get('batchID').setValue(batch.id)
    }
    // let batchSession = this.timesheetActivityForm.get('batchSession').value
    // if (batchSession) {
    //   this.timesheetActivityForm.get('batchSessionID').setValue(batchSession.id)
    // }
    let batchTopic = this.timesheetActivityForm.get('batchTopic').value
    if (batchTopic) {
      this.timesheetActivityForm.get('batchTopicID').setValue(batchTopic.id)
    }

    if (this.timesheetActivityForm.get('isCompleted')?.value === true) {
      this.timesheetActivityForm.get('nextEstimatedDate').setValue(null)
    }

    if (this.timesheetActivityForm.invalid) {
      this.timesheetActivityForm.markAllAsTouched()
      return
    }
  }

  onActivitySubmit(): void {
    this.assignTimesheetActivityValidators()
    this.patchTimesheetActivityID()
    this.updateTimesheetActivity()
    // // patch id's
    // let project = this.timesheetActivityForm.get('project').value
    // if (project) {
    //   this.timesheetActivityForm.get('projectID').setValue(project.id)
    // }
    // let subProject = this.timesheetActivityForm.get('subProject').value
    // if (subProject) {
    //   this.timesheetActivityForm.get('subProjectID').setValue(subProject.id)
    // }
    // let batch = this.timesheetActivityForm.get('batch').value
    // if (batch) {
    //   this.timesheetActivityForm.get('batchID').setValue(batch.id)
    // }
    // let batchSession = this.timesheetActivityForm.get('batchSession').value
    // if (batchSession) {
    //   this.timesheetActivityForm.get('batchSessionID').setValue(batchSession.id)
    // }

    // if (this.timesheetActivityForm.get('isCompleted')?.value === true) {
    //   this.timesheetActivityForm.get('nextEstimatedDate').setValue(null)
    // }

    // if (this.timesheetActivityForm.invalid) {
    //   this.timesheetActivityForm.markAllAsTouched()
    //   return
    // }


  }

  onSubmit(): void {
    this.assignNecessaryValidators()
    this.patchIDFromObjects()

    if (this.timesheetForm.invalid) {
      this.timesheetForm.markAllAsTouched();
      return
    }

    if (this.isOperationUpdate) {
      this.updateTimesheet()
      return
    }
    this.addTimesheet()
  }

  patchIDFromObjects(): void {
    let activityControl = (this.timesheetForm.get("activities") as FormArray).controls

    for (let index = 0; index < activityControl.length; index++) {
      let project = activityControl[index].get('project').value
      if (project) {
        activityControl[index].get('projectID').setValue(project.id)
      }
      let subProject = activityControl[index].get('subProject').value
      if (subProject) {
        activityControl[index].get('subProjectID').setValue(subProject.id)
      }
      let batch = activityControl[index].get('batch').value
      if (batch) {
        activityControl[index].get('batchID').setValue(batch.id)
      }
      // let batchSession = activityControl[index].get('batchSession').value
      // if (batchSession) {
      //   activityControl[index].get('batchSessionID').setValue(batchSession.id)
      // }
      let batchTopic = activityControl[index].get('batchTopic').value
      if (batchTopic) {
        activityControl[index].get('batchTopicID').setValue(batchTopic.id)
      }
    }

  }

  setPaginationString() {
    this.paginationString = ''
    if (this.isWeeklySearch) {
      return
    }
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (this.totalTimesheets < end) {
      end = this.totalTimesheets
    }
    if (this.totalTimesheets == 0) {
      this.paginationString = ''
      return
    }
    if (!start || !end) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start}-${end} of ${this.totalTimesheets}`
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

  // ============================================TESTING=====================================================

  createNewTimesheetForm(): FormGroup {
    return this.timesheetForm = this.formBuilder.group({
      id: new FormControl(null),
      tenantID: new FormControl(null),
      date: new FormControl(null, [Validators.required]),
      departmentID: new FormControl(this.departmentID, [Validators.required]),
      credentialID: new FormControl(this.credentialID, [Validators.required]),
      isOnLeave: new FormControl(false, [Validators.required]),
      activities: this.formBuilder.array([
        this.createTimesheetActivityForm()
      ]),
    })
  }

  createTimesheetActivityForm(): FormGroup {
    return this.timesheetActivityForm = this.formBuilder.group({
      id: new FormControl(null),
      timesheetID: new FormControl(null),
      project: new FormControl(null),
      projectID: new FormControl(null),
      subProject: new FormControl({ value: null, disabled: true }),
      subProjectID: new FormControl({ value: null, disabled: true }),
      activity: new FormControl(null, [Validators.maxLength(1000)]),
      hoursNeeded: new FormControl(null, [Validators.pattern(this.numValidator), Validators.min(0.1)]),
      batchID: new FormControl(null),
      batch: new FormControl(null),
      // batchSessionID: new FormControl({ value: null, disabled: true }),
      // batchSession: new FormControl({ value: null, disabled: true }),
      batchTopicID: new FormControl({ value: null, disabled: true }),
      batchTopic: new FormControl({ value: null, disabled: true }),
      // for edit only
      isBillable: new FormControl(null),
      isCompleted: new FormControl(null),
      workDone: new FormControl(null),
      nextEstimatedDate: new FormControl(null),
      // list
      // batchSessionList: new FormControl(null),
      // isBatchSessionLoading: new FormControl(null),
      batchTopicList: new FormControl(null),
      isBatchTopicLoading: new FormControl(null),
      subProjectList: new FormControl(null),
      isSubProjectLoading: new FormControl(null)
    })
  }

  get timesheetArrayControls() {
    return this.timesheetForm.get('activities') as FormArray
  }
}


interface IOperation {
  operation: string
  tag: string
}

interface ICredential {
  id: string
  firstName?: string
  lastName?: string
  role?: IRole
}

interface IColumn {
  name: string
  class?: string
}