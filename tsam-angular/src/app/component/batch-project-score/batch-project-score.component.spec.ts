import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchProjectScoreComponent } from './batch-project-score.component';

describe('BatchProjectScoreComponent', () => {
  let component: BatchProjectScoreComponent;
  let fixture: ComponentFixture<BatchProjectScoreComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchProjectScoreComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchProjectScoreComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
