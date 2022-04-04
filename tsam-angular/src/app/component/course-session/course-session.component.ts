import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { NgbModal, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { Location, ViewportScroller } from '@angular/common';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { UploadConstant, UrlConstant } from 'src/app/service/constant';
import { CourseService, ICourseSession, IResource } from 'src/app/service/course/course.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { pairwise, startWith } from 'rxjs/operators';
import { HttpParams } from '@angular/common/http';
import { ITechnology, TechnologyService } from 'src/app/service/technology/technology.service';
import { ResourceService } from 'src/app/service/resource/resource.service';

@Component({
  selector: 'app-course-session',
  templateUrl: './course-session.component.html',
  // encapsulation: ViewEncapsulation.None,
  styleUrls: ['./course-session.component.css']
})
export class CourseSessionComponent implements OnInit {

  // component
  technologyList: ITechnology[]
  isTechLoading: boolean
  techLimit: number
  techOffset: number

  //   course
  courseID: string
  courseName: string

  // session
  session: ICourseSession
  viewSubSessionClicked: boolean
  addNewSession: boolean
  sessionForm: FormGroup;
  sessionList: any[]
  addSessionToList: any[]
  totalSessions: number
  addSessionClicked: boolean
  addSubSessionClicked: boolean
  editIndex: number
  existingSessions: number

  //   resource
  multipleResourcesForm: FormGroup
  fileTypeList: any[]
  resourceTypeList: any[]
  resource: IResource[]
  resourceList: IResource[]
  resourceSubTypeList: any[]

  //modal
  modalAction: () => void;
  resourceHandler: () => void;
  formHandler: (index?: number) => void;
  modalButton: string
  modalHeader: string;
  resourceModalHeader: string
  resourceButton: string
  modalRef: any;

  // flags
  isViewClicked: boolean
  isEditSessionClicked: boolean
  isUpdateClicked: boolean
  isUpdateResourceClick: boolean

  // access
  permission: IPermission

    //spinner
    ;
  isResourceLoading: boolean

  // file upload
  previousFileType: string

  // constants
  readonly SWABHAV_DOCS = "Swabhav_Docs"
  readonly DOCUMENT = "Document"
  readonly AUDIO = "Audio"
  readonly IMAGE = "Image"
  readonly VIDEO = "Video"
  readonly URL = "URL"
  readonly BOOK = "Book"
  readonly PREVIEW = "Preview"

  constructor(
    private formBuilder: FormBuilder,
    public utilService: UtilityService,
    private generalService: GeneralService,
    private techService: TechnologyService,
    private fileOperationService: FileOperationService,
    private courseService: CourseService,
    private resourceService: ResourceService,
    private spinnerService: SpinnerService,
    private modalService: NgbModal,
    private urlConstant: UrlConstant,
    private viewportScroller: ViewportScroller,
    private route: ActivatedRoute,
    private location: Location,
    private uploadConstant: UploadConstant
  ) {
    this.extractCourseDetails()
    this.intializeVariables()
    this.createForms()
    this.getAllComponents()
  }

  intializeVariables(): void {
    this.permission = this.utilService.getPermission(this.urlConstant.COURSE_SESSION)

    this.sessionList = []
    this.addSessionToList = []
    this.fileTypeList = []
    this.resourceList = []
    this.resourceSubTypeList = []

    this.totalSessions = 0
    this.existingSessions = 0

    this.modalButton = "Add Session"
    this.modalHeader = "Add"
    this.resourceModalHeader = "Add Resource"
    this.previousFileType = null

    this.isViewClicked = false
    this.addNewSession = false
    this.isUpdateClicked = false
    this.isEditSessionClicked = false
    this.addSessionClicked = false
    this.addSubSessionClicked = false
    this.viewSubSessionClicked = false
    this.isUpdateResourceClick = false
    this.isResourceLoading = false
    this.techLimit = 10
    this.techOffset = 0


    this.spinnerService.loadingMessage = "Getting courses"
  }

  extractCourseDetails(): void {
    this.route.queryParamMap.subscribe(params => {
      this.courseID = params.get('courseID')
      this.courseName = params.get('courseName')
    })
  }

  createForms(): void {
    this.createSessionForm()
  }

  getAllComponents(): void {
    this.getFileType()
    this.getResourceType()
    this.getResourceSubType()
    this.getSessionForCourse()
    this.getTechnologyList()
  }

  onClickScroll(elementID: string): void {
    // console.log(elementID)
    this.viewportScroller.scrollToAnchor("editSession")
  }

  backToPreviousPage(): void {
    this.location.back()
  }


  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {
  }

  // Create Session Form Object
  createSessionForm(): void {
    this.sessionForm = this.formBuilder.group({
      id: new FormControl(null),
      name: new FormControl(null, [Validators.required]),
      hours: new FormControl(null, [Validators.required]),
      order: new FormControl(null, [Validators.required]),
      studentOutput: new FormControl(null, [Validators.required]),
      subSessions: new FormArray([]),
      resource: new FormArray([]),
    })
  }

  // Return Session Form Object.
  // returnSessionForm(): FormGroup {
  //   return this.formBuilder.group({
  //     id: new FormControl(),
  //     name: new FormControl(null, [Validators.required]),
  //     hours: new FormControl(null, [Validators.required]),
  //     order: new FormControl(null, [Validators.required]),
  //     studentOutput: new FormControl(null, [Validators.required]),
  //     subSessions: new FormArray([]),
  //     resource: new FormArray([]),
  //   });
  // }

  get subSessionControlArray() {
    return this.sessionForm.get('subSessions') as FormArray
  }

  // sub-session form
  createSubSessionForm(): void {
    this.subSessionControlArray.push(this.formBuilder.group({
      id: new FormControl(null),
      name: new FormControl(null, [Validators.required]),
      hours: new FormControl(null, [Validators.required]),
      order: new FormControl(null),
      studentOutput: new FormControl(null, [Validators.required]),
      subSessions: new FormControl(null),
      resource: new FormArray([]),
    }))
  }

  createResourceForm() {
    this.multipleResourcesForm = this.formBuilder.group({
      resources: this.formBuilder.array([], [Validators.required])
    })
    this.addResourcesToForm()
  }

  get resourceFormArray() {
    return this.multipleResourcesForm.get('resources') as FormArray
  }

  addResourcesToForm() {
    this.resourceFormArray.push(this.formBuilder.group({
      id: new FormControl(null),
      resourceName: new FormControl(null, [Validators.maxLength(100)]),
      isExistingResource: new FormControl(true, [Validators.required]),
      resourceType: new FormControl(null, [Validators.required]),
      resourceSubType: new FormControl(null),
      fileType: new FormControl(null, [Validators.required]),
      resourceURL: new FormControl(null, [Validators.maxLength(255)]),
      description: new FormControl(null, [Validators.maxLength(250)]),
      technologyID: new FormControl(null),
      docStatus: new FormControl(""),
      previewImageDocStatus: new FormControl(""),
      isFileUploading: new FormControl(false),
      isPreviewFileUploading: new FormControl(false),
      displayedFileName: new FormControl("Select File"),
      isDocumentUploadedToServer: new FormControl(false),
      isPreviewUploadedToServer: new FormControl(false),
      allowedExtension: new FormControl(null),
      previewDisplayedFileName: new FormControl("Select File"),
      // session: new FormControl(this.session)

      // control for book.
      isBook: new FormControl(false, [Validators.required]),
      author: new FormControl(null),
      publication: new FormControl(null),
    }))
    this.emitPreviousFileType(this.resourceFormArray.controls.length - 1)
  }

  // emitPreviousFileType will fetch previously selected file type.
  emitPreviousFileType(index: number): void {
    // Fill buffer with initial value.  Will emit immediately.
    this.resourceFormArray.at(index).get('fileType')
      .valueChanges.pipe(
        startWith(null as string),
        pairwise()).subscribe(([prev, next]: [any, any]) => {
          // console.log("PREV2", prev)
          this.previousFileType = prev
        })
  }

  addResourceClick(session: ICourseSession, modalContent: any): void {
    // console.log(session);
    this.isViewClicked = false
    this.isUpdateResourceClick = false
    this.session = session

    this.createResourceForm()
    this.modalAction = this.addResource
    this.resourceModalHeader = "Add Resource"
    this.resourceButton = "Add Resource"

    this.openModal(modalContent)
  }

  onViewResourceClick(session: ICourseSession, modalContent: any): void {
    this.session = session
    this.isViewClicked = true
    this.isUpdateResourceClick = false
    // this.getResourceForSession()
    // console.log(session);
    this.resource = session.resource

    this.createResourceForm()
    this.modalAction = this.addResource
    this.resourceModalHeader = "Resource Details"
    this.resourceButton = ""

    this.updateResourceForm()
    this.multipleResourcesForm.disable()
    this.openModal(modalContent)
  }

  updateResourceClick(): void {
    this.isViewClicked = false
    this.isUpdateResourceClick = true
    this.resourceModalHeader = "Update"
    this.resourceButton = "Update Resource"
    this.modalAction = this.updateResource
    this.multipleResourcesForm.enable()
  }

  updateResourceForm() {
    for (let i = 0; i < this.resource.length - 1; i++) {
      this.addResourcesToForm()
    }
    this.multipleResourcesForm.get('resources').patchValue(this.resource)
  }

  setVariableForAddResource() {
    this.updateVariable(this.createResourceForm, "Add Resource", "Add Resource", this.addResource)
  }

  setVariableForUpdateResource() {
    this.updateVariable(this.updateResourceForm, "Update Resource", "Update Resource", this.updateResource)
  }


  updateVariable(formaction: () => void, modalheader: string, modalbutton: string, modalaction: () => void): void {
    this.formHandler = formaction;
    this.modalHeader = modalheader;
    this.resourceButton = modalbutton;
    this.modalAction = modalaction;
  }


  onIsBookChange(resourceForm: FormGroup): void {
    if (!resourceForm.get("isBook")?.value) {
      resourceForm.get("author")?.setValue(null)
      resourceForm.get("publication")?.setValue(null)
    }
  }


  checkFileType(index: number) {

    this.resourceFormArray.at(index).get('allowedExtension').setValue("")
    if (this.resourceFormArray.at(index).get("fileType") != null) {
      this.setAllowedFileExtension(index)
    }

    if (this.resourceFormArray.at(index).get('resourceURL').value && !this.onFileTpeChange(index)) {
      return
    }

    if (!this.resourceFormArray.at(index).get("isExistingResource").value) {
      this.resetFileUploadFields(index)
      this.resetResourceURLField(index)
    }
    this.getResourcesByType(index)
  }

  setAllowedFileExtension(index: number): void {
    if (this.resourceFormArray.at(index).get("fileType").value == this.AUDIO) {
      this.resourceFormArray.at(index).get('allowedExtension').setValue(this.uploadConstant.ALLOWED_AUDIO_EXTENSION)
    } else if (this.resourceFormArray.at(index).get("fileType").value == this.IMAGE) {
      this.resourceFormArray.at(index).get('allowedExtension').setValue(this.uploadConstant.ALLOWED_IMAGE_EXTENSION)
    } else {
      this.resourceFormArray.at(index).get('allowedExtension').setValue(this.uploadConstant.ALLOWED_DOCUMENT_EXTENSION)
    }
  }

  onFileTpeChange(index: number): boolean {
    if (confirm("Uploaded resource will be deleted.\nAre you sure you want to continue?")) {
      this.resourceFormArray.at(index).get('resourceURL').reset()
      return true
    }
    // console.log("this.previousFileType -> ", this.previousFileType);
    this.resourceFormArray.at(index).get('fileType').setValue(this.previousFileType)
    return false
  }

  resetFileUploadFields(index: number): void {
    this.resourceFormArray.at(index).get('displayedFileName').setValue("Select File")
    this.resourceFormArray.at(index).get('previewDisplayedFileName').setValue("Select File")

    this.resourceFormArray.at(index).get('docStatus').setValue("")
    this.resourceFormArray.at(index).get('previewImageDocStatus').setValue("")

    this.resourceFormArray.at(index).get('isDocumentUploadedToServer').setValue(false)
    this.resourceFormArray.at(index).get('isPreviewUploadedToServer').setValue(false)
  }

  resetResourceURLField(index: number): void {
    const URLregex = /^(http(s)?:\/\/)[\w.-]+(\.[\w\.-]+)+[\w\-\._~:?#[\]@!\/$&'\(\)\*\+,;=.]+$/

    if (this.resourceFormArray.at(index).get('fileType').value == "Video" ||
      this.resourceFormArray.at(index).get('fileType').value == "URL") {
      this.resourceFormArray.at(index).get('resourceURL').setValidators([Validators.required,
      Validators.pattern(URLregex), Validators.maxLength(255)])
    } else {
      this.resourceFormArray.at(index).get('resourceURL').setValidators([Validators.required])
    }

    this.resourceFormArray.at(index).get('resourceURL').setValue(null)
    this.resourceFormArray.at(index).get("resourceURL").updateValueAndValidity()

  }

  //Delete uploaded resource
  // deleteUploadedResouce(fileName?: string) {
  //       this.fileOperationService.deleteUploadedFile(fileName).subscribe((data: any) => {
  //       }, (error) => {
  //             console.error(error);
  //       })
  // }

  checkAddExistingResource(index: number) {
    if (this.resourceFormArray.at(index).get('isExistingResource').value) {
      this.resourceFormArray.at(index).get('resourceName').setValidators(null)
      this.resourceFormArray.at(index).get('resourceURL').setValidators(null)
      this.resourceFormArray.at(index).get('id').setValidators([Validators.required])
    } else {
      this.resourceFormArray.at(index).get('resourceName').setValidators([Validators.required])
      this.resourceFormArray.at(index).get('resourceURL').setValidators([Validators.required])
      this.resourceFormArray.at(index).get('id').setValidators(null)
    }
    this.resourceFormArray.at(index).get('resourceName').updateValueAndValidity()
    this.resourceFormArray.at(index).get('resourceURL').updateValueAndValidity()
    this.resourceFormArray.at(index).get('id').updateValueAndValidity()
  }

  //   adds sessions to array
  addSessionToCourse(session: ICourseSession): void {
    if (this.addSubSessionClicked) {
      session.hours = 0
      for (let index = 0; index < session.subSessions.length; index++) {
        session.hours += session.subSessions[index].hours
      }
    }
    // console.log(session)

    if (this.isEditSessionClicked) {
      this.isEditSessionClicked = false
      this.modalButton = "Add Session"
      this.modalHeader = "Add"
      this.addSessionToList.splice(this.editIndex, 1)
    }
    this.addSessionToList.push(session)
    if (session.subSessions) {
      session.subSessions.sort((a, b) => (a.order > b.order) ? 1 : -1)
    }
    this.addSessionToList.sort((a, b) => (a.order > b.order) ? 1 : -1)
    this.totalSessions = this.addSessionToList.length
    this.resetSessionForm()
  }


  //   updates sessions in the list
  updateSessionInList(): void {
    this.addSessionToCourse(this.sessionForm.value)
  }

  //   deletes session from array
  deleteSessionFromList(session: ICourseSession): void {
    // console.log(session);
    if (confirm("Are you sure you want to remove the session from this course?")) {
      for (let index = 0; index < this.addSessionToList.length; index++) {
        if (session.name == this.addSessionToList[index].name) {
          this.addSessionToList.splice(index, 1)
          this.totalSessions = this.addSessionToList.length
        }
      }
    }
    // console.log(this.addSessionToList);
  }


  //   patch's session to the session form for edit
  onEditSessionClick(session: ICourseSession, index: number): void {
    this.onClickScroll("editSession")
    this.createSessionForm()
    this.isEditSessionClicked = true
    this.addNewSession = true
    this.addSubSessionClicked = false
    this.modalHeader = "Update"
    this.modalButton = "Update Session"
    // this.setVariableForUpdateSession()
    if (session.subSessions) {
      for (let index = 0; index < session.subSessions.length; index++) {
        this.addSubSessionClicked = true
        this.createSubSessionForm()
      }
    }
    // console.log(session);
    this.editIndex = index
    this.sessionForm.patchValue(session)
  }

  //   patch's existsing sessions
  onEditExistingSessionClick(session: ICourseSession, index: number): void {
    this.onClickScroll("editSession")
    this.createSessionForm()

    this.isUpdateClicked = true
    this.addNewSession = true
    this.addSubSessionClicked = false

    this.modalHeader = "Update"
    this.modalButton = "Update Session"

    // console.log(session);

    if (session.subSessions) {
      for (let index = 0; index < session.subSessions.length; index++) {
        this.addSubSessionClicked = true
        this.createSubSessionForm()
      }
    }

    if (!session.subSessions) {
      session.subSessions = []
    }

    this.editIndex = index
    this.sessionForm.patchValue(session)
  }

  onResourceSelect(event: any, fileType: string, index: number): void {
    let files = event.target.files
    if (this.PREVIEW === fileType) {
      this.previewFileSelect(files, index)
      return
    }
    this.docFileSelect(files, fileType, index)
  }

  previewFileSelect(files: any, index): void {
    this.resourceFormArray.at(index).get('previewImageDocStatus').setValue("")

    if (files && files.length) {
      let file = files[0]
      // check if file extension is valid
      let err = this.fileOperationService.isImageFileValid(file)
      if (err != null) {
        this.resourceFormArray.at(index).get('previewImageDocStatus').setValue(`<p><span>&#10060;</span> ${err}</p>`)
        return
      }
      // Upload document if it is present.
      this.uploadPreview(file, this.PREVIEW, index)
    }
  }

  docFileSelect(files: any, fileType: string, index: number): void {
    this.resourceFormArray.at(index).get('docStatus').setValue("")

    if (files && files.length) {
      let file = files[0]
      // check if file extension is valid
      let err = this.checkFileExtension(file, fileType)
      if (err != null) {
        this.resourceFormArray.at(index).get('docStatus').setValue(`<p><span>&#10060;</span> ${err}</p>`)
        return
      }
      // Upload document if it is present.
      this.uploadFile(file, fileType, index)
    }
  }

  checkFileExtension(file: any, fileType: string): Error {
    switch (fileType) {
      // case this.DOCUMENT || this.BOOK:
      //   return this.fileOperationService.isDocumentFileValid(file)
      case this.AUDIO:
        return this.fileOperationService.isAudioFileValid(file)
      case this.IMAGE:
        return this.fileOperationService.isImageFileValid(file)
      case this.PREVIEW:
        return this.fileOperationService.isImageFileValid(file)
      default:
        return this.fileOperationService.isDocumentFileValid(file)
    }
  }

  // uploads file to the server
  uploadFile(file: any, fileType: string, index: number): void {

    if (fileType == this.PREVIEW) {
      this.uploadPreview(file, fileType, index)
      return
    }

    this.resourceFormArray.at(index).get('isFileUploading').setValue(true)
    // console.log(file)
    this.fileOperationService.uploadResource(file, fileType).subscribe((data: any) => {
      // console.log(data)
      this.multipleResourcesForm.markAsDirty()
      this.resourceFormArray.at(index).patchValue({
        resourceURL: data
      })
      this.resourceFormArray.at(index).get('displayedFileName').setValue(file.name)
      this.resourceFormArray.at(index).get('isFileUploading').setValue(false)
      this.resourceFormArray.at(index).get('docStatus').setValue(`<p><span class='green'>&#10003;</span> File uploaded.</p>`)
      this.resourceFormArray.at(index).get('isDocumentUploadedToServer').setValue(true)
    }, (error: any) => {
      this.resourceFormArray.at(index).get('isFileUploading').setValue(false)
      this.resourceFormArray.at(index).get('docStatus').setValue(`<p><span>&#10060;</span> ${error}</p>`)
    })
  }

  // uploads file to the server
  uploadPreview(file: any, fileType: string, index: number): void {
    this.resourceFormArray.at(index).get('isPreviewFileUploading').setValue(true)
    this.fileOperationService.uploadResource(file, this.PREVIEW).subscribe((data: any) => {
      // console.log(data);
      this.multipleResourcesForm.markAsDirty()
      this.resourceFormArray.at(index).patchValue({
        reviewURL: data
      })

      this.resourceFormArray.at(index).get('previewDisplayedFileName').setValue(file.name)
      this.resourceFormArray.at(index).get('isPreviewFileUploading').setValue(false)
      this.resourceFormArray.at(index).get('previewImageDocStatus').setValue(`<p><span class='green'>&#10003;</span> File uploaded</p>`)
      this.resourceFormArray.at(index).get('isPreviewUploadedToServer').setValue(true)

      // this.isPreviewUploadedToServer = true
    }, (error: any) => {
      this.resourceFormArray.at(index).get('isPreviewFileUploading').setValue(false)
      this.resourceFormArray.at(index).get('previewImageDocStatus').setValue(`<p><span>&#10060;</span> ${error}</p>`)
    })
  }


  // Used to dismiss modal.
  dismissFormModal(modal: NgbModalRef) {
    for (let index = 0; index < this.resourceFormArray.controls.length; index++) {
      if (this.resourceFormArray.at(index).get('isFileUploading').value ||
        this.resourceFormArray.at(index).get('isPreviewFileUploading').value) {
        alert("Please wait till file is being uploaded")
        return
      }

      if (this.resourceFormArray.at(index).get('isDocumentUploadedToServer').value ||
        this.resourceFormArray.at(index).get('isPreviewUploadedToServer').value) {
        if (!confirm("Uploaded resource will be deleted.\nAre you sure you want to close?")) {
          return
        }
        // this.deleteAllUploadedResources()
        break
      }
    }
    modal.dismiss()
  }

  // delete sub-session from course session form.
  deleteSubSessionInForm(index: number): void {
    if (confirm("Are you sure you want to delete session?")) {
      this.subSessionControlArray.removeAt(index)
      // console.log(this.subSession.controls);
      if (this.subSessionControlArray.controls.length == 0) {
        this.addSubSessionClicked = false
      }
      this.sessionForm.markAsDirty()
    }
  }

  // deletePreviousResource(index: number): void {
  //       let file = this.resourceFormArray.at(index).get("resourceURL").value.split(this.RESOURCE)
  //       this.deleteUploadedResouce(this.RESOURCE+file[1])
  // }

  // deleteAllUploadedResources(): void {
  //       for (let index = 0; index < this.resourceFormArray.controls.length; index++) {
  //             if (this.resourceFormArray.at(index).get('isDocumentUploadedToServer').value) {
  //                   let file = this.resourceFormArray.at(index).get("resourceURL").value.split(this.RESOURCE)
  //                   this.deleteUploadedResouce(this.RESOURCE+file[1])
  //             }
  //       }
  // }

  deleteResource(index: number, modal?: NgbModalRef): void {
    // console.log(this.resourcesForm.at(index));
    if (modal) {
      if (this.resourceFormArray.at(index).get('isFileUploading').value ||
        this.resourceFormArray.at(index).get('isPreviewFileUploading').value) {
        alert("Please wait till file is being uploaded")
        return
      }
      if (this.resourceFormArray.at(index).get('isDocumentUploadedToServer').value ||
        this.resourceFormArray.at(index).get('isPreviewUploadedToServer').value) {
        if (!confirm("Uploaded file will be deleted.\nAre you sure you want to close?")) {
          return
        }
      }

      this.resourceFormArray.at(index).get('isPreviewUploadedToServer').setValue(false)
      this.resourceFormArray.at(index).get('isDocumentUploadedToServer').setValue(false)

      this.resourceFormArray.at(index).get('displayedFileName').setValue("Select file")
      this.resourceFormArray.at(index).get('previewDisplayedFileName').setValue("Select file")

      this.resourceFormArray.at(index).get('docStatus').setValue("")
      this.resourceFormArray.at(index).get('previewImageDocStatus').setValue("")

      // let file = this.resourceFormArray.at(index).get("resourceURL").value.split(this.RESOURCE)
      // this.deleteUploadedResouce(this.RESOURCE+file[1])
    }
    // if (this.resourceFormArray.at(index).valid) {
    //       if (confirm("Are you sure you want to delete the resource?")) {
    //             this.resourceFormArray.removeAt(index)
    //             this.resourceFormArray.markAsDirty()
    //             return
    //       }
    // }
    this.resourceFormArray.removeAt(index)
    this.resourceFormArray.markAsDirty()
  }

  cancelUpdate(): void {
    this.isEditSessionClicked = false
    this.modalButton = "Add Session"
    this.modalHeader = "Add"
    this.resetSessionForm()
  }

  resetSessionForm(): void {
    const subSessions = <FormArray>this.sessionForm.controls.subSessions
    subSessions.controls = []
    this.sessionForm.reset()
    this.addSubSessionClicked = false
    this.isEditSessionClicked = false
    this.isUpdateClicked = false
  }

  // validates session form
  validateSessionForm(session: ICourseSession): void {

    if (this.addSubSessionClicked) {
      this.sessionForm.get("hours").setValidators(null)
      this.utilService.updateValueAndValiditors(this.sessionForm)
    }
    // console.log(this.sessionForm.controls);

    if (this.sessionForm.invalid) {
      this.sessionForm.markAllAsTouched()
      return
    }
    if (this.isUpdateClicked) {
      this.updateSession()
      return
    }
    this.addSessionToCourse(session)
    // this.modalAction()
  }


  doesURLExist(): boolean {
    for (let index = 0; index < this.resourceFormArray.controls.length; index++) {
      if (!this.resourceFormArray.at(index).get("resourceURL").value &&
        !this.resourceFormArray.at(index).get("isExistingResource").value) {
        return false
      }
    }
    return true
  }


  validateResource(): void {
    console.log(this.multipleResourcesForm.controls);

    for (let index = 0; index < this.resourceFormArray.controls.length; index++) {
      if (this.resourceFormArray.at(index).get("isFileUploading").value ||
        this.resourceFormArray.at(index).get('isPreviewFileUploading').value) {
        alert("Please wait till file is being uploaded")
        return
      }
    }

    if (!this.doesURLExist()) {
      alert("Either file should be uploaded or link should be specified.")
      return
    }

    if (this.multipleResourcesForm.invalid) {
      this.multipleResourcesForm.markAllAsTouched()
      return
    }

    this.modalAction()
  }

  // opens modal   
  openModal(modalContent: any) {
    this.modalRef = this.modalService.open(modalContent, {
      ariaLabelledBy: 'modal-basic-title',
      size: 'xl',
      backdrop: 'static',
    })
  }

  setViewField(): void {
    for (let index = 0; index < this.sessionList.length; index++) {
      this.sessionList[index].viewSubSessionClicked = false
      this.sessionList[index].cardColumn = "col-md-6 col-sm-12 d-flex"
    }
  }

  // changes the value of cardColumn on btn click
  changeColValue(session: ICourseSession) {
    if (!session.viewSubSessionClicked) {
      session.cardColumn = "col-md-12 col-sm-12 d-flex"
      session.viewSubSessionClicked = true
      return
    }
    session.cardColumn = "col-lg-6 col-md-6 col-sm-12 d-flex"
    session.viewSubSessionClicked = false
  }

  setVariables() {
    this.isUpdateClicked = false
    this.modalButton = "Add Session"
    this.modalHeader = "Add"
  }

  getFileType(): void {
    this.generalService.getGeneralTypeByType("file_type").subscribe((respond: any) => {
      this.fileTypeList = respond;
    }, (err) => {
      console.error(err)
    })
  }

  getResourceType(): void {
    this.generalService.getGeneralTypeByType("resource_type").subscribe((respond: any) => {
      this.resourceTypeList = respond;
    }, (err) => {
      console.error(err)
    })
  }

  getResourceSubType(): void {
    this.generalService.getGeneralTypeByType("resource_sub_type").subscribe((respond: any) => {
      this.resourceSubTypeList = respond;
    }, (err) => {
      console.error(err)
    })
  }

  getResourcesByType(index: number): void {
    this.resourceFormArray.at(index).get("id").reset()
    this.resourceList = []

    if (this.resourceFormArray.at(index).get("isExistingResource").value) {
      if (this.resourceFormArray.at(index).get("resourceType").value &&
        this.resourceFormArray.at(index).get("fileType").value) {

        let params = new HttpParams()
        this.isResourceLoading = true
        params = params.append("resourceType", this.resourceFormArray.at(index).get("resourceType").value)
        params = params.append("fileType", this.resourceFormArray.at(index).get("fileType").value)

        this.resourceService.getResourcesList(params).subscribe((response: any) => {
          this.resourceList = response.body
        }, (err: any) => {
          console.error(err)
        }).add(() => {
          this.isResourceLoading = false
        })
      }
    }
  }


  // ===============================================================CRUD===============================================================   

  addSession(): void {
    this.spinnerService.loadingMessage = "Adding session"

    // console.log(this.addSessionToList);
    this.courseService.addSession(this.addSessionToList, this.courseID).subscribe((response: any) => {
      // console.log(response);
      alert("Session successfully added")
      this.addSessionToList = []
      this.totalSessions = 0
      this.resetSessionForm()

      this.getSessionForCourse()
      // this.router.navigateByUrl("/course/master")
    }, (err: any) => {
      console.error(err)

      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  updateSession(): void {
    this.spinnerService.loadingMessage = "Updating session"

    // console.log(this.sessionForm.value)
    this.courseService.updateSession(this.sessionForm.value, this.courseID).subscribe((response: any) => {
      // console.log(response);
      alert("Session successfully updated")

      this.resetSessionForm()
      this.getSessionForCourse()
      this.setVariables()
    }, (err: any) => {
      console.error(err);

      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  getSessionForCourse(): void {
    this.spinnerService.loadingMessage = "Getting sessions"

    this.courseService.getSessionsForCourse(this.courseID).subscribe((response: any) => {
      // console.log(response.body);
      this.sessionList = response.body
      this.existingSessions = response.body.length
      if (this.existingSessions == 0) {
        this.addNewSession = true
      }
      this.setViewField()

    }, (err: any) => {

      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  deleteSession(session: ICourseSession): void {
    if (confirm("Are you sure you want to delete this session?")) {
      this.spinnerService.loadingMessage = "Deleting session"

      // console.log(session);
      this.courseService.deleteSession(this.courseID, session.id).subscribe((response: any) => {
        // console.log(response);
        alert("Session successfully deleted")

        this.getSessionForCourse()
      }, (err: any) => {
        console.error(err);

        if (err.statusText.includes('Unknown')) {
          alert("No connection to server. Check internet.")
          return
        }
        alert(err.error.error)
      })
    }
  }

  addResource(): void {
    // console.log(this.multipleResourcesForm.value.resources);
    this.spinnerService.loadingMessage = "Adding Resource"

    this.courseService.addSessionResources(this.multipleResourcesForm.value.resources, this.session.id).subscribe((response: any) => {
      // console.log(response);
      this.modalRef.close()
      alert("Resource successfully added")

      this.getSessionForCourse()
    }, (err: any) => {
      console.error(err);

      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  updateResource(): void {
    // console.log(this.session);
    // console.log(this.multipleResourcesForm.value);
    this.spinnerService.loadingMessage = "Updating Resource"

    this.courseService.updateSessionResource(this.session.id, this.multipleResourcesForm.value.resources).subscribe((response: any) => {
      // console.log(response);
      alert("Resource successfully updated")
      this.modalRef.close()

      this.getSessionForCourse()
    }, (err: any) => {

      console.error(err);
      if (err.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
        return
      }
      alert(err.error.error)
    })
  }

  deleteResources(sessionID: string, resourceID: string): void {
    // console.log(sessionID);
    if (confirm("Are you sure you want to delete all resources?")) {
      this.spinnerService.loadingMessage = "Deleting resource"

      this.courseService.deleteSessionResources(sessionID, resourceID).subscribe((response: any) => {
        alert("Resource succesfully deleted")
        this.getSessionForCourse()

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


  // Get Technology List.
  getTechnologyList(event?: any): void {
    // this.generalService.getTechnologies().subscribe((respond: any[]) => {
    //       this.technologyList = respond;
    // }, (err) => {
    //       console.error(this.utilService.getErrorString(err))
    // })
    let queryParams: any = {}
    if (event && event?.term != "") {
      queryParams.language = event.term
    }
    this.isTechLoading = true
    this.techService.getAllTechnologies(this.techLimit, this.techOffset, queryParams).subscribe((response) => {
      // console.log("getTechnology -> ", response);
      this.technologyList = []
      this.technologyList = this.technologyList.concat(response.body)
    }, (err) => {
      console.error(err)
    }).add(() => {
      this.isTechLoading = false
    })
  }

}
