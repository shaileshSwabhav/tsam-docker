import { Component, OnInit, ViewChild } from '@angular/core';
import { FormGroup, FormBuilder, Validators, FormControl, FormArray } from '@angular/forms';
import { Router, ActivatedRoute, Params } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { IModule } from 'src/app/models/course/module';
import { IModuleTopic } from 'src/app/models/course/module_topic';
import { ITopicProgrammingConcept } from 'src/app/models/course/topic_programming_concept';
import { ITopicProgrammingQuestion } from 'src/app/models/course/topic_programming_question';
import { IProgrammingConcept } from 'src/app/models/programming/concept';
import { IProgrammingQuestion } from 'src/app/models/programming/question';
import { IResource } from 'src/app/models/resource/resource';
import { ConceptModuleService } from 'src/app/service/concept-module/concept-module.service';
import { Role, UploadConstant, UrlConstant } from 'src/app/service/constant';
import { FileOperationService } from 'src/app/service/file-operation/file-operation.service';
import { GeneralService } from 'src/app/service/general/general.service';
import { IPermission } from 'src/app/service/menu/menu.service';
import { ModuleService } from 'src/app/service/module/module.service';
import { ProgrammingConceptService } from 'src/app/service/programming-concept/programming-concept.service';
import { ProgrammingQuestionService } from 'src/app/service/programming-question/programming-question.service';
import { ResourceService } from 'src/app/service/resource/resource.service';
import { UtilityService } from 'src/app/service/utility/utility.service';
import { LocalService } from 'src/app/service/storage/local.service';

@Component({
	selector: 'app-module',
	templateUrl: './module.component.html',
	styleUrls: ['./module.component.css']
})
export class ModuleComponent implements OnInit {


	//assignments
	totalQuestions: number
	totalSubTopics: number

	// List.
	programmingConceptList: IProgrammingConcept[]

	// Module.
	moduleList: IModule[]
	moduleForm: FormGroup

	// Topic.
	moduleTopics: IModuleTopic[]
	moduleTopicForm: FormGroup
	previewModuleTopics: IModuleTopic[]
	topicForm: FormGroup

	// Sub Topic assignment form
	questionForm: FormGroup

	isResourceLoading: boolean

	// Search.
	moduleSearchForm: FormGroup
	searchFormValue: any
	isSearched: boolean
	isProgrammingQuestionSearched: boolean

	// logo
	logoDisplayedFileName: string
	logoDocStatus: string
	isLogoUploadedToServer: boolean
	isLogoFileUploading: boolean

	// Pagination.
	limit: number
	currentPage: number
	totalModules: number
	offset: number
	paginationStart: number
	paginationEnd: number
	totalTopics: number

	// Extra.
	isOperationUpdate: boolean
	isViewMode: boolean
	paginationString: string
	module: IModule
	isAddTopicClick: boolean
	checkSubTopicTime: boolean
	isPreviewClick: boolean
	ifFormInvalid: boolean
	cancelFormClick: boolean
	isPreviewButtonClick: boolean
	isSubTopicClick: boolean

	// Access.
	permission: IPermission
	isFaculty: boolean

	// QUESTIONS
	topicQuestionForm: FormGroup
	selectedTopicQuestionForm: FormGroup
	programmingQuestionSearchForm: FormGroup

	searchProgrammingQuestionFormValue: any

	programmingQuestions: IProgrammingQuestion[]
	totalProgrammingQuestion: number

	topicQuestions: ITopicProgrammingQuestion[]
	assignedTopicQuestions: string[]
	selectedProgrammingQuestions: IProgrammingQuestion[]
	topicID: string

	validatedTopicQuestions: any[]

	//RESOURCE
	multipleResourcesForm: FormGroup
	fileTypeList: any[]
	resourceTypeList: any[]
	resource: IResource[]
	resourceList: IResource[]
	resourceSubTypeList: any[]

	// Concept module.
	conceptModuleForm: FormGroup

	// file upload
	previousFileType: string
	isAddResourceClick: boolean;

	private readonly IGNORE_SEARCH_FIELD: string[] = ["limit", "offset"]

	// Modal.
	modalRef: any
	@ViewChild('moduleFormModal') moduleFormModal: any
	@ViewChild('topicModal') topicModal: any
	@ViewChild('resourceModal') resourceModal: any
	@ViewChild('deleteModal') deleteModal: any
	@ViewChild('addQuestionModal') addQuestionModal: any
	@ViewChild('conceptModuleModal') conceptModuleModal: any

	isViewClicked: boolean;
	modalAction: any;
	//addResource: any;
	resourceModalHeader: string;
	resourceButton: string;

	// Module resource.
	moduleResourcesForm: FormGroup
	conceptList: any[]
	conceptModuleList: any[]
	selectedModuleID: string
	showConceptModuleForm: boolean
	isConceptModuleUpdate: boolean
	conceptModuleMap: Map<number, any[]> = new Map()
	initialQuestionClick: boolean;
	selectedTopicQuestions: any[];

	// Params.
	operation: string

	constructor(
		private moduleService: ModuleService,
		private conceptService: ProgrammingConceptService,
		private generalService: GeneralService,
		private utilService: UtilityService,
		private questionService: ProgrammingQuestionService,
		private fileOps: FileOperationService,
		private urlConstant: UrlConstant,
		private formBuilder: FormBuilder,
		private spinnerService: SpinnerService,
		private modalService: NgbModal,
		private router: Router,
		private route: ActivatedRoute,
		private resourceService: ResourceService,
		private conceptModuleService: ConceptModuleService,
		private localService: LocalService,
		private role: Role,
	) {
		this.initializeVariables()
		this.getAllComponents()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit(): void {
	}

	ngAfterViewInit() {
		if (this.operation && this.operation == "add") {
			this.onAddNewModuleClick()
			this.operation = null
			return
		}
	}

	initializeVariables() {

		this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)
		if (this.isFaculty) {
			this.permission = this.utilService.getPermission(this.urlConstant.BANK_COURSE)
		}
		if (!this.isFaculty) {
			this.permission = this.utilService.getPermission(this.urlConstant.TRAINING_COURSE_MASTER)
		}

		this.programmingConceptList = []
		this.moduleList = [] as IModule[]
		this.programmingQuestions = []
		this.topicQuestions = []
		this.fileTypeList = []
		this.resourceSubTypeList = []
		this.assignedTopicQuestions = []
		this.previewModuleTopics = []
		this.validatedTopicQuestions = []
		this.resource = []

		this.isOperationUpdate = false
		this.isViewMode = false
		this.isSearched = false
		this.isProgrammingQuestionSearched = false
		this.isLogoUploadedToServer = false
		this.isLogoFileUploading = false
		this.isAddTopicClick = false
		this.isAddResourceClick = false
		this.ifFormInvalid = false
		this.cancelFormClick = false
		this.isPreviewButtonClick = false
		this.isSubTopicClick = false

		this.limit = 5
		this.offset = 0
		this.currentPage = 0
		this.paginationStart = 0
		this.paginationEnd = 0

		this.totalProgrammingQuestion = 0

		this.createModuleForm()
		this.createModuleSearchForm()

		this.spinnerService.loadingMessage = "Getting Modules"
		this.logoDisplayedFileName = "Select File"
		this.logoDocStatus = ""

		this.searchFormValue = null
		this.totalQuestions = 0
		this.totalSubTopics = 0

		// Module resource.
		this.conceptList = []
		this.conceptModuleList = []
		this.showConceptModuleForm = false
		this.isConceptModuleUpdate = false
		this.initialQuestionClick = false
		this.checkSubTopicTime = false
		this.isPreviewClick = false
	}

	getAllComponents(): void {
		this.getFileTypeList()
		this.getProgrammingConceptList()
		this.getResourceType()
		this.getResourceSubType()
		this.getFileType()
		this.searchOrGetModules()
	}

	searchOrGetModules() {
		this.route.queryParams.subscribe((queryParams: Params) => {
			console.log(" === subscribing to queryparams === ");
			console.log(queryParams);


			if (this.utilService.isObjectEmpty(queryParams)) {
				this.changePage(1)
				return
			}

			this.moduleSearchForm.patchValue(queryParams)
			this.searchModules()
		})
		// let queryParams = this.route.snapshot.queryParams
		// if (queryParams.operation){
		// 	this.operation = queryParams.operation
		// }


	}


	createModuleForm(): void {
		this.moduleForm = this.formBuilder.group({
			id: new FormControl(null),
			moduleName: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
			logo: new FormControl(null),
		})
	}

	createModuleTopicForm(): void {
		this.moduleTopicForm = new FormGroup({
			id: new FormControl(null),
			topicName: new FormControl(null, [Validators.required, Validators.maxLength(1000)]),
			totalTime: new FormControl(null),
			order: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(100)]),
			moduleID: new FormControl(this.selectedModuleID),
			topicID: new FormControl(null),
			subTopics: new FormArray([]),
			programmingConceptIDs: new FormControl(null),
			topicProgrammingConcept: new FormArray([]),
		})

		this.addSubTopicsToForm()
	}

	get moduleSubTopicControlArray(): FormArray {
		return this.moduleTopicForm.get("subTopics") as FormArray
	}


	addSubTopicsToForm(): void {
		this.moduleSubTopicControlArray.push(new FormGroup({
			id: new FormControl(null),
			topicName: new FormControl(null, [Validators.required, Validators.maxLength(1000)]),
			totalTime: new FormControl(null, [Validators.required, Validators.min(1)]),
			timeHour: new FormControl(null, [Validators.min(0), Validators.max(24)]),
			timeMin: new FormControl(null, [Validators.min(0), Validators.max(59)]),
			order: new FormControl(null, [Validators.required, Validators.min(1), Validators.max(100)]),
			moduleID: new FormControl(this.selectedModuleID),
			topicID: new FormControl(null),
			subTopics: new FormControl(null),
		}))
	}

	removeModuleSubTopic(index: number): void {
		this.moduleSubTopicControlArray.removeAt(index)
	}

	createModuleSearchForm(): void {
		this.moduleSearchForm = this.formBuilder.group({
			moduleName: new FormControl(null),
			limit: new FormControl(this.limit),
			offset: new FormControl(this.offset),
			operation: new FormControl(null),
		})
	}

	createQuestionForm(): void {
		this.questionForm = this.formBuilder.group({
			id: new FormControl(null),
			programmingQuestionID: new FormControl(null),
			subtopicID: new FormControl(null, [Validators.required]),
			programmingConceptID: new FormControl(null, [Validators.required]),
			isActive: new FormControl(null, [Validators.required]),
			order: new FormControl(null, [Validators.required]),
		})
	}

	createResourceForm(): void {
		this.multipleResourcesForm = this.formBuilder.group({
			resources: this.formBuilder.array([], [Validators.required])
		})
		this.addResourcesToForm()
	}

	// Get resources form control.
	get resourceFormArray() {
		return this.multipleResourcesForm.get('resources') as FormArray
	}

	addResourcesToForm() {
		this.resourceFormArray.push(this.formBuilder.group({
			resourceID: new FormControl(null, [Validators.required]),
			resourceType: new FormControl(null),
			resourceSubType: new FormControl(null),
			fileType: new FormControl(null),
		}))

		let len = this.resourceFormArray.length
		this.getResourcesByType(len - 1)
	}

	onManageResourceClick(module: IModule): void {
		this.isViewClicked = false
		this.isAddResourceClick = false
		this.selectedModuleID = module.id
		this.module = module
		this.resource = []
		this.createResourceForm()

		this.resourceModalHeader = "Add Resource"
		this.resourceButton = "Add Resource"

		this.openModal(this.resourceModal, "xl")
	}

	deleteResource(index: number, modal?: NgbModalRef): void {
		this.resourceFormArray.removeAt(index)
		this.resourceFormArray.markAsDirty()
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

	deleteNullValues(object: any): void {
		for (let field in object) {
			if (object[field] === null || object[field] === "" || field == "resourceID") {
				delete object[field];
			}
		}
	}

	getResourcesByType(index: number): void {
		this.resourceFormArray.at(index).get("resourceID").reset()
		this.resourceList = []
		this.isResourceLoading = true

		let queryParams: any = this.resourceFormArray.at(index)?.value
		this.deleteNullValues(queryParams)

		this.resourceService.getResourcesList(queryParams).subscribe((response: any) => {
			this.resourceList = response.body
		}, (err: any) => {
			console.error(err)
		}).add(() => {
			this.isResourceLoading = false
		})
		// }
	}

	onDeleteResourceClick(moduleID: string, resourceID: string): void {
		if (confirm("Are your sure?")) {
			this.deleteResourceByID(moduleID, resourceID)
		}
	}

	deleteResourceByID(moduleID: string, resourceID: string): void {
		this.spinnerService.loadingMessage = "Deleting module"

		this.moduleService.deleteResource(moduleID, resourceID).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllModules()
			alert(response.body)
		}, (error) => {
			// console.error(error)
			if (error.error) {
				alert(error.error)
				return
			}
			alert(error.error?.error)
		}).add(() => {

		})
	}

	getFileType(): void {
		this.generalService.getGeneralTypeByType("file_type").subscribe((respond: any) => {
			this.fileTypeList = respond;
		}, (err) => {
			console.error(err)
		})
	}

	validateResource(): void {
		if (this.multipleResourcesForm.invalid) {
			this.multipleResourcesForm.markAllAsTouched()
			return
		}
	}

	async onAddTopicClick(module: IModule): Promise<void> {
		try {
			this.module = module
			this.selectedModuleID = module.id

			this.previewModuleTopics = []
			this.createModuleTopicForm()

			this.isAddTopicClick = true
			this.isPreviewClick = false

			this.moduleTopics = await this.getModuleTopic()
			this.previewModuleTopics.push(...this.moduleTopics)
			// console.log(this.previewModuleTopics);

			if (this.moduleTopics?.length > 0) {
				this.isAddTopicClick = false
				this.isPreviewClick = true

				this.setPaginationString()
			}

			await this.setConceptModule()
			this.openModal(this.topicModal, "xl")
		} catch (err) {
			alert(err)
		} finally {
		}
	}

	async setConceptModule(): Promise<void> {
		try {
			this.conceptList = await this.getConceptModuleList()

			if (this.conceptList.length > 0) {
				this.getSelectedConceptModuleList()
				this.isConceptModuleUpdate = true
			}

			if (this.conceptList.length == 0) {
				this.isConceptModuleUpdate = false
			}

			this.makeConceptModuleMap()
		} catch (err) {
			alert(err)
		} finally {

		}
	}

	onModuleTopicChange(): void {
		for (let index = 0; index < this.module.moduleTopics.length; index++) {
			if (this.topicForm.get("topicID").value == this.module.moduleTopics[index].id) {
				this.subTopics = this.module.moduleTopics[index]?.subTopics
				this.topic = this.module.moduleTopics[index]
			}
		}
	}

	onEditPreviewClick(previewModuleTopic: IModuleTopic, index: number): void {

		this.createModuleTopicForm()
		this.isAddTopicClick = true
		this.isPreviewClick = false
		this.isSubTopicClick = true
		for (let index = 0; index < previewModuleTopic.subTopics.length - 1; index++) {
			this.addSubTopicsToForm()
		}

		this.previewModuleTopics?.splice(index, 1)
		this.moduleTopicForm.patchValue(previewModuleTopic)
		let programmingConceptIDs: string[] = []
		for (let i = 0; i < previewModuleTopic.topicProgrammingConcept.length; i++) {
			programmingConceptIDs.push(previewModuleTopic.topicProgrammingConcept[i].programmingConceptID)
		}
		this.moduleTopicForm.get("programmingConceptIDs").setValue(programmingConceptIDs)
		for (let i = 0; i < this.moduleSubTopicControlArray.length; i++) {
			let hr = Math.floor(Number(previewModuleTopic.subTopics[i].totalTime) / 60);
			let min = Number(previewModuleTopic.subTopics[i].totalTime) % 60;
			this.moduleSubTopicControlArray.at(i).get('timeHour').setValue(hr)
			this.moduleSubTopicControlArray.at(i).get('timeMin').setValue(min)

		}
		this.moduleTopicForm.markAsDirty()
	}

	assignProgrammingConcept(previewModuleTopic: IModuleTopic): void {

		for (let k = 0; k < previewModuleTopic.subTopics.length; k++) {

			if (!previewModuleTopic.subTopics[k]?.programmingConceptIDs) {
				previewModuleTopic.subTopics[k].programmingConceptIDs = []
			}

			if (previewModuleTopic.subTopics[k].programmingConceptIDs.length !=
				previewModuleTopic.subTopics[k].topicProgrammingConcept.length) {

			}

			for (let j = 0; j < previewModuleTopic.subTopics[k]?.topicProgrammingConcept?.length; j++) {
				if (!previewModuleTopic.subTopics[k].topicProgrammingConcept[j].id) {
					previewModuleTopic.subTopics[k].topicProgrammingConcept.splice(j, 1)
					j--
					continue
				}

				if (!previewModuleTopic.subTopics[k].programmingConceptIDs.includes(
					previewModuleTopic.subTopics[k].topicProgrammingConcept[j].programmingConceptID)) {
					previewModuleTopic.subTopics[k].programmingConceptIDs.push(
						previewModuleTopic.subTopics[k].topicProgrammingConcept[j].programmingConceptID
					)
				}
			}
		}
	}

	onPreviewClick(): void {
		this.isAddTopicClick = false
		this.isPreviewClick = true
	}

	onPreviewCancelClick(): void {
		this.isAddTopicClick = false
		this.isPreviewClick = false
	}

	onAddTopicButtonClick(): void {
		this.createModuleTopicForm()
		this.isAddTopicClick = true
		this.isPreviewClick = false
	}

	onRemoveFromPreviewClick(index: number): void {
		if (confirm("Are you sure")) {
			this.previewModuleTopics.splice(index, 1)
			this.sortObjectByOrder(this.previewModuleTopics)
		}
	}

	async getModuleTopic(): Promise<any> {
		try {
			return new Promise<any>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Getting All topics"

				let queryParams: any = {
					limit: -1,
					offset: 0
				}
				this.moduleService.getModuleTopic(this.selectedModuleID, queryParams).subscribe((response: any) => {
					// this.moduleTopics = response.body
					this.totalTopics = response.headers.get('X-Total-Count')
					resolve(response.body)
				}, (err: any) => {
					console.error(err)
					reject(err?.error?.error)
				})
			})
		} finally {
			// 
		}
	}


	onAddNewModuleClick(): void {
		this.isViewMode = false
		this.isOperationUpdate = false
		this.createModuleForm()
		this.openModal(this.moduleFormModal, 'md')
	}

	addModule(): void {
		this.spinnerService.loadingMessage = "Adding module"

		this.moduleService.addModule(this.moduleForm.value).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllModules()
			alert(response.body)
		}, (error) => {
			// console.error(error)
			if (error.error?.error) {
				alert(error.error?.error)
				return
			}
			alert(error.statusText)
		}).add(() => {

		})
	}

	createTopicForm(): void {
		this.topicForm = this.formBuilder.group({
			topicID: new FormControl(null)
		})
	}

	topic: IModuleTopic
	subTopics: IModuleTopic[]

	async onViewModuleClick(module: IModule): Promise<void> {
		this.isViewMode = true
		this.isOperationUpdate = false
		this.module = module
		this.selectedModuleID = module.id

		this.module.moduleTopics = await this.getModuleTopic()

		this.module.moduleTopics.filter((value) => { value.showDetails = false })

		this.createModuleForm()
		this.moduleForm.patchValue(module)
		this.moduleForm.disable()

		this.openModal(this.moduleFormModal, "xl")
	}

	onTopicChange(): void {
		this.resetTopicFlags()

		for (let index = 0; index < this.module.moduleTopics.length; index++) {
			if (this.module.moduleTopics[index].id == this.topicForm.get('topicID').value) {
				this.topic = this.module.moduleTopics[index]
				this.subTopics = this.topic.subTopics
			}
		}
	}

	onUpdateModuleClick(): void {
		this.isViewMode = false
		this.isOperationUpdate = true
		this.resetLogoVariables()
		this.moduleForm.enable()
	}

	updateModule(): void {
		this.spinnerService.loadingMessage = "Updating module"

		this.moduleService.updateModule(this.moduleForm.value).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllModules()
			alert(response.body)
		}, (error) => {
			if (error.error?.error) {
				alert(error.error.error)
				return
			}
			alert("Check connection")
		}).add(() => {

		})
	}

	onDeleteModuleClick(moduleID: string): void {
		this.openModal(this.deleteModal, 'md').result.then(() => {
			this.deleteModule(moduleID)
		}, (err) => {
			console.error(err)
			return
		})
	}

	deleteModule(moduleID: string): void {
		this.spinnerService.loadingMessage = "Deleting module"

		this.moduleService.deleteModule(moduleID).subscribe((response: any) => {
			this.modalRef.close()
			this.getAllModules()
			alert(response.body)
		}, (error) => {
			// console.error(error)
			if (error.error) {
				alert(error.error?.error)
				return
			}
			alert(error.error?.error)
		}).add(() => {

		})
	}

	onLogoSelect(event: any) {
		this.moduleForm.get("logo").reset()
		this.logoDocStatus = ""
		let files = event.target.files
		if (files && files.length) {
			let file = files[0]
			this.logoDisplayedFileName = file.name
			let err = this.fileOps.isImageFileValid(file)
			if (err != null) {
				this.logoDocStatus = `<p><span>&#10060;</span> ${err}</p>`
				return
			}
			// Upload terms and condition if it is present.]
			this.isLogoFileUploading = true

			this.fileOps.uploadLogo(file,
				this.fileOps.MODULE_FOLDER + this.fileOps.LOGO_FOLDER).subscribe((data: any) => {
					this.moduleForm.markAsDirty()
					this.moduleForm.patchValue({
						logo: data,
					})

					this.logoDisplayedFileName = file.name
					this.isLogoFileUploading = false
					this.isLogoUploadedToServer = true
					this.logoDocStatus = "<p><span class='green'>&#10003;</span> File uploaded.</p>"
				}, (error) => {
					this.isLogoFileUploading = false
					this.logoDocStatus = `<p><span>&#10060;</span> ${error}</p>`
				})
		}
	}

	searchModules(): void {
		console.log(" === searchModules === ");
		this.searchFormValue = { ...this.moduleSearchForm?.value }
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: this.searchFormValue,
		})

		let flag: boolean = true

		for (let field in this.searchFormValue) {
			if (this.searchFormValue[field] === null || this.searchFormValue[field] === "") {
				delete this.searchFormValue[field]
			} else {
				if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
					this.isSearched = true
				}
				flag = false
			}
		}

		if (flag) {
			return
		}
		this.getAllModules()
	}

	changePage(pageNumber: number): void {
		this.moduleSearchForm.get("offset").setValue(pageNumber - 1)
		this.searchModules()
	}

	getProgrammingConceptList(): void {
		this.programmingConceptList = []
		let queryParams: any = {
			limit: -1,
			offset: 0
		}
		this.conceptService.getAllConcepts(queryParams).subscribe((response: any) => {
			this.programmingConceptList = response.body
		}, (err: any) => {
			console.error(err);
		})
	}


	getAllModules() {
		this.moduleList = []
		this.totalModules = 0
		this.spinnerService.loadingMessage = "Getting modules"
		this.searchFormValue.isModuleCount = 1
		this.moduleService.getModule(this.searchFormValue).subscribe((response) => {
			this.moduleList = response.body
			// console.log(this.moduleList)
			this.totalModules = response.headers.get('X-Total-Count')
			this.resetLogoVariables()
		}, (err: any) => {
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err?.error?.error)
		}).add(() => {
			this.setPaginationString()
		})
	}

	setPaginationString() {
		this.paginationString = ''

		let limit = this.moduleSearchForm.get('limit').value
		let offset = this.moduleSearchForm.get('offset').value

		let start: number = limit * offset + 1
		let end: number = +limit + limit * offset

		if (this.totalModules < end) {
			end = this.totalModules
		}
		if (this.totalModules == 0) {
			this.paginationString = ''
			return
		}
		this.paginationString = `${start} - ${end}`
	}

	validateModuleForm(): void {
		if (this.isLogoFileUploading) {
			alert("Please wait till file is being uploaded")
			return
		}

		if (this.moduleForm.invalid) {
			this.moduleForm.markAllAsTouched()
			return
		}

		if (this.isOperationUpdate) {
			this.updateModule()
			return
		}

		this.addModule()
	}

	calculateTotalTopicTime(): void {
		let totalTime = 0

		for (let index = 0; index < this.moduleSubTopicControlArray?.controls?.length; index++) {
			let tempTotalTime = ((Number(this.moduleSubTopicControlArray.at(index).get("timeHour").value)) * 60 +
				Number(this.moduleSubTopicControlArray.at(index).get("timeMin").value))
			this.moduleSubTopicControlArray.at(index).get("totalTime").setValue(tempTotalTime)
			totalTime += tempTotalTime
		}
		this.moduleTopicForm.get("totalTime").setValue(totalTime)
	}

	addResource(resource: any, errors: string[] = []): void {
		this.spinnerService.loadingMessage = "Adding Resource"

		let resourceID: any = {
			resourceID: resource
		}
		this.moduleService.addResource(this.selectedModuleID, resourceID).subscribe((response: any) => {
			this.module.resources.push(resource)
		}, (err: any) => {
			errors.push(err.error.error)
		}).add(() => {

			if (this.ongoingOperations == 0) {
				if (errors.length == 0) {
					alert("Resources successufully added to module")
					this.modalRef.close()
					this.getAllModules()
				} else {
					let errorString = ""
					for (let index = 0; index < errors.length; index++) {
						errorString += (index == 0 ? "" : "\n") + errors[index]
					}
					alert(errorString)
				}
			}
		})
	}

	onDeleteTopicClick(topicID: string): void {
		this.openModal(this.deleteModal, 'md').result.then(() => {
			this.deleteTopic(topicID)
		}, (err) => {
			console.error(err)
			return
		})
	}

	deleteTopic(topicID: string): void {
		this.spinnerService.loadingMessage = "Deleting Topic"

		this.moduleService.deleteModuleTopic(this.selectedModuleID, topicID).subscribe((response: any) => {
			this.modalRef.close()
			this.getModuleTopic()
			alert(response.body)
		}, (error) => {
			// console.error(error)
			if (error.error) {
				alert(error.error)
				return
			}
			alert(error.error?.error)
		}).add(() => {

		})
	}

	onAddTopicProgrammingConceptClick(): void {
		let programmingConceptIDs: string = this.moduleTopicForm.value.programmingConceptIDs
		let topicProgrammingConcepts: ITopicProgrammingConcept[] = this.moduleTopicForm.value.topicProgrammingConcept
		let topicConcept: IProgrammingConcept
		for (let i = 0; i < this.programmingConceptList.length; i++) {
			if (programmingConceptIDs == this.programmingConceptList[i].id) {
				topicConcept = this.programmingConceptList[i]
				break
			}
		}
		// on add
		let tempConcept: ITopicProgrammingConcept = {
			moduleTopicID: this.moduleTopicForm.get("id").value,
			programmingConceptID: programmingConceptIDs[programmingConceptIDs.length - 1],
			programmingConcept: topicConcept
		}

		topicProgrammingConcepts.push(tempConcept)

		this.moduleTopicForm.get("topicProgrammingConcept").setValue(topicProgrammingConcepts)
	}

	onRemoveTopicProgrammingConceptClick(): void {
		let programmingConceptIDs: string[] = this.moduleTopicForm.value.programmingConceptIDs
		let topicProgrammingConcepts: ITopicProgrammingConcept[] = this.moduleTopicForm.value.topicProgrammingConcept
		// on remove
		for (let j = 0; j < topicProgrammingConcepts.length; j++) {
			if (!programmingConceptIDs.includes(topicProgrammingConcepts[j].programmingConceptID)) {
				topicProgrammingConcepts.splice(j, 1)
				j--
			}
		}

		this.moduleTopicForm.get("topicProgrammingConcept").setValue(topicProgrammingConcepts)
	}

	onSubmitResourceClick(): void {
		this.validateResource()
		let errors: string[] = []
		let tempResource: IResource = this.multipleResourcesForm.value
		this.resource.push(tempResource)
		for (let index = 0; index < this.resourceFormArray.length; index++) {
			this.addResource(this.resourceFormArray.at(index).get('resourceID').value, errors)
		}

	}

	onAddNewTopicClick(): void {
		this.calculateTotalTopicTime()

		if (this.moduleTopicForm.invalid) {
			this.moduleTopicForm.markAllAsTouched()
			return
		}

		let topic: IModuleTopic = this.moduleTopicForm.value

		if (!this.isSubTopicOrderUnique(topic)) {
			alert("Sub Topics Order Can't be Same")
			return
		}

		this.assignConceptToTopic(topic)

		this.sortObjectByOrder(topic.subTopics)
		this.previewModuleTopics.push(topic)
		this.sortObjectByOrder(this.previewModuleTopics)
		this.moduleTopicForm.reset()

		this.isAddTopicClick = false
		this.isPreviewClick = true
	}

	isSubTopicOrderUnique(topic: IModuleTopic): boolean {
		let subTopicOrder = new Map();

		for (let i = 0; i < topic.subTopics?.length; i++) {
			if (subTopicOrder.has(topic.subTopics[i].order)) {
				this.ifFormInvalid = true
				return false
			}
			subTopicOrder.set(topic.subTopics[i].order, 1)
		}
		return true
	}

	assignConceptToTopic(topic: IModuleTopic): void {
		for (let j = 0; j < topic.programmingConceptIDs?.length; j++) {
			for (let i = 0; i < this.programmingConceptList?.length; i++) {
				if (topic.programmingConceptIDs[j] == this.programmingConceptList[i].id) {
					let tempConcept: ITopicProgrammingConcept = {
						moduleTopicID: this.moduleTopicForm.get("id").value,
						programmingConceptID: topic.programmingConceptIDs[j],
						programmingConcept: this.programmingConceptList[i]
					}
					topic.topicProgrammingConcept.push(tempConcept)
					break
				}
			}
		}
	}

	async onSubmitTopicClick(): Promise<void> {

		if (this.ifFormInvalid) {
			this.ifFormInvalid = false
			this.isAddTopicClick = true
			this.isPreviewClick = false
			return
		}

		let errors: string[] = []

		await this.addOrUpdateModuleTopic(errors)

		if (errors.length == 0) {
			this.modalRef.close()
			alert("Topic successfully saved")
			this.getAllModules()
			return
		}

		alert(this.createErrorString(errors))
	}

	createErrorString(errors: string[]): string {
		let errorString = ""
		for (let index = 0; index < errors.length; index++) {
			errorString += (index == 0 ? "" : "\n") + errors[index]
		}
		console.error(errorString)
		return errorString
	}

	async addOrUpdateModuleTopic(errors: string[] = []): Promise<void> {
		try {
			for (let index = 0; index < this.previewModuleTopics.length; index++) {
				if (this.previewModuleTopics[index].isTopicAdded) {
					continue;
				}

				if (this.previewModuleTopics[index].id) {
					await this.updateModuleTopic(this.previewModuleTopics[index])
					this.previewModuleTopics.splice(index, 1)
					index = index - 1
					continue
				}

				let topic: any = await this.addModuleTopic(this.previewModuleTopics[index])
				this.previewModuleTopics[index].isTopicAdded = true
				this.previewModuleTopics.splice(index, 1)
				index = index - 1
			}
		} catch (error) {
			errors.push(error)
		}
	}

	async addModuleTopic(topic: IModuleTopic): Promise<void> {
		try {
			return await new Promise<any>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Adding module topic"

				this.moduleService.addModuleTopic(this.selectedModuleID, topic).subscribe((response: any) => {
					resolve(response)
				}, (err: any) => {
					reject(topic.topicName + " - " + err.error.error)
				})
			})
		} finally {

		}
	}

	async updateModuleTopic(topic: IModuleTopic): Promise<void> {
		try {
			return await new Promise<any>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Updating module topic"

				this.moduleService.updateModuleTopic(this.selectedModuleID, topic).subscribe((response: any) => {
					resolve(response)
				}, (err: any) => {
					reject(topic.topicName + " - " + err.error.error)
				})
			})


		} finally {

		}
	}

	createQuestionSearchForm(): void {
		this.programmingQuestionSearchForm = this.formBuilder.group({
			label: new FormControl(null),
			limit: new FormControl(10),
			offset: new FormControl(0),
			programmingConcept: new FormControl(null),
		})
	}

	createTopicQuestionForm(): void {
		this.topicQuestionForm = this.formBuilder.group({
			topicQuestions: new FormArray([])
		})
	}

	createSelectedTopicQuestionForm(): void {
		this.selectedTopicQuestionForm = this.formBuilder.group({
			selectedTopicQuestions: new FormArray([])
		})
	}

	get questionControlArray() {
		return this.topicQuestionForm.get("topicQuestions") as FormArray
	}

	get selectedQuestionControlArray() {
		return this.selectedTopicQuestionForm.get("selectedTopicQuestions") as FormArray
	}

	addQuestion(): void {
		this.questionControlArray.push(this.formBuilder.group({
			id: new FormControl(null),
			programmingQuestionID: new FormControl(null),
			topicID: new FormControl(this.topicID),
			isActive: new FormControl(null),
			isMarked: new FormControl(false),
			showDetails: new FormControl(false),
			// order: new FormControl(null, [Validators.required, Validators.min(1)]),
		}))
	}

	addSelectedQuestionToControlArray(): void {
		this.selectedQuestionControlArray.push(this.formBuilder.group({
			id: new FormControl(null),
			programmingQuestion: new FormControl(null),
			programmingQuestionID: new FormControl(null),
			topicID: new FormControl(this.topicID),
			isActive: new FormControl(null),
			isMarked: new FormControl(false),
			showDetails: new FormControl(false),
			// order: new FormControl(null, [Validators.required, Validators.min(1)]),
		}))
	}

	changeQuestionPage(pageNumber: number): void {
		this.programmingQuestionSearchForm.get("offset").setValue(pageNumber - 1)
		this.searchProgrammingQuestion()
	}

	async searchProgrammingQuestion(): Promise<void> {
		try {
			this.searchProgrammingQuestionFormValue = { ...this.programmingQuestionSearchForm?.value }
			let flag: boolean = true

			for (let field in this.searchProgrammingQuestionFormValue) {
				if (this.searchProgrammingQuestionFormValue[field] === null || this.searchProgrammingQuestionFormValue[field] === "") {
					delete this.searchProgrammingQuestionFormValue[field];
				} else {
					if (!this.IGNORE_SEARCH_FIELD.includes(field)) {
						this.isProgrammingQuestionSearched = true
					}
					flag = false
				}
			}

			// No API call on empty search.
			if (flag) {
				return
			}
			// this.getProgrammingQuestions()
			this.programmingQuestions = await this.getConceptQuestions()
			this.createAddQuestionComponents()
		} catch (err) {
			console.error(err);
		}
	}

	onTopicQuestionFormInput(questionControl: FormGroup): void {
		for (let index = 0; index < this.selectedQuestionControlArray?.controls.length; index++) {
			if (this.selectedQuestionControlArray?.at(index)?.get("programmingQuestionID").value ==
				questionControl.get("programmingQuestionID").value) {
				this.selectedQuestionControlArray?.at(index)?.patchValue(questionControl.value)
			}
		}
	}

	getSelectedTopicProgrammingQuestion(topic: any): void {
		this.spinnerService.loadingMessage = "Getting questions"
		this.selectedTopicQuestions = []

		let queryParams: any = {
			limit: -1,
			offset: 0
		}

		this.moduleService.getTopicProgrammingQuestions(topic.id, queryParams).subscribe((response: any) => {
			this.topicQuestions = response.body
			this.selectedTopicQuestions.push(...this.topicQuestions)

		}, (error: any) => {
			if (error.error?.error) {
				alert(error.error?.error)
				return
			}
			alert(error.error?.error)
		}).add(() => {

		})
	}


	async getTopicProgrammingQuestion(topicID: string): Promise<any> {
		try {
			return new Promise<any>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Getting questions"
				this.topicQuestions = []
				let queryParams: any = {
					limit: -1,
					offset: 0
				}

				this.moduleService.getTopicProgrammingQuestions(topicID, queryParams).subscribe((response: any) => {
					// console.log("topic question ->", response.body)
					resolve(response.body)
				}, (error) => {
					alert(error.error?.error)
					console.error(error);
					reject(error)
				})
			})
		} finally {
		}
	}

	getProgrammingQuestions(): void {
		this.spinnerService.loadingMessage = "Getting programming questions"
		this.programmingQuestions = []
		this.totalProgrammingQuestion = 0
		let programmingConceptID: any[] = []
		this.selectedProgrammingQuestions = []
		if (this.topic?.topicProgrammingConcept?.length == 0) {
			return
		}

		this.programmingQuestionSearchForm.get("programmingConcept").setValue(null)
		delete this.searchProgrammingQuestionFormValue.programmingConcept

		for (let index = 0; index < this.topic.topicProgrammingConcept.length; index++) {
			programmingConceptID.push(this.topic.topicProgrammingConcept[index].programmingConceptID)
		}

		this.programmingQuestionSearchForm.get("programmingConcept").setValue(programmingConceptID)
		this.searchProgrammingQuestionFormValue.programmingConcept = programmingConceptID

		let limit: number = -1
		let offset: number = 0

		this.questionService.getProgrammingQuestions(limit, offset, this.searchProgrammingQuestionFormValue).subscribe((response: any) => {
			this.selectedProgrammingQuestions = response.body
			this.totalProgrammingQuestion = response.headers.get("X-Total-Count")
			this.addSelectedTopicQuestions()
			this.createAddQuestionComponents()
		}, (err: any) => {
			alert(err.error?.error)
		}).add(() => {

		})
	}

	addSelectedTopicQuestions(): void {
		this.createSelectedTopicQuestionForm()
		for (let j = 0; j < this.topicQuestions.length; j++) {
			this.addSelectedQuestionToControlArray()
			this.selectedQuestionControlArray.at(j).patchValue(this.topicQuestions[j])
			this.selectedQuestionControlArray.at(j).get("programmingQuestionID").setValue(this.topicQuestions[j].programmingQuestion.id)
		}
	}

	createAddQuestionComponents(): void {
		this.createTopicQuestionForm()
		this.patchQuestionToForm()
	}

	patchQuestionToForm(): void {
		for (let index = 0; index < this.programmingQuestions.length; index++) {
			this.addQuestion()
			this.questionControlArray.at(index).get("isMarked").setValue(false)
			this.questionControlArray.at(index).get("programmingQuestionID").setValue(this.programmingQuestions[index].id)
			this.questionControlArray.at(index).get("isActive").disable()
		}
	}

	toggleAssignment(questionControl: FormGroup) {
		if (this.assignedTopicQuestions.includes(questionControl.get("programmingQuestionID").value)) {
			this.removeAssignedProgrammingQuestion(questionControl)
			return
		}
		// add
		if (!this.assignedTopicQuestions.includes(questionControl.get("programmingQuestionID").value)) {
			this.assignProgrammingQuestion(questionControl)
		}
	}

	assignProgrammingQuestion(questionControl: FormGroup): void {
		this.assignedTopicQuestions.push(questionControl.get("programmingQuestionID").value)

		questionControl.get("programmingQuestionID").setValidators([Validators.required])
		questionControl.get("isActive").setValidators([Validators.required])

		questionControl.get("isActive").enable()

		questionControl.get("isMarked").setValue(true)
		questionControl.get("isActive").setValue(true)

		this.utilService.updateValueAndValiditors(questionControl)
		this.assignToSelectedProgrammingQuestion(questionControl)
	}

	assignToSelectedProgrammingQuestion(questionControl: FormGroup): void {
		let tempSelectedTopicQuestions = this.selectedTopicQuestionForm.value.selectedTopicQuestions

		for (let index = 0; index < tempSelectedTopicQuestions.length; index++) {
			if (this.assignedTopicQuestions.includes(tempSelectedTopicQuestions.programmingQuestionID)) {
				return
			}
		}

		this.addSelectedQuestionToControlArray()
		let len = this.selectedQuestionControlArray.controls.length
		this.selectedQuestionControlArray.at(len - 1).patchValue(questionControl.value)
	}

	removeAssignedProgrammingQuestion(questionControl: FormGroup): void {
		let programmingQuestionID = questionControl.get("programmingQuestionID").value
		let index = this.assignedTopicQuestions.indexOf(questionControl.get("programmingQuestionID").value)
		this.assignedTopicQuestions.splice(index, 1)

		questionControl.get("programmingQuestionID").clearValidators()
		questionControl.get("isActive").clearValidators()

		questionControl.get("isActive").disable()

		questionControl.get("isActive").setValue(null)
		questionControl.get("isMarked").setValue(false)

		this.utilService.updateValueAndValiditors(questionControl)
		this.removeSelectedProgrammingQuestion(programmingQuestionID)
	}

	removeSelectedProgrammingQuestion(programmingQuestionID: string): void {

		for (let index = 0; index < this.selectedQuestionControlArray.controls.length; index++) {
			if (programmingQuestionID == this.selectedQuestionControlArray.at(index).get("programmingQuestionID").value) {
				this.selectedQuestionControlArray.removeAt(index)
				return
			}
		}
	}

	deleteExtraFormFields(assignment: ITopicProgrammingQuestion): void {
		for (let field in assignment) {
			if (assignment[field] === null || assignment[field] === "") {
				delete assignment[field];
			}
			if (field == "isMarked" || field == "showDetails") {
				delete assignment[field];
			}
		}
	}

	addTopicProgrammingQuestion(assignment: ITopicProgrammingQuestion, errors: string[] = []): Promise<void> {
		try {
			return new Promise<void>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Saving assignment"
				// this.startSpinner()
				this.deleteExtraFormFields(assignment)
				this.moduleService.addTopicProgrammingQuestion(this.topic.id, assignment).subscribe((response: any) => {
					resolve()
				}, (err: any) => {
					errors.push(err.error.error)
				})
			})
		} finally {
			// this.stopSpinner()
		}
	}

	updateTopicProgrammingAssignment(assignment: ITopicProgrammingQuestion, errors: string[] = []): Promise<void> {
		try {
			return new Promise<void>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Saving assignment"
				// this.startSpinner()
				this.deleteExtraFormFields(assignment)
				// console.log(assignment);

				this.moduleService.updateTopicProgrammingAssignment(this.topic.id, assignment).subscribe((response: any) => {
				}, (err: any) => {
					errors.push(err.error.error)
				})
			})
		} finally {
			// this.stopSpinner()
		}
	}


	deleteTopicProgrammingAssignment(topicAssignmentID: string): void {
		this.spinnerService.loadingMessage = "Deleting programming assignment"

		this.moduleService.deleteTopicProgrammingAssignment(this.selectedModuleID, topicAssignmentID).subscribe((response: any) => {
			alert("Programming Assingment removed from topic")
		}, (err: any) => {
			this.setPaginationString()
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(err.error.error)
		}).add(() => {

		})
	}

	resetTopicFlags(): void {
		this.moduleTopics?.filter((value) => {
			value.isEditQuestionClick = false
			value.isAddQuestionClick = false
		})
	}

	removeUnmarkedAssignments(topicQuestions: any[]): void {
		for (let index = 0; index < topicQuestions.length; index++) {
			if (!topicQuestions[index].isMarked) {
				topicQuestions.splice(index, 1)
				index--
			}
		}
	}


	checkAssignmentOrder(): boolean {
		let assignmentMap = new Map<number, number>()

		for (let index = 0; index < this.validatedTopicQuestions.length; index++) {
			if (assignmentMap.get(this.validatedTopicQuestions[index].order) == 1) {
				return true
			}
			assignmentMap.set(this.validatedTopicQuestions[index].order, 1)
			continue
		}
		return false
	}

	// setModuleTopicTimeRequired():void{
	// 		this.topicAssignmentForm.get('totalTime').setValue(Number(this.topicAssignmentForm.get('timeHour').value) * 60 + Number(this.topicAssignmentForm.get('timeMin').value));
	// 		console.log("time",this.topicAssignmentForm.get('totalTime').value)
	// 	  }

	async validateTopicQuestion(): Promise<void> {
		try {
			if (this.topicQuestionForm?.invalid || this.selectedTopicQuestionForm?.invalid) {
				this.topicQuestionForm?.markAllAsTouched()
				this.selectedTopicQuestionForm?.markAllAsTouched()
				return
			}

			let errors: string[] = []
			this.validatedTopicQuestions = []
			this.validatedTopicQuestions = this.selectedTopicQuestionForm.value.selectedTopicQuestions

			for (let index = 0; index < this.validatedTopicQuestions?.length; index++) {
				if (this.validatedTopicQuestions[index].id) {
					await this.updateTopicProgrammingAssignment(this.validatedTopicQuestions[index], errors)
					continue
				}
				await this.addTopicProgrammingQuestion(this.validatedTopicQuestions[index], errors)
			}

			if (errors.length == 0) {
				alert("Assignments successufully saved")
				this.resetTopicFlags()
				this.topic.topicProgrammingQuestions = this.validatedTopicQuestions
				this.getSelectedTopicProgrammingQuestion(this.topic)
				this.topic.topicProgrammingQuestions = await this.getTopicProgrammingQuestion(this.topic.id)
				return
			}

			alert(this.createErrorString(errors))

		} catch (err) {
			console.error(err);
		}
	}
	// enableAssignmentControlArray(): void {
	// 	for (let index = 0; index < this.assignmentControlArray.length; index++) {
	// 		this.assignmentControlArray.at(index).get('isMarked').enable()
	// 	}
	// }

	resetAssignmentSearchAndGetAll(): void {
		let programmingQuestionLimit = this.programmingQuestionSearchForm.get("limit").value
		let programmingQuestionOffset = this.programmingQuestionSearchForm.get("offset").value

		this.createQuestionSearchForm()
		this.programmingQuestionSearchForm.get("limit").setValue(programmingQuestionLimit)
		this.programmingQuestionSearchForm.get("offset").setValue(programmingQuestionOffset)

		this.searchProgrammingQuestionFormValue = null
		this.isProgrammingQuestionSearched = false
		this.searchProgrammingQuestion()
	}


	// Get file type list.
	getFileTypeList(): void {
		this.generalService.getGeneralTypeByType("file_type").subscribe((respond: any) => {
			this.fileTypeList = respond
		}, (err) => {
			console.error(err)
		})
	}

	// ==================================================== MODULE RESOURCE FUNCTIONS ==========================================================================

	// On clicking module resource.
	OnModuleResourceClick(): void {
		this.createModuleResourcesForm()
		this.openModal(this.resourceModal, "lg")
	}

	// Create module resources form.
	createModuleResourcesForm() {
		this.moduleResourcesForm = this.formBuilder.group({
			resources: this.formBuilder.array([], [Validators.required])
		})
		this.addResourcesToForm()
	}

	// Delete resource from module resource form.
	deleteModuleResource(index: number): void {
		if (confirm("Are you sure you want to delete the resource?")) {
			this.moduleResourcesForm.markAsDirty()
			this.resourceFormArray.removeAt(index)
		}
	}

	// ==================================================== CONCEPT TREE FUNCTIONS ==========================================================================

	// Redirect to modules concepts page.
	redirectToModuleConcepts(moduleID: string, moduleName: string): void {
		let url: string
		if (!this.isFaculty) {
			url = "/training/module/module-concept"
		}
		if (this.isFaculty) {
			url = "/bank/module/module-concept"
		}
		this.router.navigate([url], {
			queryParams: {
				"moduleID": moduleID,
				"moduleName": moduleName,
			}
		}).catch(err => {
			console.error(err)

		});
	}

	// On clicking concept module.
	onManageConceptsClick(moduleID: string): void {
		this.selectedModuleID = moduleID
		this.showConceptModuleForm = false
		this.getConceptModuleList()
		this.openModal(this.conceptModuleModal, "lg")
	}

	// Create programming concept module form.
	createConceptModuleForm() {
		this.conceptModuleForm = this.formBuilder.group({
			conceptModules: new FormArray([]),
		})
	}

	// Add new concept module to concept module form.
	addConceptModuleToConceptModuleForm(): void {
		this.conceptModuleControlArray.push(this.formBuilder.group({
			id: new FormControl(null),
			moduleID: new FormControl(this.selectedModuleID),
			programmingConceptID: new FormControl(null, [Validators.required]),
			level: new FormControl(null, [Validators.required, Validators.min(0), Validators.max(99)]),
		}))
	}

	// Get concept modules array from concept form.
	get conceptModuleControlArray(): FormArray {
		return this.conceptModuleForm.get("conceptModules") as FormArray
	}

	// Delete concept module from concept module form.
	deleteConceptModule(index: number) {
		if (confirm("Are you sure you want to delete the concept?")) {
			this.conceptModuleForm.markAsDirty()
			this.conceptModuleControlArray.removeAt(index)
		}
	}

	// Get concept module list for module.
	async getConceptModuleList(): Promise<any> {
		try {
			return new Promise<any>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Getting concepts"

				let queryParams: any = {
					limit: -1,
					offset: 0,
					moduleID: this.selectedModuleID
				}
				this.conceptModuleService.getAllModuleProgrammingConcepts(queryParams).subscribe((response: any) => {
					resolve(response.body)
				}, (error) => {
					console.error(error)
					// if (error.error?.error) {
					// 	alert(error.error?.error)
					// 	return
					// }
					reject(error?.error?.error)
				})
			})
		} finally {
			// 
		}
	}


	getSelectedConceptModuleList(): void {
		this.conceptModuleList = []
		for (let i = 0; i < this.conceptList?.length; i++) {
			for (let j = 0; j < this.programmingConceptList?.length; j++) {
				if (this.programmingConceptList[j].id == this.conceptList[i].programmingConceptID) {
					this.conceptModuleList.push(this.programmingConceptList[j])
					break;
				}
			}
		}
	}

	// On clicking add new concept module button.
	onAddNewConceptModuleClick(): void {
		this.createConceptModuleForm()
		this.addConceptModuleToConceptModuleForm()
		this.showConceptModuleForm = true
	}

	// On clicking view concept modules.
	onViewConceptModuleClick(): void {
		this.createConceptModuleForm()
		for (let i = 0; i < this.conceptList.length; i++) {
			this.addConceptModuleToConceptModuleForm()
		}
		this.conceptModuleForm.get('conceptModules').patchValue(this.conceptList)
		this.showConceptModuleForm = true
	}

	// On clicking concept module form submit button.
	onConceptModuleFormSubmit(): void {
		if (this.conceptModuleForm.invalid) {
			this.conceptModuleForm.markAllAsTouched()
			return
		}
		if (this.isConceptModuleUpdate) {
			this.updateConceptModules()
			return
		}
		this.addConceptModules()
	}

	// Add new concept modules.
	addConceptModules(): void {
		this.spinnerService.loadingMessage = "Adding Concepts to Module"


		this.conceptModuleService.addConceptModules(this.conceptModuleForm.get('conceptModules').value)
			.subscribe((response: any) => {
				this.getConceptModuleList()
				this.showConceptModuleForm = false
				alert(response)
			}, (error) => {
				console.error(error);
				;
				if (error.error?.error) {
					alert(error.error?.error);
					return;
				}
				alert(error.error?.error);
			})
	}

	// Update concept modules.
	updateConceptModules(): void {
		this.spinnerService.loadingMessage = "Updating Concepts to Module"


		this.conceptModuleService.updateConceptModules(this.selectedModuleID,
			this.conceptModuleForm.get('conceptModules').value).subscribe((response: any) => {
				this.getConceptModuleList()
				this.showConceptModuleForm = false
				alert(response)
			}, (error) => {
				console.error(error);
				;
				if (error.error?.error) {
					alert(error.error?.error);
					return;
				}
				alert(error.error?.error);
			})
	}

	// Make map of concept module level and its related concept modules.
	makeConceptModuleMap(): void {
		let tempConceptModuleArray: any[] = []
		this.conceptModuleMap = new Map()
		for (let i = 0; i < this.conceptList.length; i++) {
			if (!this.conceptModuleMap.get(this.conceptList[i].level)) {
				this.conceptModuleMap.set(this.conceptList[i].level, [])
			}
			tempConceptModuleArray = this.conceptModuleMap.get(this.conceptList[i].level)
			tempConceptModuleArray.push(this.conceptList[i])
			this.conceptModuleMap.set(this.conceptList[i].level, tempConceptModuleArray)
		}
	}

	// =========================================================================================================

	openModal(content: any, size: string = "lg"): NgbModalRef {
		let options: NgbModalOptions = {
			ariaLabelledBy: 'modal-basic-title', keyboard: false,
			backdrop: 'static', size: size
		}

		this.modalRef = this.modalService.open(content, options)
		return this.modalRef
	}

	dismissFormModal(modal: NgbModalRef): void {
		if (this.isLogoFileUploading) {
			alert("Please wait till file is being uploaded")
			return
		}

		if (this.isLogoUploadedToServer) {
			if (!confirm("Uploaded file will be deleted.\nAre you sure you want to close?")) {
				return
			}
		}

		modal.dismiss()
		this.resetLogoVariables()
	}

	resetLogoVariables(): void {
		this.logoDisplayedFileName = "Select File"
		this.logoDocStatus = ""
		this.isLogoUploadedToServer = false
		this.isLogoFileUploading = false
	}

	resetSearchAndGetAll(): void {
		this.resetSearchForm()
		this.searchFormValue = null
		this.isSearched = false
		// this.router.navigate(['/module'])
		// this.changePage(1)
		this.searchModules()
	}

	resetSearchForm(): void {

		this.limit = this.moduleSearchForm.get('limit').value
		this.offset = this.moduleSearchForm.get('offset').value

		this.moduleSearchForm.reset({
			limit: this.limit,
			offset: this.offset,
		})
	}

	sortObjectByOrder(object: any[]): void {
		object.sort((a, b) => (a.order < b.order ? -1 : 1));
	}


	// ============================================== PROGRAMMING QUESTIONS ===================================================

	async onManageQuestionClick(module: IModule): Promise<void> {
		try {
			this.module = module
			this.selectedModuleID = module.id
			this.module.moduleTopics = await this.getModuleTopic()
			this.moduleTopics = this.module.moduleTopics

			this.resetTopicFlags()
			this.addOrEditQuestion()
			this.openModal(this.addQuestionModal, "xl")
		} catch (err) {
			console.error(err)
		}
	}

	async onAddQuestionClick(topic: IModuleTopic): Promise<void> {
		try {
			if (topic.isAddQuestionClick) {
				return
			}
			this.resetTopicFlags()
			this.topic = topic
			topic.isAddQuestionClick = true
			if (topic?.topicProgrammingConcept?.length == 0) {
				return
			}

			this.addOrEditQuestion()
			this.searchProgrammingQuestion()
			// this.programmingQuestions = await this.getConceptQuestions()
			// this.createAddQuestionComponents()
		} catch (err) {
			console.error(err);
		}
	}

	async onEditQuestionClick(topic: IModuleTopic): Promise<void> {
		try {
			if (topic.isEditQuestionClick) {
				return
			}
			this.resetTopicFlags()
			this.topic = topic
			topic.isEditQuestionClick = true
			this.addOrEditQuestion()
			this.topicQuestions = await this.getTopicProgrammingQuestion(topic.id)
			this.addSelectedTopicQuestions()
		} catch (err) {
			console.error(err);
		}
	}

	addOrEditQuestion(): void {
		this.createTopicQuestionForm()
		this.createQuestionSearchForm()
		this.createSelectedTopicQuestionForm()

		this.assignedTopicQuestions = []
		this.isProgrammingQuestionSearched = false
	}

	async getConceptQuestions(): Promise<any> {
		try {
			return new Promise<any>((resolve, reject) => {
				this.spinnerService.loadingMessage = "Getting questions"
				this.programmingQuestions = []
				this.totalProgrammingQuestion = 0
				this.selectedProgrammingQuestions = []
				this.searchQuestionConcept()

				// console.log("searchProgrammingQuestionFormValue -> ", this.searchProgrammingQuestionFormValue);

				// this.questionService.getProgrammingQuestions(limit, offset, this.searchProgrammingQuestionFormValue).subscribe((response: any) => {
				this.questionService.getTopicProgrammingQuestion(this.searchProgrammingQuestionFormValue).subscribe((response: any) => {
					this.totalProgrammingQuestion = response.headers.get("X-Total-Count")
					// console.log(response.body);
					resolve(response.body)
				}, (err: any) => {
					alert(err.error?.error)
					reject(err)
				})
			})
		} finally {
			// this.stopSpinner()
		}
	}

	searchQuestionConcept(): void {
		let programmingConceptID: any[] = []
		let topicIDs: string[] = []

		this.programmingQuestionSearchForm.get("programmingConcept").setValue([])

		for (let index = 0; index < this.topic?.topicProgrammingConcept?.length; index++) {
			programmingConceptID.push(this.topic.topicProgrammingConcept[index].programmingConceptID)
		}

		this.programmingQuestionSearchForm.get("programmingConcept").setValue(programmingConceptID)
		this.searchProgrammingQuestionFormValue['programmingConcept'] = programmingConceptID

		if (!this.searchProgrammingQuestionFormValue?.programmingConcept) {
			delete this.searchProgrammingQuestionFormValue?.programmingConcept
		}

		for (let index = 0; index < this.module.moduleTopics.length; index++) {
			// if (this.topic.id != this.module.moduleTopics[index].id) {
			topicIDs.push(this.module.moduleTopics[index].id)
			// }
		}

		this.searchProgrammingQuestionFormValue['topicIDs'] = topicIDs
		// this.searchProgrammingQuestionFormValue.topicID = this.topic.id
	}

}
