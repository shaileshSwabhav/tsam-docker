import { HttpResponse } from '@angular/common/http';
import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { IBatchSessionDetail } from 'src/app/models/batch/session_detail';
import { IModule } from 'src/app/models/course/module';
import { IModuleTopic } from 'src/app/models/course/module_topic';
import { IResource } from 'src/app/models/resource/resource';
import { BatchSessionService } from 'src/app/service/batch-session/batch-session.service';
import { Role } from 'src/app/service/constant';
import { GeneralService } from 'src/app/service/general/general.service';
import { ModuleService } from 'src/app/service/module/module.service';
import { ResourceService } from 'src/app/service/resource/resource.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { UtilityService } from 'src/app/service/utility/utility.service';

@Component({
	selector: 'app-batch-session-details',
	templateUrl: './batch-session-details.component.html',
	styleUrls: ['./batch-session-details.component.css']
})
export class BatchSessionDetailsComponent implements OnInit {

	batchSessions: IBatchSessionDetail[]
	batchID: string;
	totalbatchSession: number
	totalTopicBatchSessionTime: number
	isViewClicked: boolean;
	isAddResourceClick: boolean;
	preRequisiteForm: FormGroup

	//RESOURCE
	multipleResourcesForm: FormGroup
	fileTypeList: any[]
	resourceTypeList: any[]
	resource: IResource[]
	resourceList: IResource[]
	resourceSubTypeList: any[]

	// Access.
	loginID: string

	// Flags.
	isTalent: boolean
	isFaculty: boolean

	//Modal
	modalRef: any
	@ViewChild('resourceModal') resourceModal: any
	@ViewChild('assignmentModal') assignmentModal: any
	@ViewChild('preRequisiteFormModal') preRequisiteFormModal: any

	module: IModule;
	isViewAssignmentClick: boolean;
	moduleTopic: IModuleTopic;
	currentIndex: number;
	isViewPreRequisiteClick: boolean;
	batchSession: IBatchSessionDetail;
	modules: IModule[];
	moduleTopics: IModuleTopic[];
	isViewMode: boolean;
	isOperationUpdate: boolean;


	// ckeditor
	ckConfig: any
	isResourceLoading: boolean;
	moduleID: string;

	loginName: string

	isEditSessionPlan:boolean

	constructor(
		private batchSessionService: BatchSessionService,
		private modalService: NgbModal,
		public utilService: UtilityService,
		private formBuilder: FormBuilder,
		private spinnerService: SpinnerService,
		private route: ActivatedRoute,
		private router: Router,
		private resourceService: ResourceService,
		private generalService: GeneralService,
		private moduleService: ModuleService,
		private localService: LocalService,
		private role: Role,
	) {
		this.editorConfig()
		this.initializeVariables()
		this.getAllComponent()
		this.extractQueryparams()
	}


	get ongoingOperations() {
		return this.spinnerService.ongoingOperations
	}

	ngOnInit(): void {
	}

	editorConfig(): void {
		this.ckConfig = {
			extraPlugins: 'codeTag,kbdTag',
			removePlugins: "exportpdf",
			toolbar: [
				{ name: 'styles', items: ['Styles', 'Format'] },
				{
					name: 'basicstyles', groups: ['basicstyles', 'cleanup'],
					items: ['Bold', 'Italic', 'Underline', 'Strike', 'Subscript', 'Superscript', 'RemoveFormat', 'Code', 'Kbd']
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
				{ name: 'links' },
			],
			removeButtons: "",
			language: "en",
			resize_enabled: false,
			width: "100%",
			height: "80%",
			forcePasteAsPlainText: false,
			allowedContent: true,
			extraAllowedContent: 'img',
		}
	}

	initializeVariables() {
		this.batchSessions = []
		this.totalbatchSession = 0
		this.totalTopicBatchSessionTime = 0

		this.isAddResourceClick = false
		this.isOperationUpdate = false
		this.resource = []
		this.moduleTopics = []
		this.resourceTypeList = []
		this.resourceSubTypeList = []
		this.fileTypeList = []
		this.resourceList = []

		// Access.
		this.isTalent = (this.role.TALENT == this.localService.getJsonValue("roleName") ? true : false)
		this.isFaculty = (this.role.FACULTY == this.localService.getJsonValue("roleName") ? true : false)
		this.loginID = this.localService.getJsonValue("loginID")

		this.loginName = this.localService.getJsonValue("firstName") + " " + this.localService.getJsonValue("lastName")

		this.isEditSessionPlan = false
	}

	extractQueryparams(): void {
		let queryparams = this.route.snapshot.queryParams
		console.log(queryparams);

		this.batchID = queryparams.batchID
		if (this.batchID) {
			this.getBatchSessionPlan()
		}
	}

	getAllComponent(): void {
		this.getResourceType()
		this.getResourceSubType()
		this.getFileType()
	}

	// modal
	@ViewChild("deleteModal") deleteModal: any

	onDeleteSessionPlanClick(): void {
		// console.log(subTopicID, topicAssignmentID);
		this.openModal(this.deleteModal, 'md').result.then(() => {
			this.deleteSessionPlan()
		}, (err) => {
			console.error(err);
			return
		})
	}

	deleteSessionPlan(): void {
		this.spinnerService.loadingMessage = "Deleting plan"

		this.batchSessionService.deleteSessionPlan(this.batchID, this.loginID).subscribe((response: any) => {
			// console.log(response);
			alert("Session plan successfully deleted")
			this.getBatchSessionPlan()
		}, (err) => {
			if (err.statusText.includes('Unknown')) {
				alert("No connection to server. Check internet.")
				return
			}
			alert(this.utilService.getErrorString(err))
			console.error(err)
		})
	}


	getBatchSessionPlan(): void {
		this.spinnerService.loadingMessage = "Getting Plan"

		this.batchSessions = []
		this.totalbatchSession = 0
		let queryParams: any
		if (this.isFaculty) {
			queryParams = {
				facultyID: this.loginID
			}
		}
		this.batchSessionService.getAllBatchSessionPlan(this.batchID, queryParams).subscribe((response: any) => {
			this.batchSessions = response.body
			// console.log("batchSessions", this.batchSessions)
			this.totalbatchSession = this.batchSessions.length
			this.countSubTopics()
			if (this.isTalent) {
				this.formatSessionPlanListForTalent()
			}
		}, (error) => {
			if (error.error?.error) {
				alert(error.error?.error)
				return
			}
			alert(error.statusText)
		}, () => {

		})
	}

	countSubTopics(): void {
		for (let i = 0; i < this.batchSessions.length; i++) {
			this.batchSessions[i].totalModuleSubTopic = 0
			for (let j = 0; j < this.batchSessions[i]?.module?.length; j++) {
				this.batchSessions[i].module[j].totalSubTopics = 0
				for (let k = 0; k < this.batchSessions[i]?.module[j]?.moduleTopics?.length; k++) {
					this.batchSessions[i].module[j].totalSubTopics += this.batchSessions[i]?.module[j]?.moduleTopics[k]?.subTopics?.length
					this.batchSessions[i].module[j].moduleTopics[k].totalTime = 0
					for (let x = 0; x < this.batchSessions[i]?.module[j]?.moduleTopics[k]?.subTopics?.length; x++) {
						this.batchSessions[i].module[j].moduleTopics[k].totalTime +=
							this.batchSessions[i]?.module[j].moduleTopics[k].subTopics[x].totalTime
					}
				}
				this.batchSessions[i].totalModuleSubTopic += this.batchSessions[i].module[j].totalSubTopics
			}
			// console.log("totalModuleSubTopic -> ", this.batchSessionsPlan[i].totalModuleSubTopic);
		}
	}

	moduleResources: IResource[]

	onManageResourceClick(module: IModule): void {
		this.isViewClicked = false
		this.isAddResourceClick = false
		this.module = module
		this.moduleResources = []
		this.getModuleResource()
		this.createResourceForm()
		this.openModal(this.resourceModal, "xl")
	}

	getModuleResource(): void {
		this.spinnerService.loadingMessage = "Getting resources"
		this.moduleResources = []
		this.moduleService.getModuleResource(this.module.id).subscribe((response: HttpResponse<IResource[]>) => {
			this.moduleResources = response.body
			// console.log(response.body);
		}, (err: any) => {
			console.error(err)
		})
	}

	// Format fields of session plan list for talent.
	formatSessionPlanListForTalent(): void {
		for (let i = 0; i < this.batchSessions.length; i++) {
			for (let j = 0; j < this.batchSessions[i]?.module.length; j++) {
				// if (this.batchSessions[i]?.module[j].resources.length > 0) {
				// 	this.batchSessions[i].module[j].isResources = true
				// }
				for (let k = 0; k < this.batchSessions[i]?.module[j]?.moduleTopics.length; k++) {
					let tempAssignmnetList: any = []
					for (let l = 0; l < this.batchSessions[i]?.module[j]?.moduleTopics[k]?.batchTopicAssignment?.length; l++) {
						if (this.batchSessions[i]?.module[j]?.moduleTopics[k]?.batchTopicAssignment[l]?.assignedDate) {
							tempAssignmnetList.push(this.batchSessions[i]?.module[j]?.moduleTopics[k]?.batchTopicAssignment[l])
						}
					}
					this.batchSessions[i].module[j].moduleTopics[k].batchTopicAssignment = tempAssignmnetList
				}
			}
		}
		// console.log(this.batchSessionsPlan)
	}

	onAssignmentClick(module: any, moduleTopic: any, index: number): void {
		this.isViewClicked = false
		this.isViewAssignmentClick = false
		this.currentIndex = index
		this.module = module
		this.moduleTopic = moduleTopic
		this.openModal(this.assignmentModal, "xl")
	}

	calculateDeadlineDate(days: number): string {
		let date: string = this.batchSessions[this.currentIndex].date

		return date
	}

	onPreRequisiteClick(batchSession: IBatchSessionDetail): void {
		this.batchSession = batchSession
		this.createPreRequisiteForm()
		if (this.batchSession.batchSessionPrerequisite) {
			this.preRequisiteForm.patchValue(this.batchSession.batchSessionPrerequisite)
		}
		// console.log(this.preRequisiteForm);

		this.openModal(this.preRequisiteFormModal, "xl")
	}

	createPreRequisiteForm(): void {
		this.preRequisiteForm = this.formBuilder.group({
			id: new FormControl(null),
			prerequisite: new FormControl(null, [Validators.required, Validators.maxLength(100)]),
		})
	}


	addPreRequisite(): void {
		this.spinnerService.loadingMessage = "Adding Pre-Requisite"
		this.batchSessionService.addPreRequisite(this.batchID, this.batchSession.id, this.preRequisiteForm.value).subscribe((response: any) => {
			this.modalRef.close()
			this.getPreRequisite()
			alert(response.body)
		}, (error) => {
			if (error.error?.error) {
				alert(error.error?.error)
				return
			}
			alert(error.statusText)
		}).add(() => {

		})
	}

	validatePreRequisiteForm(): void {

		if (this.preRequisiteForm.invalid) {
			this.preRequisiteForm.markAllAsTouched()
			return
		}

		if (this.isOperationUpdate) {
			this.updatePreRequisite()
			return
		}

		this.addPreRequisite()
	}

	getPreRequisite(): void {
		this.spinnerService.loadingMessage = "Getting Pre-Requisite"

		let quaryParam: any = {
			batchSessionID: this.batchSession.id
		}
		this.batchSessionService.getPreRequisite(this.batchID, quaryParam).subscribe((response: any) => {
			this.batchSession.batchSessionPrerequisite = response.body[0]
		}, (error) => {
			if (error.error?.error) {
				alert(error.error?.error)
				return
			}
			alert(error.statusText)
		}).add(() => {

		})
	}

	updatePreRequisite(): void {
		this.spinnerService.loadingMessage = "Updating Pre-Requisite"

		this.batchSessionService.updatePreRequisite(this.batchID, this.batchSession.id, this.batchSession.batchSessionPrerequisite.id, this.preRequisiteForm.value).subscribe((response: any) => {
			this.modalRef.close()
			this.getPreRequisite()
			alert(response.body)
		}, (error) => {
			if (error.error?.error) {
				alert(error.error?.error)
				return
			}
			alert(error.statusText)
		}).add(() => {

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
			resourceID: new FormControl(null),
			resourceType: new FormControl(null),
			resourceSubType: new FormControl(null),
			fileType: new FormControl(null),
		}))

		let len = this.resourceFormArray.length
		this.getResourcesByType(len - 1)
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
	}

	async addResource(resourceID: any): Promise<void> {
		return new Promise<any>((resolve, reject) => {
			this.spinnerService.loadingMessage = "Adding Resource"

			let resource: any = {
				resourceID: resourceID
			}

			this.moduleService.addResource(this.moduleID, resource).subscribe((response: any) => {
				// this.module.resources.push(resource)
				resolve(response)
			}, (err: any) => {
				reject(err.error?.error)
			})
		})
	}

	async onSubmitResourceClick(): Promise<void> {
		this.validateResource()
		let errors: string[] = []
		let tempResource: IResource = this.multipleResourcesForm.value
		this.resource.push(tempResource)
		for (let index = 0; index < this.resourceFormArray.length; index++) {
			try {
				await this.addResource(this.resourceFormArray.at(index).get('resourceID').value)
			} catch (err) {
				errors.push(err)
			}
		}

		if (errors.length == 0) {
			alert("Resources successufully added")
			this.modalRef.close()
			return
		}
		alert(this.createErrorString(errors))
	}

	createErrorString(errors: string[]): string {
		let errorString = ""
		for (let index = 0; index < errors.length; index++) {
			errorString += (index == 0 ? "" : "\n") + errors[index]
		}
		return errorString
	}

	validateResource(): void {
		if (this.multipleResourcesForm.invalid) {
			this.multipleResourcesForm.markAllAsTouched()
			return
		}
	}

	openModal(content: any, size: string = "lg"): NgbModalRef {
		let options: NgbModalOptions = {
			ariaLabelledBy: 'modal-basic-title', keyboard: false,
			backdrop: 'static', size: size
		}

		return this.modalRef = this.modalService.open(content, options)
		// return this.modalRef
	}

	dismissFormModal(modal: NgbModalRef): void {
		modal.dismiss()
	}

	redirectToMeetLink(batchMeetLink: any): void {
		window.open(batchMeetLink, "_blank");
	}

	deleteNullValues(object: any): void {
		for (let field in object) {
			if (object[field] === null || object[field] === "" || field == "resourceID") {
				delete object[field];
			}
		}
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

	getFileType(): void {
		this.generalService.getGeneralTypeByType("file_type").subscribe((respond: any) => {
			this.fileTypeList = respond;
		}, (err) => {
			console.error(err)
		})
	}

}
