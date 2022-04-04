import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CollegeBranchComponent } from './college-branch.component';

describe('CollegeBranchComponent', () => {
  let component: CollegeBranchComponent;
  let fixture: ComponentFixture<CollegeBranchComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CollegeBranchComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CollegeBranchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
