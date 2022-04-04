import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { CollegeCampusService } from 'src/app/service/college/college-campus/college-campus.service';
import { DegreeService } from 'src/app/service/degree/degree.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-campus-registration-form',
  templateUrl: './campus-registration-form.component.html',
  styleUrls: ['./campus-registration-form.component.css']
})
export class CampusRegistrationFormComponent implements OnInit {

  // Components.
  stateList: any[]
  countryList: any[]
  degreeList: any[]
  specializationList: any[]
  academicYearList: any[]
  selectedCampusDriveCollegeList: any[]
  candidateResultList: any[]

  // Flags.
  showRegisterCandidateForm: boolean
  showCampusTalentRegistrationFields: boolean
  showCandidatesList: boolean
  isLinkInvalid: boolean
  isSpecializationLoading: boolean
  showForm: boolean
  showRedirectMessage: boolean
  showInvalidLinkMeassage: boolean

  // Candidate.
  candidateList: any[]
  selectedCandidatesList: string[]
  candidateForm: FormGroup

  // Pagination.
  limitForCandidateList: number
  offsetForCandidateList: number
  currentPageForCandidateList: number
  totalCandidates: number
  paginationStringForCandidate: string

  // Modal.
  @ViewChild('candidateFormModal') candidateFormModal: any

  // Search.
  showSearchForCandidate: boolean
  isSearchedForCandidate: boolean
  candidateSearchForm: FormGroup
  searchFormValueForCandidate: any

  // Resume.
  isResumeUploadedToServer: boolean
  isFileUploading: boolean
  docStatus: string
  displayedFileName: string

  // Spinner.


  // Candidate related variables.
  currentYear: number
  selectedCampusdrive: any
  campusDriveCodeParam: string

  constructor(private formBuilder: FormBuilder,
    public utilityService: UtilityService,
    private spinnerService: SpinnerService,
    private generalService: GeneralService,
    private degreeService: DegreeService,
    private collegeCampusService: CollegeCampusService,
    private fileOperationService: FileOperationService,
    private router: Router,
    private activatedRoute: ActivatedRoute
  ) {
    this.initializeVariables()
    if (!this.isLinkInvalid) {
      this.getAllComponents()
    }
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  initializeVariables(): void {
    // Components.
    this.stateList = []
    this.countryList = []
    this.degreeList = []
    this.specializationList = []
    this.academicYearList = []
    this.selectedCampusDriveCollegeList = []
    this.candidateResultList = []

    // Flags.
    this.showRegisterCandidateForm = false
    this.showCampusTalentRegistrationFields = false
    this.showCandidatesList = false
    this.isLinkInvalid = false
    this.isSpecializationLoading = false
    this.showForm = true
    this.showRedirectMessage = false
    this.showInvalidLinkMeassage = false

    // Candidate.
    this.selectedCandidatesList = []

    // Pagination.
    this.limitForCandidateList = 5
    this.offsetForCandidateList = 0
    this.currentPageForCandidateList = 0

    // Search.
    this.showSearchForCandidate = false
    this.isSearchedForCandidate = false
    this.searchFormValueForCandidate = {}

    // Resume.
    this.docStatus = ""
    this.displayedFileName = "Select file"
    this.isResumeUploadedToServer = false
    this.isFileUploading = false

    // Candidate related variables.
    this.currentYear = new Date().getFullYear()

    // Get source, enquiry for and course from url params.
    let queryParams = this.activatedRoute.snapshot.queryParams

    // Campus drive code.
    this.campusDriveCodeParam = queryParams['code']
    if (!this.campusDriveCodeParam) {
      this.isLinkInvalid = true
      this.showForm = false
      this.showInvalidLinkMeassage = true
      alert("Invalid Link!!!")
      return
    }

    //****************************INITIALIZE FORMS*************************************** */
    this.createCandidateForm()
  }

  //*********************************************CREATE FORMS************************************************************

  // Create new candidate form.
  createCandidateForm(): void {
    this.candidateForm = this.formBuilder.group({
      firstName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
      lastName: new FormControl(null, [Validators.required, Validators.pattern(/^[a-zA-Z]*$/), Validators.maxLength(50)]),
      email: new FormControl(null, [Validators.pattern(/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$/), Validators.required,
      Validators.maxLength(100)]),
      contact: new FormControl(null, [Validators.pattern(/^[6789]\d{9}$/), Validators.required]),
      address: new FormControl(null, [Validators.pattern(/^[.0-9a-zA-Z\s,-\/]+$/)]),
      city: new FormControl(null, [Validators.pattern(/^[a-zA-Z]+([a-zA-Z ]?)+$/)]),
      state: new FormControl(null),
      country: new FormControl(null),
      pinCode: new FormControl(null, [Validators.pattern(/^[1-9][0-9]{5}$/)]),
      academicYear: new FormControl(null, [Validators.required]),
      isSwabhavTalent: new FormControl(null),
      college: new FormControl(null, [Validators.required, Validators.maxLength(200)]),
      percentage: new FormControl(null, [Validators.min(1), Validators.pattern(/(^100([.]0{1,2})?)$|(^\d{1,2}([.]\d{1,2})?)$/i), Validators.required]),
      passout: new FormControl(null, [Validators.min(1980), Validators.max(this.currentYear + 3), Validators.required]),
      degree: new FormControl(null, Validators.required),
      specialization: new FormControl(null, Validators.required),
      talentID: new FormControl(null),
      campusDriveID: new FormControl(null),
      campusTalentRegistrationID: new FormControl(null),
      talentAcademicID: new FormControl(null),
      resume: new FormControl(null),
      isTestLinkSent: new FormControl(false, [Validators.required]),
      hasAttempted: new FormControl(false, [Validators.required])
    })
    this.candidateForm.get('state').disable()
    this.candidateForm.get('specialization').disable()
  }

  //*********************************************ADD FUNCTIONS FOR CANDIDATE FORM************************************************************
  // On clicking add candidate button on candidate form.
  addCandidate(): void {
    this.spinnerService.loadingMessage = "Adding Candidate"

    let candidate: any = this.candidateForm.value
    this.patchIDFromObjectsForCandidate(candidate)
    this.collegeCampusService.addCandidate(this.candidateForm.value, this.selectedCampusdrive.id).subscribe((response) => {
      this.candidateForm.reset()
      this.showForm = false
      this.showRedirectMessage = true
    }, (error) => {
      console.error(error)
      if (typeof error.error == 'object' && error) {
        alert(this.utilityService.getErrorString(error))
        return
      }
      if (error.error == undefined) {
        alert('Campus Talent could not be registered, try again')
        return
      }
      alert(error.statusText)
    })
  }

  //*********************************************OTHER FUNCTIONS FOR CANDIDATE************************************************************
  // On clicking sumbit button in candidate form.
  onCandidateFormSubmit(): void {
    if (this.candidateForm.invalid) {
      this.candidateForm.markAllAsTouched()
      console.log(this.findInvalidControls())
      return
    }
    this.addCandidate()
  }

  // On uplaoding resume.
  onResourceSelect(event: any): void {
    this.docStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]

      // Upload resume if it is present.]
      this.isFileUploading = true
      this.fileOperationService.uploadResume(file).subscribe((data: any) => {
        this.candidateForm.markAsDirty()
        this.candidateForm.patchValue({
          resume: data
        })
        this.displayedFileName = file.name
        this.isFileUploading = false
        this.isResumeUploadedToServer = true
        this.docStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
      }, (error) => {
        this.isFileUploading = false
        this.docStatus = `<p><span>&#10060;</span> ${error}</p>`
      })
    }
  }

  // Will be called on [addTag]= "addCollegeToList" args will be passed automatically.
  addCollegeToList(option: any): Promise<any> {
    return new Promise((resolve) => {
      resolve(option)
    })
  }

  // Extract ID from objects in candidate form.
  patchIDFromObjectsForCandidate(candidate: any): void {
    if (this.candidateForm.get('degree').value) {
      candidate.degreeID = this.candidateForm.get('degree').value.id
      delete candidate['degree']
    }
    if (this.candidateForm.get('specialization').value) {
      candidate.specializationID = this.candidateForm.get('specialization').value.id
      delete candidate['specialization']
    }
  }

  //Compare for select option field.
  compareFn(optionOne: any, optionTwo: any): boolean {
    if (optionOne == null && optionTwo == null) {
      return true
    }
    if (optionTwo != undefined && optionOne != undefined) {
      return optionOne.id === optionTwo.id
    }
    return false
  }

  // Get all invalid controls in talent form.
  public findInvalidControls(): any[] {
    const invalid = []
    const controls = this.candidateForm.controls
    for (const name in controls) {
      if (controls[name].invalid) {
        invalid.push(name)
      }
    }
    return invalid
  }

  redirectToHomePage(): void {
    // this.elementRef.nativeElement.ownerDocument.body.style.backgroundColor = 'white'
    window.location.href = 'https://swabhavtechlabs.com/'
  }

  //*********************************************GET FUNCTIONS************************************************************

  // Get all lists.
  getAllComponents(): void {
    this.getCountryList()
    this.getQualificationList()
    this.getAcademicYear()
    this.getCampusDriveByCode()
  }

  getCampusDriveByCode(): void {
    this.collegeCampusService.getCampusDriveByCode(this.campusDriveCodeParam).subscribe((response: any[]) => {
      this.selectedCampusdrive = response
      this.getSelectedCampusDriveColleegList()
    }, (err) => {
      console.error(err.error)
    })
  }

  // Get country list.
  getCountryList(): void {
    this.generalService.getCountries().subscribe((response: any[]) => {
      this.countryList = response
    }, (err) => {
      console.error(err.error)
    })
  }

  // Get state list by country.
  getStateListByCountry(country: any): void {
    if (country == null) {
      this.stateList = []
      this.candidateForm.get('state').setValue(null)
      this.candidateForm.get('state').disable()
      return
    }
    if (country.name == "U.S.A." || country.name == "India") {
      this.candidateForm.get('state').enable()
      this.generalService.getStatesByCountryID(country.id).subscribe((response: any[]) => {
        this.stateList = response
      }, (err) => {
        console.error(err)
      })
    }
    else {
      this.stateList = []
      this.candidateForm.get('state').setValue(null)
      this.candidateForm.get('state').disable()
    }
  }

  // Get selected campus drive's college list.
  getSelectedCampusDriveColleegList(): void {
    this.selectedCampusDriveCollegeList = this.selectedCampusdrive.collegeBranches
  }

  // Get degree list.
  getQualificationList(): void {
    let queryParams: any = {
      limit: -1,
      offset: 0,
    }
    this.degreeService.getAllDegrees(queryParams).subscribe((response: any) => {
      this.degreeList = response.body
    }, (err) => {
      console.error(this.utilityService.getErrorString(err))
    })
  }

  getSpecializationListByDegree(degree: any): void {
    this.candidateForm.get('specialization').setValue(null)
    this.specializationList = []
    if (!degree) {
      this.candidateForm.get('specialization').setValue(null)
      this.specializationList = []
      this.candidateForm.get('specialization').disable()
      return
    }
    this.isSpecializationLoading = true
    this.generalService.getSpecializationByDegreeID(degree.id).subscribe((response: any) => {
      this.specializationList = response.body
      this.candidateForm.get('specialization').enable()
      this.isSpecializationLoading = false
    }, (err) => {
      console.error(err)
      this.isSpecializationLoading = false
    })
  }

  // Get academic year list.
  getAcademicYear(): void {
    this.generalService.getGeneralTypeByType("academic_year").subscribe((response: any[]) => {
      this.academicYearList = response
    }, (err) => {
      console.error(this.utilityService.getErrorString(err))
    })
  }
}
