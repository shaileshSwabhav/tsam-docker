import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchTopicComponent } from './batch-topic.component';

describe('BatchTopicComponent', () => {
  let component: BatchTopicComponent;
  let fixture: ComponentFixture<BatchTopicComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchTopicComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchTopicComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
