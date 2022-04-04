import { Component, OnInit } from '@angular/core';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;

@Component({
  selector: 'app-refer-and-earn',
  templateUrl: './refer-and-earn.component.html',
  styleUrls: ['./refer-and-earn.component.css']
})
export class ReferAndEarnComponent implements OnInit {

  // Spinner.



  // Refer link.
  referLink: string

  // Leader.
  leaderList: ILeader[]

  // Coin.
  coins: number

  // FAQ
  faqList: IFAQ[]

  // Description.
  description: string

  constructor(
    private spinnerService: SpinnerService,
  ) {
    this.initializeVariables()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize all global variables.
  initializeVariables() {

    // Spinner.
    this.spinnerService.loadingMessage = "Loading"


    //Refer link.
    this.referLink = "https://jbjvesbvjwbekjbvkjdbskvbskbvkbkjb/bcjhbsj"

    // Coin.
    this.coins = 1000

    // Description.
    this.description = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut " +
      "labore et dolore magna aliqua. Bibendum est ultricies integer quis. Iaculis urna id volutpat " +
      "lacus laoreet. Mauris vitae ultricies leo integer malesuada. Ac odio tempor orci dapibus " +
      "ultrices in. Egestas diam in arcu cursus euismod. Dictum fusce ut placerat orci nulla. " +
      "Tincidunt ornare massa eget egestas purus viverra accumsan in nisl. Tempor id eu nisl nunc" +
      " mi ipsum faucibus. Fusce id velit ut tortor pretium. Massa ultricies mi quis hendrerit dolor " +
      "magna eget. Nullam eget felis eget nunc lobortis. Faucibus ornare suspendisse sed nisi. " +
      "Sagittis eu volutpat odio facilisis mauris sit amet massa. Erat velit scelerisque in dictum " +
      "non consectetur a erat. Amet nulla facilisi morbi tempus iaculis urna. Egestas purus viverra " +
      "accumsan in nisl. Feugiat in ante metus dictum at tempor commodo. Convallis tellus id interdum " +
      "velit laoreet. Proin sagittis nisl rhoncus mattis rhoncus urna neque viverra. Viverra aliquet " +
      "eget sit amet tellus cras adipiscing enim eu. Ut faucibus pulvinar elementum integer enim neque " +
      "volutpat ac tincidunt. Dui faucibus in ornare quam. In iaculis nunc sed augue lacus viverra vitae " +
      "congue. Vitae tempus"

    // FAQ.
    this.faqList = [
      {
        id: "1",
        title: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore",
        description: this.description
      },
      {
        id: "2",
        title: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore",
        description: this.description
      },
      {
        id: "3",
        title: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore",
        description: this.description
      },
      {
        id: "3",
        title: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore",
        description: this.description
      }
    ]

    // Leader.
    this.leaderList = [
      {
        id: "1",
        image: "assets/images/rachel.jpg",
        name: "You",
        coins: 100,
        referrals: 100,
        rank: 1
      },
      {
        id: "2",
        image: "assets/images/rachel.jpg",
        name: "Rachel",
        coins: 200,
        referrals: 200,
        rank: 2
      },
      {
        id: "3",
        image: "assets/images/rachel.jpg",
        name: "Ross",
        coins: 300,
        referrals: 300,
        rank: 3
      },
      {
        id: "4",
        image: "assets/images/rachel.jpg",
        name: "Monica",
        coins: 400,
        referrals: 400,
        rank: 4
      }
    ]
  }

}

export interface ILeader {
  id?: string
  image: string,
  name: string
  coins: number
  referrals: number
  rank: number
}

export interface IFAQ {
  id?: string
  title: string
  description: string
  isVisible?: boolean
}
