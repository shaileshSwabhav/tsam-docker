import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchCompletionDetailsComponent } from './batch-completion-details.component';

describe('BatchCompletionDetailsComponent', () => {
  let component: BatchCompletionDetailsComponent;
  let fixture: ComponentFixture<BatchCompletionDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchCompletionDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchCompletionDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
