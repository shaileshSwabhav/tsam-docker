import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IProblemOfTheDayQuestion } from 'src/app/service/problem-of-the-day/problem-of-the-day.service';
import { ProgrammingConceptService } from 'src/app/service/programming-concept/programming-concept.service';
import { ProgrammingQuestionService } from 'src/app/service/programming-question/programming-question.service';
import { IProgrammingConcept } from 'src/app/models/programming/concept';

@Component({
  selector: 'app-practice',
  templateUrl: './practice.component.html',
  styleUrls: ['./practice.component.css']
})
export class PracticeComponent implements OnInit {

  // Spinner.



  // Programming concept.
  conceptList: IProgrammingConcept[]
  selectedConcept: any
  selectedConceptID: string

  // Sub programming concept.
  selectedSubConcept: any

  // Questions.
  problemList: IProblemOfTheDayQuestion[]
  totalProblems: number

  constructor(
    private spinnerService: SpinnerService,
    private router: Router,
    private programmingConceptService: ProgrammingConceptService,
    private questionService: ProgrammingQuestionService,
  ) {
    this.initializeVariables()
    this.getAllConcepts()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Questions.
    this.problemList = []

    // Programming concept.
    this.conceptList = []
    this.selectedConcept = {}

    // Spinner.
    this.spinnerService.loadingMessage = "Getting All Programming Concepts"


    // Sub prgramming concept.
    this.selectedSubConcept = {}
  }

  // =============================================================FORMAT FUNCTIONS==========================================================================

  // Set fields for concept list.
  formatConceptList(): void {

    for (let i = 0; i < this.conceptList.length; i++) {

      // Give default logo to concept if logo is null.
      // if (!this.conceptList[i].logo){
      //   this.conceptList[i].logo = "assets/logo/default-programming-concept-logo.png"
      // }

      // Set class for difficulty.
      if (this.conceptList[i].complexity == 1) {
        this.conceptList[i].levelName = "Easy"
      }
      if (this.conceptList[i].complexity == 2) {
        this.conceptList[i].levelName = "Medium"
      }
      if (this.conceptList[i].complexity == 3) {
        this.conceptList[i].levelName = "Hard"
      }
    }
  }

  // Set fields for problem list.
  formatProblemList(): void {
    for (let i = 0; i < this.problemList.length; i++) {

      // Set class for difficulty.
      if (this.problemList[i].level == 1) {
        this.problemList[i].levelName = "Easy"
      }
      if (this.problemList[i].level == 2) {
        this.problemList[i].levelName = "Medium"
      }
      if (this.problemList[i].level == 3) {
        this.problemList[i].levelName = "Hard"
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

  // Set the default class for sub programming concept.
  formatConceptAndSubConcept(): void {

    for (let i = 0; i < this.conceptList.length; i++) {

      // Make sub concept card class as deafult (orange).
      // for (let j = 0; j < this.conceptList[i].subProgrammingConcepts.length; j++) {
      //   this.conceptList[i].subProgrammingConcepts[j].subConceptClass = "card practice-card-style h-100"
      // }

      // // Give default logo to concept if logo is null.
      // if (!this.conceptList[i].logo){
      //   this.conceptList[i].logo = "assets/logo/default-programming-concept-logo.png"
      // }
    }
  }

  // Format concepts, sub concepts and problems on problem list selection. 
  formatOnProblemListSelection(conceptID: string, subConceptID: string): void {

    // Get selected concept.
    for (let i = 0; i < this.conceptList.length; i++) {
      if (this.conceptList[i].id == conceptID) {
        this.selectedConcept = this.conceptList[i]
        break
      }
    }

    // Get selected sub concept and its problem list.
    for (let i = 0; i < this.selectedConcept.subProgrammingConcepts.length; i++) {
      if (this.selectedConcept.subProgrammingConcepts[i].id == subConceptID) {
        this.selectedSubConcept = this.selectedConcept.subProgrammingConcepts[i]
        break
      }
    }

    for (let i = 0; i < this.problemList.length; i++) {

      // Set class for difficulty.
      if (this.problemList[i].level == 1) {
        this.problemList[i].levelName = "Easy"
        this.problemList[i].levelClass = "opacity-86-style color-212121 font-weight-bold font-sm-style easy"
      }
      if (this.problemList[i].level == 2) {
        this.problemList[i].levelName = "Medium"
        this.problemList[i].levelClass = "opacity-86-style color-212121 font-weight-bold font-sm-style medium"
      }
      if (this.problemList[i].level == 3) {
        this.problemList[i].levelName = "Hard"
        this.problemList[i].levelClass = "opacity-86-style color-212121 font-weight-bold font-sm-style hard"
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

    // Make isVisible property of concept true to make the problem list visible.
    this.selectedConcept.isVisible = true

    // Make sub concept card class as disabled (grey) except for selected sub concept.
    // for (let i = 0; i < this.conceptList.length; i++) {
    //   for (let j = 0; j < this.conceptList[i].subProgrammingConcepts.length; j++) {
    //     if (this.conceptList[i].subProgrammingConcepts[j].id == this.selectedSubConcept.id) {
    //       continue
    //     }
    //     this.conceptList[i].subProgrammingConcepts[j].subConceptClass = "card practice-diasbled-card-style h-100"
    //   }
    // }
  }

  // =============================================================REDIRECT FUNCTIONS==========================================================================

  // // Redirect to problem details page.
  // redirectToProblemDetails(problemID: string, subConceptID: string): void {
  //   this.router.navigate(['/practice/problem-details'], {
  //     queryParams: {
  //       "problemID": problemID,
  //       "subConceptID": subConceptID,
  //     }
  //   }).catch(err => {
  //     console.error(err)
  //   })
  // }

  // Redirect to problem details page.
  redirectToProblemDetails(problemID: string): void {
    this.router.navigate(['/practice/problem-details'], {
      queryParams: {
        "problemID": problemID,
        "subConceptID": this.selectedConceptID,
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to sub programming concept question list page.
  redirectToSubProgrammingConceptQuestionsPage(subConceptID: string): void {
    this.router.navigate(['/practice/practice-problem-list'], {
      queryParams: {
        "subConceptID": subConceptID,
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // =============================================================OTHER FUNCTIONS==========================================================================

  // On clicking close button.
  onCloseButtonClick(): void {
    this.selectedConcept.isVisible = false
    this.selectedConcept = {}
    this.formatConceptAndSubConcept()
  }

  // =============================================================GET FUNCTIONS==========================================================================

  // Get all programming concepts by limit and offset.
  getAllConcepts() {
    this.spinnerService.loadingMessage = "Getting All Programming Concepts"


    let queryParams: any = {
      limit: -1,
      offset: 0
    }
    this.programmingConceptService.getAllConcepts(queryParams).subscribe((response) => {

      this.conceptList = response.body
      this.formatConceptList()
    }, error => {

      console.error(error);
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get question list by sub programming concept id. 
  getProgrammingQuestionListBySubProgrammingConceptID(conceptID: string, subConceptID: string): void {
    this.spinnerService.loadingMessage = "Getting Practice Problems"


    let queryParams: any = {
      "programmingConceptID": subConceptID
    }
    this.questionService.getProgrammingQuestionsPractice(6, 0, queryParams).subscribe((response: any) => {
      this.problemList = response.body
      this.totalProblems = response.headers.get('X-Total-Count')
      this.formatOnProblemListSelection(conceptID, subConceptID)
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  // Get question list by programming concept id. 
  getProgrammingQuestionListByProgrammingConceptID(conceptID: string): void {

    // If concept is already selected then unselect it.
    if (conceptID == this.selectedConceptID) {
      this.selectedConceptID = null
      return
    }

    this.selectedConceptID = conceptID
    this.spinnerService.loadingMessage = "Getting Practice Problems"


    let queryParams: any = {
      "programmingConceptID": conceptID
    }
    this.questionService.getProgrammingQuestionsPractice(6, 0, queryParams).subscribe((response: any) => {
      this.problemList = response.body
      this.totalProblems = response.headers.get('X-Total-Count')
      this.formatProblemList()
      // this.formatOnProblemListSelection(conceptID, subConceptID)
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

}