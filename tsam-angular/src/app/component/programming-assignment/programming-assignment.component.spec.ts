import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProgrammingAssignmentComponent } from './programming-assignment.component';

describe('ProgrammingAssignmentComponent', () => {
  let component: ProgrammingAssignmentComponent;
  let fixture: ComponentFixture<ProgrammingAssignmentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProgrammingAssignmentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProgrammingAssignmentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
