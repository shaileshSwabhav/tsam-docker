import { Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { CompanyRequirementService } from 'src/app/service/company/company-requirement/company-requirement.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IWaitingList, TalentService } from 'src/app/service/talent/talent.service';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { Location } from '@angular/common';
import { GeneralService } from 'src/app/service/general/general.service';

@Component({
  selector: 'app-company-details',
  templateUrl: './company-details.component.html',
  styleUrls: ['./company-details.component.css']
})
export class CompanyDetailsComponent implements OnInit {

  // Company branch.
  companyBranch: any

  // Company.
  company: any

  // Company requirement.
  companyRequirementID: string
  companyRequirementDetails: any
  technologies: string
  packageInString: string

  // Talent.
  talentID: string
  email: string

  // Spinner.



  // Source.
  sourceList: any[]

  // Application.
  isApplied: boolean

  // Waiting list.
  waitingList: IWaitingList[]

  // Modal.
  modalRef: any
  @ViewChild('applyNowModal') applyNowModal: any

  // Background image.
  buildingPath: string
  buildingPathArray: string[]

  constructor(
    private activatedRoute: ActivatedRoute,
    private companyRequirementService: CompanyRequirementService,
    private spinnerService: SpinnerService,
    private talentService: TalentService,
    private modalService: NgbModal,
    private _location: Location,
    private generalService: GeneralService,
  ) {
    this.initializeVariables()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Company Details"


    // Application.
    this.isApplied = false

    // Waiting list.
    this.waitingList = []

    // Source.
    this.sourceList = []

    // Comapny branch.
    this.companyBranch = {}

    // Company.
    this.company = {}

    // Company requirement.
    this.companyRequirementDetails = {}

    // Background image.
    this.buildingPathArray = ["assets/images/building1.jpg", "assets/images/building2.jpg", "assets/images/building3.jpg",
      "assets/images/building4.jpg", "assets/images/building5.jpg"]
    const random = Math.floor(Math.random() * this.buildingPathArray.length);
    this.buildingPath = this.buildingPathArray[random]

    // Get query params.
    this.companyRequirementID = this.activatedRoute.snapshot.queryParamMap.get("companyRequirementID")
    this.talentID = this.activatedRoute.snapshot.queryParamMap.get("talentID")
    this.email = this.activatedRoute.snapshot.queryParamMap.get("email")
    if (this.companyRequirementID) {
      this.getCompanyRequirementDetails()
    }
    if (this.talentID) {
      this.getWaitingList()
    }
    this.getSourceList()
  }

  //*********************************************FORMAT FUNCTIONS************************************************************

  // Format the package number in terms of LPA.
  formatPackageInLPA(min: number, max: number): string {
    if (!min && !max) {
      return
    }
    let minNumber: number = min / 100000
    let maxNumber: number = max / 100000
    let output: string = minNumber + " - " + maxNumber + " Lpa"
    this.packageInString = output
  }

  // Format the experience.
  formatExperience(minExp: number, maxExp: number): void {
    if (!minExp && !maxExp) {
      this.companyRequirementDetails.experienced = "Fresher"
      return
    }
    if (!maxExp) {
      this.companyRequirementDetails.experienced = minExp + " year(s)"
      return
    }
    this.companyRequirementDetails.experienced = minExp + " - " + maxExp + " year(s)"
  }

  // Format fields of company requirement through interface.
  formatCompanyRequirementsFields(): void {
    // Format experience.
    this.formatExperience(this.companyRequirementDetails.minimumExperience, this.companyRequirementDetails.maximumExperience)

    // Format package offered in indian rupee system.
    this.formatPackageInLPA(this.companyRequirementDetails.minimumPackage, this.companyRequirementDetails.maximumPackage)
  }

  //*********************************************APPLY FUNCTIONS************************************************************

  // On clicking apply now button.
  onApplyNowButtonClick(): void {
    this.openModal(this.applyNowModal, 'md').result.then(() => {
      this.applyNow()
    }, (err) => {
      console.error(err)
      return
    })
  }

  // Apply for the requirement by adding entry into waiting list.
  applyNow(): void {
    this.spinnerService.loadingMessage = "Applying..."


    let currentSourceID: string
    for (let i = 0; i < this.sourceList.length; i++) {
      if (this.sourceList[i].name == "wl") {
        currentSourceID = this.sourceList[i].id
      }
    }
    if (!currentSourceID) {
      alert("Application unsuccessful, please try again later.")


      return
    }
    let waitingList: IWaitingList = {
      talentID: this.talentID,
      isActive: true,
      email: this.email,
      companyBranchID: this.companyRequirementDetails.companyBranch.id,
      companyRequirementID: this.companyRequirementID,
      sourceID: currentSourceID
    }
    this.talentService.addWaitingList(waitingList).subscribe((response) => {
      alert("Application successful")
      this.isApplied = true
    }, error => {
      alert("Application unsuccessful, please try again later.")
      console.error(error)
    })
  }

  //*********************************************OTHER FUNCTIONS************************************************************
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

  // Nagivate to previous url.
  backToPreviousPage(): void {
    this._location.back()
  }

  //*********************************************GET FUNCTIONS************************************************************
  // Get waiting list by the talent id and requirement id.
  getWaitingList(): void {
    this.spinnerService.loadingMessage = "Getting Company Details"


    let queryParams: any = {
      "companyRequirementID": this.companyRequirementID
    }
    this.talentService.getWaitingListByTalent(this.talentID, queryParams).subscribe((response) => {
      this.waitingList = response
      if (this.waitingList.length > 0) {
        this.isApplied = true
      }
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get one company requirement by id.
  getCompanyRequirementDetails(): void {
    this.spinnerService.loadingMessage = "Getting Company Details"


    this.companyRequirementService.getCompanyDetails(this.companyRequirementID).subscribe((response) => {
      this.companyRequirementDetails = response
      this.companyBranch = this.companyRequirementDetails.companyBranch
      this.company = this.companyRequirementDetails.companyBranch.company
      this.formatCompanyRequirementsFields()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get source list.
  getSourceList(): void {
    this.generalService.getSources().subscribe(response => {
      this.sourceList = response
    }, err => {
      console.error(err)
    })
  }
}


