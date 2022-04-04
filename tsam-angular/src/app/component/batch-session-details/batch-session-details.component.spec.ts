import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchSessionDetailsComponent } from './batch-session-details.component';

describe('BatchSessionDetailsComponent', () => {
  let component: BatchSessionDetailsComponent;
  let fixture: ComponentFixture<BatchSessionDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchSessionDetailsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchSessionDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
