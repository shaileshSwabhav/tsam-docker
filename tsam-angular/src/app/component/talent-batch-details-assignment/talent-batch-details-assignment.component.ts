import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { DomSanitizer } from '@angular/platform-browser';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchTopicAssignmentService } from 'src/app/service/batch-topic-assignment/batch-topic-assignment.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
  selector: 'app-talent-batch-details-assignment',
  templateUrl: './talent-batch-details-assignment.component.html',
  styleUrls: ['./talent-batch-details-assignment.component.css']
})
export class TalentBatchDetailsAssignmentComponent implements OnInit {

  public readonly PENDING_STATUS = "Pending"
  public readonly PENDING_ORDER = 1
  public readonly SUBMITTED_STATUS = "Submitted"
  public readonly SUBMITTED_ORDER = 2
  public readonly COMPLETED_STATUS = "Completed"
  public readonly COMPLETED_ORDER = 3

  // Batch Topic Assignment.
  batchTopicAssignmentList: any[]
  selectedBatchTopicAssignment: any
  selectedBatchTopicAssignmentID: string

  // Talent assignment submission.
  submissionForm: FormGroup
  selectedSubmission: any

  // Flags.
  isAssignmentListVisible: boolean
  isQuestionSelected: boolean

  // Spinner.



  // Talent.
  talentID: string

  // Modal.
  modalRef: any
  @ViewChild('commentModal') commentModal: any
  @ViewChild('voiceNoteModal') voiceNoteModal: any
  @ViewChild('imagesModal') imagesModal: any
  @ViewChild('submissionFormModal') submissionFormModal: any
  @ViewChild('conceptScoreModal') conceptScoreModal: any

  // Batch.
  batchID: string

  constructor(
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private batchTopicAssignmentService: BatchTopicAssignmentService,
    private route: ActivatedRoute,
    private localService: LocalService,
    private formBuilder: FormBuilder,
    private datePipe: DatePipe,
    private fileOperationService: FileOperationService,
    private router: Router,
    private domSanitizer: DomSanitizer
  ) {
    this.initializeVariables()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    // Flags.
    this.isAssignmentListVisible = true
    this.isQuestionSelected = false

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Problem Details"

    this.batchTopicAssignmentList = []

    // Talent.
    this.talentID = this.localService.getJsonValue("loginID")

    // Batch.
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      this.selectedBatchTopicAssignmentID = params.get("batchTopicAssignmentID")
      if (this.batchID) {
        this.getAllBatchTopicAssignmentsForTalent()
      }
      if (this.selectedBatchTopicAssignmentID) {
        this.isQuestionSelected = true
      }
    }, err => {
      console.error(err)
    })
  }

  //********************************************* CREATE  FORMS ***********************************************

  // Create talent assignment submission form.
  createSubmissionForm(): void {
    this.submissionForm = this.formBuilder.group({
      talentID: new FormControl(this.talentID),
      isAccepted: new FormControl(false),
      isChecked: new FormControl(false),
      githubURL: new FormControl(null, [Validators.maxLength(500),
      Validators.pattern(/^([A-Za-z0-9]+@|http(|s)\:\/\/)([A-Za-z0-9.]+(:\d+)?)(?::|\/)([\d\/\w.-]+?)(\.git)?$/i)]),
      assignmentSubmissionUploads: new FormArray([]),
    })
  }

  // Add new upload to submission form.
  addUploadToSubmissionForm(): void {
    this.uploadControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      assignmentSubmissionID: new FormControl(null),
      imageURL: new FormControl(null, [Validators.required]),
      description: new FormControl(null, [Validators.maxLength(1000)]),
      isImageUploadedToServer: new FormControl(false),
      imageDocStatus: new FormControl(""),
      imageDisplayedFileName: new FormControl("Select file"),
      isImageUploading: new FormControl(false),
    }))
  }

  // Get uploads array from submission form.
  get uploadControlArray(): FormArray {
    return this.submissionForm.get("assignmentSubmissionUploads") as FormArray
  }

  // Delete upload from submission form.
  deleteUpload(index: number) {
    if (confirm("Are you sure you want to delete the image?")) {
      this.submissionForm.markAsDirty()
      this.uploadControlArray.removeAt(index)
    }
  }

  //********************************************* TALENT ASSIGNMENT FUNCTIONS ************************************************************

  // On clicking assignment row in assignment table.
  onAssignmentRowClick(assignment: any): void {
    this.isAssignmentListVisible = false
    this.isQuestionSelected = true
    this.selectedBatchTopicAssignment = assignment
    console.log("selected batch topic assignment", this.selectedBatchTopicAssignment)
    this.formatProblemFields()
    this.redirectToSamePage()
  }

  // Format fields of problem details.
  formatProblemFields(): void {

    // Set class for difficulty.
    if (this.selectedBatchTopicAssignment.programmingQuestion?.level == 1) {
      this.selectedBatchTopicAssignment.programmingQuestion.levelName = "Easy"
      this.selectedBatchTopicAssignment.programmingQuestion.levelClass = "easy"
    }
    if (this.selectedBatchTopicAssignment.programmingQuestion?.level == 2) {
      this.selectedBatchTopicAssignment.programmingQuestion.levelName = "Medium"
      this.selectedBatchTopicAssignment.programmingQuestion.levelClass = "medium"
    }
    if (this.selectedBatchTopicAssignment.programmingQuestion?.level == 3) {
      this.selectedBatchTopicAssignment.programmingQuestion.levelName = "Hard"
      this.selectedBatchTopicAssignment.programmingQuestion.levelClass = "hard"
    }
  }

  // Format the fields of batch topic assignment list.
  formatBatchTopicAssignmentFields(): void {
    for (let i = 0; i < this.batchTopicAssignmentList.length; i++) {

      // Format submitted on date.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0) {
        this.batchTopicAssignmentList[i].submissions[0].submittedOn = this.datePipe.transform(this.batchTopicAssignmentList[i].submissions[0]?.submittedOn, 'EEE, MMM d, y')
      }

      // Format score.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0 && this.batchTopicAssignmentList[i].submissions[0]?.score != null) {
        this.batchTopicAssignmentList[i].score = this.batchTopicAssignmentList[i].submissions[0]?.score
      }

      if (this.batchTopicAssignmentList[i].submissions?.length == 0 || this.batchTopicAssignmentList[i].submissions[0]?.score == null) {
        this.batchTopicAssignmentList[i].score = null
      }

      // Format faculty remarks.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0 && this.batchTopicAssignmentList[i].submissions[0]?.facultyRemarks != null) {
        this.batchTopicAssignmentList[i].facultyRemarks = this.batchTopicAssignmentList[i].submissions[0]?.facultyRemarks
      }

      if (this.batchTopicAssignmentList[i].submissions?.length == 0 || this.batchTopicAssignmentList[i].submissions[0]?.facultyRemarks == null) {
        this.batchTopicAssignmentList[i].facultyRemarks = null
      }

      // Set completed status.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0 && this.batchTopicAssignmentList[i].submissions[0].isAccepted) {
        this.batchTopicAssignmentList[i].status = this.COMPLETED_STATUS
        this.batchTopicAssignmentList[i].statusOrder = this.COMPLETED_ORDER
        continue
      }

      // Set submitted status.
      if (this.batchTopicAssignmentList[i].submissions?.length > 0 && !this.batchTopicAssignmentList[i].submissions[0].isChecked) {
        this.batchTopicAssignmentList[i].status = this.SUBMITTED_STATUS
        this.batchTopicAssignmentList[i].statusOrder = this.SUBMITTED_ORDER
        continue
      }

      // Set pending status.
      this.batchTopicAssignmentList[i].status = this.PENDING_STATUS
      this.batchTopicAssignmentList[i].statusOrder = this.PENDING_ORDER
    }

    // Arrange he batch topic assignment list by their status.
    this.batchTopicAssignmentList.sort(this.compareByStatusOrder)
  }

  //********************************************* TALENT SUBMISSION FUNCTIONS ***********************************************

  // On clicking comment.
  onCommentClick(submission: any): void {
    this.selectedSubmission = submission
    this.openModal(this.commentModal, 'md')
  }

  //On Voice Click.
  onVoiceClick(submission: any): void {
    this.selectedSubmission = submission
    this.openModal(this.voiceNoteModal, 'md')
  }
  // On clicking images.
  onImagesClick(submission: any): void {
    this.selectedSubmission = submission
    this.openModal(this.imagesModal, 'lg')
  }

  // Set the current batch topic assignment after adding submission.
  getCurrentBatchTopicAssignment(selectBatchTopicAssignmentID: string): void {
    for (let i = 0; i < this.batchTopicAssignmentList.length; i++) {
      if (this.batchTopicAssignmentList[i].id == selectBatchTopicAssignmentID) {
        this.selectedBatchTopicAssignment = this.batchTopicAssignmentList[i]
      }
    }
  }

  // On clicking add new talent assignment submission.
  onAddNewSubmissionClick(): void {
    this.createSubmissionForm()
    this.openModal(this.submissionFormModal, 'md')
  }

  // Validate talent assignment submission form.
  validateSubmissionForm(): void {
    if (this.submissionForm.invalid) {
      alert("Github Link or Image must be specified")
      this.submissionForm.markAllAsTouched()
      return
    }
    var submissionData: any = this.submissionForm.value
    if ((submissionData?.githubURL == "" || submissionData?.githubURL == null) &&
      submissionData?.assignmentSubmissionUploads?.length == 0) {
      alert("Github Link or Image must be specified")
      return
    }
    this.addSubmissionForm()
  }

  // Add talent assignment submission.
  addSubmissionForm(): void {
    this.spinnerService.loadingMessage = "Adding Assignment Submission"

    this.batchTopicAssignmentService.addAssignmentSubmission(this.selectedBatchTopicAssignment.id, this.talentID,
      this.submissionForm.value).subscribe((response: any) => {
        this.modalRef.close()
        this.getAllBatchTopicAssignmentsForTalent()
        alert(response.body)
      }, (error) => {
        // console.error(error)
        if (error.error?.error) {
          alert(error.error?.error)
          return
        }
        alert(error.statusText)
      })
  }

  // On uplaoding image.
  onImageSelect(event: any, index: number): void {
    this.uploadControlArray.at(index).get('imageDocStatus').setValue("")
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      let err = this.fileOperationService.isImageFileValid(file)
      if (err != null) {
        this.uploadControlArray.at(index).get('imageDocStatus').setValue("<p><span>&#10060;</span>" + err + "</p>")
        return
      }
      // Upload image if it is present.
      this.uploadControlArray.at(index).get('isImageUploading').setValue(true)
      this.fileOperationService.uploadTalentSubmissionImage(file)
        .subscribe((data: any) => {
          this.submissionForm.markAsDirty()
          this.uploadControlArray.at(index).get('imageURL').setValue(data)
          this.uploadControlArray.at(index).get('imageDisplayedFileName').setValue(file.name)
          this.uploadControlArray.at(index).get('isImageUploading').setValue(false)
          this.uploadControlArray.at(index).get('isImageUploadedToServer').setValue(true)
          this.uploadControlArray.at(index).get('imageDocStatus').setValue("<p><span class='green'>&#10003;</span> File uploaded.</p>")
        }, (error) => {
          this.uploadControlArray.at(index).get('isImageUploading').setValue(false)
          this.uploadControlArray.at(index).get('imageDocStatus').setValue("<p><span>&#10060;</span>" + error + "</p>")
        })
    }
  }

  // Set the selected bacth topic assignment from query params.
  setSelectedBatchTopicAssignmentFromQueryParams(): void {
    this.route.queryParamMap.subscribe(params => {
      this.selectedBatchTopicAssignmentID = params.get("batchTopicAssignmentID")
      if (this.selectedBatchTopicAssignmentID) {
        this.getCurrentBatchTopicAssignment(this.selectedBatchTopicAssignmentID)
        this.isAssignmentListVisible = false
      }
    }, err => {
      console.error(err)
    })
  }

  // On clicking the score of submission.
  onSubmissionScoreClick(submission: any): void {
    this.selectedSubmission = submission
    console.log(submission)
    this.openModal(this.conceptScoreModal, 'md')
  }

  //********************************************* OTHER FUNCTIONS ************************************************************

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

  // Redirect to external link.
  redirectToExternalLink(url): void {
    window.open(url, "_blank");
  }

  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef): void {
    let isImageUploading: boolean
    for (let i = 0; i < this.uploadControlArray.length; i++) {
      if (this.uploadControlArray.at(i).get('isImageUploading').value == true) {
        isImageUploading = true
        break
      }
    }
    if (isImageUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    let isImageUploadedToServer: boolean
    for (let i = 0; i < this.uploadControlArray.length; i++) {
      if (this.uploadControlArray.at(i).get('isImageUploadedToServer').value == true) {
        isImageUploadedToServer = true
        break
      }
    }
    if (isImageUploadedToServer) {
      if (!confirm("Uploaded images will be deleted.\nAre you sure you want to close?")) {
        return
      }
      // this.deleteUploadedFile()
    }
    modal.dismiss()
    for (let i = 0; i < this.uploadControlArray.length; i++) {
      this.uploadControlArray.at(i).get('isImageUploadedToServer').setValue(false)
      this.uploadControlArray.at(i).get('imageDisplayedFileName').setValue("Select file")
      this.uploadControlArray.at(i).get('imageDocStatus').setValue("")
    }
  }

  // Chcek if image is uploading.
  checkIfImageIsUploading(): boolean {
    let isImageUploading: boolean
    for (let i = 0; i < this.uploadControlArray.length; i++) {
      if (this.uploadControlArray.at(i).get('isImageUploading').value == true) {
        isImageUploading = true
        break
      }
    }
    return isImageUploading
  }

  // Sort array.
  compareByStatusOrder(a, b): number {
    if (a.statusOrder < b.statusOrder) {
      return -1
    }
    if (a.statusOrder > b.statusOrder) {
      return 1
    }
    return 0
  }

  // Redirect to Same page.
  redirectToSamePage(): void {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        "batchID": this.batchID,
        "tab": "Assignments",
        "subTab": "Submit",
        "batchTopicAssignmentID": this.selectedBatchTopicAssignment.id,
      },
    })
  }

  //********************************************* GET FUNCTIONS ************************************************************

  // Get all batch topic assignments for talent.
  getAllBatchTopicAssignmentsForTalent(): void {
    this.batchTopicAssignmentList = []
    this.spinnerService.loadingMessage = "Getting All Assignments..."

    this.batchTopicAssignmentService.getAllBatchTopicAssignmentsForTalent(this.batchID, this.talentID).subscribe((response) => {
      this.batchTopicAssignmentList = response.body
      console.log("batchTopic Assignment List", this.batchTopicAssignmentList)
      if (!this.isAssignmentListVisible) {
        this.getCurrentBatchTopicAssignment(this.selectedBatchTopicAssignment.id)
      }
      this.formatBatchTopicAssignmentFields()
      this.setSelectedBatchTopicAssignmentFromQueryParams()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }

  voiceNote(url: string) {
    return this.domSanitizer.bypassSecurityTrustUrl(url);
  }

}

