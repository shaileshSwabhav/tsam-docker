import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminCareerObjectiveComponent } from './admin-career-objective.component';

describe('AdminCareerObjectiveComponent', () => {
  let component: AdminCareerObjectiveComponent;
  let fixture: ComponentFixture<AdminCareerObjectiveComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AdminCareerObjectiveComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AdminCareerObjectiveComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
