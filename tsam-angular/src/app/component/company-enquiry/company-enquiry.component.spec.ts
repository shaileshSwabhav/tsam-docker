import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CompanyEnquiryComponent } from './company-enquiry.component';

describe('CompanyEnquiryComponent', () => {
  let component: CompanyEnquiryComponent;
  let fixture: ComponentFixture<CompanyEnquiryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CompanyEnquiryComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CompanyEnquiryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
