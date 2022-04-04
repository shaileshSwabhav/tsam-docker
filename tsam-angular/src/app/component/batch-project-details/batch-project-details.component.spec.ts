import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchProjectDetailsComponent } from './batch-project-details.component';

describe('BatchProjectDetailsComponent', () => {
  let component: BatchProjectDetailsComponent;
  let fixture: ComponentFixture<BatchProjectDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchProjectDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchProjectDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
