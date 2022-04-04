import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, FormControl } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { IAssignmentSubmission } from 'src/app/models/assignment-submission';
import { IBatchTopicAssignment } from 'src/app/models/batch-topic-assignment';
import { AssignmentData } from 'src/app/providers/assignment-data';
import { BatchTalentService } from 'src/app/service/batch-talent/batch-talent.service';
import { BatchTopicAssignmentService } from 'src/app/service/batch-topic-assignment/batch-topic-assignment.service';
import { BackNavigationUrl, Role, UrlConstant } from 'src/app/service/constant';
import { IPermission } from 'src/app/service/menu/menu.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ITalent, ITalentAssignmentScore } from 'src/app/service/talent/talent.service';
import { UrlService } from 'src/app/service/url.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-batch-topic-assignment-score',
  templateUrl: './batch-topic-assignment-score.component.html',
  styleUrls: ['./batch-topic-assignment-score.component.css']
})
export class BatchTopicAssignmentScoreComponent implements OnInit {

  // search form
  searchScoreForm: FormGroup

  // batch-session-assignment-score
  talentAssignmentScores: ITalentAssignmentScore[]

  batchTopicAssignments: IBatchTopicAssignment[]

  batchTalents: ITalent[]

  // batch-session programming assingment
  sessionAssignmentList: IBatchTopicAssignment[]

  // access
  permission: IPermission

  // batch-details
  batchID: string
  batchName: string

  // other


  allAssignmentSubmissions: IAssignmentSubmission[]
  limit: number = 5
  currentPage: number
  offset: number = 0
  options: string[] = ["Not Scored", "Completed", "Pending", "All"]
  option: string = this.options[3]
  visibleAssignmentSubmissions: IAssignmentSubmission[]

  constructor(
    private formBuilder: FormBuilder,
    private urlConstant: UrlConstant,
    private btaService: BatchTopicAssignmentService,
    private batchTalentService: BatchTalentService,
    private utilService: UtilityService,
    private router: Router,
    private activatedRoute: ActivatedRoute,
    private assignmentProvider: AssignmentData,
    private urlService: UrlService,
    private localService: LocalService,
    private backNavigation: BackNavigationUrl,
    private spinnerService: SpinnerService,
		private role: Role,
  ) {
    this.extractData()
    this.createSearchForm()
    this.initializeVariables()
  }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
    this.getAllComponents()
  }

  initializeVariables(): void {
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN || this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_BATCH_MASTER_SESSION_DETAILS)
		}
		if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
			this.permission = this.utilService.getPermission(this.urlConstant.MY_BATCH_SESSION_DETAILS)
		}

    this.talentAssignmentScores = []
    this.batchTalents = []
    this.batchTopicAssignments = []
  }

  extractData(): void {
    this.activatedRoute.queryParamMap.subscribe(
      (params: any) => {
        this.batchID = params.get("batchID")
        this.batchName = params.get("batchName")
      }, (err: any) => {
        console.error(err);
      })
  }

  createSearchForm(): void {
    this.searchScoreForm = this.formBuilder.group({
      dueDate: new FormControl(null),
      // Check with role name #niranjan
      facultyID: new FormControl(this.localService.getJsonValue("loginID")),
    })
  }

  getAllComponents(): void {
    this.getAllAssignmentsWithScore()
    this.getBatchTalents()
  }

  // getTalentAssignmentScore(): void {
  //   this.spinnerService.loadingMessage = "Getting talent scores"
  //   
  //   this.talentAssignmentScores = []

  //   let queryParam: any = this.searchScoreForm.value

  //   if (!queryParam.dueDate) {
  //     delete queryParam.dueDate
  //   }

  //   this.talentService.getTalentAssignmentScores(this.batchID, this.searchScoreForm.value).subscribe((response: any) => {
  //     this.talentAssignmentScores = response.body
  //   }, (err: any) => {
  //     console.error(err)
  //     if (err.statusText.includes('Unknown')) {
  //       alert("No connection to server. Check internet.")
  //       return
  //     }
  //     alert(err.error.error)
  //   }).add(() => {
  //     
  //   })
  // }

  getAllAssignmentsWithScore(): void {
    this.spinnerService.loadingMessage = "Getting talent scores"
    this.batchTopicAssignments = []

    let queryParam: any = this.searchScoreForm.value

    if (!queryParam.dueDate) {
      delete queryParam.dueDate
    }
    this.btaService.getTopicAssignmentWithSubmissions(this.batchID, this.searchScoreForm.value).subscribe((response: any) => {
      this.batchTopicAssignments = response.body
      this.orderScores()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getBatchTalents(): void {
    this.spinnerService.loadingMessage = "Getting batch talents"
    this.batchTalents = []
    this.batchTalentService.getBatchTalentList(this.batchID).subscribe((response: any) => {
      this.batchTalents = response.body
      this.orderScores()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  // Need to change to wait #niranjan
  orderScores(): void {
    this.allAssignmentSubmissions = []
    let no = 0
    for (let index = 0; index < this.batchTopicAssignments.length; index++) {
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
    this.changeOption()
  }

  changeOption() {
    switch (this.option) {
      case "Not Scored":
        this.getNotScored()
        break;
      case "Completed":
        this.getCompleted()
        break;
      case "Pending":
        this.getPending()
        break;
      default:
        this.getAll()
    }
  }

  getNotScored() {
    this.visibleAssignmentSubmissions = []
    for (let i = 0; i < this.allAssignmentSubmissions.length; i++) {
      for (let value of this.allAssignmentSubmissions[i].talentSubmission.values()) {
        if (value && !value.isChecked) {
          this.visibleAssignmentSubmissions.push(this.allAssignmentSubmissions[i])
          break
        }
      }
    }
  }

  getCompleted() {
    this.visibleAssignmentSubmissions = []
    for (let i = 0; i < this.allAssignmentSubmissions.length; i++) {
      let num = 0
      for (let value of this.allAssignmentSubmissions[i].talentSubmission.values()) {
        if (value && value.isAccepted) {
          num++
          if (num == this.batchTalents.length)
            this.visibleAssignmentSubmissions.push(this.allAssignmentSubmissions[i])
          continue
        }
        break
      }
    }
  }

  getPending() {
    this.visibleAssignmentSubmissions = []
    for (let i = 0; i < this.allAssignmentSubmissions.length; i++) {
      for (let value of this.allAssignmentSubmissions[i].talentSubmission.values()) {
        if (!value || (value.isChecked && !value.isAccepted)) {
          this.visibleAssignmentSubmissions.push(this.allAssignmentSubmissions[i])
          break
        }
      }
    }
  }

  getDueDateCrossed(e: any) {
    if (!e.target.checked) {
      this.changeOption()
      return
    }
    this.visibleAssignmentSubmissions = []
    let today = new Date()
    for (let i = 0; i < this.allAssignmentSubmissions.length; i++) {
      let dueDate = new Date(this.allAssignmentSubmissions[i].assignment.dueDate)
      for (let value of this.allAssignmentSubmissions[i].talentSubmission.values()) {
        if ((!value || (value.isChecked && !value.isAccepted)) &&
          dueDate <= today) {
          this.visibleAssignmentSubmissions.push(this.allAssignmentSubmissions[i])
          break
        }
      }
    }
  }

  getAll() {
    this.visibleAssignmentSubmissions = this.allAssignmentSubmissions
  }

  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
  }

  redirectToDetails(talentID?: string, assignmentID?: string, isTalent: boolean = false) {
    this.assignmentProvider.allAssignmentSubmissions = this.allAssignmentSubmissions
    this.assignmentProvider.batchTalents = this.batchTalents

    this.urlService.setUrlFrom(this.backNavigation.BTA_DETAILS_TO_BTA_SCORES,
      this.router.url)
    let url : string
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN || this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
		  url = this.urlConstant.TRAINING_BATCH_MASTER_ASSIGNMENT_DETAILS
		}
		if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
			url = this.urlConstant.MY_BATCH_ASSIGNMENT_DETAILS
		}

    this.router.navigate([url], {
      // relativeTo: this.activatedRoute,
      queryParams: {
        "batchID": this.batchID,
        "assignmentID": assignmentID,
        "talentID": talentID,
        "isTalentView": isTalent
      }
    }).catch(err => {
      console.error(err)
    });
  }


}

// class batchTopicAssignment {
//   id?: string
//   batchID: string
//   programmingQuestion: IProgrammingQuestion
//   dueDate?: string
//   assignedDate: string
//   showDetails?: boolean
//   submissions: ITalentAssignmentSubmission[]
// }

// export class assignmentSubmission {
//   assignment: batchTopicAssignment
//   talentSubmission: Map<string, ITalentAssignmentSubmission>
// }