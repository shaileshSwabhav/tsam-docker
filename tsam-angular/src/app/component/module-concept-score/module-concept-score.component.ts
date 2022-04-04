import { Component, ElementRef, OnInit, QueryList, ViewChildren } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';;
import { Observable } from 'rxjs';
import { delay, map } from 'rxjs/operators';
import { BatchService } from 'src/app/service/batch/batch.service';
import { WorkflowItemScore, ConceptModuleService } from 'src/app/service/concept-module/concept-module.service';
import { AccessLevel, Role } from 'src/app/service/constant';
import { DagManagerService } from 'src/app/service/dag-manager/dag-manager.service';
import { LocalService } from 'src/app/service/storage/local.service';

declare let LeaderLine: any

@Component({
  selector: 'app-module-concept-score',
  templateUrl: './module-concept-score.component.html',
  styleUrls: ['./module-concept-score.component.css'],
  providers: [DagManagerService]
})
export class ModuleConceptScoreComponent implements OnInit {

  // Course.
  courseID: string

  // Concept module.
  conceptModuleList: any[]

  // Batch.
  batchID: string
  batchName: string

  // Batch module.
  batchModuleList: any[]
  batchModuleTabList: any[]
  selectedBatchModule: any

  // Batch talent.
  batchTalentList: any[]

  // Talent.
  talentID: string

  // Node.
  @ViewChildren('nodeList', { read: ElementRef }) nodeList: QueryList<ElementRef>
  public nodeListCount$: Observable<number>

  // Workflow.
  workflowItems: WorkflowItemScore[]

  public workflow$: Observable<WorkflowItemScore[][]>

  // Line.
  private linesArray = []

  // Login.
  loginID: string
  roleName: string

  // Access.
  access: any

  constructor(
    private _dagManagerService: DagManagerService<WorkflowItemScore>,
    private route: ActivatedRoute,
    private spinnerService: SpinnerService,
    private conceptModuleService: ConceptModuleService,
    private router: Router,
    private batchService: BatchService,
    private localService: LocalService,
    private role: Role,
    private accessLevel: AccessLevel,
  ) { }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {

    // Spinner.
    this.spinnerService.loadingMessage = "Getting Concept Tree"


    // Concept module.
    this.conceptModuleList = []

    // Module.
    this.route.queryParamMap.subscribe(params => {
      this.batchID = params.get("batchID")
      this.courseID = params.get("courseID")
      this.batchName = params.get("batchName")
      if (!this.batchID) {
        return
      }
    }, err => {
      console.error(err)
    })

    // Workflow.
    this.workflowItems = []

    // Batch module.
    this.batchModuleList = []
    this.batchModuleTabList = []

    // Talent.
    this.loginID = this.localService.getJsonValue("loginID")
    if (this.role.TALENT == this.localService.getJsonValue("roleName")) {
      this.access = this.accessLevel.ONLY_TALENT
      this.talentID = this.loginID
    }
    if (this.role.FACULTY == this.localService.getJsonValue("roleName")) {
      this.access = this.accessLevel.ONLY_FACULTY
    }
    if (this.role.ADMIN == this.localService.getJsonValue("roleName") || this.role.SALES_PERSON == this.localService.getJsonValue("roleName")) {
      this.access = this.accessLevel.ADMIN_AND_SALESPERSON
    }

    // Batch talent.
    this.batchTalentList = []

    // Get components.
    this.getBatchModuleList()
  }

  ngOnDestroy() {
    this.removeLines()
  }

  ngAfterViewInit() {
    this.nodeListCount$ = this.nodeList.changes.pipe(
      delay(0),
      map((list: QueryList<ElementRef>) => list.toArray().length)
    )
    this.nodeList.changes.subscribe((list) => {
      this.removeLines()
      this.drawLines()
    })
    this.drawLines()
  }

  // ************************************************** LINE ************************************************

  // Remove and draw lines.
  removeAndDrawLines(): void {
    setTimeout(() => {
      this.removeLines()
      this.drawLines()
    }, 0)
  }

  // Remove all lines between nodes.
  removeLines() {
    this.linesArray.forEach((line) => line.remove())
    this.linesArray = []
  }

  // Draw all lines between nodes.
  drawLines() {
    const nodeListAs1DArray = this._dagManagerService.getSingleDimensionalArrayFromModel()
    nodeListAs1DArray.forEach((box: WorkflowItemScore) => {
      box.parentIds.forEach((parentId: number) => {
        if (parentId > 0) {
          const parent: ElementRef<any> = this.nodeList.find(
            (b: ElementRef) => parentId === +b.nativeElement.children[0].id
          )
          const self: ElementRef<any> = this.nodeList.find(
            (b: ElementRef) => +b.nativeElement.children[0].id === box.stepId
          )
          if (parent && self) {
            const line = new LeaderLine(
              parent.nativeElement,
              self.nativeElement,
              {
                startSocket: 'bottom',
                endSocket: 'top',
                endPlug: 'arrow1',
                color: 'rgba(0, 0, 0, 0.3)',
                path: 'straight',
                size: 3,
                // dropShadow: true,
                // outlineColor: 'white',
                // outline: true,
              }
            )
            this.linesArray.push(line)
          }
        }
      })
    })
  }

  // On scrolling concept tree.
  onConceptTreeScroll(event: any): void {
    this.linesArray.forEach((line) => line.position())
  }

  // ************************************************** NODE ************************************************

  // Create work flow items.
  createWorkflowItems(): void {
    this.workflowItems = []
    if (this.conceptModuleList.length > 0) {
      this.convertConceptModuleListIntoWorkflowItems()
      this._dagManagerService.setNextNumber(this.conceptModuleList.length + 1)
      this._dagManagerService.setNewItemsArrayAsDagModel(this.workflowItems)
      this.workflow$ = this._dagManagerService.dagModel$
      this.workflow$.subscribe((response) => {
      })
    }
  }

  // ************************************************** CONCEPT ************************************************

  // Get step id of node by concept id.
  getStepIDByConceptID(conceptID: string): number {
    for (let i = 0; i < this.workflowItems.length; i++) {
      if (this.workflowItems[i].conceptID == conceptID) {
        return this.workflowItems[i].stepId
      }
    }
    return -1
  }

  // Convert concept module list into workflow items. 
  convertConceptModuleListIntoWorkflowItems(): void {
    let stepId: number = 0
    let level: number = 0
    let branchPath: number = 1
    this.initializeStepId()
    for (let i = 0; i < this.conceptModuleList.length; i++) {

      // Set step id.
      stepId = stepId + 1

      // Set branch path id. 
      if (level == this.conceptModuleList[i].level) {
        branchPath = branchPath + 1
      }
      if (level < this.conceptModuleList[i].level) {
        level = this.conceptModuleList[i].level
        branchPath = 1
      }

      // Set parent step ids.
      let parentIds: number[] = []
      let parentConceptIDs: string[] = []
      if (this.conceptModuleList[i].parentModuleProgrammingConcepts.length == 0) {
        parentIds.push(0)
      }
      if (this.conceptModuleList[i].parentModuleProgrammingConcepts.length > 0) {
        for (let j = 0; j < this.conceptModuleList[i].parentModuleProgrammingConcepts.length; j++) {
          for (let k = 0; k < this.workflowItems.length; k++) {
            if (this.conceptModuleList[i].parentModuleProgrammingConcepts[j].programmingConceptID == this.workflowItems[k].conceptID) {
              parentIds.push(this.workflowItems[k].stepId)
              parentConceptIDs.push(this.workflowItems[k].conceptID)
            }
          }
        }
      }
      parentIds.sort((a, b) => b - a)

      // Create work flow item.
      this.workflowItems.push({
        id: this.conceptModuleList[i].id,
        conceptID: this.conceptModuleList[i].programmingConcept.id,
        stepId: stepId,
        parentIds: parentIds,
        branchPath: branchPath,
        score: this.conceptModuleList[i].averageScore,
        conceptName: this.conceptModuleList[i].programmingConcept.name,
        showParentConcepts: false,
        complexityName: this.conceptModuleList[i].programmingConcept.complexityName
      })
    }

    for (let k = 0; k < this.workflowItems.length; k++) {
      let parentExists: boolean = false
      let stepId: number = this.workflowItems[k].stepId
      for (let i = 0; i < this.workflowItems.length; i++) {
        for (let j = 0; j < this.workflowItems[i].parentIds.length; j++) {
          if (this.workflowItems[i].parentIds[j] == stepId) {
            parentExists = true
          }
        }
      }
      if (parentExists) {
        this.workflowItems[k].isConceptIDDisabled = true
      }
    }
  }

  // Initialize step id.
  initializeStepId(): void {
    this._dagManagerService.setNextNumber(1)
  }

  // ************************************************** OTHER FUNCTIONS ************************************************

  // Redirect to batch details.
  redirectToBtchDetails(): void {
    let url: string
    let queryParams: any
    if (this.access.isAdmin || this.access.isSalesPerson) {
      url = "/training/batch/master/session/details"
      queryParams = {
        "batchID": this.batchID,
        "tab": 2,
        "subTab": "Create/Update",
        "courseID": this.courseID,
        "batchName": this.batchName
      }
    }
    if (this.access.isFaculty) {
      url = "/my/batch/session/details"
      queryParams = {
        "batchID": this.batchID,
        "tab": 2,
        "subTab": "Create/Update",
        "courseID": this.courseID,
        "batchName": this.batchName
      }
    }
    if (this.access.isTalent) {
      url = "/my-batches"
      queryParams = {
        "batchID": this.batchID,
        "tab": "My Class",
        "subTab": "Batch Details",
      }
    }

    this.router.navigate([url], { queryParams }).catch(err => {
      console.error(err)
    })
  }

  // Format the batch modules tab list.
  formatBatchModuleTabList(): void {

    // Create course module tab list.
    for (let i = 0; i < this.batchModuleList.length; i++) {
      this.batchModuleTabList.push(
        {
          moduleName: this.batchModuleList[i].module.moduleName,
          module: this.batchModuleList[i].module,
        }
      )

      if (this.batchModuleTabList[i].module.logo) {
        this.batchModuleTabList[i].imageURL = this.batchModuleList[i].module.logo
      }
      else {
        this.batchModuleTabList[i].imageURL = "assets/icon/grey-icons/Score.png"
      }
    }
    this.selectedBatchModule = this.batchModuleList[0]?.module
  }

  // On clicking module tab.
  onModuleTabClick(event: any): void {
    for (let i = 0; i < this.batchModuleList.length; i++) {
      if (i == event.index) {
        this.selectedBatchModule = this.batchModuleList[i].module
      }
    }
    this.getConceptModuleList(this.talentID)
  }

  // Format concept list.
  formatConceptList(): void {
    for (let i = 0; i < this.conceptModuleList.length; i++) {
      if (this.conceptModuleList[i].programmingConcept.complexity == 1) {
        this.conceptModuleList[i].programmingConcept.complexityName = "Easy"
      }
      if (this.conceptModuleList[i].programmingConcept.complexity == 2) {
        this.conceptModuleList[i].programmingConcept.complexityName = "Medium"
      }
      if (this.conceptModuleList[i].programmingConcept.complexity == 3) {
        this.conceptModuleList[i].programmingConcept.complexityName = "Hard"
      }
    }
  }

  // On talent click.
  onTalentClick(talentID: string): void {
    this.talentID = talentID
    this.getConceptModuleList(this.talentID)
  }

  // ************************************************** GET FUNCTIONS ************************************************

  // Get concept modules for module.
  getConceptModuleList(talentID: string): void {
    this.spinnerService.loadingMessage = "Getting Concept Tree"
    let queryParams: any = {
      limit: -1,
      offset: 0,
      moduleID: this.selectedBatchModule.id
    }
    this.conceptModuleService.getAllConceptModulesForTalentScore(talentID, queryParams).subscribe((response: any) => {
      this.conceptModuleList = response.body
      this.formatConceptList()
      this.createWorkflowItems()
    }, (error) => {
      console.error(error)
      if (error.error?.error) {
        alert(error.error?.error)
        return
      }
      alert(error.statusText)
    })
  }

  // Get batch module list.
  getBatchModuleList(): void {
    this.spinnerService.loadingMessage = "Getting Concept Tree"
    let queryParams: any = {
      limit: -1,
      offset: 0,
    }
    this.batchService.getBatchModules(this.batchID, queryParams).subscribe((response) => {
      this.batchModuleList = response.body
      this.formatBatchModuleTabList()
      if (this.access.isTalent && this.batchModuleList.length > 0) {
        this.getConceptModuleList(this.loginID)
      }
      if (!this.access.isTalent && this.batchModuleList.length > 0) {
        this.getTalentListOfBatch()
      }
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

  // Get talent list of batch.
  getTalentListOfBatch(): void {
    this.spinnerService.loadingMessage = "Getting Concept Tree"
    this.batchService.getBatchTalentList(this.batchID).subscribe((response) => {
      this.batchTalentList = response.body
      if (this.batchTalentList.length > 0) {
        this.talentID = this.batchTalentList[0].id
        this.getConceptModuleList(this.talentID)
      }
    }, error => {
      console.error(error)
      if (error.statusText.includes('Unknown')) {
        alert("No connection to server. Check internet.")
      }
    })
  }

}
