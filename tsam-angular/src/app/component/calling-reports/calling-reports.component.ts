import { formatDate } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { CallingReportService, IDaywiseEnquiryCallingReport, ILoginwiseEnquiryCallingReport, ITalentEnquiryCallingReport } from 'src/app/service/calling-report/calling-report.service';
import { UrlConstant } from 'src/app/service/constant';
import { GeneralService } from 'src/app/service/general/general.service';

@Component({
  selector: 'app-calling-reports',
  templateUrl: './calling-reports.component.html',
  styleUrls: ['./calling-reports.component.css']
})
export class CallingReportsComponent implements OnInit {

  //modal
  modalRef: any;
  @ViewChild('callingReportModal') callingReportModal: any
  @ViewChild('detailedCallingReportModal') detailedCallingReportModal: any

  @ViewChild('enquiryCallingReportModal') enquiryCallingReportModal: any
  @ViewChild('detailedEnquiryCallingReportModal') detailedEnquiryCallingReportModal: any

  loginwiseCallingReports: ILoginwiseCallingReport[]
  daywiseCallingReports: IDaywiseCallingReport[]
  totalLoginwiseReports: number
  totalCallingReports: number
  loginwiseLimit: number
  loginwiseOffset: number


  callingReportList: any[]
  callingReportSearchForm: FormGroup
  detailedCallingReportForm: FormGroup
  loginwiseCallingReportSearchForm: FormGroup
  callingReportSearchFormValue: any
  isCallingReportSearched: boolean
  callingReportOffset: number
  callingReportLimit: number
  currentCallingReportPage: number
  paginationStart: number
  paginationEnd: number

  loginwiseCallingReportSearchFormValue: any
  isLoginwiseSearched: boolean
  showCallingReportSearch: boolean

  purposeList: any[]
  outcomeList: any[]
  userCredentialList: any[]

  // talent-enquiry forms

  loginwiseEnquiryCallingReports: ILoginwiseEnquiryCallingReport[]
  daywiseEnquiryCallingReports: IDaywiseEnquiryCallingReport[]
  enquiryCallingReportList: any[]

  totalLoginwiseEnquiryReports: number
  totalEnquiryCallingReports: number
  loginwiseEnquiryLimit: number
  loginwiseEnquiryOffset: number
  enquiryCallingReportOffset: number
  enquiryCallingReportLimit: number
  currentEnquiryCallingReportPage: number

  isLoginwiseEnquirySearched: boolean
  isEnquiryCallingReportSearched: boolean
  showEnquiryCallingReportSearch: boolean

  loginwiseEnquiryCallingReportSearchFormValue: any
  enquiryCallingReportSearchFormValue: any

  loginwiseEnquiryCallingReportSearchForm: FormGroup
  enquiryCallingReportSearchForm: FormGroup
  enquiryDetailedCallingReportForm: FormGroup

  constructor(
    private formBuilder: FormBuilder,
    private callingReportService: CallingReportService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private generalService: GeneralService,
    private urlConstant: UrlConstant
  ) {

    this.getAllComponents()
    this.initializeVariables()
    this.createForms()
  }

  getAllComponents() {
    this.getLoginwiseCallingReports()
    this.getUserCredentialList()
    // this.getPurposeList()
  }

  initializeVariables() {
    this.isCallingReportSearched = false
    this.callingReportOffset = 0
    this.callingReportLimit = 5

    this.loginwiseCallingReportSearchFormValue = {}
    this.isLoginwiseSearched = false
    this.showCallingReportSearch = false

    this.isLoginwiseEnquirySearched = false
    this.isEnquiryCallingReportSearched = false
    this.showEnquiryCallingReportSearch = false

    this.loginwiseEnquiryCallingReportSearchFormValue = {}
    this.enquiryCallingReportSearchFormValue = {}

    this.enquiryCallingReportLimit = 5
    this.enquiryCallingReportOffset = 0
  }

  createForms() {
    this.createCallingSearchForm()
    this.createDetailedReportForm()
    this.createLoginwiseCallingSearchForm()

    // talent-enquiry
    this.createEnquiryCallingSearchForm()
    this.createDetailedEnquiryReportForm()
    this.createLoginwiseEnquiryCallingSearchForm()
  }

  createCallingSearchForm() {
    this.callingReportSearchForm = this.formBuilder.group({
      credentialID: new FormControl(null),
      purposeID: new FormControl(null),
      outcomeID: new FormControl(null),
      targetDate: new FormControl(null),
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
      date: new FormControl(null),
      targetFromDate: new FormControl(null),
      targetToDate: new FormControl(null),
      noticePeriod: new FormControl(null),
      expectedCtc: new FormControl(null)
    })
    // // as reset does not disable the outcome selection
    // this.callingReportSearchForm.get('outcomeID').disable()
    this.outcomeList = null
  }

  createLoginwiseCallingSearchForm() {
    this.loginwiseCallingReportSearchForm = this.formBuilder.group({
      duration: new FormControl(null),
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
    })
  }

  createDetailedReportForm() {
    this.detailedCallingReportForm = this.formBuilder.group({
      comment: new FormControl('Not mentioned'),
      dateTime: new FormControl(null),
      loginName: new FormControl(null),
      outcome: new FormControl(null),
      purpose: new FormControl(null),
      expectedCTC: new FormControl('Not mentioned'),
      noticePeriod: new FormControl('Not mentioned'),
      targetDate: new FormControl('Not mentioned'),
      talent: new FormGroup({
        contact: new FormControl(null),
        email: new FormControl(null),
        firstName: new FormControl(null),
        lastName: new FormControl(null)
      })
    })
  }

  searchCallingReports() {
    this.callingReportSearchFormValue = { ...this.callingReportSearchForm.value }
    let flag: boolean = true
    for (let field in this.callingReportSearchFormValue) {
      if (!this.callingReportSearchFormValue[field] || this.callingReportSearchFormValue[field] === "") {
        delete this.callingReportSearchFormValue[field]
      } else {
        this.isCallingReportSearched = true
        flag = false
      }
    }
    // No API call on empty search.
    if (flag) {
      return
    }
    this.spinnerService.loadingMessage = "Getting calling reports"
    this.changeCallingReportPage(1)
  }

  resetLoginwiseCallingReportSearchForm(): void {
    this.loginwiseCallingReportSearchForm.reset()
  }

  resetCallReportSearchForm(): void {
    this.callingReportSearchForm.reset()
  }

  resetLoginwiseSearchAndGetAll() {
    this.isLoginwiseSearched = false
    this.loginwiseCallingReportSearchForm.reset()
    this.loginwiseCallingReportSearchFormValue = {}
    this.getLoginwiseCallingReports()
  }

  resetCallReportSearchAndGetAll() {
    this.isCallingReportSearched = false
    this.callingReportSearchForm.reset()
    this.callingReportSearchFormValue = {}
    this.getCallingReports()
  }

  searchLoginwiseCallingReports() {
    this.loginwiseCallingReportSearchFormValue = { ...this.loginwiseCallingReportSearchForm.value }
    let flag: boolean = true
    for (let field in this.loginwiseCallingReportSearchFormValue) {
      if (!this.loginwiseCallingReportSearchFormValue[field] ||
        this.loginwiseCallingReportSearchFormValue[field] === "") {
        delete this.loginwiseCallingReportSearchFormValue[field]
      } else {
        flag = false
        this.isLoginwiseSearched = true
      }
    }
    // No API call on empty search.
    if (flag) {
      return
    }
    this.spinnerService.loadingMessage = "Getting calling reports"
    this.getLoginwiseCallingReports()
  }

  // Handles pagination
  changeCallingReportPage($event: any): void {
    // $event will be the page number & offset will be 1 less than it.
    this.callingReportOffset = $event - 1
    this.currentCallingReportPage = $event
    this.getCallingReports()
  }

  //Get purpose list
  getPurposeList(type: string) {
    this.generalService.getPurposeListByType(type).subscribe(data => {
      this.purposeList = data;
    }, err => {
      console.error(err)
    })
  }

  //Get purpose list
  getUserCredentialList() {
    this.generalService.getUserCredentialList().subscribe(data => {
      this.userCredentialList = data;
    }, err => {
      console.error(err)
    })
  }

  //Get outcome list by purpose id list
  getOutcomesByPurposeID(callingReportForm: any) {
    this.outcomeList = null;
    callingReportForm.get('outcomeID').setValue(null)
    let purposeID = callingReportForm.get('purposeID').value
    if (purposeID == null) {
      return
    }
    this.generalService.getOutcomeListByPurpose(purposeID).subscribe((data: any) => {
      callingReportForm.get('outcomeID').enable()
      this.outcomeList = data
    }, (err) => {
      console.error(err)
    })
  }

  getLoginwiseCallingReports() {
    this.spinnerService.loadingMessage = "Getting Loginwise reports"

    this.callingReportService.getLoginwiseCallingReports(this.loginwiseCallingReportSearchFormValue).
      subscribe((data) => {

        this.loginwiseCallingReports = data;
      }, (error) => {

        console.error(error);
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
        }
      });
  }

  getAllDaywiseCallingReports() {
    this.spinnerService.loadingMessage = "Getting Daywise reports"

    this.callingReportService.getDaywiseCallingReports().
      subscribe((data) => {

        this.daywiseCallingReports = data;
      }, (error) => {

        console.error(error);
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
        }
      });
  }

  onViewCallReportsClick(form: { controlName: any, value: any }) {
    this.getPurposeList("talent")
    this.callingReportSearchForm.get(form.controlName).patchValue(form.value)
    this.searchCallingReports()
    this.isCallingReportSearched = false
    this.openModal(this.callingReportModal, 'xl').result.catch(() => {
      this.callingReportSearchForm.reset()
      // this.callingReportSearchFormValue = {}
    })
  }

  onDetailedReportViewClick(report: ITalentCallingReport) {
    this.getPurposeList("talent")
    report.dateTime = formatDate(report.dateTime, 'medium', 'en_US')
    this.detailedCallingReportForm.patchValue(report)
    this.openModal(this.detailedCallingReportModal).result.catch(() => {
      // Reset form on modal close. 
      this.createDetailedReportForm()
    })
  }

  getCallingReports() {
    this.spinnerService.loadingMessage = "Getting Calling reports"

    this.callingReportService.getCallingReports(this.callingReportLimit, this.callingReportOffset, this.callingReportSearchFormValue).
      subscribe((data) => {

        this.totalCallingReports = data.headers.get('X-Total-Count');
        this.setPaginationString(this.totalCallingReports, this.callingReportLimit, this.callingReportOffset)
        this.callingReportList = data.body;
      }, (error) => {

        this.totalCallingReports = 0
        this.setPaginationString(this.totalCallingReports, this.callingReportLimit, this.callingReportOffset)
        this.callingReportList = []
        console.error(error);
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
        }
      });
  }

  openModal(content: any, size?: string): NgbModalRef {
    if (!size) {
      size = 'lg'
    }
    let options: NgbModalOptions = { ariaLabelledBy: 'modal-basic-title', keyboard: false, backdrop: 'static', size: size }
    return this.modalService.open(content, options)
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  // Set total list on current page.
  setPaginationString(total: number, limit: number, offset: number): void {
    this.paginationStart = limit * offset + 1
    this.paginationEnd = +limit + limit * offset
    if (total < this.paginationEnd) {
      this.paginationEnd = total
    }
  }

  // =======================================================TALENT-ENQUIRY=======================================================

  createEnquiryCallingSearchForm() {
    this.enquiryCallingReportSearchForm = this.formBuilder.group({
      credentialID: new FormControl(null),
      purposeID: new FormControl(null),
      outcomeID: new FormControl(null),
      targetDate: new FormControl(null),
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
      date: new FormControl(null),
      targetFromDate: new FormControl(null),
      targetToDate: new FormControl(null),
      noticePeriod: new FormControl(null),
      expectedCtc: new FormControl(null)
    })
    // // as reset does not disable the outcome selection
    // this.enquiryCallingReportSearchForm.get('outcomeID').disable()
    this.outcomeList = null
  }

  createLoginwiseEnquiryCallingSearchForm() {
    this.loginwiseEnquiryCallingReportSearchForm = this.formBuilder.group({
      duration: new FormControl(null),
      fromDate: new FormControl(null),
      toDate: new FormControl(null),
    })
  }

  createDetailedEnquiryReportForm() {
    this.enquiryDetailedCallingReportForm = this.formBuilder.group({
      comment: new FormControl('Not mentioned'),
      dateTime: new FormControl(null),
      loginName: new FormControl(null),
      outcome: new FormControl(null),
      purpose: new FormControl(null),
      expectedCTC: new FormControl('Not mentioned'),
      noticePeriod: new FormControl('Not mentioned'),
      targetDate: new FormControl('Not mentioned'),
      talentEnquiry: new FormGroup({
        contact: new FormControl(null),
        email: new FormControl(null),
        firstName: new FormControl(null),
        lastName: new FormControl(null)
      })
    })
  }

  // Handles pagination
  changeEnquiryCallingReportPage($event: any): void {
    // $event will be the page number & offset will be 1 less than it.
    this.enquiryCallingReportOffset = $event - 1
    this.currentEnquiryCallingReportPage = $event
    this.getEnquiryCallingReports()
  }

  resetLoginwiseEnquiryCallingReportSearchForm(): void {
    this.enquiryCallingReportSearchForm.reset()
  }

  resetEnquiryCallReportSearchForm(): void {
    this.enquiryCallingReportSearchForm.reset()
  }

  resetLoginwiseEnquirySearchAndGetAll() {
    this.isLoginwiseEnquirySearched = false
    this.loginwiseEnquiryCallingReportSearchForm.reset()
    this.loginwiseEnquiryCallingReportSearchFormValue = {}
    this.getLoginwiseEnquiryCallingReports()
  }

  resetEnquiryCallReportSearchAndGetAll() {
    this.isEnquiryCallingReportSearched = false
    this.enquiryCallingReportSearchForm.reset()
    this.enquiryCallingReportSearchFormValue = {}
    this.getEnquiryCallingReports()
  }

  searchEnquiryCallingReports() {
    this.enquiryCallingReportSearchFormValue = { ...this.enquiryCallingReportSearchForm.value }
    let flag: boolean = true
    for (let field in this.enquiryCallingReportSearchFormValue) {
      if (!this.enquiryCallingReportSearchFormValue[field] || this.enquiryCallingReportSearchFormValue[field] === "") {
        delete this.enquiryCallingReportSearchFormValue[field]
      } else {
        this.isEnquiryCallingReportSearched = true
        flag = false
      }
    }
    // No API call on empty search.
    if (flag) {
      return
    }
    this.spinnerService.loadingMessage = "Getting calling reports"
    this.changeEnquiryCallingReportPage(1)
  }

  searchLoginwiseEnquiryCallingReports() {
    this.loginwiseEnquiryCallingReportSearchFormValue = { ...this.loginwiseEnquiryCallingReportSearchForm.value }
    let flag: boolean = true
    for (let field in this.loginwiseEnquiryCallingReportSearchFormValue) {
      if (!this.loginwiseEnquiryCallingReportSearchFormValue[field] ||
        this.loginwiseEnquiryCallingReportSearchFormValue[field] === "") {
        delete this.loginwiseEnquiryCallingReportSearchFormValue[field]
      } else {
        flag = false
        this.isLoginwiseEnquirySearched = true
      }
    }
    // No API call on empty search.
    if (flag) {
      return
    }
    this.spinnerService.loadingMessage = "Getting calling reports"
    this.getLoginwiseEnquiryCallingReports()
  }

  onViewEnquiryCallReportsClick(form: { controlName: any, value: any }) {
    this.getPurposeList("talent_enquiry")
    this.enquiryCallingReportSearchForm.get(form.controlName).patchValue(form.value)
    this.searchEnquiryCallingReports()
    this.isEnquiryCallingReportSearched = false
    this.openModal(this.enquiryCallingReportModal, 'xl').result.catch(() => {
      this.enquiryCallingReportSearchForm.reset()
      // this.callingReportSearchFormValue = {}
    })
  }

  onEnquiryDetailedReportViewClick(report: ITalentEnquiryCallingReport) {
    console.log("whfkwfnkwnfkwfn")
    this.getPurposeList("talent_enquiry")
    report.dateTime = formatDate(report.dateTime, 'medium', 'en_US')
    this.enquiryDetailedCallingReportForm.patchValue(report)
    this.openModal(this.detailedEnquiryCallingReportModal).result.catch(() => {
      // Reset form on modal close. 
      this.createDetailedEnquiryReportForm()
    })
  }

  getLoginwiseEnquiryCallingReports() {
    this.spinnerService.loadingMessage = "Getting loginwise reports"

    this.callingReportService.getLoginwiseEnquiryCallingReports(this.loginwiseEnquiryCallingReportSearchFormValue).
      subscribe((data) => {

        this.loginwiseEnquiryCallingReports = data;
      }, (error) => {

        console.error(error);
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
        }
      });
  }

  getDaywiseEnquiryCallingReports() {
    this.spinnerService.loadingMessage = "Getting daywise reports"

    this.callingReportService.getDaywiseEnquiryCallingReports().
      subscribe((data) => {

        this.daywiseEnquiryCallingReports = data;
      }, (error) => {

        console.error(error);
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
        }
      });
  }

  getEnquiryCallingReports() {
    this.spinnerService.loadingMessage = "Getting calling reports"

    this.callingReportService.getEnquiryCallingReports(this.enquiryCallingReportLimit, this.enquiryCallingReportOffset, this.enquiryCallingReportSearchFormValue).
      subscribe((data) => {

        this.totalEnquiryCallingReports = data.headers.get('X-Total-Count');
        this.setPaginationString(this.totalEnquiryCallingReports, this.enquiryCallingReportLimit, this.enquiryCallingReportOffset)
        this.enquiryCallingReportList = data.body;
      }, (error) => {

        this.totalEnquiryCallingReports = 0
        this.setPaginationString(this.totalEnquiryCallingReports, this.enquiryCallingReportLimit, this.enquiryCallingReportOffset)
        this.enquiryCallingReportList = []
        console.error(error);
        if (error.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
        }
      });
  }

  onTabChange(event: any) {
    switch (event) {
      case 1:
        this.getLoginwiseCallingReports()
        break;
      case 2:
        this.getAllDaywiseCallingReports()
        break
      case 3:
        this.getLoginwiseEnquiryCallingReports()
        break
      case 4:
        this.getDaywiseEnquiryCallingReports()
        break
    }
    // if (event == 1) {
    // }
    // if (event == 2) {
    //   this.getAllDaywiseCallingReports()
    // }
  }
}


interface ILoginwiseCallingReport {
  credentialID: string
  firstName: string
  lastName: string
  totalCallingCount: number
  totalTalentCount: number
}

interface IDaywiseCallingReport {
  date: string
  totalCallingCount: number
  totalTalentCount: number
}

interface ITalentCallingReport {
  comment: string,
  dateTime: string,
  loginName: string,
  outcome: string,
  purpose: string,
  talent: ITalent
}

interface ITalent {
  contact: string,
  email: string,
  firstName: string,
  lastName: string
}