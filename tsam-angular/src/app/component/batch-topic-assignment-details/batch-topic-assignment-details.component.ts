import { HttpParams } from '@angular/common/http';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { NgAudioRecorderService, OutputFormat, RecorderState } from 'ng-audio-recorder';
import { IAssignmentSubmission } from 'src/app/models/assignment-submission';
import { IBatchTopicAssignment } from 'src/app/models/batch-topic-assignment';
import { IModuleConcept } from 'src/app/models/module-concept';
import { ITalent } from 'src/app/models/talent/talent';
import { ITalentAssignmentSubmission } from 'src/app/models/talent-assignment-submission';
import { AssignmentData as AssignmentData } from 'src/app/providers/assignment-data';
import { BatchTalentService } from 'src/app/service/batch-talent/batch-talent.service';
import { BatchTopicAssignmentService } from 'src/app/service/batch-topic-assignment/batch-topic-assignment.service';
import { ConceptModuleService } from 'src/app/service/concept-module/concept-module.service';
import { BackNavigationUrl, Role, UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UrlService } from 'src/app/service/url.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-batch-topic-assignment-details',
  templateUrl: './batch-topic-assignment-details.component.html',
  styleUrls: ['./batch-topic-assignment-details.component.css']
})
export class BatchTopicAssignmentDetailsComponent implements OnInit {


  // Flags.
  isQuestionSelected: boolean = true
  uncheckedTalents: ITalent[] = []
  pendingTalents: ITalent[] = []
  completedTalents: ITalent[] = []
  uncheckedAssignments: IBatchTopicAssignment[] = []
  completedAssignments: IBatchTopicAssignment[] = []
  pendingAssignments: IBatchTopicAssignment[] = []
  allAssignments: IBatchTopicAssignment[] = []


  batchID: string
  selectedTalent: ITalent = {} as ITalent
  selectedAssignment: IBatchTopicAssignment = {} as IBatchTopicAssignment
  assignmentID: string
  talentID: string
  isTalentView: boolean = false
  allAssignmentSubmissions: IAssignmentSubmission[] = []
  batchTalents: ITalent[] = []

  batchTopicAssignments: IBatchTopicAssignment[] = []

  talentAllSubmissions: ITalentAssignmentSubmission[] = []
  moduleConcepts: IModuleConcept[] = []

  showCompleted: boolean = false
  showCompletedAssignments: boolean = false
  showPending: boolean = false
  showPendingAssignments: boolean = false
  showUnchecked: boolean = false
  showUncheckedAssignments: boolean = false
  // talentAllSubmissionsMap: Map<string,TalentAssignmentSubmission[]> = new Map
  // talentAssignmentSubmissionsMap: Map<string,Map<string,TalentAssignmentSubmission[]>> = new Map

  permission: IPermission

  // Modal.
  modalRef: NgbModalRef
  @ViewChild('commentModal') commentModal: any
  @ViewChild('scoreModal') scoreModal: any

  constructor(
    private modalService: NgbModal,
    private activatedRoute: ActivatedRoute,
    private router: Router,
    private formBuilder: FormBuilder,
    private assignmentProvider: AssignmentData,
    private batchTalentService: BatchTalentService,
    private urlService: UrlService,
    private conceptModuleService: ConceptModuleService,
    private btaService: BatchTopicAssignmentService,
    private localService: LocalService,
    private utilService: UtilityService,
    private audioRecorderService: NgAudioRecorderService,
    private sanitizer: DomSanitizer,
    private fileOperationService: FileOperationService,
    private urlConstant: UrlConstant,
    private backNavigation: BackNavigationUrl,
		private role: Role,
    private spinnerService: SpinnerService,
  ) {
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN || this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER_ASSIGNMENT_DETAILS)
		}
		if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
			this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH_ASSIGNMENT_DETAILS)
		}
    this.extractData()
    this.initialize()
  }

  goBack(): void {
    this.urlService.goBack(this.backNavigation.BTA_DETAILS_TO_BTA_SCORES)
  }

  extractData(): void {
    this.allAssignmentSubmissions = this.assignmentProvider.allAssignmentSubmissions
    this.batchTalents = this.assignmentProvider.batchTalents

    this.activatedRoute.queryParamMap.subscribe(
      (params: any) => {
        this.batchID = params.get("batchID")
        this.talentID = params.get("talentID")
        this.assignmentID = params.get("assignmentID")
        let isTalView = params.get("isTalentView")
        if (isTalView != null || undefined) {
          this.isTalentView = isTalView == 'true'
        }
        let isQueSelected = params.get("isQuestionSelected")
        if (isQueSelected != null || undefined) {
          this.isQuestionSelected = isQueSelected == 'true'
        }
      }, (err: any) => {
        console.error(err);
      })
  }

  updateURL(): void {
    this.router.navigate(
      [],
      {
        relativeTo: this.activatedRoute,
        queryParams: {
          "assignmentID": this.selectedAssignment.id,
          "talentID": this.selectedTalent.id,
          "isQuestionSelected": this.isQuestionSelected,
          "isTalentView": this.isTalentView
        },
        queryParamsHandling: 'merge'
      });
  }



  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  async initialize(override: boolean = false): Promise<void> {
    this.isRecordingStored = false
    // need topic assignments for allAssignmentSubmissions.
    try {
      if (!(this.batchTopicAssignments?.length > 0) || override) {
        this.batchTopicAssignments = await this.getAllAssignmentsWithScore()
      }
      if (!(this.batchTalents?.length > 0) || override) {
        this.batchTalents = await this.getBatchTalents()
      }
    } catch (err) {
      alert(err)
      return
    }
    this.orderScores()
    this.setSelectedTalent()
    this.setSelectedAssignment()
    this.arrangeAllTalents()
    this.arrangeAllAssignments()
    if (!this.isQuestionSelected) {
      this.getAllSubmissionsForTalent()
    }
    // this.createLatestSubmissionForm()
  }

  // this will arrange all talents, categorizing them depending on their submission state: (pending,completed and unchecked)
  arrangeAllTalents(): void {
    if (this.allAssignmentSubmissions.length <= 0) {
      return
    }

    this.pendingTalents = []
    this.uncheckedTalents = []
    this.completedTalents = []

    let assignSubmission = this.allAssignmentSubmissions.find((val) => {
      if (val.assignment.id === this.selectedAssignment.id)
        return val
    })
    // if (!assignSubmission) {
    //   return
    // }
    for (let talent of this.batchTalents) {
      let value = assignSubmission.talentSubmission.get(talent.id)
      // Unchecked category.
      if (value && !value.isChecked) {
        if (this.selectedTalent.id === talent.id &&
          value?.batchTopicAssignmentID === this.selectedAssignment.id) {
          // Display the list with unchecked submissions.
          this.showUnchecked = true
          this.showStudentAnswer()
        }
        this.uncheckedTalents.push(talent)
        continue
      }
      // Completed category
      if (value && value.isAccepted) {
        if (this.selectedTalent.id === talent.id &&
          value?.batchTopicAssignmentID === this.selectedAssignment.id) {
          // Display the list with completed submissions.
          this.showCompleted = true
          this.showStudentAnswer()
        }
        this.completedTalents.push(talent)
        continue
      }
      if (this.selectedTalent.id === talent.id) {
        // Display the list with pending submissions.
        this.showPending = true
      }
      this.pendingTalents.push(talent)
    }
  }

  // this will arrange all assignments, categorizing them depending on their submission state: (pending,completed and unchecked)
  arrangeAllAssignments(): void {
    if (this.allAssignmentSubmissions.length <= 0) {
      return
    }

    this.pendingAssignments = []
    this.uncheckedAssignments = []
    this.completedAssignments = []
    this.allAssignments = []

    for (let i = 0; i < this.allAssignmentSubmissions.length; i++) {
      this.allAssignments.push(this.allAssignmentSubmissions[i].assignment)
      let value = this.allAssignmentSubmissions[i].talentSubmission.get(this.selectedTalent?.id)
      if (!value || (value.isChecked && !value.isAccepted)) {
        if (this.allAssignmentSubmissions[i].assignment.id === this.selectedAssignment.id) {
          this.showPendingAssignments = true
        }
        this.pendingAssignments.push(this.allAssignmentSubmissions[i].assignment)
        continue
      }
      if (value && !value.isChecked) {
        if (this.allAssignmentSubmissions[i].assignment.id === this.selectedAssignment.id) {
          this.showUncheckedAssignments = true
        }
        this.uncheckedAssignments.push(this.allAssignmentSubmissions[i].assignment)
        continue
      }
      if (value && value.isAccepted) {
        if (this.allAssignmentSubmissions[i].assignment.id === this.selectedAssignment.id) {
          this.showCompletedAssignments = true
        }
        this.allAssignmentSubmissions[i].assignment.submittedOn = value.submittedOn
        this.completedAssignments.push(this.allAssignmentSubmissions[i].assignment)
        continue
      }
    }
    this.sortByDueDate(this.pendingAssignments)
    this.sortByDueDate(this.uncheckedAssignments)
    this.sortByDueDate(this.completedAssignments)

    this.sortByAssignedDate(this.allAssignments)
  }

  // Redirect to external link. change #Niranjan
  redirectToexternalLink(url: string): void {
    window.open(url, "_blank");
  }

  // Used to open modal.
  openModal(content: any, options: NgbModalOptions = {
    ariaLabelledBy: 'modal-basic-title', keyboard: false,
    backdrop: 'static', size: 'lg'
  }): NgbModalRef {
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  dismissModal(reason?: any) {
    if (this.isRecording || this.isRecordingStored || this.audioRecorderService.getRecorderState() === RecorderState.RECORDING) {
      if (!confirm("Recording will be cancelled/deleted.")) {
        return
      }
      this.onDeleteRecording()
    }
    this.modalRef.dismiss(reason)
  }

  closeModal(result?: any) {
    this.modalRef.close(result)
  }

  // On clicking score.
  attemptNo: number
  selectedSubmission: ITalentAssignmentSubmission

  async openScoringModal(attemptNo: number, submission: ITalentAssignmentSubmission): Promise<void> {
    this.attemptNo = attemptNo
    this.selectedSubmission = submission
    try {
      this.moduleConcepts = await this.getModuleConceptsForAssignment()
    } catch (err) {
      alert(err)
      return
    }
    this.createLatestSubmissionForm()
    this.addTalentConceptRatingForm(this.moduleConcepts)
    this.latestSubmissionForm.patchValue(submission)
    this.calculateAssignmentScore()
    this.openModal(this.scoreModal, { size: 'md' })
  }


  @ViewChild('imagesModal') imagesModal: any

  // On clicking images.
  onImagesClick(submission: any): void {
    this.selectedSubmission = submission
    console.log(this.selectedSubmission)
    this.openModal(this.imagesModal, { keyboard: true, size: 'xl' })
  }


  onPlusButtonClick(control: FormControl) {
    // 9 becomes 10, 10 or greater remains 10
    if (control.value > 8) {
      control.setValue(10)
    } else {
      control.setValue(control.value + 1)
    }
    this.calculateAssignmentScore()
  }

  onMinusButtonClick(control: FormControl) {
    if (control.value > 10) {
      control.setValue(10)
    }
    if (control.value > 1) {
      control.setValue(control.value - 1)
    }
    this.calculateAssignmentScore()
  }

  assignmentScore: number = 0
  delayPenalty: number = 0
  calculateAssignmentScore() {
    this.assignmentScore = 0
    let outOf: number = this.moduleConcepts.length * 10

    let submittedDate: Date = new Date(this.latestSubmissionForm.get('submittedOn').value)
    let dueDate: Date = new Date(this.selectedAssignment.dueDate)
    let diff = dueDate.getTime() - submittedDate.getTime()
    let daysBetweenDates: number = Math.ceil(diff / (1000 * 60 * 60 * 24));
    if (submittedDate > dueDate) {
      // can be calculated using the above difference.
      this.delayPenalty = 20
    }

    let scored: number = 0
    for (let control of this.talentConceptRatingForm.controls) {
      if (control.get('score').invalid) {
        this.assignmentScore = 0
        return
      }
      scored += control.get('score').value
    }
    let questionScore = this.selectedAssignment?.programmingQuestion?.score
    this.assignmentScore = (scored / outOf * questionScore) * ((100 - this.delayPenalty) / 100)

    if (Number.isNaN(this.assignmentScore) || this.assignmentScore == 0) {
      this.assignmentScore = 0
      return
    }
    this.assignmentScore = Math.round(this.assignmentScore * 100 + Number.EPSILON) / 100
    this.latestSubmissionForm.get('score').setValue(this.assignmentScore)
  }

  // ----------------------------

  latestSubmissionForm: FormGroup = {} as FormGroup

  // add validation for accepted case #Niranjan
  createLatestSubmissionForm(): void {
    this.latestSubmissionForm = this.formBuilder.group({
      id: new FormControl(null),
      facultyRemarks: new FormControl(null),
      // need to add Validators.max(this.selectedAssignment.programmingQuestion.score)
      score: new FormControl(null),
      // faculty id from login #niranjan
      facultyID: this.localService.getJsonValue('loginID'),
      talentConceptRatings: new FormArray([]),
      isAccepted: new FormControl(null, Validators.required),
      isChecked: new FormControl(null),
      facultyVoiceNote: new FormControl(null),
      submittedOn: new FormControl(null),
      solution: new FormControl(null),
      githubURL: new FormControl(null),
      // assignmentSubmissionUploads: new FormArray(null),
    })
  }

  get talentConceptRatingForm(): FormArray {
    return this.latestSubmissionForm.get('talentConceptRatings') as FormArray
  }

  addTalentConceptRatingForm(moduleConcepts: IModuleConcept[]) {
    this.talentConceptRatingForm.clear()
    for (let i = 0; i < moduleConcepts.length; i++) {
      this.talentConceptRatingForm.push(this.formBuilder.group({
        id: new FormControl(null),
        // programmingConceptID: new FormControl(null, Validators.required),
        // talentID: new FormControl(null, Validators.required),
        // talentSubmissionID: new FormControl(null, Validators.required),
        programmingConceptID: new FormControl(null),
        talentID: new FormControl(null),
        talentSubmissionID: new FormControl(null),
        programmingConceptModuleID: moduleConcepts[i].id,
        score: new FormControl(null)
      }))
    }
  }

  swapView() {
    this.isTalentView = !this.isTalentView
    if (this.isTalentView) {
      this.arrangeAllTalents()
    } else {
      this.arrangeAllAssignments()
    }
    this.updateURL()
  }


  showStudentAnswer(): void {
    this.isQuestionSelected = false
    this.getAllSubmissionsForTalent()
    this.updateURL()
  }

  getAllSubmissionsForTalent(): void {
    let assignmentID: string = this.selectedAssignment.id
    let talentID: string = this.selectedTalent.id
    if (!assignmentID || !talentID) {
      return
    }
    this.spinnerService.loadingMessage = "Getting submissions"
    this.talentAllSubmissions = []

    this.btaService.getTalentAssignmentSubmissions(assignmentID, talentID).subscribe((response) => {
      this.talentAllSubmissions = response.body
      console.log("talentAllSubmissions", this.talentAllSubmissions)
      if (this.talentAllSubmissions.length > 0) {
        this.talentAllSubmissions[0].isLatestSubmission = true
      }
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  async getModuleConceptsForAssignment(): Promise<IModuleConcept[]> {
    return new Promise<IModuleConcept[]>((resolve, reject) => {
      if (this.moduleConcepts.length > 0) {
        resolve(this.moduleConcepts)
        return
      }
      this.spinnerService.loadingMessage = "Getting concepts"
      // this.moduleConcepts = []
      this.conceptModuleService.getConceptModulesForAssignment(this.selectedAssignment.id).subscribe((response) => {
        resolve(response.body)
      }, (err: any) => {
        console.error(err)
        if (err.statusText.includes('Unknown')) {
          reject("No connection to server. Check internet.")
          return
        }
        reject(err.error.error)
      })
    })
  }

  async getBatchTalents(): Promise<any> {
    return new Promise((resolve, reject) => {
      this.spinnerService.loadingMessage = "Getting talents";
      this.batchTalents = [];
      this.batchTalentService.getBatchTalentList(this.batchID).subscribe((response: any) => {
        resolve(response.body);
      }, (err: any) => {
        console.error(err);
        if (err.statusText.includes('Unknown')) {
          reject("No connection to server. Check internet.");
          return;
        }
        reject(err.error.error);
      });
    });
  }

  async onSubmit() {
    try {
      if (!this.validateCheckedSubmission()) {
        return
      }
      if (this.isRecordingStored) {
        var uploadedURL = await this.saveRecording()
      }
      // not to clear or re upload if error occurs (may use a variable to check if already uplaoded) #niranjan
      this.latestSubmissionForm.get('facultyVoiceNote').setValue(uploadedURL)
      await this.scoreAssignment()
      this.initialize(true)
      this.modalRef.close()
      alert("Submission checked")
    } catch (err) {
      alert(err)
      return
    }
  }

  validateCheckedSubmission(): boolean {
    if (!this.isRecordingStored) {
      this.latestSubmissionForm.get('facultyRemarks').setValidators([Validators.required, Validators.maxLength(1000)])
    } else {
      this.latestSubmissionForm.get('facultyRemarks').setValidators([Validators.maxLength(1000)])
    }
    if (this.latestSubmissionForm.get('isAccepted').value) {
      for (let control of this.talentConceptRatingForm.controls) {
        control.get('score').setValidators([Validators.required, Validators.max(10), Validators.min(1)])
      }
    } else {
      this.talentConceptRatingForm.clear()
    }
    this.utilService.updateValueAndValiditors(this.latestSubmissionForm)
    if (this.latestSubmissionForm.invalid) {
      console.error(this.utilService.findInvalidControlsRecursive(this.latestSubmissionForm))
      this.latestSubmissionForm.markAllAsTouched()
      return false
    }
    return true
  }

  async scoreAssignment(): Promise<any> {
    this.spinnerService.loadingMessage = "Updating assignment"
    console.log("batchID", this.batchID)
    console.log("selectedProject.id", this.selectedAssignment.id, this.selectedTalent.id)
    console.log("latestSubmissionForm", this.latestSubmissionForm.get('id').value)
    console.log("latestSubmissionForm value", this.latestSubmissionForm.value)
    return new Promise((resolve, reject) => {
      this.btaService.scoreTalentAssignment(this.selectedAssignment.id, this.selectedTalent.id,
        this.latestSubmissionForm.get('id').value, this.latestSubmissionForm.value).subscribe((response: any) => {
          resolve("Scored successfully.")
        }, (err: any) => {
          console.error(err)
          if (err.statusText.includes('Unknown')) {
            reject("No connection to server. Check internet.")
            return
          }
          reject(err.error.error)
        })
    })
  }

  async getAllAssignmentsWithScore(): Promise<IBatchTopicAssignment[]> {
    return new Promise<IBatchTopicAssignment[]>((resolve, reject) => {

      this.spinnerService.loadingMessage = "Getting assignments";
      this.batchTopicAssignments = [];
      let params = new HttpParams();
      params = params.append("facultyID", this.localService.getJsonValue("loginID"));

      this.btaService.getTopicAssignmentWithSubmissions(this.batchID, params).subscribe((response) => {
        resolve(response.body);
      }, (err: any) => {
        console.error(err);
        if (err.statusText.includes('Unknown')) {
          reject("No connection to server. Check internet.");
          return;
        }
        reject(err.error.error);
      });
    });
  }
  // this.talentService.getTopicAssignmentWithScores(this.batchID).subscribe((response: any) => {
  //   this.batchTopicAssignments = response.body
  //   // alert("#assign with score")
  //   this.temp()
  //   if (!this.isQuestionSelected) {
  //     this.getAllSubmissionsForTalent()
  //   }
  // }, (err: any) => {
  //   console.error(err)
  //   if (err.statusText.includes('Unknown')) {
  //     alert("No connection to server. Check internet.")
  //     return
  //   }
  //   alert(err.error.error)
  // }).add(() => {
  //   
  // })

  // Need to change to wait #niranjan
  // need to add latest submission condition #bug
  orderScores(): void {
    this.allAssignmentSubmissions = []
    for (let index = 0; index < this.batchTopicAssignments?.length; index++) {
      this.allAssignmentSubmissions.push({ assignment: this.batchTopicAssignments[index], talentSubmission: new Map })
      for (let talIndex = 0; talIndex < this.batchTalents.length; talIndex++) {
        let sub = this.batchTopicAssignments[index]?.submissions.find((value) => {
          if (value.talent.id === this.batchTalents[talIndex].id) {
            return value
          }
        })
        if (!sub) {
          this.allAssignmentSubmissions[index].talentSubmission.set(this.batchTalents[talIndex].id, null)
          continue
        }
        this.allAssignmentSubmissions[index].talentSubmission.set(this.batchTalents[talIndex].id, sub)
      }
    }
  }


  onAssignmentChange(assignment?: IBatchTopicAssignment): void {
    // whenever assignment changes, only then modules need to be reloaded.
    this.moduleConcepts = []
    if (assignment) {
      this.selectedAssignment = assignment
    }
    this.arrangeAllTalents()
    if (this.isTalentView && !this.isQuestionSelected) {
      this.isQuestionSelected = true
    }
    this.updateURL()
    if (!this.isQuestionSelected && !this.isTalentView) {
      this.getAllSubmissionsForTalent()
    }
  }

  // ngmodal?
  onTalentChange(talent?: ITalent): void {
    // this.selectedTalent = tal
    if (talent) {
      if (this.selectedTalent.id === talent.id) {
        return
      }
      this.selectedTalent = talent
    }
    this.arrangeAllAssignments()
    this.updateURL()
    if (!this.isQuestionSelected) {
      this.getAllSubmissionsForTalent()
    }
  }

  setSelectedTalent(): void {
    if (this.batchTalents.length <= 0) {
      return
    }
    if (!this.talentID) {
      this.selectedTalent = this.batchTalents[0]
      return
    }
    this.selectedTalent = this.batchTalents.find((val) => {
      if (val.id === this.talentID) {
        return val
      }
    })
  }

  setSelectedAssignment(): void {
    if (this.allAssignmentSubmissions.length <= 0) {
      return
    }
    if (!this.assignmentID) {
      this.selectedAssignment = this.allAssignmentSubmissions[0]?.assignment
      return
    }
    this.selectedAssignment = this.allAssignmentSubmissions?.find((val) => {
      if (val.assignment.id == this.assignmentID) {
        return val
      }
    })?.assignment
  }


  getPending(): void {
    this.pendingAssignments = []
    for (let i = 0; i < this.allAssignmentSubmissions?.length; i++) {
      let value = this.allAssignmentSubmissions[i].talentSubmission.get(this.selectedTalent.id)
      // // for (let value of this.allAssignmentSubmissions[i].talentSubmission.values()) {
      if (!value || (value.isChecked && !value.isAccepted)) {
        this.pendingAssignments.push(this.allAssignmentSubmissions[i].assignment)
      }
    }
    this.sortByDueDate(this.pendingAssignments)
  }

  getUnchecked(): void {
    this.uncheckedAssignments = []
    for (let i = 0; i < this.allAssignmentSubmissions?.length; i++) {
      let value = this.allAssignmentSubmissions[i].talentSubmission.get(this.selectedTalent.id)
      if (value && !value.isChecked) {
        this.uncheckedAssignments.push(this.allAssignmentSubmissions[i].assignment)
      }
    }
    this.sortByDueDate(this.uncheckedAssignments)

  }

  getCompleted(): void {
    this.completedAssignments = []
    for (let i = 0; i < this.allAssignmentSubmissions?.length; i++) {
      let value = this.allAssignmentSubmissions[i].talentSubmission.get(this.selectedTalent.id)
      if (value && value.isAccepted) {
        this.completedAssignments[i] = this.allAssignmentSubmissions[i].assignment
        this.completedAssignments[i].submittedOn = value.submittedOn
        // this.completedAssignments.push(this.allAssignmentSubmissions[i].assignment)
      }
    }
    this.sortByDueDate(this.completedAssignments)
  }

  sortByDueDate(assignments: IBatchTopicAssignment[], inAssending: boolean = true): void {
    assignments?.sort((val1, val2) => {
      if (val1.dueDate < val2.dueDate) {
        if (inAssending) {
          return -1
        }
        return 1
      }
      if (val1.dueDate > val2.dueDate) {
        if (inAssending) {
          return 1
        }
        return -1
      }
      return 0
    })
  }

  sortByAssignedDate(assignments: IBatchTopicAssignment[], inDescending: boolean = true): void {
    assignments?.sort((val1, val2) => {
      if (val1.assignedDate < val2.assignedDate) {
        if (inDescending) {
          return 1
        }
        return -1
      }
      if (val1.assignedDate > val2.assignedDate) {
        if (inDescending) {
          return -1
        }
        return 1
      }
      return 0
    })
  }


  // Recording
  audioUrl: SafeUrl
  isRecording: boolean = false
  time: number = 0
  timerDisplay: string = "00:00"
  file: File
  isRecordingStored: boolean = false


  onStartRecording() {
    this.audioRecorderService.startRecording();
    let totalCallbacks = 200
    let interval = setInterval(() => {
      totalCallbacks--
      if (this.audioRecorderService.getRecorderState() === RecorderState.RECORDING) {
        this.isRecording = true
        this.startTimer()
        clearInterval(interval)
        return
      }
      if (totalCallbacks <= 0) {
        if (this.audioRecorderService.getRecorderState() === RecorderState.INITIALIZING) {
          clearInterval(interval)
          alert("Need permission to start recording")
          return
        }
        clearInterval(interval)
      }
    }, 75)
  }

  sanitizeUrl(value: string): SafeUrl {
    return this.sanitizer.bypassSecurityTrustUrl(value)
  }

  onStopRecording() {
    this.isRecording = false
    this.isRecordingStored = true
    this.audioRecorderService.stopRecording(OutputFormat.WEBM_BLOB).then((output: string) => {
      var blob = new Blob([output], { type: "audio/mp3" });
      this.file = new File([blob], "my-audio.mp3");
      this.audioUrl = this.sanitizer.bypassSecurityTrustUrl(URL.createObjectURL(blob))
    }).catch(err => {
      console.error(err)
    });
  }

  onCancelRecording() {
    this.isRecording = false
    this.audioRecorderService.stopRecording(OutputFormat.WEBM_BLOB_URL).catch(err => {
      console.error(err)
    })
  }

  onDeleteRecording() {
    if (this.audioRecorderService.getRecorderState() === RecorderState.RECORDING) {
      this.onCancelRecording()
    }
    this.isRecordingStored = false
    this.audioUrl = null
    this.file = null
  }

  // starts timer and displays in view. (MM:SS)
  startTimer() {
    let x = setInterval(() => {
      if (!this.isRecording) {
        this.time = 0
        this.timerDisplay = "00:00"
        clearInterval(x)
        return
      }
      this.time++;
      this.timerDisplay = this.getDisplayTimer(this.time);
    }, 1000)
  }

  // uploads recorded audio.
  async saveRecording(): Promise<SafeUrl> {
    return new Promise((resolve, reject) => {
      this.fileOperationService.uploadFacultyVoiceNote(this.file).subscribe((uploadedURL: SafeUrl) => {
        resolve(uploadedURL)
      }, (err) => {
        reject(err)
      })
    })
  }


  getDisplayTimer(time: number): string {
    const minutes = '0' + Math.floor(time % 3600 / 60);
    const seconds = '0' + Math.floor(time % 3600 % 60);

    // MM:SS
    return `${minutes.slice(-2, -1)}${minutes.slice(-1)}:${seconds.slice(-2, -1)}${seconds.slice(-1)}`
  }
}

