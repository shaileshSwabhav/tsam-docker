import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchSessionComponent } from './batch-session.component';

describe('BatchSessionComponent', () => {
  let component: BatchSessionComponent;
  let fixture: ComponentFixture<BatchSessionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchSessionComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchSessionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
