import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CampusRegistrationFormComponent } from './campus-registration-form.component';

describe('CampusRegistrationFormComponent', () => {
  let component: CampusRegistrationFormComponent;
  let fixture: ComponentFixture<CampusRegistrationFormComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CampusRegistrationFormComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CampusRegistrationFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
