import { Component, ElementRef, OnInit, QueryList, ViewChildren } from '@angular/core';
import { DagManagerService } from 'src/app/service/dag-manager/dag-manager.service';
import { Observable } from 'rxjs';
import { delay, map } from 'rxjs/operators';
import { ConceptModuleService, WorkflowItem } from 'src/app/service/concept-module/concept-module.service';
import { ActivatedRoute, Router } from '@angular/router';
import { ProgrammingConceptService } from 'src/app/service/programming-concept/programming-concept.service';
import { SpinnerService } from 'src/app/service/spinner/spinner.service';
import { LocalService } from 'src/app/service/storage/local.service';
import { Role } from 'src/app/service/constant';

declare let LeaderLine: any

@Component({
  selector: 'app-module-concept',
  templateUrl: './module-concept.component.html',
  styleUrls: ['./module-concept.component.css'],
  providers: [DagManagerService]
})
export class ModuleConceptComponent implements OnInit {

  // Concept.
  conceptList: any[]
  parentConceptList: any[]
  selectedConceptList: string[]

  // Module
  moduleID: string
  moduleName: string

  // Concept module.
  conceptModuleList: any[]

  // Node.
  @ViewChildren('nodeList', { read: ElementRef }) nodeList: QueryList<ElementRef>
  @ViewChildren('nodeComponentList') nodeComponentList: QueryList<any>
  public nodeListCount$: Observable<number>

  // Workflow.
  workflowItems: WorkflowItem[]

  public workflow$: Observable<WorkflowItem[][]>

  // Flags.
  isOperationUpdate: boolean
  conceptTreeScroll: boolean
  isConceptTreeDirty: boolean
  isFaculty: boolean

  // Line.
  private linesArray = []

  constructor(
    private _dagManagerService: DagManagerService<WorkflowItem>,
    private route: ActivatedRoute,
    private conceptService: ProgrammingConceptService,
    private spinnerService: SpinnerService,
    private conceptModuleService: ConceptModuleService,
    private router: Router,
    private localService: LocalService,
		private role: Role,
  ) { }

  get ongoingOperations() {
    return this.spinnerService.ongoingOperations
  }

  ngOnInit(): void {

    // Concept.
    this.conceptList = []
    this.parentConceptList = []
    this.selectedConceptList = []

    // Spinner.
    this.spinnerService.loadingMessage = "Creating Concept Tree"

    // Concept module.
    this.conceptModuleList = []

    // Flags.
    this.isOperationUpdate = false
    this.conceptTreeScroll = false
    this.isConceptTreeDirty = false
		this.isFaculty = (this.localService.getJsonValue("roleName") == this.role.FACULTY ? true : false)

    // Module.
    this.route.queryParamMap.subscribe(params => {
      this.moduleID = params.get("moduleID")
      this.moduleName = params.get("moduleName")
      if (!this.moduleID && !this.moduleName) {
        return
      }
    }, err => {
      console.error(err)
    })

    // Workflow.
    this.workflowItems = []

    // Get components.
    this.getConceptList()
    this.getConceptModuleList()
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
      this.conceptChangeFromParent()
    })
    this.drawLines()
    this.conceptChangeFromParent()
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
    nodeListAs1DArray.forEach((box: WorkflowItem) => {
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
    // if (this.conceptTreeScroll == false){
    //   this.repositionLines()
    // }
    this.linesArray.forEach((line) => line.position())

    // this.conceptTreeScroll = true
    // this.repositionLines()
  }

  // Reposition lines.
  repositionLines(): void {
    setInterval(() => {
      if (this.conceptTreeScroll) {
        this.conceptTreeScroll = false;
        this.linesArray.forEach((line) => line.position())
        clearInterval()
      }
    }, 1000)
  }

  // ************************************************** NODE ************************************************

  // Create work flow items.
  createWorkflowItems(): void {

    this.workflowItems = []

    // If there is no concept module list then initialize work flow items.
    if (this.conceptModuleList.length == 0) {
      this.isOperationUpdate = false
      this.workflowItems.push({
        conceptID: null,
        id: null,
        stepId: 1,
        parentIds: [0],
        branchPath: 1,
        toBeRemoved: false,
        parentConceptIDs: [],
        isConceptIDDisabled: false,
        showParentConcepts: false
      })
      this._dagManagerService.setNextNumber(2)
    }

    // If there is concept module list then create work flow items.
    if (this.conceptModuleList.length > 0) {
      this.isOperationUpdate = true
      this.convertConceptModuleListIntoWorkflowItems()
      this._dagManagerService.setNextNumber(this.conceptModuleList.length + 1)
    }

    this._dagManagerService.setNewItemsArrayAsDagModel(this.workflowItems)
    this.workflow$ = this._dagManagerService.dagModel$
    this.workflow$.subscribe((response) => {
      // console.log(response)
    })
  }

  // Add starting node.
  addStartingNode(): void {
    this.initializeStepId()
    this._dagManagerService.addNewStep([0], 1, 1,
      {
        conceptID: null,
        toBeRemoved: false,
        parentConceptIDs: [],
        isConceptIDDisabled: false,
        id: null,
        showParentConcepts: false
      })

    // Get all the nodes.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()
  }

  // Add child node.
  addChildNode({
    parentIds,
    numberOfChildrenToBeAdded,
    parentConceptID,
  }: {
    parentIds: number[]
    numberOfChildrenToBeAdded: number
    parentConceptID: string
  }) {

    // Count the number of children of node.
    let childrenCount: number = this._dagManagerService.nodeChildrenCount(parentIds[0])

    // Get all the nodes.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()

    // Set the next step id.
    this.setNextStepId()

    // If node has children then add new child in same level. 
    if (childrenCount > 0) {
      this._dagManagerService.addNewStepAsNewPath(parentIds[0], numberOfChildrenToBeAdded,
        {
          conceptID: null,
          toBeRemoved: false,
          parentConceptIDs: [parentConceptID],
          isConceptIDDisabled: false,
          id: null,
          showParentConcepts: false
        })
    }

    // If node has no children then add new first child. 
    if (childrenCount == 0) {
      this._dagManagerService.addNewStep(parentIds, numberOfChildrenToBeAdded, 1,
        {
          conceptID: null,
          toBeRemoved: false,
          parentConceptIDs: [parentConceptID],
          isConceptIDDisabled: false,
          id: null,
          showParentConcepts: false
        })
    }

    // Make concept id of node(to whom child node is being added) disabled.
    this.nodeComponentList.forEach((tempNode) => {
      if (tempNode.node.stepId == parentIds[0]) {
        tempNode.node.isConceptIDDisabled = true
      }
    })

    // Get all the nodes.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()

    // // Hide all nodes parent concepts.
    // this.hideAllNodesParentConcepts()

    // Make is concept tree dirty.
    this.isConceptTreeDirty = true
  }

  // Remove node.
  removeNode({ stepId }: { stepId: number }) {

    // Get all the nodes.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()

    // Remove all nodes with no parents.
    this.removeNodeRecursively(stepId)

    // If the parent node has no children then enable its concept id.
    for (let K = 0; K < this.workflowItems.length; K++) {
      let parentExists: boolean = false
      let stepId: number = this.workflowItems[K].stepId
      for (let i = 0; i < this.workflowItems.length; i++) {
        for (let j = 0; j < this.workflowItems[i].parentIds.length; j++) {
          if (this.workflowItems[i].parentIds[j] == stepId) {
            parentExists = true
          }
        }
        this.nodeComponentList.forEach((tempNode) => {
          if (tempNode.node.stepId == stepId && parentExists) {
            tempNode.node.isConceptIDDisabled = true
          }
          if (tempNode.node.stepId == stepId && !parentExists) {
            tempNode.node.isConceptIDDisabled = false
          }
        })
      }
    }

    // // On changing concept of node.
    // this.conceptChange()

    // Update the dag model.
    this._dagManagerService.setNewItemsArrayAsDagModel(this.workflowItems)

    // Get all the nodes.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()

    // // Hide all nodes parent concepts.
    // this.hideAllNodesParentConcepts()

    // Make is concept tree dirty.
    this.isConceptTreeDirty = true
  }

  // Remove node recuresively.
  removeNodeRecursively(stepId: number): void {
    let toBeRemovedNodeList: WorkflowItem[] = []

    // Remove occurrence of the step id as parent id in all nodes.
    for (let i = 0; i < this.workflowItems.length; i++) {
      if (this.workflowItems[i].stepId == stepId) {
        this.workflowItems[i].toBeRemoved = true
      }
      if (this.workflowItems[i].parentIds) {
        for (let j = 0; j < this.workflowItems[i].parentIds.length; j++) {
          if (this.workflowItems[i].parentIds[j] == stepId) {
            this.workflowItems[i].parentIds.splice(j, 1)
          }
        }
      }
      if (this.workflowItems[i].parentIds.length == 0) {
        this.workflowItems[i].toBeRemoved = true
        toBeRemovedNodeList.push(this.workflowItems[i])
      }
    }

    // Only keep those nodes whose toBeRemoved is false.
    let tempWorkflowItems: WorkflowItem[] = []
    for (let i = 0; i < this.workflowItems.length; i++) {
      if (!this.workflowItems[i].toBeRemoved) {
        tempWorkflowItems.push(this.workflowItems[i])
      }
    }
    this.workflowItems = tempWorkflowItems

    if (toBeRemovedNodeList.length > 0) {
      for (let i = 0; i < toBeRemovedNodeList.length; i++) {
        this.removeNodeRecursively(toBeRemovedNodeList[i].stepId)
      }
    }
  }

  // Add sibling root node.
  addSiblingRootNode() {

    // Get all the nodes.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()

    // Get the max branchpath of the first level.
    let tempDagModel: WorkflowItem[][] = this._dagManagerService.getCurrentDagModel()
    let maxBranchPath: number = 0
    for (let i = 0; i < tempDagModel[0].length; i++) {
      if (maxBranchPath < tempDagModel[0][i].branchPath) {
        maxBranchPath = tempDagModel[0][i].branchPath
      }
    }

    // Add sibling root node.
    this.workflowItems.push(
      {
        conceptID: null,
        id: null,
        stepId: this.getNextStepId(this.workflowItems),
        parentIds: [0],
        branchPath: (maxBranchPath + 1),
        toBeRemoved: false,
        parentConceptIDs: [],
        isConceptIDDisabled: false,
        showParentConcepts: false
      })

    // Update the dag model.
    this._dagManagerService.setNewItemsArrayAsDagModel(this.workflowItems)

    // Get all the nodes.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()

    // // Hide all nodes parent concepts.
    // this.hideAllNodesParentConcepts()

    // Make is concept tree dirty.
    this.isConceptTreeDirty = true
  }

  // Add parent to node.
  addParentToNode(event: any): void {

    // Get step id of node by concept id.
    let parentStepID: number = this.getStepIDByConceptID(event.parentConceptID)

    // Add parent concept node to concept node.
    try {
      this._dagManagerService.addNewRelation(event.stepId, parentStepID)
    } catch (err) {
      alert("Not possible")
      // console.error(err)

      // Call child function to remove id from its parent concept ids.
      this.nodeComponentList.forEach((tempNode) => {
        if (tempNode.node.stepId == event.stepId) {
          tempNode.removeParentConceptIDIfError(event.parentConceptID)
        }
      })
    }

    // Disable the concept id of the parent node.
    this.nodeComponentList.forEach((tempNode) => {
      if (tempNode.node.stepId == parentStepID) {
        tempNode.node.isConceptIDDisabled = true
      }
    })

    // Get all the nodes.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()

    // // Hide all nodes parent concepts.
    // this.hideAllNodesParentConcepts()

    // Make is concept tree dirty.
    this.isConceptTreeDirty = true
  }

  // Remove parent from node.
  removeParentFromNode(event: any) {

    // Get step id of node by concept id.
    let parentStepID: number = this.getStepIDByConceptID(event.parentConceptID)

    // Remove parent from node.
    for (let i = 0; i < this.workflowItems.length; i++) {
      if (this.workflowItems[i].stepId == event.stepId) {
        for (let j = 0; j < this.workflowItems[i].parentIds.length; j++) {
          if (this.workflowItems[i].parentIds[j] == parentStepID) {
            this.workflowItems[i].parentIds.splice(j, 1)
          }
        }
      }
    }

    // If the parent node has no children then enable its concept id.
    let parentExists: boolean = false
    for (let i = 0; i < this.workflowItems.length; i++) {
      if (this.workflowItems[i].stepId == parentStepID) {
        for (let j = 0; j < this.workflowItems[i].parentIds.length; j++) {
          if (this.workflowItems[i].parentIds[j] == parentStepID) {
            parentExists = true
          }
        }
      }
    }
    if (!parentExists) {
      this.nodeComponentList.forEach((tempNode) => {
        if (tempNode.node.stepId == parentStepID) {
          tempNode.node.isConceptIDDisabled = false
        }
      })
    }

    // Update the dag model.
    this._dagManagerService.setNewItemsArrayAsDagModel(this.workflowItems)

    // Get all the nodes.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()

    // // Hide all nodes parent concepts.
    // this.hideAllNodesParentConcepts()

    // Make is concept tree dirty.
    this.isConceptTreeDirty = true
  }

  // Get next step id.
  getNextStepId(arr: WorkflowItem[]) {
    return (
      Math.max.apply(
        Math,
        arr.map((i) => i.stepId)
      ) + 1
    )
  }

  // Set next step id while adding node.
  setNextStepId(): void {
    this._dagManagerService.setNextNumber(this.getNextStepId(this.workflowItems))
  }

  // Initialize step id.
  initializeStepId(): void {
    this._dagManagerService.setNextNumber(1)
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

  // On changing concept of node.
  conceptChange(): void {

    // Iterate concept list.
    for (let i = 0; i < this.conceptList.length; i++) {

      // Make concept enabled.
      this.conceptList[i].disabled = false

      // Iterate work flow items.
      for (let j = 0; j < this.workflowItems.length; j++) {
        if (this.workflowItems[j].conceptID == this.conceptList[i].id) {
          this.conceptList[i].disabled = true
        }
      }
    }

    // Send the updated concept list to all nodes when there is no addition or removal of nodes.
    this.nodeComponentList.forEach((tempNode) => {
      tempNode.setConceptList(this.conceptList)
    })

    // Send the updated concept list to all nodes when there is addition or removal of nodes.
    this.nodeComponentList.changes.subscribe((list) => {
      this.nodeComponentList.forEach((tempNode) => {
        tempNode.setConceptList(this.conceptList)
      })
    })
  }

  // On changing concept of node from parent while loading the nodes by getting concept modules.
  conceptChangeFromParent(): void {
    this.conceptChange()
  }

  // On changing concept of node from child when child node's concept changes.
  conceptChangeFromChild(): void {
    this.conceptChange()

    // Make is concept tree dirty.
    this.isConceptTreeDirty = true
  }

  // Hide all nodes' parent concepts.
  hideAllNodesParentConcepts(): void {
    this.nodeComponentList.forEach((tempNode) => {
      tempNode.node.showParentConcepts = false
    })
  }

  // On clicking save button.
  onSaveClick(): void {
    let doesAllNodesHaveConecpts: boolean = true
    this.nodeComponentList.forEach((tempNode) => {
      if (tempNode.node.conceptID == null) {
        doesAllNodesHaveConecpts = false
      }
    })
    if (!doesAllNodesHaveConecpts) {
      alert("Please provide concept for all nodes")
      return
    }

    // Add/Update concept modules. 
    this.formatConceptModuleForAddUpdate()
  }

  // Format the nodes into concept modules for update and add.
  formatConceptModuleForAddUpdate(): void {

    // Get 2D array to give level to all nodes.
    let workflowItems2DArray: WorkflowItem[][] = this._dagManagerService.getCurrentDagModel()
    for (let i = 0; i < workflowItems2DArray.length; i++) {
      for (let j = 0; j < workflowItems2DArray[i].length; j++) {
        workflowItems2DArray[i][j]["level"] = (i + 1)
      }
    }

    // Get all the nodes in 1D array.
    this.workflowItems = this._dagManagerService.getSingleDimensionalArrayFromModel()

    // Convert the nodes into concept modules.
    this.conceptModuleList = []
    for (let i = 0; i < this.workflowItems.length; i++) {
      let conceptModule: any = {}
      conceptModule.programmingConceptID = this.workflowItems[i].conceptID
      conceptModule.id = this.workflowItems[i].id
      conceptModule.moduleID = this.moduleID
      conceptModule.level = this.workflowItems[i]["level"]
      conceptModule.parentConceptIDs = this.workflowItems[i].parentConceptIDs
      conceptModule.parentModuleProgrammingConcepts = []

      // If level is greater than 1 then give parent concpet module to all concept modules.
      if (conceptModule.level > 1) {
        for (let j = 0; j < conceptModule.parentConceptIDs.length; j++) {
          for (let k = 0; k < this.conceptModuleList.length; k++) {
            if (conceptModule.parentConceptIDs[j] == this.conceptModuleList[k].programmingConceptID) {
              let tempParentConceptModule: any = {}
              tempParentConceptModule.programmingConceptID = this.conceptModuleList[k].programmingConceptID
              tempParentConceptModule.moduleID = this.moduleID
              tempParentConceptModule.level = this.conceptModuleList[k]["level"]
              conceptModule.parentModuleProgrammingConcepts.push(tempParentConceptModule)
            }
          }
        }
      }
      this.conceptModuleList.push(conceptModule)
    }
    // console.log(this.conceptModuleList)

    // Update concept modules. 
    if (this.isOperationUpdate) {
      this.updateConceptModules()
    }

    // Add concept modules. 
    if (!this.isOperationUpdate) {
      this.addConceptModules()
    }
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
      if (this.conceptModuleList[i].parentModuleProgrammingConcepts?.length == 0) {
        parentIds.push(0)
      }
      if (this.conceptModuleList[i].parentModuleProgrammingConcepts?.length > 0) {
        for (let j = 0; j < this.conceptModuleList[i].parentModuleProgrammingConcepts?.length; j++) {
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
        conceptID: this.conceptModuleList[i].programmingConceptID,
        stepId: stepId,
        parentIds: parentIds,
        branchPath: branchPath,
        toBeRemoved: false,
        parentConceptIDs: parentConceptIDs,
        isConceptIDDisabled: false,
        showParentConcepts: false
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
    // console.log(this.workflowItems)
  }

  // Add new concept modules.
  addConceptModules(): void {
    this.spinnerService.loadingMessage = "Adding Concepts"
    this.conceptModuleService.addConceptModules(this.conceptModuleList)
      .subscribe((response: any) => {
        this.getConceptModuleList()
        alert(response)
      }, (error) => {
        console.error(error)
        if (error.error?.error) {
          alert(error.error?.error)
          return
        }
        alert(error.statusText)
      })
  }

  // Update new concept modules.
  updateConceptModules(): void {
    this.spinnerService.loadingMessage = "Updating Concepts"
    this.conceptModuleService.updateConceptModules(this.moduleID, this.conceptModuleList)
      .subscribe((response: any) => {
        this.getConceptModuleList()
        alert(response)
      }, (error) => {
        console.error(error)
        if (error.error?.error) {
          alert(error.error?.error)
          return
        }
        alert(error.statusText)
      })
  }

  // Get concept modules for module.
  getConceptModuleList(): void {
    this.spinnerService.loadingMessage = "Getting modules"
    let queryParams: any = {
      limit: -1,
      offset: 0,
      moduleID: this.moduleID
    }
    this.conceptModuleService.GetAllModuleProgrammingConceptsForConceptTree(queryParams).subscribe((response: any) => {
      this.conceptModuleList = response.body
      this.isConceptTreeDirty = false
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

  // Format concept list.
  formatConceptList(): void {
    for (let i = 0; i < this.conceptList.length; i++) {
      if (this.conceptList[i].complexity == 1) {
        this.conceptList[i].complexityName = "Easy"
      }
      if (this.conceptList[i].complexity == 2) {
        this.conceptList[i].complexityName = "Medium"
      }
      if (this.conceptList[i].complexity == 3) {
        this.conceptList[i].complexityName = "Hard"
      }
    }
  }

  // Delete whole concept tree.
  deleteConceptTree(): void {
    this.workflowItems = []
    this._dagManagerService.setNewItemsArrayAsDagModel(this.workflowItems)
    this.conceptModuleList = []
    this.updateConceptModules()
  }

  // ************************************************** OTHER FUNCTIONS ************************************************

  // Redirect to modules page.
  redirectToModules(): void {
    let url: string
		if (!this.isFaculty){
			url = "/training/module"
		}
		if (this.isFaculty){
			url = "/bank/module"
		}
    this.router.navigate([url], {
    }).catch(err => {
      console.error(err)
    })
  }

  // ************************************************** GET FUNCTIONS ************************************************

  // Get concept list.
  getConceptList(): void {
    let queryParams: any = {
      limit: -1,
      offset: 0
    }
    this.conceptService.getAllConcepts(queryParams).subscribe((response: any) => {
      this.conceptList = response.body
      this.formatConceptList()
      this.parentConceptList = JSON.parse(JSON.stringify(this.conceptList))
    }, (err: any) => {
      console.error(err)
    })
  }

}

