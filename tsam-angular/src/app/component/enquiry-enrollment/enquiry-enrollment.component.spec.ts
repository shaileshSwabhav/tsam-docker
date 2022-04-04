import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EnquiryEnrollmentComponent } from './enquiry-enrollment.component';

describe('EnquiryEnrollmentComponent', () => {
  let component: EnquiryEnrollmentComponent;
  let fixture: ComponentFixture<EnquiryEnrollmentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EnquiryEnrollmentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EnquiryEnrollmentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
