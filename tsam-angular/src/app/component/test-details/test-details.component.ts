import { Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;

@Component({
  selector: 'app-test-details',
  templateUrl: './test-details.component.html',
  styleUrls: ['./test-details.component.css']
})
export class TestDetailsComponent implements OnInit {

  // Spinner.



  // Test.
  subTestList: ISubTest[]
  selectedSubTestName: string

  // Leader.
  leaderList: ILeader[]

  constructor(
    private spinnerService: SpinnerService,
    private router: Router,
    private activatedRoute: ActivatedRoute,
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
    this.spinnerService.loadingMessage = "Getting Test Details"


    // Test.
    this.subTestList = [
      {
        id: "1",
        name: "Find the sequenceA A!",
        difficulty: "Easy",
        maxScore: 30,
      },
      {
        id: "2",
        name: "Good elements",
        difficulty: "Hard",
        maxScore: 40,
      },
      {
        id: "3",
        name: "Cut the pastries",
        difficulty: "Medium",
        maxScore: 50,
      },
      {
        id: "4",
        name: "Can you make our guest happy?",
        difficulty: "Hard",
        maxScore: 50,
      },
      {
        id: "5",
        name: "Debug the given Code",
        difficulty: "Easy",
        maxScore: 70,
      },
      {
        id: "6",
        name: "Debug the given Code",
        difficulty: "Hard",
        maxScore: 50,
      },
      {
        id: "7",
        name: "Instructions for the Next Two Questions",
        difficulty: "Medium",
        maxScore: 80,
      },
      {
        id: "8",
        name: "Delete Node - Bug Explanation, if any",
        difficulty: "Hard",
        maxScore: 50,
      },
      {
        id: "9",
        name: "Interview Prep TA Test - C++",
        difficulty: "Easy",
        maxScore: 50,
      },
      {
        id: "10",
        name: "Data Structures And Algorithms TA Test",
        difficulty: "Hard",
        maxScore: 100,
      },
      {
        id: "11",
        name: "Interview Prep TA Test - C++",
        difficulty: "Medium",
        maxScore: 50,
      },
      {
        id: "12",
        name: "Interview Prep TA Test - C++",
        difficulty: "Hard",
        maxScore: 40,
      }
    ]

    // Leader.
    this.leaderList = [
      {
        id: "1",
        image: "assets/images/rachel.jpg",
        name: "You",
        rewards: "assets/images/medals.png",
        score: 100,
        rank: 1
      },
      {
        id: "2",
        image: "assets/images/rachel.jpg",
        name: "Rachel",
        rewards: "assets/images/medals.png",
        score: 200,
        rank: 2
      },
      {
        id: "3",
        image: "assets/images/rachel.jpg",
        name: "Ross",
        rewards: "assets/images/medals.png",
        score: 300,
        rank: 3
      },
      {
        id: "4",
        image: "assets/images/rachel.jpg",
        name: "Monica",
        rewards: "assets/images/medals.png",
        score: 400,
        rank: 4
      }
    ]

    // Set class for difficulty.
    for (let i = 0; i < this.subTestList.length; i++) {
      if (this.subTestList[i].difficulty == "Easy") {
        this.subTestList[i].difficultyClass = "opacity-86-style color-212121 font-weight-bold font-sm-style easy"
      }
      if (this.subTestList[i].difficulty == "Medium") {
        this.subTestList[i].difficultyClass = "opacity-86-style color-212121 font-weight-bold font-sm-style medium"
      }
      if (this.subTestList[i].difficulty == "Hard") {
        this.subTestList[i].difficultyClass = "opacity-86-style color-212121 font-weight-bold font-sm-style hard"
      }
    }

    // Get the test name from query params.
    this.selectedSubTestName = this.activatedRoute.snapshot.queryParamMap.get("testName")
  }

  // Redirect to problem details.
  redirectToProblemDetails(subTestID: string): void {
    this.router.navigate(['/test/test-details/problem-details'], {
      queryParams: {
        "subTestID": subTestID,
      }
    }).catch(err => {
      console.error(err)
    })
  }

}

export interface ISubTest {
  id?: string
  name: string
  difficulty: string
  difficultyClass?: string
  maxScore: number
}

export interface ILeader {
  id?: string
  image: string,
  name: string
  rewards: string
  score: number
  rank: number
}

