import { Component, EventEmitter, OnInit, Output, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { NgbModalRef, NgbModalOptions, NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Role, UrlConstant } from 'src/app/service/constant';
import { CourseModuleService } from 'src/app/service/course-module/course-module.service';
import { CourseService, ICourse } from 'src/app/service/course/course.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { ITechnology, TechnologyService } from 'src/app/service/technology/technology.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
  selector: 'app-course-details',
  templateUrl: './course-details.component.html',
  styleUrls: ['./course-details.component.css']
})
export class CourseDetailsComponent implements OnInit {

  // Component.
  technologyList: ITechnology[]
  studentRatinglist: any[]
  academicYearList: any[]
  courseTypeList: any[]
  courseLevelList: any[]

  // Course.
  // minor change remove after pull
  courseID: any
  course: any
  courseDetailsLeftList: any[]
  courseDetailsRightList: any[]
  courseHeaderList: any[]
  courseModuleTabList: any[]
  courseForm: FormGroup
  technologyPreRequisite: string

  // Brochure.
  isBrochureUploadedToServer: boolean
  isBrochureUploading: boolean
  brochureDocStatus: string
  brochureDisplayedFileName: string

  // Logo.
  isLogoUploadedToServer: boolean
  isLogoUploading: boolean
  logoDocStatus: string
  logoDisplayedFileName: string

  // Modal.
  modalRef: any
  @ViewChild('courseUpdateModal') courseUpdateModal: any

  // Course Module.
  courseModuleList: any[]
  selectedCourseModule: any

  // Spinner.



  // Access.
  permission: IPermission

  // Flags.
  doesEligibilityExists: boolean
  isTechLoading: boolean
  isFaculty: boolean

  // Number.
  technologyLimit: number
  technologyOffset: number

  // cke editor.
  @ViewChild("ckeEditorPreRequisites") ckeditorPreRequisites: any
  ckeEditorConfig: any

  //emitter to Parent
  @Output() showCourseModule = new EventEmitter()
  constructor(
    private route: ActivatedRoute,
    private spinnerService: SpinnerService,
    private courseService: CourseService,
    private courseModuleService: CourseModuleService,
    private utilService: UtilityService,
    private urlConstant: UrlConstant,
    private modalService: NgbModal,
    private formBuilder: FormBuilder,
    private techService: TechnologyService,
    private generalService: GeneralService,
    private fileOperationService: FileOperationService,
    private localService: LocalService,
    private role: Role,
  ) {
    this.initializeVariables()
    this.getAllComponents()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void { }

  // Initialize global variables.
  initializeVariables(): void {

    // Component.
    this.technologyList = []
    this.studentRatinglist = []
    this.academicYearList = []
    this.courseTypeList = []
    this.courseLevelList = []

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Dashboard Ready..."


    // Access.
    this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
    if (this.isFaculty){
			this.permission = this.utilService.getPermission(this.urlConstant.BANK_COURSE)
		}
		if (!this.isFaculty){
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_COURSE_MASTER)
		}

    // Course.
    this.courseDetailsLeftList = []
    this.courseDetailsRightList = []
    this.courseModuleTabList = []
    this.courseHeaderList = []

    // Brochure.
    this.brochureDocStatus = ""
    this.brochureDisplayedFileName = "Select file"
    this.isBrochureUploadedToServer = false
    this.isBrochureUploading = false

    // Logo.
    this.logoDocStatus = ""
    this.logoDisplayedFileName = "Select file"
    this.isLogoUploadedToServer = false
    this.isLogoUploading = false

    // Flags.
    this.doesEligibilityExists = false
    this.isTechLoading = false

    // Number.
    this.technologyLimit = 10
    this.technologyOffset = 0

    // Cke editor configuration.
    this.ckeEditorCongiguration()

    // Course Module.
    this.courseModuleList = []
    this.selectedCourseModule = []
    this.getQueryParm()
  }

  //********************************************* CREATE FORM FUNCTIONS ************************************************************

  // Create course form.
  createCourseForm(): void {
    this.courseForm = this.formBuilder.group({
      id: new FormControl(),
      code: new FormControl(null),
      name: new FormControl(null, [Validators.required]),
      technologies: [Array(), Validators.required],
      courseType: new FormControl(null, [Validators.required]),
      courseLevel: new FormControl(null, [Validators.required]),
      preRequisites: new FormControl(null, [Validators.maxLength(2000)]),
      description: new FormControl(null, [Validators.maxLength(2000)]),
      eligibility: this.formBuilder.group({
        id: new FormControl(),
        technologies: [null],
        studentRating: new FormControl(null),
        experience: new FormControl(null),
        academicYear: new FormControl(null)
      }),
      price: new FormControl(null, [Validators.required]),
      durationInMonths: new FormControl(null, [Validators.required]),
      totalHours: new FormControl(null, [Validators.required]),
      totalSessions: new FormControl(null),
      brochure: new FormControl(null),
      logo: new FormControl(null),
    })
  }

  //********************************************* FORMAT FUNCTIONS ************************************************************

  // Format the course details list.
  formatCourseDetailsList(): void {
    // console.log(this.course.eligibility)
    if (this.course.eligibility && this.course.eligibility?.technologies) {
      let tempTechnologyList: string[] = []
      for (let i = 0; i < this.course.eligibility.technologies.length; i++) {
        tempTechnologyList.push(this.course.eligibility.technologies[i].language)
      }
      this.technologyPreRequisite = tempTechnologyList.join(", ")
      // console.log(this.technologyPreRequisite)
    }
    else {
      this.technologyPreRequisite = null
    }
    this.courseDetailsLeftList = []
    this.courseDetailsRightList = []
    this.courseDetailsLeftList = [
      {
        fieldName: "Name",
        fieldValue: this.course.name
      },
      {
        fieldName: "Code",
        fieldValue: this.course.code
      },
      {
        fieldName: "Type",
        fieldValue: this.course.courseType
      }
    ]

    if (!this.isFaculty) {
      this.courseDetailsLeftList.push({
        fieldName: "Price",
        fieldValue: "INR " + this.course.price
      }
      )
    }

    if (this.course.description != null) {
      this.courseDetailsRightList.push({
        fieldName: "Course Overview",
        fieldValue: this.course.description
      }
      )
    }

    this.courseDetailsRightList.push(
      {
        fieldName: "Duration (in Hours)",
        fieldValue: this.course.totalHours
      }
    )
  }

  // Format the course modules list.
  formatCourseModuleList(): void {

    // Create course module tab list.
    this.courseModuleTabList = []
    for (let i = 0; i < this.courseModuleList.length; i++) {
      this.courseModuleTabList.push(
        {
          moduleName: this.courseModuleList[i].module.moduleName,
          module: this.courseModuleList[i].module,
        }
      )

      if (this.courseModuleList[i].module.logo) {
        this.courseModuleTabList[i].imageURL = this.courseModuleList[i].module.logo
      }
      else {
        this.courseModuleTabList[i].imageURL = "assets/icon/grey-icons/Score.png"
      }
    }

    if (this.courseModuleList.length > 0) {
      this.selectedCourseModule = this.courseModuleList[0].module
    }

    this.formatCourseModuleHeaderList()
  }

  // Format the course module header list.
  formatCourseModuleHeaderList(): void {

    let conceptCount: number = 0
    let assignmentCount: number = 0
    // Count the number of concepts in all sub topics.
    console.log("selectedCourseModule", this.selectedCourseModule?.moduleTopics)
    var map = new Map();
    for (let i = 0; i < this.selectedCourseModule?.moduleTopics?.length; i++) {
      let moduleTopic: any = this.selectedCourseModule?.moduleTopics[i]
      if (moduleTopic?.topicProgrammingQuestions) {
        assignmentCount += moduleTopic?.topicProgrammingQuestions?.length
      }
      if (moduleTopic.topicProgrammingConcept?.length > 0) {
        for (let j = 0; j < moduleTopic.topicProgrammingConcept.length; j++) {
          if (map.has(moduleTopic.topicProgrammingConcept[j].programmingConceptID)) {
            continue
          } else {
            map.set(moduleTopic.topicProgrammingConcept[j].programmingConceptID, 1)
            conceptCount++
          }
        }
      }
    }

    // Create coures header list.
    this.courseHeaderList = [
      {
        image: "assets/course/Modules.png",
        header: "Modules",
        count: this.courseModuleList.length
      },
      {
        image: "assets/course/topics.png",
        header: "Topics",
        count: this.selectedCourseModule.moduleTopics?.length
      },
      {
        image: "assets/course/subtopics.png",
        header: "Concepts",
        count: conceptCount
      },
      {
        image: "assets/course/assignments.png",
        header: "Assignments",
        count: assignmentCount
      },
      // {
      //   image: "assets/course/assignments.png",
      //   header: "Projects",
      //   count: 10
      // }
    ]
  }

  //********************************************* COURSE MODULE FUNCTIONS ************************************************************

  // On clicking module tab.
  onModuleTabClick(event: any): void {
    for (let i = 0; i < this.courseModuleList.length; i++) {
      if (i == event.index) {
        this.selectedCourseModule = this.courseModuleList[i].module
        this.formatCourseModuleHeaderList()
      }
    }
    // this.selectedCourseModule = courseModule
    // this.setIsClickedModuleTopic()
    // this.formatCourseModuleHeaderList()
  }

  // Set isTopicClicked and isSubtopicClicked as false initially.
  setIsClickedModuleTopic(): void {
    for (let i = 0; i < this.selectedCourseModule.moduleTopics.length; i++) {
      this.selectedCourseModule.moduleTopics[i].isTopicClicked = false
      if (this.selectedCourseModule.moduleTopics[i].subTopics?.length > 0) {
        for (let j = 0; j < this.selectedCourseModule.moduleTopics[i].subTopics.length; j++) {
          this.selectedCourseModule.moduleTopics[i].subTopics.isSubTopicClicked = false
        }
      }
    }
  }

  //********************************************* COURSE FUNCTIONS ************************************************************

  // Emit to parent component.
  changeToCourseModule(): void {
    // this.spinnerService.loadingMessage = "Adding course module"
    this.showCourseModule.emit()
  }

  // On clicking course update button.
  onCourseUpdateButtonClick() {
    this.brochureDocStatus = ""
    this.logoDocStatus = ""
    this.isBrochureUploadedToServer = false
    this.isLogoUploadedToServer = false
    this.updateCourseForm(this.course)
    this.openModal(this.courseUpdateModal)

    // Brochure.
    this.brochureDisplayedFileName = "Select file"
    if (this.course.brochure) {
      this.brochureDisplayedFileName = `<a class='custom-file-label' href=${this.course.brochure} target="_blank">Select file</a>`
    }

    // Logo.
    this.logoDisplayedFileName = "Select file"
    if (this.course.logo) {
      this.logoDisplayedFileName = `<a class='custom-file-label' href=${this.course.logo} target="_blank">Select file</a>`
    }
  }

  // Update course form.
  updateCourseForm(course: any): void {
    this.createCourseForm()
    if (course.eligibility != null) {
      this.doesEligibilityExists = true
    } else {
      course.eligibility = {}
      this.doesEligibilityExists = false
    }
    this.courseForm.patchValue(course)

    // Make eligibility null if eligibility does not exist.
    if (Object.keys(course.eligibility).length == 0) {
      course.eligibility = null
    }
  }

  // Add or remove eligibility.
  eligilibilityChecked(event) {
    event.target.checked = true
    if (this.doesEligibilityExists) {
      if (this.courseForm.get('eligibility').valid) {
        if (confirm("This action will delete the batch eligibility. Are you sure you want to delete batch eligibility?")) {
          this.courseForm.get('eligibility.id').setValue(null)
          event.target.checked = false
        } else {
          return
        }
      }
      this.courseForm.get('eligibility.id').setValue(null)
      this.doesEligibilityExists = false
      this.setEligibility()
      return
    }
    this.doesEligibilityExists = true
    this.setEligibility()
  }


  // Set the eligibility.
  setEligibility() {
    this.technologyCheck()
    this.studentRatingCheck()
    this.experienceCheck()
    this.academicYearCheck()
    this.courseForm.markAsDirty()
  }

  // Check technology eligibility.
  technologyCheck() {
    const technologyControl = this.courseForm.get('eligibility.technologies')
    if (this.doesEligibilityExists) {
      technologyControl.setValidators([Validators.required])
    } else {
      technologyControl.setValue(null)
      technologyControl.setValidators(null)
      technologyControl.markAsUntouched()
    }
    technologyControl.updateValueAndValidity()
  }

  // Check student rating eligibility.
  studentRatingCheck() {
    const studentRatingControl = this.courseForm.get('eligibility.studentRating')
    if (this.doesEligibilityExists) {
      studentRatingControl.setValidators([Validators.required])
    } else {
      studentRatingControl.setValue(null)
      studentRatingControl.setValidators(null)
      studentRatingControl.markAsUntouched()
    }
    studentRatingControl.updateValueAndValidity()
  }

  // Check experience eligibility.
  experienceCheck() {
    const experienceControl = this.courseForm.get('eligibility.experience')
    if (this.doesEligibilityExists) {
      experienceControl.setValidators([Validators.required])
    } else {
      experienceControl.setValue(null)
      experienceControl.setValidators(null)
      experienceControl.markAsUntouched()
    }
    experienceControl.updateValueAndValidity()
  }

  // Check academic year eligibility.
  academicYearCheck() {
    const academicYearControl = this.courseForm.get('eligibility.academicYear')
    if (this.doesEligibilityExists) {
      academicYearControl.setValidators([Validators.required])
    } else {
      academicYearControl.setValue(null)
      academicYearControl.setValidators(null)
      academicYearControl.markAsUntouched()
    }
    academicYearControl.updateValueAndValidity()
  }

  // On uplaoding brochure.
  onBrochureSelect(event: any): void {
    this.brochureDocStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      let err = this.fileOperationService.isDocumentFileValid(file)
      if (err != null) {
        this.brochureDocStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // Upload brochure if it is present.
      this.isBrochureUploading = true
      this.fileOperationService.uploadBrochure(file, this.fileOperationService.COURSE_FOLDER + this.fileOperationService.BROCHURE_FOLDER)
        .subscribe((data: any) => {
          this.courseForm.markAsDirty()
          this.courseForm.patchValue({
            brochure: data
          })
          this.brochureDisplayedFileName = file.name
          this.isBrochureUploading = false
          this.isBrochureUploadedToServer = true
          this.brochureDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
        }, (error) => {
          this.isBrochureUploading = false
          this.brochureDocStatus = `<p><span>&#10060;</span> ${error}</p>`
        })
    }
  }

  // On uplaoding logo.
  onLogoSelect(event: any) {
    this.logoDocStatus = ""
    let files = event.target.files
    if (files && files.length) {
      let file = files[0]
      let err = this.fileOperationService.isImageFileValid(file)
      if (err != null) {
        this.logoDocStatus = `<p><span>&#10060;</span> ${err}</p>`
        return
      }
      // Upload brochure if it is present.
      this.isLogoUploading = true
      this.fileOperationService.uploadLogo(file,
        this.fileOperationService.BATCH_FOLDER + this.fileOperationService.LOGO_FOLDER).subscribe((data: any) => {
          this.courseForm.markAsDirty()
          this.courseForm.patchValue({
            logo: data
          })
          this.logoDisplayedFileName = file.name
          this.isLogoUploading = false
          this.isLogoUploadedToServer = true
          this.logoDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
        }, (error) => {
          this.isLogoUploading = false
          this.logoDocStatus = `<p><span>&#10060;</span> ${error}</p>`
        })
    }
  }

  // Delete uploaded filr.
  deleteUploadedFile(): void {
    this.fileOperationService.deleteUploadedFile().subscribe((data: any) => {
    }, (error) => {
      console.error(error)
    })
  }

  // Validate course form.
  validateCourseForm(): void {
    if (this.courseForm.invalid) {
      this.courseForm.markAllAsTouched()
      return
    }

    this.updateCourse()
  }

  // Update course.
  updateCourse(): void {
    this.spinnerService.loadingMessage = "Updating course"

    let data = this.courseForm.value
    this.utilService.deleteNullValueIDFromObject(data)
    this.utilService.deleteNullValuePropertyFromObject(data.eligibility)
    if (data.eligibility) {
      this.deleteEligibilityIfEmpty(data)
    }
    this.courseService.updateCourse(data).subscribe(respond => {
      this.modalRef.close()

      this.getCourseDetails()
      alert("Course successfully updated")
    }, (err) => {
      console.error(err)

      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err?.error?.error)
    })
  }

  // Delete eligibility from course if empty.
  deleteEligibilityIfEmpty(course: ICourse): void {
    if (Object.keys(course.eligibility).length === 0 && course.eligibility.constructor === Object) {
      delete course.eligibility
    }
  }

  //********************************************* OTHER FUNCTIONS ************************************************************

  // Values from query params.
  getQueryParm(): void {
    this.route.queryParamMap.subscribe(params => {
      this.courseID = params.get("courseID")
      if (this.courseID) {
        this.getCourseDetails()
        this.getCourseModuleList()
      }
    }, err => {
      console.error(err);
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

  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef): void {
    if (this.isBrochureUploading || this.isLogoUploading) {
      alert("Please wait till file is being uploaded")
      return
    }
    if (this.isBrochureUploadedToServer) {
      if (!confirm("Uploaded brochure will be deleted.\nAre you sure you want to close?")) {
        return
      }
      this.deleteUploadedFile()
    }
    if (this.isLogoUploadedToServer) {
      if (!confirm("Uploaded logo will be deleted.\nAre you sure you want to close?")) {
        return
      }
      this.deleteUploadedFile()
    }
    modal.dismiss()
    this.isBrochureUploadedToServer = false
    this.isLogoUploadedToServer = false
    this.brochureDisplayedFileName = "Select file"
    this.logoDisplayedFileName = "Select file"
    this.brochureDocStatus = ""
    this.logoDocStatus = ""
  }

  // Redirect to brochure in new tab.
  redirectToBrochureLink(): void {
    window.open(this.course.brochure, "_blank");
  }

  // cke editor congiguration.
  ckeEditorCongiguration(): void {
    this.ckeEditorConfig = {
      extraPlugins: 'codeTag',
      removePlugins: "exportpdf",
      // stylesSet: 'new_styles',
      toolbar: [
        { name: 'styles', items: ['Styles', 'Format'] },
        {
          name: 'basicstyles', groups: ['basicstyles', 'cleanup'],
          items: ['Bold', 'Italic', 'Underline', 'Strike', 'Subscript', 'Superscript', 'RemoveFormat', 'Code']
        },
        {
          name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'],
          items: ['NumberedList', 'BulletedList', '-', 'Outdent', 'Indent', '-', 'Blockquote']
        },
        { name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
      ],
      toolbarGroups: [
        { name: 'styles' },
        { name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
        { name: 'document', groups: ['mode', 'document', 'doctools'] },
        { name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
      ],
      removeButtons: "",
      language: 'en',
      resize_enabled: false,
      width: "100%", height: "150px",
      forcePasteAsPlainText: false,
    }
  }

  // Redirect to external link.
  redirectToExternalLink(url): void {
    window.open(url, "_blank")
  }

  //********************************************* GET FUNCTIONS ************************************************************

  // Get all components.
  getAllComponents(): void {
    this.getTechnologyList()
    this.getStudentRatingList()
    this.getAcademicYearList()
    this.getCourseTypeList()
    this.getCourseLevelList()
  }

  // Get course details.
  getCourseDetails(): void {
    this.spinnerService.loadingMessage = "Getting Course Details..."


    this.courseService.getCourse(this.courseID).subscribe((response) => {
      this.course = response
      this.formatCourseDetailsList()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get course module list.
  getCourseModuleList(): void {
    this.spinnerService.loadingMessage = "Getting Course Modules..."


    let queryParams: any = {
      limit: -1,
      offset: 0
    }
    this.courseModuleService.getCourseModules(this.courseID, queryParams).subscribe((response) => {
      this.courseModuleList = response.body
      this.formatCourseModuleList()
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get technology list.
  getTechnologyList(event?: any): void {
    let queryParams: any = {}
    if (event && event?.term != "") {
      queryParams.language = event.term
    }
    this.isTechLoading = true
    this.techService.getAllTechnologies(this.technologyLimit, this.technologyOffset, queryParams).subscribe((response) => {
      this.technologyList = []
      this.technologyList = this.technologyList.concat(response.body)
    }, (err) => {
      console.error(err)
    }).add(() => {
      this.isTechLoading = false
    })
  }

  // Get student rating list.
  getStudentRatingList(): void {
    this.generalService.getGeneralTypeByType("talent_rating").subscribe((respond: any[]) => {
      this.studentRatinglist = respond
    }, (err) => {
      console.error(err)
    })
  }

  // Get academic year list.
  getAcademicYearList(): void {
    this.generalService.getGeneralTypeByType("academic_year").subscribe((respond: any[]) => {
      this.academicYearList = respond
    }, (err) => {
      console.error(err)
    })
  }

  // Get course type list. 
  getCourseTypeList(): void {
    this.generalService.getGeneralTypeByType("course_type").subscribe((respond: any) => {
      this.courseTypeList = respond
    }, (err) => {
      console.error(err)
    })
  }

  // Get course level list. 
  getCourseLevelList(): void {
    this.generalService.getGeneralTypeByType("course_level").subscribe((respond: any) => {
      this.courseLevelList = respond
    }, (err) => {
      console.error(err)
    })
  }

}
