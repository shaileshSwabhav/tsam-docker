import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, FormControl } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { IBatchProject } from 'src/app/models/batch/project';
import { IProjectSubmission } from 'src/app/models/project-submission';
import { AssignmentData } from 'src/app/providers/assignment-data';
import { BatchProjectService } from 'src/app/service/batch-project/batch-project.service';
import { BatchTalentService } from 'src/app/service/batch-talent/batch-talent.service';
import { BackNavigationUrl, Role, UrlConstant } from 'src/app/service/constant';
import { IPermission } from 'src/app/service/menu/menu.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ITalent} from 'src/app/service/talent/talent.service';
import { UrlService } from 'src/app/service/url.service';
import { UtilityService } from 'src/app/service/utility/utility.service';


@Component({
  selector: 'app-batch-project-score',
  templateUrl: './batch-project-score.component.html',
  styleUrls: ['./batch-project-score.component.css']
})
export class BatchProjectScoreComponent implements OnInit {

  // search form
  searchScoreForm: FormGroup

  // batch-session-assignment-score
  //talentProjectScores: any[]

  batchProjects: IBatchProject[]

  batchTalents: ITalent[]

  // batch-session programming assingment
  // sessionAssignmentList: any[]

  // access
  permission: IPermission

  // batch-details
  batchID: string
  batchName: string

  // other


  allProjectSubmissions: IProjectSubmission[]
  limit: number = 5
  currentPage: number
  offset: number = 0
  options: string[] = ["Not Scored", "Completed", "Pending", "All"]
  option: string = this.options[3]
  visibleProjectSubmissions: IProjectSubmission[]

  constructor(
    private formBuilder: FormBuilder,
    private urlConstant: UrlConstant,
    private bproService: BatchProjectService,
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

    // this.talentProjectScores = []
    this.batchTalents = []
    this.batchProjects = []
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
    this.getAllProjectsWithScore()
    this.getBatchTalents()
  }

  getAllProjectsWithScore(): void {
    this.spinnerService.loadingMessage = "Getting talent scores"
    this.batchProjects = []

    let queryParam: any = this.searchScoreForm.value

    if (!queryParam.dueDate) {
      delete queryParam.dueDate
    }
    console.log("batchID",this.batchID)
    console.log("searchScoreForm",this.searchScoreForm)
    this.bproService.getProjectWithSubmissions(this.batchID, this.searchScoreForm.value).subscribe((response: any) => {
      this.batchProjects = response.body
      console.log("batchProjects",this.batchProjects)
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
      console.log("batchTalents",this.batchTalents)
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
    this.allProjectSubmissions = []
    let no = 0
    for (let index = 0; index < this.batchProjects.length; index++) {
      this.allProjectSubmissions.push({ project: this.batchProjects[index], talentSubmission: new Map })
      for (let talIndex = 0; talIndex < this.batchTalents.length; talIndex++) {
        let sub = this.batchProjects[index]?.submissions?.find((value) => {
          if (value.talent.id === this.batchTalents[talIndex].id) {
            return value
          }
        })
        if (!sub) {
          this.allProjectSubmissions[index].talentSubmission.set(this.batchTalents[talIndex].id, null)
          continue
        }
        this.allProjectSubmissions[index].talentSubmission.set(this.batchTalents[talIndex].id, sub)
      }
    }
    console.log("allProjectSubmissions",this.allProjectSubmissions)
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
    this.visibleProjectSubmissions = []
    for (let i = 0; i < this.allProjectSubmissions.length; i++) {
      for (let value of this.allProjectSubmissions[i].talentSubmission.values()) {
        if (value && !value.isChecked) {
          this.visibleProjectSubmissions.push(this.allProjectSubmissions[i])
          break
        }
      }
    }
  }

  getCompleted() {
    this.visibleProjectSubmissions = []
    for (let i = 0; i < this.allProjectSubmissions.length; i++) {
      let num = 0
      for (let value of this.allProjectSubmissions[i].talentSubmission.values()) {
        if (value && value.isAccepted) {
          num++
          if (num == this.batchTalents.length)
            this.visibleProjectSubmissions.push(this.allProjectSubmissions[i])
          continue
        }
        break
      }
    }
  }

  getPending() {
    this.visibleProjectSubmissions = []
    for (let i = 0; i < this.allProjectSubmissions.length; i++) {
      for (let value of this.allProjectSubmissions[i].talentSubmission.values()) {
        if (!value || (value.isChecked && !value.isAccepted)) {
          this.visibleProjectSubmissions.push(this.allProjectSubmissions[i])
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
    this.visibleProjectSubmissions = []
    let today = new Date()
    for (let i = 0; i < this.allProjectSubmissions.length; i++) {
      let dueDate = new Date(this.allProjectSubmissions[i].project.dueDate)
      for (let value of this.allProjectSubmissions[i].talentSubmission.values()) {
        if ((!value || (value.isChecked && !value.isAccepted)) &&
          dueDate <= today) {
          this.visibleProjectSubmissions.push(this.allProjectSubmissions[i])
          break
        }
      }
    }
  }

  getAll() {
    this.visibleProjectSubmissions = this.allProjectSubmissions
    console.log("visibleProjectSubmissions",this.visibleProjectSubmissions)
  }

  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
  }

  redirectToDetails(talentID?: string, projectID?: string, isTalent: boolean = false) {
    this.assignmentProvider.allProjectSubmissions = this.allProjectSubmissions
    this.assignmentProvider.batchTalents = this.batchTalents

    this.urlService.setUrlFrom(this.backNavigation.BP_DETAILS_TO_BP_SCORES,
      this.router.url)
    let url: string
    if (this.localService.getJsonValue("roleName") == this.role.ADMIN || this.localService.getJsonValue("roleName") == this.role.SALES_PERSON) {
      url = this.urlConstant.TRAINING_BATCH_MASTER_PROJECT_DETAILS
    }
    if (this.localService.getJsonValue("roleName") == this.role.FACULTY) {
      url = this.urlConstant.MY_BATCH_PROJECT_DETAILS
    }
    this.router.navigate([url], {
      // relativeTo: this.activatedRoute,
      queryParams: {
        "batchID": this.batchID,
        "projectID": projectID,
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