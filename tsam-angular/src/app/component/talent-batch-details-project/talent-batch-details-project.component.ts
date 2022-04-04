import { DatePipe } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchProjectService } from 'src/app/service/batch-project/batch-project.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { DomSanitizer } from '@angular/platform-browser';

@Component({
  selector: 'app-talent-batch-details-project',
  templateUrl: './talent-batch-details-project.component.html',
  styleUrls: ['./talent-batch-details-project.component.css']
})
export class TalentBatchDetailsProjectComponent implements OnInit {

  public readonly PENDING_STATUS = "Pending"
  public readonly PENDING_ORDER = 1
  public readonly SUBMITTED_STATUS = "Submitted"
  public readonly SUBMITTED_ORDER = 2
  public readonly COMPLETED_STATUS = "Completed"
  public readonly COMPLETED_ORDER = 3

  // Batch Topic Assignment.
  batchProjectList: any[]
  selectedBatchProject: any
  selectedBatchProjectID: string

  // Talent assignment submission.
  submissionForm: FormGroup
  selectedSubmission: any

  // Flags.
  isProjectListVisible: boolean
  isQuestionSelected: boolean

  // Spinner.



  // Talent.
  talentID: string

  //project Upload
  isProjectUploadedToServer: boolean
  isProjectUploading: boolean
  projectStatus: string
  displayedFileName: string

  // Modal.
  modalRef: any
  @ViewChild('commentModal') commentModal: any
  @ViewChild('imagesModal') imagesModal: any
  @ViewChild('voiceNoteModal') voiceNoteModal: any
  @ViewChild('submissionFormModal') submissionFormModal: any
  @ViewChild('conceptScoreModal') conceptScoreModal: any

  // Batch.
  batchID: string

  constructor(
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private batchProjectService: BatchProjectService,
    private route: ActivatedRoute,
    private localService: LocalService,
    private formBuilder: FormBuilder,
    private datePipe: DatePipe,
    private fileOperationService: FileOperationService,
    private router: Router,
    private domSanitizer: DomSanitizer,
  ) {
    this.initializeVariables()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  // Initialize all global variables.
  initializeVariables() {

    // Flags.
    this.isProjectListVisible = true
    this.isQuestionSelected = false

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Problem Details"

    this.batchProjectList = []

    // Talent.
    this.talentID = this.localService.getJsonValue("loginID")
    console.log("talent id", this.talentID)

    // Batch.
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      console.log("batchId", this.batchID)
      this.selectedBatchProjectID = params.get("batchProjectID")
      console.log("selectedBatchProjectID", this.selectedBatchProjectID)

      if (this.batchID) {
        this.getAllBatchProjectForTalent()
      }
      if (this.selectedBatchProjectID) {
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
      githubURL: new FormControl(null, [Validators.required, Validators.maxLength(500),
      Validators.pattern(/^([A-Za-z0-9]+@|http(|s)\:\/\/)([A-Za-z0-9.]+(:\d+)?)(?::|\/)([\d\/\w.-]+?)(\.git)?$/i)]),
      websiteLink: new FormControl(null),
      projectUpload: new FormControl(null),
      projectSubmissionUploads: new FormArray([]),
    })
  }

  // Add new upload to submission form.
  addUploadToSubmissionForm(): void {
    this.uploadControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      assignmentSubmissionID: new FormControl(null),
      imageURL: new FormControl(null),
      description: new FormControl(null, [Validators.maxLength(1000)]),
      isImageUploadedToServer: new FormControl(false),
      imageDocStatus: new FormControl(""),
      imageDisplayedFileName: new FormControl("Select file"),
      isImageUploading: new FormControl(false),
    }))
  }

  // Get uploads array from submission form.
  get uploadControlArray(): FormArray {
    return this.submissionForm.get("projectSubmissionUploads") as FormArray
  }

  // Delete upload from submission form.
  deleteUpload(index: number) {
    if (confirm("Are you sure you want to delete the image?")) {
      this.submissionForm.markAsDirty()
      this.uploadControlArray.removeAt(index)
    }
  }

  //********************************************* TALENT ASSIGNMENT FUNCTIONS ************************************************************

  // On clicking project row in assignment table.
  onProjectRowClick(project: any): void {
    this.isProjectListVisible = false
    this.isQuestionSelected = true
    this.selectedBatchProject = project
    console.log("selectedBatchProject", this.selectedBatchProject)
    this.formatProblemFields()
    this.redirectToSamePage()
  }

  // Format fields of problem details.
  formatProblemFields(): void {

    // Set class for difficulty.
    if (this.selectedBatchProject.programmingProject?.complexityLevel >=1 && 
      this.selectedBatchProject.programmingProject?.complexityLevel <=5) {
      this.selectedBatchProject.programmingProject.levelName = "Easy"
      this.selectedBatchProject.programmingProject.levelClass = "easy"
    }
    if (this.selectedBatchProject.programmingProject?.complexityLevel >=6 && 
      this.selectedBatchProject.programmingProject?.complexityLevel <=8) {
      this.selectedBatchProject.programmingProject.levelName = "Medium"
      this.selectedBatchProject.programmingProject.levelClass = "medium"
    }
    if (this.selectedBatchProject.programmingProject?.complexityLevel >=9 &&
       this.selectedBatchProject.programmingProject?.complexityLevel <=10) {
      this.selectedBatchProject.programmingProject.levelName = "Hard"
      this.selectedBatchProject.programmingProject.levelClass = "hard"
    }
  }

  // Format the fields of batch topic project list.
  formatBatchProjectFields(): void {
    for (let i = 0; i < this.batchProjectList.length; i++) {

      // Format submitted on date.
      if (this.batchProjectList[i].submissions?.length > 0) {
        this.batchProjectList[i].submissions[0].submittedOn = this.datePipe.transform(this.batchProjectList[i].submissions[0]?.submittedOn, 'EEE, MMM d, y')
      }

      // Format score.
      if (this.batchProjectList[i].submissions?.length > 0 && this.batchProjectList[i].submissions[0]?.score != null) {
        this.batchProjectList[i].score = this.batchProjectList[i].submissions[0]?.score
      }

      // if (this.batchProjectList[i].submissions?.length == 0 || this.batchProjectList[i].submissions[0]?.score == null) {
      //   this.batchProjectList[i].score = null
      // }

      // Format faculty remarks.
      if (this.batchProjectList[i].submissions?.length > 0 && this.batchProjectList[i].submissions[0]?.facultyRemarks != null) {
        this.batchProjectList[i].facultyRemarks = this.batchProjectList[i].submissions[0]?.facultyRemarks
      }

      // if (this.batchProjectList[i].submissions?.length == 0 || this.batchProjectList[i].submissions[0]?.facultyRemarks == null) {
      //   this.batchProjectList[i].facultyRemarks = null
      // }

      // Set completed status.
      if (this.batchProjectList[i].submissions?.length > 0 && this.batchProjectList[i].submissions[0].isAccepted) {
        this.batchProjectList[i].status = this.COMPLETED_STATUS
        this.batchProjectList[i].statusOrder = this.COMPLETED_ORDER
        continue
      }

      // Set submitted status.
      if (this.batchProjectList[i].submissions?.length > 0 && !this.batchProjectList[i].submissions[0].isChecked) {
        this.batchProjectList[i].status = this.SUBMITTED_STATUS
        this.batchProjectList[i].statusOrder = this.SUBMITTED_ORDER
        continue
      }

      // Set pending status.
      this.batchProjectList[i].status = this.PENDING_STATUS
      this.batchProjectList[i].statusOrder = this.PENDING_ORDER
    }

    // Arrange he batch topic assignment list by their status.
    this.batchProjectList.sort(this.compareByStatusOrder)
  }

  //********************************************* TALENT SUBMISSION FUNCTIONS ***********************************************

  // On clicking comment.
  onCommentClick(submission: any): void {
    this.selectedSubmission = submission
    this.openModal(this.commentModal, 'md')
  }

  // On clicking images.
  onImagesClick(submission: any): void {
    this.selectedSubmission = submission
    this.openModal(this.imagesModal, 'md')
  }

  // Set the current batch topic assignment after adding submission.
  getCurrentBatchProject(selectBatchTopicAssignmentID: string): void {
    for (let i = 0; i < this.batchProjectList.length; i++) {
      if (this.batchProjectList[i].id == selectBatchTopicAssignmentID) {
        this.selectedBatchProject = this.batchProjectList[i]
      }
    }
    console.log("selectedBatchProject", this.selectedBatchProject)
  }

  // On clicking add new talent assignment submission.
  onAddNewSubmissionClick(): void {
    this.createSubmissionForm()
    this.openModal(this.submissionFormModal, 'md')
  }

  // Validate talent assignment submission form.
  validateSubmissionForm(): void {
    if (this.submissionForm.invalid) {
      this.submissionForm.markAllAsTouched()
      return
    }
    console.log("submissionForm", this.submissionForm)
    this.addSubmissionForm()
  }

  // Add talent assignment submission.
  addSubmissionForm(): void {
    this.spinnerService.loadingMessage = "Adding Project Submission"

    console.log("batchID", this.batchID)
    console.log("talentID", this.talentID)
    console.log("selectedBatchProjectID", this.selectedBatchProjectID)


    this.batchProjectService.addProjectSubmission(this.batchID, this.talentID, this.selectedBatchProjectID,
      this.submissionForm.value).subscribe((response: any) => {
        this.modalRef.close()
        this.getAllBatchProjectForTalent()
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
      this.fileOperationService.uploadTalentProjectSubmissionImage(file)
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
  setSelectedBatchProjectFromQueryParams(): void {
    this.route.queryParamMap.subscribe(params => {
      this.selectedBatchProjectID = params.get("batchProjectID")
      if (this.selectedBatchProjectID) {
        this.getCurrentBatchProject(this.selectedBatchProjectID)
        this.isProjectListVisible = false
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
    if (this.isProjectUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    if (this.isProjectUploadedToServer) {
      if (!confirm("Uploaded project will be deleted.\nAre you sure you want to close?")) {
        return
      }
    }
    modal.dismiss()
    for (let i = 0; i < this.uploadControlArray.length; i++) {
      this.uploadControlArray.at(i).get('isImageUploadedToServer').setValue(false)
      this.uploadControlArray.at(i).get('imageDisplayedFileName').setValue("Select file")
      this.uploadControlArray.at(i).get('imageDocStatus').setValue("")
    }
    this.resetProjectUploadFields()
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
        "tab": "Projects",
        "subTab": "Submit",
        "batchProjectID": this.selectedBatchProject.id,
      },
    })
  }

  //********************************************* GET FUNCTIONS ************************************************************

  // Get all batch topic projects for talent.
  getAllBatchProjectForTalent(): void {
    this.batchProjectList = []
    this.spinnerService.loadingMessage = "Getting Projects..."

    this.batchProjectService.getAllBatchProjectForTalent(this.batchID, this.talentID, this.selectedBatchProjectID).subscribe((response) => {
      this.batchProjectList = response.body
      console.log("projectList", this.batchProjectList)
      if (!this.isProjectListVisible) {
        this.getCurrentBatchProject(this.selectedBatchProject.id)
      }
      this.formatBatchProjectFields()
      this.setSelectedBatchProjectFromQueryParams()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }


  //On uplaoding Document
  onResourceSelect(event: any) {
    this.projectStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      let err = this.fileOperationService.isProjectFileValid(file)
      if (err != null) {
        this.projectStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // console.log(file)
      // Upload project if it is present.]
      this.isProjectUploading = true
      this.fileOperationService.uploadProject(file)
        .subscribe((data: any) => {
          this.submissionForm.markAsDirty()
          this.submissionForm.patchValue({
            projectUpload: data
          })
          if (file.name.toString().length > 25) {
            this.displayedFileName = file.name.toString().substr(0, 25) + "....."
          } else {
            this.displayedFileName = file.name
          }
          this.isProjectUploading = false
          this.isProjectUploadedToServer = true
          this.projectStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
        }, (error) => {
          this.isProjectUploading = false
          this.projectStatus = `<p><span>&#10060;</span> ${error}</p>`
        })
    }
  }

  resetProjectUploadFields(): void {
    this.isProjectUploadedToServer = false
    this.isProjectUploading = false
    this.displayedFileName = "Select file"
    this.projectStatus = ""
  }

  //On Voice Click.
  onVoiceClick(submission: any): void {
    this.selectedSubmission = submission
    console.log(this.selectedSubmission)
    this.openModal(this.voiceNoteModal, 'md')
  }

  voiceNote(url: string) {
    console.log(url)
    return this.domSanitizer.bypassSecurityTrustUrl(url);
  }
}

