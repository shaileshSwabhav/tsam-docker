import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IProblemOfTheDayQuestion } from 'src/app/service/problem-of-the-day/problem-of-the-day.service';
import { ProgrammingQuestionService } from 'src/app/service/programming-question/programming-question.service';
import { Location } from '@angular/common';

@Component({
  selector: 'app-practice-problem-list',
  templateUrl: './practice-problem-list.component.html',
  styleUrls: ['./practice-problem-list.component.css']
})
export class PracticeProblemListComponent implements OnInit {

  // Programming Question.
  problemList: IProblemOfTheDayQuestion[]

  // Pagination.
  limit: number
  currentPage: number
  totalProblems: number
  offset: number
  paginationStart: number
  paginationEnd: number

  // Spinner.



  // Programming concept.
  subConceptID: string

  constructor(
    private spinnerService: SpinnerService,
    private activatedRoute: ActivatedRoute,
    private questionService: ProgrammingQuestionService,
    private router: Router,
    private _location: Location,
  ) {
    this.initializeVariables()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    // Programming Question.
    this.problemList = [] as IProblemOfTheDayQuestion[]

    // Pagination.
    this.limit = 10
    this.offset = 0
    this.currentPage = 0
    this.paginationStart = 0
    this.paginationEnd = 0

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Problems"


    // Programming concept.
    this.subConceptID = this.activatedRoute.snapshot.queryParamMap.get("subConceptID")
    if (this.subConceptID) {
      this.getProgrammingQuestionListBySubProgrammingConceptID()
    }
  }

  // =============================================================FORMAT FUNCTIONS==========================================================================

  // Format concepts, sub concepts and problems on problem list selection. 
  formatOnProblemList(): void {

    for (let i = 0; i < this.problemList.length; i++) {

      // Set class for difficulty.
      if (this.problemList[i].level == 1) {
        this.problemList[i].levelName = "Easy"
        this.problemList[i].levelClass = "easy"
      }
      if (this.problemList[i].level == 2) {
        this.problemList[i].levelName = "Medium"
        this.problemList[i].levelClass = "medium"
      }
      if (this.problemList[i].level == 3) {
        this.problemList[i].levelName = "Hard"
        this.problemList[i].levelClass = "hard"
      }

      // Set success ratio.
      if (this.problemList[i].attemptedByCount != 0) {
        this.problemList[i].successRatio = ((this.problemList[i].solvedByCount / this.problemList[i].attemptedByCount) * 100).toFixed(2)
        if (this.problemList[i].successRatio == "0.00") {
          this.problemList[i].successRatio = "0"
        }
      }
      else {
        this.problemList[i].successRatio = "0"
      }
    }
  }

  // ================================================OTHER FUNCTIONS FOR PROBLEM===============================================

  // Page change.
  changePage(pageNumber: number): void {
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    this.getProgrammingQuestionListBySubProgrammingConceptID()
  }

  // Set total talents list on current page.
  setPaginationString(): void {
    this.paginationStart = this.limit * this.offset + 1
    this.paginationEnd = +this.limit + this.limit * this.offset
    if (this.totalProblems < this.paginationEnd) {
      this.paginationEnd = this.totalProblems
    }
  }

  // On clicking solve button.
  onSolveProblemClick(problemID: string): void {
    this.router.navigate(['/practice/problem-details'], {
      queryParams: {
        "subConceptID": this.subConceptID,
        "problemID": problemID,
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Nagivate to previous url.
  backToPreviousPage(): void {
    this._location.back()
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get question list by sub programming concept id. 
  getProgrammingQuestionListBySubProgrammingConceptID(): void {
    this.spinnerService.loadingMessage = "Getting Practice Problems"


    let queryParams: any = {
      "programmingConceptID": this.subConceptID
    }
    this.questionService.getProgrammingQuestionsPractice(6, 0, queryParams).subscribe((response: any) => {
      this.problemList = response.body
      this.totalProblems = response.headers.get('X-Total-Count')
      this.formatOnProblemList()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.setPaginationString()
    })
  }

}
