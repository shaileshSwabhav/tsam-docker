import { Injectable } from '@angular/core';
import { NgbModal, NgbModalOptions, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { ProgrammingQuestionModalComponent } from 'src/app/component/programming-question-modal/programming-question-modal.component';
import { UrlConstant } from '../constant';

@Injectable({
  providedIn: 'root'
})
export class AddModalComponentService {

  // Modal.
  modalRef: any

  // Menu.
  url: string

  constructor(
    private modalService: NgbModal,
    private urlConstant: UrlConstant,
  ) { }

  // Open modal of the component by its menu name.
  openModalByMenuName(url: string): void{
    this.url = url
    if (this.url == "/bank/coding-question"){
      this.openModal(ProgrammingQuestionModalComponent, 'xl')
    }else{
      console.log("nothing")
    }
    this.dismissModal()
  }

  // Used to open modal.
  openModal(content: any, size?: string, options: NgbModalOptions = {
    ariaLabelledBy: 'modal-basic-title', keyboard: false,
    backdrop: 'static', size: size
  }): NgbModalRef {
    if (!size) {
      options.size = 'lg'
    }
    this.modalRef = this.modalService.open(content, options)
    return this.modalRef
  }

  // Dismiss modal of the component.
  dismissModal() {
    let componentInstance: any
    if (this.url == this.urlConstant.BANK_PROGRAMMING_QUESTION){
      componentInstance = <ProgrammingQuestionModalComponent>this.modalRef.componentInstance
      componentInstance.dismissModalEvent.subscribe(() => {
        this.modalRef.dismiss()
      })
      componentInstance.addEvent.subscribe((id: string) => {
        this.modalRef.close()
        alert("Question added to bank")
      }, error => {
        console.error(error)
        alert("Question could not be added to bank")
      })
    }else{
      console.log("nothing")
    }
  }
}
