import { Component, OnInit, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-test',
  templateUrl: './test.component.html',
  styleUrls: ['./test.component.css']
})
export class TestComponent implements OnInit {

  // Spinner.



  // Test.
  testList: ITest[]
  activeTestList: ITest[]
  previousTestList: ITest[]
  selectedTest: any

  // Modal.
  modalRef: any
  @ViewChild('startNowModal') startNowModal: any

  constructor(
    private spinnerService: SpinnerService,
    private router: Router,
    private modalService: NgbModal,
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
    this.spinnerService.loadingMessage = "Getting Tests"


    // Test.
    this.activeTestList = []
    this.previousTestList = []
    this.selectedTest = {}
    this.testList = [
      {
        id: "1",
        name: "Interview Prep TA Test",
        difficulty: "Easy",
        time: 180,
        deadline: "11/10/2021",
        numberOfProblems: 20
      },
      {
        id: "2",
        name: "Data Structures And Algorithms TA Test",
        difficulty: "Hard",
        time: 160,
        deadline: "12/16/2021",
        numberOfProblems: 10
      },
      {
        id: "3",
        name: "Interview Prep TA Test - C++",
        difficulty: "Medium",
        time: 120,
        deadline: "09/10/2021",
        numberOfProblems: 30
      },
      {
        id: "4",
        name: "Scholarship Test - Summer 17",
        difficulty: "Hard",
        time: 90,
        deadline: "11/12/2021",
        numberOfProblems: 20
      },
      {
        id: "5",
        name: "Advanced Course Admission Test 2",
        difficulty: "Easy",
        time: 200,
        deadline: "11/04/2021",
        numberOfProblems: 10
      },
      {
        id: "6",
        name: "Intro To Programming Practice Test",
        difficulty: "Hard",
        time: 40,
        deadline: "12/10/2021",
        numberOfProblems: 20
      },
      {
        id: "7",
        name: "DS Practice Test",
        difficulty: "Medium",
        time: 50,
        deadline: "11/07/2021",
        numberOfProblems: 80
      },
      {
        id: "8",
        name: "Scholarship Test - Summer 14",
        difficulty: "Hard",
        time: 20,
        deadline: "12/07/2020",
        numberOfProblems: 100
      },
      {
        id: "9",
        name: "Advanced Course Admission Test 2",
        difficulty: "Easy",
        time: 40,
        deadline: "12/05/2020",
        numberOfProblems: 10
      },
      {
        id: "10",
        name: "Intro To Programming Practice Test",
        difficulty: "Hard",
        time: 60,
        deadline: "12/07/2020",
        numberOfProblems: 20
      },
      {
        id: "11",
        name: "DS Practice Test",
        difficulty: "Medium",
        time: 30,
        deadline: "09/07/2020",
        numberOfProblems: 10
      },
      {
        id: "12",
        name: "Scholarship Test - Summer 14",
        difficulty: "Hard",
        time: 120,
        deadline: "11/07/2020",
        numberOfProblems: 20
      }
    ]

    // Create the active and previous test lists.
    let today: Date = new Date()
    for (let i = 0; i < this.testList.length; i++) {
      let testDate: Date = new Date(this.testList[i].deadline)
      if (today > testDate) {
        this.previousTestList.push(this.testList[i])
        continue
      }
      this.activeTestList.push(this.testList[i])
    }

    // Set class for difficulty.
    for (let i = 0; i < this.testList.length; i++) {
      if (this.testList[i].difficulty == "Easy") {
        this.testList[i].difficultyClass = "card-sub-detail easy"
      }
      if (this.testList[i].difficulty == "Medium") {
        this.testList[i].difficultyClass = "card-sub-detail medium"
      }
      if (this.testList[i].difficulty == "Hard") {
        this.testList[i].difficultyClass = "card-sub-detail hard"
      }
    }
  }

  // On clicking start now button.
  startNow(testID: string): void {
    for (let i = 0; i < this.testList.length; i++) {
      if (this.testList[i].id == testID) {
        this.selectedTest = this.testList[i]
      }
    }
    this.openModal(this.startNowModal, 'lg')
  }

  // Redirect to test details page.
  redirectToTestDetails(testName: string): void {
    this.router.navigate(['/test/test-details'], {
      queryParams: {
        "testName": testName,
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Used to open modal.
  openModal(content: any, size?: string): NgbModalRef {
    if (!size) {
      size = 'xl'
    }
    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: size
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

}

export interface ITest {
  id?: string
  name: string
  difficulty: string
  difficultyClass?: string
  time: number
  deadline: string
  numberOfProblems: number
}
