import { ChangeDetectorRef, Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { NgSelectComponent } from '@ng-select/ng-select';
import { WorkflowItem } from 'src/app/service/concept-module/concept-module.service';

@Component({
  selector: 'app-module-concept-node',
  templateUrl: './module-concept-node.component.html',
  styleUrls: ['./module-concept-node.component.css']
})
export class ModuleConceptNodeComponent {

  // Node.
  @Input() node: WorkflowItem

  // Concept.
  @Input() conceptList: any[]
  @Input() parentConceptList: any[]
  @ViewChild('ngSelectParentConceptComponent') ngSelectParentConceptComponent: NgSelectComponent;
  @ViewChild('ngSelectConceptComponent') ngSelectConceptComponent: NgSelectComponent;

  // Add node emitter.
  @Output() addNodeEmitter: EventEmitter<{
    parentIds: number[]
    numberOfChildrenToBeAdded: number
    parentConceptID: string
  }> = new EventEmitter<{ parentIds: number[]; numberOfChildrenToBeAdded: number; parentConceptID: string }>()

  // Add new parent to node emitter.
  @Output() addNewParentToNodeEmitter: EventEmitter<any> = new EventEmitter()

  // Remove node emitter.
  @Output() removeNodeEmitter: EventEmitter<{ stepId: number }> = new EventEmitter<{stepId: number }>()

  // Change concept emitter.
  @Output() changeConceptEmitter: EventEmitter<any> = new EventEmitter<any>()

  // Remove parent concept emitter.
  @Output() removeParentConceptEmitter: EventEmitter<any> = new EventEmitter<any>()

  // Redraw lines fo parent component.
  @Output() redrawParentLines: EventEmitter<any> = new EventEmitter<any>()

  constructor(
		private cdr: ChangeDetectorRef
  ) {
  }

  // ngAfterViewInit(){
  //   console.log("ng after nide")
  //   // this.cdr.detectChanges()
  // }

  // nodeDetectChanges(): void{
  //   console.log(this.node.cdr)
  //   this.cdr.detectChanges()
  // }

  // Add child node.
  addChildNode() {
    
    // If concept id is not given then cannot add child.
    if (this.node.conceptID == null){
      alert("Please provide concept first")
      return
    }
    this.addNodeEmitter.emit({
      parentIds: [this.node.stepId],
      numberOfChildrenToBeAdded: 1,
      parentConceptID: this.node.conceptID
    })
  }

  // Add new parent to node.
  tempAddNewParentToNode(stepId: number, parentIdInString: string) {
    let parentId: number = parseInt(parentIdInString)
    this.addNewParentToNodeEmitter.emit({stepId, parentId})
  }

  // Add new parent to node.
  addNewParentToNode(parentConceptID: string) {
    let stepId: number = this.node.stepId
    this.addNewParentToNodeEmitter.emit({stepId, parentConceptID})
  }

  // Remove node.
  removeNode() {
    this.removeNodeEmitter.emit({
      stepId: this.node.stepId,
    })
  }

  // On cliking show parent section.
  onShowParentSectionClick(): void{
    this.node.showParentSection = !this.node.showParentSection
  }

  // On changing concept.
  onConceptChange(): void{
    this.changeConceptEmitter.emit()
  } 

  // Set concept list.
  setConceptList(tempConceptList: any[]): void{

    // Update the concept list.
    this.conceptList = [...tempConceptList]

    // this.cdr.detectChanges()
  }

  // Set concept.
  setConcept(conceptID: string): void{

    // Set the concept.
    this.node.conceptID = conceptID
  }

  // On adding parent concept id.
  onAddingParentConceptID(parentConcept: any){

    // If concept id is null then cannot give parent concept id.
    if (this.node.conceptID == null){
      let tempParentConceptsIDs: string[] = this.node.parentConceptIDs
      let index: number = tempParentConceptsIDs.indexOf(parentConcept.id)
      if (index !== -1) {
        tempParentConceptsIDs.splice(index, 1)
      }
      this.node.parentConceptIDs = [...tempParentConceptsIDs]
      alert("Please provide concept first")
      return
    }

    // If concept id and parent concept id is same then return.
    if (this.node.parentConceptIDs.includes(this.node.conceptID)){
      let tempParentConceptsIDs: string[] = this.node.parentConceptIDs
      let index: number = tempParentConceptsIDs.indexOf(this.node.conceptID)
      if (index !== -1) {
        tempParentConceptsIDs.splice(index, 1)
      }
      this.node.parentConceptIDs = [...tempParentConceptsIDs]
      alert("Concept and parent concept cannot be same")
      return
    }

    // Add parent node to child node.
    this.addNewParentToNode(parentConcept.id)
  }

  // Remove parent concept id from parent concept ids if error is sent from parent component.
  removeParentConceptIDIfError(parentConceptID: string): void{
    let tempParentConceptsIDs: string[] = this.node.parentConceptIDs
      let index: number = tempParentConceptsIDs.indexOf(parentConceptID)
      if (index !== -1) {
        tempParentConceptsIDs.splice(index, 1)
      }
      this.node.parentConceptIDs = [...tempParentConceptsIDs]
  }

  // On removing parent concept id from parent concept ids.
  onRemoveParentConceptID(parentConcept: any): void{
    let stepId: number = this.node.stepId
    let parentConceptID: string = parentConcept.value.id
    this.removeParentConceptEmitter.emit({stepId, parentConceptID})
  }

  // On clicking show parent concepts click.
  onShowParentConceptsClick(): void{
    this.node.showParentConcepts = !this.node.showParentConcepts
    this.redrawParentLines.emit()
  }

  // On parent concept change.
  onParentConceptChange(): void{
    this.redrawParentLines.emit()
  }
}
