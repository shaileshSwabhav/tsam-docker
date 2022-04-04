import { Component, OnInit, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { AdminService, IEvent } from 'src/app/service/admin/admin.service';
import { BatchService } from 'src/app/service/batch/batch.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-events',
  templateUrl: './events.component.html',
  styleUrls: ['./events.component.css']
})
export class EventsComponent implements OnInit {

  // tabs
  eventTabs: IEventTab[]

  // Event.
  eventList: IEvent[]

  // Batch
  batchList: any[]

  // Source.
  sourceList: any[]

  // Spinner.



  // Pagination.
  totalEvents: number
  totalBatches: number
  limit: number
  currentPage: number
  offset: number
  paginationString: string

  // Flags.
  isFirstLoading: boolean

  // Modal.
  modalRef: any
  @ViewChild('applyNowModal') applyNowModal: any

  // constant
  readonly LIVE_EVENTS = 0
  readonly UPCOMING_EVENTS = 1
  readonly COMPLETED_EVENTS = 2
  readonly UPCOMING_BATCHES = 3

  constructor(
    private spinnerService: SpinnerService,
    private router: Router,
    private batchService: BatchService,
    private generalService: GeneralService,
    private adminService: AdminService,
    public utilService: UtilityService,
  ) {
    this.initializeVariables()
    this.getSourceList()
    this.getUpcomingBatches()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  // Initialize global variables.
  initializeVariables(): void {

    // Spinner.


    // Event.
    this.eventList = []

    // Batch.
    this.batchList = []

    // Source.
    this.sourceList = []

    this.createEventTab()

    // Pagination.
    this.limit = 6
    this.offset = 0
    this.currentPage = 0

    // Flags
    this.isFirstLoading = true
  }

  // NOTE: If sequence is changed here then make sure to handle callFunctionForCurrentTab function
  // and also handle isActive checks in html
  createEventTab(): void {
    this.eventTabs = [
      { tabName: "Live Events", isActive: true },
      { tabName: "Upcoming Events", isActive: false },
      { tabName: "Completed Events", isActive: false },
      { tabName: "Upcoming Batches", isActive: false },
    ]
    this.callFunctionForCurrentTab(this.LIVE_EVENTS)
  }

  // Set total entries list on current page.
  setPaginationString(total: number) {
    this.paginationString = ''
    let start: number = this.limit * this.offset + 1
    let end: number = +this.limit + this.limit * this.offset
    if (total < end) {
      end = total
    }
    if (total == 0) {
      this.paginationString = ''
      return
    }
    this.paginationString = `${start} - ${end} of ${total}`
  }

  // On page change.
  changePage(pageNumber: number): void {
    let activeIndex: number
    this.currentPage = pageNumber
    this.offset = this.currentPage - 1
    for (let index = 0; index < this.eventTabs.length; index++) {
      if (this.eventTabs[index].isActive) {
        activeIndex = index
      }
    }
    this.callFunctionForCurrentTab(activeIndex)
  }

  onEventTabClick(activeIndex: number) {
    if (this.eventTabs[activeIndex].isActive) {
      return
    }

    for (let index = 0; index < this.eventTabs.length; index++) {
      this.eventTabs[index].isActive = false
    }

    this.eventTabs[activeIndex].isActive = true
    this.callFunctionForCurrentTab(activeIndex)
  }

  callFunctionForCurrentTab(index: number): void {
    let queryParams: any

    switch (index) {
      case this.LIVE_EVENTS:
        queryParams = {
          isActive: "1",
          eventStatus: "Live"
        }
        this.getAllEvents(queryParams)
        break;
      case this.UPCOMING_EVENTS:
        queryParams = {
          isActive: "1",
          eventStatus: "Upcoming"
        }
        this.getAllEvents(queryParams)
        break;
      case this.COMPLETED_EVENTS:
        queryParams = {
          isActive: "1",
          eventStatus: "Completed"
        }
        this.getAllEvents(queryParams)
        break;
      case this.UPCOMING_BATCHES:
        this.getUpcomingBatches()
        break;
      default:
        break;
    }
  }

  // Calculate the last registration date of all upcoming batches.
  calculateFieldsOfBatchList(): void {

    for (let i = 0; i < this.batchList.length; i++) {

      // Last registration date.
      this.batchList[i].lastRegDate = new Date(this.batchList[i].startDate)
      this.batchList[i].lastRegDate.setDate(this.batchList[i].lastRegDate.getDate() - 7)

      // Sale.
      this.batchList[i].isSale = false
    }
  }

  // Get upcoming batches.
  getUpcomingBatches(): void {
    this.spinnerService.loadingMessage = "Getting Upcoming Batches"

    this.batchService.getUpcomingBatches().subscribe((response) => {
      this.totalBatches = response.headers.get('X-Total-Count')
      this.batchList = response.body
      this.calculateFieldsOfBatchList()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(error.error.error)
    }).add(() => {

      this.setPaginationString(this.totalBatches)
    })
  }

  // Get all events.
  getAllEvents(queryParams: any): void {
    this.spinnerService.loadingMessage = "Getting events"

    this.totalEvents = 0
    this.eventList = []
    this.adminService.getAllEvent(queryParams).subscribe((response: any) => {
      this.eventList = response.body
      this.totalEvents = this.eventList.length
      if (this.isFirstLoading && this.totalEvents == 0) {
        this.eventTabs[this.LIVE_EVENTS].isActive = false
        this.eventTabs[this.UPCOMING_EVENTS].isActive = true
        this.isFirstLoading = false
        this.callFunctionForCurrentTab(this.UPCOMING_EVENTS)
        return
      }
      this.isFirstLoading = false
    }, (err: any) => {
      this.totalEvents = 0
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    }).add(() => {
      this.setPaginationString(this.totalEvents)

    })
  }

  // Get source list.
  getSourceList(): void {
    this.generalService.getSources().subscribe(response => {
      this.sourceList = response
    }, error => {
      console.error(error)
    })
  }

  // Redirect to event-details page.
  redirectToEventDetails(eventID: string): void {
    this.router.navigate(['events/event-details'], {
      queryParams: {
        "eventID": eventID,
      }
    }).catch(err => {
      console.error(err)
    })
  }

  // Redirect to batch-details page.
  redirectToBatchDetails(batchID: string, courseID: string): void {
    this.router.navigate(['events/batch-details'], {
      queryParams: {
        "batchID": batchID,
        "courseID": courseID
      }
    }).catch(err => {
      console.error(err)
    })
  }


}

export interface IEventTab {
  tabName: string
  isActive: boolean
}