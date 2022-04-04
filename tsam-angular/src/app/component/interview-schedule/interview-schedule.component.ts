import { Component, OnInit, ViewChild } from '@angular/core';
import { IInterview, InterviewScheduleService } from 'src/app/service/talent/interview-schedule.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { ActivatedRoute } from '@angular/router';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { GeneralService } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { UrlConstant } from 'src/app/service/constant';
import { DatePipe, Location } from '@angular/common';

@Component({
  selector: 'app-interview-schedule',
  templateUrl: './interview-schedule.component.html',
  styleUrls: ['./interview-schedule.component.css']
})
export class InterviewScheduleComponent implements OnInit {

  //***********************************INTERVIEW SCHEDULE*************************************** */
  // Components.
  interviewScheduleStatusList: any[]

  // Flags.
  isViewMode: boolean
  isOperationInterviewScheduleUpdate: boolean

  // Interview Schedule.
  interviewSchedulesList: any[]
  interviewScheduleForm: FormGroup
  talentID: string
  talentName: string
  selectedInterviewScheduleID: string

  // Modal.
  modalRef: any
  @ViewChild('interviewScheduleFormModal') interviewScheduleFormModal: any
  @ViewChild('deleteInterviewScheduleModal') deleteInterviewScheduleModal: any

  // Spinner.


  // Pagination.
  totalInterviewSchedules: number

  // Permissions.
  permission: IPermission
  roleName: string

  //***********************************INTERVIEW*************************************** */
  // Components.
  interviewRoundList: any[]
  ratingList: any[]
  interviewerList: any[]
  selectedInterviewerList: any[]
  interviewStatusList: any[]

  // Flags.
  showInterviewForm: boolean
  isOperationInterviewUpdate: boolean

  //Interview.
  interviewList: any[]
  interviewForm: FormGroup

  // Modal.
  @ViewChild('interviewFormModal') interviewFormModal: any

  constructor(
    private interviewScheduleService: InterviewScheduleService,
    private spinnerService: SpinnerService,
    private activatedRoute: ActivatedRoute,
    private formBuilder: FormBuilder,
    private modalService: NgbModal,
    private generalService: GeneralService,
    public utilityService: UtilityService,
    private localService: LocalService,
    private urlConstant: UrlConstant,
    private _location: Location,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  initializeVariables(): void {
    // Components.
    this.interviewScheduleStatusList = []
    this.interviewSchedulesList = []
    this.interviewRoundList = []
    this.ratingList = []
    this.interviewerList = []
    this.interviewStatusList = []
    this.interviewList = []
    this.selectedInterviewerList = []

    // Flags.
    this.isViewMode = false
    this.isOperationInterviewScheduleUpdate = false
    this.isOperationInterviewUpdate = false
    this.showInterviewForm = false

    // Initialize forms.
    this.createInterviewScheduleForm()

    // Spinner.
    this.spinnerService.loadingMessage = "Getting all interview schedules"

    // Permision.
    // Get permissions from menus using utilityService function.
    this.permission = this.utilityService.getPermission(this.urlConstant.INTERVIEW_SCHEDULE)

    // Get role name for menu for calling their specific apis.
    this.roleName = this.localService.getJsonValue("roleName")

    // Get talent id from params.
    if (this.activatedRoute.snapshot.queryParamMap.keys != []) {
      this.talentName = this.activatedRoute.snapshot.queryParamMap.get("talentName")
      this.talentID = this.activatedRoute.snapshot.queryParamMap.get("talentID")
      if (this.talentID) {

        this.getInterviewSchedulesByTalent()
        return
      }
    }
  }

  //*********************************************CREATE FORMS************************************************************
  // Create new interview schedule form.
  createInterviewScheduleForm(): void {
    this.interviewScheduleForm = this.formBuilder.group({
      id: new FormControl(),
      scheduledDate: new FormControl(null, [Validators.required]),
      status: new FormControl("Scheduled", [Validators.required]),
    })
  }

  // Create new interview form.
  createInterviewForm(): void {
    this.interviewForm = this.formBuilder.group({
      id: new FormControl(),
      talentID: new FormControl(null),
      roundID: new FormControl(null, [Validators.required]),
      rating: new FormControl(null, [Validators.required]),
      status: new FormControl(null, [Validators.required]),
      comment: new FormControl(null, [Validators.maxLength(1000)]),
      takenBy: new FormControl(Array(), [Validators.required])
    })
    this.getInterviewerbyRound(this.interviewForm.get('roundID').value)
  }

  //*********************************************ADD FOR INTERVIEW SCHEDULE FUNCTIONS************************************************************
  // On add new interview schedule button click.
  onAddNewInterviewScheduleButtonClick() {
    this.isViewMode = false
    this.isOperationInterviewScheduleUpdate = false
    this.createInterviewScheduleForm()
    this.openModal(this.interviewScheduleFormModal, 'md')
  }

  // Add new interview schedule.
  addInterviewSchedule() {
    this.spinnerService.loadingMessage = "Adding Interview Schedule"

    this.interviewScheduleService.addInterviewSchedule(this.interviewScheduleForm.value, this.talentID).subscribe((response: any) => {
      this.modalRef.close('success')
      this.getInterviewSchedulesByTalent()
      alert(response)
    }, (error) => {
      console.error(error)

      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  //*********************************************UPDATE AND VIEW FOR INTERVIEW SCHEDULE FUNCTIONS************************************************************
  // On clicking view interview schedule button.
  onViewInterviewScheduleClick(interviewSchedule: any) {
    this.isViewMode = true
    this.createInterviewScheduleForm()
    this.interviewScheduleForm.patchValue(interviewSchedule)
    this.interviewScheduleForm.disable()
    this.openModal(this.interviewScheduleFormModal, 'md')
  }

  // On clicking update interview schedule button.
  onUpdateInterviewScheduleClick() {
    this.isViewMode = false
    this.isOperationInterviewScheduleUpdate = true
    this.interviewScheduleForm.enable()
  }

  // Update interview schedule.
  updateInterviewSchedule() {
    this.spinnerService.loadingMessage = "Updating Interview Schedule"

    this.interviewScheduleService.updateInterviewSchedule(this.interviewScheduleForm.value, this.talentID).subscribe((response: any) => {

      this.modalRef.close('success')
      this.getInterviewSchedulesByTalent()
      alert(response)
    }, (error) => {
      console.error(error)

      if (error.error?.error) {
        alert(error.error.error)
        return
      }
      alert("Check connection")
    })
  }

  //*********************************************DELETE FOR INTERVIEW SCHEDULE FUNCTIONS************************************************************
  // On clicking delete interview schedule button. 
  onDeleteInterviewScheduleClick(designationID: string): void {
    this.openModal(this.deleteInterviewScheduleModal, 'md').result.then(() => {
      this.deleteInterviewSchedule(designationID)
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Delete interview schedule.
  deleteInterviewSchedule(designationID: string) {
    this.spinnerService.loadingMessage = "Deleting Interview Schedule"
    this.modalRef.close()

    this.interviewScheduleService.deleteInterviewSchedule(designationID, this.talentID).subscribe((response: any) => {
      this.getInterviewSchedulesByTalent()
      alert(response)
    }, (error) => {
      console.error(error)

      if (error.error) {
        alert(error.error)
        return
      }
      alert(error.statusText)
    })
  }

  //*********************************************CRUD FUNCTIONS FOR INTERVIEW************************************************************
  // On clickig interviews button in interview schedule list, get all interviews by interview schdeule id.
  getIntreviewsForSelectedInterviewSchedule(interviewScheduleID: string): void {
    this.spinnerService.loadingMessage = "Getting all interviews"

    this.selectedInterviewScheduleID = interviewScheduleID
    this.getAllInterviews()
    this.showInterviewForm = false
    this.openModal(this.interviewFormModal, 'lg')
  }

  // Get all interviews by interview schdeule id.
  getAllInterviews(): void {
    this.interviewScheduleService.getInterviewsByInterviewSchedule(this.selectedInterviewScheduleID).subscribe(response => {
      this.interviewList = response
      this.convertRatingFromNumberToString(this.interviewList)

    }, error => {
      console.error(error)

      if (error.error) {
        alert(error.error)
        return
      }
      alert(error.statusText)
    })
  }

  // Convert rating from number to string.
  convertRatingFromNumberToString(interviews: any[]): void {
    for (let i = 0; i < interviews.length; i++) {
      interviews[i].rating = interviews[i].rating + ""
    }
  }

  // On clicking add new interview button.
  onAddNewInterviewButtonClick(): void {
    this.createInterviewForm()
    this.isOperationInterviewUpdate = false
    this.showInterviewForm = true
  }

  // Add interview.
  addInterview() {
    let interview: IInterview = this.interviewForm.value
    interview.rating = +interview.rating
    interview.talentID = this.talentID

    this.spinnerService.loadingMessage = "Adding Interview"
    this.interviewScheduleService.addInterview(interview, this.selectedInterviewScheduleID).subscribe((response: any) => {
      this.showInterviewForm = false
      this.getAllInterviews()
      alert(response)
    }, (error) => {
      console.error(error)

      if (error.error) {
        alert(error.error)
        return
      }
      alert(error.statusText)
    })
  }

  // On clicking update interview button.
  OnUpdateInterviewButtonClick(index: number) {
    this.createInterviewForm()
    this.isOperationInterviewUpdate = true
    this.showInterviewForm = true
    this.interviewForm.patchValue(this.interviewList[index])
  }

  // Update interview.
  updateInterview() {
    let interview: IInterview = this.interviewForm.value
    interview.rating = +interview.rating
    interview.talentID = this.talentID

    this.spinnerService.loadingMessage = "Updating Interview"
    this.interviewScheduleService.updateInterview(this.interviewForm.value, this.selectedInterviewScheduleID).subscribe((response: any) => {
      this.showInterviewForm = false
      this.getAllInterviews()
      alert(response)
      this.interviewForm.reset()
    }, (error) => {
      console.error(error)

      if (error.error) {
        alert(error.error)
        return
      }
      alert(error.statusText)
    })
  }

  // Delete interview.
  deleteInterview(interviewID: string) {
    if (confirm("Are you sure you want to delete the interview?")) {

      this.spinnerService.loadingMessage = "Deleting Interview"
      this.interviewScheduleService.deleteInterview(interviewID, this.selectedInterviewScheduleID).subscribe((response: any) => {
        this.interviewForm.reset()
        this.getAllInterviews()
        alert(response)
      }, (error) => {
        console.error(error)

        if (error.error) {
          alert(error.error)
          return
        }
        alert(error.statusText)
      })
    }
  }

  // Validate interview form.
  validateInterviewForm() {
    if (this.interviewForm.invalid) {
      this.interviewForm.markAllAsTouched()
      return
    }
    if (this.isOperationInterviewUpdate) {
      this.updateInterview()
      return
    }
    this.addInterview()
  }

  // Get roung name by round id.
  getRoundNameByRoundID(roundID: string): void {
    for (let i = 0; i < this.interviewRoundList.length; i++) {
      if (this.interviewRoundList[i].id == roundID) {
        return this.interviewRoundList[i].name
      }
    }
  }

  //*********************************************FUNCTIONS FOR INTERVIEW SCHEDULE FORM************************************************************
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

  // Validate interview schedule form.
  onFormSubmit(): void {
    if (this.interviewScheduleForm.invalid) {
      this.interviewScheduleForm.markAllAsTouched()
      return
    }
    if (this.isOperationInterviewScheduleUpdate) {
      this.updateInterviewSchedule()
      return
    }
    this.addInterviewSchedule()
  }

  // Nagivate to previous url.
  backToPreviousPage(): void {
    this._location.back()
  }

  //*********************************************GET FUNCTIONS************************************************************
  // Get all lists.
  getAllComponents() {
    this.getInterviewRounds()
    this.getInterviewScheduleStatusList()
    this.getInterviewStatusList()
    this.getRatingList()
    this.getAllInterviewersList()
  }

  // Get all interview schedules by talent id.
  getInterviewSchedulesByTalent() {
    this.spinnerService.loadingMessage = "Getting all interview schedules"

    this.interviewScheduleService.getInterviewSchedulesByTalent(this.talentID).subscribe(response => {
      this.interviewSchedulesList = response
      this.totalInterviewSchedules = this.interviewSchedulesList.length

    }, err => {
      console.error(err)

    })
  }

  // Get all interview rounds.
  getInterviewRounds() {
    this.interviewScheduleService.getInterviewRounds().subscribe(response => {
      this.interviewRoundList = response

    }, err => {
      console.error(err)
    })
  }

  // Get interview schedule status list.
  getInterviewScheduleStatusList() {
    this.generalService.getGeneralTypeByType("interview_schedule_status").subscribe((response: any[]) => {
      this.interviewScheduleStatusList = response
    }, (err) => {
      console.error(err)
    })
  }

  // Get interview status list.
  getInterviewStatusList() {
    this.generalService.getGeneralTypeByType("interview_status").subscribe((response: any[]) => {
      this.interviewStatusList = response
    }, (err) => {
      console.error(err)
    })
  }

  // Get rating list.
  getRatingList() {
    this.generalService.getGeneralTypeByType("interview_rating").subscribe((response: any[]) => {
      this.ratingList = response
    }, (err) => {
      console.error(err)
    })
  }

  // Get all interviewers list.
  getAllInterviewersList() {
    this.interviewScheduleService.getAllInterviewersList().subscribe(response => {
      this.interviewerList = response

    }, err => {
      console.error(err)
    })
  }

  // Get Interviewer list by round.
  getInterviewerbyRound(roundID: string): void {
    if (roundID == null) {
          this.selectedInterviewerList = []
          this.interviewForm.get('takenBy').setValue(null)
          this.interviewForm.get('takenBy').disable()
          return
    }else{
      this.interviewForm.get('takenBy').setValue(null)
      this.interviewForm.get('takenBy').enable()
      this.selectedInterviewerList = []
      let whichRound :string =""
      for(let j=0;j<this.interviewRoundList.length;j++){
        if(roundID == this.interviewRoundList[j].id){
          whichRound = this.interviewRoundList[j].name
          break;
        }

      }
      if( whichRound === "HR round" || whichRound === "Company round" ){
        for(let i=0;i<this.interviewerList.length;i++){
          if("SalesPerson" === this.interviewerList[i].role.roleName){
            this.selectedInterviewerList.push(this.interviewerList[i])
          }
        }
      }else{
        for(let i=0;i<this.interviewerList.length;i++){
          if("SalesPerson" != this.interviewerList[i].role.roleName){
            this.selectedInterviewerList.push(this.interviewerList[i])
          }
        }
      }
      
    }
  }

}
