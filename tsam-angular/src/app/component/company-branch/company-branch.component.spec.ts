import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CompanyBranchComponent } from './company-branch.component';

describe('CompanyBranchComponent', () => {
  let component: CompanyBranchComponent;
  let fixture: ComponentFixture<CompanyBranchComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CompanyBranchComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CompanyBranchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
