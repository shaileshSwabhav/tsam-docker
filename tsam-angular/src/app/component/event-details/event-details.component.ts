import { Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { DatePipe, Location } from '@angular/common';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { FormGroup, FormControl, Validators, FormBuilder } from '@angular/forms';
import { AdminService, IEvent } from 'src/app/service/admin/admin.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { ITalentEventRegistration, TalentService } from 'src/app/service/talent/talent.service';

@Component({
  selector: 'app-event-details',
  templateUrl: './event-details.component.html',
  styleUrls: ['./event-details.component.css']
})
export class EventDetailsComponent implements OnInit {

  // Spinner.



  // talent registration details.
  registrationDetails: ITalentEventRegistration

  // Event.
  eventID: string
  event: IEvent

  // Register.
  registerForm: FormGroup

  // credentials
  loginID: string

  // Modal.
  modalRef: any
  @ViewChild('registerModal') registerModal: any
  @ViewChild('successModal') successModal: any

  // Related course.
  relatedCourseList: IRelatedCourse[]

  constructor(
    private formBuilder: FormBuilder,
    private activatedRoute: ActivatedRoute,
    private _location: Location,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private localService: LocalService,
    private talentService: TalentService,
    private adminService: AdminService,
    private datePipe: DatePipe
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  // Initialize global variables.
  initializeVariables(): void {

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Event Details"


    this.loginID = this.localService.getJsonValue("loginID")

    this.extractEventIDFromURL()
  }

  getAllComponents(): void {
    this.getRegistrationDetails()
    this.getEvent()
    this.getRelatedCourse()
  }

  extractEventIDFromURL(): void {
    this.eventID = this.activatedRoute.snapshot.queryParamMap.get("eventID")

  }

  // Create register form.
  createRegisterForm(): void {

    this.registerForm = this.formBuilder.group({
      id: new FormControl(null),
      talentID: new FormControl(this.loginID),
      eventID: new FormControl(this.eventID),
      registrationDate: new FormControl(this.datePipe.transform(new Date(), 'yyyy-MM-dd')),
      hasAttended: new FormControl(false),
    })

    // this.registerForm = this.formBuilder.group({
    //   id: new FormControl(null),
    //   code: new FormControl({ value: null, disabled: true }),
    //   name: new FormControl(null, [Validators.pattern(/^[a-zA-Z ]*$/), Validators.maxLength(50)]),
    //   email: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/),
    //     Validators.maxLength(100)]),
    //   contact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/)]),
    //   isExperienced: new FormControl(null),
    //   college: new FormControl(null, [Validators.maxLength(50)]),
    //   company: new FormControl(null, [Validators.maxLength(50)]),
    //   qualification: new FormControl(null),
    //   city: new FormControl(null, [Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/), Validators.maxLength(50)]),
    //   courseRating: new FormControl(null, [Validators.min(1), Validators.max(10)]),
    //   technologies: new FormControl(null),
    // })
  }

  // Nagivate to previous url.
  backToPreviousPage(): void {
    this._location.back()
  }

  // Used to open modal.
  openModal(content: any, size?: string): NgbModalRef {

    if (!size) {
      size = "xl"
    }

    let options: NgbModalOptions = {
      ariaLabelledBy: 'modal-basic-title', keyboard: false,
      backdrop: 'static', size: size, centered: true
    }

    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  onSuccessModalClose(modal: NgbModalRef): void {
    modal.dismiss()
    this._location.back()
  }

  // On clicking register button.
  onRegisterNowButtonClick(): void {
    if (confirm("Click on yes to confirm your registration.")) {
      this.createRegisterForm()
      this.registerTalent()
    }
  }

  getEvent(): void {
    this.spinnerService.loadingMessage = "Getting event"

    this.adminService.getEvent(this.eventID).subscribe((response: any) => {
      this.event = response.body
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  registerTalent(): void {
    this.spinnerService.loadingMessage = "Registring for event"

    this.talentService.addTalentRegistration(this.registerForm.value).subscribe((response: any) => {
      // alert("Successfully registered for the event.")
      this.openModal(this.successModal, 'lg')
      // this.getRegistrationDetails()
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getRegistrationDetails(): void {
    this.spinnerService.loadingMessage = "Getting events"



    let queryParams: any = {
      talentID: this.loginID,
      eventID: this.eventID
    }

    this.talentService.getTalentRegistration(queryParams).subscribe((response: any) => {
      this.registrationDetails = response.body
      // console.log(this.registrationDetails);
    }, (err: any) => {
      console.error(err)
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  onMeetingLinkClick(): void {
    // run this commands only when talent has not attended the event.
    if (!this.registrationDetails.hasAttended) {
      this.createRegisterForm()
      this.registrationDetails.registrationDate = this.datePipe.transform(this.registrationDetails.registrationDate, 'yyyy-MM-dd')
      this.registrationDetails.hasAttended = true
      this.registerForm.patchValue(this.registrationDetails)

      // console.log(this.registerForm.value);
      this.updateTalentRegistration()
    }
  }

  updateTalentRegistration(): void {
    this.talentService.updateTalentRegistration(this.registerForm.value).subscribe((response: any) => {
      console.log(response);
    }, (err: any) => {
      console.error(err)
    })
  }

  getRelatedCourse(): void {
    // Related course.
    this.relatedCourseList = [
      {
        id: "1",
        name: "Python",
        image: "assets/images/python.png"
      },
      {
        id: "2",
        name: "c++",
        image: "assets/images/c++.png"
      }
    ]
  }

  // On clicking register button.
  register(): void {
    this.modalRef.close()
    this.openModal(this.successModal, 'lg')
  }

}

export interface IRelatedCourse {
  id?: string
  name: string
  image?: string
}
