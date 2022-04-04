import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { BatchService, IBatchTopicAssignment } from 'src/app/service/batch/batch.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { Location } from '@angular/common';

@Component({
  selector: 'app-my-sessions',
  templateUrl: './my-sessions.component.html',
  styleUrls: ['./my-sessions.component.css']
})
export class MySessionsComponent implements OnInit {

  // Session.
  sessionList: any[]
  selectedSession: any
  selectedSubSessionList: any[]
  selectedSubSession: any

  // assignment
  selectedAssignment: IBatchTopicAssignment

  // title
  title: string

  // Spinner.



  // Course.
  batchName: string
  batchID: string

  constructor(
    private activatedRoute: ActivatedRoute,
    private spinnerService: SpinnerService,
    private batchService: BatchService,
    private localService: LocalService,
    private _location: Location,
    private router: Router,
  ) {
    this.initializeVariables()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  // Initialize global variables.
  initializeVariables(): void {

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Sessions"


    // Session.
    this.sessionList = []
    this.selectedSession = null
    this.selectedSubSessionList = []
    this.selectedSubSession = null

    // Get query params.
    this.batchName = this.activatedRoute.snapshot.queryParamMap.get("batchName")
    this.batchID = this.activatedRoute.snapshot.queryParamMap.get("batchID")
    if (this.batchID) {
      this.getBatchSessionAndAssignmentList()
    }
  }

  //*********************************************GET FUNCTIONS************************************************************

  // Get batch sessions by batch id.
  getBatchSessionAndAssignmentList(): void {
    let queryParams: any = {
      roleName: this.localService.getJsonValue("roleName"),
      loginID: this.localService.getJsonValue("loginID"),
    }
    this.spinnerService.loadingMessage = "Getting Sessions"


    this.batchService.getSessionAndAssignmentForBatch(this.batchID, queryParams).subscribe((response) => {
      this.sessionList = response.body
      console.log(this.sessionList)
      if (this.sessionList.length > 0) {
        this.formatSessionsList()
      }
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  //*********************************************OTHER FUNCTIONS************************************************************

  // Format the session list on getting session list.
  formatSessionsList(): void {
    // Initial class for all the session buttons.
    for (let i = 0; i < this.sessionList.length; i++) {
      if (i == 0) {
        this.sessionList[i].class = "card selected-session-header h-100"
        continue
      }
      this.sessionList[i].class = "card session-header h-100"
    }

    // Give the first session selected button class.
    // this.sessionList[0].class = "card selected-session-header h-100"

    // Set the first session as selected session.
    this.selectedSession = this.sessionList[0]
    if (this.selectedSession.session.subSessions && this.selectedSession.session.subSessions.length > 0) {
      this.selectedSubSessionList = this.selectedSession.session.subSessions
      this.selectedSubSession = this.selectedSubSessionList[0]
      for (let i = 0; i < this.selectedSubSessionList.length; i++) {
        this.selectedSubSessionList[i].class = "btn btn-default subSession-header"
      }
      this.selectedSubSessionList[0].class = "btn btn-default selected-subSession-header"
    }
  }

  // On clicking session button.
  onSessionClick(sessionID: string): void {
    this.selectedSession = null
    this.selectedSubSessionList = []
    this.selectedSubSession = null

    // Set the selectedSession.
    for (let i = 0; i < this.sessionList.length; i++) {
      if (this.sessionList[i].id == sessionID) {
        this.selectedSession = this.sessionList[i]
        break
      }
    }

    // Give the selectedSession button highlight.
    this.selectedSession.class = "card selected-session-header h-100"

    // Give all other session buttons default.
    for (let i = 0; i < this.sessionList.length; i++) {
      if (this.sessionList[i].id == sessionID) {
        continue
      }
      this.sessionList[i].class = "card session-header h-100"
    }

    // Set selected subSession list and give its first entry highlight.
    if (this.selectedSession.session.subSessions && this.selectedSession.session.subSessions.length > 0) {
      this.selectedSubSessionList = this.selectedSession.session.subSessions
      this.selectedSubSession = this.selectedSubSessionList[0]
      for (let i = 0; i < this.selectedSubSessionList.length; i++) {
        this.selectedSubSessionList[i].class = "btn btn-default subSession-header"
      }
      this.selectedSubSessionList[0].class = "btn btn-default selected-subSession-header"
    }
  }

  onAssignmentClick(assignment: IBatchTopicAssignment): void {
    this.title = assignment.programmingQuestion.title
    this.selectedAssignment = assignment
    console.log(assignment);
  }

  // On clicking sub session button.
  onSubSessionClick(subSessionID: string): void {

    this.selectedSubSession = null

    // Set the selectedSubSession.
    this.selectedSubSessionList.find((subSession: any) => {
      if (subSession.id == subSessionID) {
        this.selectedSubSession = subSession
        this.title = subSession.name
        return
      }
    })
    console.log("exited");

    // for (let i = 0; i < this.selectedSubSessionList.length; i++) {
    //   if (this.selectedSubSessionList[i].id == subSessionID) {
    //     this.selectedSubSession = this.selectedSubSessionList[i]
    //     break
    //   }
    // }

    // Give the selectedSubSession button highlight.
    this.selectedSubSession.class = "btn btn-default selected-subSession-header"

    // Give all other sub session buttons default.
    for (let i = 0; i < this.selectedSubSessionList.length; i++) {
      if (this.selectedSubSessionList[i].id == subSessionID) {
        continue
      }
      this.selectedSubSessionList[i].class = "btn btn-default subSession-header"
    }
  }

  // Redirect to session feedback page. 
  redirectToSessionFeedBack(batchID: string, sessionID: string, sessionName: string): void {
    this.router.navigate(['/batch/master/session/feedback'], {
      queryParams: {
        "batchID": batchID,
        "sessionID": sessionID,
        "talentID": this.localService.getJsonValue("loginID"),
        "sessionName": sessionName
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Nagivate to previous url.
  backToPreviousPage(): void {
    this._location.back()
  }
}