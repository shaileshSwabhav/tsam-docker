import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CompanyMasterRequirementComponent } from './company-master-requirement.component';

describe('CompanyMasterRequirementComponent', () => {
  let component: CompanyMasterRequirementComponent;
  let fixture: ComponentFixture<CompanyMasterRequirementComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CompanyMasterRequirementComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CompanyMasterRequirementComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
