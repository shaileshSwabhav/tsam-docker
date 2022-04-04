import { Component, OnInit } from '@angular/core';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { FacultyDashboardService, IBarchart, IFacultyBatchDetails, IOngoingBatchDetails, IPiechart } from 'src/app/service/faculty-dashboard/faculty-dashboard.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { DashboardService, IFeedback, IGroupWiseKeywordName, ISessionGroupScore, ITalentFeedbackScore } from 'src/app/service/dashboard/dashboard.service';
import { FeedbackGroupService } from 'src/app/service/feedback-group/feedbak-group.service';
import { AdminService, ITimesheetActivity } from 'src/app/service/admin/admin.service';
import { Router } from '@angular/router';
import { DatePipe } from '@angular/common';
import { ChartDataSets, ChartOptions, ChartType } from 'chart.js';
import * as ChartDataLabels from 'chartjs-plugin-datalabels';
import { Label } from 'ng2-charts';
import { IFacultyReport, ReportService } from 'src/app/service/report/report.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
	selector: 'app-faculty-dashboard',
	templateUrl: './faculty-dashboard.component.html',
	styleUrls: ['./faculty-dashboard.component.css']
})
export class FacultyDashboardComponent implements OnInit {

	// components
	talentTypeList: any[]

	// faculty details
	loginID: string
	facultyID: string
	facultyList: any[]

	// batch-details
	batchDetails: IFacultyBatchDetails
	ongoingBatchDetails: IOngoingBatchDetails[]

	// batch-talents
	batchTalentPerformanceDetails: ITalentFeedbackScore[]

	// feedback-keyword
	feebackKeyword: IGroupWiseKeywordName[]
	talentGroupScore: ISessionGroupScore
	isGroupScoreVisible: boolean

	// task list
	facultyTimesheetActivity: ITimesheetActivity[]
	completedActivity: ITimesheetActivity[]
	pendingActivity: ITimesheetActivity[]

	// piechart
	timesheetPiechart: IPiechart[]
	isPiechartDataPresent: boolean

	// barchart data.
	mentorGraphData: IBarchart

	// spinner

	ongoingOperationList: number

	// nav
	studentDetailsTab: number;
	mentorTab: number;
	taskListTab: number;
	workScheduleTab: number;

	// workschedule timetable
	facultyTimeTable: IFacultyReport[]

	// workschedule flags
	isTimetableVisible: boolean
	isPiechartVisible: boolean

	// charts
	workScheduleData: number[]

	// piechart
	piechartLabels: Label[]
	barchartPlugins = []
	piechartOptions: ChartOptions
	piechartData: ChartDataSets[]
	piechartType: ChartType
	piechartLegend: boolean
	piechartPlugins: any
	piechartColors: any

	// barchart
	barchartOptions: ChartOptions
	barchartLabels: Label[]
	barchartType: ChartType
	barchartData: ChartDataSets[]
	barchartLegend: boolean
	barchartColors: any

	// loading
	workScheduleLabel: string

	public readonly ONGOINGBATCHES: string = "Ongoing"
	public readonly COMPLETEDBATCHES: string = "Finished"
	public readonly UPCOMINGBATCHES: string = "Upcoming"

	constructor(
		private facultyDashboardService: FacultyDashboardService,
		private dashboardService: DashboardService,
		private generalService: GeneralService,
		private feedbackGroupService: FeedbackGroupService,
		private reportService: ReportService,
		private adminService: AdminService,
		public utilService: UtilityService,
		private localService: LocalService,
		private spinnerService: SpinnerService,
		private router: Router,
		private datePipe: DatePipe,
	) {
		this.initializeVariables()
		this.initializePiechart()
		this.initializeBarchart()
		this.getAllComponents()
	}

	initializeVariables(): void {
		this.loginID = this.localService.getJsonValue("loginID")
		this.facultyID = this.loginID

		this.ongoingOperationList = 0
		this.facultyList = []
		this.ongoingBatchDetails = []
		this.feebackKeyword = []
		this.batchTalentPerformanceDetails = []
		this.facultyTimesheetActivity = []
		this.completedActivity = []
		this.pendingActivity = []
		this.timesheetPiechart = []
		this.talentTypeList = []

		this.studentDetailsTab = 1
		this.mentorTab = 1
		this.taskListTab = 1
		this.workScheduleTab = 1

		// chart
		this.workScheduleData = []
		this.piechartLabels = [];

		this.isTimetableVisible = true
		this.isPiechartVisible = false
		this.isPiechartDataPresent = true

		this.workScheduleLabel = "Work Schedule This Week"
	}

	redirectToBatches(params:string):void{
		let queryParams ={
			batchStatus:params
		}
		this.router.navigate(["/my/batch"], {
			queryParams: queryParams
		}).catch(err => {
			console.error(err)

		});
	}

	redirectToTalents():void{
		this.router.navigate(["/my/talent"], {
			// queryParams: queryParams
		}).catch(err => {
			console.error(err)

		});
	}

	initializePiechart(): void {
		let numSectors = 5;
		let sectorDegree = 180 / numSectors;

		// ng2-charts piechart
		this.piechartOptions = {
			responsive: true,
			legend: {
				position: "right"
			},
			// scale: {
			// 	ticks: {
			// 		fontColor: "white",
			// 		fontSize: 20,
			// 	}
			// },
			plugins: {
				datalabels: {
					formatter: (value, ctx) => {
						let total = 0
						// console.log("ctx -> ", ctx.chart.data.labels[ctx.dataIndex]);

						for (let index = 0; index < this.workScheduleData.length; index++) {
							total += this.workScheduleData[index]
						}

						// const label = ctx.chart.data.labels[ctx.dataIndex];
						let score = (this.workScheduleData[ctx.dataIndex] / total) * 100
						return Math.round(score).toString() + "%";
					},
					align: 'center',
					anchor: 'center',
					color: "#fff",
					font: {
						weight: 'bold',
						size: 20,
					}
				}
			},
			elements: {
				line: {
					fill: false
				},
				point: {
					hoverRadius: 7,
					radius: 5
				}
			},

		};

		// public piechartData: SingleDataSet = [30, 10, 20, 40];
		this.piechartData = [
			{ data: this.workScheduleData, label: 'Approved', stack: 'a' },
		];
		this.piechartType = 'pie';
		this.piechartLegend = true;
		this.piechartPlugins = [ChartDataLabels];
		this.piechartColors = [{
			backgroundColor: ['#F26A00', '#EF6C00', '#FFA726', '#FB8C00', '#F57C00']
		}]
	}

	initializeBarchart(): void {

		this.barchartOptions = {
			responsive: true,
			plugins: {
				datalabels: {
					display: "none",
				}
			},
			scales: {
				yAxes: [{
					ticks: {
						beginAtZero: true,
						max: 80
					},
					gridLines: {
						display: false
					},
				}],
				xAxes: [{
					gridLines: {
						display: false,
					}
				}]
			},
		};
		this.barchartLabels = ['Total Students', 'Fresher', 'Professional'] //, 'Students Placed'
		this.barchartType = 'bar';
		this.barchartLegend = false;
		this.barchartData = [{
			data: [],
			label: 'Students',
			backgroundColor: "#F8BAA7"
		}];
		this.barchartColors = [{
			backgroundColor: ['#F8BAA7', '#F8BAA7', '#F8BAA7'] //, '#F8BAA7'
		}]
	}

	getAllComponents(): void {
		this.getFacultyBatchDetails()
		this.getFacultyList()
		this.getTalentType()
		this.getOngoingBatchDetails()
		this.getGroupwiseKeywordNames()
		this.getTodayTimesheet()
		this.getAllFacultyReport()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit(): void {
	}

	onStudentDetailsTabChange(event: any): void {
		if (this.ongoingBatchDetails.length >= event - 1) {
			// this.getBatchTalent(this.ongoingBatchDetails[event - 1]?.batchID)
			this.getTalentFeedback(this.ongoingBatchDetails[event - 1]?.batchID)
		}
	}

	onMentorTabChange(event: any): void {
		this.getBarchartData(this.ongoingBatchDetails[event - 1].batchID)
	}

	onTaskListChange(event: any): void {
		if (this.taskListTab == 1) {
			this.getTodayTimesheet()
		}

		if (this.taskListTab == 2) {
			this.getPendingTimesheet()
		}
	}

	onTimetableClick(): void {
		if (this.isTimetableVisible) {
			return
		}
		this.workScheduleLabel = "Work Schedule This Week"
		this.isPiechartVisible = false
		this.isTimetableVisible = true
		this.getAllFacultyReport()
	}

	onPiechartClick(): void {
		if (this.isPiechartVisible) {
			return
		}
		this.isTimetableVisible = false
		this.isPiechartVisible = true
		this.workScheduleLabel = "Time Spent Analysis"
		this.getPiechartData()
	}

	/**
	 * getMonday
	 * @returns date depending on the weekIndex value, 0 is for this week,
	 *  -7 is for previous week and 7 for next week.
	 */
	getMonday(): Date {
		let date = new Date();
		let day = date.getDay();
		let finalDate = new Date();

		let offset = 1 - day;
		if (offset > 0) {
			offset = -6;
		}

		finalDate.setDate(new Date().getDate() + offset);
		return finalDate
	}

	showGroupScore(performance: ITalentFeedbackScore, groupScore: ISessionGroupScore): void {
		groupScore.showGroupDetails = !groupScore.showGroupDetails
		performance.feedbackKeywords.find((val) => {
			if (groupScore.groupName != val.groupName && val.showGroupDetails) {
				val.showGroupDetails = !val.showGroupDetails
			}
		})
	}

	redirectToTimesheet(): void {
		let queryParams: any

		if (this.taskListTab == 1) {
			queryParams = {
				fromDate: this.datePipe.transform(new Date(), "yyyy-MM-dd"),
				toDate: this.datePipe.transform(new Date(), "yyyy-MM-dd"),
				credentialID: this.localService.getJsonValue("credentialID")
			}
		}

		if (this.taskListTab == 2) {
			queryParams = {
				credentialID: this.localService.getJsonValue("credentialID"),
				isCompleted: "null",
			}
		}

		this.router.navigate(["/my/timesheet"], {
			queryParams: queryParams
		}).catch(err => {
			console.error(err)

		});
	}

	getFacultyBatchDetails(): void {
		this.spinnerService.loadingMessage = "Loading dashboard details"

		this.ongoingOperationList++

		this.facultyDashboardService.getFacultyBatchDetails(this.facultyID).subscribe((response: any) => {
			this.batchDetails = response.body
		}, (err: any) => {
			console.error(err);
		}).add(() => {
			this.ongoingOperationList--
			if (this.ongoingOperationList == 0) {
				;
			}
		})
	}

	getOngoingBatchDetails(): void {
		this.spinnerService.loadingMessage = "Loading dashboard details"

		this.ongoingOperationList++

		this.facultyDashboardService.getOngoingBatchDetails(this.facultyID).subscribe((response: any) => {
			this.ongoingBatchDetails = response.body
			if (this.ongoingBatchDetails && this.ongoingBatchDetails.length > 0) {
				// this.getBatchTalent(this.ongoingBatchDetails[0]?.batchID)
				this.getTalentFeedback(this.ongoingBatchDetails[0]?.batchID)
				this.getBarchartData(this.ongoingBatchDetails[0]?.batchID)
			}
		}, (err: any) => {
			console.error(err);
		}).add(() => {
			this.ongoingOperationList--
			if (this.ongoingOperationList == 0) {
				;
			}
		})
	}

	getTodayTimesheet(): void {
		let queryParams: any = {
			fromDate: this.datePipe.transform(new Date(), "yyyy-MM-dd"),
			toDate: this.datePipe.transform(new Date(), "yyyy-MM-dd"),
			// credentialID: this.localService.getJsonValue("credentialID")
		}
		this.getFacultyTimesheet(queryParams)
	}

	getPendingTimesheet(): void {
		let date = new Date()
		date.setDate(new Date().getDate() - 1)

		let queryParams: any = {
			nextEstimatedDate: this.datePipe.transform(new Date(), "yyyy-MM-dd"),
			toDate: this.datePipe.transform(date, "yyyy-MM-dd"),
			// credentialID: this.localService.getJsonValue("credentialID")
		}
		this.getFacultyTimesheet(queryParams)
	}

	getFacultyTimesheet(queryParams?: any): void {
		this.spinnerService.loadingMessage = "Loading..."

		this.ongoingOperationList++
		this.facultyTimesheetActivity = []
		this.facultyDashboardService.getTaskList(queryParams).subscribe((response: any) => {
			this.facultyTimesheetActivity = response.body
			// console.log(this.facultyTimesheetActivity);
		}, (err: any) => {
			console.error(err)
		}).add(() => {
			this.ongoingOperationList--
			if (this.ongoingOperationList == 0) {
				;
			}
		})
	}


	getPiechartData(): void {
		let date: Date = this.getMonday()
		date.setDate(this.getMonday().getDate() + 6)

		let queryParams: any = {
			fromDate: this.datePipe.transform(this.getMonday(), "yyyy-MM-dd"),
			toDate: this.datePipe.transform(date, "yyyy-MM-dd"),
			credentialID: this.localService.getJsonValue("credentialID")
		}
		this.spinnerService.loadingMessage = "Generating piechart..."

		this.isPiechartDataPresent = true
		// this.timesheetPiechart = []
		this.workScheduleData = []
		this.piechartLabels = []
		this.facultyDashboardService.getPiechartData(this.facultyID, queryParams).subscribe((response: any) => {
			this.timesheetPiechart = response.body
			let timesheetMap = new Map<string, number>()

			if (this.timesheetPiechart.length > 0) {
				for (let index = 0; index < this.timesheetPiechart.length; index++) {
					if (index < 4) {
						timesheetMap.set(this.timesheetPiechart[index].projectName, this.timesheetPiechart[index].hours)
						continue
					}

					if (timesheetMap.get("Other") > 0) {
						let count = timesheetMap.get("Other") + this.timesheetPiechart[index].hours
						timesheetMap.set("Other", count)
					} else {
						timesheetMap.set("Other", this.timesheetPiechart[index].hours)
					}
				}

				for (let [key, value] of timesheetMap) {
					this.workScheduleData.push(value)
					this.piechartLabels.push(key)
				}

				this.piechartData = [
					{ data: this.workScheduleData, label: 'Approved', stack: 'a' },
				];
			}
			// console.log("timesheetPiechart -> ",this.timesheetPiechart);
			// console.log("timesheetMap -> ",timesheetMap);
			// console.log("workScheduleData -> ",this.workScheduleData);
			// console.log("piechartData -> ",this.piechartData);

		}, (err: any) => {
			console.error(err)
		}).add(() => {
			if (this.timesheetPiechart.length == 0) {
				this.isPiechartDataPresent = false
			}
		})
	}

	feedback: IFeedback

	getTalentFeedback(batchID: string): void {
		this.feedback = null
		this.dashboardService.getTalentFeedbackScore(batchID).subscribe((response: any) => {
			this.feedback = response.body
			console.log("feedback -> ", this.feedback);
		}, (err: any) => {
			console.error(err)
		})
	}

	// getBatchTalent(batchID: string): void {
	// 	this.batchTalentPerformanceDetails = []
	// 	this.spinnerService.loadingMessage = "Loading..."

	// 	this.ongoingOperationList++
	// 	this.dashboardService.getBatchTalents(batchID).subscribe((response: any) => {
	// 		this.batchTalentPerformanceDetails = response.body
	// 		for (let index = 0; index < this.batchTalentPerformanceDetails.length; index++) {
	// 			this.batchTalentPerformanceDetails[index].isVisibile = false
	// 		}
	// 		console.log(this.batchTalentPerformanceDetails);
	// 		console.log(this.feebackKeyword);

	// 	}, (err: any) => {
	// 		console.error(err)
	// 	}).add(() => {
	// 		this.ongoingOperationList--
	// 		if (this.ongoingOperationList == 0) {
	// 			;
	// 		}
	// 	})
	// }

	getBarchartData(batchID: string): void {
		this.spinnerService.loadingMessage = "Loading..."

		this.ongoingOperationList++

		let queryParams: any = {
			batchID: batchID
		}

		this.facultyDashboardService.getBarchartData(this.facultyID, queryParams).subscribe((response: any) => {
			this.mentorGraphData = response.body
			// console.log(this.mentorGraphData);

			this.barchartData = [{
				data: [this.mentorGraphData.totalStudents, this.mentorGraphData.fresher, this.mentorGraphData.professional,
				this.mentorGraphData.studentsPlaced],
				label: 'Students',
				backgroundColor: "#F8BAA7",
				stack: 'a'
			}]

			this.barchartOptions = {
				responsive: true,
				plugins: {
					datalabels: {
						display: false,
					}
				},
				scales: {
					yAxes: [{
						ticks: {
							beginAtZero: true,
							max: this.mentorGraphData.totalStudents + 5
						},
						gridLines: {
							display: false
						},
					}],
					xAxes: [{
						gridLines: {
							display: false,
						}
					}]
				}
			};
			// console.log(this.barchartData);

		}, (err: any) => {
			console.error(err);
		}).add(() => {
			this.ongoingOperationList--
			if (this.ongoingOperationList == 0) {
				;
			}
		})
	}

	getAllFacultyReport(): void {
		this.spinnerService.loadingMessage = "Loading..."

		this.facultyTimeTable = []
		let queryParams: any = {
			facultyID: this.facultyID,
			date: this.getMonday().toISOString()
		}
		// console.log(queryParams);
		this.isTimetableVisible = true

		this.reportService.getAllFacultyReport(queryParams).subscribe((response: any) => {
			this.facultyTimeTable = response.body
			if (this.facultyTimeTable.length == 0) {
				this.isTimetableVisible = false
			}
			// console.log(this.facultyTimeTable);
		}, (err: any) => {
			this.facultyTimeTable = []
			console.error(err);
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
		}).add(() => {
			if (this.facultyTimeTable.length == 0) {
				this.isTimetableVisible = false
			}
		})
	}


	toggleGroupVisibility(group: IGroupWiseKeywordName): void {
		group.showGroupDetails = !group.showGroupDetails
	}

	calculateAvgGroupScore(groupScore: ISessionGroupScore): number {
		let score = 0
		for (let index = 0; index < groupScore.feedbackScore.length; index++) {
			score += groupScore.feedbackScore[index].keywordScore
		}
		return score / groupScore.feedbackScore.length
	}

	// Get All Faculty List
	getFacultyList(): void {
		this.generalService.getFacultyList().subscribe((data: any) => {
			this.facultyList = data.body
		}, (err) => {
			console.error(err);
		})
	}

	getGroupwiseKeywordNames(): void {
		this.feedbackGroupService.getGroupwiseKeywordNames().subscribe((data: any) => {
			this.feebackKeyword = data.body
		}, (err) => {
			console.error(err);
		})
	}

	// Get talent type list.
	getTalentType(): void {
		this.generalService.getGeneralTypeByType("talent_rating").subscribe((respond: any[]) => {
			this.talentTypeList = respond
		}, (err) => {
			console.error(this.utilService.getErrorString(err))
		})
	}


}
