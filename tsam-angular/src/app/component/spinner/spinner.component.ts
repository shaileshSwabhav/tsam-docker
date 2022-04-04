import { Component, Input, OnChanges, OnInit } from '@angular/core';

@Component({
  selector: 'app-spinner',
  templateUrl: './spinner.component.html',
  styleUrls: ['./spinner.component.css']
})
export class SpinnerComponent implements OnInit, OnChanges {

  @Input() loadingMessage: string
  loaderTemplate: string

  constructor() { }

  ngOnInit(): void {
  }

  ngOnChanges(): void {
    if (!this.loadingMessage) {
      this.loadingMessage = "Loading..."
    }
    this.loaderTemplate = `<div class="loading"> <span>${this.loadingMessage}</span> </div>`
  }

}
