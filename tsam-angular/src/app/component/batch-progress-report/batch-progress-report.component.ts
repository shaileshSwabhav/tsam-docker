import { Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import * as Chart from 'chart.js';
import { BatchTalentService } from 'src/app/service/batch-talent/batch-talent.service';
import { BatchTopicAssignmentService } from 'src/app/service/batch-topic-assignment/batch-topic-assignment.service';
import { BatchService } from 'src/app/service/batch/batch.service';
import { AccessLevel, Role, UrlConstant } from 'src/app/service/constant';
import { GeneralService } from 'src/app/service/general/general.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { InterviewScheduleService } from 'src/app/service/talent/interview-schedule.service';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { FacultyDashboardService } from 'src/app/service/faculty-dashboard/faculty-dashboard.service';
import { TalentDashboardService } from 'src/app/service/talent-dashboard/talent-dashboard.service';
import { ConceptModuleService } from 'src/app/service/concept-module/concept-module.service';
import { ConceptDashboardService } from 'src/app/service/concept-dashboard/concept-dashboard.service';

@Component({
  selector: 'app-batch-progress-report',
  templateUrl: './batch-progress-report.component.html',
  styleUrls: ['./batch-progress-report.component.css']
})
export class BatchProgressReportComponent implements OnInit {

  readonly EMPTY_ID = "00000000-0000-0000-0000-000000000000"
  public readonly PENDING_STATUS = "Pending"
  public readonly SUBMITTED_STATUS = "Submitted"
  public readonly COMPLETED_STATUS = "Completed"
  public readonly FACULTY_FEEDBACK_TYPE = "Talent_Session_Feedback"
  public readonly TALENT_FEEDBACK_TYPE = "Faculty_Session_Feedback"
  
  // Access.
  loginID: string
  talentID: string
  facultyID: string
  access: any
  roleFirstName: string
  isViewAllModules: boolean
  isTalentSelected: boolean

  // Batch.
  batchID: string
  batchDetails: any

  // Module.
  batchModuleList: any[]
  selectedModuleID: string
  selectedModuleForDates: any
  
  // Module concept.
  moduleConceptList: any[]
  easyModuleConceptCount: number
  mediumModuleConceptCount: number
  hardModuleConceptCount: number
  complexConceptList: any[]
  complexConceptMaxTalentList: any[]
  selectedModuleConceptID: string
  talentConceptRatingWithAssignmentList: any[]
  @ViewChild('complexConceptListModal') complexConceptListModal: any

  // Batch talent.
  batchTalentList: any[]
  talentInterview: any
  talentAverageRatingWeekly: any
  talentAverageRatingWeeklyInString: string

  // Batch topic assignment.
  batchTopicAssignmentList: any[]
  batchTopicAssignmentForModuleConceptList: any[]
  pendingCount: number 
  submittedCount: number 
  completedCount: number 
  highestScoredAssignment: any
  lowestScoredAssignment: any
  overdueAssignmentList: any[]
  @ViewChild('overdueAssignmentListModal') overdueAssignmentListModal: any

  // Talent rating.
  talentAverageRatingAllWeekByParameterList: any[]
  talentAverageRatingAllWeekList: any[]
  talentFeedbackQusetionID: string
  talentFeedbackQusetionList: any[]
  facultyFeedbackForTalentLeaderBoard: any[]
  @ViewChild('facultyFeedbackForTalentLedaerboardModal') facultyFeedbackForTalentLedaerboardModal: any

  // Faculty rating.
  facultyAverageRatingWeekly: any
  facultyAverageRatingAllWeekListByParameter: any[]
  facultyAverageRatingAllWeekList: any[]
  facultyFeedbackQusetionID: string
  facultyFeedbackQusetionList: any[]

  // Line chart.
  facultyWeeklyRatingGraph: Chart
  facultyWeeklyRatingGraphClass: string
  facultyRatingLabelList: string[]
  facultyRatingDataList: number[]
  talentWeeklyRatingGraph: Chart
  talentWeeklyRatingGraphClass: string
  talentRatingLabelList: string[]
  talentRatingDataList: any[]
  talentConceptRatingGraph: Chart
  talentConceptRatingGraphClass: string
  talentConceptRatingLabelList: string[]
  talentConceptRatingDataList: any[]

  // Donut chart.
  assignmentDonutChart: Chart
  assignmentDonutChartClass: string
  moduleConceptDonutChart: Chart
  moduleConceptDonutChartClass: string
  
  // Bar chart.
  facultyRatingBarChart: Chart
  facultyRatingBarChartClass: string
  facultyRatingLabelByParameterList: string[]
  facultyRatingDataByParameterList: number[]
  facultyRatingBackgroundByParameterList: string[]
  talentRatingBarChart: Chart
  talentRatingBarChartClass: string
  talentRatingLabelByParameterList: string[]
  talentRatingDataByParameterList: number[]
  talentRatingBackgroundByParameterList: string[]

  // Modal.
  modalRef: any

  constructor(
    private batchTopicAssignmentService: BatchTopicAssignmentService,
    private spinnerService: SpinnerService,
    private generalService: GeneralService,
    private batchService: BatchService,
    private batchTalentService: BatchTalentService,
    private route: ActivatedRoute,
    private interviewScheduleService: InterviewScheduleService,
    private localService: LocalService,
    private role: Role,
    private accessLevel: AccessLevel,
    private modalService: NgbModal,
    private facultyDashboardService: FacultyDashboardService,
    private talentDashboardService: TalentDashboardService,
    private conceptModuleService: ConceptModuleService,
    private conceptDashboardService: ConceptDashboardService,
    ) {
    this.initializeVariables()
  }

  ngOnInit(): void {

    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      this.getAllComponents()
    }, err => {
      console.error(err)
    })
  }

  // Initialize global variables.
  initializeVariables(): void {

    // Batch talent.
    this.batchTalentList = []
    this.talentID = null
    this.pendingCount = 0
    this.submittedCount = 0
    this.completedCount = 0

    // Access.
    this.isViewAllModules = true
    this.isTalentSelected = false
    this.loginID = this.localService.getJsonValue("loginID")
    if (this.role.TALENT == this.localService.getJsonValue("roleName")) {
      this.talentID = this.loginID
      this.access = this.accessLevel.ONLY_TALENT
    }
    if (this.role.FACULTY == this.localService.getJsonValue("roleName")) {
      this.facultyID = this.loginID
      this.access = this.accessLevel.ONLY_FACULTY
    }
    if (this.role.ADMIN == this.localService.getJsonValue("roleName") || this.role.SALES_PERSON == this.localService.getJsonValue("roleName")) {
      this.access = this.accessLevel.ADMIN_AND_SALESPERSON
    }
    this.roleFirstName = this.localService.getJsonValue("firstName")
    if (this.access.isFaculty){
      this.isViewAllModules = false
    }

    // Batch topic assignment.
    this.batchTopicAssignmentList = []
    this.batchTopicAssignmentForModuleConceptList = []
    this.overdueAssignmentList = []

    // Talent rating.
    this.talentAverageRatingAllWeekByParameterList = []
    this.talentAverageRatingAllWeekList = []
    this.talentRatingLabelByParameterList = []
    this.talentRatingDataByParameterList = []
    this.talentRatingBackgroundByParameterList = []
    this.talentRatingLabelList = []
    this.talentRatingDataList = []
    this.talentFeedbackQusetionList = []
    this.facultyFeedbackForTalentLeaderBoard = []

    // Faculty rating.
    this.facultyAverageRatingAllWeekListByParameter = []
    this.facultyAverageRatingAllWeekList  = []
    this.facultyRatingLabelByParameterList = []
    this.facultyRatingDataByParameterList = []
    this.facultyRatingBackgroundByParameterList = []
    this.facultyRatingLabelList = []
    this.facultyRatingDataList = []
    this.facultyFeedbackQusetionList = []

    // Module.
    this.batchModuleList = []

    // Module concept.
    this.moduleConceptList = []
    this.easyModuleConceptCount = 0
    this.mediumModuleConceptCount = 0
    this.hardModuleConceptCount = 0
    this.complexConceptList = []
    this.complexConceptMaxTalentList = []
    this.talentConceptRatingWithAssignmentList = []
    this.talentConceptRatingLabelList = []
    this.talentConceptRatingDataList = []
  }

  //********************************************* CREATE CHART FUNCTIONS ************************************************************

  // Format lists related to faculty rating graph.
  formatListsRelatedToFacultyRatingGraph(): void{
    this.facultyRatingLabelList = []
    this.facultyRatingDataList = []
    for (let i = 0; i < this.facultyAverageRatingAllWeekList.length; i++){
      this.facultyRatingLabelList.push((i+1).toString())
      this.facultyRatingDataList.push(this.facultyAverageRatingAllWeekList[i].rating)
    }
  }

  // Create faculty average rating for all weeks graph.
  createFacultyWeeklyRatingGraph() {
    this.formatListsRelatedToFacultyRatingGraph()
    if (this.facultyWeeklyRatingGraph){
      this.facultyWeeklyRatingGraph.data.labels = this.facultyRatingLabelList
      this.facultyWeeklyRatingGraph.data.datasets[0].data = this.talentConceptRatingDataList
      this.facultyWeeklyRatingGraph.update()
      return
    }
    this.facultyWeeklyRatingGraph = new Chart("facultyWeeklyRatingGraph", {
      type: 'line',
      data: {
        datasets: [
          {
            label: this.roleFirstName,
            data: this.facultyRatingDataList,
            fill: false,
            lineTension: 0.6,
            borderColor: "#1940d9",
            borderWidth: 1
          }],
        labels: this.facultyRatingLabelList,
      },
      options: {
        plugins: {
          datalabels: {
            display: false,
          }
        },
        tooltips: {
          enabled: true,
          mode: 'point',
          callbacks: {
              title: function(tooltipItems, data) {
                return 'Week ' + tooltipItems[0].xLabel
            },
            label: function(tooltipItems, data) { 
                let rating: number = Number(tooltipItems.yLabel)
                let ratingInString: string = rating.toString()
                if (rating % 1 != 0){
                ratingInString = rating.toFixed(1)
              }
              return "Rating: " + ratingInString + "/10"
            }
          }
        },
        legend:{
          display: false
        },
        title: {
          display: false
        },
        responsive: true,
        scales: {
          xAxes: [
            {
              gridLines: {
                color: "rgba(0, 0, 0, 0)",
              }
            }
          ],
          yAxes: [
            {
              ticks: {
                beginAtZero: true,
                max: 10,
                stepSize: 2
              }
            }
          ]
        }
      }
    })
  }

  // Format lists related to talent rating graph.
  formatLitsRelatedToTalentRatingBarChart(): void{
    this.talentRatingLabelList = []
    this.talentRatingDataList = []
    for (let i = 0; i < this.talentAverageRatingAllWeekList[0].weeklyAvgRating?.length; i++){
      this.talentRatingLabelList.push((i+1).toString())
    }
    for (let i = 0; i < this.talentAverageRatingAllWeekList.length; i++){
      let data: number[] = []
      for (let j = 0; j < this.talentAverageRatingAllWeekList[i].weeklyAvgRating?.length; j++){
        data.push(this.talentAverageRatingAllWeekList[i].weeklyAvgRating[j].rating)
      }
      this.talentRatingDataList.push({
        label: this.talentAverageRatingAllWeekList[i].firstName,
        data: data,
        fill: false,
        lineTension: 0.6,
        borderColor: this.getRandomRolor(),
        borderWidth: 1
      })
    }
  }

  // Create talent average rating for all weeks graph.
  createTalentWeeklyRatingGraph() {
    this.formatLitsRelatedToTalentRatingBarChart()
    if (this.talentWeeklyRatingGraph){
      this.talentWeeklyRatingGraph.data.labels = this.talentRatingLabelList
      this.talentWeeklyRatingGraph.data.datasets = this.talentRatingDataList
      this.talentWeeklyRatingGraph.update()
      return
    }
    this.talentWeeklyRatingGraph = new Chart("talentWeeklyRatingGraph", {
      type: 'line',
      data: {
        datasets: this.talentRatingDataList,
        labels: this.talentRatingLabelList,
      },
      options: {
        plugins: {
          datalabels: {
            display: false,
          }
        },
        tooltips: {
          enabled: true,
          mode: 'point',
          callbacks: {
              title: function(tooltipItems, data) {
                return 'Week ' + tooltipItems[0].xLabel;
            },
            label: function(tooltipItems, data) { 
              let rating: number = Number(tooltipItems.yLabel)
              let ratingInString: string = rating.toString()
              if (rating % 1 != 0){
                ratingInString = rating.toFixed(1)
              }
              return data.datasets[tooltipItems.datasetIndex].label + "'s Rating: " + ratingInString + "/10"
            }
          }
        },
        legend:{
          display: false
        },
        title: {
          display: false
        },
        responsive: true,
        scales: {
          xAxes: [
            {
              gridLines: {
                color: "rgba(0, 0, 0, 0)",
              }
            }
          ],
          yAxes: [
            {
              ticks: {
                beginAtZero: true,
                max: 10,
                stepSize: 2
              }
            }
          ]
        }
      }
    })
  }

  // Create assignment donut chart.
  createAssignmentDonutChart() {
    if (this.assignmentDonutChart){
      this.assignmentDonutChart.data.datasets[0].data = [this.pendingCount, this.submittedCount, this.completedCount]
      this.assignmentDonutChart.update()
      return
    }
    this.assignmentDonutChart = new Chart("assignmentDonutChart", {
      type: 'doughnut',
      data: {
        datasets: [{
          data: [this.pendingCount, this.submittedCount, this.completedCount],
          backgroundColor: [
            '#ff715b',
            '#31bacc',
            '#1940d9'
          ],
          borderWidth: 5,
          radius: 500,
          weight: 4,
        }],
        labels: [
          "Pending",
          "Submitted",
          "Completed"
        ],
      },
      options: {
        plugins: {
          datalabels: {
            display: false,
          }
        },
        legend:{
          display: false
        },
        cutoutPercentage: 85,
        responsive: true,
        maintainAspectRatio: false,
      },
    })
  }

  // Create module concept donut chart.
  createModuleConceptDonutChart() {
    if (this.moduleConceptDonutChart){
      this.moduleConceptDonutChart.data.datasets[0].data = [this.easyModuleConceptCount, this.mediumModuleConceptCount, this.hardModuleConceptCount]
      this.moduleConceptDonutChart.update()
      return
    }
    this.moduleConceptDonutChart = new Chart("moduleConceptDonutChart", {
      type: 'doughnut',
      data: {
        datasets: [{
          data: [this.easyModuleConceptCount, this.mediumModuleConceptCount, this.hardModuleConceptCount],
          backgroundColor: [
            '#5cad6d',
            '#feb72e',
            '#fd715b'
          ],
          radius: 500,
          weight: 4,
        }],
        labels: [
          "Easy",
          "Medium",
          "Hard"
        ],
      },
      options: {
        plugins: {
          datalabels: {
            display: false,
          }
        },
        legend:{
          display: false
        },
        cutoutPercentage: 60,
        responsive: true,
        maintainAspectRatio: false,
      },
    })
  }

  // Format lists related to faculty rating bar chart by parameter.
  formatListsRelatedToFacultyRatingBarChartByParameter(): void{
    this.facultyRatingLabelByParameterList = []
    this.facultyRatingDataByParameterList = []
    this.facultyRatingBackgroundByParameterList = []
    for (let i = 0; i < this.facultyAverageRatingAllWeekListByParameter.length; i++){
      this.facultyRatingLabelByParameterList.push((i+1).toString())
      this.facultyRatingDataByParameterList.push(this.facultyAverageRatingAllWeekListByParameter[i].rating)
      this.facultyRatingBackgroundByParameterList.push("#1940d9")
    }
  }

  // Create faculty rating bar chart.
  createFacultyRatingBarChart() {
    this.formatListsRelatedToFacultyRatingBarChartByParameter()
    if (this.facultyRatingBarChart){
      this.facultyRatingBarChart.data.labels = this.facultyRatingLabelByParameterList
      this.facultyRatingBarChart.data.datasets[0].data = this.facultyRatingDataByParameterList
      this.facultyRatingBarChart.update()
      return
    }
    this.facultyRatingBarChart = new Chart("facultyRatingBarChart", {
      type: 'bar',
      data: {
        labels: this.facultyRatingLabelByParameterList,
        datasets: [{
          barThickness: 10,
          data: this.facultyRatingDataByParameterList,
          backgroundColor: this.facultyRatingBackgroundByParameterList
        }]
      },
      options:{
        plugins: {
          datalabels: {
            display: false,
          }
        },
          tooltips: {
            enabled: true,
            mode: 'point',
            callbacks: {
                title: function(tooltipItems, data) {
                  return 'Week ' + tooltipItems[0].xLabel;
              },
              label: function(tooltipItems, data) { 
                let rating: number = Number(tooltipItems.yLabel)
                let ratingInString: string = rating.toString()
                if (rating % 1 != 0){
                  ratingInString = rating.toFixed(1)
                }
                return "Rating: " + ratingInString + "/10"
              }
            }
        },
        legend:{
          display: false
        },
        scales: {
          yAxes: [
            {
              ticks: {
                beginAtZero: true,
                max: 10,
                stepSize: 5
              },
              gridLines: {
                color: "rgba(0, 0, 0, 0)",
              }
            }
          ],
          xAxes: [
            {
              ticks: {
                beginAtZero: true,
              },
              gridLines: {
                color: "rgba(0, 0, 0, 0)",
              }
            },
          ]
        },
      },
    })
  }

  // Format lists related to talent rating bar chart.
  formatLitsRelatedToTalentRatingByParameterBarChart(): void{
    this.talentRatingLabelByParameterList = []
    this.talentRatingDataByParameterList = []
    this.talentRatingBackgroundByParameterList = []
    for (let i = 0; i < this.talentAverageRatingAllWeekByParameterList.length; i++){
      this.talentRatingLabelByParameterList.push(this.talentAverageRatingAllWeekByParameterList[i].firstName)
      this.talentRatingDataByParameterList.push(this.talentAverageRatingAllWeekByParameterList[i].averageRating)
      this.talentRatingBackgroundByParameterList.push("#1940d9")
    }
  }

  // Create talent rating bar chart.
  createTalentRatingBarChart() {
    this.formatLitsRelatedToTalentRatingByParameterBarChart()
    if (this.talentRatingBarChart){
      this.talentRatingBarChart.data.labels = this.talentRatingLabelByParameterList
      this.talentRatingBarChart.data.datasets[0].data = this.talentRatingDataByParameterList
      this.talentRatingBarChart.data.datasets[0].backgroundColor = this.talentRatingBackgroundByParameterList
      this.talentRatingBarChart.update()
      return
    }
    this.talentRatingBarChart = new Chart("talentRatingBarChart", {
      type: 'bar',
      data: {
        labels: this.talentRatingLabelByParameterList,
        datasets: [{
          barThickness: 10,
          data: this.talentRatingDataByParameterList,
          backgroundColor: this.talentRatingBackgroundByParameterList
        }]
      },
      options:{
        plugins: {
          datalabels: {
            display: false,
          }
        },
        tooltips: {
          enabled: true,
          mode: 'point',
          callbacks: {
              title: function(tooltipItems, data) {
                return ''
            },
            label: function(tooltipItems, data) { 
              let rating: number = Number(tooltipItems.yLabel)
              let ratingInString: string = rating.toString()
              if (rating % 1 != 0){
                ratingInString = rating.toFixed(1)
              }
              return tooltipItems.xLabel + "'s Rating: " + ratingInString + "/10"
            }
          }
        },
        legend:{
          display: false
        },
        scales: {
          yAxes: [
            {
              ticks: {
                beginAtZero: true,
                max: 10,
                stepSize: 5
              },
              gridLines: {
                color: "rgba(0, 0, 0, 0)",
              }
            }
          ],
          xAxes: [
            {
              ticks: {
                beginAtZero: true,
              },
              gridLines: {
                color: "rgba(0, 0, 0, 0)",
              }
            },
          ]
        },
      },
    })
  }

  // Format lists related to talent concept rating rating graph.
  formatLitsRelatedToTalentConceptRatingBarChart(): void{
    this.talentConceptRatingLabelList = []
    this.talentConceptRatingDataList = []
    for (let i = 0; i < this.batchTopicAssignmentForModuleConceptList.length; i++){
      this.talentConceptRatingLabelList.push(this.batchTopicAssignmentForModuleConceptList[i].programmingQuestion?.label)
    }
    for (let i = 0; i < this.talentConceptRatingWithAssignmentList.length; i++){
      let data: number[] = []
      for (let j = 0; j < this.talentConceptRatingWithAssignmentList[i].scoreList.length; j++){
        data.push(this.talentConceptRatingWithAssignmentList[i].scoreList[j])
      }
      this.talentConceptRatingDataList.push({
        label: this.talentConceptRatingWithAssignmentList[i].firstName,
        data: data,
        fill: false,
        lineTension: 0.6,
        borderColor: this.getRandomRolor(),
        borderWidth: 1
      })
    }
    if (this.talentConceptRatingLabelList.length > 0){
      this.talentConceptRatingGraphClass = ""
    }
    if (this.talentConceptRatingLabelList.length == 0){
      this.talentConceptRatingGraphClass = "hide-style"
    }
  }

  // Create talent concept rating graph.
  createTalentConceptRatingRatingGraph() {
    this.formatLitsRelatedToTalentConceptRatingBarChart()
    if (this.talentConceptRatingGraph){
      this.talentConceptRatingGraph.data.labels = this.talentConceptRatingLabelList
      this.talentConceptRatingGraph.data.datasets = this.talentConceptRatingDataList
      this.talentConceptRatingGraph.update()
      return
    }
    this.talentConceptRatingGraph = new Chart("talentConceptRatingGraph", {
      type: 'line',
      data: {
        datasets: 
        this.talentConceptRatingDataList,
        labels: this.talentConceptRatingLabelList,
      },
      options: {
        plugins: {
          datalabels: {
            display: false,
          }
        },
        tooltips: {
          enabled: true,
          mode: 'point',
          callbacks: {
              title: function(tooltipItems, data) {
                return tooltipItems[0].xLabel + ""
            },
            label: function(tooltipItems, data) { 
              return data.datasets[tooltipItems.datasetIndex].label + "'s Score: " + tooltipItems.yLabel + "/10"
            }
          }
        },
        legend:{
          display: false
        },
        title: {
          display: false
        },
        responsive: true,
        scales: {
          xAxes: [
            {
              gridLines: {
                color: "rgba(0, 0, 0, 0)",
              }
            }
          ],
          yAxes: [
            {
              ticks: {
                beginAtZero: true,
                max: 10,
                stepSize: 2
              }
            }
          ]
        }
      }
    })
  }

  //********************************************* FORMAT FUNCTIONS ************************************************************

  // Format batch details
  formatBatchDetails(): void{
    this.batchDetails.totalCompletedHours = Math.floor(this.batchDetails.totalCompletedHours/60)
    this.batchDetails.totalHours = Math.floor(this.batchDetails.totalHours/60)
    if (this.access.isTalent || this.isTalentSelected){
      this.batchDetails.attendedHours = Math.floor(this.batchDetails.attendedHours/60)
    }
  }

  // Format batch topic assignment list.
  formatBatchTopicAssignmentList(): void{
    this.overdueAssignmentList = []
    this.pendingCount = 0
    this.submittedCount = 0
    this.completedCount = 0
    this.highestScoredAssignment = null
    this.lowestScoredAssignment = null
    let highestScorePercent: number = 0
    let lowestScorePercent: number = 0
    let highestScoreDate: Date
    let lowestScoreDate: Date

    for (let i = 0; i < this.batchTopicAssignmentList.length; i++){

      // Create assignemnt submission for each student for each batch topic assignment. 
      this.batchTopicAssignmentList[i] = this.formatAssignmentMapForBatchTopicAssignmentList(this.batchTopicAssignmentList[i])

      for (let [key,value] of this.batchTopicAssignmentList[i].talentSubmissionMap){

        value.scorePercent = value.score/value.programmingQuestion.score*10
      
        if (value.score){
          
          // Get highest scored assignment.
          if (!highestScoreDate){
            highestScorePercent = value.scorePercent
            highestScoreDate = new Date(value.submittedOn)
            this.highestScoredAssignment = value
          }
          if (highestScoreDate && value.scorePercent > highestScorePercent){
            highestScorePercent = value.scorePercent
            highestScoreDate = new Date(value.submittedOn)
            this.highestScoredAssignment = value
          }
          if (highestScoreDate && value.scorePercent == highestScorePercent){
            let valueSubmittedOnDate: Date = new Date(value.submittedOn)
            if (valueSubmittedOnDate < highestScoreDate){
              highestScorePercent = value.scorePercent
              highestScoreDate = new Date(value.submittedOn)
              this.highestScoredAssignment = value
            }
          }

          // Get lowest scored assignment.
          if (!lowestScoreDate){
            lowestScorePercent = value.scorePercent
            lowestScoreDate = new Date(value.submittedOn)
            this.lowestScoredAssignment = value
          }
          if (lowestScoreDate && value.scorePercent < lowestScorePercent){
            lowestScorePercent = value.scorePercent
            lowestScoreDate = new Date(value.submittedOn)
            this.lowestScoredAssignment = value
          }
          if (lowestScoreDate && value.scorePercent == lowestScorePercent){
            let valueSubmittedOnDate: Date = new Date(value.submittedOn)
            if (valueSubmittedOnDate > lowestScoreDate){
              lowestScorePercent = value.scorePercent
              lowestScoreDate = new Date(value.submittedOn)
              this.lowestScoredAssignment = value
            }
          }
        }

        // Get the assignments that have submitted date greater than due date.
        let submittedDate: Date = new Date(value.submittedOn)
        let dueDate: Date = new Date(value.dueDate)
        let diffInDays = Math.floor((Date.UTC(submittedDate.getFullYear(), submittedDate.getMonth(), submittedDate.getDate()) 
          - Date.UTC(dueDate.getFullYear(), dueDate.getMonth(), dueDate.getDate()) ) /(1000 * 60 * 60 * 24))
        if (diffInDays > 0){
          this.overdueAssignmentList.push(value)
        }
        
        // Set completed status.
        if (value.isAccepted) {
          this.completedCount = this.completedCount + 1
          value.status = this.COMPLETED_STATUS
          continue
        }

        // Set submitted status.
        if (!value.isChecked) {
          this.submittedCount = this.submittedCount + 1
          value.status = this.SUBMITTED_STATUS
          continue
        }

        // Set pending status.
        this.pendingCount = this.pendingCount + 1
        value.status = this.PENDING_STATUS
      }
    }

    // Create assignment donut chart.
    this.createAssignmentDonutChart()
  }

  // Create assignemnt submission for each student for each batch topic assignment. 
  formatAssignmentMapForBatchTopicAssignmentList(batchTopicAssignment: any): void{
    let talentSubmissionMap: Map<string, any> = new Map()

    // If no submissions are there then increase the pending count.
    if (batchTopicAssignment.submissions && batchTopicAssignment.submissions.length == 0){
      if (!this.access.isTalent && !this.isTalentSelected){
        this.pendingCount = this.pendingCount + this.batchDetails.totalStudents
      }
      if (this.access.isTalent || this.isTalentSelected){
        this.pendingCount = this.pendingCount + 1
      }
    }
    for (let i = 0; i < batchTopicAssignment.submissions?.length; i++){
      if (!talentSubmissionMap.get(batchTopicAssignment.submissions[i].talent?.id)) {
        batchTopicAssignment.submissions[i].programmingQuestion = batchTopicAssignment.programmingQuestion
        batchTopicAssignment.submissions[i].dueDate = batchTopicAssignment.dueDate
        talentSubmissionMap.set(batchTopicAssignment.submissions[i].talent?.id, batchTopicAssignment.submissions[i])
      }
    }

    // If submissions are there but less than total students of batch then increase the pending count.
    if (!this.access.isTalent && !this.isTalentSelected && talentSubmissionMap.size > 0 && talentSubmissionMap.size < this.batchDetails.totalStudents){
      this.pendingCount = this.pendingCount + (this.batchDetails.totalStudents - talentSubmissionMap.size)
    }
    batchTopicAssignment.talentSubmissionMap = talentSubmissionMap
    return batchTopicAssignment
  }

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

  // Format module concept list.
  formatModuelConceptList(): void{
    this.easyModuleConceptCount = 0
    this.mediumModuleConceptCount = 0
    this.hardModuleConceptCount = 0

    for (let i = 0; i < this.moduleConceptList.length; i++){
      if (this.moduleConceptList[i].programmingConcept?.complexity == 1){
        this.easyModuleConceptCount = this.easyModuleConceptCount + 1
      }
      if (this.moduleConceptList[i].programmingConcept?.complexity == 2){
        this.mediumModuleConceptCount = this.mediumModuleConceptCount + 1
      }
      if (this.moduleConceptList[i].programmingConcept?.complexity == 3){
        this.hardModuleConceptCount = this.hardModuleConceptCount + 1
      }
    }
    this.createModuleConceptDonutChart()
  }

  // Format complex concept list.
  formatComplexConceptList(): void{
    this.complexConceptMaxTalentList = []
    let conceptCountMap: Map<string, any[]> = new Map()
    let halfBatchTalentCount: number = 0
    if (!this.access.isTalent && !this.isTalentSelected){
      halfBatchTalentCount = Math.ceil(this.batchDetails.totalStudents/2)
    }

    // Give complexity names to all concepts.
    for (let i = 0; i < this.complexConceptList.length; i++) {
      if (this.complexConceptList[i].complexity == 1) {
        this.complexConceptList[i].complexityName = "Easy"
      }
      if (this.complexConceptList[i].complexity == 2) {
        this.complexConceptList[i].complexityName = "Medium"
      }
      if (this.complexConceptList[i].complexity == 3) {
        this.complexConceptList[i].complexityName = "Hard"
      }
    }

    // Collect all unique concepts.
    for (let i = 0; i < this.complexConceptList.length; i++){
      if (conceptCountMap.get(this.complexConceptList[i].conceptName)) {
        let tempComplexConceptlist: any[] = conceptCountMap.get(this.complexConceptList[i].conceptName)
        tempComplexConceptlist.push(this.complexConceptList[i])
        conceptCountMap.set(this.complexConceptList[i].conceptName, tempComplexConceptlist)
      }
      if (!conceptCountMap.get(this.complexConceptList[i].conceptName)) {
        let tempComplexConceptlist: any[] = []
        tempComplexConceptlist.push(this.complexConceptList[i])
        conceptCountMap.set(this.complexConceptList[i].conceptName, tempComplexConceptlist)
      }
    }

    // Calculate how many concepts are complex for more than half the batch talents for faculty, if for talent then 
    // for himself how many concepts are present in array.
    for (let [key,value] of conceptCountMap){
      if (!this.access.isTalent && !this.isTalentSelected){
        if (value.length >= halfBatchTalentCount){
          this.complexConceptMaxTalentList.push(value)
        }
      }
      if (this.access.isTalent || this.isTalentSelected){
        this.complexConceptMaxTalentList.push(value)
      }
    }
  }

  // Format talent concept rating with assignment list.
  formatTalentConceptRatingWithAssignmentList(): void{

    // Get all batch topic assignments that have the selected module concept.
    let batchTopicAssignmentMap: Map<string, any> = new Map()
    for (let i = 0; i < this.batchTopicAssignmentList.length; i++){
      for (let j = 0; j < this.batchTopicAssignmentList[i].submissions?.length; j++){
        for (let k = 0; k < this.batchTopicAssignmentList[i].submissions[j].talentConceptRatings?.length; k++){
          if (this.batchTopicAssignmentList[i].submissions[j].talentConceptRatings[k].programmingConceptModule?.id == this.selectedModuleConceptID){
            if (batchTopicAssignmentMap.get(this.batchTopicAssignmentList[i]?.programmingQuestion?.id)) {
              continue
            }
            if (!batchTopicAssignmentMap.get(this.batchTopicAssignmentList[i].programmingQuestion?.id)) {
              batchTopicAssignmentMap.set(this.batchTopicAssignmentList[i].programmingQuestion?.id, this.batchTopicAssignmentList[i])
            }
          }
        }
      }
    }
    this.batchTopicAssignmentForModuleConceptList = []
    for (let [key,value] of batchTopicAssignmentMap){
      this.batchTopicAssignmentForModuleConceptList.push(value)
    }

    // Iterate all talents.
    for (let i = 0; i < this.talentConceptRatingWithAssignmentList.length; i++){

      // Create a bucket for collecting score for each batch topic assignment for that concept module.
      let scoreList: number[] = []

      // Iterate all batch topic assignemnts for the selected concept module.
      for (let j = 0; j < this.batchTopicAssignmentForModuleConceptList.length; j++){

        // Iterate all assignmnets of each talent.
        for (let k = 0; k < this.talentConceptRatingWithAssignmentList[i].assignments?.length; k++){

          // If same assignment then store the score in score list.
          if (this.batchTopicAssignmentForModuleConceptList[j].id == this.talentConceptRatingWithAssignmentList[i].assignments[k].assignmentID){
            scoreList.push(this.talentConceptRatingWithAssignmentList[i].assignments[k].score)
            continue
          }
        }

        // If no score for the batch topic assignment found then give 0.
        if (!scoreList[j]){
          scoreList.push(0)
        }
      }
      this.talentConceptRatingWithAssignmentList[i].scoreList = scoreList
    }
    this.createTalentConceptRatingRatingGraph()
  }

  //********************************************* OTHER FUNCTIONS ************************************************************

  // On changing batch talent selection.
  onBatchTalentChange(): void{
    if (!this.talentID){
      this.isTalentSelected = false
    }
    if (this.talentID){
      this.isTalentSelected = true
    }
    this.getAllComponents()
  }

  // On clicking view all overdue assignment list.
  onClickViewAllOverdureAssignmentList(): void{
    this.openModal(this.overdueAssignmentListModal, 'lg')
  }

  // On clicking view all complex concept list.
  onClickViewAllComplexConceptList(): void{
    this.openModal(this.complexConceptListModal, 'lg')
  }

  // On clicking view all faculty for talent ledaerboard.
  onClickViewAllFacultyFeedbackForTalentLedaerboard(): void{
    this.openModal(this.facultyFeedbackForTalentLedaerboardModal, 'lg')
  }

  // On changing faculty feedback question selection.
  onFacultyFeedbackQusetionChange(): void{
    this.getAllWeeksAverageFacultyRating(true)
  }

  // On changing talent feedback question selection.
  onTalentFeedbackQusetionChange(): void{
    this.getTalentFeedbackRatingForAllTalentsByParameter().then((data)=>{
      this.talentAverageRatingAllWeekByParameterList = data
      if (this.talentAverageRatingAllWeekByParameterList.length > 0){
        this.createTalentRatingBarChart()
      }
    }).catch((err)=>{
      console.error(err)
    })
  }

  // On batch module change.
  onBatchModuleChange(): void{
    this.getListsRelatedToBatchModules()
  }

  // On module concept change.
  onModuleConceptChange(): void{
    this.getTalentConceptRatingWithBatchTopicAssignment().then((data)=>{
      this.talentConceptRatingWithAssignmentList = data
      this.formatTalentConceptRatingWithAssignmentList()
    }).catch((err)=>{
      console.error(err)
    })
  }

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

  // Generate random dark colors.
  getRandomRolor() {
    var letters = '0123456789ABCDEF'.split('')
    var color = '#'
    for (var i = 0; i < 6; i++) {
        color += letters[Math.round(Math.random() * 15)]
    }
    return color
  } 

  // Toggle between view all modules and my modules for faculty.
  toggleIsViewAllModules(): void{
    this.isViewAllModules = !this.isViewAllModules
    this.getAllComponents()
  }

  // Toggle between talent and faculty.
  toggleTalentFaculty(): void{
    this.access.isTalent = !this.access.isTalent
    if (this.access.isTalent){
      this.isViewAllModules = true 
    }
    this.getAllComponents()
  }

  // On module for dates selection change.
  onSelectedModuleForDatesChange(): void{
    this.batchDetails.startDate = this.selectedModuleForDates.startDate
    this.batchDetails.estimatedEndDate = this.selectedModuleForDates.estimatedEndDate
  }
 
  //********************************************* GET FUNCTIONS ************************************************************

  // Get all components.
  getAllComponents(): void{
    if (this.batchID) {
      if (!this.access.isTalent && !this.isTalentSelected){
        this.getBatchDetailsAndItsRelatedLists()
        this.getTalentListOfBatch()
        this.getTalentWeeklyRating()
        this.getTalentFeedbackQusetionListAndItsRelatedLists()
        this.getFacultyFeedbackForTalentLeaderBoard()
        if (!this.isViewAllModules){
          this.getWeeklyFacultyAverageRating()
          this.getFacultyFeedbackQusetionList()
          this.getAllWeeksAverageFacultyRating(false)
        }
      }
      if (this.access.isTalent || this.isTalentSelected){
        this.getBatchTalentDetails()
        this.getTalentInterview()
        this.getTalentAverageRatingWeekly()
        this.getAllBatchTopicAssignmentsForTalent()
        this.getBatchModulesAndItsRelatedLists()
        this.getTalentWeeklyRating()
        this.getTalentFeedbackQusetionListAndItsRelatedLists()
        this.getFacultyFeedbackForTalentLeaderBoard()
      }
    }
  }

  //************************* Batch Details ************************************

  // Get batch details and its related lists by calling asynchronously.
  async getBatchDetailsAndItsRelatedLists(): Promise<void> {
		try {

      // Get batch details.
      this.batchDetails = await this.getBatchDetails()
      this.formatBatchDetails()

      // Get batch module list.
      await this.getBatchModulesAndItsRelatedLists()

      // Get all batch topic assignments.
			this.batchTopicAssignmentList = await this.getAllBatchTopicAssignmentsForOtherRoles()
      if (this.batchTopicAssignmentList.length > 0){
        this.assignmentDonutChartClass = ""
      }
      if (this.batchTopicAssignmentList.length == 0){
        this.assignmentDonutChartClass = "hide-style"
      }
      this.formatBatchTopicAssignmentList()

		} catch (error) {
			console.error(error)
		}
	}

  // Get batch details.
  async getBatchDetails(): Promise<any[]> {
		try {
			return await new Promise<any[]>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Getting Progress Report"
        let queryParms: any = {}
        if (!this.isViewAllModules){
          queryParms.facultyID = this.facultyID
        }
				this.batchService.getBatchDetails(this.batchID, queryParms).subscribe((response) => {
					resolve(response)
				}, (err: any) => {
					console.error(err)
					if (err.statusText.includes('Unknown')) {
						reject("No connection to server. Check internet.");
						return
					}
					reject(err.error.error)
				})
			})
		} finally {
		}
	}

  // Get batch talent details.
  getBatchTalentDetails(): void {
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."
    let queryParams: any = {}
    if (!this.isViewAllModules){
      queryParams.facultyID = this.facultyID
    }
    this.batchTalentService.getOneBatchTalentDetails(this.batchID, this.talentID, queryParams).subscribe((response) => {
      this.batchDetails = response
      this.formatBatchDetails()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get talent list of batch.
  getTalentListOfBatch(): void {
    this.spinnerService.loadingMessage = "Getting Progress Report"
    this.batchTalentService.getBatchTalentList(this.batchID).subscribe((response) => {
      this.batchTalentList = response.body
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get talent intreview.
  getTalentInterview(): void {
    this.spinnerService.loadingMessage = "Getting Progress Report"
    let queryParams: any = {
      talentID: this.talentID,
      status: "Selected"
    }
    this.interviewScheduleService.getInterview(queryParams).subscribe((response) => {
      this.talentInterview = response.body
      if (this.talentInterview.id == this.EMPTY_ID){
        this.talentInterview = null
      }
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get talent average rating for batch for current week.
  getTalentAverageRatingWeekly(): void {
    this.spinnerService.loadingMessage = "Getting Progress Report"
    let queryParams: any = {
      talentID: this.talentID
    }
    this.batchService.GetAverageRatingForTalent(this.batchID, this.talentID, queryParams).subscribe((response) => {
      this.talentAverageRatingWeekly = response.body
      if (this.talentAverageRatingWeekly && this.talentAverageRatingWeekly.averageRating){
        this.talentAverageRatingWeeklyInString = (this.talentAverageRatingWeekly.averageRating/this.talentAverageRatingWeekly.maxScore*10).toFixed(1)
      }
      if (this.talentAverageRatingWeekly.maxScore == null){
        this.talentAverageRatingWeekly = null
      }
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  //************************* Assignment Dashboard ************************************

  // Get all batch topic assignments for other roles.
  async getAllBatchTopicAssignmentsForOtherRoles(): Promise<any[]> {
		try {
			return await new Promise<any[]>((resolve, reject) => {
        this.batchTopicAssignmentList = []
        this.spinnerService.loadingMessage = "Getting Progress Report"
        let queryParams: any = {}
        if (!this.isViewAllModules){
          queryParams.facultyID = this.facultyID
        }
				this.batchService.getSessionAndAssignmentForBatch(this.batchID, queryParams).subscribe((response) => {
					resolve(response.body)
				}, (err: any) => {
					console.error(err)
					if (err.statusText.includes('Unknown')) {
						reject("No connection to server. Check internet.");
						return
					}
					reject(err.error.error)
				})
			})
		} finally {
		}
	}

  // Get all batch topic assignments for talent.
  getAllBatchTopicAssignmentsForTalent(): void {
    this.batchTopicAssignmentList = []
    this.spinnerService.loadingMessage = "Getting Progress Report"
    let queryParams: any = {}
    if (!this.isViewAllModules){
      queryParams.facultyID = this.facultyID
    }
    this.batchTopicAssignmentService.getAllBatchTopicAssignmentsForTalent(this.batchID, this.talentID, queryParams).subscribe((response) => {
      this.batchTopicAssignmentList = response.body
      this.formatBatchTopicAssignmentList()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }

  //************************* Faculty Feedback Dashboard ************************************

  // Get faculty feedback question list.
  getFacultyFeedbackQusetionList(): void {
    this.spinnerService.loadingMessage = "Getting Progress Report"
    this.generalService.getFeedbackQuestionByType(this.FACULTY_FEEDBACK_TYPE).subscribe((response) => {
      this.facultyFeedbackQusetionList = response.body
      if (this.facultyFeedbackQusetionList.length > 0){
        this.facultyFeedbackQusetionID = this.facultyFeedbackQusetionList[0].id
        this.getAllWeeksAverageFacultyRating(true)
      }
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }

	// Get avaerage weekly rating for talent to faculty.
  getWeeklyFacultyAverageRating(): void {
    this.spinnerService.loadingMessage = "Getting Progress Report"
    let queryParams: any = {
      facultyID: this.facultyID
    }
    this.facultyDashboardService.getWeeklyAverageRating(this.batchID, queryParams).subscribe((response) => {
      this.facultyAverageRatingWeekly = response.body
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }

	// Get avaerage weekly rating for talent to faculty for all weeks.
  getAllWeeksAverageFacultyRating(isSingleParameter: boolean): void {
    this.spinnerService.loadingMessage = "Getting Progress Report"
    let queryParams: any = {
      facultyID: this.facultyID
    }
    if (isSingleParameter){
      queryParams.questionID = this.facultyFeedbackQusetionID
    }
    this.facultyDashboardService.getAllWeeksAverageRating(this.batchID, queryParams).subscribe((response) => {
      if (isSingleParameter){
        this.facultyAverageRatingAllWeekListByParameter = response.body
        if (this.facultyAverageRatingAllWeekListByParameter.length > 0){
          this.facultyRatingBarChartClass = ""
          this.createFacultyRatingBarChart()
        }
        if (this.facultyAverageRatingAllWeekListByParameter.length == 0){
          this.facultyRatingBarChartClass = "hide-style"
        }
      }
      if (!isSingleParameter){
        this.facultyAverageRatingAllWeekList = response.body
        if (this.facultyAverageRatingAllWeekList.length > 0){
          this.facultyWeeklyRatingGraphClass = ""
          this.createFacultyWeeklyRatingGraph()
        }
        if (this.facultyAverageRatingAllWeekList.length == 0){
          this.facultyWeeklyRatingGraphClass = "hide-style"
        }
      }
     
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }

  //************************* Student Performance Dashboard ************************************

  // Get talent feedback qusetion list and its related lists.
  async getTalentFeedbackQusetionListAndItsRelatedLists(): Promise<void> {
		try {

      // Get talent feedback question list.
      this.talentFeedbackQusetionList = await this.getTalentFeedbackQusetionList()
      if (this.talentFeedbackQusetionList.length > 0){
        this.talentFeedbackQusetionID = this.talentFeedbackQusetionList[0].id
      }

      // Get talent feedback rating for all talents by parameter.
      this.talentAverageRatingAllWeekByParameterList = await this.getTalentFeedbackRatingForAllTalentsByParameter()
      if (this.talentAverageRatingAllWeekByParameterList.length > 0){
        this.talentRatingBarChartClass = ""
        this.createTalentRatingBarChart()
      }
      if (this.talentAverageRatingAllWeekByParameterList.length == 0){
        this.talentRatingBarChartClass = "hide-style"
      }

		} catch (error) {
			console.error(error)
		}
	}

  // Get talent feedback question list.
  async getTalentFeedbackQusetionList(): Promise<any[]> {
    try {
      return new Promise<any[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Progress Report"
        this.generalService.getFeedbackQuestionByType(this.TALENT_FEEDBACK_TYPE).subscribe((response) => {
          resolve(response.body)
        }, (err: any) => {
          console.error(err)
          if (err.statusText.includes('Unknown')) {
            reject("No connection to server. Check internet.");
            return
          }
          reject(err.error.error)
        })
      })
    } finally {
    }
  }

  // Get talent feedback rating for all talents by parameter.
  async getTalentFeedbackRatingForAllTalentsByParameter(): Promise<any[]> {
    try {
      return new Promise<any[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Progress Report"
        let queryParams: any = {
          questionID: this.talentFeedbackQusetionID
        }
        if (this.access.isTalent || this.isTalentSelected){
          queryParams.talentID = this.talentID
        }
        if (!this.access.isTalent && !this.isViewAllModules){
          queryParams.facultyID = this.facultyID
        }
        this.talentDashboardService.getTalentFeedbackRatingForAllTalents(this.batchID, queryParams).subscribe((response) => {
          resolve(response.body)
        }, (err: any) => {
          console.error(err)
          if (err.statusText.includes('Unknown')) {
            reject("No connection to server. Check internet.");
            return
          }
          reject(err.error.error)
        })
      })
    } finally {
    }
  }

  // Get talent feedback rating for all talents.
  getTalentWeeklyRating(): void {
    this.spinnerService.loadingMessage = "Getting Progress Report"
    let queryParams: any = {}
    if (this.access.isTalent || this.isTalentSelected){
      queryParams.talentID = this.talentID
    }
    if (!this.access.isTalent && !this.isViewAllModules){
      queryParams.facultyID = this.facultyID
    }
    this.talentDashboardService.getTalentWeeklyRatingForAllTalents(this.batchID, queryParams).subscribe((response) => {
      this.talentAverageRatingAllWeekList = response.body
      if (this.talentAverageRatingAllWeekList.length > 0){
        this.talentWeeklyRatingGraphClass = ""
        this.createTalentWeeklyRatingGraph()
      }
      if (this.talentAverageRatingAllWeekList.length == 0){
        this.talentWeeklyRatingGraphClass = "hide-style"
      }
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }

  // Get faculty feedback for talent leader board.
  getFacultyFeedbackForTalentLeaderBoard(): void {
    this.spinnerService.loadingMessage = "Getting Progress Report"
    let queryParams: any = {}
    if (!this.access.isTalent && !this.isViewAllModules){
      queryParams.facultyID = this.facultyID
    }
    this.talentDashboardService.getFacultyFeedbackForTalentLeaderBoard(this.batchID, queryParams).subscribe((response: any) => {
      this.facultyFeedbackForTalentLeaderBoard = response
      this.formatFacultyFeedbackForTalentLeaderBoard()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }

  //************************* Concepts Dashboard ************************************

  // Get batch module list and its related lists.
  async getBatchModulesAndItsRelatedLists(): Promise<void> {
		try {

      // Get batch module list.
      this.batchModuleList = await this.getBatchModuleList()
      if (this.batchModuleList.length > 0){
        this.selectedModuleID = this.batchModuleList[0].module?.id 
        this.selectedModuleForDates = this.batchModuleList[0]
        if (!this.access.isTalent && !this.isViewAllModules){
          this.batchDetails.startDate = this.selectedModuleForDates.startDate
          this.batchDetails.estimatedEndDate = this.selectedModuleForDates.estimatedEndDate
        }
      }

      // Get lists related to batch modules by calling asynchronously.
      await this.getListsRelatedToBatchModules()

		} catch (error) {
			console.error(error)
		}
	}

  // Get lists related to batch modules by calling asynchronously.
  async getListsRelatedToBatchModules(): Promise<void> {
		try {

      // Get module concept list.
			this.moduleConceptList = await this.getModuleConceptList()
      if (this.moduleConceptList.length > 0){
        this.selectedModuleConceptID = this.moduleConceptList[0].id
        this.formatModuelConceptList()
        this.moduleConceptDonutChartClass = ""
      }
      if (this.moduleConceptList.length == 0){
        this.moduleConceptDonutChartClass = "hide-style"
      }

      // Get complex concept list.
			this.complexConceptList = await this.getComplexConceptList()
      this.formatComplexConceptList()

      // GetTalentConceptRatingWithBatchTopicAssignment will get talent concept ratings for each batch topic assignemnt for all talents.
			this.talentConceptRatingWithAssignmentList = await this.getTalentConceptRatingWithBatchTopicAssignment()
      this.formatTalentConceptRatingWithAssignmentList()
      
		} catch (error) {
			console.error(error)
		}
	}

  // Get batch module list.
  async getBatchModuleList(): Promise<any[]> {
		try {
			return await new Promise<any[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Progress Report"
        let queryParams: any ={
          limit: -1,
          offset: 0
        }
        if (!this.access.isTalent && !this.isViewAllModules){
          queryParams.facultyID = this.facultyID
        }
				this.batchService.getBatchModules(this.batchID, queryParams).subscribe((response) => {
					resolve(response.body)
				}, (err: any) => {
					console.error(err)
					if (err.statusText.includes('Unknown')) {
						reject("No connection to server. Check internet.");
						return
					}
					reject(err.error.error)
				})
			})
		} finally {
		}
	}

  // Get module concept list.
  async getModuleConceptList(): Promise<any[]> {
		try {
			return await new Promise<any[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Progress Report"
        let queryParams: any = {
          batchID: this.batchID,
          moduleID: this.selectedModuleID,
          limit: -1,
          offset: 0
        }
				this.conceptModuleService.getAllModuleProgrammingConcepts(queryParams).subscribe((response) => {
					resolve(response.body)
				}, (err: any) => {
					console.error(err)
					if (err.statusText.includes('Unknown')) {
						reject("No connection to server. Check internet.");
						return
					}
					reject(err.error.error)
				})
			})
		} finally {
		}
	}

  // Get complex concept list.
  async getComplexConceptList(): Promise<any[]> {
		try {
			return await new Promise<any[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Progress Report"
        let queryParams: any = {
          moduleID: this.selectedModuleID,
        }
        if (this.access.isTalent || this.isTalentSelected){
          queryParams.talentID = this.talentID
        }
				this.conceptDashboardService.getComplexConcepts(queryParams).subscribe((response) => {
					resolve(response.body)
				}, (err: any) => {
					console.error(err)
					if (err.statusText.includes('Unknown')) {
						reject("No connection to server. Check internet.");
						return
					}
					reject(err.error.error)
				})
			})
		} finally {
		}
	}

  // GetTalentConceptRatingWithBatchTopicAssignment will get talent concept ratings for each batch topic assignemnt for all talents.
  async getTalentConceptRatingWithBatchTopicAssignment(): Promise<any[]> {
    try {
      return new Promise<any[]>((resolve, reject) => {
        this.spinnerService.loadingMessage = "Getting Progress Report"
        let queryParams: any ={
          moduleConceptID: this.selectedModuleConceptID
        }
        if (this.access.isTalent || this.isTalentSelected){
          queryParams.talentID = this.talentID
        }
        this.talentDashboardService.getTalentConceptRatingWithBatchTopicAssignment(this.batchID, queryParams).subscribe((response) => {
          resolve(response.body)
        }, (err: any) => {
          console.error(err)
          if (err.statusText.includes('Unknown')) {
            reject("No connection to server. Check internet.");
            return
          }
          reject(err.error.error)
        })
      })
    } finally {
    }
  }
}
