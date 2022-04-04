import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CourseProgrammingAssignmentComponent } from './course-programming-assignment.component';

describe('CourseProgrammingAssignmentComponent', () => {
  let component: CourseProgrammingAssignmentComponent;
  let fixture: ComponentFixture<CourseProgrammingAssignmentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CourseProgrammingAssignmentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CourseProgrammingAssignmentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
