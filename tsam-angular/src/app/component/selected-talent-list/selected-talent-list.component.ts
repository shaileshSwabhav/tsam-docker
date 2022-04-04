import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';

@Component({
      selector: 'app-selected-talent-list',
      templateUrl: './selected-talent-list.component.html',
      styleUrls: ['./selected-talent-list.component.css']
})
export class SelectedTalentListComponent implements OnInit {

      @Input() talentlist: any[];
      @Input() control: any;
      @Output() data = new EventEmitter<any>();
      private classarray: boolean[];
      constructor() {
            this.control = false;
      }

      ngOnInit() {
            this.classarray = Array(false, false, false);
            this.delayCall()
      }

      changeClass(index, status?: boolean) {
            if (this.classarray[index] == true) {
                  this.classarray[index] = false;
                  return;
            }
            if (status == true) {
                  return;
            }
            this.classarray[index] = true
      }

      getCssClass(index): boolean {
            return this.classarray[index];
      }

      delayCall() {
            setTimeout(() => {
                  console.log(this.control);
            }, 10);
      }

}

