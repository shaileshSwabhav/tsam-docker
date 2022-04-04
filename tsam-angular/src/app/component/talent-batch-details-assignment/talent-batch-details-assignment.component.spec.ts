import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TalentBatchDetailsAssignmentComponent } from './talent-batch-details-assignment.component';

describe('TalentBatchDetailsAssignmentComponent', () => {
  let component: TalentBatchDetailsAssignmentComponent;
  let fixture: ComponentFixture<TalentBatchDetailsAssignmentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TalentBatchDetailsAssignmentComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TalentBatchDetailsAssignmentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
