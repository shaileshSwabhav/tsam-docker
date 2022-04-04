import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchProjectComponent } from './batch-project.component';

describe('BatchProjectComponent', () => {
  let component: BatchProjectComponent;
  let fixture: ComponentFixture<BatchProjectComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchProjectComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchProjectComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
