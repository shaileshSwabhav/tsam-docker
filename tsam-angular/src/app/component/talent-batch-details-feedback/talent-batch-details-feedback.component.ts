import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchSessionService } from 'src/app/service/batch-session/batch-session.service';
import { BatchService } from 'src/app/service/batch/batch.service';
import { FeedbackService } from 'src/app/service/feedback/feedback.service';
import { IFeedbackQuestion } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { TalentDashboardService } from 'src/app/service/talent-dashboard/talent-dashboard.service';

@Component({
  selector: 'app-talent-batch-details-feedback',
  templateUrl: './talent-batch-details-feedback.component.html',
  styleUrls: ['./talent-batch-details-feedback.component.css']
})
export class TalentBatchDetailsFeedbackComponent implements OnInit {

  readonly MAX_SCORE = 10

  // Batch session.
  batchID: string
  batchSessionList: any[]
  selectedBatchSessionForTalentFeedback: any
  batchRelatedDetailsCount: number

  // Feedback.
  selectedBatchSessionFaculty: any
  selectedBatch
  talentFeedbackForFacultyList: any[]
  @ViewChild('talentFeedbackToFacultyModal') talentFeedbackToFacultyModal: any
  talentFeedbackToFaculytyForm: FormGroup
  talentTofacultyFeedbackQuestionList: IFeedbackQuestion[]
  submittedTalentSingleSessionFeedbackList: any[]
  @Input() isSubmission: boolean

  // Leaderboard.
  facultyFeedbackForTalentLeaderBoard: any[]

  // Talent.
  talentID: string

  // Modal.
  modalRef: any

  // Batch session talent.
  batchSessionTalentList: any[]
  currentBatchSession: any

  constructor(
    private route: ActivatedRoute,
    private spinnerService: SpinnerService,
    private talentDashboardService: TalentDashboardService,
    private localService: LocalService,
    private modalService: NgbModal,
    private formBuilder: FormBuilder,
    private feedbackService: FeedbackService,
    private batchService: BatchService,
    private batchSessionService: BatchSessionService,
  ) {
    this.initializeVariables()
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {

    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")

      // If is submission is true then show feedback list.
      if (this.batchID && this.isSubmission) {
        this.getAllComponents()
      }

      // If is submission is false then show leaderboard.
      if (this.batchID && !this.isSubmission) {
        this.getFacultyFeedbackForTalentLeaderBoard()
      }
    }, err => {
      console.error(err)
    })
  }

  // Initialize all global variables.
  initializeVariables() {

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Feedback List..."

    // Talent.
    this.talentID = this.localService.getJsonValue("loginID")

    // Batch session.
    this.batchSessionList = []
    this.batchRelatedDetailsCount = 3

    // Feedback.
    this.talentFeedbackForFacultyList = []
    this.talentTofacultyFeedbackQuestionList = []
    this.submittedTalentSingleSessionFeedbackList = []

    // Batch session talent.
    this.batchSessionTalentList = []

    // Leaderboard.
    this.facultyFeedbackForTalentLeaderBoard = []
  }

  //********************************************* TALENT FEEDACK FOR FACULTY FUNCTIONS ************************************************************

  // Format batch session list.
  formatBatchSessionDateList(): void {

    // Give attendance of talent to all batch session list.
    for (let i = 0; i < this.batchSessionList.length; i++) {
      this.batchSessionList[i].sessionNumber = (i + 1)
      this.batchSessionList[i].isFeedbackGiven = false

      // Set the topic and sub topics for the batch session list.
      let topics: any[] = []
      for (let j = 0; j < this.batchSessionList[i].batchSessionTopics.length; j++) {
        if (j == 0) {
          topics.push(this.batchSessionList[i].batchSessionTopics[j].topic)
          topics[0].subTopics = []
          topics[0].subTopics.push(this.batchSessionList[i].batchSessionTopics[j].subTopic)
          continue
        }
        if (j > 0 && this.batchSessionList[i].batchSessionTopics[j].topic.id == topics[topics.length - 1].id) {
          topics[topics.length - 1].subTopics.push(this.batchSessionList[i].batchSessionTopics[j].subTopic)
        }
        if (j > 0 && this.batchSessionList[i].batchSessionTopics[j].topic.id != topics[topics.length - 1].id) {
          topics.push(this.batchSessionList[i].batchSessionTopics[j].topic)
          topics[topics.length - 1].subTopics = []
          topics[topics.length - 1].subTopics.push(this.batchSessionList[i].batchSessionTopics[j].subTopic)
        }
      }
      this.batchSessionList[i].topics = topics

      // Iterate batch session list.
      for (let j = 0; j < this.batchSessionTalentList.length; j++) {
        if (this.batchSessionList[i].id == this.batchSessionTalentList[j].batchSessionID
          && this.batchSessionTalentList[j].isPresent == true) {
          this.batchSessionList[i].isPresent = true
        }
        if (this.batchSessionList[i].id == this.batchSessionTalentList[j].batchSessionID
          && this.batchSessionTalentList[j].isPresent == false) {
          this.batchSessionList[i].isPresent = false
        }
      }
    }

    // Iterate batch session list.
    for (let i = 0; i < this.batchSessionList.length; i++) {

      // Iterate talent feedback for faculty list.
      for (let j = 0; j < this.talentFeedbackForFacultyList.length; j++) {

        // If batch session is complete, attendance is present and feedback is not given then show 
        // option for feedback.
        if (this.batchSessionList[i].id == this.talentFeedbackForFacultyList[j].batchSessionID &&
          this.batchSessionList[i].isSessionTaken == true && this.batchSessionList[i].isPresent == true) {
          this.batchSessionList[i].feedbacks = this.talentFeedbackForFacultyList[j].sessionFeedbacks
          if (this.talentFeedbackForFacultyList[j].sessionFeedbacks?.length > 0) {
            this.batchSessionList[i].isFeedbackGiven = true
          }
          if (this.talentFeedbackForFacultyList[j].sessionFeedbacks?.length == 0) {
            this.batchSessionList[i].isFeedbackGiven = false
          }
        }
      }

      // Keep pending batch sessions at the top.
      let pendingBatchSessionList: any = []
      let submittedBatchSessionList: any = []
      for (let i = 0; i < this.batchSessionList.length; i++) {
        if (this.batchSessionList[i].isFeedbackGiven == true && this.batchSessionList[i].isSessionTaken == true
          && this.batchSessionList[i].isPresent == true) {
          submittedBatchSessionList.push(this.batchSessionList[i])
        }
        if (this.batchSessionList[i].isFeedbackGiven == false && this.batchSessionList[i].isSessionTaken == true
          && this.batchSessionList[i].isPresent == true) {
          pendingBatchSessionList.push(this.batchSessionList[i])
        }
      }
      this.batchSessionList = []
      for (let i = 0; i < pendingBatchSessionList.length; i++) {
        this.batchSessionList.push(pendingBatchSessionList[i])
      }
      for (let i = 0; i < submittedBatchSessionList.length; i++) {
        this.batchSessionList.push(submittedBatchSessionList[i])
      }
    }
  }

  //********************************************* TALENT FEEDACK FOR FACULTY FUNCTIONS ************************************************************

  // On clicking talent feedback to faculty.
  onTalentFeeddbackToFacultyClick(session: any): void {
    this.currentBatchSession = session
    for (let i = 0; i < this.batchSessionList.length; i++) {
      if (this.batchSessionList[i].date == session.date) {
        this.selectedBatchSessionFaculty = this.batchSessionList[i].faculty
      }
    }
    this.openModal(this.talentFeedbackToFacultyModal, 'lg')
  }

  // Receive is feedback add successful from child.
  receiveIsFeedbackAddSuccessful(isSuccessful: any): void {

    // If modal is closed.
    if (isSuccessful == null) {
      this.modalRef.close()
    }

    // If add is successful.
    if (isSuccessful == true) {
      this.modalRef.close()
      this.batchRelatedDetailsCount = 3
      this.getAllComponents()
    }
  }

  //********************************************* LEADERBOARD FUNCTIONS************************************************************

  // Format the faculty feedback for talent ledaer board fields.
  formatFacultyFeedbackForTalentLeaderBoard(): void {
    let rank: number = 1
    this.facultyFeedbackForTalentLeaderBoard.sort(this.sortLeaderBoardByRating)
    for (let i = 0; i < this.facultyFeedbackForTalentLeaderBoard.length; i++) {

      // Give rank.
      if (i == 0){
        this.facultyFeedbackForTalentLeaderBoard[i].rank = rank
      }
      if (i > 0){
        if (this.facultyFeedbackForTalentLeaderBoard[i-1].rating != this.facultyFeedbackForTalentLeaderBoard[i].rating){
          rank = rank + 1
        }
        this.facultyFeedbackForTalentLeaderBoard[i].rank = rank
      }

      // Set the images.
      if (i > 2){
        continue
      }
      this.facultyFeedbackForTalentLeaderBoard[i].imageURL = "assets/icon/talent-dashboard/badge.png"
      if (i == 0) {
        this.facultyFeedbackForTalentLeaderBoard[i].imageURL = "assets/icon/talent-dashboard/first.png"
      }
      if (i == 1) {
        this.facultyFeedbackForTalentLeaderBoard[i].imageURL = "assets/icon/talent-dashboard/second.png"
      }
      if (i == 2) {
        this.facultyFeedbackForTalentLeaderBoard[i].imageURL = "assets/icon/talent-dashboard/third.png"
      }
    }
  }

  // Sort leader borad by rating.
  sortLeaderBoardByRating(a, b) {
    if (a.rating > b.rating) {
      return -1
    }
    if (a.rating < b.rating) {
      return 1
    }
    return 0
  }

  //********************************************* OTHER FUNCTIONS************************************************************

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

  // Decrement batch related details count.
  decrementBatchRelatedDetailsCount(): void {
    this.batchRelatedDetailsCount = this.batchRelatedDetailsCount - 1
    if (this.batchRelatedDetailsCount == 0) {
      this.formatBatchSessionDateList()
    }
  }

  //*********************************************GET FUNCTIONS************************************************************

  // Get all components.
  getAllComponents(): void {
    this.getBatchSessionList()
    this.getTalentFeedbackForFacultyList()
    this.getBatchSessionTalentList()
  }

  // Get all session dates for batch.
  getBatchSessionList(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.batchSessionService.getBatchSessionWithTopicNameList(this.batchID).subscribe((response: any) => {
      this.batchSessionList = response.body
    }, (err: any) => {
      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.decrementBatchRelatedDetailsCount()
    })
  }

  // Get talent feedback for faculty list.
  getTalentFeedbackForFacultyList(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.batchService.getTalentBatchFeedback(this.batchID, this.talentID).subscribe((response) => {
      this.talentFeedbackForFacultyList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.decrementBatchRelatedDetailsCount()
    })
  }

  // Get batch session talent list.
  getBatchSessionTalentList(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    this.batchService.getAllBatchSessionTalentsForTalent(this.batchID, this.talentID).subscribe((response) => {
      this.batchSessionTalentList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    }).add(() => {
      this.decrementBatchRelatedDetailsCount()
    })
  }

  // Get submitted talent feedback for single batch session.
  getSubmittedTalentFeedbackList(session: any): void {
    this.batchService.getSpecifiedTalentBatchSessionFeedback(this.batchID, session.id, this.talentID).subscribe((response: any) => {
      this.submittedTalentSingleSessionFeedbackList = response.body
      this.onTalentFeeddbackToFacultyClick(session)
    }, (err: any) => {
      console.error(err)
    })
  }

  // Get faculty feedback for talent leader board.
  getFacultyFeedbackForTalentLeaderBoard(): void {
    this.talentDashboardService.getFacultyFeedbackForTalentLeaderBoard(this.batchID).subscribe((response: any) => {
      this.facultyFeedbackForTalentLeaderBoard = response
      this.formatFacultyFeedbackForTalentLeaderBoard()
    }, (err: any) => {
      console.error(err)
    })
  }

}
