import { ChangeDetectorRef, Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { NgSelectComponent } from '@ng-select/ng-select';
import { WorkflowItemScore } from 'src/app/service/concept-module/concept-module.service';

@Component({
  selector: 'app-module-concept-score-node',
  templateUrl: './module-concept-score-node.component.html',
  styleUrls: ['./module-concept-score-node.component.css']
})
export class ModuleConceptScoreNodeComponent {

  // Node.
  @Input() node: WorkflowItemScore

  // Redraw lines fo parent component.
  @Output() redrawParentLines: EventEmitter<any> = new EventEmitter<any>()

  constructor(
		private cdr: ChangeDetectorRef
  ) {
  }

  // On cliking show parent section.
  onShowParentSectionClick(): void{
    this.node.showParentSection = !this.node.showParentSection
  }

  // On clicking show parent concepts click.
  onShowParentConceptsClick(): void{
    this.node.showParentConcepts = !this.node.showParentConcepts
    this.redrawParentLines.emit()
  }

}
