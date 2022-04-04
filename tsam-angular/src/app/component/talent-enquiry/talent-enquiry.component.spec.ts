import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TalentEnquiryComponent } from './talent-enquiry.component';

describe('TalentEnquiryComponent', () => {
  let component: TalentEnquiryComponent;
  let fixture: ComponentFixture<TalentEnquiryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TalentEnquiryComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TalentEnquiryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
